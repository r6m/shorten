package store

import (
	"errors"

	"github.com/r6m/shorten/models"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("already exists")
)

type Store interface {
	Save(*models.URL) error
	Load(string, *models.Detail) (*models.URL, error)
}
