package actions

import (
	"github.com/ariaci/kopia-provisioner/identity"
	"github.com/ariaci/kopia-provisioner/kopia"
)

func (context UserActionContext) updateIdentity(identity identity.Identity, info identity.IdentityInfo) error {
	args, err := context.kopiaArguments(kopia.KopiaActionUpdate, identity, info)
	if err != nil {
		return err
	}

	return kopia.Run(kopia.Discard, context.ConfigFile, args...)
}
