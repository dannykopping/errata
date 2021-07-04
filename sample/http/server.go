package http

import (
	"encoding/json"
	"fmt"

	"github.com/dannykopping/errata"
	"github.com/dannykopping/errata/sample/errors"
	"github.com/dannykopping/errata/sample/login"
	"github.com/gofiber/fiber/v2"
)

func NewServer() *fiber.App {
	app := fiber.New()

	app.Use(errataMiddleware)
	app.Post("/login", func(c *fiber.Ctx) error {
		var req login.Request

		err := c.BodyParser(&req)
		if err != nil {
			return errata.New(errors.InvalidRequest)
		}

		if code := login.Validate(req); code != "" {
			return errata.New(code)
		}

		return c.SendString(fmt.Sprintf("Logged in successfully as: %s", req.EmailAddress))
	})

	return app
}

func errataMiddleware(c *fiber.Ctx) error {
	err := c.Next()

	if e, ok := err.(*errata.Error); e != nil && ok {
		statusCode := e.HTTPStatusCode(fiber.StatusInternalServerError)
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