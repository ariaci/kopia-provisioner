package identity

import (
	"kopia-provisioner/kopia"
	"log"
	"strings"
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

func LoadKopiaIdentities(kopiaRepoConfigPath string) KopiaIdentities {
	var identities = make(KopiaIdentities)

	out, err := kopia.Run(kopiaRepoConfigPath, "server", "users", "list")
	if err != nil {
		log.Fatalf("kopia users list failed: %v\n%s", err, out)
	}

	for i, line := range out {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		identities[newIdentityFromIdentity(line)] = KopiaIdentityConfig{Line: i + 1}
	}

	return identities
}
