package cmd_add_user

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/nahtann/trancome/internal/database"
)

var dbManager *database.DatabaseManager

var AddUserCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new user to the database",
	Long:  `Add a new user to the database with a unique ID and name.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.WithDatabase(dbManager, func(db *sql.DB) error {
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return fmt.Errorf("failed to get name flag: %w", err)
			}
			email, err := cmd.Flags().GetString("email")
			if err != nil {
				return fmt.Errorf("failed to get email flag: %w", err)
			}

			query := `INSERT INTO users (id, name, email) VALUES (?, ?, ?)`
			uuid, err := uuid.NewV7()
			if err != nil {
				return fmt.Errorf("failed to generate UUID: %w", err)
			}

			if email == "" {
				_, err = db.Exec(query, uuid, name, nil)
				if err != nil {
					return fmt.Errorf("failed to insert user into database: %w", err)
				}
			} else {
				_, err = db.Exec(query, uuid, name, email)
				if err != nil {
					return fmt.Errorf("failed to insert user into database: %w", err)
				}
			}

			fmt.Printf("User '%s' created with ID %s\n", name, uuid.String())
			fmt.Println("Application initialized successfully.")

			return nil
		})
	},
}

func init() {
	dbManager = database.NewDatabaseManager(nil)

	AddUserCmd.PersistentFlags().StringP("name", "n", "", "Name of the user to add")
	AddUserCmd.MarkFlagRequired("name") // Mark the username flag as required

	AddUserCmd.PersistentFlags().StringP("email", "e", "", "Email of the user to add")
}
