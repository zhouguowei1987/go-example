package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	fmt.Println(now)

	secs := now.Unix() //秒
	fmt.Println(secs)

	nanos := now.UnixNano() //纳秒
	fmt.Println(nanos)

	//UnixMillis是不存在的，所有要得到毫秒的话，要手动从纳秒转化一下
	millis := nanos / 1000000
	fmt.Println(millis)

	fmt.Println(time.Unix(secs, 0))
	fmt.Println(time.Unix(0, nanos))
}
