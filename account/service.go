package account

import (
	"github.com/blackplayerten/IdealVisual_backend/database"
)

type Service struct {
	db *database.Database
}

func New(dataSource *database.Database) *Service {
	return &Service{
		db: dataSource,
	}
}
