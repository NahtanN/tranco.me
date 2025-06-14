package cmd

import (
	"io/fs"
	"os"

	"github.com/spf13/cobra"

	cmd_user "github.com/nahtann/trancome/cmd/user"
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
	rootCmd.AddCommand(cmd_user.UserCmd)
}
