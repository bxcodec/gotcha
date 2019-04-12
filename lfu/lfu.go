package lfu

import (
	"sync"

	"github.com/bxcodec/gotcha/cache"
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
	return &Cache{
		Option: option,
	}
}

// Cache ...
type Cache struct {
	sync.RWMutex
	repo   Repository
	Option cache.Option
}

// Set ...
func (c *Cache) Set(key string, value interface{}) error {
	panic("TODO: (bxcodec)")
}

// Get ...
func (c *Cache) Get(key string) (val interface{}, err error) {
	panic("TODO: (bxcodec)")
}

// Delete ...
func (c *Cache) Delete(key string) (err error) {
	panic("TODO: (bxcodec)")
}

// GetKeys ...
func (c *Cache) GetKeys() (keys []string, err error) {
	panic("TODO: (bxcodec)")
}

// ClearCache ...
func (c *Cache) ClearCache() (err error) {
	panic("TODO: (bxcodec)")
}
