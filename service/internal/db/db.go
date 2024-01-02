package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"github.com/Nahbox/streamed-order-viewer/service/internal/config"
)

type Database struct {
	conn *sql.DB
}

func Initialize(conf *config.PgConfig) (*Database, error) {
	dsn := conf.Dsn()

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	log.Info("database connection established")

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	migrationsFilePath := fmt.Sprintf("file://%s", conf.PgMigrationsPath)

	m, err := migrate.NewWithDatabaseInstance(migrationsFilePath, "postgres", driver)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	log.Info("database migrated")

	return &Database{
		conn: conn,
	}, nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}
