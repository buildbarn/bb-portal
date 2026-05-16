package buildeventrecorder

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/klauspost/compress/zstd"
)

// decodeStream is a test helper: zstd-decompresses then peels off
// length-prefixed records, returning the inner payload bytes in order.
func decodeStream(t *testing.T, compressed []byte) [][]byte {
	t.Helper()
	if len(compressed) == 0 {
		return nil
	}
	dec, err := zstd.NewReader(nil)
	if err != nil {
		t.Fatalf("zstd reader: %v", err)
	}
	defer dec.Close()
	raw, err := dec.DecodeAll(compressed, nil)
	if err != nil {
		t.Fatalf("zstd decompress: %v", err)
	}
	var out [][]byte
	for off := 0; off < len(raw); {
		n, consumed := binary.Uvarint(raw[off:])
		if consumed <= 0 {
			t.Fatalf("invalid varint at offset %d", off)
		}
		off += consumed
		if off+int(n) > len(raw) {
			t.Fatalf("truncated record at offset %d (len=%d, remaining=%d)", off, n, len(raw)-off)
		}
		out = append(out, raw[off:off+int(n)])
		off += int(n)
	}
	return out
}

func TestArtifactGraphBuffer_EmptyFinalize(t *testing.T) {
	buf, err := newArtifactGraphBuffer()
	if err != nil {
		t.Fatalf("new buffer: %v", err)
	}
	payload, uncompressed, err := buf.Finalize()
	if err != nil {
		t.Fatalf("finalize: %v", err)
	}
	if uncompressed != 0 {
		t.Fatalf("expected uncompressed=0; got %d", uncompressed)
	}
	// A finalized empty zstd stream may still be a few bytes of frame headers;
	// the only contract is uncompressed=0. Just verify nothing terrible happened.
	if payload != nil && len(payload) > 32 {
		t.Fatalf("expected small/empty payload; got %d bytes", len(payload))
	}
}

func TestArtifactGraphBuffer_RoundTrip(t *testing.T) {
	buf, err := newArtifactGraphBuffer()
	if err != nil {
		t.Fatalf("new buffer: %v", err)
	}
	inputs := [][]byte{
		[]byte("first"),
		[]byte("a longer second record that has more content"),
		[]byte("third"),
	}
	for _, in := range inputs {
		if err := buf.AppendBuildEvent(in); err != nil {
			t.Fatalf("append: %v", err)
		}
	}
	payload, uncompressed, err := buf.Finalize()
	if err != nil {
		t.Fatalf("finalize: %v", err)
	}

	got := decodeStream(t, payload)
	if diff := cmp.Diff(inputs, got); diff != "" {
		t.Fatalf("round-tripped records differ (-want +got):\n%s", diff)
	}

	// Check uncompressed includes both varint headers and payload bytes.
	var wantUncompressed int64
	for _, in := range inputs {
		var v [binary.MaxVarintLen64]byte
		n := binary.PutUvarint(v[:], uint64(len(in)))
		wantUncompressed += int64(n + len(in))
	}
	if uncompressed != wantUncompressed {
		t.Fatalf("uncompressed accounting wrong: want %d got %d", wantUncompressed, uncompressed)
	}
}

func TestArtifactGraphBuffer_DoubleFinalize(t *testing.T) {
	buf, err := newArtifactGraphBuffer()
	if err != nil {
		t.Fatalf("new buffer: %v", err)
	}
	if _, _, err := buf.Finalize(); err != nil {
		t.Fatalf("first finalize: %v", err)
	}
	if _, _, err := buf.Finalize(); err == nil {
		t.Fatalf("expected error on second finalize, got nil")
	}
}

func TestArtifactGraphBuffer_AppendAfterFinalize(t *testing.T) {
	buf, err := newArtifactGraphBuffer()
	if err != nil {
		t.Fatalf("new buffer: %v", err)
	}
	if _, _, err := buf.Finalize(); err != nil {
		t.Fatalf("finalize: %v", err)
	}
	if err := buf.AppendBuildEvent([]byte("late")); err == nil {
		t.Fatalf("expected error on append-after-finalize, got nil")
	}
}

// TestArtifactGraphBuffer_MemoryCap exercises the cap path by appending a
// single record larger than the cap; subsequent appends should be dropped
// and finalize should still return a usable payload (containing nothing
// beyond what was accepted).
func TestArtifactGraphBuffer_MemoryCap(t *testing.T) {
	buf, err := newArtifactGraphBuffer()
	if err != nil {
		t.Fatalf("new buffer: %v", err)
	}

	// Manually set cap to a small value by hijacking via append size logic.
	// We can't override the const at runtime, so instead push an oversize
	// record directly: it should NOT be capped on its own because the cap
	// is per-cumulative bytes; once cumulative exceeds the cap, further
	// appends are dropped.
	//
	// Strategy: append the largest single record the buffer will accept
	// (just below cap), then attempt another append. The second should
	// silently drop and Capped() should report true.
	big := bytes.Repeat([]byte("x"), 1<<20) // 1 MiB
	// 256 MiB cap / 1 MiB = 256 accepted records before the cap.
	for i := 0; i < 257; i++ {
		if err := buf.AppendBuildEvent(big); err != nil {
			t.Fatalf("append #%d: %v", i, err)
		}
	}
	if !buf.Capped() {
		t.Fatalf("expected Capped() to be true after exceeding the cap")
	}
	payload, uncompressed, err := buf.Finalize()
	if err != nil {
		t.Fatalf("finalize: %v", err)
	}
	if uncompressed > artifactGraphMaxUncompressed {
		t.Fatalf("uncompressed (%d) should be <= cap (%d)", uncompressed, artifactGraphMaxUncompressed)
	}
	if len(payload) == 0 {
		t.Fatalf("expected non-empty payload")
	}
	// Decoded records: every one should be the 1 MiB string of 'x'.
	got := decodeStream(t, payload)
	for i, rec := range got {
		if !bytes.Equal(rec, big) {
			t.Fatalf("record %d does not match expected (len got=%d want=%d)", i, len(rec), len(big))
		}
	}
	// Should be strictly fewer than the 257 we tried to push.
	if len(got) >= 257 {
		t.Fatalf("expected cap to drop some records; got all %d back", len(got))
	}
}

// TestArtifactGraphBuffer_ConcurrentAppends sanity-checks that the mutex
// in the buffer serializes appends without corruption.
func TestArtifactGraphBuffer_ConcurrentAppends(t *testing.T) {
	buf, err := newArtifactGraphBuffer()
	if err != nil {
		t.Fatalf("new buffer: %v", err)
	}
	const n = 100
	done := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer func() { done <- struct{}{} }()
			rec := []byte(strings.Repeat("a", i+1))
			if err := buf.AppendBuildEvent(rec); err != nil {
				t.Errorf("append: %v", err)
			}
		}()
	}
	for i := 0; i < n; i++ {
		<-done
	}
	payload, _, err := buf.Finalize()
	if err != nil {
		t.Fatalf("finalize: %v", err)
	}
	got := decodeStream(t, payload)
	if len(got) != n {
		t.Fatalf("expected %d records, got %d", n, len(got))
	}
}
