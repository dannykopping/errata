package shell

import (
	"errors"
	"fmt"

	"github.com/dannykopping/errata/sample/errata"
	"github.com/dannykopping/errata/sample/login"
	"github.com/urfave/cli/v2"
)

var (
	request login.Request
)

const (
	SuccessCode      = 0
	InvalidCode      = 1
	UnsuccessfulCode = 2
	BlockedCode      = 3

	UnhandledErrorCode = 127
)

func NewApp() *cli.App {
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
						Destination: &request.EmailAddress,
					},
					&cli.StringFlag{
						Name:        "password",
						Destination: &request.Password,
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
	err := login.Validate(request)

	if err != nil {
		var e errata.Error
		if errors.As(err, &e) {
			return cli.Exit(fmt.Sprintf("%s: %q", e.Code, e.Message), e.Interfaces.ShellExitCode)
		}

		return cli.Exit(fmt.Sprintf("unhandled error: %s", e), UnhandledErrorCode)

	}

	return cli.Exit(fmt.Sprintf("Logged in successfully as: %s", request.EmailAddress), SuccessCode)
}
