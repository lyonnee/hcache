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
	cacheq       sync.Map
	historyq     map[K]*LRUKKeypair[K, V]
	condition    uint64
	cap          uint64
	len          uint64
	head         *LRUKKeypair[K, V]
	historyqLock sync.RWMutex
	tail         *LRUKKeypair[K, V]
}

func (lc *LRUKCache[K, V]) Get(key K) (V, bool) {
	v, ok := lc.cacheq.Load(key)
	var n *LRUKKeypair[K, V]
	if ok {
		n = v.(*LRUKKeypair[K, V])
		lc.toHead(n)
		return n.Value, true
	}

	lc.historyqLock.RLock()
	n, ok = lc.historyq[key]
	lc.historyqLock.RUnlock()
	if ok {
		n.visits++
		lc.moveNodeFromHistoryqToCacheq(n)
		return n.Value, true
	}

	var res V
	return res, false
}

func (lc *LRUKCache[K, V]) Put(key K, value V) error {
	v, ok := lc.cacheq.Load(key)
	var n *LRUKKeypair[K, V]
	if ok {
		n = v.(*LRUKKeypair[K, V])
		n.Value = value
	} else {
		lc.historyqLock.RLock()
		n, ok = lc.historyq[key]
		lc.historyqLock.RUnlock()
		if ok {
			n.Value = value
			n.visits++
		} else {
			n = &LRUKKeypair[K, V]{
				Key:    key,
				Value:  value,
				visits: 1,
			}
			lc.historyqLock.Lock()
			lc.historyq[key] = n
			lc.historyqLock.Unlock()
		}
		lc.moveNodeFromHistoryqToCacheq(n)
		return nil
	}

	lc.toHead(n)
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
		lc.cacheq.Store(n.Key, n)
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
	lc.historyqLock.Lock()
	defer lc.historyqLock.Unlock()
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

	lc.cacheq.Delete(n.Key)
}

func newLRUKCache[K comparable, V any](cacheqCap uint64, historyqCap uint64, condition uint64) *LRUKCache[K, V] {
	return &LRUKCache[K, V]{
		cap:       cacheqCap,
		condition: condition,
		cacheq:    sync.Map{},
		historyq:  make(map[K]*LRUKKeypair[K, V], historyqCap),
	}
}
