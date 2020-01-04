package api

import (
	"time"
)

type Config struct {
	HTTP struct {
		Addr    string
		Timeout time.Duration
	}
	Static struct {
		Root             string
		KeepOriginalName bool
	}
}
