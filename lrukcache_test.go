package hcache

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLRU1_Cache(t *testing.T) {
	cache := newLRUKCache[string, int](10, 1, 1)
	for i := 1; i <= 10; i++ {
		cache.Put(strconv.FormatInt(int64(i), 10), i)
	}
	assert.Equal(t, 10, cache.head.Value)
	assert.Equal(t, 1, cache.tail.Value)

	cache.Get("6")
	assert.Equal(t, 6, cache.head.Value)

	cache.Put(strconv.FormatInt(11, 10), 11)
	cache.Put(strconv.FormatInt(12, 10), 12)
	assert.Equal(t, 12, cache.head.Value)
	assert.Equal(t, 3, cache.tail.Value)
}

func TestLRU2_Cache(t *testing.T) {
	cache := newLRUKCache[string, int](10, 10, 2)
	for i := 1; i <= 10; i++ {
		cache.Put(strconv.FormatInt(int64(i), 10), i)
	}

	for i := 1; i <= 10; i++ {
		cache.Get(strconv.FormatInt(int64(i), 10))
	}
	assert.Equal(t, 10, cache.head.Value)
	assert.Equal(t, 1, cache.tail.Value)

	cache.Get("6")
	assert.Equal(t, 6, cache.head.Value)

	cache.Put(strconv.FormatInt(11, 10), 11)
	cache.Put(strconv.FormatInt(12, 10), 12)
	assert.Equal(t, 6, cache.head.Value)
	assert.Equal(t, 1, cache.tail.Value)
}

func BenchmarkLRUKCache(b *testing.B) {
	cacheqCapacity := 1000
	historyqCapacity := 1000
	condition := 10
	cache := newLRUKCache[uint64, string](cacheqCapacity, historyqCapacity, condition)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := uint64(time.Now().UnixNano())
			value := fmt.Sprintf("value-%d", key)
			cache.Put(key, value)
			cache.Get(key)
		}
	})
}
