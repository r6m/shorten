package models

import "github.com/kaibox-git/randstring"

// URL model
type URL struct {
	Key         string    `json:"key"`
	OriginalURL string    `json:"original_url"`
	Details     []*Detail `json:"details,omityempty"`
}

// GenerateKey sets a random 6 characters string as url key
func (u *URL) GenerateKey() {
	u.Key = randstring.Create(6)
}
