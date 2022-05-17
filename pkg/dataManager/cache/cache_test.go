package cache_test

import (
	"github.com/neda1985/getir/model"
	"github.com/neda1985/getir/pkg/dataManager/cache"
	"testing"
)

var (
	testInMemoryInput = model.InMemoryDataIO{
		Key:   "test-key",
		Value: "test-value",
	}
	inMemoryService = cache.NewInMemoryService()
)

func TestHolder_Fetch(t *testing.T) {
	inMemoryService.Insert(testInMemoryInput)
	v, e := inMemoryService.Fetch("test-key")
	if e != nil {
		t.Fail()
	}
	if v != "test-value" {
		t.Fail()
	}
}

func TestHolder_Insert(t *testing.T) {
	inMemoryService.Insert(testInMemoryInput)
	if k, ok := inMemoryService.StorageMap["test-key"]; !ok {
		t.Fail()
	} else if k != "test-value" {
		t.Fail()
	}
}
