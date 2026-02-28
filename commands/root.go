package commands

import (
	"github.com/spf13/cobra"
)

var (
	configFile string
)

var rootCmd = &cobra.Command{
	Use:   "kopia-provisioner",
	Short: "Provisioner for Kopia identities based on YAML configuration",
}

func addCommonRootFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&configFile, "config-file", "c", "", "Specify the kopia repository configuration file to use")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
