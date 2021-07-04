package main

import (
	"log"
	"os"

	"github.com/dannykopping/errata"
	"github.com/dannykopping/errata/sample/http"
	"github.com/dannykopping/errata/sample/shell"
)

func main() {
	ds, err := errata.NewFileDatasource("errata.yml")
	if err != nil {
		log.Fatal(err)
	}

	if err := errata.RegisterDataSource(ds); err != nil {
		log.Fatal(err)
	}

	if len(os.Args) <= 1 {
		showHelp()
	}

	mode := os.Args[1]
	switch mode {
	case "http":
		log.Fatal(runHTTP())
	case "shell":
		runShell()
	default:
		showHelp()
	}
}

func showHelp() {
	log.Fatal(`select a mode: "http" or "shell"`)
}

func runHTTP() error {
	server := http.NewServer()
	return server.Listen(":3000")
}

func runShell() error {
	app := shell.NewApp()
	return app.Run(os.Args[1:])
}
