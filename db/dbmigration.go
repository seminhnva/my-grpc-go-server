package dbmigration

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
	postgres "github.com/golang-migrate/migrate/v4/database/postgres"
)

func Migrate(conn *sql.DB) error {
	log.Println("Database migration start ")
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create postgres migration driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("initialize database migrations: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("apply database migrations: %w", err)
		}
	}
	log.Println("DB migration end")
	return nil
}
