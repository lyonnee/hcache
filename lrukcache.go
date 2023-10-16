package hcache

import (
	"sync"
	"sync/atomic"
)

type LRUKKeypair[K comparable, V any] struct {
	Key    K
	Value  V
	visits uint64
	prev   *LRUKKeypair[K, V]
	next   *LRUKKeypair[K, V]
}

type LRUKCache[K comparable, V any] struct {
	cacheq    map[K]*LRUKKeypair[K, V]
	historyq  map[K]*LRUKKeypair[K, V]
	condition uint64
	cap       uint64
	len       atomic.Uint64
	head      *LRUKKeypair[K, V]
	headlock  sync.Mutex
	tail      *LRUKKeypair[K, V]
}

func (lc *LRUKCache[K, V]) Get(key K) V {
	n, ok := lc.cacheq[key]
	if ok {
		lc.toHeadNode(n)
		return n.Value
	}

	n, ok = lc.historyq[key]
	if ok {
		n.visits++
		lc.moveNodeFromHistoryqToCacheq(n)
	}

	var res V
	return res
}

func (lc *LRUKCache[K, V]) Put(key K, value V) error {
	n, ok := lc.cacheq[key]
	if ok {
		n.Value = value
	} else {
		n, ok = lc.historyq[key]
		if ok {
			n.Value = value
			n.visits++
		} else {
			n = &LRUKKeypair[K, V]{
				Key:    key,
				Value:  value,
				visits: 1,
			}
			lc.historyq[key] = n
		}
		lc.moveNodeFromHistoryqToCacheq(n)
		return nil
	}

	lc.toHeadNode(n)
	return nil
}

func (lc *LRUKCache[K, V]) moveNodeFromHistoryqToCacheq(n *LRUKKeypair[K, V]) {
	if n.visits >= lc.condition {
		// 先从historyquene删除节点
		lc.deleteHistoryqNode(n)

		// 如果cachequene内存已满
		if lc.len.Load() == lc.cap {
			lc.deleteTail()
		}
		// 存放到cachequene中
		lc.cacheq[n.Key] = n
		lc.len.Add(1)

		lc.toHeadNode(n)
	}
}

func (lc *LRUKCache[K, V]) toHeadNode(n *LRUKKeypair[K, V]) {
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

func (lc *LRUKCache[K, V]) deleteHistoryqNode(n *LRUKKeypair[K, V]) {
	delete(lc.historyq, n.Key)
}

func (lc *LRUKCache[K, V]) deleteTail() {
	n := lc.tail

	n.prev.next = nil
	lc.tail = n.prev

	delete(lc.cacheq, n.Key)
	lc.len.Store(lc.len.Load() - 1)
}

func newLRUKCache[K comparable, V any](cacheqCap uint64, historyqCap uint64, condition uint64) *LRUKCache[K, V] {
	return &LRUKCache[K, V]{
		cap:       cacheqCap,
		condition: condition,
		cacheq:    make(map[K]*LRUKKeypair[K, V], cacheqCap),
		historyq:  make(map[K]*LRUKKeypair[K, V], historyqCap),
	}
}
