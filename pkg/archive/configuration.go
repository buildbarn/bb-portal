package archive

import (
	"os"

	pb "github.com/buildbarn/bb-portal/pkg/proto/configuration/archive"
	"github.com/buildbarn/bb-storage/pkg/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	folderPermission  = 0o750
)

// NewBlobArchiverFromConfiguration creates a BlobArchiver based on
// parameters provided in a configuration file.
func NewBlobArchiverFromConfiguration(configurations []*pb.BlobArchiverConfiguration) (BlobMultiArchiver, error) {
	if len(configurations) == 0 {
		return BlobMultiArchiver{}, status.Error(codes.InvalidArgument, "Archiver configuration not specified")
	}

	blobArchiver := NewBlobMultiArchiver()
	for _, configuration := range configurations {
		switch backend := configuration.Backend.(type) {
		case *pb.BlobArchiverConfiguration_Local:
			folder := backend.Local.Directory
			err := os.MkdirAll(folder, folderPermission)
			if err != nil {
				return BlobMultiArchiver{}, util.StatusWrapf(err, "failed to create blob archive folder %s", folder)
			}
			localFileArchiver := NewLocalFileArchiver(folder)
			blobArchiver.RegisterArchiver("file", localFileArchiver)
		}
	}

	return blobArchiver, nil
}
