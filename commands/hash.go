package commands

import (
	"fmt"
	"strings"

	"github.com/ariaci/kopia-provisioner/kopia"

	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash <text>",
	Short: "Calculates hash of the given text for given kopia repository configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		args[0] = strings.TrimSpace(args[0])
		if args[0] == "" {
			return nil
		}

		hash := ""
		err := kopia.Run(func(line string) {
			line = strings.TrimSpace(line)
			if line != "" {
				hash = line
			}
		}, configFile, "server", "users", "hash-password", "--user-password", args[0])

		if err != nil {
			return fmt.Errorf("kopia users hash-password failed: %w", err)
		}
		if hash == "" {
			return fmt.Errorf("kopia users hash-password failed: no output")
		}

		fmt.Println(hash)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)
	addCommonRootFlags(hashCmd)
}
