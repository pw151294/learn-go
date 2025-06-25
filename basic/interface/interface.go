package main

import "fmt"

type Car interface {
	Drive()
}

type TrafficTool interface {
	Drive()
}

type Truck struct {
	Model string
}

func (t Truck) Drive() {
	fmt.Printf("Truck %s is driving\n", t.Model)
}

func main() {
	var c Car = Truck{Model: "Truck"}

	t := c.(Truck)
	fmt.Println(t.Model)

	tt := c.(TrafficTool)
	tt.Drive()

	switch c.(type) {
	case TrafficTool:
		fmt.Println("TrafficTool is driving")
	case Truck:
		fmt.Println("Truck is driving")
	}

	var a interface{}
	fmt.Println(a == nil) // true
	var b *int
	a = b                 // a的底层是eface 赋值操作之后eface内的_type值变为了非空
	fmt.Println(b == nil) // true
	fmt.Println(a == nil) // false
}
