package actions

import (
	"kopia-provisioner/identity"
	"kopia-provisioner/kopia"
)

func (context UserActionContext) updateIdentity(identity identity.Identity, info identity.IdentityInfo) error {
	_, err := kopia.Run(context.ConfigFile, context.kopiaArguments(kopia.KopiaActionUpdate, identity, info)...)
	return err
}
