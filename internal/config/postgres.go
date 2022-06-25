package config

type Postgres struct {
	User         string `env:"PG_COMM_USER"`
	Password     string `env:"PG_COMM_PASSWORD"`
	Host         string `env:"PG_COMM_HOST"`
	Port         int    `env:"PG_COMM_PORT"`
	Database     string `env:"PG_COMM_DATABASE"`
	TestDatabase string `env:"PG_COMM_TEST_DATABASE"`
}
