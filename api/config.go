package api

import (
	"time"
)

type Config struct {
	HTTP struct {
		Addr        string
		Timeout     time.Duration
		BodyLimitMB int
	}
	Static struct {
		Root             string
		KeepOriginalName bool
	}
}
