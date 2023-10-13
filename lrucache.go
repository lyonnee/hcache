package hcache

import (
	"sync"
	"sync/atomic"
)

type Keypair[T any] struct {
	Key   string
	Value T
	prev  *Keypair[T]
	next  *Keypair[T]
}

type LRUCache[T any] struct {
	cacheq   map[string]*Keypair[T]
	cap      uint64
	len      atomic.Uint64
	head     *Keypair[T]
	headlock sync.Mutex
	tail     *Keypair[T]
}

func (lc *LRUCache[T]) Get(key string) T {
	n, ok := lc.cacheq[key]
	if !ok {
		var res T
		return res
	}

	lc.toHeadNode(n)
	return n.Value
}

func (lc *LRUCache[T]) Put(key string, value T) error {
	n, ok := lc.cacheq[key]
	if !ok {
		n = &Keypair[T]{
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

func (lc *LRUCache[T]) toHeadNode(n *Keypair[T]) {
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

func (lc *LRUCache[T]) deleteTail() {
	n := lc.tail

	n.prev.next = nil
	lc.tail = n.prev

	delete(lc.cacheq, n.Key)
	lc.len.Store(lc.len.Load() - 1)
}

func newLRUCache[T any](cap uint64) *LRUCache[T] {
	return &LRUCache[T]{
		cap:    cap,
		cacheq: make(map[string]*Keypair[T], cap),
	}
}
