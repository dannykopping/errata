package shell

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dannykopping/errata/sample/errata"
	"github.com/dannykopping/errata/sample/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestErrorResponses(t *testing.T) {
	store, err := store.NewUsersStore("../users.sqlite3")
	require.NoError(t, err)

	app := NewApp(store)

	requests := []struct {
		email              string
		password           string
		expectedExitCode   int
		expectedErrataCode string
		expectedStdout     string
	}{
		{"valid@email.com", "password", SuccessCode, "", "Logged in successfully as: valid@email.com"},
		{"valid@email.com", "wrong", UnsuccessfulCode, errata.IncorrectCredentialsErrCode, errata.IncorrectCredentialsErrCode},
		{"valid@email.com", "", InvalidCode, errata.MissingValuesErrCode, errata.MissingValuesErrCode},
		{"", "", InvalidCode, errata.MissingValuesErrCode, errata.MissingValuesErrCode},
		{"", "pass", InvalidCode, errata.MissingValuesErrCode, errata.MissingValuesErrCode},
		{"spam@email.com", "password", BlockedCode, errata.AccountBlockedSpamErrCode, errata.AccountBlockedSpamErrCode},
		{"abuse@email.com", "password", BlockedCode, errata.AccountBlockedAbuseErrCode, errata.AccountBlockedAbuseErrCode},
		{"abuse@email.com", "password", BlockedCode, errata.AccountBlockedAbuseErrCode, errata.AccountBlockedAbuseErrCode},
		{"invalid.email", "password", InvalidCode, errata.InvalidEmailErrCode, errata.InvalidEmailErrCode},
		{"missing@email.com", "password", UnsuccessfulCode, errata.IncorrectCredentialsErrCode, errata.IncorrectCredentialsErrCode},
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
			assert.Contains(t, err.Error(), request.expectedStdout)
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
