package commands

import (
	"kopia-provisioner/actions"

	"github.com/spf13/cobra"
)

var usersUpdateCmd = &cobra.Command{
	Use:   "update <config.yaml>",
	Short: "Updates Kopia identities based on the YAML configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return makeUserActionContext(args).Execute(actions.UserActionUpdate)
	},
}

func init() {
	usersCmd.AddCommand(usersUpdateCmd)
	addCommonUserFlags(usersUpdateCmd)
}
