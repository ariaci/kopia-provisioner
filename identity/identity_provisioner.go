package identity

import (
	"fmt"
	"maps"
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	Users map[string]userConfig `yaml:"users"`
}

type ProvisionerIdentityHosts map[string]ProvisionerIdentityConfig

type userConfig struct {
	Default ProvisionerIdentityConfig `yaml:"default"`
	Hosts   ProvisionerIdentityHosts  `yaml:"hosts"`
}

type ProvisionerIdentityConfig struct {
	Password PasswordWrapper `yaml:"password"`
}

type provisionerIdentities map[Identity]ProvisionerIdentityConfig

func (c config) makeIdentities() provisionerIdentities {
	identities := make(provisionerIdentities)

	for user, userConfig := range c.Users {
		// set default password to NilPassword if not set, so we have no problems later
		// with nil interfaces when trying to call KopiaArguments() on them
		if !userConfig.Default.Password.IsSet() {
			userConfig.Default.Password = PasswordWrapper{
				Password: NilPassword{},
			}
		}
		// fill in identities with corresponding (default) passwords
		for host, hostConfig := range userConfig.Hosts {
			if !hostConfig.Password.IsSet() {
				hostConfig.Password = userConfig.Default.Password
			}
			identities[newIdentityFromUserAndHost(user, host)] = hostConfig
		}
	}

	return identities
}

func (h *ProvisionerIdentityHosts) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// First try: Unmarshal as List
	var list []string

	if err := unmarshal(&list); err == nil {
		m := make(ProvisionerIdentityHosts)
		for _, host := range list {
			m[host] = ProvisionerIdentityConfig{}
		}
		*h = m
		return nil
	}

	// Second try: Unmarshal as Map
	var m2 map[string]ProvisionerIdentityConfig
	if err := unmarshal(&m2); err == nil {
		m := make(ProvisionerIdentityHosts)
		maps.Copy(m, m2)
		*h = m
		return nil
	}

	return fmt.Errorf("hosts must be either a list of strings or a map[string]ProvisionerIdentityConfig")
}

func loadProvisionerIdentities(configPath string) provisionerIdentities {
	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	var cfg config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	return cfg.makeIdentities()
}
