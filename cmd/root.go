/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io/fs"
	"os"

	"github.com/spf13/cobra"
)

var migrations fs.FS

// rootCmd represents the base command when called without any subcommands
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(migrationSource fs.FS) {
	migrations = migrationSource

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.trancome.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
