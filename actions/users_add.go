package actions

import (
	"kopia-provisioner/identity"
	"kopia-provisioner/kopia"
)

func (context UserActionContext) addIdentity(identity identity.Identity, info identity.IdentityInfo) error {
	_, err := kopia.Run(context.ConfigFile, context.kopiaArguments(kopia.KopiaActionAdd, identity, info)...)
	return err
}
