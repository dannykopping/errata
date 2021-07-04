package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dannykopping/errata"
	"github.com/dannykopping/errata/sample/http"
	"github.com/dannykopping/errata/sample/shell"
)

func main() {
	db, err := readDatabaseFromFile("errata.yml")
	if err != nil {
		log.Fatal(err)
	}

	if err := errata.RegisterSource(db); err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 1 {
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

func readDatabaseFromFile(file string) (*errata.Database, error) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("db open error: %s\n", err)
		return nil, errata.DatabaseFileOpen
	}

	db, err := errata.Parse(f)
	if err != nil {
		fmt.Printf("db parse error: %s\n", err)
		return nil, errata.DatabaseFileParse
	}

	return db, nil
}
