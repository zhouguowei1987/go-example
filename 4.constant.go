package main

import (
	"fmt"
	"math"
)

const s4 string = "constant"

func main() {
	fmt.Println(s4)

	const n = 500000000
	const d = 3e20 / n
	fmt.Println(d)

	fmt.Println(int64(d))

	fmt.Println(float64(n))
	fmt.Println(math.Sin(n))
}
