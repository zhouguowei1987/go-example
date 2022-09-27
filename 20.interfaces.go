package main

import (
	"fmt"
	"math"
)

type geometry interface {
	area() float64
	perim() float64
}

type rect20 struct {
	width, height float64
}

func (r *rect20) area() float64 {
	return r.width * r.height
}

func (r *rect20) perim() float64 {
	return 2*r.width + 2*r.height
}

type circle20 struct {
	radius float64
}

func (c *circle20) area() float64 {
	return math.Pi * c.radius * c.radius
}
func (c *circle20) perim() float64 {
	return 2 * math.Pi * c.radius
}

func measure(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.perim())
}

func main() {
	r := &rect20{width: 3, height: 4}
	measure(r)

	c := &circle20{radius: 5}
	measure(c)
}
