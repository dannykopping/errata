package shell

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/dannykopping/errata"
	"github.com/dannykopping/errata/sample/errors"
	"github.com/stretchr/testify/assert"
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
		{"valid@email.com", "wrong", UnsuccessfulCode, errors.IncorrectPassword, stdoutResponse(errors.IncorrectPassword)},
		{"valid@email.com", "", InvalidCode, errors.MissingValues, stdoutResponse(errors.MissingValues)},
		{"", "", InvalidCode, errors.MissingValues, stdoutResponse(errors.MissingValues)},
		{"", "pass", InvalidCode, errors.MissingValues, stdoutResponse(errors.MissingValues)},
		{"spam@email.com", "1234", BlockedCode, errors.AccountBlockedSpam, stdoutResponse(errors.AccountBlockedSpam)},
		{"abuse@email.com", "1234", BlockedCode, errors.AccountBlockedAbuse, stdoutResponse(errors.AccountBlockedAbuse)},
		{"abuse@email.com", "1234", BlockedCode, errors.AccountBlockedAbuse, stdoutResponse(errors.AccountBlockedAbuse)},
		{"invalid.email", "1234", InvalidCode, errors.InvalidEmail, stdoutResponse(errors.InvalidEmail)},
		{"missing@email.com", "1234", UnsuccessfulCode, errors.IncorrectEmail, stdoutResponse(errors.IncorrectEmail)},
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
	// TODO don't duplicate this code
	f, err := os.Open("../errata.yml")
	assert.NoError(t, err)

	db, err := errata.Parse(f)
	assert.NoError(t, err)

	assert.NoError(t, errata.RegisterSource(db))

	server := NewApp()
	return server, err
}
