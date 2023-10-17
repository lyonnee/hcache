package hcache

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListCache(t *testing.T) {
	cache := NewListCache[string, int](10)
	for i := 1; i <= 10; i++ {
		cache.Put(strconv.FormatInt(int64(i), 10), i)
	}
	assert.Equal(t, 10, cache.list.Front().Value.(*LCKeypair[string, int]).Value)
	assert.Equal(t, 1, cache.list.Back().Value.(*LCKeypair[string, int]).Value)

	data, _ := cache.Get("6")
	assert.Equal(t, data, cache.list.Front().Value.(*LCKeypair[string, int]).Value)

	cache.Put(strconv.FormatInt(11, 10), 11)
	cache.Put(strconv.FormatInt(12, 10), 12)
	assert.Equal(t, 12, cache.list.Front().Value.(*LCKeypair[string, int]).Value)
	assert.Equal(t, 3, cache.list.Back().Value.(*LCKeypair[string, int]).Value)
}

func BenchmarkListCache_Rand(b *testing.B) {
	l := NewListCache[int64, int64](8192)

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

func BenchmarkListCache_Freq(b *testing.B) {
	l := NewListCache[int64, int64](8192)

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
