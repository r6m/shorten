package models

import "time"

// Detail represents url access details
type Detail struct {
	UserAgent string    `json:"ua"`
	CreatedAt time.Time `json:"created_at"`
}
