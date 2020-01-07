package post

import (
	"time"
)

//go:generate easyjson models.go

//easyjson:json
type Post struct {
	ID         string     `json:"id"`
	Acc        uint64     `json:"-"`
	Photo      string     `json:"photo"`
	PhotoIndex *int       `json:"photo_index,omitempty" db:"photo_index"`
	Date       *time.Time `json:"date,omitempty"`
	Place      *string    `json:"place,omitempty"`
	Text       *string    `json:"text,omitempty"`
}
