package database

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"strings"

	_ "github.com/tursodatabase/go-libsql"
)

const (
	FILE_UP_MIGRATION_INDEX = 1

	dbFileName = "shared.db"
)

type DatabaseManager struct {
	Migrations fs.FS
}

func NewDatabaseManager(migrations fs.FS) *DatabaseManager {
	return &DatabaseManager{
		Migrations: migrations,
	}
}

func (dm *DatabaseManager) GetConnection() (*sql.DB, error) {
	if _, err := os.Stat(dbFileName); os.IsNotExist(err) {
		file, err := os.Create(dbFileName)
		if err != nil {
			return nil, fmt.Errorf("failed to create database file: %v", err)
		}
		file.Close()
	}

	db, err := sql.Open("libsql", "file:"+dbFileName+"?journal_mode=WAL&sync=full&foreign_keys=on")
	if err != nil {
		panic(fmt.Sprintf("Failed to open database: %v", err))
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		db.Close()
		panic(fmt.Sprintf("Failed to connect to the database: %v", err))
	}

	return db, nil
}

func (dm *DatabaseManager) InitializeDatabase() {
	db, err := dm.GetConnection()
	if err != nil {
		panic(fmt.Sprintf("Failed to get database connection: %v", err))
	}
	defer db.Close()

	files, err := fs.ReadDir(dm.Migrations, "migrations/shared")
	if err != nil {
		panic(fmt.Sprintf("Failed to read migration files: %v", err))
	}

	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}

		fileSections := strings.Split(file.Name(), ".")
		if fileSections[FILE_UP_MIGRATION_INDEX] != "up" {
			continue // skip files that are not "up" migrations
		}

		migration, err := fs.ReadFile(dm.Migrations, "migrations/shared/"+file.Name())
		if err != nil {
			panic(fmt.Sprintf("Failed to read migration file %s: %v", file.Name(), err))
		}

		_, err = db.Exec(string(migration))
		if err != nil {
			panic(fmt.Sprintf("Failed to execute migration %s: %v", file.Name(), err))
		}
	}
}
