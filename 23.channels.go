package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	go func() {
		time.Sleep(time.Second * 5)
		ch <- "ping"
	}()

	fmt.Println(<-ch)
}
