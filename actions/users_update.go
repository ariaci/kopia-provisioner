package actions

import (
	"github.com/ariaci/kopia-provisioner/identity"
	"github.com/ariaci/kopia-provisioner/kopia"
)

func (ctx UserActionContext) updateIdentity(identity identity.Identity, info identity.IdentityInfo) error {
	args, err := ctx.kopiaArguments(kopia.KopiaActionUpdate, identity, info)
	if err != nil {
		return err
	}

	return kopia.Run(kopia.Discard, ctx.ConfigFile, args...)
}
