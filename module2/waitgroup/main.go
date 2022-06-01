package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	waitByWaitGroup()
}

func waitBySleep() {
	// 因為無法確定goroutine執行的時間，故不適合使用time.sleep
	for i := 0; i < 100; i++ {
		go fmt.Println(i)
	}
	time.Sleep(time.Second)
}

func waitByChannel() {
	c := make(chan bool, 100) // channel是額外開銷
	for i := 0; i < 100; i++ {
		go func(i int) {
			fmt.Println(i)
			c <- true
		}(i)
	}
	for i := 0; i < 100; i++ {
		<-c
	}
}

func waitByWaitGroup() {
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			fmt.Println(i)
			wg.Done()
		}(i)
	}
	wg.Wait() // 等待加入的goroutine結束
}
