package main

import (
	"log"
	"os"

	"github.com/dannykopping/errata/sample/http"
	"github.com/dannykopping/errata/sample/shell"
)

func main() {
	if len(os.Args) <= 1 {
		showHelp()
	}

	mode := os.Args[1]
	switch mode {
	case "http":
		log.Fatal(runHTTP())
	case "shell":
		log.Fatal(runShell())
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
