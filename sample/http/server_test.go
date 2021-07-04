package http

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"net/url"
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
	require.NoError(t, err)

	requests := []struct {
		email              string
		password           string
		expectedStatus     int
		expectedErrataCode string
		expectedBody       string
	}{
		{"valid@email.com", "1234", fiber.StatusOK, "", "Logged in successfully as: valid@email.com"},
		{"valid@email.com", "wrong", fiber.StatusForbidden, errors.IncorrectPassword, jsonBodyResponse(errors.IncorrectPassword)},
		{"valid@email.com", "", fiber.StatusBadRequest, errors.MissingValues, jsonBodyResponse(errors.MissingValues)},
		{"", "", fiber.StatusBadRequest, errors.MissingValues, jsonBodyResponse(errors.MissingValues)},
		{"", "pass", fiber.StatusBadRequest, errors.MissingValues, jsonBodyResponse(errors.MissingValues)},
		{"spam@email.com", "1234", fiber.StatusForbidden, errors.AccountBlockedSpam, jsonBodyResponse(errors.AccountBlockedSpam)},
		{"abuse@email.com", "1234", fiber.StatusForbidden, errors.AccountBlockedAbuse, jsonBodyResponse(errors.AccountBlockedAbuse)},
		{"abuse@email.com", "1234", fiber.StatusForbidden, errors.AccountBlockedAbuse, jsonBodyResponse(errors.AccountBlockedAbuse)},
		{"invalid.email", "1234", fiber.StatusBadRequest, errors.InvalidEmail, jsonBodyResponse(errors.InvalidEmail)},
		{"missing@email.com", "1234", fiber.StatusForbidden, errors.IncorrectEmail, jsonBodyResponse(errors.IncorrectEmail)},
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
	db, err := errata.NewFileDatasource("../errata.yml")
	require.NoError(t, err)

	assert.NoError(t, errata.RegisterDataSource(db))

	server := NewServer()
	return server, nil
}
