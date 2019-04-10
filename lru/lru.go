package lru

import (
	"container/list"
	"sync"
	"time"

	"github.com/bxcodec/gotcha"
)

// Repository ...
type Repository interface {
	Set(doc *gotcha.Document) (err error)
	Get(key string) (res *gotcha.Document, err error)
	Clear() (err error)
	Contains(key string) (ok bool)
	Peek(key string) (res *gotcha.Document, ok bool)
	Delete(key string) (ok bool, err error)
	RemoveOldest() (res *gotcha.Document, ok bool)
	GetOldest() (res *gotcha.Document, ok bool)
	Keys() []interface{}
	Len()
	MemoryUsage() (size int64, err error)
}

// Cache ...
type Cache struct {
	sync.RWMutex
	DocumentList list.List
	Documents    map[string]gotcha.Document
	// Fragments                       []gotcha.Document
	// DocumentFragmentsPositionMapper map[int]gotcha.Document
	Option gotcha.CacheOption
}

// Set ...
func (c *Cache) Set(key string, value interface{}) (err error) {
	document := gotcha.Document{
		Key:        key,
		Value:      value,
		StoredTime: time.Now(),
	}
	c.Lock()
	if sizeFit(c.Documents, c.Option.MaxSizeItem) {
		// document.Position =
		c.Documents[key] = document
	}
	c.Unlock()

	panic("TODO")
	return
}

func sizeFit(docs map[string]gotcha.Document, maxSize uint64) (ok bool) {
	ok = uint64(len(docs)) < maxSize
	return
}

// Get ...
func (c *Cache) Get(key string) (value interface{}, err error) {
	panic("TODO")
	return
}

// Delete ...
func (c *Cache) Delete(key string) (err error) {
	panic("TODO")
	return
}
