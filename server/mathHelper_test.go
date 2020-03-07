package server_test

import (
	"context"
	"fmt"
	pb "github.com/jon-wade/oriServer/ori"
	"github.com/jon-wade/oriServer/server"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"testing"
)

func TestMathHelperServer_Factorial(t *testing.T) {
	var tests = []struct {
		in     *pb.FactorialInput
		result *pb.FactorialResult
		err    error
	}{
		{
			&pb.FactorialInput{Base: 3},
			&pb.FactorialResult{Result: 6, Base: 3},
			nil,
		},
		{
			&pb.FactorialInput{Base: 5},
			&pb.FactorialResult{Result: 120, Base: 5},
			nil,
		},
		{
			&pb.FactorialInput{Base: 21},
			nil,
			status.Errorf(codes.OutOfRange, "factorial result exceeds maximum integer value"),
		},
		{
			&pb.FactorialInput{Base: -5},
			nil,
			status.Errorf(codes.OutOfRange, "factorial result exceeds maximum integer value"),
		},
	}

	s := server.MathHelperServer{}
	ctx := context.Background()

	for idx, testData := range tests {
		testName := fmt.Sprintf("in=%v,result=%v,err=%v", testData.in, testData.result, testData.err)
		t.Run(testName, func(t *testing.T) {
			result, err := s.Factorial(ctx, testData.in)

			switch idx {
			case 0, 1:
				if err != testData.err {
					t.Errorf("Expected err=%v, got %v", testData.err, err)
				}
				if err == nil {
					if result.Base != testData.result.Base {
						t.Errorf("Expected result.Base=%v, got %v", testData.result.Base, result.Base)
					}
					if result.Result != testData.result.Result {
						t.Errorf("Expected result.Base=%v, got %v", testData.result.Result, result.Result)
					}
				}
			case 2, 3:
				if result != testData.result {
					t.Errorf("Expected result=%v, got %v", testData.result, result)
				}
				if err != nil {
					if err.Error() != testData.err.Error() {
						t.Errorf("Expected err=%v, got %v", testData.err, err)
					}
				}
			}
		})
	}
}

func TestMathHelperServer_Summation(t *testing.T) {
	var tests = []struct {
		in     *pb.SummationInput
		result *pb.SummationResult
		err    error
	}{
		{
			&pb.SummationInput{First: 2, Last: 2},
			&pb.SummationResult{Result: 4, First: 2, Last: 2},
			nil,
		},
		{
			&pb.SummationInput{First: 250, Last: -250},
			&pb.SummationResult{Result: 0, First: 250, Last: -250},
			nil,
		},
		{
			&pb.SummationInput{First: math.MaxInt64, Last: 1},
			nil,
			status.Errorf(codes.OutOfRange, "summation result exceeds maximum integer value"),
		},
	}

	s := server.MathHelperServer{}
	ctx := context.Background()

	for idx, testData := range tests {
		testName := fmt.Sprintf("in=%v,result=%v,err=%v", testData.in, testData.result, testData.err)
		t.Run(testName, func(t *testing.T) {
			result, err := s.Summation(ctx, testData.in)

			switch idx {
			case 0, 1:
				if err != testData.err {
					t.Errorf("Expected err=%v, got %v", testData.err, err)
				}
				if err == nil {
					if result.First != testData.result.First {
						t.Errorf("Expected result.First=%v, got %v", testData.result.First, result.First)
					}
					if result.Last != testData.result.Last {
						t.Errorf("Expected result.Last=%v, got %v", testData.result.Last, result.Last)
					}
					if result.Result != testData.result.Result {
						t.Errorf("Expected result.Base=%v, got %v", testData.result.Result, result.Result)
					}
				}
			case 2:
				if result != testData.result {
					t.Errorf("Expected result=%v, got %v", testData.result, result)
				}
				if err != nil {
					if err.Error() != testData.err.Error() {
						t.Errorf("Expected err=%v, got %v", testData.err, err)
					}
				}
			}
		})
	}
}
