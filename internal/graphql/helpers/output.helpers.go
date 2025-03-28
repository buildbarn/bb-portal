package helpers

import (
	"context"
	"fmt"
	"log/slog"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/blob"
	"github.com/buildbarn/bb-portal/internal/graphql/model"
	"github.com/buildbarn/bb-portal/pkg/events"
)

// FileLookup A file lookup type.
type FileLookup func(ctx context.Context) (*bes.File, error)

// BlobReferenceForFile Blob Reference for File function.
func BlobReferenceForFile(ctx context.Context, db *ent.Client, fileLookup FileLookup) (*model.BlobReference, error) {
	file, err := fileLookup(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not find bes.File: %w", err)
	}
	blobRecord, err := findBlob(ctx, db, file)
	if err != nil {
		return nil, fmt.Errorf("could not find blob: %w", err)
	}
	if blobRecord == nil {
		return nil, nil
	}
	return &model.BlobReference{
		Name: file.GetName(),
		Blob: blobRecord,
	}, nil
}

// find blob function.
func findBlob(ctx context.Context, db *ent.Client, file *bes.File) (*ent.Blob, error) {
	uri := file.GetUri()
	if uri == "" {
		return nil, nil
	}

	return db.Blob.Query().Where(blob.URI(uri)).First(ctx)
}

// GetAction Get an Action.
func GetAction(ctx context.Context, problem *ent.BazelInvocationProblem) (*bes.ActionExecuted, error) {
	bepEvents, err := events.FromJSONArray(problem.BepEvents)
	if err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal action problem events", "problem", problem, "err", err)
		return nil, fmt.Errorf("failed to parse action problem: %w", err)
	}
	for _, event := range bepEvents {
		if event.IsActionCompleted() {
			return event.GetAction(), nil
		}
	}
	return nil, errActionNotFound
}
