package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var BiaoZhunEnableHttpProxy = false
var BiaoZhunHttpProxyUrl = ""
var BiaoZhunHttpProxyUrlArr = make([]string, 0)

func BiaoZhunHttpProxy() error {
	pageMax := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	//pageMax := []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
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
					BiaoZhunHttpProxyUrlArr = append(BiaoZhunHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					BiaoZhunHttpProxyUrlArr = append(BiaoZhunHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func BiaoZhunSetHttpProxy() (httpclient *http.Client) {
	if BiaoZhunHttpProxyUrl == "" {
		if len(BiaoZhunHttpProxyUrlArr) <= 0 {
			err := BiaoZhunHttpProxy()
			if err != nil {
				BiaoZhunSetHttpProxy()
			}
		}
		BiaoZhunHttpProxyUrl = BiaoZhunHttpProxyUrlArr[0]
		if len(BiaoZhunHttpProxyUrlArr) >= 2 {
			BiaoZhunHttpProxyUrlArr = BiaoZhunHttpProxyUrlArr[1:]
		} else {
			BiaoZhunHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(BiaoZhunHttpProxyUrl)
	ProxyURL, _ := url.Parse(BiaoZhunHttpProxyUrl)
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

//var BiaoZhunNextDownloadSleep = 2

// ychEduSpider 获取国家职业技能标准
// @Title 获取国家职业技能标准
// @Description http://biaozhun.osta.org.cn/，获取国家职业技能标准
func main() {
	page := 1
	isPageListGo := true
	for isPageListGo {
		requestUrl := fmt.Sprintf("http://biaozhun.osta.org.cn/api/v1/profession/list?pageNum=%d&pageSize=20", page)
		biaoZhunListResponse, err := BiaoZhunList(requestUrl)
		if err != nil {
			fmt.Println(err)
			page = 1
			isPageListGo = false
			break
		}
		if len(biaoZhunListResponse.Rows) <= 0 {
			fmt.Println(err)
			page = 1
			isPageListGo = false
			break
		}
		for _, row := range biaoZhunListResponse.Rows {
			fmt.Println("============================================================================")
			fmt.Println("=======当前页URL", page, "========")

			id := row.Id
			fmt.Println(id)

			name := row.Name
			fmt.Println(name)

			code := row.Code
			fmt.Println(code)

			filePath := "F:\\workspace\\biaozhun.osta.org.cn\\" + name + "（" + code + "）.pdf"
			_, err = os.Stat(filePath)
			if err != nil {
				fmt.Println("=======开始下载" + strconv.Itoa(page) + "========")
				// 获取pdf文件流
				biaoZhunPdfViewDetailResponse, err := BiaoZhunPdfViewDetail(strconv.Itoa(id))
				if err != nil {
					fmt.Println(err)
					continue
				}
				pdfBytes, err := base64.StdEncoding.DecodeString(biaoZhunPdfViewDetailResponse.Data)
				if err != nil {
					fmt.Println(err)
					continue
				}

				fileDiv := filepath.Dir(filePath)
				if _, err = os.Stat(fileDiv); err != nil {
					if os.MkdirAll(fileDiv, 0777) != nil {
						fmt.Println(err)
						continue
					}
				}

				err = ioutil.WriteFile(filePath, pdfBytes, 0644)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")
				//for i := 1; i <= BiaoZhunNextDownloadSleep; i++ {
				//	time.Sleep(time.Second)
				//	fmt.Println("===========操作结束，暂停", BiaoZhunNextDownloadSleep, "秒，倒计时", i, "秒===========")
				//}
			}
		}
		page++
		isPageListGo = true
	}
}

type BiaoZhunListResponse struct {
	Code  int                        `json:"code"`
	Msg   string                     `json:"msg"`
	Rows  []BiaoZhunListResponseRows `json:"rows"`
	Total int                        `json:"total"`
}
type BiaoZhunListResponseRows struct {
	Id         int    `json:"id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	IssueNum   string `json:"issueNum"`
	IssueDate  string `json:"issueDate"`
	FileName   string `json:"fileName"`
	Attachment string `json:"attachment"`
}

func BiaoZhunList(requestUrl string) (biaoZhunListResponse BiaoZhunListResponse, err error) {
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
	if BiaoZhunEnableHttpProxy {
		client = BiaoZhunSetHttpProxy()
	}
	biaoZhunListResponse = BiaoZhunListResponse{}
	postData := url.Values{}
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return biaoZhunListResponse, err
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", "Hm_lvt_e85984af56dd04582a569a53719e397f=1703646400; _gid=GA1.3.1849962298.1703646402; _gscu_486005091=036596422fvcte80; _gscbrs_486005091=1; _gscs_486005091=036596420qvqm080|pv:2; Hm_lpvt_e85984af56dd04582a569a53719e397f=1703659702; _ga=GA1.1.346077914.1703646402; _ga_34B604LFFQ=GS1.1.1703658285.2.1.1703659792.60.0.0")
	req.Header.Set("Host", "biaozhun.osta.org.cn")
	req.Header.Set("Origin", "http://biaozhun.osta.org.cn")
	req.Header.Set("Referer", "http://biaozhun.osta.org.cn/")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return biaoZhunListResponse, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return biaoZhunListResponse, err
	}
	err = json.Unmarshal(respBytes, &biaoZhunListResponse)
	if err != nil {
		return biaoZhunListResponse, err
	}
	return biaoZhunListResponse, nil
}

type BiaoZhunPdfViewDetailResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func BiaoZhunPdfViewDetail(code string) (biaoZhunPdfViewDetailResponse BiaoZhunPdfViewDetailResponse, err error) {
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
	if BiaoZhunEnableHttpProxy {
		client = BiaoZhunSetHttpProxy()
	}
	biaoZhunPdfViewDetailResponse = BiaoZhunPdfViewDetailResponse{}
	postData := url.Values{}
	postData.Add("code", code)
	req, err := http.NewRequest("POST", "http://biaozhun.osta.org.cn/api/v1/profession/detail", strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return biaoZhunPdfViewDetailResponse, err
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", "Hm_lvt_e85984af56dd04582a569a53719e397f=1703646400; Hm_lpvt_e85984af56dd04582a569a53719e397f=1703646400; _gid=GA1.3.1849962298.1703646402; _ga=GA1.1.346077914.1703646402; _ga_34B604LFFQ=GS1.1.1703658285.2.1.1703659537.60.0.0")
	req.Header.Set("Host", "biaozhun.osta.org.cn")
	req.Header.Set("Origin", "http://biaozhun.osta.org.cn")
	req.Header.Set("Referer", "http://biaozhun.osta.org.cn/pdfview.html?code="+code)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return biaoZhunPdfViewDetailResponse, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return biaoZhunPdfViewDetailResponse, err
	}
	err = json.Unmarshal(respBytes, &biaoZhunPdfViewDetailResponse)
	if err != nil {
		return biaoZhunPdfViewDetailResponse, err
	}
	return biaoZhunPdfViewDetailResponse, nil
}
