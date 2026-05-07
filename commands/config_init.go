package commands

import (
	"kopia-provisioner/actions"
	"strings"

	"github.com/spf13/cobra"
)

type scopeSliceValue struct {
	target *[]actions.ScopeEntry
}

var (
	scopeEntries []actions.ScopeEntry
)

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes provisioner YAML configuration based on Kopia identities",
	RunE: func(cmd *cobra.Command, args []string) error {
		context, err := makeConfigInitActionContext()
		if err != nil {
			return err
		}

		return context.Execute()
	},
}

func init() {
	configCmd.AddCommand(configInitCmd)
	addCommonConfigFlags(configInitCmd)

	configInitCmd.Flags().VarP(
		newScopeSliceValue(&scopeEntries),
		"scope",
		"s",
		"Set scope values (password:user, password:identity)",
	)
}

func newScopeSliceValue(target *[]actions.ScopeEntry) *scopeSliceValue {
	return &scopeSliceValue{target: target}
}

func (s *scopeSliceValue) String() string {
	parts := make([]string, len(*s.target))
	for i, e := range *s.target {
		parts[i] = e.String()
	}
	return strings.Join(parts, ",")
}

func (s *scopeSliceValue) Set(input string) error {
	var entry actions.ScopeEntry
	if err := entry.Set(input); err != nil {
		return err
	}
	*s.target = append(*s.target, entry)
	return nil
}

func (s *scopeSliceValue) Type() string {
	return "scope"
}
