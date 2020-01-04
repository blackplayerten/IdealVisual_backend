package session

import (
	"time"
)

type Config struct {
	Database   int
	ConnString string

	Expiration time.Duration
}
