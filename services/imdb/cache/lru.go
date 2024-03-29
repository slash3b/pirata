package cache

import (
	"container/list"
	"fmt"
	"sync"

	"imdb/metrics"
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
		metrics.HitMissCache.Dec()
		metrics.CacheEvent.WithLabelValues("cache_miss").Inc()
		// https://stackoverflow.com/questions/70585852/return-default-value-for-generic-type
		return *new(V), fmt.Errorf("no cache entry under key : %v", key)
	}

	metrics.HitMissCache.Inc()
	metrics.CacheEvent.WithLabelValues("cache_hit").Inc()
	lr.mutex.Lock()
	defer lr.mutex.Unlock()
	lr.l.MoveToFront(el)

	return el.Value.(LRUEntry[K, V]).V, nil
}

func (lr *LRU[K, V]) Set(k K, v V) V {
	if el, ok := lr.m[k]; ok {
		metrics.CacheEvent.WithLabelValues("already_set").Inc()
		lr.mutex.Lock()
		defer lr.mutex.Unlock()
		lr.l.MoveToFront(el)

		return v
	}

	if lr.l.Len() == lr.cap {
		metrics.CacheEvent.WithLabelValues("full_cache_lru_element_evicted").Inc()
		backEl := lr.l.Back()

		pair := backEl.Value.(LRUEntry[K, V])

		delete(lr.m, pair.K)
		lr.mutex.Lock()
		lr.l.Remove(backEl)
		lr.mutex.Unlock()
	}

	metrics.CacheEvent.WithLabelValues("cache_set").Inc()
	pair := LRUEntry[K, V]{
		K: k,
		V: v,
	}

	lr.mutex.Lock()
	lr.m[k] = lr.l.PushFront(pair)
	lr.mutex.Unlock()

	return v
}
