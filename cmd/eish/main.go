package main

import (
	"errors"
	"fmt"
	"os"

	errata "github.com/dannykopping/errata/pkg/errors"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "EISH",
		Usage:    "Errata Interactive SHell",
		HideHelp: true,
		Authors: []*cli.Author{
			{
				Name:  "Danny Kopping",
				Email: "dannykopping@gmail.com",
			},
		},
		Commands: []*cli.Command{
			{
				Name:     "generate",
				HideHelp: true,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "eds.file",
						Required:    true,
						Destination: &edsFile,
					},
					&cli.StringFlag{
						Name:        "language",
						Value:       "golang",
						Destination: &lang,
					},
					&cli.StringFlag{
						Name:        "package",
						Value:       "errors",
						Destination: &pkg,
					},
				},
				Action: generate,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)

		var e errata.ErrorDefinition
		if errors.As(err, &e) {
			// TODO: use exit code defined in definition
			os.Exit(1)
		}
	}
}
