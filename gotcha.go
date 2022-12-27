package gotcha

import (
	"sync"
	"time"

	"github.com/bxcodec/gotcha/cache"
	"github.com/bxcodec/gotcha/internal"
	"github.com/bxcodec/gotcha/internal/lfu"
	"github.com/bxcodec/gotcha/internal/lru"
)

var (
	// DefaultCache use for default cache client
	DefaultCache = New()
)

// New will create a new cache client. If the options not set, the cache will use the default options
func New(options ...*cache.Option) (c cache.Cache) {
	option := mergeOptions(options...)
	if option.MaxSizeItem == 0 {
		// Use default
		option.MaxSizeItem = cache.DefaultSize
	}
	if option.AlgorithmType == "" {
		// Use LRU Default
		option.AlgorithmType = cache.LRUAlgorithm
	}
	if option.ExpiryTime == 0 {
		// Use default expiry time
		option.ExpiryTime = cache.DefaultExpiryTime
	}

	c = &Cache{
		repo:  NewRepository(*option),
		mutex: &sync.RWMutex{},
	}
	return
}

// NewOption return an empty option
func NewOption() (op *cache.Option) {
	return &cache.Option{}
}

func mergeOptions(options ...*cache.Option) (opts *cache.Option) {
	opts = new(cache.Option)
	for _, op := range options {
		if op.AlgorithmType != "" {
			opts.AlgorithmType = op.AlgorithmType
		}
		if op.ExpiryTime != 0 {
			opts.ExpiryTime = op.ExpiryTime
		}
		if op.MaxMemory != 0 {
			opts.MaxMemory = op.MaxMemory
		}
		if op.MaxSizeItem != 0 {
			opts.MaxSizeItem = op.MaxSizeItem
		}
	}
	return
}

// Set will set an item to cache using default option
func Set(key string, value interface{}) (err error) {
	return DefaultCache.Set(key, value)
}

// Get will get an item from cache using default option
func Get(key string) (value interface{}, err error) {
	return DefaultCache.Get(key)
}

// Delete will delete an item from the cache using default option
func Delete(key string) (err error) {
	return DefaultCache.Delete(key)
}

// GetKeys will get all keys from the cache using default option
func GetKeys() (keys []string, err error) {
	return DefaultCache.GetKeys()
}

// ClearCache will Clear the cache using default option
func ClearCache() (err error) {
	return DefaultCache.ClearCache()
}

// NewRepository return the implementations of repository cache
func NewRepository(option cache.Option) internal.Repository {
	var repo internal.Repository
	switch option.AlgorithmType {
	case cache.LRUAlgorithm:
		repo = lru.New(option.MaxSizeItem, option.MaxMemory, option.ExpiryTime)
	case cache.LFUAlgorithm:
		repo = lfu.New(option.MaxSizeItem, option.MaxMemory, option.ExpiryTime)
	}
	return repo
}

// Cache represent the Cache handler
type Cache struct {
	mutex *sync.RWMutex
	repo  internal.Repository
}

// Set used for setting the item to cache
// TODO: (bxcodec)
// Add Test for this function
func (c *Cache) Set(key string, value interface{}) (err error) {
	document := &cache.Document{
		Key:        key,
		Value:      value,
		StoredTime: time.Now().Unix(),
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	err = c.repo.Set(document)
	return
}

// Get will retrieve the item from cache
// TODO: (bxcodec)
// Add Test for this function
func (c *Cache) Get(key string) (value interface{}, err error) {
	c.mutex.RLock()
	doc, err := c.repo.Get(key)
	c.mutex.RUnlock()
	if err != nil {
		return
	}
	value = doc.Value
	return
}

// Delete will remove the item from cache
// TODO: (bxcodec)
// Add Test for this function
func (c *Cache) Delete(key string) (err error) {
	c.mutex.Lock()
	_, err = c.repo.Delete(key)
	c.mutex.Unlock()
	if err != nil {
		return
	}
	return
}

// GetKeys will retrieve all keys from cache
// TODO: (bxcodec)
// Add Test for this function
func (c *Cache) GetKeys() (keys []string, err error) {
	c.mutex.RLock()
	keys, err = c.repo.Keys()
	c.mutex.RUnlock()
	return keys, err
}

// ClearCache will cleanup all the cache
// TODO: (bxcodec)
// Add Test for this function
func (c *Cache) ClearCache() (err error) {
	c.mutex.Lock()
	err = c.repo.Clear()
	c.mutex.Unlock()
	return
}
