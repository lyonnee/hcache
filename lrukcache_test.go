package hcache

import (
	"fmt"
	"testing"
)

func TestLRUCache1(t *testing.T) {
	cache := newLRUKCache[string, string](5, 5, 1)

	// 添加键值对
	cache.Put("key1", "value1")
	cache.Put("key2", "value2")
	cache.Put("key3", "value3")

	// 获取键值对
	value, exists := cache.Get("key1")
	if !exists || value != "value1" {
		t.Errorf("Expected key1 to exist with value 'value1'")
	}

	// 测试LRU淘汰
	cache.Put("key4", "value4")
	cache.Put("key5", "value5")
	cache.Put("key6", "value6") // 这将导致"key1"被淘汰

	value, exists = cache.Get("key1")
	if exists {
		t.Errorf("Expected key1 to be evicted from the cache")
	}

	// 更新值
	cache.Put("key3", "new_value3")
	value, _ = cache.Get("key3")
	if value != "new_value3" {
		t.Errorf("Expected key3 to be updated with 'new_value3'")
	}

	// 测试不存在的键
	_, exists = cache.Get("non_existent_key")
	if exists {
		t.Errorf("Expected non_existent_key not to exist in the cache")
	}
}

func TestConcurrentLRU1Cache(t *testing.T) {
	cache := newLRUKCache[string, string](100, 100, 1)

	numKeys := 100
	numReaders := 10
	numWriters := 5

	// 添加一些初始键值对
	for i := 0; i < numKeys; i++ {
		cache.Put(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	// 启动读取者
	for i := 0; i < numReaders; i++ {
		go func() {
			for j := 0; j < numKeys; j++ {
				key := fmt.Sprintf("key%d", j)
				value, exists := cache.Get(key)
				if !exists {
					t.Errorf("Reader: Expected key %s to exist in the cache", key)
				}
				expectedValue := fmt.Sprintf("value%d", j)
				if value != expectedValue {
					t.Errorf("Reader: Expected key %s to have value %s, but got %s", key, expectedValue, value)
				}
			}
		}()
	}

	// 启动写入者
	for i := 0; i < numWriters; i++ {
		go func() {
			for j := 0; j < numKeys; j++ {
				key := fmt.Sprintf("key%d", j)
				newValue := fmt.Sprintf("new_value%d", j)
				cache.Put(key, newValue)
				value, _ := cache.Get(key)
				if value != newValue {
					t.Errorf("Writer: Expected key %s to be updated with value %s", key, newValue)
				}
			}
		}()
	}
}

func TestLRU2QCache(t *testing.T) {
	cache := newLRUKCache[string, string](5, 5, 2)

	// 添加键值对
	cache.Put("key1", "value1")
	cache.Put("key2", "value2")
	cache.Put("key3", "value3")

	// 获取键值对
	value, exists := cache.Get("key1")
	if !exists || value != "value1" {
		t.Errorf("Expected key1 to exist with value 'value1'")
	}

	// 测试LRU淘汰
	cache.Put("key4", "value4")
	cache.Put("key5", "value5")
	cache.Put("key6", "value6") // 这将导致从historyq中的"key1"被淘汰

	// "key1" 仍然应该存在，因为它在2Q队列中
	value, exists = cache.Get("key1")
	if !exists || value != "value1" {
		t.Errorf("Expected key1 to exist with value 'value1'")
	}

	// 更新值
	cache.Put("key3", "new_value3")
	value, _ = cache.Get("key3")
	if value != "new_value3" {
		t.Errorf("Expected key3 to be updated with 'new_value3'")
	}

	// 测试不存在的键
	_, exists = cache.Get("non_existent_key")
	if exists {
		t.Errorf("Expected non_existent_key not to exist in the cache")
	}
}
func Benchmark2Q_Rand(b *testing.B) {
	l := newLRUKCache[int64, int64](1000, 1000, 2)

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
	l := newLRUKCache[int64, int64](1000, 1000, 2)

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
