package login

import (
	"strings"

	"github.com/dannykopping/errata/sample/errata"
)

type Request struct {
	EmailAddress string `form:"email"`
	Password     string `form:"password"`
}

// Validate given request against database, returning error if present
func Validate(req Request) error {
	if req.EmailAddress == "" || req.Password == "" {
		return errata.NewMissingValuesErr(nil)
	}

	if strings.Index(req.EmailAddress, "@") < 0 {
		return errata.NewInvalidEmailErr(nil)
	}

	switch req.EmailAddress {
	case "spam@email.com":
		return errata.NewAccountBlockedSpamErr(nil)
	case "abuse@email.com":
		return errata.NewAccountBlockedAbuseErr(nil)
	case "valid@email.com":
		if req.Password != "password" {
			return errata.NewIncorrectPasswordErr(nil)
		}
		return nil
	}

	return errata.NewIncorrectEmailErr(nil)
}
