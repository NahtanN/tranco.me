package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nahtann/trancome/internal/styles"
	"github.com/nahtann/trancome/utils"
)

type Config struct {
	DatabaseDir string `mapstructure:"database_dir"`
	SharedDB    string `mapstructure:"shared_db"`
	UserDBDir   string `mapstructure:"user_db_dir"`
	noConfig    bool
}

func NewConfig() *Config {
	return &Config{
		DatabaseDir: "",
		SharedDB:    "shared.db",
		UserDBDir:   "users",
		noConfig:    true,
	}
}

func (c *Config) Load() *Config {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".trancome")

	viper.AutomaticEnv()

	// Set default values for configuration
	setDefaults()

	if err := viper.ReadInConfig(); err == nil {
		c.noConfig = false
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	} else {
		c.noConfig = true
		log.Fatalf("Error reading config file: %v. Please run 'trancome init' to create a default configuration.", err)
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

func (c *Config) CreateDefault() *Config {
	if c.noConfig {
		if err := createDefaultConfig(); err != nil {
			log.Printf("Warning: Could not read config file: %v", err)
		}
	}

	return c.Load()
}

func (c *Config) CheckConsistency() (*Config, error) {
	// validate if share database file exists
	sharedDBPath := filepath.Join(c.DatabaseDir, c.SharedDB)
	if _, err := os.Stat(sharedDBPath); os.IsNotExist(err) {
		return nil, fmt.Errorf(
			styles.Red(fmt.Sprintf("Shared database file does not exist: %s", sharedDBPath)),
		)
	}

	// validate if user database directory exists
	userDBDirPath := filepath.Join(c.DatabaseDir, c.UserDBDir)
	if _, err := os.Stat(userDBDirPath); os.IsNotExist(err) {
		log.Fatalf(
			"User database directory does not exist: %s. Please run 'trancome init' to create the user database directory.",
			userDBDirPath,
		)
	}

	// validate if user database directory is a directory
	info, err := os.Stat(userDBDirPath)
	if err != nil {
		log.Fatalf("Error checking user database directory: %v", err)
	}
	if !info.IsDir() {
		log.Fatalf(
			"User database path is not a directory: %s. Please run 'trancome init' to create the user database directory.",
			userDBDirPath,
		)
	}

	// validate if user database directory is writable
	file, err := os.CreateTemp(userDBDirPath, "tempfile")
	if err != nil {
		log.Fatalf(
			"User database directory is not writable: %s. Please check permissions.",
			userDBDirPath,
		)
	}
	file.Close()
	os.Remove(file.Name())

	return c, nil
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
