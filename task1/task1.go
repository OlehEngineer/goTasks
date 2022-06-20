package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

//timer start to measure 10 seconds
//after 10 seconds Goodbye message appear and program exit
func main() {
	fmt.Println("Hello World")
	start := time.Now()
	go time.AfterFunc(10*time.Second, func() {
		fmt.Println("Goodbye world")
		os.Exit(0)
	})
	Reader := bufio.NewReader(os.Stdin) // read user input
	Userinput, _ := Reader.ReadString('\n')
	if Userinput != "" {
		fmt.Printf("Stopped by user after %0.1f seconds", time.Since(start).Seconds())
	}

}
