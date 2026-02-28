package identity

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Password interface {
	KopiaArguments() []string
}

type PasswordWrapper struct {
	Password
}

type PlainPassword struct {
	Password string
}

type KopiaPasswordHash struct {
	Hash string
}

func (w PasswordWrapper) IsSet() bool {
	return w.Password != nil
}

func (p PlainPassword) KopiaArguments() []string {
	return []string{"--user-password", p.Password}
}

func (p KopiaPasswordHash) KopiaArguments() []string {
	return []string{"--user-password-hash", p.Hash}
}

func (w *PasswordWrapper) UnmarshalYAML(node *yaml.Node) error {
	// password is always a string like "plain:foo" or "kopia-hash:bar"
	if node.Kind != yaml.ScalarNode {
		return fmt.Errorf("password must be a scalar string")
	}

	raw := node.Value

	parts := strings.SplitN(raw, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid password format, expected type:value")
	}

	typ := parts[0]
	value := parts[1]

	pw, err := newPassword(typ, value)
	if err != nil {
		return err
	}

	w.Password = pw
	return nil
}

func newPassword(pwType, value string) (Password, error) {
	switch pwType {
	case "plain":
		return PlainPassword{Password: value}, nil
	case "kopia-hash":
		return KopiaPasswordHash{Hash: value}, nil
	default:
		return nil, fmt.Errorf("unknown password type: %s", pwType)
	}
}
