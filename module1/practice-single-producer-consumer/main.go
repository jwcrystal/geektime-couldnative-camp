package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	ch := make(chan string)
	go Producer(ch)
	Consumer(ch)
}

func Producer(c chan string) {
	for i := 0; i < 10; i++ {
		c <- strconv.Itoa(i)
		time.Sleep(time.Second)
		fmt.Printf("Producer[%d]\n", i)
	}
	close(c)
}

func Consumer(c chan string) {
	//fmt.Printf("Concumer[%s]\n", <-c)
	for v := range c {
		fmt.Printf("Concumer[%s]\n", v)
	}
}
