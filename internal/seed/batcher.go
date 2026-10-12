package seed

import (
	"fmt"
	"log/slog"
	"time"
)

const batchSize = 4096

// batcher accumulates rows and flushes them to the database in chunks.
// Identifiers are assigned client side in a deterministic fashion, meaning
// that execute does not need to return any results.
type batcher[T any] struct {
	chunk     []T
	execute   func(chunk []T) error
	lastPrint time.Time
	processed int
	total     int
	name      string
}

func newBatcher[T any](name string, total int, execute func(chunk []T) error) *batcher[T] {
	return &batcher[T]{
		chunk:     make([]T, 0, batchSize),
		execute:   execute,
		lastPrint: time.Now(),
		processed: 0,
		total:     total,
		name:      name,
	}
}

func (b *batcher[T]) add(item T) error {
	b.chunk = append(b.chunk, item)
	if len(b.chunk) >= batchSize {
		return b.flush()
	}
	return nil
}

func (b *batcher[T]) flush() error {
	if len(b.chunk) > 0 {
		if err := b.execute(b.chunk); err != nil {
			return err
		}
		b.processed += len(b.chunk)
		b.chunk = make([]T, 0, batchSize)

		if time.Since(b.lastPrint) > 1*time.Second {
			slog.Info(fmt.Sprintf("%s: processed %d/%d", b.name, b.processed, b.total))
			b.lastPrint = time.Now()
		}
	}
	return nil
}

func (b *batcher[T]) result() error {
	return b.flush()
}
