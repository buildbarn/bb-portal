package buildeventrecorder

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/digest"
	"github.com/buildbarn/bb-portal/ent/gen/ent/file"
	"github.com/buildbarn/bb-portal/ent/gen/ent/filepath"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-portal/pkg/invocation/files"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// SaveSingleFile saves a single file to the database
func SaveSingleFile(ctx context.Context, tx database.Handle, instanceNameDbID int64, parsedFile files.ParsedBepFile) (int64, error) {
	filePathDbID, err := tx.Ent().FilePath.Create().
		SetBepInstanceNameID(instanceNameDbID).
		SetPath(parsedFile.Path).
		OnConflictColumns(filepath.BepInstanceNameColumn, filepath.FieldPath).
		Ignore().
		ID(ctx)
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to upsert file path")
	}

	digestDbID, err := tx.Ent().Digest.Create().
		SetRev2InstanceName(parsedFile.InstanceName).
		SetDigestFunction(parsedFile.DigestFunction).
		SetHash(parsedFile.Hash).
		SetSizeBytes(parsedFile.SizeBytes).
		OnConflictColumns(digest.FieldRev2InstanceName, digest.FieldDigestFunction, digest.FieldHash, digest.FieldSizeBytes).
		Ignore().
		ID(ctx)
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to upsert digest")
	}

	fileDbID, err := tx.Ent().File.Create().
		SetFilePathID(filePathDbID).
		SetDigestID(digestDbID).
		OnConflictColumns(file.FilePathColumn, file.DigestColumn).
		Ignore().
		ID(ctx)
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to upsert file")
	}
	return fileDbID, nil
}

func saveFilesBatch(ctx context.Context, tx database.Handle, instanceNameDbID int64, files []*files.ParsedBepFile) error {
	if len(files) == 0 {
		return nil
	}

	filePathsMap := make(map[string]struct{})

	type DigestKey struct {
		instanceName   string
		digestFunction int16
		hashStr        string
		sizeBytes      int64
	}
	digestsMap := make(map[DigestKey]struct{})

	type FileKey struct {
		filePath       string
		instanceName   string
		digestFunction int16
		hashStr        string
		sizeBytes      int64
	}
	filesMap := make(map[FileKey]struct{})

	for _, file := range files {
		hashStr := string(file.Hash)

		filePathsMap[file.Path] = struct{}{}
		digestsMap[DigestKey{instanceName: file.InstanceName, digestFunction: int16(file.DigestFunction), hashStr: hashStr, sizeBytes: file.SizeBytes}] = struct{}{}
		filesMap[FileKey{filePath: file.Path, instanceName: file.InstanceName, digestFunction: int16(file.DigestFunction), hashStr: hashStr, sizeBytes: file.SizeBytes}] = struct{}{}
	}

	// Upsert file paths
	filePathParams := sqlc.InsertMissingFilePathsParams{
		BepInstanceNameID: instanceNameDbID,
		Paths:             make([]string, 0, len(filePathsMap)),
	}
	for key := range filePathsMap {
		filePathParams.Paths = append(filePathParams.Paths, key)
	}
	if err := tx.Sqlc().InsertMissingFilePaths(ctx, filePathParams); err != nil {
		return util.StatusWrap(err, "Failed to upsert file paths in batch")
	}

	// Upsert digests
	digestsParams := sqlc.InsertMissingDigestsParams{
		Rev2InstanceNames: make([]string, 0, len(digestsMap)),
		DigestFunctions:   make([]int16, 0, len(digestsMap)),
		Hashes:            make([][]byte, 0, len(digestsMap)),
		SizeBytes:         make([]int64, 0, len(digestsMap)),
	}
	for key := range digestsMap {
		digestsParams.Rev2InstanceNames = append(digestsParams.Rev2InstanceNames, key.instanceName)
		digestsParams.DigestFunctions = append(digestsParams.DigestFunctions, key.digestFunction)
		digestsParams.Hashes = append(digestsParams.Hashes, []byte(key.hashStr))
		digestsParams.SizeBytes = append(digestsParams.SizeBytes, key.sizeBytes)
	}
	if err := tx.Sqlc().InsertMissingDigests(ctx, digestsParams); err != nil {
		return util.StatusWrap(err, "Failed to upsert digests in batch")
	}

	// Upsert files
	filesParams := sqlc.InsertMissingFilesParams{
		BepInstanceNameID: instanceNameDbID,
		FilePaths:         make([]string, 0, len(filesMap)),
		Rev2InstanceNames: make([]string, 0, len(filesMap)),
		DigestFunctions:   make([]int16, 0, len(filesMap)),
		Hashes:            make([][]byte, 0, len(filesMap)),
		SizeBytes:         make([]int64, 0, len(filesMap)),
	}
	for key := range filesMap {
		filesParams.FilePaths = append(filesParams.FilePaths, key.filePath)
		filesParams.Rev2InstanceNames = append(filesParams.Rev2InstanceNames, key.instanceName)
		filesParams.DigestFunctions = append(filesParams.DigestFunctions, key.digestFunction)
		filesParams.Hashes = append(filesParams.Hashes, []byte(key.hashStr))
		filesParams.SizeBytes = append(filesParams.SizeBytes, key.sizeBytes)
	}
	if err := tx.Sqlc().InsertMissingFiles(ctx, filesParams); err != nil {
		return util.StatusWrap(err, "Failed to upsert files in batch")
	}
	return nil
}
