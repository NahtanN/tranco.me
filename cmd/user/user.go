package cmd_user

import (
	"github.com/spf13/cobra"

	cmd_add_user "github.com/nahtann/trancome/cmd/user/add"
	"github.com/nahtann/trancome/internal/database"
)

var (
	dbManager *database.DatabaseManager
	name      string
)

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users in the application",
	Long:  `The user command allows you to manage users in the application.`,
}

func init() {
	UserCmd.AddCommand(cmd_add_user.AddUserCmd)
}
