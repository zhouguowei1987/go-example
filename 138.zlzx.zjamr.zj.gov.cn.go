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

const (
	ZjAmrEnableHttpProxy = false
	ZjAmrHttpProxyUrl    = "111.225.152.186:8089"
)

func ZjAmrSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(ZjAmrHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

var ZjAmrCookie = "node_id=nginx_1; node_id=nginx_1; _d_id=30fbbb89f633bbccfddce4c793e63e"

// 获取浙江标准在线
// @Title 获取浙江标准在线
// @Description http://zlzx.zjamr.zj.gov.cn/，获取浙江标准在线
func main() {
	maxPage := 3197
	page := 1
	isPageListGo := true
	for isPageListGo {
		requestUrl := fmt.Sprintf("https://zlzx.zjamr.zj.gov.cn/bzzx/public/news/list/BZBP/ALL/%d.html", page)
		fmt.Println(requestUrl)

		pageDoc, err := htmlquery.LoadURL(requestUrl)
		if err != nil {
			fmt.Println(err)
		}
		aNodes := htmlquery.Find(pageDoc, `//div[@class="xwlbs-cc"]/div[@class="xwlbs"]/ul[@class="tzgg_bj"]/a`)
		if len(aNodes) <= 0 {
			isPageListGo = false
			break
		}

		for _, aNode := range aNodes {
			title := htmlquery.InnerText(htmlquery.FindOne(aNode, `./li/span[@class="news-title"]`))
			titleArray := strings.Split(title, "  |  ")
			title = strings.TrimSpace(titleArray[1])
			fmt.Println(title)

			code := strings.TrimSpace(titleArray[0])
			code = strings.ReplaceAll(code, "/", "-")
			fmt.Println(code)

			filePath := "../zlzx.zjamr.zj.gov.cn/地方标准/" + title + "(" + code + ")" + ".pdf"
			fmt.Println(filePath)
			if _, err := os.Stat(filePath); err != nil {

				aHrefNode := htmlquery.FindOne(aNode, `./@href`)
				if aHrefNode == nil {
					continue
				}
				aHrefUrl := htmlquery.InnerText(aHrefNode)
				aHrefUrl = "https://zlzx.zjamr.zj.gov.cn" + aHrefUrl
				fmt.Println(aHrefUrl)

				detailDoc, err := ZjAmrDetailDoc(aHrefUrl, "https://zlzx.zjamr.zj.gov.cn/bzzx/public/news/list/lstd/ALL/1.html")
				if err != nil {
					fmt.Println(err)
					break
				}

				// /html/body/div[2]/div/div/ul/table/tbody/tr[2]/td[2]/a
				downloadNode := htmlquery.FindOne(detailDoc, `//html/body/div[2]/div/div/ul/table/tbody/tr[2]/td[2]/a/@onclick`)
				if downloadNode == nil {
					continue
				}

				downloadText := htmlquery.InnerText(downloadNode)
				downloadTextArray := strings.Split(downloadText, "\"")
				// 				fmt.Println(downloadTextArray)
				downloadUrl := downloadTextArray[1]
				downloadUrl = strings.ReplaceAll(downloadUrl, "\\/", "/")
				fmt.Println(downloadUrl)

				fmt.Println("=======开始下载" + title + "========")
				err = downloadZjAmr(downloadUrl, aHrefUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")
				//DownLoadZjAmrTimeSleep := 10
				DownLoadZjAmrTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadZjAmrTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page="+strconv.Itoa(page)+"==========="+filePath+"下载成功，暂停", DownLoadZjAmrTimeSleep, "秒，倒计时", i, "秒===========")
				}

			}
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}
func ZjAmrDetailDoc(requestUrl string, referer string) (doc *html.Node, err error) {
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
	if ZjAmrEnableHttpProxy {
		client = ZjAmrSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", ZjAmrCookie)
	req.Header.Set("Referer", referer)
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

func downloadZjAmr(attachmentUrl string, referer string, filePath string) error {
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
	if ZjAmrEnableHttpProxy {
		client = ZjAmrSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", ZjAmrCookie)
	req.Header.Set("Host", "zlzx.zjamr.zj.gov.cn")
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
