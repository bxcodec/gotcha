package gotcha

import (
	"github.com/bxcodec/gotcha/cache"
	"github.com/bxcodec/gotcha/lru"
)

// New ...
func New(options ...*cache.Option) (c cache.Interactor) {
	option := mergeOptions(options...)
	if option.MaxMemory == 0 { // Unlimited
		// option.MaxMemory = (get max memory)
	}

	if option.MaxSizeItem == 0 {
		// Use default
		option.MaxSizeItem = cache.DefaultCacheSize
	}

	if option.AlgorithmType == "" {
		// Use LRU Default
		option.AlgorithmType = cache.LRUCacheAlgorithm
	}

	if option.ExpiryTime == 0 {
		// Use default expiry time
		option.ExpiryTime = cache.DefaultExpiryTime
	}

	switch option.AlgorithmType {
	case cache.LRUCacheAlgorithm:
		c = lru.NewCache(*option)
	case cache.LFUCacheAlgorithm:
	}

	panic("TODO")
	return
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
