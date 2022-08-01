package memory

import (
	"sync"

	"github.com/r6m/shorten/models"
	"github.com/r6m/shorten/store"
	"github.com/sirupsen/logrus"
)

type memoryStore struct {
	urls sync.Map
}

func NewStore() *memoryStore {
	logrus.Info("setting up memory storage")
	return &memoryStore{
		urls: sync.Map{},
	}
}

func (s *memoryStore) Save(url *models.URL) error {
	_, loaded := s.urls.LoadOrStore(url.Key, url)

	if loaded {
		return store.ErrDuplicate
	}

	return nil
}

func (s *memoryStore) Load(key string, detail *models.Detail) (*models.URL, error) {
	v, ok := s.urls.Load(key)
	if ok {
		url := v.(*models.URL)
		if detail != nil {
			url.Details = append(url.Details, detail)
			s.urls.Store(key, url)
		}
		return url, nil
	}

	return nil, store.ErrNotFound
}
