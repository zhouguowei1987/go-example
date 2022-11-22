package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	//appkey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	//secret := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	//
	//Timestamp := time.Now().UnixNano() / 1e6
	//fmt.Println("Timestamp 的数据类型是:", reflect.TypeOf(Timestamp))
	//fmt.Println(Timestamp)
	////十进制
	//Time := strconv.FormatInt(Timestamp, 10)
	//fmt.Println("Time 的数据类型是:", reflect.TypeOf(Time))

	Sign := "20221031T073525Z"

	hash := sha256.New()
	hash.Write([]byte(Sign))
	sum := hash.Sum(nil)

	sign := hex.EncodeToString(sum)
	fmt.Println(sign)

	fmt.Println(len(sign))

}
