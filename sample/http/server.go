package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/dannykopping/errata/sample/errata"
	"github.com/dannykopping/errata/sample/login"
	"github.com/gofiber/fiber/v2"
)

func NewServer() *fiber.App {
	app := fiber.New()

	app.Use(errataMiddleware)
	app.Post("/login", func(c *fiber.Ctx) error {
		var req login.Request

		if err := c.BodyParser(&req); err != nil {
			return errata.NewInvalidRequestErr(err)
		}

		if err := login.Validate(req); err != nil {
			return err
		}

		return c.SendString(fmt.Sprintf("Logged in successfully as: %s", req.EmailAddress))
	})

	return app
}

func errataMiddleware(c *fiber.Ctx) error {
	err := c.Next()

	var e errata.Error
	if err != nil && errors.As(err, &e) {
		statusCode, ex := getHTTPStatusCode(e)
		if ex != nil || statusCode <= 0 {
			statusCode = fiber.StatusInternalServerError
		}

		c.Response().Header.Add("X-Errata-Code", e.Code)

		body, err := formatError(e)
		if err != nil {
			e := err.(errata.Error)
			return fiber.NewError(statusCode, e.Message)
		}

		return fiber.NewError(statusCode, body)
	}

	return err
}

func formatError(e errata.Error) (string, error) {
	s := struct {
		Code string `json:"code"`
	}{
		Code: e.Code,
	}

	r, err := json.Marshal(&s)
	if err != nil {
		return "", errata.NewResponseFormattingErr(err)
	}

	return string(r), nil
}

func getHTTPStatusCode(err errata.Error) (int, error) {
	c, ok := err.Labels["http_response_code"]
	if ok {
		code, e := strconv.Atoi(c)
		if e != nil {
			return 0, e
		}

		return code, nil
	}

	// no exit code defined
	return 0, nil
}
