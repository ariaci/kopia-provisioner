package identity

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Users map[string]UserConfig `yaml:"users"`
}

type ProvisionerIdentityHost struct {
	Host   string
	Config ProvisionerIdentityConfig
}
type ProvisionerIdentityHosts map[string]ProvisionerIdentityConfig

type UserConfig struct {
	Default ProvisionerIdentityConfig `yaml:"default"`
	Hosts   ProvisionerIdentityHosts  `yaml:"hosts"`
}

type ProvisionerIdentityConfig struct {
	Line     int             `yaml:"-"`
	Password PasswordWrapper `yaml:"password,omitempty"`
}
type ProvisionerIdentityConfigRaw ProvisionerIdentityConfig

type provisionerIdentities map[Identity]ProvisionerIdentityConfig

func findLineForKey(r *yaml.Node, key string) int {
	for i := 0; i < len(r.Content); i += 2 {
		if r.Content[i].Value == key {
			return r.Content[i].Line
		}
	}

	return 0
}

func loadProvisionerIdentities(configPath string, allowEmpty bool) (provisionerIdentities, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if !allowEmpty && len(cfg.Users) == 0 {
		return nil, fmt.Errorf("refusing to continue: users list is empty (possible YAML error)")
	}

	ctx, err := newProviderPipelineContext(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider pipeline context: %w", err)
	}

	return cfg.makeIdentities(ctx)
}

func (c Config) makeIdentities(context ProviderPipelineContext) (provisionerIdentities, error) {
	identities := make(provisionerIdentities)

	for user, userConfig := range c.Users {
		// instantiate root password of default password pipeline for the user
		err := userConfig.Default.Password.resolve(context, NilPassword{})
		if err != nil {
			return nil, fmt.Errorf("failed to resolve default password for user %q defined at line %d: %w", user, userConfig.Default.Line, err)
		}

		for host, hostConfig := range userConfig.Hosts {
			id := newIdentityFromUserAndHost(user, host)

			// check for duplicate identities in the provisioner configuration
			if prevConfig, exists := identities[id]; exists {
				return nil, fmt.Errorf(
					"duplicate identity %q found (first defined at line %d, duplicate at line %d)",
					id,
					prevConfig.Line,
					hostConfig.Line,
				)
			}

			// instantiate root password of password pipeline for the identity
			err := hostConfig.Password.resolve(context, userConfig.Default.Password)
			if err != nil {
				return nil, fmt.Errorf("failed to instantiate root password for identity %q defined at line %d: %w", id, hostConfig.Line, err)
			}

			identities[id] = hostConfig
		}
	}

	return identities, nil
}

func (c *ProvisionerIdentityConfig) UnmarshalYAML(node *yaml.Node) error {
	var raw ProvisionerIdentityConfigRaw
	if err := node.Decode(&raw); err != nil {
		return err
	}
	c.Password = raw.Password
	c.Line = node.Line

	return nil
}

func (c ProvisionerIdentityConfig) Compare(other ProvisionerIdentityConfig) int {
	switch {
	case c.Line < other.Line:
		return -1
	case c.Line > other.Line:
		return 1
	default:
		return 0
	}
}

func (h *ProvisionerIdentityHost) UnmarshalYAML(node *yaml.Node) error {
	if err := node.Decode(&h.Host); err == nil {
		h.Config = ProvisionerIdentityConfig{
			Line: node.Line,
		}
		return nil
	}

	return fmt.Errorf("host must be a string, not a map")
}

func (h *ProvisionerIdentityHosts) fixMissingLines(node *yaml.Node) {
	for k, v := range *h {
		if v.Line != 0 {
			continue
		}

		v.Line = findLineForKey(node, k)
		(*h)[k] = v
	}
}

func (h *ProvisionerIdentityHosts) UnmarshalYAML(node *yaml.Node) error {
	// First try: Unmarshal as List
	var l []ProvisionerIdentityHost
	if err := node.Decode(&l); err == nil {
		*h = make(ProvisionerIdentityHosts)
		for _, e := range l {
			(*h)[e.Host] = e.Config
		}
		return nil
	}

	// Second try: Unmarshal as Map
	var m map[string]ProvisionerIdentityConfig
	if err := node.Decode(&m); err == nil {
		*h = m
		h.fixMissingLines(node)

		return nil
	}

	return fmt.Errorf("hosts must be either a list of strings or a map[string]ProvisionerIdentityConfig")
}

func (w PasswordWrapper) MarshalYAML() (any, error) {
	return fmt.Sprint(w.Type, ":", w.Value), nil
}
