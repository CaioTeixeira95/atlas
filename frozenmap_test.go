package atlas

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrozenMapSet(t *testing.T) {
	m := NewFrozenMap[string, string]()
	assert.Empty(t, m.mp)

	t.Run("sets key successfully", func(t *testing.T) {
		err := m.Set("key", "value")
		require.NoError(t, err)

		assert.Contains(t, m.mp, "key")
		assert.Equal(t, "value", m.mp["key"])

		err = m.Set("hi", "hi")
		require.NoError(t, err)

		assert.Contains(t, m.mp, "hi")
		assert.Equal(t, "hi", m.mp["hi"])
	})

	t.Run("returns error when key is already set", func(t *testing.T) {
		err := m.Set("key", "another")
		assert.EqualError(t, err, ErrKeyAlreadySet.Error())
	})

	t.Run("ensures no race condition", func(t *testing.T) {
		var wg sync.WaitGroup

		wg.Add(2)
		go func() {
			defer wg.Done()
			err := m.Set("go1", "go1-value")
			require.NoError(t, err)
		}()

		go func() {
			defer wg.Done()
			err := m.Set("go2", "go2-value")
			require.NoError(t, err)
		}()

		wg.Wait()

		assert.Contains(t, m.mp, "go1")
		assert.Equal(t, "go1-value", m.mp["go1"])

		assert.Contains(t, m.mp, "go2")
		assert.Equal(t, "go2-value", m.mp["go2"])
	})
}

func TestFrozenMapGet(t *testing.T) {
	m := NewFrozenMap[string, float32]()
	assert.Empty(t, m.mp)

	t.Run("returns false when key doesn't exist", func(t *testing.T) {
		value, ok := m.Get("some-key")
		assert.False(t, ok)
		assert.Empty(t, value)
	})

	t.Run("returns value successfully", func(t *testing.T) {
		err := m.Set("key", 3.14)
		require.NoError(t, err)

		value, ok := m.Get("key")
		assert.True(t, ok)
		assert.Equal(t, float32(3.14), value)
	})
}

func TestFrozenMapHas(t *testing.T) {
	m := NewFrozenMap[int64, int64]()
	assert.Empty(t, m.mp)

	assert.False(t, m.Has(42))

	err := m.Set(42, 42)
	require.NoError(t, err)

	assert.True(t, m.Has(42))
}

func TestFrozenMapDelete(t *testing.T) {
	m := NewFrozenMap[string, string]()
	assert.Empty(t, m.mp)

	// Nothing happens
	m.Delete("key")

	err := m.Set("key", "value")
	require.NoError(t, err)

	assert.Contains(t, m.mp, "key")
	assert.Equal(t, "value", m.mp["key"])

	m.Delete("key")
	assert.Empty(t, m.mp)
}

func TestFrozenMapKeys(t *testing.T) {
	m := NewFrozenMap[string, string]()
	assert.Empty(t, m.mp)

	err := m.Set("key1", "value1")
	require.NoError(t, err)

	err = m.Set("key2", "value2")
	require.NoError(t, err)

	err = m.Set("key3", "value3")
	require.NoError(t, err)

	assert.ElementsMatch(
		t,
		[]string{
			"key1",
			"key2",
			"key3",
		},
		m.Keys(),
	)
}

func TestFrozenMapValues(t *testing.T) {
	m := NewFrozenMap[string, string]()
	assert.Empty(t, m.mp)

	err := m.Set("key1", "value1")
	require.NoError(t, err)

	err = m.Set("key2", "value2")
	require.NoError(t, err)

	err = m.Set("key3", "value3")
	require.NoError(t, err)

	assert.ElementsMatch(
		t,
		[]string{
			"value1",
			"value2",
			"value3",
		},
		m.Values(),
	)
}

func TestFrozenMapSize(t *testing.T) {
	m := NewFrozenMap[string, string]()
	assert.Empty(t, m.mp)

	assert.Equal(t, 0, m.Size())

	err := m.Set("key1", "value1")
	require.NoError(t, err)

	err = m.Set("key2", "value2")
	require.NoError(t, err)

	err = m.Set("key3", "value3")
	require.NoError(t, err)

	assert.Equal(t, 3, m.Size())
}

func TestFrozenMapToMap(t *testing.T) {
	m := NewFrozenMap[string, string]()
	assert.Empty(t, m.mp)

	err := m.Set("key1", "value1")
	require.NoError(t, err)

	err = m.Set("key2", "value2")
	require.NoError(t, err)

	err = m.Set("key3", "value3")
	require.NoError(t, err)

	gotMap := m.ToMap()

	assert.Equal(t, m.mp, gotMap)
	assert.NotSame(t, m.mp, gotMap) // ensure is not the same pointer

	gotMap["key4"] = "value4"
	assert.NotEqual(t, m.mp, gotMap)
}
