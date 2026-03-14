package identity

import (
	"kopia-provisioner/kopia"
	"log"
	"strings"
)

type KopiaIdentityConfig struct {
	Line int
}
type kopiaIdentities map[Identity]KopiaIdentityConfig

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

func loadKopiaIdentities(kopiaRepoConfigPath string) kopiaIdentities {
	var identities = make(kopiaIdentities)

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
