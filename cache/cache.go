package cache

import (
	"errors"
	"time"
)

var (
	// ErrMissed ...
	ErrMissed = errors.New("Cache item's missing")
)

const (
	// Byte ...
	Byte uint64 = 1
	// KB ...
	KB = Byte * 1024
	// MB ...
	MB = KB * 1024
	// LRUAlgorithm ...
	LRUAlgorithm = "lru"
	// LFUAlgorithm ...
	LFUAlgorithm = "lfu"
	// DefaultSize ..
	DefaultSize = 100
	// DefaultExpiryTime ...
	DefaultExpiryTime = time.Second * 10
	// DefaultAlgorithm ...
	DefaultAlgorithm = LRUAlgorithm
	// DefaultMaxMemory ...
	DefaultMaxMemory = 10 * MB
)

// Document represent the Document structure stored in the cache
type Document struct {
	Key        string
	Value      interface{}
	StoredTime int64 // timestamp
}

// Option used for Cache configuration
type Option struct {
	AlgorithmType string        // represent the algorithm type
	ExpiryTime    time.Duration // represent the expiry time of each stored item
	MaxSizeItem   uint64        // Max size of item for eviction
	MaxMemory     uint64        // Max Memory of item stored for eviction
}

// SetAlgorithm will set the algorithm value
func (o *Option) SetAlgorithm(algorithm string) *Option {
	o.AlgorithmType = algorithm
	return o
}

// SetExpiryTime will set the expiry time
func (o *Option) SetExpiryTime(expiry time.Duration) *Option {
	o.ExpiryTime = expiry
	return o
}

// SetMaxSizeItem will set the maximum size of item in cache
func (o *Option) SetMaxSizeItem(size uint64) *Option {
	o.MaxSizeItem = size
	return o
}

// SetMaxMemory will set the maximum memory will used for cache
func (o *Option) SetMaxMemory(memory uint64) *Option {
	o.MaxMemory = memory
	return o
}

// Cache represent the public API that will available used by user
type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) (val interface{}, err error)
	Delete(key string) (err error)
	GetKeys() (keys []string, err error)
	ClearCache() (err error)
}
