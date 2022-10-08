package main

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
)

func main() {
	s := "sha1 this string"

	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	fmt.Println(s)
	fmt.Printf("% x\n", bs)

	ms := "md5 this string"

	mh := md5.New()
	mh.Write([]byte(ms))
	mbs := mh.Sum(nil)
	fmt.Println(ms)
	fmt.Printf("% x\n", mbs)

}
