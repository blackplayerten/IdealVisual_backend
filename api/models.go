package api

import (
	"github.com/blackplayerten/IdealVisual_backend/account"
	"github.com/blackplayerten/IdealVisual_backend/post"
)

//go:generate easyjson models.go

const (
	AlreadyExists = "already_exists"
)

//easyjson:json
type Errors struct {
	Errors []*FieldError `json:"errors"`
}

//easyjson:json
type FieldError struct {
	Field   string   `json:"field"`
	Reasons []string `json:"reasons"`
}

//easyjson:json
type AccWithToken struct {
	Token string `json:"token"`
	*account.Account
}

//easyjson:json
type UploadedPath struct {
	Path string `json:"path"`
}

//easyjson:json
type Posts []post.Post
