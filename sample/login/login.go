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
		return errata.NewMissingValues()
	}

	if strings.Index(req.EmailAddress, "@") < 0 {
		return errata.NewInvalidEmail()
	}

	switch req.EmailAddress {
	case "spam@email.com":
		return errata.NewAccountBlockedSpam()
	case "abuse@email.com":
		return errata.NewAccountBlockedAbuse()
	case "valid@email.com":
		return nil
	}

	return errata.NewIncorrectEmail()
}
