package main

import (
	"errors"
	"fmt"
)

func main() {
	//	Go語言無內置exception機制，只提供error接口定義錯誤
	//	type error interface {
	// 		Error() string
	//	}
	// 可通過 errors.New (需import errors package) 或 fmt.Errorf 創建新的error
	err := fmt.Errorf("this is a error")
	err2 := errors.New("error as well")
	fmt.Println(err)
	fmt.Println(err2)
}
