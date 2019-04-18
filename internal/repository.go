package internal

import (
	"github.com/bxcodec/gotcha/cache"
)

type Repository interface {
	Set(doc *cache.Document) (err error)
	Get(key string) (res *cache.Document, err error)
	Clear() (err error)
	Contains(key string) (ok bool)
	Delete(key string) (ok bool, err error)
	Keys() (keys []string, err error)
}
