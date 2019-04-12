package repository

import (
	"container/list"
	"fmt"
	"reflect"
	"time"

	"github.com/bxcodec/gotcha/cache"
)

// Repository ...
type Repository struct {
	frequencyList *list.List // will store list of frequencyItem
	byKey         map[string]*lfuItem
	maxSize       uint64
	maxMemory     uint64
}

type lfuItem struct {
	freqParent *list.Element
	data       *cache.Document
}

type frequencyItem struct {
	frequency uint64
	// TODO: (bxcodec) Change to Set type if possible
	// In the paper of Prof. Ketan Shah this items using SET
	// since SET is not available in Golang, I just use Map here
	items map[*lfuItem]bool
}

func NewRepository(maxSize, maxMemory uint64) (repo *Repository) {
	repo = &Repository{
		frequencyList: list.New(),
		byKey:         make(map[string]*lfuItem),
		maxMemory:     maxMemory,
		maxSize:       maxSize,
	}
	return
}

func (r *Repository) Get(key string) (res *cache.Document, err error) {
	tmp := r.byKey[key]
	if tmp == nil {
		err = cache.ErrMissed
		return
	}
	res = tmp.data
	freq := tmp.freqParent
	nextFreq := freq.Next()
	if nextFreq == nil {
		nextFreq = freq
	}

	freqVal := freq.Value.(*frequencyItem)
	nextFreqVal := nextFreq.Value.(*frequencyItem)
	headFreq := r.frequencyList.Front()
	if nextFreq == headFreq || nextFreqVal.frequency != (freqVal.frequency+1) {
		newNodeFreq := &frequencyItem{
			frequency: freqVal.frequency + 1,
		}
		nextFreq = r.frequencyList.InsertAfter(newNodeFreq, freq)
	}

	nextFreqVal = nextFreq.Value.(*frequencyItem)
	if len(nextFreqVal.items) == 0 {
		nextFreqVal.items = make(map[*lfuItem]bool)
	}
	nextFreqVal.items[tmp] = true
	tmp.freqParent = nextFreq
	delete(freqVal.items, tmp)
	if len(freqVal.items) == 0 {
		r.frequencyList.Remove(freq)
	}

	// TODO: (bxcodec)
	//  Check Expiry and Remove expired item
	return
}

func (r *Repository) Set(doc *cache.Document) (err error) {
	if _, ok := r.byKey[doc.Key]; ok {
		// TODO: (bxcodec)
		// Re-insert the document
		return
	}

	freq := r.frequencyList.Front() // Front will always be the least frequently used
	freqVal := &frequencyItem{}

	if freq == nil {
		newNodeFreq := &frequencyItem{
			frequency: 1,
		}
		freq = r.frequencyList.PushFront(newNodeFreq)
		freqVal = freq.Value.(*frequencyItem)
		item := &lfuItem{
			freqParent: freq,
			data:       doc,
		}

		freqVal.items = map[*lfuItem]bool{
			item: true,
		}
		r.byKey[doc.Key] = item
		return
	}

	freqVal = freq.Value.(*frequencyItem)
	if freqVal.frequency != 1 {
		newNodeFreq := &frequencyItem{
			frequency: 1,
			items:     make(map[*lfuItem]bool),
		}

		freq = r.frequencyList.InsertAfter(newNodeFreq, freq)
		fmt.Println("Kye", doc.Key)
	}

	freqVal = freq.Value.(*frequencyItem)
	item := &lfuItem{
		freqParent: freq,
		data:       doc,
	}

	freqVal.items[item] = true
	r.byKey[doc.Key] = item

	// Remove oldest
	if uint64(len(r.byKey)) > r.maxSize {
		r.removeLfuOldest()
	}
	return
}

func (r *Repository) printList() {
	for elem := r.frequencyList.Front(); elem != nil; elem = elem.Next() {
		first := elem.Value.(*frequencyItem)
		fmt.Printf("Elem Freq: %+v\n", first.frequency)
		for item, _ := range first.items {
			fmt.Printf("\tElem Doc: %+v\n", item)
		}
	}
}

func (r *Repository) removeLfuOldest() (oldestItem *lfuItem) {
	lfuList := r.frequencyList.Front()
	freqItem := lfuList.Value.(*frequencyItem)

	minStoreTime := time.Now().Unix()
	// Search for the oldest one with store time
	for item := range freqItem.items {

		if item.data.StoredTime < minStoreTime {
			minStoreTime = item.data.StoredTime
			oldestItem = item
		}
	}

	if oldestItem == nil && len(freqItem.items) > 0 {
		//  Get randomly
		oldestItem = reflect.ValueOf(freqItem.items).MapKeys()[0].Interface().(*lfuItem)
	}

	// Remove from Cache
	delete(freqItem.items, oldestItem)
	delete(r.byKey, oldestItem.data.Key)
	if len(freqItem.items) == 0 {
		r.frequencyList.Remove(lfuList)
	}
	return
}
