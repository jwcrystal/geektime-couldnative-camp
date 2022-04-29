package main

import "fmt"

func main() {
	interfaces := []IF{}
	h := new(Human)
	h.firstName = "first"
	h.lastName = "last"
	interfaces = append(interfaces, h)
	c := new(Car)
	//c := Car{} //Car does not implement IF (getName method has pointer receiver)
	c.factory = "benz"
	c.model = "s"
	interfaces = append(interfaces, c) // struct 可以加到同一個array （slice）
	// GO lang 有能力自動判斷對應接口的實現， interface可能為"nil"，一定要預判空，否則會crash（nil panic）
	for _, f := range interfaces {
		println(f.getName())
	}
	p := Plane{}
	p.vendor = "testVendor"
	p.model = "testModel"
	fmt.Println(p.getName())
}

type IF interface {
	getName() string
}

type Human struct { // Struct初始化意味著空間分配， 所以struct的引用不會出現空指針
	firstName, lastName string
}

type Plane struct {
	vendor, model string
}

type Car struct {
	factory string
	model   string
}

// 用多個結構體（struct）去實現接口（interface）
func (h *Human) getName() string {
	return h.firstName + "," + h.lastName
}

func (p *Plane) getName() string {
	return fmt.Sprintf("vendor: %s, model: %s", p.vendor, p.model)
}

func (c *Car) getName() string {
	return c.factory + "-" + c.model
}
