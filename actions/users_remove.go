package actions

import (
	"kopia-provisioner/identity"
	"kopia-provisioner/kopia"
)

func (context UserActionContext) deleteIdentity(identity identity.Identity, info identity.IdentityInfo) error {
	_, err := kopia.Run(context.ConfigFile, context.kopiaArguments(kopia.KopiaActionRemove, identity, info)...)
	return err
}
