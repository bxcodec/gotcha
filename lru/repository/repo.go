package repository

import (
	"container/list"

	"github.com/bxcodec/gotcha"
)

// Repository implements the Repository cache
type Repository struct {
	maxsize              uint64
	maxmemory            uint64
	fragmentPositionList *list.List
	items                map[string]*list.Element
}

// New constructs an Repository of the given size
func New(size uint64, memory uint64) *Repository {
	c := &Repository{
		maxsize:              size,
		fragmentPositionList: list.New(),
		items:                make(map[string]*list.Element),
	}
	return c
}

// Set adds a value to the cache.  Returns true if an eviction occurred.
func (r *Repository) Set(doc *gotcha.Document) (err error) {
	// Check for existing item
	if elem, ok := r.items[doc.Key]; ok {
		// TODO: (bxcodec)
		// Check the expiry item
		r.fragmentPositionList.MoveToFront(elem)
		elem.Value.(*gotcha.Document).Value = doc.Value
		return nil
	}

	elem := r.fragmentPositionList.PushFront(doc)
	r.items[doc.Key] = elem

	// Remove the oldest if the fragment is full
	if uint64(r.fragmentPositionList.Len()) > r.maxsize {
		r.removeOldest()
	}

	// TODO: (bxcodec)
	// Remove the oldest if the memory is reach the maximun
	return nil
}

// Get looks up a key's value from the cache.
func (r *Repository) Get(key string) (res *gotcha.Document, err error) {
	if elem, ok := r.items[key]; ok {
		r.fragmentPositionList.MoveToFront(elem)
		return elem.Value.(*gotcha.Document), nil
	}
	err = gotcha.ErrCacheMissed
	return
}

// GetOldest returns the oldest entry
func (r *Repository) GetOldest() (res *gotcha.Document, err error) {
	elem := r.fragmentPositionList.Back()
	if elem != nil {
		res = elem.Value.(*gotcha.Document)
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
func (r *Repository) Peek(key string) (res *gotcha.Document, err error) {
	if elem, ok := r.items[key]; ok {
		res = elem.Value.(*gotcha.Document)
		return
	}
	err = gotcha.ErrCacheMissed
	return
}

// Delete removes the provided key from the cache, returning if the
// key was contained.
func (r *Repository) Delete(key string) (ok bool, err error) {
	elem, ok := r.items[key]
	if ok {
		r.removeElement(elem)
		return
	}
	return false, nil
}

// removeElement is used to remove a given list element from the cache
func (r *Repository) removeElement(e *list.Element) {
	r.fragmentPositionList.Remove(e)
	doc := e.Value.(*gotcha.Document)
	delete(r.items, doc.Key)
}

// RemoveOldest removes the oldest item from the cache.
func (r *Repository) RemoveOldest() (res *gotcha.Document, err error) {
	elem := r.fragmentPositionList.Back()
	if elem != nil {
		r.removeElement(elem)
		res = elem.Value.(*gotcha.Document)
		return
	}
	return
}

// removeOldest removes the oldest item from the cache.
func (r *Repository) removeOldest() {
	elem := r.fragmentPositionList.Back()
	if elem != nil {
		r.removeElement(elem)
	}
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (r *Repository) Keys() (keys []string, err error) {
	keys = make([]string, len(r.items))
	i := 0
	for elem := r.fragmentPositionList.Back(); elem != nil; elem = elem.Prev() {
		keys[i] = elem.Value.(*gotcha.Document).Key
		i++
	}
	return
}

// Len returns the number of items in the cache.
func (r *Repository) Len() (len int64, err error) {
	len = int64(r.fragmentPositionList.Len())
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
	r.fragmentPositionList.Init()
	return
}
