package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
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
	Tc260EnableHttpProxy = false
	Tc260HttpProxyUrl    = "111.225.152.186:8089"
)

func Tc260SetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(Tc260HttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

// ychEduSpider 获取全国信息安全标准化技术委员会文档
// @Title 获取全国信息安全标准化技术委员会文档
// @Description https://www.tc260.org.cn/，获取全国信息安全标准化技术委员会文档
func main() {
	page := 1
	length := 10
	isPageListGo := true
	for isPageListGo {
		requestUrl := "https://www.tc260.org.cn/front/bzcx/yfgbcx.html"
		fmt.Println(requestUrl)
		start := (page - 1) * length
		YfGbCxDoc, err := YfGbCxList(requestUrl, start, length)
		if err != nil {
			fmt.Println(err)
			break
		}
		YfGbCxTableTrNodes := htmlquery.Find(YfGbCxDoc, `//div[@class="search_res_tab"]/table/tbody/tr`)
		if len(YfGbCxTableTrNodes) > 0 {
			for _, trNode := range YfGbCxTableTrNodes {
				// 标准编号
				standardNoTdNode := htmlquery.FindOne(trNode, `./td[1]`)
				standardNo := htmlquery.InnerText(standardNoTdNode)
				standardNo = strings.ReplaceAll(standardNo, "/", "-")
				fmt.Println(standardNo)

				// 中文标题
				chineseTitleTdNode := htmlquery.FindOne(trNode, `./td[2]`)
				chineseTitle := htmlquery.InnerText(chineseTitleTdNode)
				chineseTitle = strings.ReplaceAll(chineseTitle, "/", "-")
				chineseTitle = strings.ReplaceAll(chineseTitle, " ", "")
				chineseTitle = strings.ReplaceAll(chineseTitle, "：", ":")
				fmt.Println(chineseTitle)

				// 下载文档URL
				downLoadTdNode := htmlquery.FindOne(trNode, `./td[9]/a/@href`)
				downLoadUrl := htmlquery.InnerText(downLoadTdNode)
				downLoadUrl = "https://www.tc260.org.cn" + downLoadUrl
				fmt.Println(downLoadUrl)

				// 文件格式
				attachmentFormat := strings.Split(downLoadUrl, ".")
				filePath := "../www.tc260.org.cn/www.tc260.org.cn/" + chineseTitle + "(" + standardNo + ")" + "." + attachmentFormat[len(attachmentFormat)-1]
				if _, err := os.Stat(filePath); err != nil {
					fmt.Println("=======开始下载========")
					err = downloadTc260(downLoadUrl, requestUrl, filePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("=======开始完成========")
				}
				time.Sleep(time.Millisecond * 100)
			}
			page++
		} else {
			fmt.Println("没有更多分页了")
			isPageListGo = false
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func YfGbCxList(requestUrl string, start int, length int) (doc *html.Node, err error) {
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
	if Tc260EnableHttpProxy {
		client = Tc260SetHttpProxy()
	}

	postData := url.Values{}
	postData.Add("start", strconv.Itoa(start))
	postData.Add("length", strconv.Itoa(length))
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Length", "213")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "www.tc260.org.cn")
	req.Header.Set("Origin", "https://www.tc260.org.cn/")
	req.Header.Set("Referer", "https://www.tc260.org.cn/front/bzcx/yfgbcx.html")
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
		return doc, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
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

func downloadTc260(attachmentUrl string, referer string, filePath string) error {
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
	if Tc260EnableHttpProxy {
		client = Tc260SetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.tc260.org.cn")
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
