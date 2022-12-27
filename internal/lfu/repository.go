package lfu

import (
	"container/list"
	"encoding/json"
	"reflect"
	"time"

	"github.com/bxcodec/gotcha/cache"
)

// Repository represent the data repository for inernal cache
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

// New will initialize the LFU memory cache
func New(maxSize, maxMemory uint64, expiryTreshold time.Duration) (repo *Repository) {
	repo = &Repository{
		frequencyList:  list.New(),
		byKey:          make(map[string]*lfuItem),
		maxMemory:      maxMemory,
		maxSize:        maxSize,
		expiryTreshold: expiryTreshold,
	}
	return
}

// Get will retrieve the item from cache
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
		_, _ = r.Delete(key)
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

	return res, nil
}

// Set wil save the item to cache
func (r *Repository) Set(doc *cache.Document) (err error) {
	if _, ok := r.byKey[doc.Key]; ok {
		// TODO: (bxcodec)
		// Re-insert the document
		return
	}

	freq := r.frequencyList.Front() // Front will always be the least frequently used
	if freq == nil {
		newNodeFreq := &frequencyItem{
			Frequency: 1,
		}
		freq = r.frequencyList.PushFront(newNodeFreq)
		freqVal := freq.Value.(*frequencyItem)
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

	freqVal := freq.Value.(*frequencyItem)
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
	// Remove oldest if the max-size reached
	if uint64(len(r.byKey)) > r.maxSize {
		r.removeLfuOldest()
	}

	// Avoid memory limit if set zero to increase performances
	if r.maxMemory == 0 {
		return
	}

	byteMap, err := json.Marshal(r.byKey)
	if err != nil {
		_, _ = r.Delete(doc.Key)
		return
	}
	// Remove oldest if the maxmemory reached
	if uint64(len(byteMap)) > r.maxMemory {
		r.removeLfuOldest()
	}
	return nil
}

func (r *Repository) removeLfuOldest() {
	lfuList := r.frequencyList.Front()
	if r.frequencyList.Len() == 0 {
		return
	}
	freqItem := lfuList.Value.(*frequencyItem)

	minStoreTime := time.Now().Unix()
	var oldestItem *lfuItem
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
}

// Clear will clear up the item from cache
func (r *Repository) Clear() (err error) {
	for k := range r.byKey {
		delete(r.byKey, k)
	}
	r.frequencyList.Init()
	return
}

// Len return the total items in the cache
func (r *Repository) Len() int {
	return len(r.byKey)
}

// Contains check if any item with the given key exist in the cache
func (r *Repository) Contains(key string) (ok bool) {
	_, ok = r.byKey[key]
	return
}

// Delete will delete the item from cache
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

// Keys return all keys from cache
func (r *Repository) Keys() (keys []string, err error) {
	for k := range r.byKey {
		keys = append(keys, k)
	}
	return
}
