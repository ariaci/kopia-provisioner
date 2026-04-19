package actions

import (
	"iter"
	"kopia-provisioner/identity"
	"kopia-provisioner/utils"
	"slices"
)

type ConfigActionContext struct {
	ConfigFile string
	Identities identity.Identities
}

type ConfigInitActionContext ConfigActionContext

func (c ConfigInitActionContext) iterateIdentities() iter.Seq[identity.IdentityEntries] {
	ids := c.Identities.MakeEntries()
	slices.SortFunc(ids, func(a, b identity.IdentityEntry) int {
		return a.Identity.Compare(b.Identity)
	})

	return utils.EqIterate(ids, func(a, b identity.IdentityEntry) bool {
		return a.Identity.User == b.Identity.User
	})
}
