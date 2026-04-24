package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	// "math/rand"
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

var NewTtBzEnableHttpProxy = false
var NewTtBzHttpProxyUrl = "111.225.152.186:8089"
var NewTtBzHttpProxyUrlArr = make([]string, 0)

func NewTtBzHttpProxy() error {
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
					NewTtBzHttpProxyUrlArr = append(NewTtBzHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					NewTtBzHttpProxyUrlArr = append(NewTtBzHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func NewTtBzSetHttpProxy() (httpclient *http.Client) {
	if NewTtBzHttpProxyUrl == "" {
		if len(NewTtBzHttpProxyUrlArr) <= 0 {
			err := NewTtBzHttpProxy()
			if err != nil {
				NewTtBzSetHttpProxy()
			}
		}
		NewTtBzHttpProxyUrl = NewTtBzHttpProxyUrlArr[0]
		if len(NewTtBzHttpProxyUrlArr) >= 2 {
			NewTtBzHttpProxyUrlArr = NewTtBzHttpProxyUrlArr[1:]
		} else {
			NewTtBzHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(NewTtBzHttpProxyUrl)
	ProxyURL, _ := url.Parse(NewTtBzHttpProxyUrl)
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

type QueryNewTtBzDetailFormData struct {
	id int
}

type QueryNewTtBzPdfWaterMarkedFormData struct {
	operateType int
	id          string
	fileLang    string
}

var NewTtBzCookie = "__jsluid_s=e07d2bef452c475e086d267d5d1af922; JSESSIONID=C6BB3993FC73B5158F513B6A154B7B3F; GOV_SHIRO_SESSION_ID=95304da2-b0ee-4d6a-b273-4af9699f6067; IS_LOGGED_IN=1; MANAGER_SESSION_ID=95304da2-b0ee-4d6a-b273-4af9699f6067"

// 下载新版团体标准文档
// @Title 下载新版团体标准文档
// @Description https://www.ttbz.org.cn/，下载新版团体标准文档
func main() {
	var startId = 166237
	var endId = 166510
	var pageDetailUrl = "https://www.ttbz.org.cn/cms-proxy/ms/portal/standardInfo/getPortalStandardById"
	for id := startId; id <= endId; id++ {
		fmt.Println(id)
		queryNewTtBzDetailFormData := QueryNewTtBzDetailFormData{
			id: id,
		}
		queryNewTtBzDetailResponseData, err := QueryNewTtBzDetail(pageDetailUrl, queryNewTtBzDetailFormData)
		if err != nil {
			NewTtBzHttpProxyUrl = ""
			fmt.Println(err)
			continue
		}
		fmt.Println("=====================开始处理数据=========================")

		// if queryNewTtBzDetailResponseData.IsOpen != 1 {
		// 	fmt.Println("文档不可预览，跳过")
		// 	continue
		// }

		if len(queryNewTtBzDetailResponseData.StandardNo) == 0 || len(queryNewTtBzDetailResponseData.StandardTitleCn) == 0 {
			fmt.Println("文档标准号或标题为空，跳过")
			continue
		}

		code := queryNewTtBzDetailResponseData.StandardNo
		code = strings.ReplaceAll(code, "/", "-")
		fmt.Println(code)

		title := queryNewTtBzDetailResponseData.StandardTitleCn
		title = strings.TrimSpace(title)
		title = strings.ReplaceAll(title, " ", "")
		title = strings.ReplaceAll(title, "　", "")
		title = strings.ReplaceAll(title, "/", "-")
		title = strings.ReplaceAll(title, "《", "")
		title = strings.ReplaceAll(title, "》", "")
		title = strings.ReplaceAll(title, "--", "-")
		title = strings.ReplaceAll(title, "——", "-")
		fmt.Println(title)

		filePath = "../www.ttbz.org.cn/" + strconv.Itoa(id) + "-" + title + "(" + code + ").pdf"
		fmt.Println(filePath)

		_, err = os.Stat(filePath)
		if err == nil {
			fmt.Println("文档已下载过，跳过")
			continue
		}

		// 获取加水印的pdf文件
		fmt.Println("=======获取加水印的pdf文件========")
		pdfWaterMarkedUrl := "https://www.ttbz.org.cn/cms-proxy/ms/bus/standardInfo/getStdPdfWatermarked"
		queryNewTtBzPdfWaterMarkedFormData := QueryNewTtBzPdfWaterMarkedFormData{
			operateType: 1,
			id:          queryNewTtBzDetailResponseData.Id,
			fileLang:    "cn",
		}
		queryNewTtBzPdfWaterMarkedResponse, err := QueryNewTtBzPdfWaterMarked(pdfWaterMarkedUrl, queryNewTtBzPdfWaterMarkedFormData)
		if err != nil {
			NewTtBzHttpProxyUrl = ""
			fmt.Println(err)
			continue
		}
		downloadUrl := "https://www.ttbz.org.cn" + queryNewTtBzPdfWaterMarkedResponse.Data
		fmt.Println(downloadUrl)

		fmt.Println("=======开始下载========")

		err = downloadNewTtBz(downloadUrl, filePath)
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
		tempFilePath := strings.ReplaceAll(filePath, "www.ttbz.org.cn", "temp-www.ttbz.org.cn")
		err = copyNewTtBzFile(filePath, tempFilePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("=======下载完成========")
		DownLoadNewTtBzTimeSleep := 10
		// DownLoadNewTtBzTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadNewTtBzTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("filePath="+filePath+"===========下载成功 暂停", DownLoadNewTtBzTimeSleep, "秒 倒计时", i, "秒===========")
		}
	}
}

type QueryNewTtBzDetailResponse struct {
	Data   QueryNewTtBzDetailResponseData `json:"data"`
	Code   int                            `json:"code"`
	Result bool                           `json:"result"`
}

type QueryNewTtBzDetailResponseData struct {
	IsOpen          int    `json:"isOpen"`
	Id              string `json:"id"`
	StandardNo      string `json:"standardNo"`
	StandardTitleCn string `json:"standardTitleCn"`
}

func QueryNewTtBzDetail(requestUrl string, queryNewTtBzDetailFormData QueryNewTtBzDetailFormData) (queryNewTtBzDetailResponseData QueryNewTtBzDetailResponseData, err error) {
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
	if NewTtBzEnableHttpProxy {
		client = NewTtBzSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("id", strconv.Itoa(queryNewTtBzDetailFormData.id))
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	queryNewTtBzDetailResponse := QueryNewTtBzDetailResponse{}
	if err != nil {
		return queryNewTtBzDetailResponseData, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Host", "www.ttbz.org.cn")
	req.Header.Set("Origin", "https://www.ttbz.org.cn")
	req.Header.Set("Referer", "https://www.ttbz.org.cn/standard.html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryNewTtBzDetailResponseData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryNewTtBzDetailResponseData, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryNewTtBzDetailResponseData, err
	}
	err = json.Unmarshal(respBytes, &queryNewTtBzDetailResponse)
	if err != nil {
		return queryNewTtBzDetailResponseData, err
	}
	queryNewTtBzDetailResponseData = queryNewTtBzDetailResponse.Data
	return queryNewTtBzDetailResponseData, nil
}

type QueryNewTtBzPdfWaterMarkedResponse struct {
	Result bool   `json:"result"`
	Code   int    `json:"code"`
	Data   string `json:"data"`
}

func QueryNewTtBzPdfWaterMarked(requestUrl string, queryNewTtBzPdfWaterMarkedFormData QueryNewTtBzPdfWaterMarkedFormData) (queryNewTtBzPdfWaterMarkedResponse QueryNewTtBzPdfWaterMarkedResponse, err error) {
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
	if NewTtBzEnableHttpProxy {
		client = NewTtBzSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("operateType", strconv.Itoa(queryNewTtBzPdfWaterMarkedFormData.operateType))
	postData.Add("id", queryNewTtBzPdfWaterMarkedFormData.id)
	postData.Add("fileLang", queryNewTtBzPdfWaterMarkedFormData.fileLang)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	queryNewTtBzPdfWaterMarkedResponse = QueryNewTtBzPdfWaterMarkedResponse{}
	if err != nil {
		return queryNewTtBzPdfWaterMarkedResponse, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", NewTtBzCookie)
	req.Header.Set("Host", "www.ttbz.org.cn")
	req.Header.Set("Origin", "https://www.ttbz.org.cn")
	req.Header.Set("Referer", "https://www.ttbz.org.cn/standard.html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryNewTtBzPdfWaterMarkedResponse, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryNewTtBzPdfWaterMarkedResponse, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryNewTtBzPdfWaterMarkedResponse, err
	}
	err = json.Unmarshal(respBytes, &queryNewTtBzPdfWaterMarkedResponse)
	if err != nil {
		return queryNewTtBzPdfWaterMarkedResponse, err
	}
	return queryNewTtBzPdfWaterMarkedResponse, nil
}

func downloadNewTtBz(attachmentUrl string, filePath string) error {
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
	if NewTtBzEnableHttpProxy {
		client = NewTtBzSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.ttbz.org.cn")
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

func copyNewTtBzFile(src, dst string) (err error) {
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
