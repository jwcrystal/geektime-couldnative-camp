package main

import "fmt"

func main() {
	name := "testing"
	fmt.Printf("%s\n", name)
	// go vet 提示
	// ./main.go:8:2: fmt.Printf format %s reads arg #2, but call has 1 arg
	fmt.Printf("%s%s\n", name)

}
