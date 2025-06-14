package cmd_add_user

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/nahtann/trancome/internal/database"
)

var (
	dbManager *database.DatabaseManager
	name      string
	email     string
)

var AddUserCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new user to the database",
	Long:  `Add a new user to the database with a unique ID and name.`,
	Run: func(cmd *cobra.Command, args []string) {
		database.WithDatabase(dbManager, func(db *sql.DB) error {
			query := `INSERT INTO users (id, name) VALUES (?, ?)`
			uuid, err := uuid.NewV7()
			if err != nil {
				fmt.Println("Error generating UUID:", err)
			}

			result, err := db.Exec(query, uuid, name)
			if err != nil {
				fmt.Println("Error inserting user:", err)
			}

			id, _ := result.LastInsertId()
			fmt.Printf("User '%s' created with ID %d\n", name, id)
			fmt.Println("Application initialized successfully.")

			return nil
		})
	},
}

func init() {
	dbManager = database.NewDatabaseManager(nil)

	AddUserCmd.Flags().
		StringVarP(&name, "name", "n", "", "Name for the root user (required)")
	AddUserCmd.MarkFlagRequired("name") // Mark the username flag as required

	AddUserCmd.Flags().
		StringVarP(&email, "email", "e", "", "Email address for the root user (optional)")
}
