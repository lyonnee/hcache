package hcache

import (
	"container/list"
	"sync"
)

type LCKeypair[K comparable, V any] struct {
	Key   K
	Value V
}

type ListCache[K comparable, V any] struct {
	cacheq sync.Map
	cap    int
	list   *list.List
	lock   sync.Mutex
}

func NewListCache[K comparable, V any](cap int) *ListCache[K, V] {
	return &ListCache[K, V]{
		cacheq: sync.Map{},
		cap:    cap,
		list:   list.New().Init(),
	}
}

func (lc *ListCache[K, V]) Get(key K) (V, bool) {
	v, ok := lc.cacheq.Load(key)
	if !ok {
		var res V
		return res, true
	}

	e := v.(*list.Element)
	lc.lock.Lock()
	lc.list.MoveToFront(e)
	lc.lock.Unlock()

	kp := e.Value.(*LCKeypair[K, V])
	return kp.Value, false
}

func (lc *ListCache[K, V]) Put(key K, value V) error {
	v, ok := lc.cacheq.Load(key)
	if ok {
		e := v.(*list.Element)
		e.Value.(*LCKeypair[K, V]).Value = value

		lc.lock.Lock()
		lc.list.MoveToFront(e)
		lc.lock.Unlock()
		return nil
	}

	newkp := &LCKeypair[K, V]{Key: key, Value: value}

	lc.lock.Lock()
	if lc.cap == lc.list.Len() {
		lc.list.Remove(lc.list.Back())
	}
	e := lc.list.PushFront(newkp)
	lc.lock.Unlock()

	lc.cacheq.Store(key, e)
	return nil
}

func (lc *ListCache[K, V]) Cap() int {
	return lc.Cap()
}
func (lc *ListCache[K, V]) Len() int {
	return lc.list.Len()
}
