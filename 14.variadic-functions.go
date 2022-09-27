package main

import "fmt"

func sum(total int, nums ...int) {
	fmt.Print(nums, " ")
	for _, num := range nums {
		total += num
	}
	fmt.Println(total)
}

func main() {
	sum(0, 1, 2)
	sum(0, 1, 2, 3)

	nums := []int{1, 2, 3, 4}
	sum(0, nums...)
}
