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

var FsStandardEnableHttpProxy = false
var FsStandardHttpProxyUrl = "111.225.152.186:8089"
var FsStandardHttpProxyUrlArr = make([]string, 0)

func FsStandardHttpProxy() error {
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
					FsStandardHttpProxyUrlArr = append(FsStandardHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					FsStandardHttpProxyUrlArr = append(FsStandardHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func FsStandardSetHttpProxy() (httpclient *http.Client) {
	if FsStandardHttpProxyUrl == "" {
		if len(FsStandardHttpProxyUrlArr) <= 0 {
			err := FsStandardHttpProxy()
			if err != nil {
				FsStandardSetHttpProxy()
			}
		}
		FsStandardHttpProxyUrl = FsStandardHttpProxyUrlArr[0]
		if len(FsStandardHttpProxyUrlArr) >= 2 {
			FsStandardHttpProxyUrlArr = FsStandardHttpProxyUrlArr[1:]
		} else {
			FsStandardHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(FsStandardHttpProxyUrl)
	ProxyURL, _ := url.Parse(FsStandardHttpProxyUrl)
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

var FsStandardCookie = "JSESSIONID=e7d8e0b6-c6fd-474e-9373-025517f6a09e"

// 下载佛山地方标准文档
// @Title 下载佛山地方标准文档
// @Description http://www.fsstandard.org.cn/，下载佛山地方标准文档
func main() {
	pageListUrl := "http://61.142.177.93:8083/organizationStandard/queryStandard?PID=-1&page=1&limit=1000"
	fmt.Println(pageListUrl)
	queryFsStandardListResponseData, err := QueryFsStandardList(pageListUrl)
	if err != nil {
		FsStandardHttpProxyUrl = ""
		fmt.Println(err)
	}
	for id_index, fsStandardData := range queryFsStandardListResponseData {
		fmt.Println("=====================开始处理数据 id_index = ", id_index, "=========================")

		title := fsStandardData.StdCn
		title = strings.TrimSpace(title)
		title = strings.ReplaceAll(title, " ", "-")
		title = strings.ReplaceAll(title, "　", "-")
		title = strings.ReplaceAll(title, "/", "-")
		title = strings.ReplaceAll(title, "--", "-")
		fmt.Println(title)

		code := fsStandardData.StdId
		code = strings.TrimSpace(code)
		code = strings.ReplaceAll(code, "/", "-")
		code = strings.ReplaceAll(code, "—", "-")
		fmt.Println(code)

		filePath := "../www.fsstandard.org.cn/" + title + "(" + code + ")" + ".pdf"
		fmt.Println(filePath)

		_, err = os.Stat(filePath)
		if err == nil {
			fmt.Println("文档已下载过，跳过")
			continue
		}

		fsStandardDownloadHref := fmt.Sprintf("http://61.142.177.93:8083/PDFViewer/PDF/%s", fsStandardData.PdfName)
		fmt.Println(fsStandardDownloadHref)
		fsStandardDownloadReferer := "http://www.fsstandard.org.cn/contents/18/306.html"

		fmt.Println("=======开始下载========")
		err = downloadFsStandard(fsStandardDownloadHref, fsStandardDownloadReferer, filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//复制文件
		tempFilePath := strings.ReplaceAll(filePath, "../www.fsstandard.org.cn", "../upload.doc88.com/dbba.sacinfo.org.cn")
		err = FsStandardCopyFile(filePath, tempFilePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("=======完成下载========")
		// 		DownLoadFsStandardTimeSleep := 10
		DownLoadFsStandardTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadFsStandardTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("title="+title+"===========下载", title, "成功，暂停", DownLoadFsStandardTimeSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

type QueryFsStandardListResponse struct {
	Code  int                               `json:"code"`
	Count int                               `json:"count"`
	Data  []QueryFsStandardListResponseData `json:"data"`
	Msg   string                            `json:"msg"`
}

type QueryFsStandardListResponseData struct {
	PdfName string `json:"pdfname"`
	StdCn   string `json:"stdcn"`
	StdId   string `json:"stdid"`
}

func QueryFsStandardList(requestUrl string) (queryFsStandardListResponseData []QueryFsStandardListResponseData, err error) {
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
	if FsStandardEnableHttpProxy {
		client = FsStandardSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return queryFsStandardListResponseData, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", FsStandardCookie)
	req.Header.Set("Host", "61.142.177.93:8083")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryFsStandardListResponseData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryFsStandardListResponseData, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryFsStandardListResponseData, err
	}
	queryFsStandardListResponse := &QueryFsStandardListResponse{}
	err = json.Unmarshal(respBytes, queryFsStandardListResponse)
	if err != nil {
		return queryFsStandardListResponseData, err
	}
	queryFsStandardListResponseData = queryFsStandardListResponse.Data
	return queryFsStandardListResponseData, nil
}

func downloadFsStandard(attachmentUrl string, referer string, filePath string) error {
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
	if FsStandardEnableHttpProxy {
		client = FsStandardSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", FsStandardCookie)
	req.Header.Set("Host", "61.142.177.93:8083")
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

func FsStandardCopyFile(src, dst string) (err error) {
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
