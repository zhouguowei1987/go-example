package main

import (
	"fmt"
	"log"

	"github.com/otiai10/gosseract/v2"
)

func main() {
	// 创建OCR客户端
	client := gosseract.NewClient()
	defer client.Close()

	// 设置验证码图片路径
	client.SetImage("182.jpg")

	// 执行OCR识别
	text, err := client.Text()
	if err != nil {
		log.Fatalf("识别失败: %v", err)
	}

	// 输出识别结果
	fmt.Printf("识别的验证码是: %s\n", text)
}
