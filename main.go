package main

import (
	pb "github.com/jon-wade/oriServer/ori"
	"github.com/jon-wade/oriServer/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

// TODO: pop this into an external config
const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMathHelperServer(s, &server.server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}


