package identity

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Password interface {
	KopiaArguments() ([]string, error)
}

type PasswordWrapper struct {
	Password `yaml:"-"`
	Type     string `yaml:"-"`
	Value    string `yaml:"-"`
}

type NilPassword struct {
}

type PlainPassword struct {
	Password string
}

type KopiaPasswordHash struct {
	Hash string
}

func (p NilPassword) KopiaArguments() ([]string, error) {
	return nil, fmt.Errorf("no password provided")
}

func (p PlainPassword) KopiaArguments() ([]string, error) {
	if p.Password == "" {
		return nil, fmt.Errorf("Kopia plain password cannot be empty")
	}

	return []string{"--user-password", p.Password}, nil
}

func (p KopiaPasswordHash) KopiaArguments() ([]string, error) {
	if p.Hash == "" {
		return nil, fmt.Errorf("Kopia password hash cannot be empty")
	}

	return []string{"--user-password-hash", p.Hash}, nil
}

func (w *PasswordWrapper) UnmarshalYAML(node *yaml.Node) error {
	// password is always a string like "plain:foo" or "kopia-hash:bar"
	if node.Kind != yaml.ScalarNode {
		return fmt.Errorf("password must be a scalar string")
	}

	raw := node.Value

	parts := strings.SplitN(raw, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid password format, expected [provider>...>]type:value")
	}

	w.Type = parts[0]
	w.Value = parts[1]

	return nil
}

func (w *PasswordWrapper) resolve(context ProviderPipelineContext, defaultPassword Password) error {
	if w.Password != nil {
		return nil // already resolved
	}

	if w.Type == "" {
		*w = PasswordWrapper{Password: defaultPassword}
		return nil
	}

	password, err := newPassword(context, w.Type, w.Value)
	if err != nil {
		return err
	}

	*w = PasswordWrapper{Password: password}
	return nil
}

func newPassword(context ProviderPipelineContext, pwType, value string) (Password, error) {
	parts := strings.SplitN(pwType, ">", 2)
	switch parts[0] {
	// --- Password types ---
	case "nil":
		return NilPassword{}, nil
	case "plain":
		return PlainPassword{Password: value}, nil
	case "kopia-hash":
		return KopiaPasswordHash{Hash: value}, nil
	// --- Provider types or unknown ---
	default:
		if factory, ok := providerRegistry[parts[0]]; ok {
			if len(parts) < 2 || parts[1] == "" {
				return nil, fmt.Errorf("missing next stage for provider '%s', expected %s[>provider]>type", pwType, parts[0])
			}
			return factory(context, value, parts[1])
		}

		return nil, fmt.Errorf("unknown password type or provider: %s", pwType)
	}
}
