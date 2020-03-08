package main

import (
	pb "github.com/jon-wade/oriServer/ori"
	"github.com/jon-wade/oriServer/server"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	lis, err := net.Listen("tcp", os.Getenv("ORI_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMathHelperServer(s, &server.MathHelperServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
