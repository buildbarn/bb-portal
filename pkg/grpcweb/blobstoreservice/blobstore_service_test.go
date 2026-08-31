package blobstoreservice

import (
	"context"
	"testing"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/blobstore"
	"github.com/buildbarn/bb-storage/pkg/blobstore/buffer"
	"github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type recordingBlobAccess struct {
	blobstore.BlobAccess
	findMissingDigests digest.Set
}

func (ba *recordingBlobAccess) FindMissing(_ context.Context, digests digest.Set) (digest.Set, error) {
	ba.findMissingDigests = digests
	return digest.EmptySet, nil
}

func TestRecorderContentAddressableStorage(t *testing.T) {
	ctx := context.Background()
	base := &recordingBlobAccess{}
	allowAuthorizer := auth.NewStaticAuthorizer(func(digest.InstanceName) bool { return true })
	denyAuthorizer := auth.NewStaticAuthorizer(func(digest.InstanceName) bool { return false })
	contentAddressableStorage := newRecorderContentAddressableStorage(base, allowAuthorizer, denyAuthorizer)

	actionDigest := digest.MustNewDigest(
		"",
		remoteexecution.DigestFunction_SHA256,
		"4048aad102bbf0ee98cdfc2dc9797d9b01afab79ee77738134828ae637f5be07",
		145,
	)
	actionDigests := actionDigest.ToSingletonSet()

	missing, err := contentAddressableStorage.FindMissing(ctx, actionDigests)
	require.NoError(t, err)
	require.True(t, missing.Empty())
	require.Equal(t, actionDigests, base.findMissingDigests)

	err = contentAddressableStorage.Put(
		ctx,
		actionDigest,
		buffer.NewValidatedBufferFromByteSlice([]byte("action")),
	)
	require.Equal(t, codes.PermissionDenied, status.Code(err))
}
