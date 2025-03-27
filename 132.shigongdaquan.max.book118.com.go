package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"math/rand"
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
	ShiGongDaQuanEnableHttpProxy = false
	ShiGongDaQuanHttpProxyUrl    = "27.42.168.46:55481"
)

func ShiGongDaQuanSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(ShiGongDaQuanHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

// eduYSpider 获取施工大全文档
// @Title 获取施工大全文档
// @Description http://shigongdaquan.max.book118.com/，获取施工大全文档
func main() {
	// 第一步：顶级分类
	var topListUrl = "http://shigongdaquan.max.book118.com/"
	topListDoc, err := htmlquery.LoadURL(topListUrl)
	if err != nil {
		fmt.Println(err)
	}
	// /html/body/table/tbody/tr
	topCategoryListNodes := htmlquery.Find(topListDoc, `//html/body/table/tbody/tr`)
	if len(topCategoryListNodes) >= 1 {
		for _, topCategoryNode := range topCategoryListNodes {
			topDirNodeA := htmlquery.FindOne(topCategoryNode, `./td/a`)
			if topDirNodeA == nil {
				fmt.Println("不是顶级类别，跳过")
				continue
			}
			topDirTitle := htmlquery.InnerText(topDirNodeA)
			fmt.Println(topDirTitle)

			topDirUrlA := htmlquery.FindOne(topDirNodeA, `./@href`)
			if topDirUrlA == nil {
				fmt.Println("连接不存在")
				continue
			}

			secondCategoryUrl := "http://shigongdaquan.max.book118.com" + htmlquery.InnerText(topDirUrlA)
			fmt.Println(secondCategoryUrl)
			secondListDoc, err := htmlquery.LoadURL(secondCategoryUrl)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// /html/body/table[2]/tbody/tr[1]
			listNodes := htmlquery.Find(secondListDoc, `//html/body/table[2]/tbody/tr`)
			if len(listNodes) >= 1 {
				for _, listNode := range listNodes {
					wordA := htmlquery.FindOne(listNode, `./td/a`)
					if wordA == nil {
						fmt.Println("不是要提取的内容，跳过")
						continue
					}
					title := htmlquery.InnerText(wordA)
					title = strings.ToLower(title)
					if !strings.Contains(title, "doc") && !strings.Contains(title, "docx") && !strings.Contains(title, "pdf") {
						fmt.Println("不是要提取的内容类型，跳过")
						continue
					}
					fmt.Println(title)

					wordUrlA := htmlquery.FindOne(listNode, `./td/a/@href`)
					if wordUrlA == nil {
						fmt.Println("连接不存在")
						continue
					}

					attachmentUrl := "http://shigongdaquan.max.book118.com" + htmlquery.InnerText(wordUrlA)
					fmt.Println(attachmentUrl)
					filePath := "F:\\workspace\\shigongdaquan.max.book118.com\\shigongdaquan.max.book118.com\\" + title
					_, err = os.Stat(filePath)
					if err != nil {
						fmt.Println("=======开始下载========")
						err := downloadShiGongDaQuan(attachmentUrl, filePath, secondCategoryUrl)
						if err != nil {
							fmt.Println(err)
							continue
						}
						fmt.Println("=======完成下载========")
						DownLoadShiGongDaQuanTimeSleep := rand.Intn(10)
						for i := 1; i <= DownLoadShiGongDaQuanTimeSleep; i++ {
							time.Sleep(time.Second)
							fmt.Println("topDirTitle="+topDirTitle+"===========下载", title, "成功，暂停", DownLoadShiGongDaQuanTimeSleep, "秒，倒计时", i, "秒===========")
						}
					}
				}
			}
		}
	}
}

func downloadShiGongDaQuan(attachmentUrl string, filePath string, referer string) error {
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
	if ShiGongDaQuanEnableHttpProxy {
		client = ShiGongDaQuanSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "shigongdaquan.max.book118.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", referer)
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
