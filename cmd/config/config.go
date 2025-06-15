package cmd_config

import (
	"github.com/spf13/cobra"

	cmd_config_set_dir "github.com/nahtann/trancome/cmd/config/set_dir"
	cmd_config_show "github.com/nahtann/trancome/cmd/config/show"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	ConfigCmd.AddCommand(cmd_config_show.ShowCmd)
	ConfigCmd.AddCommand(cmd_config_set_dir.SetDirCmd)
}
