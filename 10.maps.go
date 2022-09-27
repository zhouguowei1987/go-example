package main

import "fmt"

func main() {
	var i map[string]int
	fmt.Println("map:", i)

	m := make(map[string]int)
	m["k1"] = 7
	m["k2"] = 8
	fmt.Println("map:", m)

	v1 := m["k1"]
	fmt.Println("v1:", v1)

	fmt.Println("len:", len(m))

	delete(m, "k1")
	fmt.Println("del:", m)

	_, prs := m["k2"]
	fmt.Println("prs:", prs)

	n := map[string]int{"foo": 1, "bar": 2}
	fmt.Println("map:", n)
}
