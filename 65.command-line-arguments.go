package main

import (
	"fmt"
	"os"
)

func main() {
	argsWithProg := os.Args
	argsWithOutProg := os.Args[1:]

	arg := os.Args[3]

	fmt.Println(argsWithProg)
	fmt.Println(argsWithOutProg)
	fmt.Println(arg)
}
