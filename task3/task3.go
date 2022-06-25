package main

import (
	"fmt"
	"time"
)

func main() {
	incomimgTime := "12:55:25AM"

	fmt.Println(TimeConverter(incomimgTime))
}
func TimeConverter(clock string) string {
	layoutIn := "03:04:05PM"
	layoutOut := "15:04:05"

	result, err := time.Parse(layoutIn, clock)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return result.Format(layoutOut)
}
