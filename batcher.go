package batchpool

type Batcher[T any] struct {
	pool      *Pool[Batch[T]]
	batchSize int
	output    chan *Batch[T]
}

// NewBatcher creates a new batcher. The batchesBufSize parameter defines the size of the Output channel buffer.
func NewBatcher[T any](batchSize int, batchesBufSize int) *Batcher[T] {
	newBatch := func(p *Pool[Batch[T]]) *Batch[T] {
		return NewBatch[T](batchSize, func(b *Batch[T]) { p.Put(b) })
	}

	return &Batcher[T]{
		pool:      NewPool[Batch[T]](newBatch),
		batchSize: batchSize,
		output:    make(chan *Batch[T], batchesBufSize),
	}
}

// Range calls the callback for each element in the batch, automatically creating new batches when needed.
func (b *Batcher[T]) Range(cb func(v *T) (next bool, err error)) error {
	finished := false

	innerCallback := func(v *T) (next bool, err error) {
		next, err = cb(v)
		finished = !next

		return next, err
	}

	for !finished {
		batch := b.pool.Get()

		if err := batch.Range(innerCallback); err != nil {
			return err
		}

		b.output <- batch
	}

	close(b.output)

	return nil
}

// Batches returns a channel of batches. The channel is closed when the Range callback returns false.
func (b *Batcher[T]) Batches() <-chan *Batch[T] {
	return b.output
}
