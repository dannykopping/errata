package main

import (
	"github.com/dannykopping/errata/sample/backend"
	"github.com/dannykopping/errata/sample/frontend"
)

func main() {
	go backend.Start()
	frontend.Start()
}
