package shell

import (
	"fmt"
	"strconv"

	"github.com/dannykopping/errata/sample/errata"
	"github.com/dannykopping/errata/sample/login"
	"github.com/dannykopping/errata/sample/store"
	"github.com/urfave/cli/v2"
)

var (
	email    string
	password string

	db store.Store
)

const (
	SuccessCode      = 0
	InvalidCode      = 1
	UnsuccessfulCode = 2
	BlockedCode      = 3

	UnhandledErrorCode = 127
)

func NewApp(store store.Store) *cli.App {
	db = store

	return &cli.App{
		Commands: []*cli.Command{
			{
				Description: "Sample login application",
				UsageText:   usageText(),
				Name:        "login",
				Action:      loginAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "email",
						Destination: &email,
					},
					&cli.StringFlag{
						Name:        "password",
						Destination: &password,
					},
				},
			},
		},
	}
}

func usageText() string {
	return fmt.Sprintf(`possible exit codes:
		%d: invalid
		%d: unsuccessful
		%d: blocked
		%d: unhandled`, InvalidCode, UnsuccessfulCode, BlockedCode, UnhandledErrorCode)
}

func loginAction(_ *cli.Context) error {
	// attempt login
	err := login.Validate(db, email, password)

	if err == nil {
		return cli.Exit(fmt.Sprintf("Logged in successfully as: %s", email), SuccessCode)
	}

	exitCode := UnhandledErrorCode
	if e, ok := err.(HasShellExitCode); ok {
		if code, cerr := strconv.Atoi(e.GetShellExitCode()); cerr == nil {
			exitCode = code
		}
	}

	return cli.Exit(err.Error(), exitCode)
}

type HasShellExitCode interface {
	errata.Erratum
	GetShellExitCode() string
}
