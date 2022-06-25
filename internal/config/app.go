package config

type Application struct {
	LogLevel  string `env:"LOG_LEVEL" envDefault:"INFO"`
	Postgres  Postgres
	Server    Server
	Moderator Moderator
}
