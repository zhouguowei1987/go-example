package main

import "fmt"

type rect19 struct {
	width, height int
}

func (r *rect19) area() int {
	return r.width * r.height
}

func (r rect19) perim() int {
	return 2*r.width + 2*r.height
}

func main() {
	r := rect19{width: 10, height: 5}

	fmt.Println("area:", r.area())
	fmt.Println("perim:", r.perim())

	rp := &r
	fmt.Println("area:", rp.area())
	fmt.Println("perim:", rp.perim())
}
