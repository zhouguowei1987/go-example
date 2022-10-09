package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func check58(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	dat, err := ioutil.ReadFile("dat")
	check58(err)
	fmt.Print(string(dat))

	f, err := os.Open("dat")
	check58(err)
	defer f.Close()

	b1 := make([]byte, 5)
	n1, err := f.Read(b1)
	check58(err)
	fmt.Printf("%d bytes : %s\n", n1, string(b1))

	o2, err := f.Seek(6, 0)
	check58(err)
	b2 := make([]byte, 2)
	n2, err := f.Read(b2)
	check58(err)
	fmt.Printf("%d bytes @ %d:%s\n", n2, o2, string(b2))

	o3, err := f.Seek(6, 0)
	check58(err)
	b3 := make([]byte, 2)
	n3, err := io.ReadAtLeast(f, b3, 2)
	check58(err)
	fmt.Printf("%d bytes @ %d:%s\n", n3, o3, string(b3))

	_, err = f.Seek(0, 0)
	check58(err)

	r4 := bufio.NewReader(f)
	b4, err := r4.Peek(5)
	check58(err)
	fmt.Printf("5 bytes: %s\n", string(b4))
}
