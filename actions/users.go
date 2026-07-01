package actions

import (
	"fmt"
	"log"
	"slices"

	"github.com/ariaci/kopia-provisioner/identity"
	"github.com/ariaci/kopia-provisioner/kopia"
)

type UserActionFlags uint8

const UserActionFlagNone UserActionFlags = 0
const (
	UserActionFlagDryRun UserActionFlags = 1 << iota
	UserActionFlagUpdate
	UserActionFlagVerbose
	UserActionFlagForce
)

type UserActionContext struct {
	Flags      UserActionFlags
	ConfigFile string
	Identities identity.Identities
}

type UserActionClassification uint8

const (
	UserActionClassificationNone UserActionClassification = iota
	UserActionClassificationAdd
	UserActionClassificationDelete
	UserActionClassificationUpdate

	UserActionClassificationSkip  = 1 << 6
	UserActionClassificationError = 1 << 7
)

type UserActionClassificationResult struct {
	Classification UserActionClassification
	Reason         error
}

type UserAction uint8

const (
	UserActionSync UserAction = iota
	UserActionAdd
	UserActionRemove
	UserActionUpdate
)

func (ctx UserActionContext) hasConfigFile() bool {
	return len(ctx.ConfigFile) > 0
}

func (ctx UserActionContext) isDryRun() bool {
	return ctx.Flags&UserActionFlagDryRun != 0
}

func (ctx UserActionContext) shouldUpdate() bool {
	return ctx.Flags&UserActionFlagUpdate != 0
}

func (ctx UserActionContext) isVerbose() bool {
	return ctx.Flags&UserActionFlagVerbose != 0
}

func (ctx UserActionContext) classify(action UserAction, entry identity.IdentityEntry) UserActionClassificationResult {
	switch {
	// Case 1) Identity only available in kopia -> remove Identity from kopia
	case (action == UserActionSync || action == UserActionRemove) && entry.Info.HasKopia() && !entry.Info.HasProvisioner():
		// check if there are snapshots associated with the identity, if so, skip deletion unless force flag is set
		if s := len(entry.Info.Kopia.Snapshots); ctx.Flags&UserActionFlagForce == 0 && s > 0 {
			return UserActionClassificationResult{
				Classification: UserActionClassificationDelete | UserActionClassificationSkip,
				Reason:         fmt.Errorf("cannot remove identity %s: %d associated snapshot(s)", entry.Identity, s)}
		}
		return UserActionClassificationResult{Classification: UserActionClassificationDelete}
	// Case 2) Identity only available in provisioner -> add Identity to kopia
	case (action == UserActionSync || action == UserActionAdd) && !entry.Info.HasKopia() && entry.Info.HasProvisioner():
		return UserActionClassificationResult{Classification: UserActionClassificationAdd}
	// Case 3) Identity available in kopia and provisioner -> update Identity in kopia
	case (action == UserActionUpdate || ctx.shouldUpdate()) && entry.Info.HasKopia() && entry.Info.HasProvisioner():
		return UserActionClassificationResult{Classification: UserActionClassificationUpdate}

	default:
		return UserActionClassificationResult{Classification: UserActionClassificationNone}
	}
}

func (ctx UserActionContext) kopiaArguments(action kopia.KopiaAction, identity identity.Identity, info identity.IdentityInfo) ([]string, error) {
	args := []string{"server", "users", action.String(), identity.String()}
	if action == kopia.KopiaActionRemove {
		return args, nil
	}

	actionArgs, err := info.Provisioner.Password.KopiaArguments()
	if err != nil {
		return args, fmt.Errorf("error occurred while generating Kopia arguments for identity %s: %w", identity, err)
	}

	return append(args, actionArgs...), nil
}

func (ctx UserActionContext) report(identity identity.Identity, c UserActionClassification, err error) {
	if ctx.isVerbose() || c != UserActionClassificationNone || err != nil {
		fmt.Printf("%s\t%s\n", c, identity)
	}

	if err != nil {
		log.Println(err)
	}
}

func (ctx UserActionContext) Execute(action UserAction) error {
	ids := ctx.Identities.MakeEntries()
	slices.SortFunc(ids, identity.IdentityEntry.Compare)

	for _, entry := range ids {
		cr := ctx.classify(action, entry)

		err := cr.verify(entry.Identity, entry.Info)
		if err == nil && !ctx.isDryRun() {
			switch cr.Classification {
			case UserActionClassificationAdd:
				err = ctx.addIdentity(entry.Identity, entry.Info)
			case UserActionClassificationUpdate:
				err = ctx.updateIdentity(entry.Identity, entry.Info)
			case UserActionClassificationDelete:
				err = ctx.deleteIdentity(entry.Identity, entry.Info)
			}
		}

		if err != nil && cr.Classification&UserActionClassificationSkip == 0 {
			cr.Classification |= UserActionClassificationError
		}

		ctx.report(entry.Identity, cr.Classification, err)
	}

	return nil
}

func (cr UserActionClassificationResult) verify(identity identity.Identity, info identity.IdentityInfo) error {
	if cr.Reason != nil {
		return cr.Reason
	}
	if cr.Classification != UserActionClassificationAdd && cr.Classification != UserActionClassificationUpdate {
		return nil
	}

	if _, err := info.Provisioner.Password.KopiaArguments(); err != nil {
		return fmt.Errorf("error occurred while generating Kopia arguments for identity %s: %w", identity, err)
	}

	return nil
}

func (c UserActionClassification) String() string {
	base := c &^ (UserActionClassificationSkip | UserActionClassificationError)

	var a string
	switch base {
	case UserActionClassificationAdd:
		a = "AI"
	case UserActionClassificationDelete:
		a = "DI"
	case UserActionClassificationUpdate:
		a = "UI"
	default:
		a = "-I"
	}

	if c&UserActionClassificationSkip != 0 {
		return "-" + a
	}
	if c&UserActionClassificationError != 0 {
		return "!" + a
	}

	return " " + a
}
