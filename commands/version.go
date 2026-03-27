package commands

import (
	"fmt"
	"kopia-provisioner/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version of the kopia-provisioner",
	RunE: func(_ *cobra.Command, args []string) error {
		fmt.Println(version.Get())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
