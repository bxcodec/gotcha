package repository

import (
	"container/list"
	"errors"

	"github.com/bxcodec/gotcha"
)

// Repository implements the Repository cache
type Repository struct {
	maxsize   uint64
	maxmemory uint64
	evictList *list.List
	items     map[string]*list.Element
}

// NewRepository constructs an Repository of the given size
func NewRepository(size uint64, memory uint64) (*Repository, error) {
	if size <= 0 {
		return nil, errors.New("Must provide a positive size")
	}
	c := &Repository{
		maxsize:   size,
		evictList: list.New(),
		items:     make(map[string]*list.Element),
	}
	return c, nil
}

// Clear is used to completely clear the cache.
func (r *Repository) Clear() (err error) {
	for k := range r.items {
		delete(r.items, k)
	}
	r.evictList.Init()
	return
}

// Set adds a value to the cache.  Returns true if an eviction occurred.
func (r *Repository) Set(doc *gotcha.Document) (err error) {
	// Check for existing item
	if ent, ok := r.items[doc.Key]; ok {
		// TODO: (bxcodec)
		// Check the expiry item
		r.evictList.MoveToFront(ent)
		ent.Value.(*gotcha.Document).Value = doc.Value
		return nil
	}

	entry := r.evictList.PushFront(doc)
	r.items[doc.Key] = entry
	evict := uint64(r.evictList.Len()) > r.maxsize
	// Verify size not exceeded
	if evict {
		r.removeOldest()
	}
	return nil
}

// Get looks up a key's value from the cache.
func (r *Repository) Get(key string) (res *gotcha.Document, err error) {
	if ent, ok := r.items[key]; ok {
		r.evictList.MoveToFront(ent)
		return ent.Value.(*gotcha.Document), nil
	}
	err = gotcha.ErrCacheMissed
	return
}

// Contains checks if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (r *Repository) Contains(key string) (ok bool) {
	_, ok = r.items[key]
	return ok
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (r *Repository) Peek(key string) (res *gotcha.Document, ok bool) {
	var ent *list.Element
	if ent, ok = r.items[key]; ok {
		return ent.Value.(*gotcha.Document), true
	}
	return nil, ok
}

// Delete removes the provided key from the cache, returning if the
// key was contained.
func (r *Repository) Delete(key string) (ok bool, err error) {
	ent, ok := r.items[key]
	if ok {
		r.removeElement(ent)
		return
	}
	return false, nil
}

// RemoveOldest removes the oldest item from the cache.
func (r *Repository) RemoveOldest() (res *gotcha.Document, ok bool) {
	ent := r.evictList.Back()
	if ent != nil {
		r.removeElement(ent)
		res, ok = ent.Value.(*gotcha.Document)
		return
	}
	return
}

// GetOldest returns the oldest entry
func (r *Repository) GetOldest() (res *gotcha.Document, ok bool) {
	ent := r.evictList.Back()
	if ent != nil {
		res, ok = ent.Value.(*gotcha.Document)
		return
	}
	return
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (r *Repository) Keys() []interface{} {
	keys := make([]interface{}, len(r.items))
	i := 0
	for ent := r.evictList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*gotcha.Document).Key
		i++
	}
	return keys
}

// Len returns the number of items in the cache.
func (r *Repository) Len() int {
	return r.evictList.Len()
}

// MemoryUsage returns the number of memory usage for all cache item
func (r *Repository) MemoryUsage() (size int64, err error) {
	panic("TODO: (bxcodec)")
	return
}

// removeOldest removes the oldest item from the cache.
func (r *Repository) removeOldest() {
	ent := r.evictList.Back()
	if ent != nil {
		r.removeElement(ent)
	}
}

// removeElement is used to remove a given list element from the cache
func (r *Repository) removeElement(e *list.Element) {
	r.evictList.Remove(e)
	kv := e.Value.(*gotcha.Document)
	delete(r.items, kv.Key)
}
