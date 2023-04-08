package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/djimenez/iconv-go"
	"golang.org/x/net/html"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	ZuoWenWangEnableHttpProxy = false
	ZuoWenWangHttpProxyUrl    = "111.225.152.186:8089"
)

func ZuoWenWangSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(ZuoWenWangHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

// ychEduSpider 获取满分作文网文档
// @Title 获取满分作文网文档
// @Description https://www.zuowenwang.net/，获取满分作文网文档
func main() {
	var startId = 1
	var endId = 279370
	var id = startId
	var isGoGo = true
	for isGoGo {
		fmt.Println("========================")
		attachmentUrl := fmt.Sprintf("https://pay.biyan8.com/index/scan/down?url=https://www.zuowenwang.net/p/%d.html", id)
		detailUrl := fmt.Sprintf("https://www.zuowenwang.net/p/%d.html", id)
		fmt.Println(detailUrl)

		detailDoc, _ := getZuoWenWangDetail(detailUrl)
		//fileName := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="article-t"]/h1`))
		fileName := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="relates"]/ul/p[1]/strong`))
		fileName = strings.ReplaceAll(fileName, "/", "-")
		fileName = strings.ReplaceAll(fileName, " ", "")
		fileName, _ = iconv.ConvertString(fileName, "gb2312", "utf-8")
		fmt.Println(fileName)

		filePath := "../www.zuowenwang.net/ " + strconv.Itoa(id%28) + "/"
		fileName = strconv.Itoa(id) + "-" + fileName + ".docx"
		if _, err := os.Stat(filePath + fileName); err != nil {
			err := downloadZuoWenWang(attachmentUrl, detailUrl, filePath, fileName)
			//time.Sleep(time.Second * 1)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		id++
		if id >= endId {
			isGoGo = false
		}
	}
}

func getZuoWenWangDetail(url string) (doc *html.Node, err error) {
	client := &http.Client{}                     //初始化客户端
	req, err := http.NewRequest("GET", url, nil) //建立连接
	if err != nil {
		return doc, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.zuowenwang.net")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
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

func downloadZuoWenWang(attachmentUrl string, referer string, filePath string, fileName string) error {
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
			ResponseHeaderTimeout: time.Second * 3,
		},
	}
	if ZuoWenWangEnableHttpProxy {
		client = ZuoWenWangSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "https://www.zuowenwang.net")
	req.Header.Set("Referer", referer)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
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
	out, err := os.Create(filePath + fileName)
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
