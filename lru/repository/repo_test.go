package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bxcodec/gotcha/cache"
	"github.com/bxcodec/gotcha/lru/repository"
)

func TestSet(t *testing.T) {
	repo := repository.New(10, 100, time.Minute*5)
	doc := &cache.Document{
		Key:        "key-2",
		Value:      "Hello World",
		StoredTime: time.Now().Unix(),
	}
	err := repo.Set(doc)
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	// Check if the item is exists
	item, err := repo.Peek("key-2")

	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	if item == nil {
		t.Fatalf("expected %v, actual %v", "Hello World", err)
	}
}

func TestSetMultiple(t *testing.T) {
	repo := repository.New(5, 100, time.Minute*5)
	for i := 1; i <= 10; i++ {
		doc := &cache.Document{
			Key:        fmt.Sprintf("key:%d", i),
			Value:      i,
			StoredTime: time.Now().Unix(),
		}
		err := repo.Set(doc)
		if err != nil {
			t.Fatalf("expected %v, actual %v", nil, err)
		}
	}

	// Since the size is 5, so the first 5 should be not exists
	// Assert the key:1 - key:5 is not exists
	for i := 1; i <= 5; i++ {
		item, err := repo.Peek(fmt.Sprintf("key:%d", i))
		if err == nil {
			t.Fatalf("expected %v, actual %v", cache.ErrMissed, err)
		}
		if item != nil {
			t.Fatalf("expected %v, actual %v", nil, item)
		}
	}

	// Assert the key:6 - key:10 is exists
	for i := 6; i <= 10; i++ {
		item, err := repo.Peek(fmt.Sprintf("key:%d", i))
		if err != nil {
			t.Fatalf("expected %v, actual %v", nil, err)
		}

		if item == nil {
			t.Fatalf("expected %v, actual %v", i, err)
		}
	}
}

func TestSetWithExistingKey(t *testing.T) {
	arrDoc := []*cache.Document{
		&cache.Document{
			Key:        "key-1",
			Value:      "Hello World 1",
			StoredTime: time.Now().Unix(),
		},
		&cache.Document{
			Key:        "key-2",
			Value:      "Hello World 2",
			StoredTime: time.Now().Unix(),
		},
		&cache.Document{
			Key:        "key-1",
			Value:      "Hello World 1 Modified",
			StoredTime: time.Now().Unix(),
		},
		&cache.Document{
			Key:        "key-3",
			Value:      "Hello World 3 Modified",
			StoredTime: time.Now().Unix(),
		},
		&cache.Document{
			Key:        "key-1",
			Value:      "Hello World 1 Modified Twice",
			StoredTime: time.Now().Unix(),
		},
	}

	repo := repository.New(10, 100, time.Minute*5)

	for _, doc := range arrDoc {
		err := repo.Set(doc)
		if err != nil {
			t.Fatalf("expected %v, actual %v", nil, err)
		}
	}

	len := repo.Len()
	// Since the key is only 3 are different even the item to be set are 5
	if len != 3 {
		t.Fatalf("expected %v, actual %v", 3, len)
	}
}

func TestGet(t *testing.T) {
	repo := repository.New(10, 100, time.Minute*5)
	doc := &cache.Document{
		Key:        "key-2",
		Value:      "Hello World",
		StoredTime: time.Now().Unix(),
	}
	err := repo.Set(doc)
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	// Check if the item is exists
	item, err := repo.Get("key-2")

	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}

	if item == nil {
		t.Fatalf("expected %v, actual %v", "Hello World", err)
	}
}

func TestGetOldest(t *testing.T) {
	arrDoc := []*cache.Document{
		&cache.Document{
			Key:        "key-1",
			Value:      "Hello World 1",
			StoredTime: time.Now().Unix(),
		},
		&cache.Document{
			Key:        "key-2",
			Value:      "Hello World 2",
			StoredTime: time.Now().Unix(),
		},
		&cache.Document{
			Key:        "key-1",
			Value:      "Hello World 1 Modified",
			StoredTime: time.Now().Unix(),
		},
		&cache.Document{
			Key:        "key-3",
			Value:      "Hello World 3 Modified",
			StoredTime: time.Now().Unix(),
		},
		&cache.Document{
			Key:        "key-1",
			Value:      "Hello World 1 Modified Twice",
			StoredTime: time.Now().Unix(),
		},
	}

	repo := repository.New(10, 100, time.Minute*5)

	for _, doc := range arrDoc {
		err := repo.Set(doc)
		if err != nil {
			t.Fatalf("expected %v, actual %v", nil, err)
		}
	}

	res, err := repo.Get("key-3")
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}
	if res.Value != arrDoc[3].Value {
		t.Fatalf("expected %v, actual %v", arrDoc[3].Value, res.Value)
	}

	doc, err := repo.GetOldest()
	if err != nil {
		t.Fatalf("expected %v, actual %v", nil, err)
	}
	if doc.Value != arrDoc[1].Value {
		t.Fatalf("expected %v, actual %v", arrDoc[1].Value, doc.Value)
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

func BenchmarkSetItem(b *testing.B) {
	repo := repository.New(10, 100, time.Minute*5)
	preDoc := &cache.Document{
		Key:        "key-1",
		Value:      "Hello World",
		StoredTime: time.Now().Unix(),
	}
	err := repo.Set(preDoc)
	if err != nil {
		b.Fatalf("expected %v, actual %v", nil, err)
	}

	doc := &cache.Document{
		Key:        "key-2",
		Value:      "Hello World",
		StoredTime: time.Now().Unix(),
	}

	for i := 0; i < b.N; i++ {
		tmp := *doc
		tmp.Key = fmt.Sprintf("key-%d", i)
		err := repo.Set(&tmp)
		if err != nil {
			b.Fatalf("expected %v, actual %v", nil, err)
		}
	}
}
