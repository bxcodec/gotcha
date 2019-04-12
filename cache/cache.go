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
	// LRUAlgorithm ...
	LRUAlgorithm = "lru"
	// LFUAlgorithm ...
	LFUAlgorithm = "lfu"
	// DefaultSize ..
	DefaultSize = 100
	// DefaultExpiryTime ...
	DefaultExpiryTime = time.Second * 10
)

// Document ...
type Document struct {
	Key        string
	Value      interface{}
	StoredTime int64 //timestamp
}

// Option ...
type Option struct {
	AlgorithmType string        // represent the algorithm type
	ExpiryTime    time.Duration // represent the expiry time of each stored item
	MaxSizeItem   uint64        // Max size of item for eviction
	MaxMemory     uint64        // Max Memory of item stored for eviction
}

// Interactor ...
type Interactor interface {
	Set(key string, value interface{}) error
	Get(key string) (val interface{}, err error)
	Delete(key string) (err error)
	GetKeys() (keys []string, err error)
	ClearCache() (err error)
}
