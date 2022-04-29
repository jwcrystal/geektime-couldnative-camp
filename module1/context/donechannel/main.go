package main

import (
	"fmt"
	"time"
)

func main() {
	message := make(chan int, 10)
	done := make(chan bool)

	defer close(message)
	// Consumer
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for _ = range ticker.C {
			select {
			case <-done:
				fmt.Println("child process interrupt...")
			default:
				//fmt.Println("receive message:", <-message)
				fmt.Printf("receive meassage: %d\n", <-message)
			}
		}
	}()
	// Producer
	for i := 0; i < 10; i++ {
		message <- i
	}
	time.Sleep(5 * time.Second)
	close(done)
	time.Sleep(1 * time.Second)
	fmt.Println("main process exit")
}
