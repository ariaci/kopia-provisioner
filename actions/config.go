package actions

import (
	"fmt"
	"iter"
	"kopia-provisioner/identity"
	"kopia-provisioner/utils"
	"slices"
	"strings"
)

type ScopeEntry struct {
	Scope string
	Value string
}

type PasswordScopeValue string

const (
	PasswordScopeUser     PasswordScopeValue = "user"
	PasswordScopeIdentity PasswordScopeValue = "identity"
)

type ScopeConfig struct {
	Password PasswordScopeValue
}

type ConfigActionContext struct {
	ConfigFile string
	Identities identity.Identities
}

type ConfigInitActionContext struct {
	ConfigActionContext
	Scope ScopeConfig
}

func (s *ScopeEntry) String() string {
	if s.Scope == "" {
		return ""
	}
	return s.Scope + ":" + s.Value
}

func (s *ScopeEntry) Type() string {
	return "scope"
}

func (s *ScopeEntry) Set(input string) error {
	parts := strings.SplitN(input, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid scope format %q, expected <scope>:<value>", input)
	}

	s.Scope = parts[0]
	s.Value = parts[1]
	return nil
}

func (c *ScopeConfig) Apply(entry ScopeEntry) error {
	switch entry.Scope {

	case "password":
		v := PasswordScopeValue(entry.Value)
		if !v.IsValid() {
			return fmt.Errorf("invalid password scope value %q", entry.Value)
		}
		c.Password = v
		return nil

	default:
		return fmt.Errorf("unknown scope %q", entry.Scope)
	}
}

func (v PasswordScopeValue) IsValid() bool {
	switch v {
	case PasswordScopeUser, PasswordScopeIdentity:
		return true
	}
	return false
}

func (c ConfigInitActionContext) iterateIdentities() iter.Seq[identity.IdentityEntries] {
	ids := c.Identities.MakeEntries()
	slices.SortFunc(ids, func(a, b identity.IdentityEntry) int {
		return a.Identity.Compare(b.Identity)
	})

	return utils.EqIterate(ids, func(a, b identity.IdentityEntry) bool {
		return a.Identity.User == b.Identity.User
	})
}
