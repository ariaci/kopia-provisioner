package commands

import (
	"strings"

	"github.com/ariaci/kopia-provisioner/actions"
	"github.com/ariaci/kopia-provisioner/identity"

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
	var context = actions.ConfigInitActionContext{
		ConfigActionContext: actions.ConfigActionContext{
			Identities: identity.BuildIdentities(
				identity.BuildIdentitiesContext{
					Sources:             identity.SourceKopia,
					KopiaRepoConfigPath: configFile,
				}),
		},
		Scope: actions.ScopeConfig{
			Password: actions.PasswordScopeIdentity,
		},
	}

	if cfgFile := strings.TrimSpace(configFile); len(cfgFile) > 0 {
		context.ConfigFile = cfgFile
	}

	for _, entry := range scopeEntries {
		if err := context.Scope.Apply(entry); err != nil {
			return context, err
		}
	}
	return context, nil
}
