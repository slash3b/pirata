package cache

import (
	"container/list"
	"fmt"
	"sync"
)

type LRU[K comparable, V any] struct {
	l     *list.List
	cap   int
	m     map[K]*list.Element
	mutex sync.Mutex
}

type LRUEntry[K comparable, V any] struct {
	K K
	V V
}

func NewLRU[K comparable, V any](capacity int) *LRU[K, V] {
	return &LRU[K, V]{
		l:     list.New(),
		cap:   capacity,
		m:     make(map[K]*list.Element),
		mutex: sync.Mutex{},
	}
}

func (lr *LRU[K, V]) Get(key K) (V, error) {
	el, ok := lr.m[key]
	if !ok {
		// https://stackoverflow.com/questions/70585852/return-default-value-for-generic-type
		return *new(V), fmt.Errorf("no cache entry under key : %v", key)
	}

	lr.mutex.Lock()
	defer lr.mutex.Unlock()
	lr.l.MoveToFront(el)

	return el.Value.(LRUEntry[K, V]).V, nil
}

func (lr *LRU[K, V]) Set(k K, v V) {
	if el, ok := lr.m[k]; ok {
		lr.mutex.Lock()
		defer lr.mutex.Unlock()
		lr.l.MoveToFront(el)
		return
	}

	if lr.l.Len() == lr.cap {
		backEl := lr.l.Back()

		pair := backEl.Value.(LRUEntry[K, V])

		delete(lr.m, pair.K)
		lr.mutex.Lock()
		lr.l.Remove(backEl)
		lr.mutex.Unlock()
	}

	pair := LRUEntry[K, V]{
		K: k,
		V: v,
	}

	lr.mutex.Lock()
	lr.m[k] = lr.l.PushFront(pair)
	lr.mutex.Unlock()
}
