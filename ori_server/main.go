package main

import (
	"context"
	"github.com/JohnCGriffin/overflow"
	pb "github.com/jon-wade/oriServer/ori"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

// TODO: pop this into an external config
const (
	port = ":50051"
)

func fact(n int64) (int64, bool) {
	var total = n
	var ok bool
	for i := n - 1; i > 1; i-- {
		total, ok = overflow.Mul64(total, i)
		if !ok {
			return 0, ok
		}
	}
	return total, ok
}

type server struct {
	pb.MathHelperServer
}

func (s *server) Summation(_ context.Context, in *pb.SummationInput) (*pb.SummationResult, error) {
	result, ok := overflow.Add64(in.First, in.Last)
	if !ok {
		return nil, status.Errorf(codes.OutOfRange, "summation result exceeds maximum integer value")
	}
	return &pb.SummationResult{Result: result, First: in.First, Last: in.Last}, nil
}

func (s *server) Factorial(_ context.Context, in *pb.FactorialInput) (*pb.FactorialResult, error) {
	result, ok := fact(in.Base)
	if !ok {
		return nil, status.Errorf(codes.OutOfRange, "factorial result exceeds maximum integer value")
	}
	return &pb.FactorialResult{Result: result, Base: in.Base}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMathHelperServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}


