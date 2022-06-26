package http

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dannykopping/errata/sample/errata"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorResponses(t *testing.T) {
	server := NewServer()

	requests := []struct {
		email              string
		password           string
		expectedStatus     int
		expectedErrataCode string
		expectedBody       string
	}{
		{"valid@email.com", "password", fiber.StatusOK, "", "Logged in successfully as: valid@email.com"},
		{"valid@email.com", "wrong", fiber.StatusForbidden, errata.ErrIncorrectPassword, jsonBodyResponse(errata.ErrIncorrectPassword)},
		{"valid@email.com", "", fiber.StatusBadRequest, errata.ErrMissingValues, jsonBodyResponse(errata.ErrMissingValues)},
		{"", "", fiber.StatusBadRequest, errata.ErrMissingValues, jsonBodyResponse(errata.ErrMissingValues)},
		{"", "pass", fiber.StatusBadRequest, errata.ErrMissingValues, jsonBodyResponse(errata.ErrMissingValues)},
		{"spam@email.com", "1234", fiber.StatusForbidden, errata.ErrAccountBlockedSpam, jsonBodyResponse(errata.ErrAccountBlockedSpam)},
		{"abuse@email.com", "1234", fiber.StatusForbidden, errata.ErrAccountBlockedAbuse, jsonBodyResponse(errata.ErrAccountBlockedAbuse)},
		{"abuse@email.com", "1234", fiber.StatusForbidden, errata.ErrAccountBlockedAbuse, jsonBodyResponse(errata.ErrAccountBlockedAbuse)},
		{"invalid.email", "1234", fiber.StatusBadRequest, errata.ErrInvalidEmail, jsonBodyResponse(errata.ErrInvalidEmail)},
		{"missing@email.com", "1234", fiber.StatusForbidden, errata.ErrIncorrectEmail, jsonBodyResponse(errata.ErrIncorrectEmail)},
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
