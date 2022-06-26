package http

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dannykopping/errata/sample/errata"
	"github.com/dannykopping/errata/sample/store"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorResponses(t *testing.T) {
	store, err := store.NewUsersStore("../users.sqlite3")
	require.NoError(t, err)

	server := NewServer(store)

	requests := []struct {
		email              string
		password           string
		expectedStatus     int
		expectedErrataCode string
		expectedBody       string
	}{
		{"valid@email.com", "password", fiber.StatusOK, "", "Logged in successfully as: valid@email.com"},
		{"valid@email.com", "wrong", fiber.StatusForbidden, errata.IncorrectCredentialsErrCode, jsonBodyResponse(errata.IncorrectCredentialsErrCode)},
		{"valid@email.com", "", fiber.StatusBadRequest, errata.MissingValuesErrCode, jsonBodyResponse(errata.MissingValuesErrCode)},
		{"", "", fiber.StatusBadRequest, errata.MissingValuesErrCode, jsonBodyResponse(errata.MissingValuesErrCode)},
		{"", "pass", fiber.StatusBadRequest, errata.MissingValuesErrCode, jsonBodyResponse(errata.MissingValuesErrCode)},
		{"spam@email.com", "password", fiber.StatusForbidden, errata.AccountBlockedSpamErrCode, jsonBodyResponse(errata.AccountBlockedSpamErrCode)},
		{"abuse@email.com", "password", fiber.StatusForbidden, errata.AccountBlockedAbuseErrCode, jsonBodyResponse(errata.AccountBlockedAbuseErrCode)},
		{"abuse@email.com", "password", fiber.StatusForbidden, errata.AccountBlockedAbuseErrCode, jsonBodyResponse(errata.AccountBlockedAbuseErrCode)},
		{"invalid.email", "password", fiber.StatusBadRequest, errata.InvalidEmailErrCode, jsonBodyResponse(errata.InvalidEmailErrCode)},
		{"missing@email.com", "password", fiber.StatusForbidden, errata.IncorrectCredentialsErrCode, jsonBodyResponse(errata.IncorrectCredentialsErrCode)},
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
		assert.Contains(t, string(body), request.expectedErrataCode)
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
