package cmd

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/spf13/cobra"

	cmd_config "github.com/nahtann/trancome/cmd/config"
	cmd_user "github.com/nahtann/trancome/cmd/user"
	"github.com/nahtann/trancome/config"
	"github.com/nahtann/trancome/internal/database"
)

var (
	dbManager  *database.DatabaseManager
	migrations fs.FS

	cfgFile    string
	configEnvs *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "trancome",
	Short: "A brief description of your application",
	Long: `Tranco is a service that provides financial control and management tools for individuals. 

It offers features such as expense tracking, budget management, 
and financial insights to help users make informed decisions about their finances.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Tranco!")
	},
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
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(cmd_user.UserCmd)
	rootCmd.AddCommand(cmd_config.ConfigCmd)
}

func initConfig() {
	configEnvs = config.Load("")
}
