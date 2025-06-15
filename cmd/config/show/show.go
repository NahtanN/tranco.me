package cmd_config_show

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseDir string `mapstructure:"database_dir"`
	SharedDB    string `mapstructure:"shared_db"`
	UserDBDir   string `mapstructure:"user_db_dir"`
}

var config Config

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Current Configuration:")
		fmt.Printf("  Database Directory: %s\n", config.DatabaseDir)
		fmt.Printf("  Shared Database: %s\n", config.SharedDB)
		fmt.Printf("  User Database Directory: %s\n", config.UserDBDir)
		fmt.Printf("  Config File: %s\n", viper.ConfigFileUsed())
	},
}

func init() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".trancome")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}
}
