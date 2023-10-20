# HCache (Happy Cache)

A high-performance LRU cache library implemented in Go.

## Introduction

HCache is a Golang library that provides two LRU cache implementations - LRU and LRU-K algorithm. It aims to deliver higher performance and hit rates compared to other existing cache solutions.

Key features:

- Thread-safe LRU and LRU-K algorithm implementation
- Generic type support
- Superior benchmark performance and cache hit ratio. In the fixed access scenario, the cache hit ratio is very close to 100%
- Configurable cache size

## Benchmarks

HCache benchmark results:

```
cpu: Intel(R) Core(TM) i5-7400 CPU @ 3.00GHz
BenchmarkLRUCache_Rand-4   	 4523179	       264.1 ns/op	       4 B/op	       0 allocs/op
--- BENCH: BenchmarkLRUCache_Rand-4
    lrucache_test.go:52: hit: 0 miss: 1 ratio: 0.000000
    lrucache_test.go:52: hit: 1 miss: 99 ratio: 0.010000
    lrucache_test.go:52: hit: 1448 miss: 8552 ratio: 0.144800
    lrucache_test.go:52: hit: 961487 miss: 38513 ratio: 0.961487
    lrucache_test.go:52: hit: 4451955 miss: 71224 ratio: 0.984254
BenchmarkLRUCache_Freq-4   	 4518084	       233.9 ns/op	       2 B/op	       0 allocs/op
--- BENCH: BenchmarkLRUCache_Freq-4
    lrucache_test.go:88: hit: 1 miss: 0 ratio: 1.000000
    lrucache_test.go:88: hit: 100 miss: 0 ratio: 1.000000
    lrucache_test.go:88: hit: 9850 miss: 150 ratio: 0.985000
    lrucache_test.go:88: hit: 1000000 miss: 0 ratio: 1.000000
    lrucache_test.go:88: hit: 4518084 miss: 0 ratio: 1.000000
Benchmark2Q_Rand-4         	 4439372	       307.4 ns/op	       4 B/op	       0 allocs/op
--- BENCH: Benchmark2Q_Rand-4
    lrukcache_test.go:70: hit: 0 miss: 1 ratio: 0.000000
    lrukcache_test.go:70: hit: 0 miss: 100 ratio: 0.000000
    lrukcache_test.go:70: hit: 1290 miss: 8710 ratio: 0.129000
    lrukcache_test.go:70: hit: 957844 miss: 42156 ratio: 0.957844
    lrukcache_test.go:70: hit: 3436005 miss: 33289 ratio: 0.990405
    lrukcache_test.go:70: hit: 4376889 miss: 62483 ratio: 0.985925
Benchmark2Q_Freq-4         	 4161651	       249.7 ns/op	       2 B/op	       0 allocs/op
--- BENCH: Benchmark2Q_Freq-4
    lrukcache_test.go:98: hit: 1 miss: 0 ratio: 1.000000
    lrukcache_test.go:98: hit: 100 miss: 0 ratio: 1.000000
    lrukcache_test.go:98: hit: 7621 miss: 2379 ratio: 0.762100
    lrukcache_test.go:98: hit: 1000000 miss: 0 ratio: 1.000000
    lrukcache_test.go:98: hit: 4161651 miss: 0 ratio: 1.000000
PASS
```

golang-lru benchmark results:

```
cpu: Intel(R) Core(TM) i5-7400 CPU @ 3.00GHz
BenchmarkLRU_Rand-4   	 4440745	       280.3 ns/op	      60 B/op	       0 allocs/op
--- BENCH: BenchmarkLRU_Rand-4
    lru_test.go:36: hit: 0 miss: 1 ratio: 0.000000
    lru_test.go:36: hit: 0 miss: 100 ratio: 0.000000
    lru_test.go:36: hit: 1337 miss: 8663 ratio: 0.133700
    lru_test.go:36: hit: 248373 miss: 751627 ratio: 0.248373
    lru_test.go:36: hit: 1109774 miss: 3330971 ratio: 0.249907
BenchmarkLRU_Freq-4   	 4444165	       266.3 ns/op	      55 B/op	       0 allocs/op
--- BENCH: BenchmarkLRU_Freq-4
    lru_test.go:67: hit: 1 miss: 0 ratio: 1.000000
    lru_test.go:67: hit: 100 miss: 0 ratio: 1.000000
    lru_test.go:67: hit: 9859 miss: 141 ratio: 0.985900
    lru_test.go:67: hit: 311872 miss: 688128 ratio: 0.311872
    lru_test.go:67: hit: 1369485 miss: 3074680 ratio: 0.308154
Benchmark2Q_Rand-4    	 2108774	       576.6 ns/op	     103 B/op	       1 allocs/op
--- BENCH: Benchmark2Q_Rand-4
    2q_test.go:35: hit: 0 miss: 1 ratio: 0.000000
    2q_test.go:35: hit: 0 miss: 100 ratio: 0.000000
    2q_test.go:35: hit: 1445 miss: 8555 ratio: 0.144500
    2q_test.go:35: hit: 248338 miss: 751662 ratio: 0.248338
    2q_test.go:35: hit: 526742 miss: 1582032 ratio: 0.249786
Benchmark2Q_Freq-4    	 2475675	       478.1 ns/op	      91 B/op	       1 allocs/op
--- BENCH: Benchmark2Q_Freq-4
    2q_test.go:66: hit: 1 miss: 0 ratio: 1.000000
    2q_test.go:66: hit: 100 miss: 0 ratio: 1.000000
    2q_test.go:66: hit: 9827 miss: 173 ratio: 0.982700
    2q_test.go:66: hit: 332386 miss: 667614 ratio: 0.332386
    2q_test.go:66: hit: 810082 miss: 1665593 ratio: 0.327217
```

HCache performs better on throughput and allocations.

See test files for the full benchmark code.

## Installation
```shell
go get github.com/lyonnee/hcache@latest
```

## Usage & Example
### LRU(1) 
```go
import "github.com/lyonnee/hcache"

opts := &hcache.Options{
		CacheqCap: 10,
	}

cache := hcache.New[int, int](opts) 
cache.Set("key", "value")
value := cache.Get("key")
```

### LRU-K
```go
import "github.com/lyonnee/hcache"

opts := &hcache.Options{
    CacheqCap:   10,
    HistoryqCap: 10, // nullable, default equal to CacheqCap value
    Condition:   3,
}

cache := hcache.New[int, int](opts)
cache.Set("key", "value")
value := cache.Get("key")
```

## Contributing

Contributions are welcome! Open an issue or submit a pull request. 

## License

HCache is released under the [MIT License](LICENSE).
