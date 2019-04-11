package repository

import (
	"container/list"
	"fmt"

	"github.com/bxcodec/gotcha/cache"
)

// Repository ...
type Repository struct {
	frequencyList *list.List
	byKey         map[string]*cacheItem
	freqHead      *list.Element
}

type cacheItem struct {
	data   *cache.Document
	parent *list.Element
}

type frequencyItem struct {
	key       string
	frequency uint64
}

func NewRepository() (repo *Repository) {
	repo = &Repository{
		frequencyList: list.New(),
		byKey:         make(map[string]*cacheItem),
		freqHead: &list.Element{
			Value: &frequencyItem{},
		},
	}
	return
}

// Get ...
func (r *Repository) Get(key string) (res *cache.Document, err error) {
	elem, ok := r.byKey[key]
	if elem == nil || !ok {
		err = cache.ErrMissed
		return
	}
	res = elem.data

	freq := elem.parent
	nextFreq := freq.Next()
	freqVal := (freq.Value.(*frequencyItem))
	if nextFreq == nil {
		freqVal.frequency++
		freq.Value = freqVal
		r.frequencyList.MoveAfter(freq, freq)
		fmt.Println("++++++++")
		r.printList()
		fmt.Println("++++++++")
		return
	}

	nextFreqVal := (nextFreq.Value.(*frequencyItem))
	if nextFreq == r.freqHead || freqVal.frequency != (nextFreqVal.frequency+1) {
		freqVal.frequency = (nextFreqVal.frequency + 1)
		freq.Value = freqVal
		r.frequencyList.MoveAfter(freq, nextFreq)
	}

	// TODO: (bxcodec)
	// Check Expiry Time

	fmt.Println("++++++++")
	r.printList()
	fmt.Println("++++++++")
	return
}

func (r *Repository) printList() {
	for elem := r.frequencyList.Back(); elem != nil; elem = elem.Prev() {
		first := elem.Value.(*frequencyItem)
		fmt.Printf("Elem Freq: %+v\n", first.frequency)
		fmt.Printf("Elem Doc: %+v\n", first.key)
	}
}

// Set ...
func (r *Repository) Set(doc *cache.Document) (err error) {
	// Check for existing item
	if elem, ok := r.byKey[doc.Key]; ok {
		elem.data = doc
		return nil
	}

	freq := r.freqHead.Next()
	freqVal := &frequencyItem{}
	if freq == nil {
		freqVal = &frequencyItem{
			// doc:       doc,
			key:       doc.Key,
			frequency: 1,
		}
		freq = r.frequencyList.PushFront(freqVal)
	}

	freqVal, _ = freq.Value.(*frequencyItem)
	if freqVal.frequency != 1 {
		freqVal.frequency = 1
		r.frequencyList.MoveAfter(freq, r.freqHead)
	}

	freqVal.key = doc.Key
	cacheItem := &cacheItem{
		data:   doc,
		parent: freq,
	}
	r.byKey[doc.Key] = cacheItem
	// elem := r.fragmentPositionList.PushFront(doc)
	// r.items[doc.Key] = elem

	// fmt.Printf("By Key: %+v\n", r.byKey)
	// fmt.Println("Len", r.frequencyList.Len())
	r.printList()
	fmt.Println("=======")
	return
}

// GetLFU ...
func (r *Repository) GetLFU() (res *cache.Document, err error) {
	elem := r.frequencyList.Front()
	freq := elem.Value.(*frequencyItem)

	cacheItem, ok := r.byKey[freq.key]
	if !ok {
		err = cache.ErrMissed
		return
	}
	res = cacheItem.data
	return
}
