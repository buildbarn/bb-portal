package processing

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/blob"
	"github.com/buildbarn/bb-portal/ent/gen/ent/build"
	"github.com/buildbarn/bb-portal/pkg/summary"
	"github.com/buildbarn/bb-portal/pkg/summary/detectors"
)

type SaveActor struct {
	db           *ent.Client
	blobArchiver BlobMultiArchiver
}

func (act SaveActor) SaveSummary(ctx context.Context, summary *summary.Summary) (*ent.BazelInvocation, error) {
	eventFile, err := act.saveEventFile(ctx, summary)
	if err != nil {
		return nil, fmt.Errorf("could not save EventFile: %w", err)
	}

	buildRecord, err := act.findOrCreateBuild(ctx, summary)
	if err != nil {
		return nil, err
	}

	bazelInvocation, err := act.saveBazelInvocation(ctx, summary, eventFile, buildRecord)
	if err != nil {
		return nil, fmt.Errorf("could not save BazelInvocation: %w", err)
	}

	var detectedBlobs []detectors.BlobURI

	err = act.db.BazelInvocationProblem.MapCreateBulk(summary.Problems, func(create *ent.BazelInvocationProblemCreate, i int) {
		problem := summary.Problems[i]
		detectedBlobs = append(detectedBlobs, problem.DetectedBlobs...)
		create.
			SetProblemType(string(problem.ProblemType)).
			SetLabel(problem.Label).
			SetBepEvents(problem.BEPEvents).
			SetBazelInvocation(bazelInvocation)
	}).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not save BazelInvocationProblems: %w", err)
	}

	missingBlobs, err := act.determineMissingBlobs(ctx, detectedBlobs)
	if err != nil {
		return nil, err
	}

	err = act.db.Blob.MapCreateBulk(missingBlobs, func(create *ent.BlobCreate, i int) {
		b := missingBlobs[i]
		create.SetURI(string(b))
		// Leave defaults for other fields, all updated during archiving if it is enabled:
		// 	- size_bytes: 0
		// 	- archiving_status: QUEUED
		// 	- reason: null
		// 	- archive_url: null
	}).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not save Blobs: %w", err)
	}

	var archivedBlobs []ent.Blob
	archivedBlobs, err = act.blobArchiver.ArchiveBlobs(ctx, missingBlobs)
	if err != nil {
		return nil, fmt.Errorf("failed to archive blobs: %w", err)
	}
	for _, archivedBlob := range archivedBlobs {
		act.updateBlobRecord(ctx, archivedBlob)
	}

	return bazelInvocation, nil
}

func (act SaveActor) determineMissingBlobs(ctx context.Context, detectedBlobs []detectors.BlobURI) ([]detectors.BlobURI, error) {
	detectedBlobURIs := make([]string, 0, len(detectedBlobs))
	blobMap := make(map[string]struct{}, len(detectedBlobs))
	for _, detectedBlob := range detectedBlobs {
		detectedBlobURIs = append(detectedBlobURIs, string(detectedBlob))
	}
	foundInDB, err := act.db.Blob.Query().Where(blob.URIIn(detectedBlobURIs...)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not query Blobs: %w", err)
	}

	for _, foundBlob := range foundInDB {
		blobMap[foundBlob.URI] = struct{}{}
	}
	missingBlobs := make([]detectors.BlobURI, 0, len(detectedBlobs)-len(foundInDB))
	for _, detectedBlob := range detectedBlobs {
		if _, ok := blobMap[string(detectedBlob)]; ok {
			continue
		}
		missingBlobs = append(missingBlobs, detectedBlob)
	}
	return missingBlobs, nil
}

func (act SaveActor) saveBazelInvocation(ctx context.Context, summary *summary.Summary, eventFile *ent.EventFile, buildRecord *ent.Build) (*ent.BazelInvocation, error) {
	create := act.db.BazelInvocation.Create().
		SetInvocationID(uuid.MustParse(summary.InvocationID)).
		SetStartedAt(summary.StartedAt).
		SetNillableEndedAt(summary.EndedAt).
		SetChangeNumber(int32(summary.ChangeNumber)).
		SetPatchsetNumber(int32(summary.PatchsetNumber)).
		SetSummary(*summary.InvocationSummary).
		SetBepCompleted(summary.BEPCompleted).
		SetStepLabel(summary.StepLabel).
		SetRelatedFiles(summary.RelatedFiles).
		SetEventFile(eventFile)

	if buildRecord != nil {
		create = create.SetBuild(buildRecord)
	}

	return create.
		Save(ctx)
}

func (act SaveActor) saveEventFile(ctx context.Context, summary *summary.Summary) (*ent.EventFile, error) {
	eventFile, err := act.db.EventFile.Create().
		SetURL(summary.EventFileURL).
		SetModTime(time.Now()).              // TODO: Save modTime in summary?
		SetProtocol("BEP").                  // Legacy: used to detect other protocols, e.g. for codechecks.
		SetMimeType("application/x-ndjson"). // NOTE: Only ndjson supported right now, but we should be able to add binary support.
		SetStatus("SUCCESS").                // TODO: Keep workflow of DETECTED->IMPORTING->...?
		Save(ctx)
	return eventFile, err
}

func (act SaveActor) findOrCreateBuild(ctx context.Context, summary *summary.Summary) (*ent.Build, error) {
	var err error
	var buildRecord *ent.Build

	if summary.BuildURL == "" {
		return nil, nil
	}

	slog.Info("Querying for build", "url", summary.BuildURL, "uuid", summary.BuildUUID)
	buildRecord, err = act.db.Build.Query().
		Where(build.BuildUUID(summary.BuildUUID)).First(ctx)

	if ent.IsNotFound(err) {
		slog.Info("Creating build", "url", summary.BuildURL, "uuid", summary.BuildUUID)
		buildRecord, err = act.db.Build.Create().
			SetBuildURL(summary.BuildURL).
			SetBuildUUID(summary.BuildUUID).
			SetEnv(buildEnvVars(summary.EnvVars)).
			Save(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("could not find or create build: %w", err)
	}
	return buildRecord, nil
}

func (act SaveActor) updateBlobRecord(ctx context.Context, b ent.Blob) {
	update := act.db.Blob.Update().Where(blob.URI(b.URI)).SetArchivingStatus(b.ArchivingStatus)
	if b.ArchiveURL != "" {
		update = update.SetArchiveURL(b.ArchiveURL)
	}
	if b.Reason != "" {
		update = update.SetReason(b.Reason)
	}
	if b.SizeBytes != 0 {
		update = update.SetSizeBytes(b.SizeBytes)
	}
	if _, err := update.Save(ctx); err != nil {
		slog.Error("failed to save archived blob", "uri", b.URI, "err", err)
	}
}

// buildEnvVars filters the input so it only contains well known environment
// variables injected into a CI build (e.g. a Jenkins build). These are well-known
// Jenkins, etc. environment variables and/or environment variables associated
// with plugins for GitHub, Gerrit, etc.
func buildEnvVars(env map[string]string) map[string]string {
	buildEnv := make(map[string]string)
	for k, v := range env {
		if !summary.IsBuildEnvKey(k) {
			continue
		}
		buildEnv[k] = v
	}

	return buildEnv
}
