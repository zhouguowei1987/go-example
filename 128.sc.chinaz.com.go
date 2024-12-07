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
	"strings"
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

var ChinaZPpCookie = "cz_statistics_visitor=c885a4ff-3789-bf0d-cd54-0414aebde75b; Hm_lvt_398913ed58c9e7dfe9695953fb7b6799=1731474256; HMACCOUNT=487EF362690A1D5D; _clck=g33rm5%7C2%7Cfqu%7C0%7C1778; Hm_lvt_dc79411433d5171fc5e72914df433002=1731476185; Hm_lvt_ca96c3507ee04e182fb6d097cb2a1a4c=1731476983; HMACCOUNT=487EF362690A1D5D; _clsk=algv0w%7C1731478472282%7C19%7C1%7Cp.clarity.ms%2Fcollect; Hm_lpvt_ca96c3507ee04e182fb6d097cb2a1a4c=1731478524; ucvalidate=2b64f270-1efb-3ef4-d943-ec0373599f69; Access-Token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1bmlxdWVfbmFtZSI6ImNoaW5hel84MjQ2ODkxIiwibmFtZWlkIjoiNTM5NzE4OSIsIlRpY2tldCI6IjlhYzg3MDg5LWNkZjgtOWVkZi02MTYxLTI2ODkwNDA3MWFiYiIsIm5iZiI6MTczMTQ3ODU3OCwiZXhwIjoxNzMyMDgzMzc4LCJpYXQiOjE3MzE0Nzg1NzgsImlzcyI6Imh0dHA6Ly9tLnNjLmNoaW5hei5jb20vIiwiYXVkIjoiaHR0cDovL20uc2MuY2hpbmF6LmNvbS8ifQ.fDrjvH6oi4VWbLUddRp1yfrKV3oiTYfEUbSNbK-mxL0; Hm_lpvt_dc79411433d5171fc5e72914df433002=1731478577; CzScCookie=e51f1396-c270-2ca4-ae97-2eb3eda5dcae; Hm_lpvt_398913ed58c9e7dfe9695953fb7b6799=1731478577"

var DownChinaZPptNextPageSleep = 10

// ychEduSpider 下载站长素材ppt模板
// @Title 下载站长素材ppt模板
// @Description https://m.sc.chinaz.com/，下载站长素材ppt模板
func main() {
	curPage := 1
	isPageListGo := true
	for isPageListGo {
		pageListUrl := "https://m.sc.chinaz.com/ppt/"
		if curPage > 1 {
			pageListUrl = fmt.Sprintf("https://m.sc.chinaz.com/ppt/?page=%d", curPage)
		}
		fmt.Println(pageListUrl)
		pageListChinaZPpt, err := ListChinaZPpt(pageListUrl, "https://m.sc.chinaz.com/")
		if err != nil {
			ChinaZPptHttpProxyUrl = ""
			fmt.Println(err)
			continue
		}
		liNodes := htmlquery.Find(pageListChinaZPpt, `//div[@class="index-box"]/div[@class="ppt-list"]/div[@class="item"]`)
		if len(liNodes) <= 0 {
			break
		}
		for _, liNode := range liNodes {

			TitleNode := htmlquery.FindOne(liNode, `./a`)
			Title := htmlquery.SelectAttr(TitleNode, "title")
			fmt.Println(Title)

			ViewChinaZPptDetailUrl := htmlquery.SelectAttr(TitleNode, "href")
			// 转化为pc端url
			ViewChinaZPptDetailUrl = strings.Replace(ViewChinaZPptDetailUrl, "https://m.sc.chinaz.com", "https://sc.chinaz.com", 1)
			ViewChinaZPptDetailUrl = strings.Replace(ViewChinaZPptDetailUrl, ".html", ".htm", 1)
			fmt.Println(ViewChinaZPptDetailUrl)
			ViewChinaZPptDoc, err := ViewChinaZPpt(ViewChinaZPptDetailUrl, pageListUrl)
			if err != nil {
				ChinaZPptHttpProxyUrl = ""
				fmt.Println(err)
				continue
			}

			DownloadNode := htmlquery.FindOne(ViewChinaZPptDoc, `//div[@class="right-div"]/div[@class="new-btn-div"]/a[@class="a-download-btn"]`)
			if DownloadNode == nil {
				fmt.Println("没有下载按钮，跳过")
				continue
			}
			downloadUrl := htmlquery.SelectAttr(DownloadNode, "href")
			fmt.Println(downloadUrl)

			// 获取文件后缀
			downloadUrlSplitArray := strings.Split(downloadUrl, ".")
			fileSuffix := downloadUrlSplitArray[len(downloadUrlSplitArray)-1]
			filePath := "../sc.chinaz.com/ppt/" + Title + "." + fileSuffix
			fmt.Println(filePath)
			if _, err := os.Stat(filePath); err != nil {

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

func ListChinaZPpt(requestUrl string, referer string) (doc *html.Node, err error) {
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

	req.Header.Set("authority", "m.sc.chinaz.com")
	req.Header.Set("method", "GET")
	path := strings.Replace(requestUrl, "https://m.sc.chinaz.com", "", 1)
	fmt.Println(path)
	req.Header.Set("path", path)
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "ext/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", ChinaZPpCookie)
	req.Header.Set("Referer", referer)
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("sec-ch-ua", "Google Chrome\";v=\"129\", \"Not=A?Brand\";v=\"8\", \"Chromium\";v=\"129\"")
	req.Header.Set("sec-ch-ua-mobile", "?1")
	req.Header.Set("sec-ch-ua-platform", "\"Android\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Mobile Safari/537.36")
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

func ViewChinaZPpt(requestUrl string, referer string) (doc *html.Node, err error) {
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

	req.Header.Set("authority", "sc.chinaz.com")
	req.Header.Set("method", "GET")
	path := strings.Replace(requestUrl, "https://sc.chinaz.com", "", 1)
	fmt.Println(path)
	req.Header.Set("path", path)
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "ext/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", ChinaZPpCookie)
	req.Header.Set("Referer", referer)
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("sec-ch-ua", "Google Chrome\";v=\"129\", \"Not=A?Brand\";v=\"8\", \"Chromium\";v=\"129\"")
	req.Header.Set("sec-ch-ua-mobile", "?1")
	req.Header.Set("sec-ch-ua-platform", "\"Android\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Mobile Safari/537.36")
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
	req.Header.Set("authority", "sc.chinaz.com")
	req.Header.Set("method", "GET")
	req.Header.Set("path", "/ppt/")
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "sc.chinaz.com")
	req.Header.Set("Origin", "https://sc.chinaz.com")
	req.Header.Set("Referer", referer)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Mobile Safari/537.36")
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
