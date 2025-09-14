package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
)

// BiaoGeSpider 获取表格网文档
// @Title 获取表格网文档
// @Description https://www.biaoge.com/，将表格网文档入库
func main() {
	var startId = 3430
	var endId = 20949
	for id := startId; id <= endId; id++ {
		detailUrl := fmt.Sprintf("https://www.biaoge.com/cat/%d.html", id)
		fmt.Println(detailUrl)
		detailDoc, err := htmlquery.LoadURL(detailUrl)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 查看文件类型
		fileTypeNode := htmlquery.FindOne(detailDoc, `//html/body/div/div[2]/div[2]/div[2]/ul/li[1]/ul/li[2]/span`)
		if fileTypeNode == nil {
			fmt.Println("没有文件类型节点")
			continue
		}
		fileType := htmlquery.InnerText(fileTypeNode)
		fileType = strings.TrimSpace(fileType)
		fileType = strings.ReplaceAll(fileType, "文件格式：", "")
		fileType = strings.ToLower(fileType)
		fmt.Println(fileType)
		if strings.Index(fileType, "doc") == -1 && strings.Index(fileType, "xls") == -1 {
			fmt.Println("不是doc、xls文档，跳过")
			continue
		}

		// 文档名称
		titleNode := htmlquery.FindOne(detailDoc, `//html/body/div/div[2]/div[2]/div[2]/h1`)
		title := htmlquery.InnerText(titleNode)
		title = strings.TrimSpace(title)
		title = strings.ReplaceAll(title, "（", "(")
		title = strings.ReplaceAll(title, "）", ")")
		title = strings.ReplaceAll(title, "/", "-")
		fmt.Println(title)

		// 文档类目名称
		catNode := htmlquery.FindOne(detailDoc, `//html/body/div/div[2]/div[1]/ul/li[5]/a`)
		cat := htmlquery.InnerText(catNode)
		cat = strings.TrimSpace(cat)
		fmt.Println(cat)

		filePath := "../www.biaoge.com/www.biaoge.com/" + title + "(" + cat + ")" + "." + fileType
		if strings.Index(filePath, "会计学堂软件") != -1 {
			fmt.Println("含有“会计学堂软件”字样，跳过")
			continue
		}
		fmt.Println(filePath)
		_, err = os.Stat(filePath)
		if err == nil {
			fmt.Println("文档已下载过，跳过")
			continue
		}

		// 查看是否有下载链接
		downloadNode := htmlquery.FindOne(detailDoc, `//html/body/div/div[2]/div[2]/div[2]/ul/li[2]/ul/li[1]/span`)
		if downloadNode == nil {
			fmt.Println("没有下载链接")
			continue
		}
		// window.open('https://al3.acc5.com/25中级协议班精讲课讲义xe9.zip')
		clickText := htmlquery.SelectAttr(downloadNode, "onclick")
		// 下载文档URL
		downLoadUrl := strings.ReplaceAll(clickText, "window.open('", "")
		downLoadUrl = strings.ReplaceAll(downLoadUrl, "')", "")
		fmt.Println(downLoadUrl)

		fmt.Println("=======开始下载========")
		err = downloadBiaoGe(downLoadUrl, filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("=======完成下载========")

		//复制文件
		tempFilePath := strings.ReplaceAll(filePath, "../www.biaoge.com/www.biaoge.com", "../www.biaoge.com/temp-www.biaoge.com")
		err = copyBiaoGeFile(filePath, tempFilePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("=======完成下载========")

		// 设置倒计时
		// DownLoadBiaoGeTimeSleep := 10
		DownLoadBiaoGeTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadBiaoGeTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("id="+strconv.Itoa(id)+"===========操作完成，", "暂停", DownLoadBiaoGeTimeSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

func downloadBiaoGe(pdfUrl string, filePath string) error {
	client := &http.Client{}                        //初始化客户端
	req, err := http.NewRequest("GET", pdfUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.biaoge.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}

	// 创建一个文件用于保存
	fileDiv := filepath.Dir(filePath)
	if _, err = os.Stat(fileDiv); err != nil {
		if os.MkdirAll(fileDiv, 0644) != nil {
			return err
		}
	}
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func copyBiaoGeFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer func(in *os.File) {
		err := in.Close()
		if err != nil {
			return
		}
	}(in)

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)

	_, err = io.Copy(out, in)
	return
}
