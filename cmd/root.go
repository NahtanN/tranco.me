package cmd

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/nahtann/trancome/internal/database"
)

var (
	dbManager  *database.DatabaseManager
	migrations fs.FS
)

var rootCmd = &cobra.Command{
	Use:   "trancome",
	Short: "A brief description of your application",
	Long: `Tranco is a service that provides financial control and management tools for individuals. 

It offers features such as expense tracking, budget management, 
and financial insights to help users make informed decisions about their finances.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute(migrationSource fs.FS) {
	migrations = migrationSource

	dbManager = database.NewDatabaseManager(migrations)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}

// WithDatabase is a helper function that manages database connections for commands
func WithDatabase(fn func(*sql.DB) error) error {
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

// WithTransaction provides automatic transaction management
// If the function returns an error, the transaction is rolled back
// If the function completes successfully, the transaction is committed
func WithTransaction(fn func(*sql.Tx) error) error {
	db, err := dbManager.GetConnection()
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

// TransactionManager provides manual transaction control for complex operations
type TransactionManager struct {
	db *sql.DB
	tx *sql.Tx
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager() (*TransactionManager, error) {
	db, err := dbManager.GetConnection()
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
