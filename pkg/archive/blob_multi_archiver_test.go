package archive_test

import (
	"context"
	"testing"

	"github.com/buildbarn/bb-portal/pkg/archive"
	"github.com/buildbarn/bb-portal/pkg/summary/detectors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"

	"go.uber.org/mock/gomock"
)

func TestArchiveBlobs(t *testing.T) {
	_, ctx := gomock.WithContext(context.Background(), t)

	archiver := archive.NewBlobMultiArchiver()

	t.Run("NoURIs", func(t *testing.T) {
		_, err := archiver.ArchiveBlobs(ctx, nil)
		require.NoError(t, err)
	})

	t.Run("NoArchiverForScheme", func(t *testing.T) {
		_, err := archiver.ArchiveBlobs(ctx, []detectors.BlobURI{"file://test"})
		require.Error(t, err)
		require.Equal(t, status.Error(codes.Internal, "Failed to clean temporary directory"), err)
	})
}
