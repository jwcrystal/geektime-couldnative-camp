package main

func main() {
	DoOperation(1, Increase)
	DoOperation(1, Decrease)
	func() { println("Anonymous") }()
}

func DoOperation(y int, f func(int, int)) {
	f(y, 1)
}

func Increase(a, b int) {
	println("increase result is:", a+b)
}

func Decrease(a, b int) {
	println("decrease result is:", a-b)
}
