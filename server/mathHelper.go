package server

import (
	"context"
	"fmt"
	"github.com/JohnCGriffin/overflow"
	pb "github.com/jon-wade/oriServer/ori"
	"github.com/jon-wade/oriServer/server/helpers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MathHelperServer struct {
	pb.MathHelperServer
}

var summationError = status.Errorf(codes.OutOfRange, "summation result exceeds maximum integer value")
var factorialError = status.Errorf(codes.OutOfRange, "factorial result exceeds maximum integer value")

func (s *MathHelperServer) Summation(_ context.Context, in *pb.SummationInput) (*pb.SummationResult, error) {
	result, ok := overflow.Add64(in.First, in.Last)
	if !ok {
		fmt.Printf("summation error: %s\n", summationError.Error())
		return nil, summationError
	}

	resp := &pb.SummationResult{Result: result, First: in.First, Last: in.Last}
	fmt.Printf("summation resp: %v\n", resp)

	return resp, nil
}

func (s *MathHelperServer) Factorial(_ context.Context, in *pb.FactorialInput) (*pb.FactorialResult, error) {
	result, ok := helpers.Factorial(in.Base)
	if !ok {
		fmt.Printf("factorial error: %s\n", factorialError.Error())
		return nil, factorialError
	}

	resp := &pb.FactorialResult{Result: result, Base: in.Base}
	fmt.Printf("factorial resp: %v\n", resp)

	return resp, nil
}