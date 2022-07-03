package http

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/dannykopping/errata/sample/errata"
	"github.com/dannykopping/errata/sample/exec"
	"github.com/gofiber/fiber/v2"
)

type Request struct {
	Command string `form:"command"`
	Args    string `form:"args"`
}

func NewServer() *fiber.App {
	app := fiber.New()

	app.Post("/exec", func(c *fiber.Ctx) error {
		var req Request

		if err := c.BodyParser(&req); err != nil {
			//return errata.NewInvalidRequestErr(err)
			return err
		}

		statusCode := fiber.StatusOK

		res, err := exec.Execute(req.Command, strings.Split(req.Args, " "))

		// if an error occurs, use its HTTP response code (if available) and set the errata code
		if err != nil {

			// if the error is an errata error with an "http_response_code" label (with a GetHttpResponseCode() getter), use its defined HTTP code
			statusCode = fiber.StatusInternalServerError
			if e, ok := err.(HTTPResponder); ok {
				if code, cerr := strconv.Atoi(e.GetHttpResponseCode()); cerr == nil {
					statusCode = code
				}
			}

			if e, ok := err.(errata.Erratum); ok {
				c.Response().Header.Add("X-Errata-Code", e.Code())
			}
		}

		body, _ := json.Marshal(res)
		c.Response().SetStatusCode(statusCode)

		return c.Send(body)
	})

	return app
}

type HTTPResponder interface {
	errata.Erratum
	GetHttpResponseCode() string
}
