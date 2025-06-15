package database

import (
	"database/sql"
	"fmt"
	"log"
)

// TransactionManager provides manual transaction control for complex operations
type TransactionManager struct {
	db *sql.DB
	tx *sql.Tx
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(dbManager *DatabaseManager, dbPath string) (*TransactionManager, error) {
	db, err := dbManager.GetConnection(dbPath)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &TransactionManager{
		db: db,
		tx: tx,
	}, nil
}

// Tx returns the underlying transaction
func (tm *TransactionManager) Tx() *sql.Tx {
	return tm.tx
}

// Commit commits the transaction and closes the database connection
func (tm *TransactionManager) Commit() error {
	defer tm.db.Close()

	if err := tm.tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Rollback rolls back the transaction and closes the database connection
func (tm *TransactionManager) Rollback() error {
	defer tm.db.Close()

	if err := tm.tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	return nil
}

// Close properly closes the transaction and database connection
// It will rollback if the transaction hasn't been committed
func (tm *TransactionManager) Close() error {
	defer tm.db.Close()

	// Try to rollback first (will fail if already committed, which is fine)
	if rollbackErr := tm.tx.Rollback(); rollbackErr != nil {
		// Only log if it's not the "transaction already committed" error
		if rollbackErr != sql.ErrTxDone {
			log.Printf("Error during transaction cleanup: %v", rollbackErr)
		}
	}

	return nil
}
