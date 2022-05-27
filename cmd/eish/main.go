package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/dannykopping/errata"
	"github.com/urfave/cli/v2"
)

func main() {
	var codeGen errata.CodeGen

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
						Destination: &codeGen.File,
					},
					&cli.StringFlag{
						Name:        "language",
						Value:       "golang",
						Destination: &codeGen.Lang,
					},
					&cli.StringFlag{
						Name:        "package",
						Value:       "errors",
						Destination: &codeGen.Package,
					},
				},
				Action: func(_ *cli.Context) error {
					return errata.Generate(codeGen, os.Stdout)
				},
			},
			{
				Name:     "serve",
				HideHelp: true,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "eds.file",
						Required:    true,
						Destination: &codeGen.File,
					},
					&cli.StringFlag{
						Name:        "language",
						Value:       "golang",
						Destination: &codeGen.Lang,
					},
					&cli.StringFlag{
						Name:        "package",
						Value:       "errors",
						Destination: &codeGen.Package,
					},
				},
				Action: func(_ *cli.Context) error {
					// TODO: don't use CodeGen type
					return errata.Serve(&errata.Server{
						File:    codeGen.File,
						Package: codeGen.Package,
					})
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		var e errata.Error
		if errors.As(err, &e) {
			if code := e.Interfaces.ShellExitCode; code > 0 {
				os.Exit(code)
			}
		}

		os.Exit(1)
	}
}
