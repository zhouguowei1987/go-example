package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"

	// "math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var CrsaCookie = "acw_tc=276077d817786491955958780ed676d5e46c7a6bbadc484ade3233d0e7a9dd"

// CrsaSpider 获取中国道路交通安全协会标准文档
// @Title 获取中国道路交通安全协会标准文档
// @Description https://www.crsa.net/，将中国道路交通安全协会标准文档入库
func main() {
	var startId = 12
	var endId = 49
	for id := startId; id <= endId; id++ {
		fmt.Println(id)
		detailUrl := fmt.Sprintf("https://www.crsa.net/bz/article/news/%d", id)
		detailDoc, err := CrsaDetailDoc(detailUrl)
		if err != nil {
			fmt.Println(err)
			continue
		}

		titleNode := htmlquery.FindOne(detailDoc, `//html/body/div/div/div/div/div/div/div[2]/div/div/div[1]/div[1]`)
		if titleNode == nil {
			fmt.Println("标题不存在，跳过")
			continue
		}
		title := strings.TrimSpace(htmlquery.InnerText(titleNode))
		title = strings.TrimSpace(title)
		title = strings.ReplaceAll(title, "/", "-")
		title = strings.ReplaceAll(title, "／", "-")
		title = strings.ReplaceAll(title, "/", "-")
		title = strings.ReplaceAll(title, "　", "-")
		title = strings.ReplaceAll(title, " ", "-")
		title = strings.ReplaceAll(title, "：", ":")
		title = strings.ReplaceAll(title, "—", "-")
		title = strings.ReplaceAll(title, "－", "-")
		title = strings.ReplaceAll(title, "（", "(")
		title = strings.ReplaceAll(title, "）", ")")
		title = strings.ReplaceAll(title, "《", "")
		title = strings.ReplaceAll(title, "》", "")
		fmt.Println(title)

		codeNode := htmlquery.FindOne(detailDoc, `//html/body/div/div/div/div/div/div/div[2]/div/div/div[3]/div[1]/div[1]/div[2]`)
		if codeNode == nil {
			fmt.Println("标准号不存在，跳过")
			continue
		}
		code := strings.TrimSpace(htmlquery.InnerText(codeNode))
		code = strings.ReplaceAll(code, "/", "-")
		code = strings.ReplaceAll(code, "—", "-")
		fmt.Println(code)

		filePath := "../www.crsa.net/" + title + "(" + code + ")" + ".pdf"
		fmt.Println(filePath)

		_, err = os.Stat(filePath)
		if err == nil {
			fmt.Println("文档已下载过，跳过")
			continue
		}

		downloadNode := htmlquery.FindOne(detailDoc, `//html/body/div/div/div/div/div/div/div[2]/div/div/div[3]/div[4]/div[2]/div[2]/a/@href`)
		if downloadNode == nil {
			fmt.Println("文档下载地址不存在，跳过")
			continue
		}
		downloadUrl := strings.TrimSpace(htmlquery.InnerText(downloadNode))
		fmt.Println(downloadUrl)
		os.Exit(1)

		fmt.Println("=======开始下载========")
		err = downloadCrsa(downloadUrl, detailUrl, filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//复制文件
		tempFilePath := strings.ReplaceAll(filePath, "www.crsa.net", "temp-hbba.sacinfo.org.cn")
		err = copyCrsaFile(filePath, tempFilePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("=======完成下载========")

		// 设置倒计时
		DownLoadTCrsaTimeSleep := 10
		// DownLoadTCrsaTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadTCrsaTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("id="+strconv.Itoa(id)+"===========操作完成，", "暂停", DownLoadTCrsaTimeSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

func CrsaDetailDoc(url string) (doc *html.Node, err error) {
	// 创建一个自定义的http.Transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 忽略证书验证
		},
	}
	client := &http.Client{Transport: tr}        //初始化客户端
	req, err := http.NewRequest("GET", url, nil) //建立连接
	if err != nil {
		return doc, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", CrsaCookie)
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	req.Header.Set("Host", "www.crsa.net")
	req.Header.Set("Referer", "https://www.crsa.net/sat/standard/standardlist/0")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return doc, err
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return doc, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	doc, err = htmlquery.Parse(resp.Body)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func downloadCrsa(requestUrl string, referer string, filePath string) error {
	// 初始化客户端
	var client *http.Client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*3)
				if err != nil {
					fmt.Println("dail timeout", err)
					return nil, err
				}
				return c, nil

			},
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 30,
		},
	} //初始化客户端
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", CrsaCookie)
	req.Header.Set("Host", "www.crsa.net")
	req.Header.Set("Origin", "https://www.crsa.net")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
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
		if os.MkdirAll(fileDiv, 0777) != nil {
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

func copyCrsaFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(in *os.File) {
		err := in.Close()
		if err != nil {
			return
		}
	}(in)

	// 创建一个文件用于保存
	fileDiv := filepath.Dir(dst)
	if _, err = os.Stat(fileDiv); err != nil {
		if os.MkdirAll(fileDiv, 0o777) != nil {
			return err
		}
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)

	_, err = io.Copy(out, in)
	return nil
}
