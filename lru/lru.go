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
	Delete(key string) (ok bool, err error)
	Keys() (keys []string, err error)
}

// NewCache return the implementations of cache with LRU algorithm
func NewCache(option cache.Option) cache.Interactor {
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
// TODO: (bxcodec)
// Add Test for this function
func (c *Cache) Set(key string, value interface{}) (err error) {
	document := &cache.Document{
		Key:        key,
		Value:      value,
		StoredTime: time.Now().Unix(),
	}
	c.Lock()
	c.repo.Set(document)
	c.Unlock()
	return
}

// Get ...
// TODO: (bxcodec)
// Add Test for this function
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
// TODO: (bxcodec)
// Add Test for this function
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
// TODO: (bxcodec)
// Add Test for this function
func (c *Cache) GetKeys() (keys []string, err error) {
	c.RLock()
	keys, err = c.repo.Keys()
	c.RUnlock()
	return
}

// ClearCache ...
// TODO: (bxcodec)
// Add Test for this function
func (c *Cache) ClearCache() (err error) {
	c.Lock()
	err = c.repo.Clear()
	c.Unlock()
	return
}
