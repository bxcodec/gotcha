package lru

import (
	"sync"
	"time"

	"github.com/bxcodec/gotcha"
)

// Repository ...
type Repository interface {
	Set(doc *gotcha.Document) (err error)
	Get(key string) (res *gotcha.Document, err error)
	Clear() (err error)
	Contains(key string) (ok bool)
	Peek(key string) (res *gotcha.Document, ok bool)
	Delete(key string) (ok bool, err error)
	RemoveOldest() (res *gotcha.Document, err error)
	GetOldest() (res *gotcha.Document, err error)
	Keys() (keys []string, err error)
	Len() (len int64, err error)
	MemoryUsage() (size int64, err error)
}

// NewLRUCache return the implementations of cache with LRU algorithm
func NewLRUCache(option gotcha.CacheOption) gotcha.CacheInteractor {
	return &Cache{
		Option: option,
	}
}

// Cache ...
type Cache struct {
	sync.RWMutex
	repo   Repository
	Option gotcha.CacheOption
}

// Set ...
func (c *Cache) Set(key string, value interface{}) (err error) {
	document := &gotcha.Document{
		Key:        key,
		Value:      value,
		StoredTime: time.Now(),
	}
	c.Lock()
	c.repo.Set(document)
	c.Unlock()
	return
}

// Get ...
func (c *Cache) Get(key string) (value interface{}, err error) {
	c.RLock()
	doc, err := c.repo.Get(key)
	c.RUnlock()
	if err != nil {
		return
	}
	value = doc.Value
	return
}

// Delete ...
func (c *Cache) Delete(key string) (err error) {
	c.Lock()
	_, err = c.repo.Delete(key)
	c.Unlock()
	if err != nil {
		return
	}
	return
}

// GetKeys ...
func (c *Cache) GetKeys() (keys []string, err error) {
	c.RLock()
	keys, err = c.repo.Keys()
	c.RUnlock()
	return
}
