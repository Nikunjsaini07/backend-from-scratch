package db

import (
	"container/heap"
	"sync"
	"time"
)

type expiryItem struct {
	key       string
	expiresAt time.Time
}

type expiryHeap []expiryItem

func (h expiryHeap) Len() int {
	return len(h)
}

func (h expiryHeap) Less(i, j int) bool {
	return h[i].expiresAt.Before(h[j].expiresAt)
}

func (h expiryHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *expiryHeap) Push(x any) {
	*
	h = append(*h, x.(expiryItem))
}

func (h *expiryHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

type DB struct {
	mu          sync.Mutex
	data        map[string]string
	expiresAt   map[string]time.Time
	expiryQueue expiryHeap
}

func NewDB() *DB {
	return &DB{
		data:      make(map[string]string),
		expiresAt: make(map[string]time.Time),
	}
}

func (d *DB) Set(key string, value string, ttl time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if ttl > 0 {
		expiresAt := time.Now().Add(ttl)
		d.expiresAt[key] = expiresAt
		heap.Push(&d.expiryQueue, expiryItem{key: key, expiresAt: expiresAt})
	} else {
		delete(d.expiresAt, key)
	}

	d.data[key] = value
}

func (d *DB) Get(key string) (string, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	value, exists := d.data[key]
	if !exists {
		return "", false
	}

	expiresAt, hasExpiry := d.expiresAt[key]
	if hasExpiry && time.Now().After(expiresAt) {
		delete(d.data, key)
		delete(d.expiresAt, key)
		return "", false
	}

	return value, true
}

func (d *DB) Del(keys ...string) int64 {
	d.mu.Lock()
	defer d.mu.Unlock()

	var deleted int64
	for _, key := range keys {
		_, exists := d.data[key]
		if !exists {
			continue
		}

		delete(d.data, key)
		delete(d.expiresAt, key)
		deleted++
	}

	return deleted
}

func (d *DB) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			d.CleanupExpired()
		}
	}()
}

func (d *DB) CleanupExpired() int64 {
	d.mu.Lock()
	defer d.mu.Unlock()

	var deleted int64
	now := time.Now()

	for d.expiryQueue.Len() > 0 {
		next := d.expiryQueue[0]
		if next.expiresAt.After(now) {
			break
		}

		heap.Pop(&d.expiryQueue)

		currentExpiry, hasExpiry := d.expiresAt[next.key]
		if !hasExpiry || !currentExpiry.Equal(next.expiresAt) {
			continue
		}

		delete(d.data, next.key)
		delete(d.expiresAt, next.key)
		deleted++
	}

	return deleted
}

func (d *DB) Expire(key string, ttl time.Duration) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, exists := d.data[key]
	if !exists {
		return false
	}

	expiresAt, hasExpiry := d.expiresAt[key]
	if hasExpiry && time.Now().After(expiresAt) {
		delete(d.data, key)
		delete(d.expiresAt, key)
		return false
	}

	if ttl <= 0 {
		delete(d.data, key)
		delete(d.expiresAt, key)
		return true
	}

	expiresAt = time.Now().Add(ttl)
	d.expiresAt[key] = expiresAt
	heap.Push(&d.expiryQueue, expiryItem{key: key, expiresAt: expiresAt})
	return true
}

// -2 if key does not exist or expired
// -1 if key lives forever
// positive number representing seconds left
func (d *DB) TTL(key string) int64 {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, exists := d.data[key]
	if !exists {
		return -2
	}

	expiresAt, hasExpiry := d.expiresAt[key]
	if hasExpiry && time.Now().After(expiresAt) {
		delete(d.data, key)
		delete(d.expiresAt, key)
		return -2
	}

	if !hasExpiry {
		return -1
	}

	remaining := time.Until(expiresAt).Seconds()
	return int64(remaining)
}
