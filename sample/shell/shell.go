package shell

import (
	"fmt"

	"github.com/dannykopping/errata"
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
		%d: blocked`, InvalidCode, UnsuccessfulCode, BlockedCode)
}

func loginAction(c *cli.Context) error {
	err, ok := login.Validate(request).(errata.Error)
	if !ok {
		return cli.Exit(fmt.Sprintf("Logged in successfully as: %s", request.EmailAddress), SuccessCode)
	}

	return cli.Exit(fmt.Sprintf("%s: %q", err.Code, err.Message), err.ShellExitCode(1))
}
