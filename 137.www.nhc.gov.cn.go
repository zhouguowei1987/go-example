package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/html"
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

var NhcGovEnableHttpProxy = false
var NhcGovHttpProxyUrl = "111.225.152.186:8089"
var NhcGovHttpProxyUrlArr = make([]string, 0)

func NhcGovHttpProxy() error {
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
					NhcGovHttpProxyUrlArr = append(NhcGovHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					NhcGovHttpProxyUrlArr = append(NhcGovHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func NhcGovSetHttpProxy() (httpclient *http.Client) {
	if NhcGovHttpProxyUrl == "" {
		if len(NhcGovHttpProxyUrlArr) <= 0 {
			err := NhcGovHttpProxy()
			if err != nil {
				NhcGovSetHttpProxy()
			}
		}
		NhcGovHttpProxyUrl = NhcGovHttpProxyUrlArr[0]
		if len(NhcGovHttpProxyUrlArr) >= 2 {
			NhcGovHttpProxyUrlArr = NhcGovHttpProxyUrlArr[1:]
		} else {
			NhcGovHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(NhcGovHttpProxyUrl)
	ProxyURL, _ := url.Parse(NhcGovHttpProxyUrl)
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

// _yfxkpy_firsttime 1748789542542 2025-07-16 11:14:03
// _yfxkpy_lasttime 1752635422911 2025-07-16 11:10:22
// _yfxkpy_visittime 1752635422911 2025-07-16 11:10:22
var NhcGovCookie = "sVoELocvxVW0S=5ZCewgg5x.lwIRUfd5Os_9gZNFxtFa3nuYdA6vgRiwPcOclyvF_tK2qwg5N98WYFDgyyIOHSycpJgCLrjkmhdJq; 5uRo8RWcod0KO=60kIcUTd7w2UNAs9nmzRXcLX.ZG6P3v7yRuSgPR7Q1p9UIri_QA93ERV7uhjN.wrrRbU7idYR6nBhkaKcYXkfshG; enable_5uRo8RWcod0K=true; JSESSIONID=9C480E27499298CB3CBA83E5128D83FB; ariawapForceOldFixed=false; arialoadData=true; 5uRo8RWcod0KS=60HGjkwJ4ZXxqmcF7HAlk9ozCjyXTF2WM2OLgPyy.KcwOGsF0W.atGC9CJfZOy1nCxfY.idZiUCtCDpRFTeHxI9q; 5uRo8RWcod0KT=0qMWrH0FvMch3bjl5J3EfYSJ0nMKpET4uyIfrAbMpCu7e3VcDNS93HUUHD4q30b4xkORvzvSitYyHSj087Grmftv9N1R64xesST_hd8IchFtvL.JG7hEeO0cmEYoZIyLqT6wHlsT0lKwdHcc_PL8OLQt613ufCZxmdOSprfFQTxTj_dgWAnp7wxapG0XbJAgbrWmyDzA26dwLwoAmwYdOakXpMCoYKQxiM1YgiZnzwBg7ABoEfoQJY6yfH7lHWL_yjO6buMS27jPhxoC.YkYdtq; _yfxkpy_ssid_10006654=%7B%22_yfxkpy_firsttime%22%3A%221737352968986%22%2C%22_yfxkpy_lasttime%22%3A%221752595494940%22%2C%22_yfxkpy_visittime%22%3A%221752596019158%22%2C%22_yfxkpy_cookie%22%3A%2220250120140248997354195004763129%22%2C%22_yfxkpy_returncount%22%3A%225%22%7D; ariauseGraymode=false; ariaappid=e2dcc6e6e04fd7a68ae6f8b8a9be7f7d; 5uRo8RWcod0KP=01q4HO.QJZ1vXEaAfEXQie6FW_ODtRTU84FKLZuN7ZMlCiol_i5VaA0pHEyVvqprzVfpnqKN._E2Ksd7s9RPK1j9bZkCa3dmuUvfDnb4r2AgktbNa3_v9DOWZwvCraS3O.LgOReWU3hNyRPd1IoRCzpqAYb2SPAFlrQhIolStS36JL9lDaKIAEinx1hQ_fVAiuGp7NzfJB4hCOa1WVnVmh0aFcp6SQBuK1yf2khUisvvOuOO_bHXC17zFA0S36zj8ackGVxDHHku2VTEeRBbVN_tP9pR252MpayDVq59idRd_oLWX9GWkVRZ.GO38apcYaaATAJKbNZCDxUEhFfaIzY2cY7cgLFAaeyZdwT21cS1A1hHxM_KphaJBiS7PgOWNyA1b1jqa34JKe6LPJ8VA9C6vV9L.7Sru6FafaqXHWEG"

type struct _yfxkpy_ssid_10006654{
}

// 下载国家卫生标准文档
// @Title 下载国家卫生标准文档
// @Description https://www.nhc.gov.cn/，下载国家卫生标准文档
func main() {
	pageListUrl := "https://www.nhc.gov.cn/search/74ef62665780458e8e43027d6b5d98aa?_isAgg=true&_isJson=true&_pageSize=9999&_template=index&_rangeTimeGte=&_channelName=&page=1"
	fmt.Println(pageListUrl)
	//queryNhcGovListResponseDataResults, err := QueryNhcGovList(pageListUrl)
	queryNhcGovListResponseDataResults, err := QueryNhcGovListJson()
	if err != nil {
		NhcGovHttpProxyUrl = ""
		fmt.Println(err)
	}
	for id_index, nhcGovResult := range queryNhcGovListResponseDataResults {
		fmt.Println("=====================开始处理数据 id_index = ", id_index, "=========================")

		title := nhcGovResult.Title
		title = strings.TrimSpace(title)
		title = strings.ReplaceAll(title, "-", "")
		title = strings.ReplaceAll(title, " ", "")
		title = strings.ReplaceAll(title, "|", "-")
		fmt.Println(title)

		nhcGovDetailUrl := nhcGovResult.Url
		fmt.Println(nhcGovDetailUrl)

		nhcGovDetailDoc, err := NhcGovDetailDoc(nhcGovDetailUrl, "https://www.nhc.gov.cn/wjw/wsbzxx/wsbz.shtml")
		if err != nil {
			fmt.Println(err)
			break
		}

		codeNode := htmlquery.FindOne(nhcGovDetailDoc, `//div[@class="w1140 bgfff p20"]/div[@class="list"]/table[@class="mt20 mb20"]/tbody/tr[1]/td[@class="zhupei"]`)
		if codeNode == nil {
			fmt.Println("没有code节点，跳过")
			break
		}

		code := htmlquery.InnerText(codeNode)
		code = strings.TrimSpace(code)
		code = strings.ReplaceAll(code, "/", "-")
		code = strings.ReplaceAll(code, "\n", "")
		code = strings.ReplaceAll(code, "\r", "")
		fmt.Println(code)

		nhcGovDownloadHrefNode := htmlquery.FindOne(nhcGovDetailDoc, `//div[@class="w1140 bgfff p20"]/div[@class="list"]/div[@class="con"]/p/a/@href`)
		nhcGovDownloadHref := htmlquery.InnerText(nhcGovDownloadHrefNode)
		nhcGovDetailUrlHandleArray := strings.Split(nhcGovDetailUrl, "/")
		nhcGovDetailUrlHandleArray = nhcGovDetailUrlHandleArray[:len(nhcGovDetailUrlHandleArray)-1]
		nhcGovDownloadHref = strings.Join(nhcGovDetailUrlHandleArray, "/") + "/" + nhcGovDownloadHref
		fmt.Println(nhcGovDownloadHref)

		filePath := "../www.nhc.gov.cn/" + title + "(" + code + ")" + ".pdf"
		fmt.Println(filePath)

		_, err = os.Stat(filePath)
		if err == nil {
			fmt.Println("文档已下载过，跳过")
			break
		}

		fmt.Println("=======开始下载========")
		err = downloadNhcGov(nhcGovDownloadHref, nhcGovDetailUrl, filePath)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("=======完成下载========")
		//DownLoadNhcGovTimeSleep := 10
		DownLoadNhcGovTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadNhcGovTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("title="+title+"===========下载", title, "成功，暂停", DownLoadNhcGovTimeSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

type QueryNhcGovListResponse struct {
	ChannelName string                      `json:"channelName"`
	Data        QueryNhcGovListResponseData `json:"data"`
	LocationUrl string                      `json:"locationUrl"`
}

type QueryNhcGovListResponseData struct {
	ChannelId         string                               `json:"channelId"`
	Page              int64                                `json:"page"`
	RelateSubChannels string                               `json:"relateSubChannels"`
	Results           []QueryNhcGovListResponseDataResults `json:"results"`
	Rows              int64                                `json:"rows"`
	Total             int64                                `json:"total"`
}

type QueryNhcGovListResponseDataResults struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func QueryNhcGovListJson() (queryNhcGovListResponseDataResults []QueryNhcGovListResponseDataResults, err error) {
	nhcGovListJsonData, err := os.ReadFile("./www.nhc.gov.cn.json")
	if err != nil {
		return queryNhcGovListResponseDataResults, err
	}

	queryNhcGovListResponse := QueryNhcGovListResponse{}
	err = json.Unmarshal(nhcGovListJsonData, &queryNhcGovListResponse)
	if err != nil {
		return queryNhcGovListResponseDataResults, err
	}
	queryNhcGovListResponseDataResults = queryNhcGovListResponse.Data.Results
	return queryNhcGovListResponseDataResults, nil
}

func QueryNhcGovList(requestUrl string) (queryNhcGovListResponseDataResults []QueryNhcGovListResponseDataResults, err error) {
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
	if NhcGovEnableHttpProxy {
		client = NhcGovSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return queryNhcGovListResponseDataResults, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", NhcGovCookie)
	req.Header.Set("Host", "www.nhc.gov.cn")
	req.Header.Set("Referer", "https://www.nhc.gov.cn/wjw/wsbzxx/wsbz.shtml")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryNhcGovListResponseDataResults, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryNhcGovListResponseDataResults, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryNhcGovListResponseDataResults, err
	}
	queryNhcGovListResponse := &QueryNhcGovListResponse{}
	err = json.Unmarshal(respBytes, queryNhcGovListResponse)
	if err != nil {
		return queryNhcGovListResponseDataResults, err
	}
	queryNhcGovListResponseDataResults = queryNhcGovListResponse.Data.Results
	return queryNhcGovListResponseDataResults, nil
}

func NhcGovDetailDoc(requestUrl string, referer string) (doc *html.Node, err error) {
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
	if NhcGovEnableHttpProxy {
		client = NhcGovSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", NhcGovCookie)
	req.Header.Set("Referer", referer)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Mobile Safari/537.36")
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

func downloadNhcGov(attachmentUrl string, referer string, filePath string) error {
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
	if NhcGovEnableHttpProxy {
		client = NhcGovSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.nhc.gov.cn")
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
