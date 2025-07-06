package cmd_add_user

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/nahtann/trancome/config"
	"github.com/nahtann/trancome/internal/database"
	"github.com/nahtann/trancome/internal/styles"
)

var cmdWarning = styles.Yellow(
	"Please run 'trancome init' to create the shared database.",
)

var AddUserCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new user to the database",
	Long:  `Add a new user to the database with a unique ID and name.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		configEnvs, err := cfg.Load().CheckConsistency()
		if err != nil {
			fmt.Printf("%s\n%s", err, cmdWarning)
			return
		}

		db := configEnvs.SharedDB
		if db == "" {
			fmt.Printf("shared database path is not configured")
			return
		}

		dbPath := filepath.Join(configEnvs.DatabaseDir, db)

		dbManager := database.NewDatabaseManager(nil)

		err = database.WithDatabase(dbManager, dbPath, func(db *sql.DB) error {
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

			userDbPath := filepath.Join(configEnvs.DatabaseDir, configEnvs.UserDBDir)
			database.CreateUserDatabase(dbManager, userDbPath, uuid.String(), name)

			style := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#04B575"))

			fmt.Println(
				style.Render(
					fmt.Sprintf("User '%s' created successfully with ID %s", name, uuid.String()),
				),
			)

			return nil
		})
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	},
}

func init() {
	AddUserCmd.PersistentFlags().StringP("name", "n", "", "Name of the user to add")
	AddUserCmd.MarkFlagRequired("name") // Mark the username flag as required

	AddUserCmd.PersistentFlags().StringP("email", "e", "", "Email of the user to add")
}
