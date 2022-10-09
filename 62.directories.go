package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func check62(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	err := os.Mkdir("subdir", 0755)
	check62(err)

	defer os.RemoveAll("subdir")

	createEmptyFile := func(name string) {
		d := []byte("")
		check62(ioutil.WriteFile(name, d, 0644))
	}

	createEmptyFile("subdir/file1")

	err = os.MkdirAll("subdir/parent/child", 0755)
	check62(err)

	createEmptyFile("subdir/parent/file2")
	createEmptyFile("subdir/parent/file3")
	createEmptyFile("subdir/parent/child/file4")

	c, err := ioutil.ReadDir("subdir/parent")
	check62(err)
	fmt.Println("Listing subdir/parent")
	for _, entry := range c {
		fmt.Println(entry.Name(), entry.IsDir())
	}

	err = os.Chdir("subdir/parent/child")
	check62(err)
	c, err = ioutil.ReadDir(".")
	check62(err)
	fmt.Println("Listing subdir/parent/child")
	for _, entry := range c {
		fmt.Println(entry.Name(), entry.IsDir())
	}

	err = os.Chdir("../../..")
	check62(err)
	fmt.Println("Visiting subdir")
	visit := func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(p, info.IsDir())
		return nil
	}
	err = filepath.Walk("subdir", visit)
}
