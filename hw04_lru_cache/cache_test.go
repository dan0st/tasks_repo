package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		t.Run("capacity overflow for queue with equal elements using", func(t *testing.T) {
			c := NewCache(3)

			c.Set("A", 1)
			c.Set("B", 2)
			c.Set("C", 3)
			c.Set("D", 4)

			element, exists := c.Get("B")
			require.True(t, exists)
			require.Equal(t, 2, element)

			element, exists = c.Get("C")
			require.True(t, exists)
			require.Equal(t, 3, element)

			element, exists = c.Get("D")
			require.True(t, exists)
			require.Equal(t, 4, element)

			element, exists = c.Get("A")
			require.False(t, exists)
			require.Nil(t, element)
		})

		t.Run("capacity overflow for queue with not equal elements using", func(t *testing.T) {
			c := NewCache(3)

			c.Set("A", 1)
			c.Set("B", 2)
			c.Set("C", 3)

			c.Get("A")
			c.Get("B")
			c.Get("A")
			c.Set("D", 4)

			element, exists := c.Get("A")
			require.True(t, exists)
			require.Equal(t, 1, element)

			element, exists = c.Get("B")
			require.True(t, exists)
			require.Equal(t, 2, element)

			element, exists = c.Get("D")
			require.True(t, exists)
			require.Equal(t, 4, element)

			element, exists = c.Get("C")
			require.False(t, exists)
			require.Nil(t, element)
		})
	})

	t.Run("clear cache", func(t *testing.T) {
		c := NewCache(4)

		c.Clear()
		c.Set("A", 1)
		c.Set("B", 2)

		element, exists := c.Get("A")
		require.Equal(t, true, exists)
		require.Equal(t, 1, element)

		element, exists = c.Get("B")
		require.Equal(t, true, exists)
		require.Equal(t, 2, element)

		c.Clear()
		element, exists = c.Get("A")
		require.Equal(t, false, exists)
		require.Equal(t, nil, element)

		element, exists = c.Get("B")
		require.Equal(t, false, exists)
		require.Equal(t, nil, element)

		c.Set("A", 1)
		element, exists = c.Get("A")
		require.Equal(t, true, exists)
		require.Equal(t, 1, element)
	})
}

func TestCacheMultithreading(t *testing.T) {
	// t.Skip() // Remove me if task with asterisk completed.
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Clear()
		}
	}()

	wg.Wait()
}
