package cmd

import (
	"github.com/spf13/cobra"

	set_up "github.com/nahtann/trancome/internal/domain/init/use_case"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up the initial configuration for the application",
	Long: `The init command initializes the application by setting up the necessary configuration. 

Set username and other required parameters to get started with the application.`,
	Run: func(cmd *cobra.Command, args []string) {
		setUpUseCase := set_up.NewSetUpUseCase(migrations)
		setUpUseCase.Execute()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// initCmd.Flags().
	// 	StringVarP(&username, "username", "u", "", "Username for the application (required)")
	// initCmd.MarkFlagRequired("username") // Mark the username flag as required
}
