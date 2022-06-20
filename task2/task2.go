package main

import (
	"sort"
)

func MinMaxSumOfSlice(incomingSlice []int) (MinSum, MaxSum int) {
	sort.Ints(incomingSlice)
	MinSum, MaxSum = 0, 0

	for i := 0; i < len(incomingSlice)-1; i++ {
		MinSum += incomingSlice[i]
	}
	for j := 1; j < len(incomingSlice); j++ {
		MaxSum += incomingSlice[j]
	}
	return MinSum, MaxSum
}
