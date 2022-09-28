package main

import (
	"fmt"
	"time"
)

func worker25(done chan bool) {
	fmt.Println("working...")
	time.Sleep(time.Second)
	fmt.Println("done")
	done <- true
}
func main() {
	done := make(chan bool, 1)
	go worker25(done)
	<-done
}
