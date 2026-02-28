package commands

import (
	"kopia-provisioner/actions"

	"github.com/spf13/cobra"
)

var usersRemoveCmd = &cobra.Command{
	Use:   "remove <config.yaml>",
	Short: "Removes and optionally updates Kopia identities based on the YAML configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return makeUserActionContext(args).Execute(actions.UserActionRemove)
	},
}

func init() {
	usersCmd.AddCommand(usersRemoveCmd)
	addCommonUserFlags(usersRemoveCmd)
	addUpdateUserFlags(usersRemoveCmd)
}
