package login

import (
	"strings"

	"github.com/dannykopping/errata/sample/errata"
	"github.com/dannykopping/errata/sample/store"
)

// Validate given request against database, returning error if present
func Validate(store store.Store, email, password string) error {
	if email == "" {
		return errata.NewMissingValuesErr(nil, "email")
	}
	if password == "" {
		return errata.NewMissingValuesErr(nil, "password")
	}

	if strings.Index(email, "@") < 0 {
		return errata.NewInvalidEmailErr(nil, email)
	}

	user, err := store.GetUser(email, password)
	if err != nil {
		return errata.NewIncorrectCredentialsErr(err)
	}

	switch true {
	case user.Abuse:
		return errata.NewAccountBlockedAbuseErr(nil)
	case user.Spam:
		return errata.NewAccountBlockedSpamErr(nil)
	}

	return nil
}
