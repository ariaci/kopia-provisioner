package actions

import (
	"fmt"
	"kopia-provisioner/identity"
	"kopia-provisioner/kopia"
	"log"
)

type UserActionFlags uint8

const UserActionFlagNone UserActionFlags = 0
const (
	UserActionFlagDryRun UserActionFlags = 1 << iota
	UserActionFlagUpdate
	UserActionFlagVerbose
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
	UserActionClassificationError = 1 << 7
)

type UserAction uint8

const (
	UserActionSync UserAction = iota
	UserActionAdd
	UserActionRemove
	UserActionUpdate
)

func (context UserActionContext) hasConfigFile() bool {
	return len(context.ConfigFile) > 0
}

func (context UserActionContext) isDryRun() bool {
	return context.Flags&UserActionFlagDryRun != 0
}

func (context UserActionContext) shouldUpdate() bool {
	return context.Flags&UserActionFlagUpdate != 0
}

func (context UserActionContext) isVerbose() bool {
	return context.Flags&UserActionFlagVerbose != 0
}

func (context UserActionContext) classify(action UserAction, info identity.IdentityInfo) UserActionClassification {
	switch {
	// Case 1) Identity only available in kopia -> remove Identity from kopia
	case (action == UserActionSync || action == UserActionRemove) && info.HasKopia() && !info.HasProvisioner():
		return UserActionClassificationDelete
	// Case 2) Identity only available in provisioner -> add Identity to kopia
	case (action == UserActionSync || action == UserActionAdd) && !info.HasKopia() && info.HasProvisioner():
		return UserActionClassificationAdd
	// Case 3) Identity available in kopia and provisioner -> update Identity in kopia
	case (action == UserActionUpdate || context.shouldUpdate()) && info.HasKopia() && info.HasProvisioner():
		return UserActionClassificationUpdate

	default:
		return UserActionClassificationNone
	}
}

func (context UserActionContext) kopiaArguments(action kopia.KopiaAction, identity identity.Identity, info identity.IdentityInfo) ([]string, error) {
	args := []string{"server", "users", action.String(), identity.String()}
	if action == kopia.KopiaActionRemove {
		return args, nil
	}

	actionArgs, err := info.Provisioner.Password.KopiaArguments()
	if err != nil {
		return args, fmt.Errorf("error occurred while generating Kopia arguments for identity %s: %w", identity, err)
	}

	args = append(args, actionArgs...)
	return args, nil
}

func (context UserActionContext) report(identity identity.Identity, c UserActionClassification, err error) {
	if context.isVerbose() || c != UserActionClassificationNone || err != nil {
		fmt.Printf("%s\t%s\n", c, identity)
	}

	if err != nil {
		log.Println(err)
	}
}

func (context UserActionContext) Execute(action UserAction) error {
	for identity, info := range context.Identities {
		c := context.classify(action, info)

		err := c.verify(identity, info)
		if err == nil && !context.isDryRun() {
			switch c {
			case UserActionClassificationAdd:
				err = context.addIdentity(identity, info)
			case UserActionClassificationUpdate:
				err = context.updateIdentity(identity, info)
			case UserActionClassificationDelete:
				err = context.deleteIdentity(identity, info)
			}
		}

		if err != nil {
			c |= UserActionClassificationError
		}

		context.report(identity, c, err)
	}

	return nil
}

func (c UserActionClassification) verify(identity identity.Identity, info identity.IdentityInfo) error {
	if c != UserActionClassificationAdd && c != UserActionClassificationUpdate {
		return nil
	}

	if _, err := info.Provisioner.Password.KopiaArguments(); err != nil {
		return fmt.Errorf("error occurred while generating Kopia arguments for identity %s: %w", identity, err)
	}

	return nil
}

func (c UserActionClassification) String() string {
	base := c &^ UserActionClassificationError

	var a string
	switch base {
	case UserActionClassificationAdd:
		a = "A"
	case UserActionClassificationDelete:
		a = "D"
	case UserActionClassificationUpdate:
		a = "U"
	default:
		a = "-"
	}

	if c&UserActionClassificationError != 0 {
		return "!" + a
	}

	return " " + a
}
