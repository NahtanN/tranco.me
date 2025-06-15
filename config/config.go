package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nahtann/trancome/utils"
)

type Config struct {
	DatabaseDir string `mapstructure:"database_dir"`
	SharedDB    string `mapstructure:"shared_db"`
	UserDBDir   string `mapstructure:"user_db_dir"`
}

func Load(cfgFile string) *Config {
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

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	// Override with command line flag if provided
	if dbDir := viper.GetString("database_dir"); dbDir != "" {
		if expandedDir, err := utils.ExpandPath(dbDir); err == nil {
			config.DatabaseDir = expandedDir
		} else {
			config.DatabaseDir = dbDir
		}
	}

	return &config
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
