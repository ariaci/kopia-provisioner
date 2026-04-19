package commands

import (
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

func makeConfigInitActionContext() actions.ConfigInitActionContext {
	var context = actions.ConfigInitActionContext{
		Identities: identity.BuildIdentities(
			identity.BuildIdentitiesContext{
				Sources:             identity.SourceKopia,
				KopiaRepoConfigPath: configFile,
			}),
	}

	if cfgFile := strings.TrimSpace(configFile); len(cfgFile) > 0 {
		context.ConfigFile = cfgFile
	}

	return context
}
