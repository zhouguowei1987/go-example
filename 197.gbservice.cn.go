package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	// 	"math/rand"
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

var GbServiceEnableHttpProxy = false
var GbServiceHttpProxyUrl = "111.225.152.186:8089"
var GbServiceHttpProxyUrlArr = make([]string, 0)

func GbServiceHttpProxy() error {
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
					GbServiceHttpProxyUrlArr = append(GbServiceHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					GbServiceHttpProxyUrlArr = append(GbServiceHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func GbServiceSetHttpProxy() (httpclient *http.Client) {
	if GbServiceHttpProxyUrl == "" {
		if len(GbServiceHttpProxyUrlArr) <= 0 {
			err := GbServiceHttpProxy()
			if err != nil {
				GbServiceSetHttpProxy()
			}
		}
		GbServiceHttpProxyUrl = GbServiceHttpProxyUrlArr[0]
		if len(GbServiceHttpProxyUrlArr) >= 2 {
			GbServiceHttpProxyUrlArr = GbServiceHttpProxyUrlArr[1:]
		} else {
			GbServiceHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(GbServiceHttpProxyUrl)
	ProxyURL, _ := url.Parse(GbServiceHttpProxyUrl)
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

var GbServiceCookie = "Hm_lvt_25c82fafe27bde86759c4030e5427888=1781236836; HMACCOUNT=9C0CD19686802BBF; Hm_lpvt_25c82fafe27bde86759c4030e5427888=1781237124"

// 下载江苏建设科技服务网标准文档
// @Title 下载江苏建设科技服务网标准文档
// @Description https://gbservice.cn/，下载江苏建设科技服务网标准文档
func main() {
	// 1332
	var startId = 1332
	var endId = 2094
	for id := startId; id <= endId; id++ {
		var pageDetailUrl = fmt.Sprintf("https://gbservice.cn/api/standard/openinfo?id=%d", id)
		var pageDetailPathUrl = fmt.Sprintf("/api/standard/openinfo?id=%d", id)
		fmt.Println(pageDetailUrl)
		queryGbServiceDetailResponseData, err := QueryGbServiceDetail(pageDetailUrl, pageDetailPathUrl)
		fmt.Println(queryGbServiceDetailResponseData)
		if len(queryGbServiceDetailResponseData.Files) <= 0 {
			fmt.Println("没有数据，跳过")
			continue
		}
		if err != nil {
			GbServiceHttpProxyUrl = ""
			fmt.Println(err)
			continue
		}
		fmt.Println("=====================开始处理数据=========================")

		code := queryGbServiceDetailResponseData.Basicid
		code = strings.ReplaceAll(code, "/", "-")
		fmt.Println(code)

		title := queryGbServiceDetailResponseData.Name
		title = strings.TrimSpace(title)
		title = strings.ReplaceAll(title, " ", "")
		title = strings.ReplaceAll(title, "　", "")
		title = strings.ReplaceAll(title, "/", "-")
		title = strings.ReplaceAll(title, "《", "")
		title = strings.ReplaceAll(title, "》", "")
		title = strings.ReplaceAll(title, "--", "-")
		title = strings.ReplaceAll(title, "——", "-")
		fmt.Println(title)

		filePath := "../gbservice.cn/" + title + "(" + code + ").pdf"
		if len(code) <= 0 {
			filePath = "../gbservice.cn/" + title + ".pdf"
		}
		fmt.Println(filePath)

		_, err = os.Stat(filePath)
		if err == nil {
			fmt.Println("文档已下载过，跳过")
			continue
		}
		downloadUrl := queryGbServiceDetailResponseData.Files[0].Url
		fmt.Println(downloadUrl)

		fmt.Println("=======开始下载========")

		err = downloadGbService(downloadUrl, filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// 查看文件大小，如果是空文件，则删除
		fileInfo, err := os.Stat(filePath)
		if err == nil && fileInfo.Size() == 0 || fileInfo.Size() == 228896 {
			fmt.Println("空文件删除")
			err = os.Remove(filePath)
		}
		if err != nil {
			continue
		}
		//复制文件
		tempFilePath := strings.ReplaceAll(filePath, "gbservice.cn", "temp-gbservice.cn")
		err = copyGbServiceFile(filePath, tempFilePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("=======下载完成========")
		DownLoadGbServiceTimeSleep := 10
		// DownLoadGbServiceTimeSleep := rand.Intn(6)
		for i := 1; i <= DownLoadGbServiceTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("filePath="+filePath+"===========下载成功 暂停", DownLoadGbServiceTimeSleep, "秒 倒计时", i, "秒===========")
		}
	}
}

type QueryGbServiceDetailResponse struct {
	Data    QueryGbServiceDetailResponseData `json:"data"`
	Code    int                              `json:"code"`
	Message string                           `json:"message"`
}

type QueryGbServiceDetailResponseData struct {
	Files   []QueryGbServiceDetailResponseDataFiles `json:"files"`
	Basicid string                                  `json:"basicid"`
	Name    string                                  `json:"name"`
}
type QueryGbServiceDetailResponseDataFiles struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

func QueryGbServiceDetail(requestUrl string, path string) (queryGbServiceDetailResponseData QueryGbServiceDetailResponseData, err error) {
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
	if GbServiceEnableHttpProxy {
		client = GbServiceSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	queryGbServiceDetailResponse := QueryGbServiceDetailResponse{}
	if err != nil {
		return queryGbServiceDetailResponseData, err
	}

	req.Header.Set("authority", "gbservice.cn")
	req.Header.Set("method", "GET")
	req.Header.Set("path", path)
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", GbServiceCookie)
	req.Header.Set("Host", "gbservice.cn")
	req.Header.Set("Origin", "https://gbservice.cn")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryGbServiceDetailResponseData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryGbServiceDetailResponseData, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryGbServiceDetailResponseData, err
	}
	err = json.Unmarshal(respBytes, &queryGbServiceDetailResponse)
	if err != nil {
		return queryGbServiceDetailResponseData, err
	}
	queryGbServiceDetailResponseData = queryGbServiceDetailResponse.Data
	return queryGbServiceDetailResponseData, nil
}

func downloadGbService(attachmentUrl string, filePath string) error {
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
	if GbServiceEnableHttpProxy {
		client = GbServiceSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "gbservice.cn")
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

func copyGbServiceFile(src, dst string) (err error) {
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
