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

type QueryNewTtBzListFormData struct {
	pageNo         int
	pageSize int
	standardStatus       int
}

type QueryNewTtBzPdfWaterMarkedFormData struct {
	operateType         int
	id string
	fileLang       string
}

var NewTtBzCookie = "__jsluid_s=3a36791cf70fd9f13dbcb2343308cbdf; HMACCOUNT=4E5B3419A3141A8E; ASP.NET_SessionId=ocqozseqhqrtjo4rbdl3sdqe; Hm_lvt_8c446e9fafe752e4975210bc30d7ab9d=1775528389; Hm_lpvt_8c446e9fafe752e4975210bc30d7ab9d=1776385138; JSESSIONID=6B92A293174566A2707CD5415213ED59; GOV_SHIRO_SESSION_ID=298bcdc5-ccf2-40b7-a883-f540ef16413b; IS_LOGGED_IN=1; MANAGER_SESSION_ID=298bcdc5-ccf2-40b7-a883-f540ef16413b"

// 下载新版团体标准文档
// @Title 下载新版团体标准文档
// @Description https://www.ttbz.org.cn/，下载新版团体标准文档
func main() {
	pageListUrl := "https://www.ttbz.org.cn/cms-proxy/ms/portal/standardInfo/getPortalStandardList"
	fmt.Println(pageListUrl)
	page := 1
	maxPage := 10
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		queryNewTtBzListFormData := QueryNewTtBzListFormData{
			pageNo:         page,
			pageSize: 10,
			standardStatus:       1,
		}
		queryNewTtBzListResponseDataRows, err := QueryNewTtBzList(pageListUrl, queryNewTtBzListFormData)
		if err != nil {
			NewTtBzHttpProxyUrl = ""
			fmt.Println(err)
		}
		for _, newTtBz := range queryNewTtBzListResponseDataRows {
			fmt.Println("=====================开始处理数据 page = ", page, "=========================")
			code := newTtBz.StandardNo
			code = strings.ReplaceAll(code, "/", "-")
			fmt.Println(code)

			title := newTtBz.StandardTitleCn
			title = strings.TrimSpace(title)
            title = strings.ReplaceAll(title, " ", "")
            title = strings.ReplaceAll(title, "　", "")
            title = strings.ReplaceAll(title, "/", "-")
            title = strings.ReplaceAll(title, "--", "-")
			fmt.Println(title)

			filePath = "../www.ttbz.org.cn/" + newTtBz.Id + "-" + title + "(" + code + ").pdf"
	        fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			// 获取加水印的pdf文件
			pdfWaterMarkedUrl := "https://www.ttbz.org.cn/cms-proxy/ms/bus/standardInfo/getStdPdfWatermarked"
			queryNewTtBzPdfWaterMarkedFormData := QueryNewTtBzPdfWaterMarkedFormData{
                operateType:         1,
                id: newTtBz.Id,
                fileLang:       "cn",
            }
            queryNewTtBzPdfWaterMarkedResponse, err := QueryNewTtBzPdfWaterMarked(pdfWaterMarkedUrl, queryNewTtBzPdfWaterMarkedFormData)
            if err != nil {
                NewTtBzHttpProxyUrl = ""
                fmt.Println(err)
            }

			fmt.Println("=======开始下载========")

			downloadUrl := "https://www.ttbz.org.cn" + queryNewTtBzPdfWaterMarkedResponse.Data
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")
			err = downloadNewTtBz(downloadUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "www.ttbz.org.cn", "temp-hbba.sacinfo.org.cn")
			err = copyNewTtBzFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			//DownLoadNewTtBzTimeSleep := 10
			DownLoadNewTtBzTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadNewTtBzTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadNewTtBzTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadNewTtBzPageTimeSleep := 10
		// DownLoadNewTtBzPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadNewTtBzPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadNewTtBzPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

type QueryNewTtBzListResponse struct {
	Data          QueryNewTtBzListResponseData `json:"data"`
	Code           int                             `json:"code"`
	Result bool                             `json:"result"`
}

type QueryNewTtBzListResponseData struct {
	Rows           []QueryNewTtBzListResponseDataRows    `json:"rows"`
	Total       int `json:"total"`
}

type QueryNewTtBzListResponseDataRows struct {
	Id           string    `json:"id"`
	StandardNo       string `json:"standardNo"`
	StandardTitleCn       string `json:"standardTitleCn"`
}

func QueryNewTtBzList(requestUrl string, queryNewTtBzListFormData QueryNewTtBzListFormData) (queryNewTtBzListResponseDataRows []QueryNewTtBzListResponseDataRows, err error) {
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
	postData.Add("pageNo", strconv.Itoa(queryNewTtBzListFormData.pageNo))
	postData.Add("pageSize", strconv.Itoa(queryNewTtBzListFormData.pageSize))
	postData.Add("standardStatus", strconv.Itoa(queryNewTtBzListFormData.standardStatus))
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	queryNewTtBzListResponse := QueryNewTtBzListResponse{}
	if err != nil {
		return queryNewTtBzListResponseDataRows, err
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
		return queryNewTtBzListResponseDataRows, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryNewTtBzListResponseDataRows, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryNewTtBzListResponseDataRows, err
	}
	err = json.Unmarshal(respBytes, &queryNewTtBzListResponse)
	if err != nil {
		return queryNewTtBzListResponseDataRows, err
	}
	queryNewTtBzListResponseDataRows = queryNewTtBzListResponse.Data.Rows
	return queryNewTtBzListResponseDataRows, nil
}

type QueryNewTtBzPdfWaterMarkedResponse struct {
	Result           bool    `json:"result"`
	Code       int `json:"code"`
	Data       string `json:"data"`
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

	queryNewTtBzPdfWaterMarkedResponse := QueryNewTtBzPdfWaterMarkedResponse{}
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
