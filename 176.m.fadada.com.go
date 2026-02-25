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

var FaDaDaEnableHttpProxy = false
var FaDaDaHttpProxyUrl = "111.225.152.186:8089"
var FaDaDaHttpProxyUrlArr = make([]string, 0)

func FaDaDaHttpProxy() error {
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
					FaDaDaHttpProxyUrlArr = append(FaDaDaHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					FaDaDaHttpProxyUrlArr = append(FaDaDaHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func FaDaDaSetHttpProxy() (httpclient *http.Client) {
	if FaDaDaHttpProxyUrl == "" {
		if len(FaDaDaHttpProxyUrlArr) <= 0 {
			err := FaDaDaHttpProxy()
			if err != nil {
				FaDaDaSetHttpProxy()
			}
		}
		FaDaDaHttpProxyUrl = FaDaDaHttpProxyUrlArr[0]
		if len(FaDaDaHttpProxyUrlArr) >= 2 {
			FaDaDaHttpProxyUrlArr = FaDaDaHttpProxyUrlArr[1:]
		} else {
			FaDaDaHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(FaDaDaHttpProxyUrl)
	ProxyURL, _ := url.Parse(FaDaDaHttpProxyUrl)
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

type QueryFaDaDaListFormData struct {
	currentPageNo int
	templateName  string
	commonTypeId  string
	templateType  string
	pageSize      int
}

var FaDaDaCookie = "gr_user_id=c3abe450-dfee-4004-9aa6-9e428055f1ee; a77c4e3f47ba1ba5_gr_last_sent_cs1=54066b4bdb80aa490dfd91fb42e26cc4; __jsluid_s=efe2ec4862c805d6df67e09d884e9756; a77c4e3f47ba1ba5_gr_cs1=54066b4bdb80aa490dfd91fb42e26cc4; Hm_lvt_3f254cfc2bb960a7048945fd36d3450e=1772021306; Hm_lpvt_3f254cfc2bb960a7048945fd36d3450e=1772021306; HMACCOUNT=1CCD0111717619C6; Qs_lvt_476583=1772021305; Qs_pv_476583=4418630159139847700; SESSION=OTcxY2Y2ODYtOWI4NS00Y2VjLWI1ZjctOTExZjYzYTY1OGE0; tgc_=D9dYTjNQ7fbOUbjrWElanpRc6WdGYlrQ1lyyqkJ3uNGrxdKWKR3iZU682EMIcFMh3/2werCq5tA5ANziDUC8Tw==; tc_=suQ62ZsNkoLOevJcabbvwQQgUlTDbknUjRD+wNT1cMlkWgdPHweRWLGdAEBFKCgYtbOn3xbd/fLXh6uecYy/Hh2IkFtCoS/Iy//f4EKLu+Q="

// 下载法大大合同模板文档
// @Title 下载法大大合同模板文档
// @Description https://m.fadada.com/，下载法大大合同模板文档
func main() {
	pageListUrl := "https://cloud.fadada.com/api/portal/contractTemplate"
	fmt.Println(pageListUrl)
	page := 1
	maxPage := 87
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		queryFaDaDaListFormData := QueryFaDaDaListFormData{
			currentPageNo: page,
			templateName:  "",
			commonTypeId:  "",
			templateType:  "",
			pageSize:      12,
		}
		queryFaDaDaListResponseDataListPageDataList, err := QueryFaDaDaList(pageListUrl, queryFaDaDaListFormData)
		if err != nil {
			FaDaDaHttpProxyUrl = ""
			fmt.Println(err)
		}
		for _, faDaDa := range queryFaDaDaListResponseDataListPageDataList {
			fmt.Println("=====================开始处理数据 page = ", page, "=========================")

			title := faDaDa.TemplateName
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, "-", "")
			title = strings.ReplaceAll(title, " ", "")
			title = strings.ReplaceAll(title, "|", "-")
			fmt.Println(title)

			detailUrl := fmt.Sprintf("https://www.fadada.com/hetongmuban/detail-%s/", faDaDa.Id)
			fmt.Println(detailUrl)

			filePath := "../m.fadada.com/" + title + ".doc"
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载========")

			downloadUrl := fmt.Sprintf("https://cloud.fadada.com/api/portal/downloadTemplate?contractTemplateId=%s", faDaDa.Id)
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")
			err = downloadFaDaDa(downloadUrl, detailUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "m.fadada.com", "temp-m.fadada.com")
			err = copyFaDaDaFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			//DownLoadFaDaDaTimeSleep := 10
			DownLoadFaDaDaTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadFaDaDaTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadFaDaDaTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadFaDaDaPageTimeSleep := 10
		// DownLoadFaDaDaPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadFaDaDaPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadFaDaDaPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

type QueryFaDaDaListResponse struct {
	Code    string                      `json:"code"`
	Data    QueryFaDaDaListResponseData `json:"data"`
	Message string                      `json:"message"`
	Success bool                        `json:"success"`
}

type QueryFaDaDaListResponseData struct {
	ListPage QueryFaDaDaListResponseDataListPage `json:listPage`
}

type QueryFaDaDaListResponseDataListPage struct {
	DataList []QueryFaDaDaListResponseDataListPageDataList `json:dataList`
}

type QueryFaDaDaListResponseDataListPageDataList struct {
	Id           string `json:"id"`
	TemplateName string `json:"templateName"`
}

func QueryFaDaDaList(requestUrl string, queryFaDaDaListFormData QueryFaDaDaListFormData) (queryFaDaDaListResponseDataListPageDataList []QueryFaDaDaListResponseDataListPageDataList, err error) {
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
	if FaDaDaEnableHttpProxy {
		client = FaDaDaSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("currentPageNo", strconv.Itoa(queryFaDaDaListFormData.currentPageNo))
	postData.Add("pageSize", strconv.Itoa(queryFaDaDaListFormData.pageSize))
	postData.Add("templateName", queryFaDaDaListFormData.templateName)
	postData.Add("commonTypeId", queryFaDaDaListFormData.commonTypeId)
	postData.Add("templateType", queryFaDaDaListFormData.templateType)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	queryFaDaDaListResponse := QueryFaDaDaListResponse{}
	if err != nil {
		return queryFaDaDaListResponseDataListPageDataList, err
	}

	req.Header.Set("authority", "cloud.fadada.com")
	req.Header.Set("method", "POST")
	req.Header.Set("path", "/api/portal/contractTemplate")
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", FaDaDaCookie)
	req.Header.Set("Host", "m.fadada.com")
	req.Header.Set("Origin", "https://m.fadada.com")
	req.Header.Set("Referer", "https://m.fadada.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryFaDaDaListResponseDataListPageDataList, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryFaDaDaListResponseDataListPageDataList, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryFaDaDaListResponseDataListPageDataList, err
	}
	err = json.Unmarshal(respBytes, &queryFaDaDaListResponse)
	if err != nil {
		return queryFaDaDaListResponseDataListPageDataList, err
	}
	queryFaDaDaListResponseDataListPageDataList = queryFaDaDaListResponse.Data.ListPage.DataList
	return queryFaDaDaListResponseDataListPageDataList, nil
}

func downloadFaDaDa(attachmentUrl string, referer string, filePath string) error {
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
	if FaDaDaEnableHttpProxy {
		client = FaDaDaSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "m.fadada.com")
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

func copyFaDaDaFile(src, dst string) (err error) {
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
