package atlas

import (
	"errors"
	"sync"

	"golang.org/x/exp/maps"
)

var (
	ErrKeyAlreadySet = errors.New("key already set")
)

type frozenMap[K comparable, V any] struct {
	mu sync.Mutex
	mp map[K]V
}

func NewFrozenMap[K comparable, V any]() *frozenMap[K, V] {
	return &frozenMap[K, V]{mp: make(map[K]V)}
}

// Set sets the value for the given key.
func (m *frozenMap[K, V]) Set(key K, value V) error {
	if m.Has(key) {
		return ErrKeyAlreadySet
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.mp[key] = value

	return nil
}

// Has checks whether a key exists in the map.
func (m *frozenMap[K, V]) Has(key K) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.mp[key]
	return ok
}

// Get gets the value of the given key.
func (m *frozenMap[K, V]) Get(key K) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, ok := m.mp[key]
	return value, ok
}

// Delete deletes the map key.
func (m *frozenMap[K, V]) Delete(key K) {
	if !m.Has(key) {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.mp, key)
}

// Keys returns the map keys.
func (m *frozenMap[K, V]) Keys() []K {
	return maps.Keys(m.mp)
}

// Values returns the map values.
func (m *frozenMap[K, V]) Values() []V {
	return maps.Values(m.mp)
}

// Size returns the map size.
func (m *frozenMap[K, V]) Size() int {
	return len(m.mp)
}

// ToMap returns a shallow clone of the original map.
func (m *frozenMap[K, V]) ToMap() map[K]V {
	return maps.Clone(m.mp)
}
