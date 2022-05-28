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
		return errata.NewMissingValues(nil)
	}

	if strings.Index(req.EmailAddress, "@") < 0 {
		return errata.NewInvalidEmail(nil)
	}

	switch req.EmailAddress {
	case "spam@email.com":
		return errata.NewAccountBlockedSpam(nil)
	case "abuse@email.com":
		return errata.NewAccountBlockedAbuse(nil)
	case "valid@email.com":
		if req.Password != "password" {
			return errata.NewIncorrectPassword(nil)
		}
		return nil
	}

	return errata.NewIncorrectEmail(nil)
}
