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
	CmaEnableHttpProxy = false
	CmaHttpProxyUrl    = "111.225.152.186:8089"
)

func CmaSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(CmaHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

// 获取中国气象局标准
// @Title 获取中国气象局标准
// @Description https://www.cma.gov.cn/ 获取中国气象局标准
func main() {
	maxPage := 79
	page := 0
	isPageListGo := true
	for isPageListGo {
		requestUrl := "https://www.cma.gov.cn/zfxxgk/gknr/flfgbz/bz/index.html"
		if page > 0 {
			requestUrl = fmt.Sprintf("https://www.cma.gov.cn/zfxxgk/gknr/flfgbz/bz/index_%d.html", page)
		}
		fmt.Println(requestUrl)

		pageDoc, err := htmlquery.LoadURL(requestUrl)
		if err != nil {
			fmt.Println(err)
		}
		liNodes := htmlquery.Find(pageDoc, `//div[@class="boxcenter"]/div[@class="mainbox clearfix"]/div[@class="mainCont clearfix"]/div[@class="rightBox2"]/div[@id="demo"]/ul[@class="mesgopen2 list"]/li[@class="list-item"]`)
		if len(liNodes) <= 0 {
			isPageListGo = false
			break
		}

		for _, liNode := range liNodes {

			aHrefNode := htmlquery.FindOne(liNode, `./a/@href`)
			if aHrefNode == nil {
				continue
			}
			detailUrl := htmlquery.InnerText(aHrefNode)
			detailUrl = "https://www.cma.gov.cn/zfxxgk/gknr/flfgbz/bz/" + strings.ReplaceAll(detailUrl, "./", "")
			fmt.Println(detailUrl)

			detailDoc, err := htmlquery.LoadURL(detailUrl)
			if err != nil {
				fmt.Println(err)
				continue
			}

			titleNode := htmlquery.FindOne(detailDoc, `//div[@class="boxcenter"]/div[@class="mainbox clearfix"]/div[@class="mainCont clearfix"]/div[@class="rightBox rightbox5"]/h1[@class="title"]`)
			if titleNode == nil {
				fmt.Println("未找到标题节点，跳过")
				continue
			}
			title := strings.TrimSpace(htmlquery.InnerText(titleNode))
			title = strings.ReplaceAll(title, "/", "-")
			fmt.Println(title)

			// /html/body/div/div[2]/div/div/div[1]/table/tbody/tr[2]/td[3]/span
			codeNode := htmlquery.FindOne(detailDoc, `//div[@class="boxcenter"]/div[@class="mainbox clearfix"]/div[@class="mainCont clearfix"]/div[@class="rightBox rightbox5"]/div[@class="fu"]/table/tbody/tr[2]/td[3]/span`)
			if codeNode == nil {
				fmt.Println("未找到标准号节点，跳过")
				continue
			}
			code := strings.TrimSpace(htmlquery.InnerText(codeNode))
			code = strings.ReplaceAll(code, "/", "-")
			fmt.Println(code)

			filePath := "../www.cma.gov.cn/" + title + "(" + code + ")" + ".pdf"
			fmt.Println(filePath)
			if _, err := os.Stat(filePath); err != nil {
				downloadNode := htmlquery.FindOne(detailDoc, `//div[@class="boxcenter"]/div[@class="mainbox clearfix"]/div[@class="mainCont clearfix"]/div[@class="rightBox rightbox5"]/div[@class="relList"]/ul[@class="fujian"]/li[@class="1"]/a/@href`)
				if downloadNode == nil {
					fmt.Println("未找到下载文件节点，跳过")
					continue
				}

				detailUrlArray := strings.Split(detailUrl, "/")
				downloadUrlArray := detailUrlArray[:len(detailUrlArray)-1]
				downloadUrl := strings.Join(downloadUrlArray, "/") + strings.ReplaceAll(htmlquery.InnerText(downloadNode), "./", "/")
				fmt.Println(downloadUrl)

				fmt.Println("=======开始下载" + title + "========")
				err = downloadCma(downloadUrl, detailUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "../www.cma.gov.cn", "../upload.doc88.com/www.cma.gov.cn")
				err = copyCmaFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")
				//DownLoadCmaTimeSleep := 10
				DownLoadCmaTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadCmaTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadCmaTimeSleep, "秒 倒计时", i, "秒===========")
				}

			}
		}
		DownLoadCmaPageTimeSleep := 10
		// DownLoadCmaPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadCmaPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadCmaPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

func downloadCma(attachmentUrl string, referer string, filePath string) error {
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
	if CmaEnableHttpProxy {
		client = CmaSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.cma.gov.cn")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 如果访问失败 就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}

	// 创建一个文件用于保存
	fileDiv := filepath.Dir(filePath)
	if _, err = os.Stat(fileDiv); err != nil {
		if os.MkdirAll(fileDiv, 0o777) != nil {
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

func copyCmaFile(src, dst string) (err error) {
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
