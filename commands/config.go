package commands

import (
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
