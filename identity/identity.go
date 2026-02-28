package identity

import "strings"

type IdentitySource uint8

const SourceNone IdentitySource = 0
const (
	SourceKopia IdentitySource = 1 << iota
	SourceProvisioner
)

type Identity struct {
	User string
	Host string
}

type IdentityInfo struct {
	Source      IdentitySource
	Kopia       KopiaIdentityConfig
	Provisioner ProvisionerIdentityConfig
}

type Identities map[Identity]IdentityInfo

func newIdentityFromIdentity(identity string) Identity {
	parts := strings.SplitN(identity, "@", 2)
	return Identity{
		User: strings.ToLower(parts[0]),
		Host: strings.ToLower(parts[1]),
	}
}

func newIdentityFromUserAndHost(user, host string) Identity {
	return Identity{
		User: strings.ToLower(user),
		Host: strings.ToLower(host),
	}
}

func (id Identity) String() string {
	return id.User + "@" + id.Host
}

func (info IdentityInfo) HasKopia() bool {
	return info.Source&SourceKopia != 0
}

func (info IdentityInfo) HasProvisioner() bool {
	return info.Source&SourceProvisioner != 0
}

func merge[T any](
	i Identities,
	src map[Identity]T,
	apply func(IdentityInfo, T) IdentityInfo,
) {
	for id, value := range src {
		i[id] = apply(i[id], value)
	}
}

func (i Identities) addKopiaIdentities(kids kopiaIdentities) {
	merge(i, kids, func(info IdentityInfo, config KopiaIdentityConfig) IdentityInfo {
		return IdentityInfo{
			Kopia:       config,
			Provisioner: info.Provisioner,
			Source:      info.Source | SourceKopia,
		}
	})
}

func (i Identities) addProvisionerIdentities(pids provisionerIdentities) {
	merge(i, pids, func(info IdentityInfo, config ProvisionerIdentityConfig) IdentityInfo {
		return IdentityInfo{
			Kopia:       info.Kopia,
			Provisioner: config,
			Source:      info.Source | SourceProvisioner,
		}
	})
}

type BuildIdentitiesContext struct {
	KopiaRepoConfigPath   string
	ProvisionerConfigPath string
}

func BuildIdentities(context BuildIdentitiesContext) Identities {
	identities := make(Identities)
	identities.addKopiaIdentities(loadKopiaIdentities(context.KopiaRepoConfigPath))
	identities.addProvisionerIdentities(loadProvisionerIdentities(context.ProvisionerConfigPath))

	return identities
}
