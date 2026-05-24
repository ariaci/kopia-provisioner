package identity

import (
	"bytes"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

type Config struct {
	Default ProvisionerIdentityConfig `yaml:"default"`
	Users   map[string]UserConfig     `yaml:"users"`
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

func findLineForKey(m ast.MappingNode, key string) (int, error) {
	for _, t := range m.Values {
		var k string
		if err := yaml.NodeToValue(t.Key, &k); err != nil {
			return 0, err
		}

		if key == k {
			return t.Key.GetToken().Position.Line, nil
		}
	}

	return 0, fmt.Errorf("key %q not found", key)
}

func loadProvisionerIdentities(configPath string, allowEmpty bool) (provisionerIdentities, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data), yaml.DisallowUnknownField())

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

	// instantiate root password of default password pipeline for all users
	err := c.Default.Password.resolve(context, NilPassword{})
	if err != nil {
		return nil, fmt.Errorf("failed to resolve global default password defined at line %d: %w", c.Default.Line, err)
	}

	for user, userConfig := range c.Users {
		// instantiate root password of default password pipeline for the user
		err := userConfig.Default.Password.resolve(context, c.Default.Password)
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

func (c *ProvisionerIdentityConfig) UnmarshalYAML(node ast.Node) error {
	var raw ProvisionerIdentityConfigRaw
	if err := yaml.NodeToValue(node, &raw, yaml.DisallowUnknownField()); err != nil {
		return err
	}
	c.Password = raw.Password
	c.Line = node.GetToken().Position.Line

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

func (h *ProvisionerIdentityHost) UnmarshalYAML(node ast.Node) error {
	if err := yaml.NodeToValue(node, &h.Host, yaml.DisallowUnknownField()); err == nil {
		h.Config = ProvisionerIdentityConfig{
			Line: node.GetToken().Position.Line,
		}
		return nil
	}

	return fmt.Errorf("host must be a string, not a map")
}

func (h *ProvisionerIdentityHosts) fixMissingLines(m ast.MappingNode) error {
	for k, v := range *h {
		if v.Line != 0 {
			continue
		}

		l, err := findLineForKey(m, k)
		if err != nil {
			return err
		}

		v.Line = l
		(*h)[k] = v
	}

	return nil
}

func (h *ProvisionerIdentityHosts) UnmarshalYAML(node ast.Node) error {
	var err error
	switch n := node.(type) {
	case *ast.SequenceNode:
		var l []ProvisionerIdentityHost
		if err = yaml.NodeToValue(node, &l, yaml.DisallowUnknownField()); err == nil {
			*h = make(ProvisionerIdentityHosts)
			for _, e := range l {
				(*h)[e.Host] = e.Config
			}
			return nil
		}
	case *ast.MappingNode:
		var m map[string]ProvisionerIdentityConfig
		if err = yaml.NodeToValue(node, &m, yaml.DisallowUnknownField()); err == nil {
			*h = m
			if err = h.fixMissingLines(*n); err != nil {
				return err
			}

			return nil
		}
	}

	if err == nil {
		return fmt.Errorf("hosts must be either a list of strings or a map[string]ProvisionerIdentityConfig")
	}

	return fmt.Errorf("failed to parse hosts: %w", err)
}

func (w PasswordWrapper) MarshalYAML() (any, error) {
	return fmt.Sprint(w.Type, ":", w.Value), nil
}
