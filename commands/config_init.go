package commands

import (
	"github.com/spf13/cobra"
)

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes provisioner YAML configuration based on Kopia identities",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	configCmd.AddCommand(configInitCmd)
	addCommonConfigFlags(configInitCmd)
}
