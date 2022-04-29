package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, World ")
	var a int = 2
	b := 1
	fmt.Println(a, b)
	GoroutineExample()
}

func GoroutineExample() {
	// 協程（線程）順序不定，依CPU執行
	go fmt.Println("1")
	go fmt.Println("2")
	go fmt.Println("3")
	time.Sleep(time.Second)
}
