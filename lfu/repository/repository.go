package repository

import (
	"container/list"

	"github.com/bxcodec/gotcha/cache"
)

// Repository ...
type Repository struct {
	frequencyList *list.List
	byKey         map[string]*list.Element
}

type cacheItem struct {
	doc      *cache.Document
	freqHead *list.Element
}

type frequencyItem struct {
	entries   map[string]bool
	frequency int
}

// Get ...
func (r *Repository) Get(key string) (res *cache.Document, err error) {
	elem, ok := r.byKey[key]
	if elem == nil || !ok {
		err = cache.ErrMissed
		return
	}
	r.addFrequency(elem)
	res = elem.Value.(*cache.Document)
	return
}

func (r *Repository) addFrequency(elem *list.Element) {
	// TODO: (bxcodec)
}

// Set ...
func (r *Repository) Set(doc *cache.Document) (err error) {
	// Check for existing item
	if elem, ok := r.byKey[doc.Key]; ok {
		elem.Value = doc
		return nil
	}

	// elem := r.fragmentPositionList.PushFront(doc)
	// r.items[doc.Key] = elem

	panic("TODO")
	return
}
