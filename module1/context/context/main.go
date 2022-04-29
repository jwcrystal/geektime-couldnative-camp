package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	baseCtx := context.Background()
	ctx := context.WithValue(baseCtx, "a", "b")
	go func(c context.Context) {
		fmt.Println(c.Value("a"))
	}(ctx)
	timeoutCtx, cancel := context.WithTimeout(baseCtx, time.Second)
	defer cancel()
	go func(c context.Context) {
		ticker := time.NewTicker(1 * time.Second)
		for _ = range ticker.C {
			select {
			case <-c.Done():
				fmt.Println("child process interrupt...")
				break
			default:
				fmt.Println("enter default")
			}
		}
	}(timeoutCtx)
	time.Sleep(1 * time.Second)
	select {
	case <-timeoutCtx.Done():
		time.Sleep(time.Second)
		fmt.Println("main process exit")
	}
}
