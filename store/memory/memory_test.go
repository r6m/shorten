package memory

import (
	"testing"

	"github.com/r6m/shorten/models"
	"github.com/r6m/shorten/store"
)

func Test_memoryStore_Load(t *testing.T) {
	s := NewStore()

	// preload urls
	url1 := &models.URL{Key: "111111", OriginalURL: "https://google.com", Details: []*models.Detail{}}
	url2 := &models.URL{Key: "222222", OriginalURL: "https://google.com", Details: []*models.Detail{}}
	s.Save(url1)
	s.Save(url2)

	type args struct {
		key    string
		detail *models.Detail
	}
	type result struct {
		url *models.URL
		err error
	}

	tests := []struct {
		name    string
		args    args
		result  result
		wantErr bool
	}{
		{
			name: "load valid key",
			args: args{
				key:    url1.Key,
				detail: &models.Detail{UserAgent: "Bot"},
			},
			result: result{
				url: url1,
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "load invalid key",
			args: args{
				key:    "InV4l!d",
				detail: &models.Detail{UserAgent: "Bot"},
			},
			result: result{
				url: nil,
				err: store.ErrNotFound,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := s.Load(tt.args.key, tt.args.detail)
			if err != tt.result.err {
				t.Errorf("memoryStore.Load() error = %v, wantErr %v", err, tt.result.err)
			}

			if url != nil && url.Key != tt.result.url.Key {
				t.Errorf("url key is %s but exptect to be %s", url.Key, tt.result.url.Key)
			}
		})
	}
}

func Test_memoryStore_Save(t *testing.T) {
	s := NewStore()
	tests := []struct {
		name string
		url  *models.URL
		err  error
	}{
		{
			name: "add valid url",
			url: &models.URL{
				Key:         "111111",
				OriginalURL: "https://google.com",
				Details:     make([]*models.Detail, 0),
			},
			err: nil,
		},
		{
			name: "add duplicate url",
			url: &models.URL{
				Key:         "111111",
				OriginalURL: "https://google.com",
				Details:     make([]*models.Detail, 0),
			},
			err: store.ErrDuplicate,
		},
		{
			name: "add second valid item",
			url: &models.URL{
				Key:         "222222",
				OriginalURL: "https://google.com",
				Details:     make([]*models.Detail, 0),
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Save(tt.url); err != tt.err {
				t.Errorf("memoryStore.Save() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
