package gotcha

import (
	"github.com/bxcodec/gotcha/cache"
	"github.com/bxcodec/gotcha/lfu"
	"github.com/bxcodec/gotcha/lru"
)

var (
	// DefaultCache use for default cache client
	DefaultCache = New()
)

// New will create a new cache client. If the options not set, the cache will use the default options
func New(options ...*cache.Option) (c cache.Cache) {
	option := mergeOptions(options...)
	if option.MaxMemory < cache.DefaultMaxMemory { // Unlimited
		option.MaxMemory = cache.DefaultMaxMemory
	}
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

	switch option.AlgorithmType {
	case cache.LRUAlgorithm:
		c = lru.NewCache(*option)
	case cache.LFUAlgorithm:
		c = lfu.NewCache(*option)
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
