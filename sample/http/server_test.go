package http

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/dannykopping/errata"
	"github.com/dannykopping/errata/sample/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorResponses(t *testing.T) {
	server, err := prepareServer(t)
	assert.NoError(t, err)

	requests := []struct {
		email              string
		password           string
		expectedStatus     int
		expectedErrataCode string
		expectedBody       string
	}{
		{"valid@email.com", "1234", 200, "", "Logged in successfully as: valid@email.com"},
		{"valid@email.com", "wrong", 403, errors.IncorrectPassword, jsonBodyResponse(errors.IncorrectPassword)},
		{"valid@email.com", "", 400, errors.MissingValues, jsonBodyResponse(errors.MissingValues)},
		{"", "", 400, errors.MissingValues, jsonBodyResponse(errors.MissingValues)},
		{"", "pass", 400, errors.MissingValues, jsonBodyResponse(errors.MissingValues)},
		{"spam@email.com", "1234", 403, errors.AccountBlockedSpam, jsonBodyResponse(errors.AccountBlockedSpam)},
		{"abuse@email.com", "1234", 403, errors.AccountBlockedAbuse, jsonBodyResponse(errors.AccountBlockedAbuse)},
		{"abuse@email.com", "1234", 403, errors.AccountBlockedAbuse, jsonBodyResponse(errors.AccountBlockedAbuse)},
		{"invalid.email", "1234", 400, errors.InvalidEmail, jsonBodyResponse(errors.InvalidEmail)},
		{"missing@email.com", "1234", 403, errors.IncorrectEmail, jsonBodyResponse(errors.IncorrectEmail)},
	}

	for _, request := range requests {
		form := url.Values{}
		form.Set("email", request.email)
		form.Set("password", request.password)

		req := httptest.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := server.Test(req, 2000)
		require.NoError(t, err)

		assert.Equal(t, request.expectedStatus, resp.StatusCode)
		assert.Equal(t, request.expectedErrataCode, resp.Header.Get("X-Errata-Code"))
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, request.expectedBody, string(body))
	}
}

func jsonBodyResponse(code string) string {
	val := struct {
		Code string `json:"code"`
	}{
		Code: code,
	}

	r, _ := json.Marshal(&val)
	return string(r)
}

func prepareServer(t *testing.T) (*fiber.App, error) {
	// TODO don't duplicate this code
	f, err := os.Open("../errata.yml")
	assert.NoError(t, err)

	db, err := errata.Parse(f)
	assert.NoError(t, err)

	assert.NoError(t, errata.RegisterSource(db))

	server := NewServer()
	return server, err
}
