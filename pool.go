package batchpool

import "sync"

type Pool[T any] struct {
	sync.Pool
}

func NewPool[T any](constructor func(p *Pool[T]) *T) *Pool[T] {
	pool := new(Pool[T])
	pool.New = func() interface{} {
		return constructor(pool)
	}

	return pool
}

func (p *Pool[T]) Get() *T {
	return p.Pool.Get().(*T)
}

func (p *Pool[T]) Put(v *T) {
	p.Pool.Put(v)
}
