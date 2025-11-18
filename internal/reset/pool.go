package reset

import (
	"sync"
)

// Интерфейс с методом Reset
type Resettable interface {
	Reset()
}

// Pool — generics контейнер с ограничением по Reset()
type Pool[T Resettable] struct {
	pool sync.Pool
}

// NewPool — Создаем Pool с поддержкой интерфейса Resettable
func NewPool[T Resettable](newFunc func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return newFunc()
			},
		},
	}
}

// Get — получить объект из пула (или создать новый)
func (p *Pool[T]) Get() T {
	v := p.pool.Get()
	if v == nil {
		var zero T
		return zero
	}
	return v.(T)
}

// Put — вернуть объект в пул (предварительно сбросить Reset)
func (p *Pool[T]) Put(obj T) {
	obj.Reset()
	p.pool.Put(obj)
}
