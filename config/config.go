package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/spf13/viper"

	"github.com/blackplayerten/IdealVisual_backend/api"
	"github.com/blackplayerten/IdealVisual_backend/session"
)

type Config struct {
	Server  *api.Config
	Session *session.Config
}

func NewConfig() (*Config, error) {
	path := flag.String("c", "", "path to config file")

	flag.Parse()

	if *path != "" {
		viper.SetConfigFile(*path)
	} else {
		viper.SetConfigName("config.yaml")
		viper.AddConfigPath("etc")
		viper.AddConfigPath(".")
	}

	viper.SetDefault("server", map[string]interface{}{
		"http": map[string]interface{}{
			"addr":    ":8080",
			"timeout": 1 * time.Minute,
		},
		"static": map[string]interface{}{
			"root":             "/var/www/ideal-visual/static",
			"keepOriginalName": false,
		},
	})
	viper.SetDefault("session", map[string]interface{}{
		"database":   0,
		"connString": "user@localhost:6379",
		"expiration": 2160 * time.Hour,
	})

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	cfg := new(Config)
	if err := viper.UnmarshalExact(cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config: %w", err)
	}

	return cfg, nil
}
