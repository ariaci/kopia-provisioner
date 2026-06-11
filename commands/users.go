package commands

import (
	"fmt"
	"strings"

	"github.com/ariaci/kopia-provisioner/actions"
	"github.com/ariaci/kopia-provisioner/identity"

	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Synchronizes, adds or removes Kopia identities based on provisioner YAML configuration",
}

var (
	optUpdate     bool
	optDryRun     bool
	optVerbose    bool
	optAllowEmpty bool
	optForce      bool
)

func addCommonUserFlags(cmd *cobra.Command) {
	addCommonRootFlags(cmd)
	cmd.Flags().BoolVarP(&optDryRun, "dry-run", "d", false, "Simply shows what would happen without making any changes")
	cmd.Flags().BoolVarP(&optVerbose, "verbose", "v", false, "Enables verbose output")
	cmd.Flags().BoolVarP(&optAllowEmpty, "allow-empty", "E", false, "Allows empty provisioner configuration (use with caution, may lead to unintended consequences)")
}

func addForceFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&optForce, "force", "f", false, "Force deletion of identities even if snapshots exist (use with caution)")
}

func addUpdateUserFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&optUpdate, "update", "u", false, "Updates password(-hash) of all available Kopia identities")
}

func init() {
	rootCmd.AddCommand(usersCmd)
}

func makeUserActionContext(args []string) (actions.UserActionContext, error) {
	ids, err := identity.BuildIdentities(
		identity.BuildIdentitiesContext{
			Sources:                     identity.SourceKopia | identity.SourceProvisioner,
			KopiaRepoConfigPath:         configFile,
			ProvisionerConfigPath:       args[0],
			AllowEmptyProvisionerConfig: optAllowEmpty,
		})
	if err != nil {
		return actions.UserActionContext{}, fmt.Errorf("failed to build identities: %w", err)
	}

	ctx := actions.UserActionContext{Identities: ids}

	if cfgFile := strings.TrimSpace(configFile); len(cfgFile) > 0 {
		ctx.ConfigFile = cfgFile
	}

	if optDryRun {
		ctx.Flags |= actions.UserActionFlagDryRun
	}
	if optUpdate {
		ctx.Flags |= actions.UserActionFlagUpdate
	}
	if optVerbose {
		ctx.Flags |= actions.UserActionFlagVerbose
	}
	if optForce {
		ctx.Flags |= actions.UserActionFlagForce
	}

	return ctx, nil
}
