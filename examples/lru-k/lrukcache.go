package lruk

import (
	"fmt"

	"github.com/lyonnee/hcache"
)

func main() {
	opts := &hcache.Options{
		CacheqCap:   10,
		HistoryqCap: 10, // nullable, default equal to CacheqCap value
		Condition:   3,
	}

	cache := hcache.New[int, int](opts)
	cache.Put(1, 1)
	cache.Put(2, 2)

	fmt.Println(cache.Get(2))
}
