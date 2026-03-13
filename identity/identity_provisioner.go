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

func (c config) makeIdentities(context ProviderPipelineContext) provisionerIdentities {
	identities := make(provisionerIdentities)

	for user, userConfig := range c.Users {
		// instantiate root password of default password pipeline for the user
		err := userConfig.Default.Password.resolve(context, NilPassword{})
		if err != nil {
			panic(err)
		}

		for host, hostConfig := range userConfig.Hosts {
			// instantiate root password of password pipeline for the identity
			err := hostConfig.Password.resolve(context, userConfig.Default.Password)
			if err != nil {
				panic(err)
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

	return cfg.makeIdentities(newProviderPipelineContext(configPath))
}
