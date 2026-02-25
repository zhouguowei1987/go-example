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
	_ "os"
	"path/filepath"
	_ "path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	_ "golang.org/x/net/html"
)

var LnSiEnableHttpProxy = false
var LnSiHttpProxyUrl = "111.225.152.186:8089"
var LnSiHttpProxyUrlArr = make([]string, 0)

func LnSiHttpProxy() error {
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
					LnSiHttpProxyUrlArr = append(LnSiHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					LnSiHttpProxyUrlArr = append(LnSiHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func LnSiSetHttpProxy() (httpclient *http.Client) {
	if LnSiHttpProxyUrl == "" {
		if len(LnSiHttpProxyUrlArr) <= 0 {
			err := LnSiHttpProxy()
			if err != nil {
				LnSiSetHttpProxy()
			}
		}
		LnSiHttpProxyUrl = LnSiHttpProxyUrlArr[0]
		if len(LnSiHttpProxyUrlArr) >= 2 {
			LnSiHttpProxyUrlArr = LnSiHttpProxyUrlArr[1:]
		} else {
			LnSiHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(LnSiHttpProxyUrl)
	ProxyURL, _ := url.Parse(LnSiHttpProxyUrl)
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

var LnSiCookie = "JSESSIONID=4AF8CD01D324AFB45518D3213598F2A4"

// 下载辽宁省地方标准文档
// @Title 下载辽宁省地方标准文档
// @Description https://www.lnsi.org:8082/，下载辽宁省地方标准文档
func main() {
	pageListUrl := "https://www.lnsi.org:8082/StandardSystem/fullText/list?pagenum=1&pagesize=99999&plannum=-1&bzmc=&bzbh=&zhidingxiuding=-1&hangye=-1&zyqcdw=&gkdw=&bzzt=-1&suoshulingyu=-1&ics=&ccs="
	fmt.Println(pageListUrl)
	queryLnSiListResponseData, err := QueryLnSiList(pageListUrl)
	if err != nil {
		LnSiHttpProxyUrl = ""
		fmt.Println(err)
	}
	for id_index, lnSiData := range queryLnSiListResponseData {
		fmt.Println("=====================开始处理数据 id_index = ", id_index, "=========================")

		title := lnSiData.BzMc
		title = strings.TrimSpace(title)
        title = strings.ReplaceAll(title, " ", "-")
        title = strings.ReplaceAll(title, "　", "-")
        title = strings.ReplaceAll(title, "：", "-")
        title = strings.ReplaceAll(title, "/", "-")
        title = strings.ReplaceAll(title, "--", "-")
        fmt.Println(title)

		code := lnSiData.BzBh
		code = strings.TrimSpace(code)
		code = strings.ReplaceAll(code, "/", "-")
		code = strings.ReplaceAll(code, "—", "-")
		fmt.Println(code)

		filePath := "../www.lnsi.org/" + title + "(" + code + ")" + ".pdf"
		fmt.Println(filePath)

		_, err = os.Stat(filePath)
		if err == nil {
			fmt.Println("文档已下载过，跳过")
			continue
		}

        handleBzBh := lnSiData.BzBh
        handleBzBh = strings.ReplaceAll(handleBzBh, "/", "_")
		handleBzBh = strings.ReplaceAll(handleBzBh, " ", "%20")
		lnSiDownloadHref := fmt.Sprintf("https://www.lnsi.org:8082/StandardSystem/profiles/SY/%s.pdf", handleBzBh)
		fmt.Println(lnSiDownloadHref)
		lnSiDownloadReferer := "https://www.lnsi.org:8082/StandardSystem/full_text.html"

		fmt.Println("=======开始下载========")
		err = downloadLnSi(lnSiDownloadHref, lnSiDownloadReferer, filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//复制文件
		tempFilePath := strings.ReplaceAll(filePath, "www.lnsi.org", "temp-dbba.sacinfo.org.cn")
		err = LnSiCopyFile(filePath, tempFilePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("=======完成下载========")
		// 		DownLoadLnSiTimeSleep := 10
		DownLoadLnSiTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadLnSiTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("title="+title+"===========下载", title, "成功，暂停", DownLoadLnSiTimeSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

type QueryLnSiListResponse struct {
	Code  int                               `json:"code"`
	Total int                               `json:"total"`
	Data  []QueryLnSiListResponseData `json:"data"`
	Message   string                            `json:"message"`
}

type QueryLnSiListResponseData struct {
	BzMc   string `json:"bzmc"`
	BzBh   string `json:"bzbh"`
}

func QueryLnSiList(requestUrl string) (queryLnSiListResponseData []QueryLnSiListResponseData, err error) {
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
	if LnSiEnableHttpProxy {
		client = LnSiSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return queryLnSiListResponseData, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", LnSiCookie)
	req.Header.Set("Host", "61.142.177.93:8083")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryLnSiListResponseData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryLnSiListResponseData, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryLnSiListResponseData, err
	}
	queryLnSiListResponse := &QueryLnSiListResponse{}
	err = json.Unmarshal(respBytes, queryLnSiListResponse)
	if err != nil {
		return queryLnSiListResponseData, err
	}
	queryLnSiListResponseData = queryLnSiListResponse.Data
	return queryLnSiListResponseData, nil
}

func downloadLnSi(attachmentUrl string, referer string, filePath string) error {
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
	if LnSiEnableHttpProxy {
		client = LnSiSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", LnSiCookie)
	req.Header.Set("Host", "www.lnsi.org:8082")
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

func LnSiCopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer func(in *os.File) {
		err := in.Close()
		if err != nil {
			return
		}
	}(in)

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)

	_, err = io.Copy(out, in)
	return
}
