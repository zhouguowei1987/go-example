package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	//注册这个给定的通道用于接受特定的信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("goroutine1...")
	}()

	go func() {
		fmt.Println("goroutine2...")
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()

	go func() {
		fmt.Println("goroutine3...")
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
