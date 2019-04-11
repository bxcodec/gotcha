package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bxcodec/gotcha/cache"
	"github.com/bxcodec/gotcha/lfu/repository"
)

/*
func TestSet(t *testing.T) {
	repo := repository.NewRepository()
	doc := &cache.Document{
		Key:        "key-2",
		Value:      "Hello World",
		StoredTime: time.Now(),
	}

	err := repo.Set(doc)
	if err != nil {
		t.Errorf("expected %v, actual %v", nil, err)
	}
}

func TestSetWithMultipleKeyExists(t *testing.T) {
	arrDoc := []*cache.Document{
		&cache.Document{
			Key:        "key-1",
			Value:      "Hello World 1",
			StoredTime: time.Now(),
		},
		&cache.Document{
			Key:        "key-2",
			Value:      "Hello World 2",
			StoredTime: time.Now(),
		},
		&cache.Document{
			Key:        "key-1",
			Value:      "Hello World 1 Modified",
			StoredTime: time.Now(),
		},
		&cache.Document{
			Key:        "key-3",
			Value:      "Hello World 3 Modified",
			StoredTime: time.Now(),
		},
		&cache.Document{
			Key:        "key-1",
			Value:      "Hello World 1 Modified Twice",
			StoredTime: time.Now(),
		},
	}
	repo := repository.NewRepository()
	for _, doc := range arrDoc {
		err := repo.Set(doc)
		if err != nil {
			t.Errorf("expected %v, actual %v", nil, err)
		}
	}
}
*/

func TestGet(t *testing.T) {
	repo := repository.NewRepository()
	arrDoc := []*cache.Document{
		// &cache.Document{
		// 	Key:        "key-1",
		// 	Value:      "Hello World 1",
		// 	StoredTime: time.Now(),
		// },
		&cache.Document{
			Key:        "key-2",
			Value:      "Hello World 2",
			StoredTime: time.Now(),
		},
		// &cache.Document{
		// 	Key:        "key-1",
		// 	Value:      "Hello World 1 Modified",
		// 	StoredTime: time.Now(),
		// },
		// &cache.Document{
		// 	Key:        "key-3",
		// 	Value:      "Hello World 3 Modified",
		// 	StoredTime: time.Now(),
		// },
		// &cache.Document{
		// 	Key:        "key-1",
		// 	Value:      "Hello World 1 Modified Twice",
		// 	StoredTime: time.Now(),
		// },
	}

	for _, doc := range arrDoc {
		err := repo.Set(doc)
		if err != nil {
			t.Errorf("expected %v, actual %v", nil, err)
		}
	}

	res, err := repo.Get("key-2")
	if err != nil {
		t.Errorf("expected %v, actual %v", nil, err)
	}

	fmt.Println("Res: ", res.Value)

	res2, err := repo.Get("key-2")
	if err != nil {
		t.Errorf("expected %v, actual %v", nil, err)
	}

	fmt.Println("Res: ", res2.Value)
}
