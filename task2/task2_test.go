package task2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinMaxSumOfSlice(t *testing.T) {
	testingInput := [][]int{
		{5, 9, 7, 2, 3, 8},
		{1, 2, 3, 4, 5},
		{4, 9, 8, 7, 6, 5},
	}
	testingOutput := [][]int{
		{25, 32}, {10, 14}, {30, 35},
	}
	for i := range testingInput {
		min, max := MinMaxSumOfSlice(testingInput[i])
		assert.Equal(t, testingOutput[i][0], min)
		assert.Equal(t, testingOutput[i][1], max)
	}

}
