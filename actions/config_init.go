package actions

import (
	"fmt"
	"kopia-provisioner/identity"

	"gopkg.in/yaml.v3"
)

func (context ConfigInitActionContext) Execute() error {
	c := identity.Config{Users: make(map[string]identity.UserConfig)}
	for ids := range context.iterateIdentities() {
		u := identity.UserConfig{}
		u.Hosts = make(map[string]identity.ProvisionerIdentityConfig)
		u.Default.Password.Type = "nil"
		for _, id := range ids {
			u.Hosts[id.Identity.Host] = identity.ProvisionerIdentityConfig{}
		}
		c.Users[ids[0].Identity.User] = u
	}

	out, _ := yaml.Marshal(c)

	fmt.Println(string(out))
	return nil
}
