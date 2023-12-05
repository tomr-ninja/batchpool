package batchpool

type Batch[T any] struct {
	data       []*T
	afterReset func(batch *Batch[T])
}

// NewBatch creates a new batch of the given size. The afterReset callback is called when the batch is closed.
func NewBatch[T any](size int, afterReset func(batch *Batch[T])) *Batch[T] {
	data := make([]*T, size)
	for i := range data {
		data[i] = new(T)
	}

	return &Batch[T]{
		data:       data,
		afterReset: afterReset,
	}
}

// Len returns the number of elements in the batch. It is the same as the batch size,
// unless the Range callback returned false before the end of the batch.
func (b *Batch[T]) Len() int {
	return len(b.data)
}

// Data returns the underlying slice of batch elements
func (b *Batch[T]) Data() []*T {
	return b.data
}

// Range calls the callback for each element in the batch
func (b *Batch[T]) Range(cb func(v *T) (next bool, err error)) error {
	for i, v := range b.data {
		if next, err := cb(v); err != nil {
			b.data = b.data[:i]

			return err
		} else if !next {
			b.data = b.data[:i+1]

			return nil
		}
	}

	return nil
}

// Close resets batch values and calls the afterReset callback (usually to return batch to the pool)
func (b *Batch[T]) Close() {
	b.data = b.data[:cap(b.data)] // original size
	for _, v := range b.data {
		*v = *new(T) // reset value
	}

	b.afterReset(b)
}
