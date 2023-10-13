package hcache

type HCache[T any] interface {
	Get(key string) T
	Put(key string, value T) error
}

type Options struct {
	CacheqCap   uint64
	HistoryqCap uint64
	Condition   uint64
}

func New[T any](opts *Options) HCache[T] {
	if opts.Condition > 1 {
		if opts.HistoryqCap == 0 {
			opts.HistoryqCap = opts.CacheqCap
		}
		return newLRUKCache[T](opts.CacheqCap, opts.HistoryqCap, opts.Condition)
	}

	return newLRUCache[T](opts.CacheqCap)
}
