package config

import "fmt"

type PgConfig struct {
	PgHost           string `envconfig:"POSTGRES_HOST" required:"true"`
	PgUser           string `envconfig:"POSTGRES_USER" required:"true"`
	PgPassword       string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	PgDB             string `envconfig:"POSTGRES_DB" required:"true"`
	PgPort           int    `envconfig:"POSTGRES_PORT" required:"true"`
	PgMigrationsPath string `envconfig:"POSTGRES_MIGRATIONS_PATH" required:"true"`
}

func (p *PgConfig) Dsn() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.PgHost, p.PgPort, p.PgUser, p.PgPassword, p.PgDB)
}
