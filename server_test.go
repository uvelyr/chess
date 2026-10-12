package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	type testCase struct{
		number int
		expected int
		err error
	}
	t.Run("abs func", func(t *testing.T) {
		tests := []testCase{
			{number: 5, expected: 5},
			{number: 0, expected: 0},
			{number: -5, expected: 5},
		}
		for _, test := range tests {
			actual := abs(test.number)
			assert.Equal(t, test.expected, actual)
		}
	})
}

func TestForward(t *testing.T) {
	type testCase struct{
		color Color
		expected int
		expectedError bool
	}
	t.Run("forward method", func(t *testing.T) {
		tests := []testCase{
			{color: White, expected: -1},
			{color: Black, expected: 1},
			{color: NoColor, expected: 0, expectedError: true},
		}
		for _, test := range tests {
			actual, err := test.color.forward()
			if test.expectedError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expected, actual)
		}
	})
}
			
