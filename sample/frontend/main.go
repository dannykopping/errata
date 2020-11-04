package frontend

import (
	"context"
	"github.com/dannykopping/errata/sample/backend"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	Address = ":8080"
)

func Start() {
	s := BuildServer()
	err := s.Start(Address)
	if err != nil {
		panic(err)
	}
}

func BuildServer() *echo.Echo {
	server := echo.New()
	server.Use(middleware.Recover())

	/* API */
	server.POST("/auth", authHandler)

	return server
}

func authHandler(c echo.Context) (err error) {
	username := c.FormValue("username")
	password := c.FormValue("password")

	return c.String(200, authenticateWithBackend(username, password).GetMessage())
}

func authenticateWithBackend(username, password string) *backend.AuthenticationResponse {
	conn, err := grpc.Dial(backend.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := backend.NewAuthenticatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Authenticate(ctx, &backend.AuthenticationRequest{Username: username, Password: password})
	if err != nil {
		log.Fatalf("could not authenticate: %v", err)
	}

	return r
}