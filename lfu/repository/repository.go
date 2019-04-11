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

// Get ...
func (r *Repository) Get(key string) (res *cache.Document, err error) {
	panic("TODO: (bxcodec)")
	return
}
