package servefiles

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"log"
	"net/http"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-portal/internal/api/common"
	"github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/buildbarn/bb-storage/pkg/filesystem/path"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s FileServerService) generateTarballDirectory(ctx context.Context, w *tar.Writer, digestFunction digest.Function, directory *remoteexecution.Directory, directoryPath *path.Trace, getDirectory func(context.Context, digest.Digest) (*remoteexecution.Directory, error), filesSeen map[string]string) error {
	// Emit child directories.
	for _, directoryNode := range directory.Directories {
		childName, ok := path.NewComponent(directoryNode.Name)
		if !ok {
			return status.Errorf(codes.InvalidArgument, "Directory %#v in directory %#v has an invalid name", directoryNode.Name, directoryPath.GetUNIXString())
		}
		childPath := directoryPath.Append(childName)

		if err := w.WriteHeader(&tar.Header{
			Typeflag: tar.TypeDir,
			Name:     childPath.GetUNIXString(),
			Mode:     0o777,
		}); err != nil {
			return err
		}
		childDigest, err := digestFunction.NewDigestFromProto(directoryNode.Digest)
		if err != nil {
			return err
		}
		childDirectory, err := getDirectory(ctx, childDigest)
		if err != nil {
			return err
		}
		if err := s.generateTarballDirectory(ctx, w, digestFunction, childDirectory, childPath, getDirectory, filesSeen); err != nil {
			return err
		}
	}

	// Emit symlinks.
	for _, symlinkNode := range directory.Symlinks {
		childName, ok := path.NewComponent(symlinkNode.Name)
		if !ok {
			return status.Errorf(codes.InvalidArgument, "Symbolic link %#v in directory %#v has an invalid name", symlinkNode.Name, directoryPath.GetUNIXString())
		}
		childPath := directoryPath.Append(childName)

		if err := w.WriteHeader(&tar.Header{
			Typeflag: tar.TypeSymlink,
			Name:     childPath.GetUNIXString(),
			Linkname: symlinkNode.Target,
			Mode:     0o777,
		}); err != nil {
			return err
		}
	}

	// Emit regular files.
	for _, fileNode := range directory.Files {
		childName, ok := path.NewComponent(fileNode.Name)
		if !ok {
			return status.Errorf(codes.InvalidArgument, "File %#v in directory %#v has an invalid name", fileNode.Name, directoryPath.GetUNIXString())
		}
		childPath := directoryPath.Append(childName)
		childPathString := childPath.GetUNIXString()

		childDigest, err := digestFunction.NewDigestFromProto(fileNode.Digest)
		if err != nil {
			return err
		}

		childKey := childDigest.GetKey(digest.KeyWithoutInstance)
		if fileNode.IsExecutable {
			childKey += "+x"
		} else {
			childKey += "-x"
		}

		if linkPath, ok := filesSeen[childKey]; ok {
			// This file was already returned previously.
			// Emit a hardlink pointing to the first
			// occurrence.
			//
			// Not only does this reduce the size of the
			// tarball, it also makes the directory more
			// representative of what it looks like when
			// executed through bb_worker.
			if err := w.WriteHeader(&tar.Header{
				Typeflag: tar.TypeLink,
				Name:     childPathString,
				Linkname: linkPath,
			}); err != nil {
				return err
			}
		} else {
			// This is the first time we're returning this
			// file. Actually add it to the archive.
			mode := int64(0o666)
			if fileNode.IsExecutable {
				mode = 0o777
			}
			if err := w.WriteHeader(&tar.Header{
				Typeflag: tar.TypeReg,
				Name:     childPathString,
				Size:     fileNode.Digest.SizeBytes,
				Mode:     mode,
			}); err != nil {
				return err
			}

			if err := s.contentAddressableStorage.Get(ctx, childDigest).IntoWriter(w); err != nil {
				return err
			}

			filesSeen[childKey] = childPathString
		}
	}
	return nil
}

func (s FileServerService) generateTarball(ctx context.Context, w http.ResponseWriter, digest digest.Digest, directory *remoteexecution.Directory, getDirectory func(context.Context, digest.Digest) (*remoteexecution.Directory, error)) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.tar.gz\"", digest.GetHashString()))
	w.Header().Set("Content-Type", "application/gzip")
	gzipWriter := gzip.NewWriter(w)
	tarWriter := tar.NewWriter(gzipWriter)
	filesSeen := map[string]string{}
	if err := s.generateTarballDirectory(ctx, tarWriter, digest.GetDigestFunction(), directory, nil, getDirectory, filesSeen); err != nil {
		// TODO(edsch): Any way to propagate this to the client?
		log.Print(err)
		panic(http.ErrAbortHandler)
	}
	if err := tarWriter.Close(); err != nil {
		log.Print(err)
		panic(http.ErrAbortHandler)
	}
	if err := gzipWriter.Close(); err != nil {
		log.Print(err)
		panic(http.ErrAbortHandler)
	}
}

// HandleDirectory serves a directory as a tarball.
func (s FileServerService) HandleDirectory(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if req.URL.Query().Get("format") != "tar" {
		http.Error(w, "Invalid format. Only supports \"tar\"", http.StatusNotFound)
		return
	}

	directoryDigest, err := getDigestFromRequest(req)
	if err != nil {
		http.Error(w, "Digest not found", http.StatusNotFound)
		return
	}

	ctx := common.ExtractContextFromRequest(req)
	directoryMessage, err := s.contentAddressableStorage.Get(ctx, directoryDigest).ToProto(&remoteexecution.Directory{}, s.maximumMessageSizeBytes)
	if err != nil {
		http.Error(w, "Digest not found", http.StatusNotFound)
		return
	}
	directory := directoryMessage.(*remoteexecution.Directory)

	s.generateTarball(ctx, w, directoryDigest, directory, func(ctx context.Context, digest digest.Digest) (*remoteexecution.Directory, error) {
		directoryMessage, err := s.contentAddressableStorage.Get(ctx, digest).ToProto(&remoteexecution.Directory{}, s.maximumMessageSizeBytes)
		if err != nil {
			return nil, err
		}
		return directoryMessage.(*remoteexecution.Directory), nil
	})
}
