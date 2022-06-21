package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//timer start to measure 10 seconds
//after 10 seconds Goodbye message appear and program exit
func main() {
	fmt.Println("Hello World")
	start := time.Now()
	// two channels for monitoring user's input
	sigs := make(chan os.Signal, 1)
	userStop := make(chan string)

	//count 10 second before closing
	go func() {
		time.AfterFunc(10*time.Second, func() {
			fmt.Println("Goodbye world")
			os.Exit(0)
		})
	}()
	// monitoring Ctrl+C pressing
	go func() {
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		// reading of user input with Enter key
	}()
	go func() {
		Reader := bufio.NewReader(os.Stdin) // read user input
		Userinput, _ := Reader.ReadString('\n')
		userStop <- Userinput
	}()
	//choose first input from user
	for i := 0; i < 2; i++ {
		select {
		case <-sigs:
			fmt.Printf("Stopped by user after %0.1f seconds", time.Since(start).Seconds())
			os.Exit(0)
		case <-userStop:
			fmt.Printf("Stopped by user after %0.1f seconds", time.Since(start).Seconds())
			os.Exit(0)
		}
	}
}
