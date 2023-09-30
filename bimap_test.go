package atlas

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBiMapSet(t *testing.T) {
	m := NewBiMap[string, string]()
	assert.Empty(t, m.mp)
	assert.Empty(t, m.inverse)

	t.Run("sets key successfully", func(t *testing.T) {
		err := m.Set("key", "value")
		require.NoError(t, err)

		assert.Contains(t, m.mp, "key")
		assert.Equal(t, "value", m.mp["key"])
		assert.Contains(t, m.inverse, "value")
		assert.Equal(t, "key", m.inverse["value"])

		err = m.Set("hi", "hi")
		require.NoError(t, err)

		assert.Contains(t, m.mp, "hi")
		assert.Equal(t, "hi", m.mp["hi"])
		assert.Contains(t, m.inverse, "hi")
		assert.Equal(t, "hi", m.inverse["hi"])
	})

	t.Run("returns error when key is duplicated", func(t *testing.T) {
		err := m.Set("key", "another")
		assert.EqualError(t, err, ErrDuplicatedKey.Error())
	})

	t.Run("returns error when value is duplicated", func(t *testing.T) {
		err := m.Set("another", "value")
		assert.EqualError(t, err, ErrDuplicatedValue.Error())
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
		assert.Contains(t, m.inverse, "go1-value")
		assert.Equal(t, "go1", m.inverse["go1-value"])

		assert.Contains(t, m.mp, "go2")
		assert.Equal(t, "go2-value", m.mp["go2"])
		assert.Contains(t, m.inverse, "go2-value")
		assert.Equal(t, "go2", m.inverse["go2-value"])
	})

	t.Run("works with different types", func(t *testing.T) {
		mInt := NewBiMap[int, int]()
		err := mInt.Set(1, 2)
		require.NoError(t, err)

		assert.Contains(t, mInt.mp, 1)
		assert.Equal(t, 2, mInt.mp[1])
		assert.Contains(t, mInt.inverse, 2)
		assert.Equal(t, 1, mInt.inverse[2])

		mStrInt := NewBiMap[string, int]()
		err = mStrInt.Set("key", 42)
		require.NoError(t, err)

		assert.Contains(t, mStrInt.mp, "key")
		assert.Equal(t, 42, mStrInt.mp["key"])
		assert.Contains(t, mStrInt.inverse, 42)
		assert.Equal(t, "key", mStrInt.inverse[42])
	})
}

func TestBiMapGet(t *testing.T) {
	m := NewBiMap[string, float32]()
	assert.Empty(t, m.mp)
	assert.Empty(t, m.inverse)

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

func TestBiMapGetInverse(t *testing.T) {
	m := NewBiMap[int64, string]()
	assert.Empty(t, m.mp)
	assert.Empty(t, m.inverse)

	t.Run("returns false when value doesn't exist", func(t *testing.T) {
		key, ok := m.GetInverse("some-value")
		assert.False(t, ok)
		assert.Empty(t, key)
	})

	t.Run("returns key successfully", func(t *testing.T) {
		err := m.Set(42, "value")
		require.NoError(t, err)

		key, ok := m.GetInverse("value")
		assert.True(t, ok)
		assert.Equal(t, int64(42), key)
	})
}

func TestBiMapHas(t *testing.T) {
	m := NewBiMap[int64, int64]()
	assert.Empty(t, m.mp)
	assert.Empty(t, m.inverse)

	assert.False(t, m.Has(42))

	err := m.Set(42, 42)
	require.NoError(t, err)

	assert.True(t, m.Has(42))
}

func TestBiMapHasInverse(t *testing.T) {
	m := NewBiMap[int64, int64]()
	assert.Empty(t, m.mp)
	assert.Empty(t, m.inverse)

	assert.False(t, m.HasInverse(42))

	err := m.Set(42, 42)
	require.NoError(t, err)

	assert.True(t, m.HasInverse(42))
}

func TestBiMapDelete(t *testing.T) {
	m := NewBiMap[string, string]()
	assert.Empty(t, m.mp)
	assert.Empty(t, m.inverse)

	// Nothing happens
	m.Delete("key")

	err := m.Set("key", "value")
	require.NoError(t, err)

	assert.Contains(t, m.mp, "key")
	assert.Equal(t, "value", m.mp["key"])
	assert.Contains(t, m.inverse, "value")
	assert.Equal(t, "key", m.inverse["value"])

	m.Delete("key")
	assert.Empty(t, m.mp)
	assert.Empty(t, m.inverse)
}

func TestBiMapKeys(t *testing.T) {
	m := NewBiMap[string, string]()
	assert.Empty(t, m.mp)
	assert.Empty(t, m.inverse)

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

func TestBiMapValues(t *testing.T) {
	m := NewBiMap[string, string]()
	assert.Empty(t, m.mp)
	assert.Empty(t, m.inverse)

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

func TestBiMapSize(t *testing.T) {
	m := NewBiMap[string, string]()
	assert.Empty(t, m.mp)
	assert.Empty(t, m.inverse)

	assert.Equal(t, 0, m.Size())

	err := m.Set("key1", "value1")
	require.NoError(t, err)

	err = m.Set("key2", "value2")
	require.NoError(t, err)

	err = m.Set("key3", "value3")
	require.NoError(t, err)

	assert.Equal(t, 3, m.Size())
}

func TestBiMapToMap(t *testing.T) {
	m := NewBiMap[string, string]()
	assert.Empty(t, m.mp)
	assert.Empty(t, m.inverse)

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

func TestBiMapToMapInverse(t *testing.T) {
	m := NewBiMap[string, string]()
	assert.Empty(t, m.mp)
	assert.Empty(t, m.inverse)

	err := m.Set("key1", "value1")
	require.NoError(t, err)

	err = m.Set("key2", "value2")
	require.NoError(t, err)

	err = m.Set("key3", "value3")
	require.NoError(t, err)

	gotMap := m.ToMapInverse()

	assert.Equal(t, m.inverse, gotMap)
	assert.NotSame(t, m.inverse, gotMap) // ensure is not the same pointer

	gotMap["value4"] = "key4"
	assert.NotEqual(t, m.inverse, gotMap)
}
