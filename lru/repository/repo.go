package repository

import (
	"container/list"

	"github.com/bxcodec/gotcha"
)

// Repository implements the Repository cache
type Repository struct {
	maxsize   uint64
	maxmemory uint64
	evictList *list.List
	items     map[string]*list.Element
}

// New constructs an Repository of the given size
func New(size uint64, memory uint64) *Repository {
	c := &Repository{
		maxsize:   size,
		evictList: list.New(),
		items:     make(map[string]*list.Element),
	}
	return c
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

// GetOldest returns the oldest entry
func (r *Repository) GetOldest() (res *gotcha.Document, err error) {
	ent := r.evictList.Back()
	if ent != nil {
		res = ent.Value.(*gotcha.Document)
		return
	}
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

// removeElement is used to remove a given list element from the cache
func (r *Repository) removeElement(e *list.Element) {
	r.evictList.Remove(e)
	doc := e.Value.(*gotcha.Document)
	delete(r.items, doc.Key)
}

// RemoveOldest removes the oldest item from the cache.
func (r *Repository) RemoveOldest() (res *gotcha.Document, err error) {
	ent := r.evictList.Back()
	if ent != nil {
		r.removeElement(ent)
		res = ent.Value.(*gotcha.Document)
		return
	}
	return
}

// removeOldest removes the oldest item from the cache.
func (r *Repository) removeOldest() {
	ent := r.evictList.Back()
	if ent != nil {
		r.removeElement(ent)
	}
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (r *Repository) Keys() (keys []string, err error) {
	keys = make([]string, len(r.items))
	i := 0
	for ent := r.evictList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*gotcha.Document).Key
		i++
	}
	return
}

// Len returns the number of items in the cache.
func (r *Repository) Len() (len int64, err error) {
	len = int64(r.evictList.Len())
	return
}

// MemoryUsage returns the number of memory usage for all cache item
func (r *Repository) MemoryUsage() (size int64, err error) {
	panic("TODO: (bxcodec)")
	return
}

// Clear is used to completely clear the cache.
func (r *Repository) Clear() (err error) {
	for k := range r.items {
		delete(r.items, k)
	}
	r.evictList.Init()
	return
}
