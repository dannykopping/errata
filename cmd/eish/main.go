package main

import (
	"fmt"
	"os"

	"github.com/dannykopping/errata"
	"github.com/go-kit/log"
	"github.com/urfave/cli/v2"
)

func main() {
	var (
		codeGen errata.CodeGenConfig
		webUI   errata.WebUIConfig
	)

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))

	app := &cli.App{
		Name:  "EISH",
		Usage: "Errata Interactive SHell",
		Authors: []*cli.Author{
			{
				Name:  "Danny Kopping",
				Email: "dannykopping@gmail.com",
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "generate",
				Description: "Generate errata from source",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "source",
						Required:    true,
						Destination: &codeGen.Source,
					},
					&cli.StringFlag{
						Name:        "template",
						Value:       "golang",
						Destination: &codeGen.Template,
					},
					&cli.StringFlag{
						Name:        "package",
						Value:       "errors",
						Destination: &codeGen.Package,
					},
				},
				Action: func(_ *cli.Context) error {
					return errata.Generate(logger, codeGen, os.Stdout)
				},
			},
			{
				Name:     "serve",
				HideHelp: true,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "source",
						Required:    true,
						Destination: &webUI.Source,
					},
					&cli.StringFlag{
						Name:        "bind-addr",
						Value:       "0.0.0.0:37707",
						Destination: &webUI.BindAddr,
					},
				},
				Action: func(_ *cli.Context) error {
					srv, err := errata.NewServer(logger, webUI)
					if err != nil {
						return errata.NewServeWebUiErr(err, webUI.Source)
					}

					return errata.Serve(srv)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
