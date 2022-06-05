package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var signal = false

func main() {
	multiChan := make(chan string, 10)

	wgProducer := new(sync.WaitGroup)
	wgConsumer := new(sync.WaitGroup)

	for i := 0; i < 5; i++ {
		wgProducer.Add(1)
		go Producer(i, wgProducer, multiChan)
	}
	for i := 0; i < 3; i++ {
		wgConsumer.Add(1)
		go Consumer(wgConsumer, multiChan)
	}
	// timeout
	go func() {
		time.Sleep(time.Second * 2) // producer stop due to timeout after 2 sec
		signal = true
	}()

	wgProducer.Wait()
	close(multiChan)
	wgConsumer.Wait()
}

func Producer(id int, wg *sync.WaitGroup, c chan string) {
	count := 0
	for !signal {
		time.Sleep(time.Second)
		count++
		data := strconv.Itoa(id) + "-" + strconv.Itoa(count)
		fmt.Printf("producer[%s]\n", data)
		c <- data
	}
	wg.Done()
}
func Consumer(wg *sync.WaitGroup, c chan string) {
	for data := range c {
		time.Sleep(time.Second)
		fmt.Printf("comsumer[%s]\n", data)
	}
	wg.Done()
}
