package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func check59(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	d1 := []byte("hello\ngo\n")
	err := ioutil.WriteFile("dat1", d1, 0644)
	check59(err)

	f, err := os.Create("dat2")
	check59(err)
	defer f.Close()

	d2 := []byte{115, 111, 109, 101, 10}
	n2, err := f.Write(d2)
	check59(err)
	fmt.Printf("wrote %d bytes\n", n2)

	n3, err := f.WriteString("writes\n")
	check59(err)
	fmt.Printf("wrote %d bytes\n", n3)

	f.Sync()

	w := bufio.NewWriter(f)
	n4, err := w.WriteString("buffered\n")
	check59(err)
	fmt.Printf("wrote %d bytes\n", n4)
	w.Flush()
}
