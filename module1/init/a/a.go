package a

import (
	_ "examples/module1/init/b"
	"fmt"
)

func init() {
	fmt.Println("init from a")
}
