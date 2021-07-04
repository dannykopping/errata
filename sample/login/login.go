package login

import (
	"strings"

	"github.com/dannykopping/errata"
	"github.com/dannykopping/errata/sample/errors"
)

type Request struct {
	EmailAddress string `form:"email"`
	Password     string `form:"password"`
}

var database = map[string]map[string]string{
	"spam@email.com": {
		"1234": errors.AccountBlockedSpam,
	},
	"abuse@email.com": {
		"1234": errors.AccountBlockedAbuse,
	},
	"valid@email.com": {
		"1234": "",
	},
}

// Validate given request against database, returning error if present
func Validate(req Request) error {
	code := validate(req)
	if code == "" {
		return nil
	}

	return errata.New(code)
}

func validate(req Request) string {
	if req.EmailAddress == "" || req.Password == "" {
		return errors.MissingValues
	}

	if strings.Index(req.EmailAddress, "@") < 0 {
		return errors.InvalidEmail
	}

	if account, found := database[req.EmailAddress]; found {
		if code, found := account[req.Password]; found {
			// valid login, email & password combo found
			return code
		}

		return errors.IncorrectPassword
	}

	return errors.IncorrectEmail
}
