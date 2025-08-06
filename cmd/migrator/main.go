package main

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is required")
	}

	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	absStoragePath, err := filepath.Abs(storagePath)
	if err != nil {
		panic("failed to get absolute path for storage: " + err.Error())
	}

	dbURL := fmt.Sprintf("sqlite3://file:%s?_fk=true&x-migrations-table=%s",
		absStoragePath, migrationsTable)

	if runtime.GOOS == "windows" {
		dbURL = fmt.Sprintf("sqlite3://file:/%s?_fk=true&x-migrations-table=%s",
			filepath.ToSlash(absStoragePath), migrationsTable)
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL,
	)
	if err != nil {
		panic("failed to create migrate instance: " + err.Error())
	}

	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}

		panic("failed to apply migrations: " + err.Error())
	}

	fmt.Println("migrations applied successfully!")
}
