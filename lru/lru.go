package lru

import (
	"sync"
	"time"

	"github.com/bxcodec/gotcha/cache"
	"github.com/bxcodec/gotcha/lru/repository"
)

// Repository ...
type Repository interface {
	Set(doc *cache.Document) (err error)
	Get(key string) (res *cache.Document, err error)
	Clear() (err error)
	Contains(key string) (ok bool)
	Peek(key string) (res *cache.Document, err error)
	Delete(key string) (ok bool, err error)
	RemoveOldest() (res *cache.Document, err error)
	GetOldest() (res *cache.Document, err error)
	Keys() (keys []string, err error)
	Len() (len int64, err error)
	MemoryUsage() (size int64, err error)
}

// NewCache return the implementations of cache with LRU algorithm
func NewCache(option cache.Option) cache.Interactor {
	repo := repository.New(option.MaxSizeItem, option.MaxMemory)
	return &Cache{
		Option: option,
		repo:   repo,
	}
}

// Cache ...
type Cache struct {
	sync.RWMutex
	repo   Repository
	Option cache.Option
}

// Set ...
func (c *Cache) Set(key string, value interface{}) (err error) {
	document := &cache.Document{
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

// ClearCache ...
func (c *Cache) ClearCache() (err error) {
	c.Lock()
	err = c.repo.Clear()
	c.Unlock()
	return
}
