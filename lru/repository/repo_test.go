package repository_test

import (
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
