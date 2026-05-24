package actions

import (
	"fmt"

	"github.com/ariaci/kopia-provisioner/identity"

	"github.com/goccy/go-yaml"
)

func (context ConfigInitActionContext) buildUserConfig(ids identity.IdentityEntries) (identity.UserConfig, error) {
	u := identity.UserConfig{}
	u.Hosts = make(map[string]identity.ProvisionerIdentityConfig)

	c := identity.ProvisionerIdentityConfig{}

	switch context.Scope.Password {
	case PasswordScopeGlobal:
		// no-op, global password is set at the config level
	case PasswordScopeUser:
		u.Default.Password.Type = "nil"
	default:
		c.Password.Type = "nil"
	}

	for _, id := range ids {
		u.Hosts[id.Identity.Host] = c
	}

	return u, nil
}

func (context ConfigInitActionContext) buildConfig() (identity.Config, error) {
	c := identity.Config{Users: make(map[string]identity.UserConfig)}

	if context.Scope.Password == PasswordScopeGlobal {
		c.Default.Password.Type = "nil"
	}

	for ids := range context.iterateIdentities() {
		u, err := context.buildUserConfig(ids)
		if err != nil {
			return identity.Config{}, fmt.Errorf("failed to build user config: %w", err)
		}

		user := ids[0].Identity.User // guaranteed by iterateIdentities grouping
		c.Users[user] = u
	}

	return c, nil
}

func (context ConfigInitActionContext) Execute() error {
	c, err := context.buildConfig()
	if err != nil {
		return fmt.Errorf("failed to build config: %w", err)
	}

	out, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	fmt.Println(string(out))
	return nil
}
