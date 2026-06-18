package identity

import (
	"fmt"
	"strings"

	"github.com/ariaci/kopia-provisioner/kopia"
)

type KopiaIdentityConfig struct {
	Line int
}
type KopiaIdentities map[Identity]KopiaIdentityConfig

type KopiaIdentityEntry struct {
	Identity Identity
	Config   KopiaIdentityConfig
}
type KopiaIdentityEntries []KopiaIdentityEntry

func (e KopiaIdentityEntry) Compare(other KopiaIdentityEntry) int {
	if r := e.Identity.Compare(other.Identity); r != 0 {
		return r
	}
	return e.Config.Compare(other.Config)
}

func (c KopiaIdentityConfig) Compare(other KopiaIdentityConfig) int {
	switch {
	case c.Line < other.Line:
		return -1
	case c.Line > other.Line:
		return 1
	default:
		return 0
	}
}

func (i KopiaIdentities) MakeEntries() KopiaIdentityEntries {
	var entries = make(KopiaIdentityEntries, 0, len(i))
	for identity, config := range i {
		entries = append(entries, KopiaIdentityEntry{
			Identity: identity,
			Config:   config,
		})
	}

	return entries
}

func LoadKopiaIdentities(kopiaRepoConfigPath string) (KopiaIdentities, error) {
	var identities = make(KopiaIdentities)

	err := kopia.Run(func(line string) {
		line = strings.TrimSpace(line)
		if line == "" {
			return
		}

		identities[newIdentityFromIdentity(line)] = KopiaIdentityConfig{Line: len(identities) + 1}
	}, kopiaRepoConfigPath, "server", "users", "list")

	if err != nil {
		return nil, fmt.Errorf("kopia users list failed: %w", err)
	}

	return identities, nil
}
