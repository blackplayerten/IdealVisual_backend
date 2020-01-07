package account

//go:generate easyjson models.go

//easyjson:json
type Credentials struct {
	// ID       uint64 `json:"-"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//easyjson:json
type Account struct {
	ID       uint64  `json:"id"`
	Email    string  `json:"email"`
	Username string  `json:"username"`
	Avatar   *string `json:"avatar,omitempty"`
}

//easyjson:json
type FullAccount struct {
	Account
	Password string `json:"password"`
}
