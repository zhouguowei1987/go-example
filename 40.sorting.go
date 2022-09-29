package main

import (
	"fmt"
	"sort"
)

func main() {
	strs := []string{"c", "a", "b"}
	fmt.Println(sort.StringsAreSorted(strs))
	sort.Strings(strs)
	fmt.Println("Strings", strs)
	fmt.Println(sort.StringsAreSorted(strs))

	ints := []int{7, 2, 4}
	sort.Ints(ints)
	fmt.Println("Ints", ints)
	fmt.Println("Sorted:", sort.IntsAreSorted(ints))
}
