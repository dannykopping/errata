package backend

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	Address = ":50051"
)

type server struct {
	UnimplementedAuthenticatorServer
}

func (s *server) Authenticate(ctx context.Context, in *AuthenticationRequest) (*AuthenticationResponse, error) {
	log.Printf("Received: %v %v", in.GetUsername(), in.GetPassword())
	return &AuthenticationResponse{Message: fmt.Sprintf("Howzit %s", in.GetUsername())}, nil
}

func Start() {
	lis, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterAuthenticatorServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}