package main

import (
	"bytes"
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

var StdMiItEnableHttpProxy = false
var StdMiItHttpProxyUrl = "111.225.152.186:8089"
var StdMiItHttpProxyUrlArr = make([]string, 0)

func StdMiItHttpProxy() error {
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
					StdMiItHttpProxyUrlArr = append(StdMiItHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					StdMiItHttpProxyUrlArr = append(StdMiItHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func StdMiItSetHttpProxy() (httpclient *http.Client) {
	if StdMiItHttpProxyUrl == "" {
		if len(StdMiItHttpProxyUrlArr) <= 0 {
			err := StdMiItHttpProxy()
			if err != nil {
				StdMiItSetHttpProxy()
			}
		}
		StdMiItHttpProxyUrl = StdMiItHttpProxyUrlArr[0]
		if len(StdMiItHttpProxyUrlArr) >= 2 {
			StdMiItHttpProxyUrlArr = StdMiItHttpProxyUrlArr[1:]
		} else {
			StdMiItHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(StdMiItHttpProxyUrl)
	ProxyURL, _ := url.Parse(StdMiItHttpProxyUrl)
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

var StdMiItCookie = "__jsluid_s=8355c056978084a6d5cc9c4ac8b7faf5; ariauseGraymode=false; SID=F497C745BDEC466B98C5C12CCF4C49C1"

type QueryStdMiItListRequestPayload struct {
	BpiBzno       string `json:"bpiBzno"`
	Gbhb          string `json:"gbhb"`
	PageNo        int    `json:"pageNo"`
	PageSize      int    `json:"pageSize"`
	PiProjectname string `json:"piProjectname"`
	PiStatus      string `json:"piStatus"`
}

// 下载工业和信息化标准信息服务平台文档
// @Title 下载工业和信息化标准信息服务平台文档
// @Description https://std.miit.gov.cn/，下载工业和信息化标准信息服务平台文档
func main() {
	pageListUrl := "https://std.miit.gov.cn/kjsStandproject/front/project/queryStandardsByPage"
	fmt.Println(pageListUrl)
	//page := 550
	page := 1
	maxPage := 2611
	rows := 15
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		queryStdMiItListRequestPayload := QueryStdMiItListRequestPayload{
			BpiBzno:       "",
			Gbhb:          "ALL",
			PageNo:        page,
			PageSize:      rows,
			PiProjectname: "",
			PiStatus:      "xx",
		}
		queryStdMiItListResponseDataObjects, err := QueryStdMiItList(pageListUrl, queryStdMiItListRequestPayload)
		if err != nil {
			fmt.Println(err)
			break
		}
		for _, stdMiIt := range queryStdMiItListResponseDataObjects {
			fmt.Println("=====================开始处理数据 page = ", page, "=========================")

			code := stdMiIt.BpiBzno
			code = strings.ReplaceAll(code, "/", "-")
			code = strings.ReplaceAll(code, "—", "-")
			fmt.Println(code)

			title := stdMiIt.PiProjectname
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "--", "-")
			fmt.Println(title)

			filePath := "../std.miit.gov.cn/" + "(" + code + ")" + title + ".pdf"
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载========")
			downloadUrl := fmt.Sprintf("https://std.miit.gov.cn/kjsStandproject/front/project/pdf?bzNo=%s", stdMiIt.BpiBzno)
			downloadUrl = strings.ReplaceAll(downloadUrl, " ", "%20")
			fmt.Println(downloadUrl)
			err = downloadStdMiIt(downloadUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// 查看文件大小，如果是空文件，则删除
			fileInfo, err := os.Stat(filePath)
			if err == nil && fileInfo.Size() == 0 {
				fmt.Println("空文件删除")
				err = os.Remove(filePath)
			}
			if err != nil {
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "std.miit.gov.cn", "temp-std.miit.gov.cn")
			err = copyStdMiItFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			//DownLoadStdMiItTimeSleep := 10
			DownLoadStdMiItTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadStdMiItTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadStdMiItTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadStdMiItPageTimeSleep := 10
		// DownLoadStdMiItPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadStdMiItPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadStdMiItPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

type QueryStdMiItListResponse struct {
	Code int                          `json:"code"`
	Data QueryStdMiItListResponseData `json:"data"`
	Msg  string                       `json:"msg"`
}

type QueryStdMiItListResponseData struct {
	Objects []QueryStdMiItListResponseDataObjects `json:"objects"`
	Total   int                                   `json:"total"`
}

type QueryStdMiItListResponseDataObjects struct {
	BpiBzno       string `json:"bpiBzno"`
	PiProjectname string `json:"piProjectname"`
}

func QueryStdMiItList(requestUrl string, queryStdMiItListRequestPayload QueryStdMiItListRequestPayload) (queryStdMiItListResponseDataObjects []QueryStdMiItListResponseDataObjects, err error) {
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
	if StdMiItEnableHttpProxy {
		client = StdMiItSetHttpProxy()
	}
	// 将数据编码为JSON格式
	queryStdMiItListRequestPayloadJson, err := json.Marshal(queryStdMiItListRequestPayload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 创建字符串读取器，这是http.Request需要的类型
	body := bytes.NewReader(queryStdMiItListRequestPayloadJson)
	req, err := http.NewRequest("POST", requestUrl, body) //建立连接

	queryStdMiItListResponse := QueryStdMiItListResponse{}
	if err != nil {
		return queryStdMiItListResponseDataObjects, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", StdMiItCookie)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Host", "std.miit.gov.cn")
	req.Header.Set("Origin", "https://std.miit.gov.cn")
	req.Header.Set("Referer", "https://std.miit.gov.cn/")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryStdMiItListResponseDataObjects, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryStdMiItListResponseDataObjects, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryStdMiItListResponseDataObjects, err
	}
	err = json.Unmarshal(respBytes, &queryStdMiItListResponse)
	if err != nil {
		return queryStdMiItListResponseDataObjects, err
	}
	queryStdMiItListResponseDataObjects = queryStdMiItListResponse.Data.Objects
	return queryStdMiItListResponseDataObjects, nil
}

func downloadStdMiIt(attachmentUrl string, filePath string) error {
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
	if StdMiItEnableHttpProxy {
		client = StdMiItSetHttpProxy()
	} //初始化客户端                     //初始化客户端
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", StdMiItCookie)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "std.miit.gov.cn")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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

func copyStdMiItFile(src, dst string) (err error) {
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
