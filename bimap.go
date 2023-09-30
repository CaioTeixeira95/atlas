package atlas

import (
	"errors"
	"sync"

	"golang.org/x/exp/maps"
)

var (
	ErrDuplicatedKey   = errors.New("duplicated key")
	ErrDuplicatedValue = errors.New("duplicated value")
)

type biMap[K comparable, V comparable] struct {
	mu      sync.Mutex
	mp      map[K]V
	inverse map[V]K
}

// NewBiMap returns a new bidirectional map.
func NewBiMap[K comparable, V comparable]() *biMap[K, V] {
	return &biMap[K, V]{mp: make(map[K]V), inverse: make(map[V]K)}
}

// Set sets the value for the given key. Returns error if key or the value already exists.
func (m *biMap[K, V]) Set(key K, value V) error {
	if m.Has(key) {
		return ErrDuplicatedKey
	}

	if _, ok := m.GetInverse(value); ok {
		return ErrDuplicatedValue
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.mp[key] = value
	m.inverse[value] = key

	return nil
}

// Get returns the value of a given key, also returns if it's present or not.
func (m *biMap[K, V]) Get(key K) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, ok := m.mp[key]
	return value, ok
}

// GetInverse returns the key of a given value, also returns if it's present or not.
func (m *biMap[K, V]) GetInverse(value V) (K, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key, ok := m.inverse[value]
	return key, ok
}

// Has checks whether a key exists in the map.
func (m *biMap[K, V]) Has(key K) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.mp[key]
	return ok
}

// HasInverse checks whether a value exists in the map.
func (m *biMap[K, V]) HasInverse(value V) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.inverse[value]
	return ok
}

// Delete deletes the map key, also removes the inverse map key.
func (m *biMap[K, V]) Delete(key K) {
	if !m.Has(key) {
		return
	}

	value, _ := m.Get(key)

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.mp, key)
	delete(m.inverse, value)
}

// Keys returns the map keys.
func (m *biMap[K, V]) Keys() []K {
	return maps.Keys(m.mp)
}

// Values returns the map values.
func (m *biMap[K, V]) Values() []V {
	return maps.Values(m.mp)
}

// Size returns the map size.
func (m *biMap[K, V]) Size() int {
	return len(m.mp)
}

// ToMap returns a shallow clone of the original map.
func (m *biMap[K, V]) ToMap() map[K]V {
	return maps.Clone(m.mp)
}

// ToMapInverse returns a shallow clone of the original inverse map.
func (m *biMap[K, V]) ToMapInverse() map[V]K {
	return maps.Clone(m.inverse)
}
