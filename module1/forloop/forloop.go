package main

import (
	"fmt"
)

func main() {
	//課後練習1.1
	//给定一个字符串数组
	//[“I”,“am”,“stupid”,“and”,“weak”]
	//用 for 循环遍历该数组并修改为
	//[“I”,“am”,“smart”,“and”,“strong”]
	str := []string{"I", "am", "stupid", "and", "weak"}
	// normal for loop
	fmt.Println("normal for loop in range")
	fmt.Println(GerneralForLoop(str))
	// index for loop in range
	fmt.Println("index for loop in range")
	fmt.Println(IndexForRangeLoop(str))
	// value for loop in range
	fmt.Println("value for loop in range")
	fmt.Println(ValueForRangeLoop(str))
	// map method
	fmt.Println("by map method")
	fmt.Println(ByMap(str))
}
func GerneralForLoop(str []string) []string {
	for i := 0; i < len(str); i++ {
		//if i == 2 {
		//	str[i] = "smart"
		//} else if i == 4 {
		//	str[i] = "strong"
		//}
		// or use switch
		switch i {
		case 2:
			str[i] = "smart"
		case 4:
			str[i] = "strong"
		}
	}
	return str
}
func IndexForRangeLoop(str []string) []string {
	for index, value := range str {
		switch index {
		case 2:
			str[index] = "smart"
		case 4:
			str[index] = "strong"
		}
		fmt.Println(index, value)
	}
	return str
}
func ValueForRangeLoop(str []string) []string {
	for index, value := range str {
		switch value {
		case "weak":
			str[index] = "strong"
		case "stupid":
			str[index] = "smart"
		}
		fmt.Println(index, value)
	}
	return str
}
func ByMap(str []string) []string {
	myMap := make(map[string]string, 2)
	myMap["stupid"] = "smart"
	myMap["weak"] = "strong"
	fmt.Println(myMap)
	for i := 0; i < len(str); i++ {
		for k, v := range myMap {
			if str[i] == k {
				str[i] = v
			}
		}
	}
	return str
}
