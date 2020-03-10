package helpers_test

import (
	"fmt"
	"github.com/jon-wade/oriServer/server/helpers"
	"testing"
)

func TestFactorial(t *testing.T) {
	var tests = []struct {
		base   int64
		result int64
		ok     bool
	}{
		{3, 6, true},
		{5, 120, true},
		{21, 0, false},
		{-5, 0, false},
	}

	for _, testData := range tests {
		testName := fmt.Sprintf("base=%d,result=%d,ok=%v", testData.base, testData.result, testData.ok)
		t.Run(testName, func(t *testing.T) {
			result, ok := helpers.Factorial(testData.base)
			if ok != testData.ok {
				t.Errorf("Expected ok=%v, got %v", testData.ok, ok)
			}
			if result != testData.result {
				t.Errorf("Expected result=%d, got %d", testData.result, result)
			}
		})
	}
}
