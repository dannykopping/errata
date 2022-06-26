package main

import (
	"log"
	"os"

	"github.com/dannykopping/errata/sample/http"
	"github.com/dannykopping/errata/sample/shell"
	"github.com/dannykopping/errata/sample/store"
	_ "github.com/glebarez/go-sqlite"
)

func main() {
	if len(os.Args) <= 1 {
		showHelp()
	}

	db, err := store.NewUsersStore("users.sqlite3")
	if err != nil {
		// TODO handle error
		log.Fatal(err)
	}

	mode := os.Args[1]
	switch mode {
	case "http":
		log.Fatal(runHTTP(db))
	case "shell":
		log.Fatal(runShell(db))
	default:
		showHelp()
	}
}

func showHelp() {
	log.Fatal(`select a mode: "http" or "shell"`)
}

func runHTTP(store store.Store) error {
	server := http.NewServer(store)
	return server.Listen(":0")
}

func runShell(store store.Store) error {
	app := shell.NewApp(store)
	return app.Run(os.Args[1:])
}
