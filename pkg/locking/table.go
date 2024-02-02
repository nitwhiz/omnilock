package locking

import "sync"

type keyType interface {
	uint64 | string
}

type Table[K keyType, V comparable] struct {
	mu      *sync.RWMutex
	entries map[K]V
}

func New[K keyType, V comparable]() *Table[K, V] {
	return &Table[K, V]{
		mu:      &sync.RWMutex{},
		entries: map[K]V{},
	}
}

func (t *Table[K, V]) TryPut(k K, v V) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.entries[k]; ok {
		return false
	}

	t.entries[k] = v

	return true
}

func (t *Table[K, V]) Exists(k K) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	_, ok := t.entries[k]

	return ok
}

func (t *Table[K, V]) Remove(k K) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.entries[k]; ok {
		delete(t.entries, k)

		return true
	}

	return false
}

func (t *Table[K, V]) RemoveIf(k K, callback func(v V) bool) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if v, ok := t.entries[k]; ok && callback(v) {
		delete(t.entries, k)

		return true
	}

	return false
}

func (t *Table[K, V]) RemoveByValue(v V) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for k, tv := range t.entries {
		if tv == v {
			delete(t.entries, k)
		}
	}
}

func (t *Table[K, V]) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return len(t.entries)
}
