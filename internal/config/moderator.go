package config

import "time"

type Moderator struct {
	CheckInterval time.Duration `env:"MODER_CHECK_INTERVAL" envDefault:"30s"`
}
