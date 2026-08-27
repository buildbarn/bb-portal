package buildeventrecorder

import (
	"context"
	"errors"
	"testing"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-storage/pkg/blobstore"
	storagedigest "github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/stretchr/testify/require"
)

const (
	availableFileURI = "bytestream://cache.example.com/blobs/4048aad102bbf0ee98cdfc2dc9797d9b01afab79ee77738134828ae637f5be07/42"
	missingFileURI   = "bytestream://cache.example.com/blobs/5048aad102bbf0ee98cdfc2dc9797d9b01afab79ee77738134828ae637f5be07/61"
)

type actionFileBlobAccess struct {
	blobstore.BlobAccess
	missingDigests storagedigest.Set
	err            error
	queriedDigests storagedigest.Set
}

func (ba *actionFileBlobAccess) FindMissing(_ context.Context, digests storagedigest.Set) (storagedigest.Set, error) {
	ba.queriedDigests = digests
	return ba.missingDigests, ba.err
}

func actionFile(uri string) *bes.File {
	return &bes.File{File: &bes.File_Uri{Uri: uri}}
}

func actionEvent(primaryOutputURI, stdoutURI string) BuildEventWithInfo {
	return BuildEventWithInfo{Event: &bes.BuildEvent{
		Payload: &bes.BuildEvent_Action{Action: &bes.ActionExecuted{
			PrimaryOutput: actionFile(primaryOutputURI),
			Stdout:        actionFile(stdoutURI),
		}},
	}}
}

func mustBytestreamDigest(t *testing.T, uri string) storagedigest.Digest {
	t.Helper()
	digest, isBytestreamURI := getBytestreamDigest(uri)
	require.True(t, isBytestreamURI)
	require.NotEqual(t, storagedigest.BadDigest, digest)
	return digest
}

func TestVerifyActionFileAvailability(t *testing.T) {
	missingDigest := mustBytestreamDigest(t, missingFileURI)
	blobAccess := &actionFileBlobAccess{missingDigests: missingDigest.ToSingletonSet()}
	recorder := &buildEventRecorder{contentAddressableStorage: blobAccess}

	availability := recorder.verifyActionFileAvailability(
		context.Background(),
		[]BuildEventWithInfo{actionEvent(availableFileURI, missingFileURI)},
	)

	require.Equal(t, 2, blobAccess.queriedDigests.Length())
	availableFile := availability.getAvailableFile(actionFile(availableFileURI), "bazel-out/bin/output")
	require.NotNil(t, availableFile)
	require.Equal(t, "bazel-out/bin/output", availableFile.Path)
	require.Nil(t, availability.getAvailableFile(actionFile(missingFileURI), "bazel-out/bin/missing"))
}

func TestActionFileAvailabilityDropsBytestreamLinksWhenVerificationFails(t *testing.T) {
	blobAccess := &actionFileBlobAccess{err: errors.New("verification failed")}
	recorder := &buildEventRecorder{contentAddressableStorage: blobAccess}

	availability := recorder.verifyActionFileAvailability(
		context.Background(),
		[]BuildEventWithInfo{actionEvent(availableFileURI, "")},
	)

	require.Nil(t, availability.getAvailableFile(actionFile(availableFileURI), "output"))
	require.Nil(t, availability.getAvailableFile(actionFile("https://example.com/output"), "output"))
}

func TestActionFileAvailabilityWithoutCASPreservesValidBytestreamFile(t *testing.T) {
	availability := (&buildEventRecorder{}).verifyActionFileAvailability(
		context.Background(),
		[]BuildEventWithInfo{actionEvent(availableFileURI, "")},
	)

	availableFile := availability.getAvailableFile(actionFile(availableFileURI), "output")
	require.NotNil(t, availableFile)
	require.Equal(t, "output", availableFile.Path)
	require.Nil(t, availability.getAvailableFile(actionFile("bytestream://cache.example.com/blobs/not-a-digest/42"), "output"))
	require.Nil(t, availability.getAvailableFile(actionFile("file:///tmp/output"), "output"))
	require.Nil(t, availability.getAvailableFile(actionFile("https://example.com/output"), "output"))
}
