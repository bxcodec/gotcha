package repository_test

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/bxcodec/gotcha/cache"
	"github.com/bxcodec/gotcha/lfu/repository"
)

func TestSet(t *testing.T) {
	repo := repository.New(5, 10*cache.MB, time.Second*5)
	doc := &cache.Document{
		Key:        "key-2",
		Value:      "Hello World",
		StoredTime: time.Now().Add(time.Second * -1).Unix(),
	}

	err := repo.Set(doc)
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}
}

func TestSetWithMultipleKeyExists(t *testing.T) {
	arrDoc := []*cache.Document{
		&cache.Document{
			Key:        "key-1",
			Value:      "Hello World 1",
			StoredTime: time.Now().Add(time.Second * -1).Unix(),
		},
		&cache.Document{
			Key:        "key-2",
			Value:      "Hello World 2",
			StoredTime: time.Now().Add(time.Second * -1).Unix(),
		},
		&cache.Document{
			Key:        "key-1",
			Value:      "Hello World 1 Modified",
			StoredTime: time.Now().Add(time.Second * -1).Unix(),
		},
		&cache.Document{
			Key:        "key-3",
			Value:      "Hello World 3 Modified",
			StoredTime: time.Now().Add(time.Second * -1).Unix(),
		},
		&cache.Document{
			Key:        "key-1",
			Value:      "Hello World 1 Modified Twice",
			StoredTime: time.Now().Add(time.Second * -1).Unix(),
		},
	}
	repo := repository.New(5, 10*cache.MB, time.Second*5)
	for _, doc := range arrDoc {
		err := repo.Set(doc)
		if err != nil {
			t.Fatalf("expected %v, actual %v", nil, err)
		}
	}
}

func TestGetOne(t *testing.T) {
	repo := repository.New(5, 100, time.Second*5)
	doc := &cache.Document{
		Key:        "key-2",
		Value:      "Hello World",
		StoredTime: time.Now().Add(time.Second * -1).Unix(),
	}

	err := repo.Set(doc)
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	res, err := repo.Get("key-2")
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	if res.Value != doc.Value {
		t.Fatalf("expected %v, actual %v", doc.Value, res.Value)
	}
}

func TestGetWithMultipleSet(t *testing.T) {
	repo := repository.New(4, 500, time.Minute*5)
	arrDoc := []*cache.Document{
		&cache.Document{
			Key:        "key-1",
			Value:      "A",
			StoredTime: time.Now().Add(time.Minute * -1).Unix(),
		},
		&cache.Document{
			Key:        "key-1",
			Value:      "A'",
			StoredTime: time.Now().Add(time.Second * -40).Unix(),
		},
		&cache.Document{
			Key:        "key-3",
			Value:      "C",
			StoredTime: time.Now().Add(time.Second * -1).Unix(),
		},
		&cache.Document{
			Key:        "key-1",
			Value:      "A''",
			StoredTime: time.Now().Add(time.Second * -10).Unix(),
		},
		&cache.Document{
			Key:        "key-4",
			Value:      "D",
			StoredTime: time.Now().Add(time.Second * -5).Unix(),
		},
	}

	for _, doc := range arrDoc {
		err := repo.Set(doc)
		if err != nil {
			t.Fatalf("expected %v, actual %v", nil, err)
		}
	}

	doc2 := &cache.Document{
		Key:        "key-2",
		Value:      "B",
		StoredTime: time.Now().Add(time.Second * -2).Unix(),
	}

	err := repo.Set(doc2)
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	res, err := repo.Get("key-2")
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	if res.Value != doc2.Value {
		t.Fatalf("expected %v, actual %v", doc2.Value, res.Value)
	}

	res3, err := repo.Get("key-3")
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	if res3.Value != arrDoc[2].Value {
		t.Fatalf("expected %v, actual %v", arrDoc[2].Value, res3.Value)
	}

	res2, err := repo.Get("key-2")
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	if res2.Value != doc2.Value {
		t.Fatalf("expected %v, actual %v", doc2.Value, res2.Value)
	}

	docLast := &cache.Document{
		Key:        "key-5",
		Value:      "E",
		StoredTime: time.Now().Unix(),
	}

	err = repo.Set(docLast)
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	res5, err := repo.Get("key-5")
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	if res5.Value != docLast.Value {
		t.Fatalf("expected %v, actual %v", docLast.Value, res5.Value)
	}
}

func TestGetExpiredItem(t *testing.T) {
	repo := repository.New(4, 100, time.Second*5)
	arrDoc := []*cache.Document{
		&cache.Document{
			Key:        "key-1",
			Value:      "A",
			StoredTime: time.Now().Add(time.Minute * -1).Unix(),
		},
		&cache.Document{
			Key:        "key-3",
			Value:      "C",
			StoredTime: time.Now().Add(time.Second * -1).Unix(),
		},
		&cache.Document{
			Key:        "key-4",
			Value:      "D",
			StoredTime: time.Now().Add(time.Second * -3).Unix(),
		},
	}

	for _, doc := range arrDoc {
		err := repo.Set(doc)
		if err != nil {
			t.Fatalf("expected %v, actual %v", nil, err)
		}
	}

	// Ensure the current length == 3
	if repo.Len() != 3 {
		t.Fatalf("expected %v, actual %v", 3, repo.Len())
	}

	// Get the expired
	res, err := repo.Get("key-1")
	if err == nil {
		t.Fatalf("expected %v, actual %v", "error", err)
	}
	if res != nil {
		t.Fatalf("expected %v, actual %v", nil, res)
	}

	// Get the non expired
	res, err = repo.Get("key-4")
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	if res.Value != arrDoc[2].Value {
		t.Fatalf("expected %v, actual %v", arrDoc[2].Value, res.Value)
	}

	// Get total length of exist item should be == 2
	if repo.Len() != 2 {
		t.Fatalf("expected %v, actual %v", 2, repo.Len())
	}
}

func TestGetExpiredButSingleSetItemInList(t *testing.T) {
	repo := repository.New(4, 100, time.Second*5)
	arrDoc := []*cache.Document{
		&cache.Document{
			Key:        "key-1",
			Value:      "A",
			StoredTime: time.Now().Add(time.Minute * -1).Unix(),
		},
		&cache.Document{
			Key:        "key-4",
			Value:      "D",
			StoredTime: time.Now().Add(time.Second * -3).Unix(),
		},
	}

	for _, doc := range arrDoc {
		err := repo.Set(doc)
		if err != nil {
			t.Fatalf("expected %v, actual %v", nil, err)
		}
	}

	// Get the non expired first to increase the frequency,
	// and let the expired alone in the set of its parent list
	res, err := repo.Get("key-4")
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	if res.Value != arrDoc[1].Value {
		t.Fatalf("expected %v, actual %v", arrDoc[1].Value, res.Value)
	}

	// Ensure the current length == 2
	if repo.Len() != 2 {
		t.Fatalf("expected %v, actual %v", 2, repo.Len())
	}

	// Get the expired
	res, err = repo.Get("key-1")
	if err == nil {
		t.Fatalf("expected %v, actual %v", "error", err)
	}
	if res != nil {
		t.Fatalf("expected %v, actual %v", nil, res)
	}

	// Get total length of exist item should be == 2
	if repo.Len() != 1 {
		t.Fatalf("expected %v, actual %v", 2, repo.Len())
	}
}

func TestSetWithFrequency1IsNotExists(t *testing.T) {
	repo := repository.New(5, 100, time.Second*5)
	doc := &cache.Document{
		Key:        "key-2",
		Value:      "Hello World",
		StoredTime: time.Now().Add(time.Second * -1).Unix(),
	}

	err := repo.Set(doc)
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	res, err := repo.Get("key-2")
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	if res.Value != doc.Value {
		t.Fatalf("expected %v, actual %v", doc.Value, res.Value)
	}

	docLast := &cache.Document{
		Key:        "key-5",
		Value:      "E",
		StoredTime: time.Now().Unix(),
	}

	err = repo.Set(docLast)
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

}

func TestClearCache(t *testing.T) {
	repo := repository.New(4, 100, time.Second*5)
	arrDoc := []*cache.Document{
		&cache.Document{
			Key:        "key-1",
			Value:      "A",
			StoredTime: time.Now().Add(time.Minute * -1).Unix(),
		},
		&cache.Document{
			Key:        "key-1",
			Value:      "A'",
			StoredTime: time.Now().Add(time.Second * -40).Unix(),
		},
		&cache.Document{
			Key:        "key-3",
			Value:      "C",
			StoredTime: time.Now().Add(time.Second * -30).Unix(),
		},
		&cache.Document{
			Key:        "key-1",
			Value:      "A''",
			StoredTime: time.Now().Add(time.Second * -10).Unix(),
		},
		&cache.Document{
			Key:        "key-4",
			Value:      "D",
			StoredTime: time.Now().Add(time.Second * -5).Unix(),
		},
	}

	for _, doc := range arrDoc {
		err := repo.Set(doc)
		if err != nil {
			t.Fatalf("expected %v, actual %v", nil, err)
		}
	}

	if repo.Len() != 3 {
		t.Fatalf("expect %v got %v", 3, repo.Len())
	}

	err := repo.Clear()
	if err != nil {
		t.Fatalf("expect %v got %v", nil, err)
	}

	if repo.Len() != 0 {
		t.Fatalf("expect %v got %v", 0, repo.Len())
	}
}

func TestContains(t *testing.T) {
	repo := repository.New(4, 100, time.Second*5)
	arrDoc := []*cache.Document{
		&cache.Document{
			Key:        "key-1",
			Value:      "A",
			StoredTime: time.Now().Add(time.Minute * -1).Unix(),
		},
		&cache.Document{
			Key:        "key-1",
			Value:      "A'",
			StoredTime: time.Now().Add(time.Second * -40).Unix(),
		},
		&cache.Document{
			Key:        "key-3",
			Value:      "C",
			StoredTime: time.Now().Add(time.Second * -30).Unix(),
		},
		&cache.Document{
			Key:        "key-1",
			Value:      "A''",
			StoredTime: time.Now().Add(time.Second * -10).Unix(),
		},
		&cache.Document{
			Key:        "key-4",
			Value:      "D",
			StoredTime: time.Now().Add(time.Second * -5).Unix(),
		},
	}

	for _, doc := range arrDoc {
		err := repo.Set(doc)
		if err != nil {
			t.Fatalf("expected %v, actual %v", nil, err)
		}
	}

	if !repo.Contains("key-3") {
		t.Fatalf("expected %v, actual %v", true, repo.Contains("key-3"))
	}
}

func TestDelete(t *testing.T) {
	repo := repository.New(4, 100, time.Second*5)
	arrDoc := []*cache.Document{
		&cache.Document{
			Key:        "key-1",
			Value:      "A",
			StoredTime: time.Now().Add(time.Minute * -1).Unix(),
		},
		&cache.Document{
			Key:        "key-1",
			Value:      "A'",
			StoredTime: time.Now().Add(time.Second * -40).Unix(),
		},
		&cache.Document{
			Key:        "key-3",
			Value:      "C",
			StoredTime: time.Now().Add(time.Second * -30).Unix(),
		},
		&cache.Document{
			Key:        "key-1",
			Value:      "A''",
			StoredTime: time.Now().Add(time.Second * -10).Unix(),
		},
		&cache.Document{
			Key:        "key-4",
			Value:      "D",
			StoredTime: time.Now().Add(time.Second * -5).Unix(),
		},
	}

	for _, doc := range arrDoc {
		err := repo.Set(doc)
		if err != nil {
			t.Fatalf("expected %v, actual %v", nil, err)
		}
	}

	if !repo.Contains("key-3") {
		t.Fatalf("expected %v, actual %v", true, repo.Contains("key-3"))
	}

	ok, err := repo.Delete("key-3")
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	if !ok {
		t.Fatalf("expected %v, actual %v", true, ok)
	}

	if repo.Contains("key-3") {
		t.Fatalf("expected %v, actual %v", false, repo.Contains("key-3"))
	}
}

func TestGetKeys(t *testing.T) {
	repo := repository.New(4, 100, time.Second*5)
	arrDoc := []*cache.Document{
		&cache.Document{
			Key:        "key-1",
			Value:      "A",
			StoredTime: time.Now().Add(time.Minute * -1).Unix(),
		},
		&cache.Document{
			Key:        "key-3",
			Value:      "C",
			StoredTime: time.Now().Add(time.Second * -30).Unix(),
		},
		&cache.Document{
			Key:        "key-4",
			Value:      "D",
			StoredTime: time.Now().Add(time.Second * -5).Unix(),
		},
	}

	for _, doc := range arrDoc {
		err := repo.Set(doc)
		if err != nil {
			t.Fatalf("expected %v, actual %v", nil, err)
		}
	}

	keys, err := repo.Keys()
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	var contains = func(kesy []string, item string) (ok bool) {
		for _, k := range keys {
			ok = item == k
			if ok {
				return
			}
		}
		return
	}

	expectedKeys := []string{"key-1", "key-3", "key-4"}
	for _, k := range keys {
		if !contains(expectedKeys, k) {
			t.Fatalf("expected %v, actual %v", true, contains(expectedKeys, k))
		}
	}
}

// This benchmark code below also used for profiling to get the memory and CPU usage
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func BenchmarkSetItem(b *testing.B) {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	repo := repository.New(1000, 100, time.Minute*40)
	preDoc := &cache.Document{
		Key:        "key-1",
		Value:      "Hello World",
		StoredTime: time.Now().Add(time.Second * -1).Unix(),
	}
	err := repo.Set(preDoc)
	if err != nil {
		b.Fatalf("expected %v, actual %v", nil, err)
	}
	doc := &cache.Document{
		Key:        "key-2",
		Value:      `Hello World`,
		StoredTime: time.Now().Add(time.Second * -1).Unix(),
	}

	counterMiss := 0
	counterHit := 0
	b.N = 100
	for i := 0; i < b.N; i++ {
		temp := *doc
		temp.Key = fmt.Sprintf("key-%d", i)
		err := repo.Set(&temp)
		if err != nil {
			b.Fatalf("expected %v, actual %v", nil, err)
		}

		randVal := 1
		if i != 0 {
			randVal = rand.Intn(i)
		}
		res, err := repo.Get(fmt.Sprintf("key-%d", randVal))
		if res == nil || err != nil {
			counterMiss++
		} else {
			counterHit++
		}
	}
	fmt.Println("Counter hit: ", counterHit)
	fmt.Println("Counter miss: ", counterMiss)
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
