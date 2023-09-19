// Package cellar
// @author tabuyos
// @since 2023/9/11
// @description cellar
package cellar

import (
	"sync"
)

type MemoryPersist[T any] struct {
	lks LockupService[EvalInfo[T]]
	mu  sync.Mutex
}

func (p *MemoryPersist[T]) Persistence(sqlKey string, info *EvalInfo[T]) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.lks.Put(sqlKey, info)
}

func (p *MemoryPersist[T]) Lookup(sqlKey string) *EvalInfo[T] {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.lks.Get(sqlKey)
}

func NewMemoryPersist[T any](service LockupService[EvalInfo[T]]) PersistService[T] {
	return &MemoryPersist[T]{lks: service}
}
