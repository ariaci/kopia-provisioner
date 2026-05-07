package commands

import (
	"strings"

	"kopia-provisioner/actions"
	"kopia-provisioner/identity"

	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Synchronizes, adds or removes Kopia identities based on provisioner YAML configuration",
}

var (
	optUpdate  bool
	optDryRun  bool
	optVerbose bool
)

func addCommonUserFlags(cmd *cobra.Command) {
	addCommonRootFlags(cmd)
	cmd.Flags().BoolVarP(&optDryRun, "dry-run", "d", false, "Simply shows what would happen without making any changes")
	cmd.Flags().BoolVarP(&optVerbose, "verbose", "v", false, "Enables verbose output")
}

func addUpdateUserFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&optUpdate, "update", "u", false, "Updates password(-hash) of all available Kopia identities")
}

func init() {
	rootCmd.AddCommand(usersCmd)
}

func makeUserActionContext(args []string) actions.UserActionContext {
	var context = actions.UserActionContext{
		Identities: identity.BuildIdentities(
			identity.BuildIdentitiesContext{
				Sources:               identity.SourceKopia | identity.SourceProvisioner,
				KopiaRepoConfigPath:   configFile,
				ProvisionerConfigPath: args[0],
			}),
	}

	if cfgFile := strings.TrimSpace(configFile); len(cfgFile) > 0 {
		context.ConfigFile = cfgFile
	}

	if optDryRun {
		context.Flags |= actions.UserActionFlagDryRun
	}
	if optUpdate {
		context.Flags |= actions.UserActionFlagUpdate
	}
	if optVerbose {
		context.Flags |= actions.UserActionFlagVerbose
	}

	return context
}
