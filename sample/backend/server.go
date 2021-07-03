package backend

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dannykopping/errata"
	"github.com/dannykopping/errata/sample/backend/errors"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
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

func NewServer() *fiber.App {
	app := fiber.New()

	app.Use(errataMiddleware)
	app.Post("/login", func(c *fiber.Ctx) error {
		var req LoginRequest

		err := c.BodyParser(&req)
		if err != nil {
			return errata.New(errors.InvalidRequest)
		}

		if code := validate(req); code != "" {
			return errata.New(code)
		}

		return c.SendString(fmt.Sprintf("Logged in successfully as: %s", req.EmailAddress))
	})

	return app
}

func errataMiddleware(c *fiber.Ctx) error {
	err := c.Next()

	if e, ok := err.(*errata.Error); e != nil && ok {
		statusCode := fiber.StatusInternalServerError
		if e.HTTP != nil {
			statusCode = e.HTTP.Code
		}

		c.Response().Header.Add("X-Errata-Code", e.Code)

		body, err := formatError(e)
		if err != nil {
			fmt.Printf("formatting error: %q\n", err)
			return fiber.NewError(fiber.StatusInternalServerError, errors.ResponseFormattingFailure)
		}

		return fiber.NewError(statusCode, body)
	}

	return err
}

func formatError(e *errata.Error) (string, error) {
	if e == nil {
		return "", nil
	}

	s := struct {
		Code string `json:"code"`
	}{
		Code: e.Code,
	}

	r, err := json.Marshal(&s)
	if err != nil {
		return "", err
	}

	return string(r), nil
}

func validate(req LoginRequest) string {
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
