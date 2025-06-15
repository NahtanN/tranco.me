package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nahtann/trancome/internal/database"
)

var (
	name  string
	email string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up the initial configuration for the application",
	Long: `The init command initializes the application by setting up the necessary configuration. 

Set username and other required parameters to get started with the application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbManager.InitializeDatabase(&configEnvs)

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
			fmt.Println("Application initialized successfully.")

			userDbPath := filepath.Join(configEnvs.DatabaseDir, configEnvs.UserDBDir)
			database.CreateUserDatabase(dbManager, userDbPath, uuid.String(), name)

			return nil
		})
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().
		StringVarP(&name, "name", "n", "", "Name for the root user (required)")
	initCmd.MarkFlagRequired("name") // Mark the username flag as required

	initCmd.Flags().
		StringVarP(&email, "email", "e", "", "Email address for the root user (optional)")

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".trancome")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
	}

	if err := viper.Unmarshal(&configEnvs); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}
}
