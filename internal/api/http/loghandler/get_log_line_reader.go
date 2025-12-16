package loghandler

import (
	"bytes"
	"context"
	"io"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/buildlogchunk"
	"github.com/buildbarn/bb-portal/internal/api/grpc/bes"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/google/uuid"
	"github.com/klauspost/compress/zstd"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *LogHandler) getLogLineReader(ctx context.Context, invocationID uuid.UUID, startLine, endLine int64) (io.Reader, error) {
	exists, err := h.client.BazelInvocation.Query().
		Where(bazelinvocation.InvocationID(invocationID)).
		Exist(ctx)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to query invocation existence")
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "Invocation %s not found", invocationID)
	}

	if startLine >= endLine {
		return bytes.NewReader(nil), nil
	}

	chunks, err := h.client.BuildLogChunk.Query().
		Where(
			buildlogchunk.HasBazelInvocationWith(bazelinvocation.InvocationID(invocationID)),
			buildlogchunk.FirstLineIndexLT(int64(endLine)),
			buildlogchunk.LastLineIndexGTE(int64(startLine)),
		).
		Order(ent.Asc(buildlogchunk.FieldChunkIndex)).
		Select(buildlogchunk.FieldID, buildlogchunk.FieldFirstLineIndex).
		All(ctx)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to fetch log chunks")
	}

	if len(chunks) == 0 {
		// For some reason there are no logs, check if they exist as non
		// normalized incompleted logs.
		logs, err := bes.GetNormalizedIncompleteBuildLogs(ctx, h.client, invocationID)
		if err != nil {
			return nil, util.StatusWrap(err, "Failed to get fallback construct")
		}
		return subsliceLogBufferByLines(logs, startLine, endLine), nil
	}

	reader, writer := io.Pipe()
	go func() {
		var retErr error
		defer func() {
			if retErr != nil {
				writer.CloseWithError(retErr)
			} else {
				writer.Close()
			}
		}()

		decoder, err := zstd.NewReader(nil)
		if err != nil {
			retErr = util.StatusWrap(err, "Failed to create zstd decoder")
		}

		decodedData := make([]byte, 0, 8*1024*1024)
		for _, meta := range chunks {
			chunk, err := h.client.BuildLogChunk.Get(ctx, meta.ID)
			if err != nil {
				retErr = util.StatusWrapf(err, "Failed to fetch chunk (id=%d)", meta.ID)
				return
			}
			decodedData, err = decoder.DecodeAll(chunk.Data, decodedData[:0])
			if err != nil {
				retErr = util.StatusWrapf(err, "Failed to decode chunk (id=%d)", meta.ID)
				return
			}

			line := chunk.FirstLineIndex
			prevByte := byte(0)
			var start, end int = 0, len(decodedData)
			for i, b := range decodedData {
				if prevByte == '\n' {
					line++
					if line == startLine {
						start = i
					}
					if line == endLine {
						end = i
						break
					}
				}
				prevByte = b
			}
			writer.Write(decodedData[start:end])
		}
	}()

	return reader, nil
}

func subsliceLogBufferByLines(buffer []byte, startLine, endLine int64) io.Reader {
	start, end := 0, len(buffer)
	if startLine != 0 {
		start = end
	}
	lineNum := int64(0)
	for i, b := range buffer {
		if b == '\n' {
			lineNum++
			if lineNum == startLine {
				start = i + 1
			}
			if lineNum == endLine {
				end = i + 1
				break
			}
		}
	}
	return bytes.NewReader(buffer[start:end])
}
