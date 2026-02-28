package commands

import (
	"kopia-provisioner/actions"

	"github.com/spf13/cobra"
)

var usersAddCmd = &cobra.Command{
	Use:   "add <config.yaml>",
	Short: "Adds and optionally updates Kopia identities based on the YAML configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return makeUserActionContext(args).Execute(actions.UserActionAdd)
	},
}

func init() {
	usersCmd.AddCommand(usersAddCmd)
	addCommonUserFlags(usersAddCmd)
	addUpdateUserFlags(usersAddCmd)
}
