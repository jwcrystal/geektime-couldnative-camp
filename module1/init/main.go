package main

import (
	_ "examples/module1/init/a"
	_ "examples/module1/init/b" // already import, same package import once
	"fmt"
)

func init() {
	fmt.Println("main init")
}
func main() {
	fmt.Println("main function")
}
