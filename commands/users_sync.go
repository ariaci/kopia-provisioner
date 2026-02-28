package commands

import (
	"kopia-provisioner/actions"

	"github.com/spf13/cobra"
)

var usersSyncCmd = &cobra.Command{
	Use:   "sync <config.yaml>",
	Short: "Synchronizes and optionally updates Kopia identities based on the YAML configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return makeUserActionContext(args).Execute(actions.UserActionSync)
	},
}

func init() {
	usersCmd.AddCommand(usersSyncCmd)
	addCommonUserFlags(usersSyncCmd)
	addUpdateUserFlags(usersSyncCmd)
}
