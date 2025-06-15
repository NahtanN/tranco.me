package cmd_config_set_dir

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nahtann/trancome/utils"
)

var SetDirCmd = &cobra.Command{
	Use:   "set-dir",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		newDir := args[0]

		// Expand path if it contains ~
		if expandedDir, err := utils.ExpandPath(newDir); err == nil {
			newDir = expandedDir
		}

		// Convert to absolute path
		absDir, err := filepath.Abs(newDir)
		if err != nil {
			fmt.Printf("Error converting to absolute path: %v\n", err)
			return
		}

		// Update the configuration file with the new directory
		viper.Set("database_dir", absDir)
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("Error writing configuration: %v\n", err)
			return
		}

		fmt.Printf("Database directory set to: %s\n", absDir)
		fmt.Println("Run 'myapp init' to initialize the new directory.")
	},
}

func init() {
}
