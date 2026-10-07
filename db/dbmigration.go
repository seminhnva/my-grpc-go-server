package dbmigration

import (
	"database/sql"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
	postgres "github.com/golang-migrate/migrate/v4/database/postgres"
)

func Migrate(conn *sql.DB) {
	log.Println("Database migration start ")
	driver, _ := postgres.WithInstance(conn, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		log.Println("DB migration failed:", err)
		return

	}
	if err := m.Down(); err != nil {
		log.Println("DB migration (down) failed", err)
		return

	}
	if err := m.Up(); err != nil {
		log.Println("DB migration (up) failed", err)
		return

	}
	log.Println("DB migration end")
}
