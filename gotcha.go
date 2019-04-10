package gotcha

import (
	"errors"
	"time"
)

// Document ...
type Document struct {
	Key        string
	Value      interface{}
	StoredTime time.Time
}

// CacheInteractor ...
type CacheInteractor interface {
	Set(key string, value interface{}) error
	Get(key string) (val interface{}, err error)
	Delete(key string) (err error)
	GetKeys() (keys []string, err error)
}

const (
	// LRUCacheAlgorithm ...
	LRUCacheAlgorithm = "lru"
	// LFUCacheAlgorithm ...
	LFUCacheAlgorithm = "lfu"
)

var (
	// ErrCacheMissed ...
	ErrCacheMissed = errors.New("Cache item's missing")
)

// CacheOption ...
type CacheOption struct {
	AlgorithmType string        // represent the algorithm type
	ExpiryTime    time.Duration // represent the expiry time of each stored item
	MaxSizeItem   uint64        // Max size of item for eviction
	MaxMemory     uint64        // Max Memory of item stored for eviction
}

// New ...
func New(option *CacheOption) (c CacheInteractor) {
	panic("TODO")
	return
}
