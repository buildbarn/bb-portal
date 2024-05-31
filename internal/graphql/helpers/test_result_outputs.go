package helpers

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/graphql/model"
	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-portal/third_party/bazel/gen/bes"
)

func GetTestResultActionLogOutput(ctx context.Context, client *ent.Client, obj *model.TestResult) (*model.BlobReference, error) {
	return getTestResultActionOutputByName(ctx, client, obj, "test.log")
}

func GetTestResultUndeclaredTestOutputs(ctx context.Context, client *ent.Client, obj *model.TestResult) (*model.BlobReference, error) {
	return getTestResultActionOutputByName(ctx, client, obj, events.UndeclaredTestOutputsName)
}

func getTestResultActionOutputByName(ctx context.Context, client *ent.Client, obj *model.TestResult, name string) (*model.BlobReference, error) {
	fileLookup := func(_ context.Context) (*bes.File, error) {
		var file *bes.File
		for _, output := range obj.BESTestResult.GetTestActionOutput() {
			if output.GetName() == name {
				file = output
				break
			}
		}

		if file == nil {
			return nil, nil
		}
		return file, nil
	}

	return BlobReferenceForFile(ctx, client, fileLookup)
}
