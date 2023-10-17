package hcache

import (
	"container/list"
	"sync"
)

type Keypair[K comparable, V any] struct {
	Key   K
	Value V
}

type LRUCache[K comparable, V any] struct {
	cacheq sync.Map
	cap    int
	list   *list.List
	lock   sync.Mutex
}

func newLRUCache[K comparable, V any](cap int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		cacheq: sync.Map{},
		cap:    cap,
		list:   list.New().Init(),
	}
}

func (lc *LRUCache[K, V]) Get(key K) (V, bool) {
	v, ok := lc.cacheq.Load(key)
	if !ok {
		var res V
		return res, false
	}

	e := v.(*list.Element)
	lc.moveElemToHead(e)

	kp := e.Value.(*Keypair[K, V])
	return kp.Value, true
}

func (lc *LRUCache[K, V]) Put(key K, value V) error {
	v, ok := lc.cacheq.Load(key)
	if ok {
		e := v.(*list.Element)
		e.Value.(*Keypair[K, V]).Value = value

		lc.moveElemToHead(e)
		return nil
	}

	newkp := &Keypair[K, V]{Key: key, Value: value}
	lc.newKpToHead(newkp)
	return nil
}

func (lc *LRUCache[K, V]) Cap() int {
	return lc.cap
}

func (lc *LRUCache[K, V]) Len() int {
	return lc.list.Len()
}

func (lc *LRUCache[K, V]) newKpToHead(kp *Keypair[K, V]) {
	lc.lock.Lock()
	if lc.cap == lc.list.Len() {
		lc.list.Remove(lc.list.Back())
	}
	e := lc.list.PushFront(kp)
	lc.lock.Unlock()

	lc.cacheq.Store(kp.Key, e)
}

func (lc *LRUCache[K, V]) moveElemToHead(e *list.Element) {
	lc.lock.Lock()
	lc.list.MoveToFront(e)
	lc.lock.Unlock()
}
