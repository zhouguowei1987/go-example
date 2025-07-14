package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var SpPtCfSaEnableHttpProxy = false
var SpPtCfSaHttpProxyUrl = "111.225.152.186:8089"
var SpPtCfSaHttpProxyUrlArr = make([]string, 0)

func SpPtCfSaHttpProxy() error {
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
					SpPtCfSaHttpProxyUrlArr = append(SpPtCfSaHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					SpPtCfSaHttpProxyUrlArr = append(SpPtCfSaHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func SpPtCfSaSetHttpProxy() (httpclient *http.Client) {
	if SpPtCfSaHttpProxyUrl == "" {
		if len(SpPtCfSaHttpProxyUrlArr) <= 0 {
			err := SpPtCfSaHttpProxy()
			if err != nil {
				SpPtCfSaSetHttpProxy()
			}
		}
		SpPtCfSaHttpProxyUrl = SpPtCfSaHttpProxyUrlArr[0]
		if len(SpPtCfSaHttpProxyUrlArr) >= 2 {
			SpPtCfSaHttpProxyUrlArr = SpPtCfSaHttpProxyUrlArr[1:]
		} else {
			SpPtCfSaHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(SpPtCfSaHttpProxyUrl)
	ProxyURL, _ := url.Parse(SpPtCfSaHttpProxyUrl)
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

type QuerySpPtCfSaListFormData struct {
	isLength  int
	num_tn int
	standard_type   string
	keyword       string
}

var SpPtCfSaCookie = "name=value; cookieName=cookieValue; JSESSIONID=1939457719D7E496EDAAB20F581FA15F"

// 下载食品安全国家标准数据文档
// @Title 下载食品安全国家标准数据文档
// @Description https://sppt.cfsa.net.cn:8086/db/，下载食品安全国家标准数据文档
func main() {
	pageListUrl := "https://sppt.cfsa.net.cn:8086/db?task=indexSearch"
    fmt.Println(pageListUrl)
    querySpPtCfSaListFormData := QuerySpPtCfSaListFormData{
        isLength:  9999,
        num_tn: 2,
        standard_type:    "",
        keyword:  "",
    }
    querySpPtCfSaListResponse, err := QuerySpPtCfSaList(pageListUrl, querySpPtCfSaListFormData)
    if err != nil {
        SpPtCfSaHttpProxyUrl = ""
        fmt.Println(err)
        continue
    }
    fmt.Println(querySpPtCfSaListResponse)
}

type QuerySpPtCfSaListFJ struct {
	FACT_NAME    string   `json:"FACT_NAME"`
	ID_F    string   `json:"ID_F"`
}

type QuerySpPtCfSaListResponse struct {
	CODE    string   `json:"CODE"`
	FJ    []QuerySpPtCfSaListFJ   `json:"FJ"`
	ID        string `json:"ID"`
	PDATE       string  `json:"PDATE"`
	SSRQ       string  `json:"SSRQ"`
	TABLENAME  string    `json:"TABLENAME"`
	TITLE string    `json:"TITLE"`
}

func QuerySpPtCfSaList(requestUrl string, querySpPtCfSaListFormData QuerySpPtCfSaListFormData) (querySpPtCfSaListResponse []QuerySpPtCfSaListResponse, err error) {
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
	if SpPtCfSaEnableHttpProxy {
		client = SpPtCfSaSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("isLength", strconv.Itoa(querySpPtCfSaListFormData.isLength))
	postData.Add("num_tn", querySpPtCfSaListFormData.num_tn)
	postData.Add("standard_type", querySpPtCfSaListFormData.standard_type)
	postData.Add("keyword", querySpPtCfSaListFormData.keyword)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	//req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Cookie", SpPtCfSaCookie)
	req.Header.Set("Host", "sppt.cfsa.net.cn:8086")
	req.Header.Set("Origin", "https://sppt.cfsa.net.cn:8086")
	req.Header.Set("Referer", "https://sppt.cfsa.net.cn:8086/db")
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
		fmt.Println(err)
		return querySpPtCfSaListResponse, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return querySpPtCfSaListResponse, err
	}
	err = json.Unmarshal(respBytes, querySpPtCfSaListResponse)
	if err != nil {
		fmt.Println(err)
		return querySpPtCfSaListResponse, err
	}
	return querySpPtCfSaListResponse, nil
}

func QuerySpPtCfSaDownLoadUrl(requestUrl string) (doc *html.Node, err error) {
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
	if SpPtCfSaEnableHttpProxy {
		client = SpPtCfSaSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", DownLoadUrlCookie)
	req.Header.Set("Host", "sppt.cfsa.net.cn:8086/db")
	req.Header.Set("Origin", "https://sppt.cfsa.net.cn:8086/db")
	req.Header.Set("Referer", "https://sppt.cfsa.net.cn:8086/db/uc/doc_manager.php?act=doc_list&state=all")
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

func DownLoadDoc88(attachmentUrl string, referer string, filePath string) error {
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
	if SpPtCfSaEnableHttpProxy {
		client = SpPtCfSaSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "sppt.cfsa.net.cn:8086/db")
	req.Header.Set("Origin", "https://sppt.cfsa.net.cn:8086/db")
	req.Header.Set("Referer", referer)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
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
