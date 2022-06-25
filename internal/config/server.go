package config

import "time"

type Server struct {
	Listen       string        `env:"API_COMM_LISTEN" envDefault:":8080"`
	ReadTimeout  time.Duration `env:"API_COMM_READ_TIMEOUT" envDefault:"30s"`
	WriteTimeout time.Duration `env:"API_COMM_WRITE_TIMEOUT" envDefault:"30s"`
}
