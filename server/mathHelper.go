package server

import (
	"context"
	"github.com/JohnCGriffin/overflow"
	pb "github.com/jon-wade/oriServer/ori"
	"github.com/jon-wade/oriServer/server/helpers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MathHelperServer struct {
	pb.MathHelperServer
}

func (s *MathHelperServer) Summation(_ context.Context, in *pb.SummationInput) (*pb.SummationResult, error) {
	result, ok := overflow.Add64(in.First, in.Last)
	if !ok {
		return nil, status.Errorf(codes.OutOfRange, "summation result exceeds maximum integer value")
	}
	return &pb.SummationResult{Result: result, First: in.First, Last: in.Last}, nil
}

func (s *MathHelperServer) Factorial(_ context.Context, in *pb.FactorialInput) (*pb.FactorialResult, error) {
	result, ok := helpers.Factorial(in.Base)
	if !ok {
		return nil, status.Errorf(codes.OutOfRange, "factorial result exceeds maximum integer value")
	}
	return &pb.FactorialResult{Result: result, Base: in.Base}, nil
}