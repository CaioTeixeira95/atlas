package atlas

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(t *testing.T, size int) (string, error) {
	t.Helper()

	b := make([]rune, size)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		if err != nil {
			return "", fmt.Errorf("generating random integer: %w", err)
		}
		b[i] = letterRunes[n.Int64()]
	}

	return string(b), nil
}

func TestNewDefaultMap(t *testing.T) {
	assert.Panics(t, func() {
		NewDefaultMap[string, string](nil)
	})
	genFunc := func() (string, error) {
		return "default", nil
	}
	m := NewDefaultMap[string, string](genFunc)
	assert.NotNil(t, m)
}

func TestDefaultMapSet(t *testing.T) {
	genFunc := func() (string, error) {
		return "default", nil
	}
	m := NewDefaultMap[string, string](genFunc)
	assert.Empty(t, m.mp)

	m.Set("key", "value")
	assert.Contains(t, m.mp, "key")
	assert.Equal(t, "value", m.mp["key"])

	t.Run("ensures no race condition", func(t *testing.T) {
		var wg sync.WaitGroup

		wg.Add(2)
		go func() {
			defer wg.Done()
			m.Set("go1", "go1-value")
		}()

		go func() {
			defer wg.Done()
			_, err := m.Get("go2")
			require.NoError(t, err)
		}()

		wg.Wait()

		assert.Contains(t, m.mp, "go1")
		assert.Equal(t, "go1-value", m.mp["go1"])

		assert.Contains(t, m.mp, "go2")
		assert.Equal(t, "default", m.mp["go2"])
	})
}

func TestDefaultMapHas(t *testing.T) {
	genFunc := func() (int64, error) {
		return 99, nil
	}
	m := NewDefaultMap[int64, int64](genFunc)
	assert.Empty(t, m.mp)

	assert.False(t, m.Has(42))

	m.Set(42, 42)
	assert.True(t, m.Has(42))
}

func TestDefaultMapGet(t *testing.T) {
	genFunc := func() (string, error) {
		return "default", nil
	}
	m := NewDefaultMap[string, string](genFunc)
	assert.Empty(t, m.mp)

	t.Run("returns the key value successfully", func(t *testing.T) {
		m.Set("key", "value")
		value, err := m.Get("key")
		require.NoError(t, err)

		assert.Equal(t, "value", value)
	})

	t.Run("returns the default value successfully", func(t *testing.T) {
		value, err := m.Get("some-key")
		require.NoError(t, err)

		assert.Equal(t, "default", value)
	})

	t.Run("returns error when fails getting a default value from genFunc", func(t *testing.T) {
		gFunc := func() (string, error) {
			return "", errors.New("unexpected")
		}
		m := NewDefaultMap[string, string](gFunc)
		assert.Empty(t, m.mp)

		value, err := m.Get("fail-key")
		assert.EqualError(t, err, "getting value from genFunc: unexpected")
		assert.Empty(t, value)
	})
}

func TestDefaultMapDelete(t *testing.T) {
	m := NewDefaultMap[string, string](func() (string, error) { return randStringRunes(t, 6) })
	assert.Empty(t, m.mp)

	// Nothing happens
	m.Delete("key")

	_, err := m.Get("key")
	require.NoError(t, err)

	assert.Contains(t, m.mp, "key")
	assert.Len(t, m.mp["key"], 6)

	m.Delete("key")
	assert.Empty(t, m.mp)
}

func TestDefaultMapKeys(t *testing.T) {
	genFunc := func() (string, error) {
		return randStringRunes(t, 6)
	}
	m := NewDefaultMap[string, string](genFunc)
	assert.Empty(t, m.mp)

	m.Set("key1", "value1")
	m.Set("key2", "value2")
	m.Set("key3", "value3")

	_, err := m.Get("key4")
	require.NoError(t, err)

	_, err = m.Get("key5")
	require.NoError(t, err)

	assert.ElementsMatch(
		t,
		[]string{
			"key1",
			"key2",
			"key3",
			"key4",
			"key5",
		},
		m.Keys(),
	)
}

func TestDefaultMapValues(t *testing.T) {
	genFunc := func() (string, error) {
		return "default", nil
	}
	m := NewDefaultMap[string, string](genFunc)
	assert.Empty(t, m.mp)

	m.Set("key1", "value1")
	m.Set("key2", "value2")
	m.Set("key3", "value3")

	_, err := m.Get("key4")
	require.NoError(t, err)

	_, err = m.Get("key5")
	require.NoError(t, err)

	assert.ElementsMatch(
		t,
		[]string{
			"value1",
			"value2",
			"value3",
			"default",
			"default",
		},
		m.Values(),
	)
}

func TestDefaultMapSize(t *testing.T) {
	genFunc := func() (string, error) {
		return randStringRunes(t, 6)
	}
	m := NewDefaultMap[string, string](genFunc)
	assert.Empty(t, m.mp)

	assert.Equal(t, 0, m.Size())

	m.Set("key1", "value1")
	m.Set("key2", "value2")
	m.Set("key3", "value3")

	_, err := m.Get("key4")
	require.NoError(t, err)

	_, err = m.Get("key5")
	require.NoError(t, err)

	assert.Equal(t, 5, m.Size())
}

func TestDefaultMapToMap(t *testing.T) {
	genFunc := func() (string, error) {
		return randStringRunes(t, 6)
	}
	m := NewDefaultMap[string, string](genFunc)
	assert.Empty(t, m.mp)

	m.Set("key1", "value1")
	m.Set("key2", "value2")
	m.Set("key3", "value3")

	_, err := m.Get("key4")
	require.NoError(t, err)

	_, err = m.Get("key5")
	require.NoError(t, err)

	gotMap := m.ToMap()

	assert.Equal(t, m.mp, gotMap)
	assert.NotSame(t, m.mp, gotMap) // ensure is not the same pointer

	gotMap["key4"] = "value4"
	assert.NotEqual(t, m.mp, gotMap)
}
