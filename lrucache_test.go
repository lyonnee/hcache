package hcache

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRUCache(t *testing.T) {
	cache := newLRUCache[int](10)
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
