package db

import (
	"sync"
	"time"
)

type dbEntry struct {
	value     string
	expiresAt time.Time
}

type DB struct {
	mu   sync.Mutex
	data map[string]dbEntry
}

func NewDB() *DB {
	return &DB{
		data: make(map[string]dbEntry),
	}
}

func (d *DB) Set(key string, value string, ttl time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	d.data[key] = dbEntry{
		value:     value,
		expiresAt: expiresAt,
	}
}

func (d *DB) Get(key string) (string, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	entry, exists := d.data[key]
	if !exists {
		return "", false
	}

	if !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
		delete(d.data, key) 
		return "", false
	}

	return entry.value, true
}

// -2 if key does not exist or expired
// -1 if key lives forever
// positive number representing seconds left
func (d *DB) TTL(key string) int64 {
	d.mu.Lock()
	defer d.mu.Unlock()

	entry, exists := d.data[key]
	if !exists {
		return -2
	}

	if !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
		delete(d.data, key)
		return -2
	}

	if entry.expiresAt.IsZero() {
		return -1
	}

	remaining := time.Until(entry.expiresAt).Seconds()
	return int64(remaining)
}