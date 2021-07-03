package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dannykopping/errata"
	"github.com/dannykopping/errata/sample/backend"
)

func main() {
	f, err := os.Open("errata.yml")
	if err != nil {
		panic(err)
	}

	db, err := errata.Parse(f)
	if err != nil {
		panic(err)
	}

	for _, e := range db.Errors {
		fmt.Printf("err: %q %p\n", e.Code, e)
		if e.External != nil {
			fmt.Printf("\texternal err: %q %p\n", e.External.Code, &e.External)
		}
	}

	if err := errata.RegisterSource(db); err != nil {
		log.Fatal(err)
	}

	server := backend.NewServer()
	log.Fatal(server.Listen(":3000"))
}
