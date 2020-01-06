package api

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
