package hcache

import (
	"sync"
	"sync/atomic"
)

type Keypair[K comparable, V any] struct {
	Key   K
	Value V
	prev  *Keypair[K, V]
	next  *Keypair[K, V]
}

type LRUCache[K comparable, V any] struct {
	cacheq   map[K]*Keypair[K, V]
	cap      uint64
	len      atomic.Uint64
	head     *Keypair[K, V]
	headlock sync.Mutex
	tail     *Keypair[K, V]
}

func (lc *LRUCache[K, V]) Get(key K) V {
	n, ok := lc.cacheq[key]
	if !ok {
		var res V
		return res
	}

	lc.toHeadNode(n)
	return n.Value
}

func (lc *LRUCache[K, V]) Put(key K, value V) error {
	n, ok := lc.cacheq[key]
	if !ok {
		n = &Keypair[K, V]{
			Key:   key,
			Value: value,
		}
		// 内存已满
		if lc.len.Load() == lc.cap {
			lc.deleteTail()
		}
		lc.cacheq[key] = n
		lc.len.Add(1)
	} else {
		n.Value = value
	}

	lc.toHeadNode(n)
	return nil
}

func (lc *LRUCache[K, V]) toHeadNode(n *Keypair[K, V]) {
	lc.headlock.Lock()
	defer lc.headlock.Unlock()
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
	n := lc.tail

	n.prev.next = nil
	lc.tail = n.prev

	delete(lc.cacheq, n.Key)
	lc.len.Store(lc.len.Load() - 1)
}

func newLRUCache[K comparable, V any](cap uint64) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		cap:    cap,
		cacheq: make(map[K]*Keypair[K, V], cap),
	}
}
