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

func (lc *LRUCache[K, V]) Get(key K) V {
	v, ok := lc.cacheq.Load(key)
	n := v.(*Keypair[K, V])
	if !ok {
		var res V
		return res
	}

	lc.toHeadNode(n)
	return n.Value
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

	lc.toHeadNode(n)
	return nil
}

func (lc *LRUCache[K, V]) toHeadNode(n *Keypair[K, V]) {
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

	lc.cacheq.Delete(n.Key)
	lc.len--
}

func newLRUCache[K comparable, V any](cap uint64) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		cap:    cap,
		cacheq: sync.Map{},
	}
}
