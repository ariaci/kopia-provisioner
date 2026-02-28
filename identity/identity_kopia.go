package identity

import (
	"kopia-provisioner/kopia"
	"log"
	"strings"
)

type KopiaIdentityConfig struct{}
type kopiaIdentities map[Identity]KopiaIdentityConfig

func loadKopiaIdentities(kopiaRepoConfigPath string) kopiaIdentities {
	var identities = make(kopiaIdentities)

	out, err := kopia.Run(kopiaRepoConfigPath, "server", "users", "list")
	if err != nil {
		log.Fatalf("kopia users list failed: %v\n%s", err, out)
	}

	for _, line := range out {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		identities[newIdentityFromIdentity(line)] = KopiaIdentityConfig{}
	}

	return identities
}
