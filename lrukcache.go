package hcache

import (
	"sync"
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
	len       uint64
	head      *LRUKKeypair[K, V]
	locker    sync.Locker
	tail      *LRUKKeypair[K, V]
}

func (lc *LRUKCache[K, V]) Get(key K) (V, bool) {
	lc.locker.Lock()
	defer lc.locker.Unlock()

	v, ok := lc.cacheq[key]
	if ok {
		lc.toHead(v)
		return v.Value, true
	}

	v, ok = lc.historyq[key]
	if ok {
		v.visits++
		lc.moveNodeFromHistoryqToCacheq(v)
		return v.Value, true
	}

	var res V
	return res, false
}

func (lc *LRUKCache[K, V]) Put(key K, value V) error {
	lc.locker.Lock()
	defer lc.locker.Unlock()

	v, ok := lc.cacheq[key]
	if ok {
		v.Value = value
	} else {
		v, ok = lc.historyq[key]
		if ok {
			v.Value = value
			v.visits++
		} else {
			v = &LRUKKeypair[K, V]{
				Key:    key,
				Value:  value,
				visits: 1,
			}
			lc.historyq[key] = v
		}
		lc.moveNodeFromHistoryqToCacheq(v)
		return nil
	}

	lc.toHead(v)
	return nil
}

func (lc *LRUKCache[K, V]) moveNodeFromHistoryqToCacheq(n *LRUKKeypair[K, V]) {
	if n.visits >= lc.condition {
		// 先从historyquene删除节点
		lc.deleteHistoryqNode(n)

		// 如果cachequene内存已满
		if lc.len == lc.cap {
			lc.deleteTail()
		}
		// 存放到cachequene中
		lc.cacheq[n.Key] = n
		lc.len++

		lc.toHead(n)
	}
}

func (lc *LRUKCache[K, V]) Cap() uint64 {
	return lc.cap
}

func (lc *LRUKCache[K, V]) Len() uint64 {
	return lc.len
}

func (lc *LRUKCache[K, V]) toHead(n *LRUKKeypair[K, V]) {
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
	lc.len--
	if lc.tail == nil {
		return
	}
	n := lc.tail

	if n.prev != nil {
		n.prev.next = nil
	}
	lc.tail = n.prev

	delete(lc.cacheq, n.Key)
}

func newLRUKCache[K comparable, V any](cacheqCap uint64, historyqCap uint64, condition uint64) *LRUKCache[K, V] {
	return &LRUKCache[K, V]{
		cap:       cacheqCap,
		condition: condition,
		locker:    &sync.Mutex{},
		cacheq:    make(map[K]*LRUKKeypair[K, V], cacheqCap),
		historyq:  make(map[K]*LRUKKeypair[K, V], historyqCap),
	}
}
