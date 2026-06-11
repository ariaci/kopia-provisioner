package commands

import (
	"github.com/ariaci/kopia-provisioner/actions"

	"github.com/spf13/cobra"
)

var usersSyncCmd = &cobra.Command{
	Use:   "sync <user-definitions-file>",
	Short: "Synchronizes and optionally updates Kopia identities based on the YAML configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := makeUserActionContext(args)
		if err != nil {
			return err
		}
		return ctx.Execute(actions.UserActionSync)
	},
}

func init() {
	usersCmd.AddCommand(usersSyncCmd)
	addCommonUserFlags(usersSyncCmd)
	addUpdateUserFlags(usersSyncCmd)
	addForceFlag(usersSyncCmd)
}
