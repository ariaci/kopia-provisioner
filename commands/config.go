package commands

import (
	"fmt"
	"kopia-provisioner/actions"
	"kopia-provisioner/identity"
	"strings"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manages provisioner YAML configuration",
}

func addCommonConfigFlags(cmd *cobra.Command) {
	addCommonRootFlags(cmd)
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func makeConfigInitActionContext() (actions.ConfigInitActionContext, error) {
	ids, err := identity.BuildIdentities(
		identity.BuildIdentitiesContext{
			Sources:             identity.SourceKopia,
			KopiaRepoConfigPath: configFile,
		})
	if err != nil {
		return actions.ConfigInitActionContext{}, fmt.Errorf("failed to build identities: %w", err)
	}

	ctx := actions.ConfigInitActionContext{
		ConfigActionContext: actions.ConfigActionContext{Identities: ids},
		Scope:               actions.ScopeConfig{Password: actions.PasswordScopeIdentity},
	}

	if cfgFile := strings.TrimSpace(configFile); len(cfgFile) > 0 {
		ctx.ConfigFile = cfgFile
	}

	for _, entry := range scopeEntries {
		if err := ctx.Scope.Apply(entry); err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}
