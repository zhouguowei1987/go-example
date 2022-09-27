package main

import "fmt"

func main() {
	ch := make(chan string, 2)

	ch <- "buffered"
	ch <- "channel"

	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
