package gotcha

import "time"

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
	GetKeys() (keys string, err error)
}

const (
	// LRUCacheAlgorithm ...
	LRUCacheAlgorithm = "lru"
	// LFUCacheAlgorithm ...
	LFUCacheAlgorithm = "lfu"
)

// CacheOption ...
type CacheOption struct {
	AlgorithmType string        // represent the algorithm type
	ExpiryTime    time.Duration // represent the expiry time of each stored item
	MaxSizeItem   int64         // Max size of item for eviction
	MaxMemory     int64         // Max Memory of item stored for eviction
}

// New ...
func New(option *CacheOption) (c CacheInteractor) {
	panic("TODO")
	return
}
