package main

import (
	"fmt"
	"reflect"
)

func main() {
	// basic type
	myMap := make(map[string]string, 10)
	myMap["a"] = "b"
	t := reflect.TypeOf(myMap)
	fmt.Println("type: ", t)
	v := reflect.ValueOf(myMap)
	fmt.Println("value: ", v)
	// struct
	myStruct := T{A: "a"}
	v1 := reflect.ValueOf(myStruct)
	for i := 0; i < v1.NumField(); i++ {
		fmt.Printf("Field #{i}: #{v1.Field(i)}\n")
		fmt.Println("Field: ", v1.Field(i))
	}
	for i := 0; i < v1.NumMethod(); i++ {
		fmt.Printf("Method #{i}: #{v1.Method(i)}\n")
		fmt.Println("Method: ", v1.Method(i))
	}
	// 需要注意receive是 struct 還是 pointer
	result := v1.Method(0).Call(nil) // = myStruct.String(), 這邊Method(0)指 呼叫 第一個func
	fmt.Println("result: ", result)
}

type T struct {
	A string
}

// 需要注意receive是 struct 還是 pointer
func (t T) String() string {
	return t.A + "1"
}
