package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ch := make(chan int) // make(chan datatype, buffer), default為 0，則代表如果沒有讀取數據，則寫入會阻塞
	go func() {
		fmt.Println("Hello from channel")
		ch <- 0
	}()
	i := <-ch
	fmt.Println("Get from channel:", i)

	GoThroughChannelBuffer()
}

func GoThroughChannelBuffer() {
	ch := make(chan int, 10)
	go func() {
		for i := 0; i < 10; i++ {
			rand.Seed(time.Now().UnixNano())
			n := rand.Intn(10) // n will be between 0 and 10
			fmt.Println("putting:", n)
			ch <- n
		}
		close(ch)
	}()
	fmt.Println("hello from GoThroughChannelBuffer()")
	for v := range ch {
		fmt.Println("receiving:", v)
	}
}
