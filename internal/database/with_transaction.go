package database

import (
	"database/sql"
	"fmt"
	"log"
)

// WithTransaction provides automatic transaction management
// If the function returns an error, the transaction is rolled back
// If the function completes successfully, the transaction is committed
func WithTransaction(dbManager *DatabaseManager, dbPath string, fn func(*sql.Tx) error) error {
	db, err := dbManager.GetConnection(dbPath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database: %v", closeErr)
		}
	}()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure transaction is handled properly
	defer func() {
		if r := recover(); r != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Failed to rollback transaction after panic: %v", rollbackErr)
			}
			panic(r)
		}
	}()

	// Execute the function
	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf(
				"transaction failed and rollback failed: %v (original error: %w)",
				rollbackErr,
				err,
			)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
