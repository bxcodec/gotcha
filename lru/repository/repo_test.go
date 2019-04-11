package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bxcodec/gotcha"
	"github.com/bxcodec/gotcha/lru/repository"
)

func TestSet(t *testing.T) {
	repo := repository.New(10, 100)
	doc := &gotcha.Document{
		Key:        "key-2",
		Value:      "Hello World",
		StoredTime: time.Now(),
	}
	err := repo.Set(doc)
	if err != nil {
		t.Errorf("expected %v, actual %v", nil, err)
	}

	// Check if the item is exists
	item, err := repo.Peek("key-2")

	if err != nil {
		t.Errorf("expected %v, actual %v", nil, err)
	}

	if item == nil {
		t.Errorf("expected %v, actual %v", "Hello World", err)
	}
}

func TestSetMultiple(t *testing.T) {
	repo := repository.New(5, 100)
	for i := 1; i <= 10; i++ {
		doc := &gotcha.Document{
			Key:        fmt.Sprintf("key:%d", i),
			Value:      i,
			StoredTime: time.Now(),
		}
		err := repo.Set(doc)
		if err != nil {
			t.Errorf("expected %v, actual %v", nil, err)
		}
	}

	// Since the size is 5, so the first 5 should be not exists
	// Assert the key:1 - key:5 is not exists
	for i := 1; i <= 5; i++ {
		item, err := repo.Peek(fmt.Sprintf("key:%d", i))
		if err == nil {
			t.Errorf("expected %v, actual %v", gotcha.ErrCacheMissed, err)
		}
		if item != nil {
			t.Errorf("expected %v, actual %v", nil, item)
		}
	}

	// Assert the key:6 - key:10 is exists
	for i := 6; i <= 10; i++ {
		item, err := repo.Peek(fmt.Sprintf("key:%d", i))
		if err != nil {
			t.Errorf("expected %v, actual %v", nil, err)
		}

		if item == nil {
			t.Errorf("expected %v, actual %v", "Hello World", err)
		}
	}
}

func TestSetWithExistingKey(t *testing.T) {
	arrDoc := []*gotcha.Document{
		&gotcha.Document{
			Key:        "key-1",
			Value:      "Hello World 1",
			StoredTime: time.Now(),
		},
		&gotcha.Document{
			Key:        "key-2",
			Value:      "Hello World 2",
			StoredTime: time.Now(),
		},
		&gotcha.Document{
			Key:        "key-1",
			Value:      "Hello World 1 Modified",
			StoredTime: time.Now(),
		},
		&gotcha.Document{
			Key:        "key-3",
			Value:      "Hello World 3 Modified",
			StoredTime: time.Now(),
		},
		&gotcha.Document{
			Key:        "key-1",
			Value:      "Hello World 1 Modified Twice",
			StoredTime: time.Now(),
		},
	}

	repo := repository.New(5, 100)

	for _, doc := range arrDoc {
		err := repo.Set(doc)
		if err != nil {
			t.Errorf("expected %v, actual %v", nil, err)
		}
	}

	len, err := repo.Len()
	if err != nil {
		t.Errorf("expected %v, actual %v", nil, err)
	}

	// Since the key is only 3 are different even the item to be set are 5
	if len != 3 {
		t.Errorf("expected %v, actual %v", 3, len)
	}
}

func TestGet(t *testing.T) {
	repo := repository.New(10, 100)
	doc := &gotcha.Document{
		Key:        "key-2",
		Value:      "Hello World",
		StoredTime: time.Now(),
	}
	err := repo.Set(doc)
	if err != nil {
		t.Errorf("expected %v, actual %v", nil, err)
	}

	// Check if the item is exists
	item, err := repo.Get("key-2")

	if err != nil {
		t.Errorf("expected %v, actual %v", nil, err)
	}

	if item == nil {
		t.Errorf("expected %v, actual %v", "Hello World", err)
	}
}
