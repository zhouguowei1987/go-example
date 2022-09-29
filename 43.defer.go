package main

import (
	"fmt"
	"os"
)

func createFile(s string) *os.File {
	fmt.Println("creating...")
	f, err := os.Create(s)
	if err != nil {
		panic(err)
	}
	return f
}

func writeFile(f *os.File) {
	fmt.Println("writing...")
	fmt.Fprintf(f, "data....")
}

func closeFile(f *os.File) {
	fmt.Println("closing...")
	f.Close()
}

func main() {
	f := createFile("defer.txt")
	defer closeFile(f)
	writeFile(f)
}
