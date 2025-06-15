package database

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/tursodatabase/go-libsql"

	"github.com/nahtann/trancome/config"
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

func (dm *DatabaseManager) GetConnection(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("libsql", "file:"+dbPath+"?journal_mode=WAL&sync=full&foreign_keys=on")
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

func (dm *DatabaseManager) InitializeDatabase(config *config.Config) {
	if err := os.MkdirAll(config.DatabaseDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create database directory: %v", err))
	}

	// Create user database subdirectory
	userDBDir := filepath.Join(config.DatabaseDir, config.UserDBDir)
	if err := os.MkdirAll(userDBDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create user database directory: %v", err))
	}

	// Initialize shared database (create if doesn't exist)
	sharedDBPath := filepath.Join(config.DatabaseDir, config.SharedDB)
	if err := initializeSharedDB(sharedDBPath); err != nil {
		panic(fmt.Sprintf("Failed to initialize shared database: %v", err))
	}

	db, err := dm.GetConnection(sharedDBPath)
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

func initializeSharedDB(sharedDBPath string) error {
	if _, err := os.Stat(sharedDBPath); os.IsNotExist(err) {
		file, err := os.Create(sharedDBPath)
		if err != nil {
			return fmt.Errorf("failed to create shared database file: %v", err)
		}
		file.Close()
	}

	return nil
}

func CreateUserDatabase(dbManager *DatabaseManager, dbPath string, userID string, username string) {
	fmt.Println("dbPath", dbPath)
	userDBPath := filepath.Join(dbPath, userID+"_"+username+".db")
	if _, err := os.Stat(userDBPath); os.IsNotExist(err) {
		file, err := os.Create(userDBPath)
		if err != nil {
			panic(fmt.Sprintf("Failed to create user database file: %v", err))
		}
		file.Close()
	}

	fmt.Println("userDBPath", userDBPath)

	userDB, err := dbManager.GetConnection(userDBPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to user database: %v", err))
	}
	defer userDB.Close()

	query := `CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT
  )`
	if _, err := userDB.Exec(query); err != nil {
		panic(fmt.Sprintf("Failed to create users table: %v", err))
	}

	fmt.Printf("User database for '%s' created successfully at %s\n", username, userDBPath)
}
