package hcache

import (
	"sync"
)

type Keypair[K comparable, V any] struct {
	Key   K
	Value V
	prev  *Keypair[K, V]
	next  *Keypair[K, V]
}

type LRUCache[K comparable, V any] struct {
	cacheq sync.Map
	cap    uint64
	len    uint64
	head   *Keypair[K, V]
	tail   *Keypair[K, V]
}

func (lc *LRUCache[K, V]) Get(key K) (V, bool) {
	v, ok := lc.cacheq.Load(key)
	if !ok {
		var res V
		return res, false
	}

	n := v.(*Keypair[K, V])
	lc.toHead(n)
	return n.Value, true
}

func (lc *LRUCache[K, V]) Put(key K, value V) error {
	v, ok := lc.cacheq.Load(key)
	var n *Keypair[K, V]
	if !ok {
		n = &Keypair[K, V]{
			Key:   key,
			Value: value,
		}
		// 内存已满
		if lc.len == lc.cap {
			lc.deleteTail()
		}
		lc.cacheq.Store(key, n)
		lc.len++
	} else {
		n = v.(*Keypair[K, V])
		n.Value = value
	}

	lc.toHead(n)
	return nil
}

func (lc *LRUCache[K, V]) Cap() uint64 {
	return lc.cap
}

func (lc *LRUCache[K, V]) Len() uint64 {
	return lc.len
}

func (lc *LRUCache[K, V]) toHead(n *Keypair[K, V]) {
	if lc.head == nil {
		lc.head = n
		lc.tail = n
		return
	}

	// 非新node
	if n.prev != nil && n.next != nil {
		n.prev.next = n.next
		n.next.prev = n.prev
	}

	// 更新n的前后指针
	n.next = lc.head
	n.prev = nil

	lc.head.prev = n
	// 更新head
	lc.head = n
}

func (lc *LRUCache[K, V]) deleteTail() {
	lc.len--
	if lc.tail == nil {
		return
	}
	n := lc.tail

	if n.prev != nil {
		n.prev.next = nil
	}
	lc.tail = n.prev

	lc.cacheq.Delete(n.Key)
}

func newLRUCache[K comparable, V any](cap uint64) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		cap:    cap,
		cacheq: sync.Map{},
	}
}
