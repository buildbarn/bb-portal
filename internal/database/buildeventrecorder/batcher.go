package buildeventrecorder

import (
	"sync"
)

// Batcher is a utility that atomically accumulates items to be
// processed in a batch.
type Batcher[T any] struct {
	mutex  sync.Mutex
	batch  []T
	signal chan struct{}
}

// NewBatcher returns a new batcher.
func NewBatcher[T any]() *Batcher[T] {
	return &Batcher[T]{
		signal: make(chan struct{}, 1),
	}
}

// Ready signals that there may be data available. Consumers must be
// able to handle there being no as it might have already been grabbed.
func (b *Batcher[T]) Ready() chan struct{} {
	return b.signal
}

// Add adds an item to the batch and signals that data is available to
// any consumer.
func (b *Batcher[T]) Add(item T) {
	b.mutex.Lock()
	b.batch = append(b.batch, item)
	b.mutex.Unlock()

	select {
	case b.signal <- struct{}{}:
	default:
	}
}

// Swap the underlying batch buffer.
func (b *Batcher[T]) Swap(buffer []T) []T {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	batch := b.batch
	b.batch = buffer
	return batch
}
