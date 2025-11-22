package config

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func RunMigrations(dbUrl string) error {
	migration, err := migrate.New("file://migrations", dbUrl)
	if err != nil {
		return err
	}
	defer func() {
		if err, _ = migration.Close(); err != nil {
			log.Printf("Error closing migration: %s\n", err)
		}
	}()

	version, dirty, err := migration.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return err
	}

	if dirty {
		log.Printf("Fixing dirty database version %d", version)
		if err := migration.Force(int(version)); err != nil {
			return fmt.Errorf("force version failed: %w", err)
		}
	}

	if err := migration.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Database is up to date")
			return nil
		}
		return fmt.Errorf("migration up failed: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}

func CreateDBConnection(dbUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}
	return db, nil
}
