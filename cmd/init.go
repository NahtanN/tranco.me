package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/nahtann/trancome/config"
	"github.com/nahtann/trancome/internal/database"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up the initial configuration for the application",
	Long: `The init command initializes the application by setting up the necessary configuration. 

Set username and other required parameters to get started with the application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		configEnvs := cfg.Load().CreateDefault()

		dbManager.InitializeDatabase(configEnvs)

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatalf("Error getting name flag: %v", err)
			return
		}

		email, err := cmd.Flags().GetString("email")
		if err != nil {
			log.Fatalf("Error getting email flag: %v", err)
			return
		}

		db := configEnvs.SharedDB
		if db == "" {
			log.Fatal("Shared database path is not configured.")
			return
		}

		dbPath := filepath.Join(configEnvs.DatabaseDir, db)

		database.WithDatabase(dbManager, dbPath, func(db *sql.DB) error {
			query := `INSERT INTO users (id, name, email) VALUES (?, ?, ?)`
			uuid, err := uuid.NewV7()
			if err != nil {
				fmt.Println("Error generating UUID:", err)
			}

			result, err := db.Exec(query, uuid, name, email)
			if err != nil {
				fmt.Println("Error inserting user:", err)
			}

			id, _ := result.LastInsertId()
			fmt.Printf("User '%s' created with ID %d\n", name, id)

			style := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#04B575"))

			fmt.Println(style.Render("Application initialized successfully."))

			userDbPath := filepath.Join(configEnvs.DatabaseDir, configEnvs.UserDBDir)
			database.CreateUserDatabase(dbManager, userDbPath, uuid.String(), name)

			return nil
		})
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.PersistentFlags().StringP("name", "n", "", "Name for the root user (required)")
	initCmd.MarkPersistentFlagRequired("name")

	initCmd.PersistentFlags().
		StringP("email", "e", "", "Email address for the root user (optional)")
}
