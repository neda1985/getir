package cache

import (
	"errors"
	"github.com/neda1985/getir/model"
	"sync"
)

type InMemoryManager interface {
	Insert(input model.InMemoryDataIO) error
	Fetch(key string) (string, error)
}

const NotFound = "key not found"

var errorNotFound = errors.New(NotFound)

type Storage struct {
	StorageMap map[string]string
	Mutex      *sync.Mutex
}

func NewInMemoryService() *Storage {
	return &Storage{
		StorageMap: make(map[string]string),
		Mutex:      &sync.Mutex{},
	}
}

func (s *Storage) Insert(input model.InMemoryDataIO) error {
	s.Mutex.Lock()
	s.StorageMap[input.Key] = input.Value
	s.Mutex.Unlock()
	return nil
}

func (s *Storage) Fetch(key string) (string, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if val, ok := s.StorageMap[key]; ok {
		return val, nil
	}
	return "", errorNotFound
}
