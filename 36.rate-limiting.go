package main

import (
	"fmt"
	"time"
)

func main() {
	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	//for j := 1; j <= 5; j++ {
	//	fmt.Println(<-requests)
	//}

	//for req := range requests {
	//	fmt.Println(req)
	//}

	tickerLimiter := time.NewTicker(time.Millisecond * 200)
	//tickLimiter := time.Tick(time.Millisecond * 200)
	for req := range requests {
		<-tickerLimiter.C
		//<-tickLimiter
		fmt.Println("request", req, time.Now())
	}

	fmt.Println("============================")

	burstyLimiter := make(chan time.Time, 3)
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	go func() {
		for t := range time.Tick(time.Millisecond * 200) {
			burstyLimiter <- t
		}
	}()

	burstyRequests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		burstyRequests <- i
	}
	close(burstyRequests)
	for req := range burstyRequests {
		<-burstyLimiter
		fmt.Println("request", req, time.Now())
	}
}
