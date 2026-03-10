package actions

import (
	"kopia-provisioner/identity"
	"kopia-provisioner/kopia"
)

func (context UserActionContext) updateIdentity(identity identity.Identity, info identity.IdentityInfo) error {
	args, err := context.kopiaArguments(kopia.KopiaActionUpdate, identity, info)
	if err != nil {
		return err
	}

	_, err = kopia.Run(context.ConfigFile, args...)
	return err
}
