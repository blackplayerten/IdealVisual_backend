package post

import (
	"github.com/blackplayerten/IdealVisual_backend/database"
)

type Service struct {
	db *database.Database
}

func New(db *database.Database) *Service {
	return &Service{
		db: db,
	}
}
