package lfu

import (
	"sync"
	"time"

	"github.com/bxcodec/gotcha/cache"
	"github.com/bxcodec/gotcha/lfu/repository"
)

// Repository ...
type Repository interface {
	Set(doc *cache.Document) (err error)
	Get(key string) (res *cache.Document, err error)
	Clear() (err error)
	Contains(key string) (ok bool)
	Delete(key string) (ok bool, err error)
	Keys() (keys []string, err error)
}

// NewCache return the implementations of cache with LRU algorithm
func NewCache(option cache.Option) cache.Cache {
	repo := repository.New(option.MaxSizeItem, option.MaxMemory, option.ExpiryTime)
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
	doc := &cache.Document{
		Key:        key,
		Value:      value,
		StoredTime: time.Now().Unix(),
	}

	c.Lock()
	err = c.repo.Set(doc)
	c.Unlock()
	return
}

// Get ...
func (c *Cache) Get(key string) (val interface{}, err error) {
	c.RLock()
	doc, err := c.repo.Get(key)
	c.RUnlock()
	if err != nil {
		return
	}
	val = doc.Value
	return
}

// Delete ...
func (c *Cache) Delete(key string) (err error) {
	c.Lock()
	_, err = c.repo.Delete(key)
	c.Unlock()
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
