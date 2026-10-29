package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var ShiJuanEnableHttpProxy = false
var ShiJuanHttpProxyUrl = ""
var ShiJuanHttpProxyUrlArr = make([]string, 0)

func ShiJuanHttpProxy() error {
	pageMax := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// pageMax := []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
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
					ShiJuanHttpProxyUrlArr = append(ShiJuanHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					ShiJuanHttpProxyUrlArr = append(ShiJuanHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func ShiJuanSetHttpProxy() (httpclient *http.Client) {
	if ShiJuanHttpProxyUrl == "" {
		if len(ShiJuanHttpProxyUrlArr) <= 0 {
			err := ShiJuanHttpProxy()
			if err != nil {
				ShiJuanSetHttpProxy()
			}
		}
		if len(ShiJuanHttpProxyUrlArr) > 1 {
			ShiJuanHttpProxyUrl = ShiJuanHttpProxyUrlArr[0]
		}
		if len(ShiJuanHttpProxyUrlArr) >= 2 {
			ShiJuanHttpProxyUrlArr = ShiJuanHttpProxyUrlArr[1:]
		} else {
			ShiJuanHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(ShiJuanHttpProxyUrl)
	ProxyURL, _ := url.Parse(ShiJuanHttpProxyUrl)
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
			ResponseHeaderTimeout: time.Second * 30,
		},
	}
	return httpclient
}

// ShiJuanSpider 获取第一试卷网文档
// @Title 获取第一试卷网文档
// @Description https://www.1shijuan.cn/，获取第一试卷网文档
func main() {
	// 2252
	var startId = 37131
	var endId = 75638
	for id := startId; id <= endId; id++ {
		err := ShiJuanSpider(id)
		if err != nil {
			fmt.Println(err)
		}
	}
	// err := ShiJuanSpider(75625)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func ShiJuanSpider(id int) error {
	detailDocUrl := fmt.Sprintf("https://www.1shijuan.cn/shijuan/%d.shtml", id)
	fmt.Println(detailDocUrl)
	detailPathUrl := fmt.Sprintf("/shijuan/%d.shtml", id)
	detailDoc, err := ShiJuanDetailDoc(detailDocUrl, detailPathUrl)

	if err != nil {
		return err
	}
	fmt.Println("============================================================================")
	fmt.Println("=======当前页URL", detailDocUrl, "========")

	titleNode := htmlquery.FindOne(detailDoc, `//html/body/div[2]/div[2]/div[1]/div[1]/h1`)
	if titleNode == nil {
		return errors.New("标题不存在")
	}
	title := htmlquery.InnerText(titleNode)
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, "免费", "")
	title = strings.ReplaceAll(title, "-", "")
	title = strings.ReplaceAll(title, " ", "")
	title = strings.ReplaceAll(title, "|", "-")
	fmt.Println(title)

	// 年级
	stageNode := htmlquery.FindOne(detailDoc, `//html/body/div[2]/div[2]/div[1]/div[1]/div[2]/div[1]/span[2]`)
	stage := htmlquery.InnerText(stageNode)
	fmt.Println("====年级:", stage, "======")
	// 科目
	subjNode := htmlquery.FindOne(detailDoc, `//html/body/div[2]/div[2]/div[1]/div[1]/div[2]/div[1]/span[3]`)
	subj := htmlquery.InnerText(subjNode)
	fmt.Println("====科目:", subj, "======")

	// 获取下载地址
	requestShiJuanDownloadUrl := fmt.Sprintf("https://www.1shijuan.cn/download/%d.html", id)
	requestShiJuanDownloadPathUrl := fmt.Sprintf("/download/%d.html", id)
	queryShiJuanDownloadUrlResponse, err := QueryShiJuanDownloadUrl(requestShiJuanDownloadUrl, requestShiJuanDownloadPathUrl, detailDocUrl)
	if err != nil {
		return err
	}
	ShiJuanDownloadUrl := queryShiJuanDownloadUrlResponse.Url
	fmt.Println("=====下载地址：", ShiJuanDownloadUrl, "=======")

	u, err := url.Parse(ShiJuanDownloadUrl)
	if err != nil {
		return err
	}
	fileType := path.Ext(u.Path)
	filePath := "E:\\workspace\\www.1shijuan.cn\\www.1shijuan.cn\\" + stage + "\\" + strconv.Itoa(id) + "-" + title + fileType
	fmt.Println(filePath)
	_, err = os.Stat(filePath)
	if err == nil {
		return errors.New("文档已下载过，跳过")
	}

	fmt.Println("=======开始下载========")
	err = downloadShiJuan(ShiJuanDownloadUrl, filePath, detailDocUrl)
	if err != nil {
		return err
	}
	//复制文件
	tempFilePath := strings.ReplaceAll(filePath, "www.1shijuan.cn\\www.1shijuan.cn", "www.1shijuan.cn\\temp-www.1shijuan.cn")
	err = copyShiJuanFile(filePath, tempFilePath)
	if err != nil {
		return err
	}
	fmt.Println("=======完成下载========")
	DownLoadShiJuanTimeSleep := rand.Intn(5)
	// DownLoadShiJuanTimeSleep := 10
	for i := 1; i <= DownLoadShiJuanTimeSleep; i++ {
		time.Sleep(time.Second)
		fmt.Println("id="+strconv.Itoa(id)+"===========下载", title, "成功，暂停", DownLoadShiJuanTimeSleep, "秒，倒计时", i, "秒===========")
	}
	return nil
}

func ShiJuanDetailDoc(requestUrl string, path string) (doc *html.Node, err error) {
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
	}
	if ShiJuanEnableHttpProxy {
		client = ShiJuanSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}
	req.Header.Set("authority", "www.1shijuan.cn")
	req.Header.Set("method", "GET")
	req.Header.Set("path", path)
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	// req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.1shijuan.cn")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
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

type QueryShiJuanDownloadUrlResponse struct {
	Code int    `json:"code"`
	Url  string `json:"url"`
}

func QueryShiJuanDownloadUrl(requestUrl string, path string, referer string) (queryShiJuanDownloadUrlResponse QueryShiJuanDownloadUrlResponse, err error) {
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
	}
	if ShiJuanEnableHttpProxy {
		client = ShiJuanSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接
	if err != nil {
		return queryShiJuanDownloadUrlResponse, err
	}
	req.Header.Set("authority", "www.1shijuan.cn")
	req.Header.Set("method", "GET")
	req.Header.Set("path", path)
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("Referer", referer)
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryShiJuanDownloadUrlResponse, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryShiJuanDownloadUrlResponse, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryShiJuanDownloadUrlResponse, err
	}
	err = json.Unmarshal(respBytes, &queryShiJuanDownloadUrlResponse)
	if err != nil {
		return queryShiJuanDownloadUrlResponse, err
	}
	return queryShiJuanDownloadUrlResponse, nil
}

func downloadShiJuan(attachmentUrl string, filePath string, referer string) error {
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
	}
	if ShiJuanEnableHttpProxy {
		client = ShiJuanSetHttpProxy()
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
	req.Header.Set("Host", "www.1shijuan.cn")
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

func copyShiJuanFile(src, dst string) (err error) {
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
