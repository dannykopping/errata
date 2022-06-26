package http

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/dannykopping/errata/sample/errata"
	"github.com/dannykopping/errata/sample/login"
	"github.com/dannykopping/errata/sample/store"
	"github.com/gofiber/fiber/v2"
)

type Request struct {
	EmailAddress string `form:"email"`
	Password     string `form:"password"`
}

type Response struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	HelpURL string `json:"help_url,omitempty"`
}

func NewServer(store store.Store) *fiber.App {
	app := fiber.New()

	app.Use(errataMiddleware)
	app.Post("/login", func(c *fiber.Ctx) error {
		var req Request

		if err := c.BodyParser(&req); err != nil {
			return errata.NewInvalidRequestErr(err)
		}

		if err := login.Validate(store, req.EmailAddress, req.Password); err != nil {
			return err
		}

		resp := Response{
			Message: fmt.Sprintf("Logged in successfully as: %s", req.EmailAddress),
		}

		body, err := json.Marshal(resp)
		if err != nil {
			return errata.NewResponseFormattingErr(err)
		}

		return c.SendString(string(body))
	})

	return app
}

type HasHTTPResponseCode interface {
	errata.Erratum
	GetHttpResponseCode() string
}

func errataMiddleware(c *fiber.Ctx) error {
	err := c.Next()

	if err == nil {
		return nil
	}

	statusCode := fiber.StatusInternalServerError
	if e, ok := err.(HasHTTPResponseCode); ok {
		if code, cerr := strconv.Atoi(e.GetHttpResponseCode()); cerr == nil {
			statusCode = code
		}
	}

	resp := Response{
		Message: err.Error(),
	}

	if e, ok := err.(errata.Erratum); ok {
		c.Response().Header.Add("X-Errata-Code", e.Code())

		resp.Code = e.Code()
		resp.HelpURL = e.HelpURL()
	}

	body, err := json.Marshal(resp)
	if err != nil {
		body = []byte("could not marshal body")
	}

	return fiber.NewError(statusCode, string(body))
}
