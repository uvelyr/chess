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

func TestInBounds(t *testing.T) {
	type testCase struct {
		coord    Coord
		expected bool
	}

	t.Run("InBounds method", func(t *testing.T) {
		tests := []testCase{
			{coord: Coord{X: 0, Y: 0}, expected: true},
			{coord: Coord{X: 7, Y: 7}, expected: true},
			{coord: Coord{X: 4, Y: 3}, expected: true},
			{coord: Coord{X: -1, Y: 0}, expected: false},
			{coord: Coord{X: 0, Y: 8}, expected: false},
			{coord: Coord{X: 8, Y: 8}, expected: false},
		}

		for _, test := range tests {
			actual := test.coord.InBounds()
			assert.Equal(t, test.expected, actual)
		}
	})
}

func TestOpponent(t *testing.T) {
	type testCase struct {
		color    Color
		expected Color
	}

	t.Run("opponent method", func(t *testing.T) {
		tests := []testCase{
			{color: White, expected: Black},
			{color: Black, expected: White},
		}

		for _, test := range tests {
			actual := test.color.opponent()
			assert.Equal(t, test.expected, actual)
		}
	})
}

func TestIsMoveLegal(t *testing.T) {
	type testCase struct {
		target       Coord
		legalTargets []Coord
		expected     bool
	}

	t.Run("isMoveLegal func", func(t *testing.T) {
		tests := []testCase{
			{
				target:       Coord{X: 4, Y: 4},
				legalTargets: []Coord{{X: 4, Y: 4}, {X: 5, Y: 4}},
				expected:     true,
			},
			{
				target:       Coord{X: 5, Y: 4},
				legalTargets: []Coord{{X: 4, Y: 4}, {X: 5, Y: 4}},
				expected:     true,
			},
			{
				target:       Coord{X: 3, Y: 3},
				legalTargets: []Coord{{X: 4, Y: 4}, {X: 5, Y: 4}},
				expected:     false,
			},
			{
				target:       Coord{X: 3, Y: 3},
				legalTargets: nil,
				expected:     false,
			},
		}

		for _, test := range tests {
			actual := isMoveLegal(test.target, test.legalTargets)
			assert.Equal(t, test.expected, actual)
		}
	})
}
			
