package main

import (
	"fmt"
	"log"

	"github.com/otiai10/gosseract/v2"
)

func main() {
	// 创建Tesseract客户端
	client := gosseract.NewClient()
	defer client.Close()

	// 设置Tesseract的路径（如果你的环境变量已设置，可以省略）
	client.SetTessdataPath("C:\\Program Files\\Tesseract-OCR\\tessdata")

	// 设置OCR语言为英文，并设置字符白名单为字母和数字
	client.SetLanguage("eng")
	client.SetVariable("tessedit_char_whitelist", "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	// 设置要识别的验证码图像路径
	imagePath := "./dbba-validate-code/validate-code.png" // 这里替换为你的验证码图片路径
	// 识别图片中的文本
	err := client.SetImage(imagePath)
	if err != nil {
		log.Fatalf("设置图片出错: %v", err)
	}
	text, err := client.Text()
	if err != nil {
		log.Fatalf("识别出错: %v", err)
	}
	fmt.Printf("识别的验证码是: %s\n", text)
}
