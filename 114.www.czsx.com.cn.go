package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
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
	CzSxEnableHttpProxy = false
	CzSxHttpProxyUrl    = "111.225.152.186:8089"
)

func CzSxSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(CzSxHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

// ychEduSpider 获取初中数学免费资源
// @Title 获取初中数学免费资源
// @Description http://www.czsx.com.cn/，获取初中数学免费资源
func main() {
	maxPage := 287
	page := 1
	isPageListGo := true
	for isPageListGo {
		requestUrl := fmt.Sprintf("http://www.czsx.com.cn/sort.asp?AClassID=5&NClassID=0&GClassID=0&Page=%d", page)
		fmt.Println(requestUrl)

		pageDoc, err := htmlquery.LoadURL(requestUrl)
		if err != nil {
			fmt.Println(err)
		}
		tableNodes := htmlquery.Find(pageDoc, `//html/body/center/table[5]/tbody/tr/td[3]/table`)
		if len(tableNodes) <= 3 {
			isPageListGo = false
			break
		}

		for i, tableNode := range tableNodes {
			if i%2 != 0 || i == 0 {
				continue
			}
			aHrefNode := htmlquery.FindOne(tableNode, `./tbody/tr/td[1]/a/@href`)
			aHrefUrl := htmlquery.InnerText(aHrefNode)
			fileIdStr := strings.ReplaceAll(aHrefUrl, "download.asp?id=", "")
			fileId, _ := strconv.Atoi(fileIdStr)

			title := htmlquery.InnerText(htmlquery.FindOne(tableNode, `./tbody/tr/td[1]/a/font`))
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, "/", "-")

			filePath := "../www.rar_czsx.com.cn/" + title + ".rar"
			if _, err := os.Stat(filePath); err != nil {
				detailUrl := fmt.Sprintf("http://www.czsx.com.cn/down.asp?id=%d", fileId)
				fmt.Println(detailUrl)

				downloadUrl := fmt.Sprintf("http://www.czsx.com.cn/downloadcheck.asp?id=%d&Ad=wt0", fileId)
				fmt.Println(downloadUrl)

				fmt.Println("=======开始下载" + title + "========")
				err := downloadCzSx(downloadUrl, detailUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")
			}

		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}
func downloadCzSx(attachmentUrl string, referer string, filePath string) error {
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
	if CzSxEnableHttpProxy {
		client = CzSxSetHttpProxy()
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
	req.Header.Set("Host", "http://www.czsx.com.cn")
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
