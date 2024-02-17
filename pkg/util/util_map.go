package util

import (
	"sync"
)

type keyType interface {
	uint64 | string
}

// Map is a goroutine-safe map
type Map[K keyType, V comparable] struct {
	mu      *sync.RWMutex
	entries map[K]V
}

func NewMap[K keyType, V comparable]() *Map[K, V] {
	return &Map[K, V]{
		mu:      &sync.RWMutex{},
		entries: map[K]V{},
	}
}

func (m *Map[K, V]) TryPut(k K, v V) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.entries[k]; ok {
		return false
	}

	m.entries[k] = v

	return true
}

func (m *Map[K, V]) Exists(k K) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.entries[k]

	return ok
}

func (m *Map[K, V]) Remove(k K) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.entries[k]; ok {
		delete(m.entries, k)
		return true
	}

	return false
}

func (m *Map[K, V]) RemoveIf(k K, callback func(v V) bool) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if v, ok := m.entries[k]; ok && callback(v) {
		delete(m.entries, k)
		return true
	}

	return false
}

func (m *Map[K, V]) RemoveByValue(v V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k, tv := range m.entries {
		if tv == v {
			delete(m.entries, k)
		}
	}
}

func (m *Map[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.entries)
}
