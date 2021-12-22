package shell

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dannykopping/errata/sample/errata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestErrorResponses(t *testing.T) {
	app, err := prepareShell(t)
	assert.NoError(t, err)

	requests := []struct {
		email              string
		password           string
		expectedExitCode   int
		expectedErrataCode string
		expectedStdout     string
	}{
		{"valid@email.com", "1234", SuccessCode, "", "Logged in successfully as: valid@email.com"},
		{"valid@email.com", "wrong", UnsuccessfulCode, errata.ErrIncorrectPassword, stdoutResponse(errata.ErrIncorrectPassword)},
		{"valid@email.com", "", InvalidCode, errata.ErrMissingValues, stdoutResponse(errata.ErrMissingValues)},
		{"", "", InvalidCode, errata.ErrMissingValues, stdoutResponse(errata.ErrMissingValues)},
		{"", "pass", InvalidCode, errata.ErrMissingValues, stdoutResponse(errata.ErrMissingValues)},
		{"spam@email.com", "1234", BlockedCode, errata.ErrAccountBlockedSpam, stdoutResponse(errata.ErrAccountBlockedSpam)},
		{"abuse@email.com", "1234", BlockedCode, errata.ErrAccountBlockedAbuse, stdoutResponse(errata.ErrAccountBlockedAbuse)},
		{"abuse@email.com", "1234", BlockedCode, errata.ErrAccountBlockedAbuse, stdoutResponse(errata.ErrAccountBlockedAbuse)},
		{"invalid.email", "1234", InvalidCode, errata.ErrInvalidEmail, stdoutResponse(errata.ErrInvalidEmail)},
		{"missing@email.com", "1234", UnsuccessfulCode, errata.ErrIncorrectEmail, stdoutResponse(errata.ErrIncorrectEmail)},
	}

	for _, request := range requests {
		cli.OsExiter = func(code int) {
			assert.Equal(t, request.expectedExitCode, code)
		}

		app.ExitErrHandler = func(context *cli.Context, err error) {
			// unfortunately, urfave/cli doesn't export its exitError type, so we have to hack the exit code out
			exitCode := int(reflect.ValueOf(err).Elem().FieldByName("exitCode").Int())
			assert.Equal(t, request.expectedExitCode, exitCode)

			assert.Error(t, err)
			assert.Equal(t, request.expectedStdout, err.Error())
		}

		err := app.Run([]string{
			"",
			"login",
			fmt.Sprintf("--email=%s", request.email),
			fmt.Sprintf("--password=%s", request.password),
		})
		assert.Error(t, err)
	}
}

func stdoutResponse(code string) string {
	err := errata.New(code)
	return fmt.Sprintf("%s: %q", err.Code, err.Message)
}

func prepareShell(t *testing.T) (*cli.App, error) {
	db, err := errata.NewFileDatasource("../errata.yml")
	require.NoError(t, err)

	assert.NoError(t, errata.RegisterDataSource(db))

	server := NewApp()
	return server, nil
}
