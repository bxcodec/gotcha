package repository

import (
	"bytes"
	"container/list"
	"encoding/gob"
	"reflect"
	"time"

	"github.com/bxcodec/gotcha/cache"
)

// Repository ...
type Repository struct {
	frequencyList  *list.List // will store list of frequencyItem
	byKey          map[string]*lfuItem
	maxSize        uint64
	maxMemory      uint64
	expiryTreshold time.Duration
}

type lfuItem struct {
	FreqParent *list.Element
	Data       *cache.Document
}

type frequencyItem struct {
	Frequency uint64
	// TODO: (bxcodec) Change to Set Data structures if possible
	// In the paper of Prof. Ketan Shah this items using SET
	// since SET is not available in Golang, I just use Map here
	items map[*lfuItem]bool
}

// New ...
func New(maxSize, maxMemory uint64, expiryTreshold time.Duration) (repo *Repository) {
	gob.Register(frequencyItem{}) // this is for serialization to count total documents (in bytes) saved in memory
	repo = &Repository{
		frequencyList:  list.New(),
		byKey:          make(map[string]*lfuItem),
		maxMemory:      maxMemory,
		maxSize:        maxSize,
		expiryTreshold: expiryTreshold,
	}
	return
}

// Get ...
func (r *Repository) Get(key string) (res *cache.Document, err error) {
	tmp := r.byKey[key]
	if tmp == nil {
		err = cache.ErrMissed
		return
	}
	res = tmp.Data

	//  Check Expiry and Remove the expired item
	storedTime := time.Unix(res.StoredTime, 0)
	if time.Since(storedTime) > r.expiryTreshold {
		r.Delete(key)
		return nil, cache.ErrMissed
	}

	freq := tmp.FreqParent
	nextFreq := freq.Next()
	if nextFreq == nil {
		nextFreq = freq
	}

	freqVal := freq.Value.(*frequencyItem)
	nextFreqVal := nextFreq.Value.(*frequencyItem)
	headFreq := r.frequencyList.Front()
	if nextFreq == headFreq || nextFreqVal.Frequency != (freqVal.Frequency+1) {
		newNodeFreq := &frequencyItem{
			Frequency: freqVal.Frequency + 1,
		}
		nextFreq = r.frequencyList.InsertAfter(newNodeFreq, freq)
	}

	nextFreqVal = nextFreq.Value.(*frequencyItem)
	if len(nextFreqVal.items) == 0 {
		nextFreqVal.items = make(map[*lfuItem]bool)
	}
	nextFreqVal.items[tmp] = true
	tmp.FreqParent = nextFreq
	delete(freqVal.items, tmp)
	if len(freqVal.items) == 0 {
		r.frequencyList.Remove(freq)
	}

	return
}

// Set ...
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
			Frequency: 1,
		}
		freq = r.frequencyList.PushFront(newNodeFreq)
		freqVal = freq.Value.(*frequencyItem)
		item := &lfuItem{
			FreqParent: freq,
			Data:       doc,
		}

		freqVal.items = map[*lfuItem]bool{
			item: true,
		}
		r.byKey[doc.Key] = item
		return
	}

	freqVal = freq.Value.(*frequencyItem)
	if freqVal.Frequency != 1 {
		newNodeFreq := &frequencyItem{
			Frequency: 1,
			items:     make(map[*lfuItem]bool),
		}
		freq = r.frequencyList.PushFront(newNodeFreq)
	}

	freqVal = freq.Value.(*frequencyItem)
	item := &lfuItem{
		FreqParent: freq,
		Data:       doc,
	}

	freqVal.items[item] = true
	r.byKey[doc.Key] = item

	// TODO: (bxcodec)
	// Move this to go-routine if possible
	// Remove oldest if the maxmemory reached
	byteMap := encodeGob(r.byKey)
	for uint64(len(byteMap)) > r.maxMemory {
		r.removeLfuOldest()
	}

	// Remove oldest if the maxmemory reached
	if uint64(len(r.byKey)) > r.maxSize {
		r.removeLfuOldest()
	}
	return
}

func encodeGob(v interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func (r *Repository) removeLfuOldest() (oldestItem *lfuItem) {
	lfuList := r.frequencyList.Front()
	freqItem := lfuList.Value.(*frequencyItem)

	minStoreTime := time.Now().Unix()
	// Search for the oldest one with store time
	for item := range freqItem.items {
		if item.Data.StoredTime < minStoreTime {
			minStoreTime = item.Data.StoredTime
			oldestItem = item
		}
	}

	if oldestItem == nil && len(freqItem.items) > 0 {
		//  Get randomly
		oldestItem = reflect.ValueOf(freqItem.items).MapKeys()[0].Interface().(*lfuItem)
	}

	// Remove from Cache
	delete(freqItem.items, oldestItem)
	delete(r.byKey, oldestItem.Data.Key)
	if len(freqItem.items) == 0 {
		r.frequencyList.Remove(lfuList)
	}
	return
}

// Clear ...
func (r *Repository) Clear() (err error) {
	for k := range r.byKey {
		delete(r.byKey, k)
	}
	r.frequencyList.Init()
	return
}

// Len ...
func (r *Repository) Len() int {
	return len(r.byKey)
}

// Contains ...
func (r *Repository) Contains(key string) (ok bool) {
	_, ok = r.byKey[key]
	return
}

// Delete ...
func (r *Repository) Delete(key string) (ok bool, err error) {
	lfuItem, ok := r.byKey[key]
	if !ok {
		return
	}

	freqItem := lfuItem.FreqParent.Value.(*frequencyItem)
	delete(freqItem.items, lfuItem)
	if len(freqItem.items) == 0 {
		r.frequencyList.Remove(lfuItem.FreqParent)
	}
	delete(r.byKey, key)
	return
}

// Keys ...
func (r *Repository) Keys() (keys []string, err error) {
	for k := range r.byKey {
		keys = append(keys, k)
	}
	return
}
