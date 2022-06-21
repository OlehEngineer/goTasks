package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please input numbers for slice you are going to check separated by space")
	text, _ := reader.ReadString('\n')
	UserSlice := strings.Fields(text)

	fmt.Println(MinMaxSumOfSlice(inputConverToInt(UserSlice)))
}

// calculate minimum and maximum sum of n-1 slice elements
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

//convert input value to integer and check is value could be converted to integer
func inputConverToInt(UserSlice []string) []int {
	ConvertedSlice := []int{}
	for _, item := range UserSlice {
		intValue, err := strconv.Atoi(item)
		if err == nil {
			ConvertedSlice = append(ConvertedSlice, intValue)
		} else {
			fmt.Printf("value %v is not an integer. Please try input once more", item)
			os.Exit(0)
		}
	}
	return ConvertedSlice
}
