package backend

import (
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/dannykopping/errata"
	"github.com/dannykopping/errata/pkg/model"
	"github.com/dannykopping/errata/sample/backend/errors"
	"github.com/gofiber/fiber/v2"
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
	}{
		{"valid@email.com", "1234", 200, ""},
		{"valid@email.com", "wrong", 403, errors.IncorrectPassword},
		{"valid@email.com", "", 400, errors.MissingValues},
		{"", "", 400, errors.MissingValues},
		{"", "pass", 400, errors.MissingValues},
		{"spam@email.com", "1234", 403, errors.AccountBlockedSpam},
		{"abuse@email.com", "1234", 403, errors.AccountBlockedAbuse},
		{"abuse@email.com", "1234", 403, errors.AccountBlockedAbuse},
		{"invalid.email", "1234", 400, errors.InvalidEmail},
		{"missing@email.com", "1234", 403, errors.IncorrectEmail},
	}

	for _, request := range requests {
		form := url.Values{}
		form.Set("email", request.email)
		form.Set("password", request.password)

		req := httptest.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := server.Test(req, 2000)
		require.NoError(t, err)

		require.Equal(t, request.expectedStatus, resp.StatusCode)
		require.Equal(t, request.expectedErrataCode, resp.Header.Get("X-Errata-Code"))
	}
}

func prepareServer(t *testing.T) (*fiber.App, error) {
	// TODO don't duplicate this code
	f, err := os.Open("../errata.yml")
	require.NoError(t, err)

	db, err := model.Parse(f)
	require.NoError(t, err)

	require.NoError(t, errata.RegisterSource(db))

	server := NewServer()
	return server, err
}
