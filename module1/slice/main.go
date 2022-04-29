package main

import "fmt"

func main() {
	SliceExample()
}

func SliceExample() {
	myArry := [5]int{1, 2, 3, 4, 5}
	mySlice := myArry[1:3]
	fmt.Printf("mySlice %+v\n", mySlice)
	fullSlice := myArry[:]
	remove3rdItem := deleteItem(fullSlice, 2)
	fmt.Printf("remove3rdItem %+v\n", remove3rdItem)
}
func deleteItem(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...) // no remove specific item, used by slice combination ("..." can't ignore)
}
