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

var GoodKeJianEnableHttpProxy = false
var GoodKeJianHttpProxyUrl = ""
var GoodKeJianHttpProxyUrlArr = make([]string, 0)

func GoodKeJianHttpProxy() error {
	pageMax := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	//pageMax := []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for _, page := range pageMax {
		freeProxyUrl := "https://www.beesproxy.com/free"
		if page > 1 {
			freeProxyUrl = fmt.Sprintf("https://www.beesproxy.com/free/page/%d", page)
		}
		beesProxyDoc, err := htmlquery.LoadURL(freeProxyUrl)
		if err != nil {
			return err
		}
		trNodes := htmlquery.Find(beesProxyDoc, `//figure[@class="wp-block-table"]/table[@class="table table-bordered bg--secondary"]/tbody/tr`)
		if len(trNodes) > 0 {
			for _, trNode := range trNodes {
				ipNode := htmlquery.FindOne(trNode, "./td[1]")
				if ipNode == nil {
					continue
				}
				ip := htmlquery.InnerText(ipNode)

				portNode := htmlquery.FindOne(trNode, "./td[2]")
				if portNode == nil {
					continue
				}
				port := htmlquery.InnerText(portNode)

				protocolNode := htmlquery.FindOne(trNode, "./td[5]")
				if protocolNode == nil {
					continue
				}
				protocol := htmlquery.InnerText(protocolNode)

				switch protocol {
				case "HTTP":
					GoodKeJianHttpProxyUrlArr = append(GoodKeJianHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					GoodKeJianHttpProxyUrlArr = append(GoodKeJianHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func GoodKeJianSetHttpProxy() (httpclient *http.Client) {
	if GoodKeJianHttpProxyUrl == "" {
		if len(GoodKeJianHttpProxyUrlArr) <= 0 {
			err := GoodKeJianHttpProxy()
			if err != nil {
				GoodKeJianSetHttpProxy()
			}
		}
		GoodKeJianHttpProxyUrl = GoodKeJianHttpProxyUrlArr[0]
		if len(GoodKeJianHttpProxyUrlArr) >= 2 {
			GoodKeJianHttpProxyUrlArr = GoodKeJianHttpProxyUrlArr[1:]
		} else {
			GoodKeJianHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(GoodKeJianHttpProxyUrl)
	ProxyURL, _ := url.Parse(GoodKeJianHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
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
	return httpclient
}

type GoodKeJianPaper struct {
	name string
	url  string
}

var GoodKeJianPapers = []GoodKeJianPaper{
	{
		name: "语文试卷",
		url:  "http://www.goodkejian.com/search.asp?keyword=%d3%ef%ce%c4%2b%ca%d4",
	},
	{
		name: "数学试卷",
		url:  "http://www.goodkejian.com/search.asp?keyword=%ca%fd%d1%a7%2b%ca%d4",
	},
	{
		name: "英语试卷",
		url:  "http://www.goodkejian.com/search.asp?keyword=%D3%A2%D3%EF+%CA%D4",
	},
}

var GoodKeJianNextDownloadSleep = 2

// ychEduSpider 获取好课件吧试卷
// @Title 获取好课件吧试卷
// @Description http://www.goodkejian.com/，获取好课件吧试卷
func main() {
	for _, paper := range GoodKeJianPapers {
		current := 1
		isPageListGo := true
		// 计算最大页数
		paperIndexUrl := paper.url
		paperIndexDoc, err := htmlquery.LoadURL(paperIndexUrl)
		if err != nil {
			fmt.Println(err)
			continue
		}
		paperTotalNode := htmlquery.FindOne(paperIndexDoc, `//div[@id="main_right1"]/div[@class="list_title"]/span[@class="l_t_r"]/font[3]`)
		paperTotalText := htmlquery.InnerText(paperTotalNode)
		paperMaxPages, err := strconv.Atoi(paperTotalText)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for isPageListGo {
			paperListUrl := paper.url + fmt.Sprintf("&page=%d", current)
			paperListDoc, err := htmlquery.LoadURL(paperListUrl)
			if err != nil {
				fmt.Println(err)
				current = 1
				isPageListGo = false
				continue
			}
			liNodes := htmlquery.Find(paperListDoc, `//div[@id="main_right1"]/div[@id="list_c"]/div[@class="list_soft"]`)
			if len(liNodes) <= 0 {
				fmt.Println(err)
				current = 1
				isPageListGo = false
				continue
			}
			for _, liNode := range liNodes {
				fmt.Println("============================================================================")
				fmt.Println("=======当前页URL", paperListUrl, "========")

				title := htmlquery.InnerText(htmlquery.FindOne(liNode, `./span[@class="soft_name"]/strong/a`))
				fmt.Println(title)

				viewHref := "http://www.goodkejian.com/" + htmlquery.InnerText(htmlquery.FindOne(liNode, `./span[@class="soft_name"]/strong/a/@href`))
				fmt.Println(viewHref)

				// 查看是否有下载按钮
				viewDoc, err := htmlquery.LoadURL(viewHref)
				if err != nil {
					fmt.Println(err)
					continue
				}
				downloadButtonNodes := htmlquery.Find(viewDoc, `//div[@id="main_right"]/div[@id="soft"]/div[@id="soft_down"]/div[@id="soft_zhulixiazai"]/ul/li`)
				if len(downloadButtonNodes) <= 0 {
					fmt.Println("没有下载链接")
					continue
				}
				downloadUrl := htmlquery.InnerText(htmlquery.FindOne(downloadButtonNodes[0], `./a/@href`))
				downloadUrl = "http://www.goodkejian.com" + strings.ReplaceAll(downloadUrl, "..", "")

				filePath := "F:\\workspace\\www.goodkejian.com\\www.rar_goodkejian.com\\" + title + ".rar"
				_, err = os.Stat(filePath)
				if err != nil {

					fmt.Println("=======开始下载"+strconv.Itoa(current)+"-", paperMaxPages, "========")
					err = downloadGoodKeJian(downloadUrl, viewHref, filePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("=======下载完成========")
					for i := 1; i <= GoodKeJianNextDownloadSleep; i++ {
						time.Sleep(time.Second)
						fmt.Println("===========操作结束，暂停", GoodKeJianNextDownloadSleep, "秒，倒计时", i, "秒===========")
					}
				}
			}
			current++
			if current > paperMaxPages {
				fmt.Println("没有更多分页")
				break
			}
			isPageListGo = true
		}
	}
}

func downloadGoodKeJian(attachmentUrl string, referer string, filePath string) error {
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
	if GoodKeJianEnableHttpProxy {
		client = GoodKeJianSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "__51cke__=; _gid=GA1.2.944241540.1703144883; PHPSESSID=22d6kpnbct0v3iopp7qctbk081; __tins__21123451=%7B%22sid%22%3A%201703148480266%2C%20%22vd%22%3A%207%2C%20%22expires%22%3A%201703151933380%7D; __51laig__=27; _gat=1; _ga_34B604LFFQ=GS1.1.1703148480.2.1.1703150135.57.0.0; _ga=GA1.1.1587097358.1703144883")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.goodkejian.com")
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
