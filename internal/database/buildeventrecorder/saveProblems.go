package buildeventrecorder

import (
	"context"
	"log/slog"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/blob"
	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-portal/pkg/summary/detectors"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *BuildEventRecorder) determineMissingBlobs(ctx context.Context, tx *ent.Tx, detectedBlobs []detectors.BlobURI) ([]detectors.BlobURI, error) {
	detectedBlobURIs := make([]string, 0, len(detectedBlobs))
	blobMap := make(map[string]struct{}, len(detectedBlobs))
	for _, detectedBlob := range detectedBlobs {
		detectedBlobURIs = append(detectedBlobURIs, string(detectedBlob))
	}
	foundInDB, err := tx.Blob.Query().Where(blob.URIIn(detectedBlobURIs...)).All(ctx)
	if err != nil {
		return nil, util.StatusWrap(err, "failed to query blobs from database")
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

func (r *BuildEventRecorder) updateBlobRecord(ctx context.Context, tx *ent.Tx, b ent.Blob) {
	update := tx.Blob.Update().Where(blob.URI(b.URI)).SetArchivingStatus(b.ArchivingStatus)
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

func (r *BuildEventRecorder) saveBazelInvocationProblems(
	ctx context.Context,
	tx *ent.Tx,
	buildEvent *events.BuildEvent,
) error {
	if buildEvent == nil {
		return nil
	}

	problems, err := r.problemDetector.GetProblems(buildEvent)
	if err != nil {
		return util.StatusWrap(err, "Failed to detect problems")
	}

	if len(problems) == 0 {
		return nil
	}

	invocationDb, err := tx.BazelInvocation.Query().Where(bazelinvocation.ID(r.InvocationDbID)).Select(bazelinvocation.FieldInvocationID, bazelinvocation.FieldInstanceName).Only(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to find invocation for problems")
	}

	var detectedBlobs []detectors.BlobURI
	err = tx.BazelInvocationProblem.MapCreateBulk(problems, func(create *ent.BazelInvocationProblemCreate, i int) {
		problem := problems[i]
		detectedBlobs = append(detectedBlobs, problem.DetectedBlobs...)
		create.
			SetProblemType(string(problem.ProblemType)).
			SetLabel(problem.Label).
			SetBepEvents(problem.BEPEvents).
			SetBazelInvocationID(invocationDb.ID)
	}).Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save Bazel invocation problems")
	}
	missingBlobs, err := r.determineMissingBlobs(ctx, tx, detectedBlobs)
	if err != nil {
		return util.StatusWrap(err, "failed to determine missing blobs")
	}
	err = tx.Blob.MapCreateBulk(missingBlobs, func(create *ent.BlobCreate, i int) {
		b := missingBlobs[i]
		create.
			SetURI(string(b)).
			SetInstanceName(r.InstanceName)
	}).Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "failed to save blobs to database")
	}
	var archivedBlobs []ent.Blob
	archivedBlobs, err = r.blobArchiver.ArchiveBlobs(ctx, missingBlobs)
	if err != nil {
		return util.StatusWrap(err, "Archiving blobs failed")
	}
	for _, archivedBlob := range archivedBlobs {
		r.updateBlobRecord(ctx, tx, archivedBlob)
	}
	return nil
}
