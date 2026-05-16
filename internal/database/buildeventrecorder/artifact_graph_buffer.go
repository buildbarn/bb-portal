package buildeventrecorder

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sync"

	"github.com/klauspost/compress/zstd"
)

// artifactGraphMaxUncompressed bounds the uncompressed bytes accepted into
// the buffer. Beyond this, further appends are silently dropped so a
// pathological build doesn't pin gigabytes per concurrent invocation.
// BuildFinished still finalizes whatever was buffered.
const artifactGraphMaxUncompressed = 256 << 20 // 256 MiB

var errArtifactGraphBufferClosed = errors.New("artifact graph buffer already finalized")

// artifactGraphBuffer is a per-invocation streaming zstd buffer that
// concatenates length-prefixed serialized BEP messages — specifically
// the BuildEvent variants for NamedSetOfFiles and TargetCompleted that
// together describe the build's per-file artifact graph.
//
// Wire format inside the compressed stream:
//
//	repeat: uvarint(len) | <len bytes of serialized bes.BuildEvent>
//
// A single buffer is owned by one buildEventRecorder; events for that
// invocation are serialized through the recorder's transactional batch
// handlers. The mutex guards against accidental concurrent appends and
// the late-arriving finalize/append race.
type artifactGraphBuffer struct {
	mu                sync.Mutex
	enc               *zstd.Encoder
	out               *bytes.Buffer
	uncompressedTotal int64
	capped            bool
	closed            bool
}

func newArtifactGraphBuffer() (*artifactGraphBuffer, error) {
	out := &bytes.Buffer{}
	enc, err := zstd.NewWriter(out, zstd.WithEncoderLevel(zstd.SpeedDefault))
	if err != nil {
		return nil, err
	}
	return &artifactGraphBuffer{enc: enc, out: out}, nil
}

// AppendBuildEvent writes one length-prefixed serialized BuildEvent.
// Silently drops further entries once the uncompressed cap is reached.
func (b *artifactGraphBuffer) AppendBuildEvent(payload []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return errArtifactGraphBufferClosed
	}
	if b.capped {
		return nil
	}
	var lenBuf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(lenBuf[:], uint64(len(payload)))
	addedUncompressed := int64(n + len(payload))
	if b.uncompressedTotal+addedUncompressed > artifactGraphMaxUncompressed {
		b.capped = true
		return nil
	}
	if _, err := b.enc.Write(lenBuf[:n]); err != nil {
		return err
	}
	if _, err := b.enc.Write(payload); err != nil {
		return err
	}
	b.uncompressedTotal += addedUncompressed
	return nil
}

// Finalize closes the encoder and returns the compressed payload and the
// total uncompressed byte count. Idempotent error on second call.
func (b *artifactGraphBuffer) Finalize() ([]byte, int64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return nil, 0, errArtifactGraphBufferClosed
	}
	b.closed = true
	if err := b.enc.Close(); err != nil {
		return nil, 0, err
	}
	return b.out.Bytes(), b.uncompressedTotal, nil
}

// Capped reports whether the uncompressed-size cap was reached. Test
// hook; recorders may surface this via logging at finalize time.
func (b *artifactGraphBuffer) Capped() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.capped
}
