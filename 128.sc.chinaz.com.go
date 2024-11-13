package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var ChinaZPptEnableHttpProxy = false
var ChinaZPptHttpProxyUrl = "111.225.152.186:8089"
var ChinaZPptHttpProxyUrlArr = make([]string, 0)

func ChinaZPptHttpProxy() error {
	pageMax := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
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
					ChinaZPptHttpProxyUrlArr = append(ChinaZPptHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					ChinaZPptHttpProxyUrlArr = append(ChinaZPptHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func ChinaZPptSetHttpProxy() (httpclient *http.Client) {
	if ChinaZPptHttpProxyUrl == "" {
		if len(ChinaZPptHttpProxyUrlArr) <= 0 {
			err := ChinaZPptHttpProxy()
			if err != nil {
				ChinaZPptSetHttpProxy()
			}
		}
		ChinaZPptHttpProxyUrl = ChinaZPptHttpProxyUrlArr[0]
		if len(ChinaZPptHttpProxyUrlArr) >= 2 {
			ChinaZPptHttpProxyUrlArr = ChinaZPptHttpProxyUrlArr[1:]
		} else {
			ChinaZPptHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(ChinaZPptHttpProxyUrl)
	ProxyURL, _ := url.Parse(ChinaZPptHttpProxyUrl)
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

var ChinaZPptListCookie = "Hm_lvt_52c68b6edf9e9c811ffb29aafb14cd06=1731469471; Hm_lpvt_52c68b6edf9e9c811ffb29aafb14cd06=1731469566; Hm_lvt_398913ed58c9e7dfe9695953fb7b6799=1731469134; HMACCOUNT=2CEC63D57647BCA5; ASP.NET_SessionId=hkkpslgmloi04mxxkbm30tlj; _clck=1uyxgz2%7C2%7Cfqu%7C0%7C1778; cz_statistics_visitor=fbb756bb-351a-1934-6bac-264e556d7aeb; Hm_lvt_52c68b6edf9e9c811ffb29aafb14cd06=1731469471,1731469856; Hm_lpvt_52c68b6edf9e9c811ffb29aafb14cd06=1731469856; Hm_lpvt_398913ed58c9e7dfe9695953fb7b6799=1731472242; _clsk=1jl5zoh%7C1731472243396%7C11%7C1%7Cp.clarity.ms%2Fcollect"

var DownChinaZPptNextPageSleep = 10

// ychEduSpider 下载站长素材ppt模板
// @Title 下载站长素材ppt模板
// @Description https://sc.chinaz.com/，下载站长素材ppt模板
func main() {
	curPage := 1
	isPageListGo := true
	for isPageListGo {
		pageListUrl := "https://sc.chinaz.com/ppt/"
		if curPage > 1 {
			pageListUrl = fmt.Sprintf("https://sc.chinaz.com/ppt/index_%d.html", curPage)
		}
		fmt.Println(pageListUrl)
		pageListDoc, err := QueryChinaZPptList(pageListUrl)
		if err != nil {
			ChinaZPptHttpProxyUrl = ""
			fmt.Println(err)
			continue
		}
		liNodes := htmlquery.Find(pageListDoc, `//div[@class="ppt-list  masonry"]/div[@class="item masonry-brick"]`)
		if len(liNodes) <= 0 {
			break
		}
		for _, liNode := range liNodes {

			TitleNode := htmlquery.FindOne(liNode, `./div[@class="bot-div"]/a[@class="name"]`)
			Title := htmlquery.SelectAttr(TitleNode, "title")
			fmt.Println(Title)

			filePath := "../sc.chinaz.com/ppt/" + Title + ".pdf"
			if _, err := os.Stat(filePath); err != nil {
				ViewChinaZPptDetailUrl := "https://sc.chinaz.com" + htmlquery.SelectAttr(TitleNode, "href")
				QueryChinaZPptDownLoadUrlDoc, err := ViewChinaZPptDetail(ViewChinaZPptDetailUrl, pageListUrl)
				if err != nil {
					ChinaZPptHttpProxyUrl = ""
					fmt.Println(err)
					continue
				}

				DownloadNode := htmlquery.FindOne(QueryChinaZPptDownLoadUrlDoc, `//div[@class="container ppt-detail clearfix"]/div[@class="right-div"]/div[@class="new-btn-div"]/a[@class="a-download-btn"]`)
				if DownloadNode == nil {
					fmt.Println("没有下载按钮，跳过")
					continue
				}
				downloadUrl := htmlquery.SelectAttr(DownloadNode, "href")
				fmt.Println(downloadUrl)

				fmt.Println("=======开始下载" + Title + "========")
				err = DownLoadChinaZPpt(downloadUrl, ViewChinaZPptDetailUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")
				DownLoadChinaZPptTimeSleep := rand.Intn(20)
				for i := 1; i <= DownLoadChinaZPptTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page="+strconv.Itoa(curPage)+"===========下载", Title, "成功，暂停", DownLoadChinaZPptTimeSleep, "秒，倒计时", i, "秒===========")
				}
			}
		}
		curPage++
		for i := 1; i <= DownChinaZPptNextPageSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("===========翻", curPage, "页，暂停", DownChinaZPptNextPageSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

func QueryChinaZPptList(requestUrl string) (doc *html.Node, err error) {
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
	if ChinaZPptEnableHttpProxy {
		client = ChinaZPptSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", ChinaZPptListCookie)
	req.Header.Set("Host", "sc.chinaz.com")
	req.Header.Set("Origin", "https://sc.chinaz.com")
	req.Header.Set("Referer", "https://sc.chinaz.com/ppt/")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
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

func ViewChinaZPptDetail(requestUrl string, referer string) (doc *html.Node, err error) {
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
	if ChinaZPptEnableHttpProxy {
		client = ChinaZPptSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", ChinaZPptListCookie)
	req.Header.Set("Host", "sc.chinaz.com")
	req.Header.Set("Origin", "https://sc.chinaz.com")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
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

func DownLoadChinaZPpt(attachmentUrl string, referer string, filePath string) error {
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
	if ChinaZPptEnableHttpProxy {
		client = ChinaZPptSetHttpProxy()
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
	req.Header.Set("Host", "sc.chinaz.com")
	req.Header.Set("Origin", "https://sc.chinaz.com")
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
