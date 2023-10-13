package hcache

import (
	"sync"
	"sync/atomic"
)

type LRUKKeypair[T any] struct {
	Key    string
	Value  T
	visits uint64
	prev   *LRUKKeypair[T]
	next   *LRUKKeypair[T]
}

type LRUKCache[T any] struct {
	cacheq    map[string]*LRUKKeypair[T]
	historyq  map[string]*LRUKKeypair[T]
	condition uint64
	cap       uint64
	len       atomic.Uint64
	head      *LRUKKeypair[T]
	headlock  sync.Mutex
	tail      *LRUKKeypair[T]
}

func (lc *LRUKCache[T]) Get(key string) T {
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

	var res T
	return res
}

func (lc *LRUKCache[T]) Put(key string, value T) error {
	n, ok := lc.cacheq[key]
	if ok {
		n.Value = value
	} else {
		n, ok = lc.historyq[key]
		if ok {
			n.Value = value
			n.visits++
		} else {
			n = &LRUKKeypair[T]{
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

func (lc *LRUKCache[T]) moveNodeFromHistoryqToCacheq(n *LRUKKeypair[T]) {
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

func (lc *LRUKCache[T]) toHeadNode(n *LRUKKeypair[T]) {
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

func (lc *LRUKCache[T]) deleteHistoryqNode(n *LRUKKeypair[T]) {
	delete(lc.historyq, n.Key)
}

func (lc *LRUKCache[T]) deleteTail() {
	n := lc.tail

	n.prev.next = nil
	lc.tail = n.prev

	delete(lc.cacheq, n.Key)
	lc.len.Store(lc.len.Load() - 1)
}

func newLRUKCache[T any](cacheqCap uint64, historyqCap uint64, condition uint64) *LRUKCache[T] {
	return &LRUKCache[T]{
		cap:       cacheqCap,
		condition: condition,
		cacheq:    make(map[string]*LRUKKeypair[T], cacheqCap),
		historyq:  make(map[string]*LRUKKeypair[T], historyqCap),
	}
}
