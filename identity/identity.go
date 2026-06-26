package identity

import (
	"fmt"
	"strings"
)

type IdentitySource uint8

const SourceNone IdentitySource = 0
const (
	SourceKopia IdentitySource = 1 << iota
	SourceKopiaSnapshots
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

type IdentityEntry struct {
	Identity Identity
	Info     IdentityInfo
}
type IdentityEntries []IdentityEntry

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

func (e Identity) Compare(other Identity) int {
	if r := strings.Compare(e.User, other.User); r != 0 {
		return r
	}
	return strings.Compare(e.Host, other.Host)
}

func (s IdentitySource) Validate() error {
	if s&SourceKopiaSnapshots != 0 && s&SourceKopia == 0 {
		return fmt.Errorf("SourceKopiaSnapshots requires SourceKopia")
	}

	return nil
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

func (i Identities) addKopiaIdentities(kids KopiaIdentities) {
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

func (i Identities) MakeEntries() IdentityEntries {
	var entries = make(IdentityEntries, 0, len(i))
	for identity, info := range i {
		entries = append(entries, IdentityEntry{
			Identity: identity,
			Info:     info,
		})
	}

	return entries
}

func (e IdentityEntry) Compare(other IdentityEntry) int {
	switch {
	case e.Info.HasProvisioner() && !other.Info.HasProvisioner():
		return -1
	case !e.Info.HasProvisioner() && other.Info.HasProvisioner():
		return 1
	case e.Info.HasProvisioner() && other.Info.HasProvisioner():
		return e.Info.Provisioner.Compare(other.Info.Provisioner)
	case e.Info.HasKopia() && other.Info.HasKopia():
		return e.Info.Kopia.Compare(other.Info.Kopia)
	}

	return 0
}

type BuildIdentitiesContext struct {
	Sources                     IdentitySource
	KopiaRepoConfigPath         string
	ProvisionerConfigPath       string
	AllowEmptyProvisionerConfig bool
}

func BuildIdentities(context BuildIdentitiesContext) (Identities, error) {
	if err := context.Sources.Validate(); err != nil {
		return nil, err
	}

	identities := make(Identities)

	if context.Sources&SourceKopia != 0 {
		ids, err := LoadKopiaIdentities(context.KopiaRepoConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load kopia identities: %w", err)
		}

		if context.Sources&SourceKopiaSnapshots != 0 {
			err = ids.AssignSnapshots(context.KopiaRepoConfigPath)
			if err != nil {
				return nil, fmt.Errorf("failed to assign snapshots: %w", err)
			}
		}

		identities.addKopiaIdentities(ids)
	}
	if context.Sources&SourceProvisioner != 0 {
		ids, err := loadProvisionerIdentities(context.ProvisionerConfigPath, context.AllowEmptyProvisionerConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to load provisioner identities: %w", err)
		}
		identities.addProvisionerIdentities(ids)
	}

	return identities, nil
}
