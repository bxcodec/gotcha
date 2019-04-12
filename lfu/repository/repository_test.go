package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bxcodec/gotcha/cache"
	"github.com/bxcodec/gotcha/lfu/repository"
)

func TestSet(t *testing.T) {
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
	repo := repository.New(5, 100, time.Second*5)
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

	// _, err = repo.Get("key-3")
	// if err == nil {
	// 	t.Fatalf("expected %v, actual %v", nil, err)
	// }

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

func BenchmarkSetItem(b *testing.B) {
	repo := repository.New(1000, 100, time.Second*40)
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
		Value:      "Hello World",
		StoredTime: time.Now().Add(time.Second * -1).Unix(),
	}
	for i := 0; i < b.N; i++ {
		temp := *doc
		temp.Key = fmt.Sprintf("key-%d", i)
		err := repo.Set(&temp)
		if err != nil {
			b.Fatalf("expected %v, actual %v", nil, err)
		}
	}
}
