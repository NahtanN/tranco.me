package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cmd_config "github.com/nahtann/trancome/cmd/config"
	cmd_user "github.com/nahtann/trancome/cmd/user"
	"github.com/nahtann/trancome/config"
	"github.com/nahtann/trancome/internal/database"
	"github.com/nahtann/trancome/utils"
)

var (
	dbManager  *database.DatabaseManager
	migrations fs.FS

	cfgFile    string
	configEnvs config.Config
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
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".trancome")
	}

	viper.AutomaticEnv()

	// Set default values for configuration
	setDefaults()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	} else {
		if err := createDefaultConfig(); err != nil {
			log.Printf("Warning: Could not read config file: %v", err)
		}
	}

	if err := viper.Unmarshal(&configEnvs); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	// Override with command line flag if provided
	if dbDir := viper.GetString("database_dir"); dbDir != "" {
		if expandedDir, err := utils.ExpandPath(dbDir); err == nil {
			configEnvs.DatabaseDir = expandedDir
		} else {
			configEnvs.DatabaseDir = dbDir
		}
	}
}

func setDefaults() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	defaultDBDir := filepath.Join(home, ".trancome", "databases")

	viper.SetDefault("database_dir", defaultDBDir)
	viper.SetDefault("shared_db", "shared.db")
	viper.SetDefault("user_db_dir", "users")
}

func createDefaultConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user home directory: %w", err)
	}

	configPath := filepath.Join(home, ".trancome.yaml")

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	viper.SetConfigFile(configPath)
	return viper.WriteConfigAs(configPath)
}
