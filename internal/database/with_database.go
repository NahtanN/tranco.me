package database

import (
	"database/sql"
	"log"
)

// WithDatabase is a helper function that manages database connections for commands
func WithDatabase(dbManager *DatabaseManager, fn func(*sql.DB) error) error {
	db, err := dbManager.GetConnection()
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database: %v", closeErr)
		}
	}()

	return fn(db)
}
