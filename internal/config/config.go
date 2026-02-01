package config

import (
	"time"
)

type (
	Config struct {
		Version  int      `mapstructure:"version"`
		Debug    bool     `mapstrucute:"debug"`
		Timeouts Timeouts `mapstructure:"timeouts"`
		Routes   []*Route `mapstructure:"routes"`
	}

	Route struct {
		Match    []any     `mapstructure:"match"`
		Backends []Backend `mapstructure:"backends"`
	}

	Timeouts struct {
		ReadTimeout  time.Duration `mapstructure:"readTimeout"`
		WriteTimeout time.Duration `mapstructure:"writeTimeout"`
		IdleTimeout  time.Duration `mapstructure:"idleTimeout"`
	}

	Backend struct {
		Host     string `mapstructure:"host"`
		Endpoint string `mapstructure:"endpoint"`
	}
)
