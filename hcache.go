package hcache

type HCache[K comparable, V any] interface {
	Get(key K) (V, bool)
	Put(key K, value V) error
	Cap() int
	Len() int
}

type Options struct {
	CacheqCap   int
	HistoryqCap int
	Condition   int
}

func New[K comparable, V any](opts *Options) HCache[K, V] {
	if opts.Condition > 1 {
		if opts.HistoryqCap == 0 {
			opts.HistoryqCap = opts.CacheqCap
		}
		return newLRUKCache[K, V](opts.CacheqCap, opts.HistoryqCap, opts.Condition)
	}

	return newLRUCache[K, V](opts.CacheqCap)
}
