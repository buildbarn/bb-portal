package dbcleanupservice_test

import (
	"context"
	"testing"

	"github.com/buildbarn/bb-portal/internal/database/buildeventrecorder"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/pkg/invocation/files"
	"github.com/buildbarn/bb-portal/test/testutils"
	"github.com/buildbarn/bb-storage/pkg/clock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/mock/gomock"
)

func TestRemoveFiles(t *testing.T) {
	_, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoFiles", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()

		cleanup, err := getNewDbCleanupService(db, clock.SystemClock, traceProvider)
		require.NoError(t, err)
		deleted, err := cleanup.RemoveUnusedFiles(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 0, deleted)
		deleted, err = cleanup.RemoveUnusedFilePaths(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 0, deleted)
		deleted, err = cleanup.RemoveUnusedDigests(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 0, deleted)

		count, err := client.File.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
		count, err = client.FilePath.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
		count, err = client.Digest.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("UsedAndUnusedFiles", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()

		instanceName := testutils.CreateInstanceName(ctx, t, client, "testInstance")

		f1, err := buildeventrecorder.SaveSingleFile(ctx, db, instanceName.ID, files.ParsedBepFile{Path: "1", InstanceName: "1", DigestFunction: 1, Hash: []byte{1}, SizeBytes: 1})
		require.NoError(t, err)
		f2, err := buildeventrecorder.SaveSingleFile(ctx, db, instanceName.ID, files.ParsedBepFile{Path: "2", InstanceName: "2", DigestFunction: 2, Hash: []byte{2}, SizeBytes: 2})
		require.NoError(t, err)
		f3, err := buildeventrecorder.SaveSingleFile(ctx, db, instanceName.ID, files.ParsedBepFile{Path: "3", InstanceName: "3", DigestFunction: 3, Hash: []byte{3}, SizeBytes: 3})
		require.NoError(t, err)
		_, err = buildeventrecorder.SaveSingleFile(ctx, db, instanceName.ID, files.ParsedBepFile{Path: "4", InstanceName: "4", DigestFunction: 4, Hash: []byte{4}, SizeBytes: 4})
		require.NoError(t, err)

		inv := testutils.StartCreateInvocation(client, instanceName).AddBuildToolLogIDs(f1).SaveX(ctx)
		conf := client.Configuration.Create().SetConfigurationID("1").SetBazelInvocation(inv).SaveX(ctx)
		client.Action.Create().SetBazelInvocation(inv).SetConfiguration(conf).SetLabel("foo").SetStdoutID(f2).SetStderrID(f3).SaveX(ctx)

		cleanup, err := getNewDbCleanupService(db, clock.SystemClock, traceProvider)
		require.NoError(t, err)
		deleted, err := cleanup.RemoveUnusedFilePaths(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 0, deleted)
		deleted, err = cleanup.RemoveUnusedDigests(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 0, deleted)
		deleted, err = cleanup.RemoveUnusedFiles(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 1, deleted)
		deleted, err = cleanup.RemoveUnusedFilePaths(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 1, deleted)
		deleted, err = cleanup.RemoveUnusedDigests(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 1, deleted)

		count, err := client.File.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 3, count)
		count, err = client.FilePath.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 3, count)
		count, err = client.Digest.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 3, count)
	})
}
