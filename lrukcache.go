package hcache

import (
	"container/list"
	"sync"
)

type LRUKKeypair[K comparable, V any] struct {
	Key    K
	Value  V
	visits int
}

type LRUKCache[K comparable, V any] struct {
	cacheq       sync.Map
	historyq     map[K]*LRUKKeypair[K, V]
	condition    int
	cap          int
	historyqLock sync.RWMutex
	list         *list.List
	listLock     sync.Mutex
}

func (lc *LRUKCache[K, V]) Get(key K) (V, bool) {
	v, ok := lc.cacheq.Load(key)
	if ok {
		e := v.(*list.Element)
		lc.moveElemToHead(e)
		n := e.Value.(*LRUKKeypair[K, V])
		return n.Value, true
	}

	lc.historyqLock.RLock()
	n, ok := lc.historyq[key]
	lc.historyqLock.RUnlock()
	if ok {
		n.visits++
		lc.newKpToHead(n)
		return n.Value, true
	}

	var res V
	return res, false
}

func (lc *LRUKCache[K, V]) Put(key K, value V) error {
	v, ok := lc.cacheq.Load(key)
	if ok {
		e := v.(*list.Element)
		e.Value.(*LRUKKeypair[K, V]).Value = value
		lc.moveElemToHead(e)
		return nil
	}

	lc.historyqLock.RLock()
	n, ok := lc.historyq[key]
	lc.historyqLock.RUnlock()
	if ok {
		n.Value = value
		n.visits++
		lc.newKpToHead(n)
		return nil
	}

	newKp := &LRUKKeypair[K, V]{
		Key:    key,
		Value:  value,
		visits: 1,
	}
	lc.historyqLock.Lock()
	lc.historyq[key] = newKp
	lc.historyqLock.Unlock()
	lc.newKpToHead(newKp)
	return nil
}

func (lc *LRUKCache[K, V]) Cap() int {
	return lc.cap
}

func (lc *LRUKCache[K, V]) Len() int {
	return lc.list.Len()
}

func (lc *LRUKCache[K, V]) newKpToHead(n *LRUKKeypair[K, V]) {
	if n.visits < lc.condition {
		return
	}
	lc.listLock.Lock()
	if lc.cap == lc.list.Len() {
		lc.list.Remove(lc.list.Back())
	}
	e := lc.list.PushFront(n)
	lc.listLock.Unlock()
	lc.cacheq.Store(n.Key, e)
}

func (lc *LRUKCache[K, V]) moveElemToHead(e *list.Element) {
	lc.listLock.Lock()
	lc.list.MoveToFront(e)
	lc.listLock.Unlock()
}

func newLRUKCache[K comparable, V any](cacheqCap int, historyqCap int, condition int) *LRUKCache[K, V] {
	return &LRUKCache[K, V]{
		cap:       cacheqCap,
		condition: condition,
		cacheq:    sync.Map{},
		historyq:  make(map[K]*LRUKKeypair[K, V], historyqCap),
		list:      list.New().Init(),
	}
}
