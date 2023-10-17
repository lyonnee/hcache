package hcache

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRU1_Cache(t *testing.T) {
	cache := newLRUKCache[string, int](10, 1, 1)
	for i := 1; i <= 10; i++ {
		cache.Put(strconv.FormatInt(int64(i), 10), i)
	}
	assert.Equal(t, 10, cache.list.Front().Value.(*LRUKKeypair[string, int]).Value)
	assert.Equal(t, 1, cache.list.Back().Value.(*LRUKKeypair[string, int]).Value)

	cache.Get("6")
	assert.Equal(t, 6, cache.list.Front().Value.(*LRUKKeypair[string, int]).Value)

	cache.Put(strconv.FormatInt(11, 10), 11)
	cache.Put(strconv.FormatInt(12, 10), 12)
	assert.Equal(t, 12, cache.list.Front().Value.(*LRUKKeypair[string, int]).Value)
	assert.Equal(t, 3, cache.list.Back().Value.(*LRUKKeypair[string, int]).Value)

	cache.Put("6", 66)
	v, _ := cache.Get("6")
	assert.Equal(t, 66, v)
}

func TestLRU2_Cache(t *testing.T) {
	cache := newLRUKCache[string, int](10, 10, 2)
	for i := 1; i <= 10; i++ {
		cache.Put(strconv.FormatInt(int64(i), 10), i)
	}

	for i := 1; i <= 10; i++ {
		cache.Get(strconv.FormatInt(int64(i), 10))
	}
	assert.Equal(t, 10, cache.list.Front().Value.(*LRUKKeypair[string, int]).Value)
	assert.Equal(t, 1, cache.list.Back().Value.(*LRUKKeypair[string, int]).Value)

	cache.Get("6")
	assert.Equal(t, 6, cache.list.Front().Value.(*LRUKKeypair[string, int]).Value)

	cache.Put(strconv.FormatInt(11, 10), 11)
	cache.Put(strconv.FormatInt(12, 10), 12)
	assert.Equal(t, 6, cache.list.Front().Value.(*LRUKKeypair[string, int]).Value)
	assert.Equal(t, 1, cache.list.Back().Value.(*LRUKKeypair[string, int]).Value)
}

func Benchmark2Q_Rand(b *testing.B) {
	l := newLRUKCache[int64, int64](8192, 8192, 2)
	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		trace[i] = getRand(b) % 32768
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < 2*b.N; i++ {
		if i%2 == 0 {
			l.Put(trace[i], trace[i])
		} else {
			if _, ok := l.Get(trace[i]); ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(hit+miss))
}

func Benchmark2Q_Freq(b *testing.B) {
	l := newLRUKCache[int64, int64](8192, 8192, 2)

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		if i%2 == 0 {
			trace[i] = getRand(b) % 16384
		} else {
			trace[i] = getRand(b) % 32768
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Put(trace[i], trace[i])
	}
	var hit, miss int
	for i := 0; i < b.N; i++ {
		if _, ok := l.Get(trace[i]); ok {
			hit++
		} else {
			miss++
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(hit+miss))
}
