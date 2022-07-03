package main

import (
	"log"

	"github.com/dannykopping/errata/sample/http"
)

func main() {
	server := http.NewServer()
	log.Fatal(server.Listen(":8080"))
}
