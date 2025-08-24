package main

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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

var CcSnEnableHttpProxy = false
var CcSnHttpProxyUrl = "111.225.152.186:8089"
var CcSnHttpProxyUrlArr = make([]string, 0)

func CcSnHttpProxy() error {
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
					CcSnHttpProxyUrlArr = append(CcSnHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					CcSnHttpProxyUrlArr = append(CcSnHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func CcSnSetHttpProxy() (httpclient *http.Client) {
	if CcSnHttpProxyUrl == "" {
		if len(CcSnHttpProxyUrlArr) <= 0 {
			err := CcSnHttpProxy()
			if err != nil {
				CcSnSetHttpProxy()
			}
		}
		CcSnHttpProxyUrl = CcSnHttpProxyUrlArr[0]
		if len(CcSnHttpProxyUrlArr) >= 2 {
			CcSnHttpProxyUrlArr = CcSnHttpProxyUrlArr[1:]
		} else {
			CcSnHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(CcSnHttpProxyUrl)
	ProxyURL, _ := url.Parse(CcSnHttpProxyUrl)
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

type QueryCcSnListFormData struct {
	__EVENTTARGET                    string
	__EVENTARGUMENT                  string
	__VIEWSTATE                      string
	__VIEWSTATEGENERATOR             string
	ID_ucZbbzList_txtKeyWord         string
	ID_ucZbbzList_ucPager1_listPage  int
	ID_ucZbbzList_ucPager1_btnPaging string
}

var CcSnCookie = "ASP.NET_SessionId=3gmamifoetqt0d55qepxrd55"

// 下载工程建设标准文档
// @Title 下载工程建设标准文档
// @Description https://www.ccsn.org.cn/，下载工程建设标准文档
func main() {
	pageListUrl := "https://www.ccsn.org.cn/Zbbz/ZbbzList.aspx"
	fmt.Println(pageListUrl)
	page := 1
	maxPage := 303
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		queryCcSnListFormData := QueryCcSnListFormData{
			__EVENTTARGET:                    "",
			__EVENTARGUMENT:                  "",
			__VIEWSTATE:                      "/wEPDwULLTEyOTc3OTk1MzUPFiAeDF9DdXJyZW50VmlldwUKdWNaYmJ6TGlzdB4hYWN0Rm9yY2VTdGFuZGFyZExpc3RfVHJhZGVvclBsYWNlZB4hYWN0Rm9yY2VTdGFuZGFyZExpc3RfU3RhbmRhcmRUeXBlBQblhajpg6geIWFjdEZvcmNlU3RhbmRhcmRMaXN0X3N0YW5kYXJkQ29kZWQeFElEX3VjWmJiekxpc3RPcmRlckJ5ZB4cQWN0QmFzZV9DbGllbnRSZXR1cm5Gb3JFcnJvcmQeEEFjdEJhc2VfUmV0dXJuVUNkHhRJRF91Y1piYnpMaXN0S2V5V29yZGQeIWFjdEZvcmNlU3RhbmRhcmRMaXN0X3N0YW5kYXJkTmFtZWQeG2FjdEZvcmNlU3RhbmRhcmRMaXN0X0lzRnpCemQeImFjdEZvcmNlU3RhbmRhcmRMaXN0X3N0YW5kYXJkQ2xhc3NkHhZJRF91Y1piYnpMaXN0UGFnZUluZGV4AgEeFklEX3VjWmJiekxpc3RQYWdlQ291bnQCrwIeCV9OZXh0Vmlld2QeEF9QcmVwYXJhdGl2ZVZpZXdkHiFhY3RGb3JjZVN0YW5kYXJkTGlzdF9Jc0Nob29zVHJhZGVkFgICAw9kFgICAQ9kFgJmD2QWBmYPZBYCAgEPFgIeBWNsYXNzBQZhY3RpdmVkAgcPFgIeC18hSXRlbUNvdW50AgoWFGYPZBYMAgEPDxYCHgRUZXh0BQIxMWRkAgMPDxYEHxIFHuKFouOAgeKFo+e6p+mTgei3r+iuvuiuoeinhOiMgx4LTmF2aWdhdGVVcmwFNi4uL1Nob3cuYXNweD9HdWlkPTg5NjY2OTNhLTQ2MDEtNDAxMy1hMzRhLTcwMmY3MDFiMjZiMWRkAgUPDxYCHxIFKkNvZGUgZm9yIGRlc2lnbiBvZiBjbGFzcyDihaLjgIHihaMgcmFpbHdheWRkAgcPDxYCHxIFDEdCNTAwMTItMjAxMmRkAgkPDxYCHxIFCjIwMTIvMTAvMTFkZAILDw8WAh8SBQkyMDEyLzEyLzFkZAIBD2QWDAIBDw8WAh8SBQIxMmRkAgMPDxYEHxIFGOWupOWklue7meawtOiuvuiuoeagh+WHhh8TBTYuLi9TaG93LmFzcHg/R3VpZD1mZjU3ZTFlZC1jNTk1LTQwNzEtOGE1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MTU3AzE1OAMxNTkDMTYwAzE2MQMxNjIDMTYzAzE2NAMxNjUDMTY2AzE2NwMxNjgDMTY5AzE3MAMxNzEDMTcyAzE3MwMxNzQDMTc1AzE3NgMxNzcDMTc4AzE3OQMxODADMTgxAzE4MgMxODMDMTg0AzE4NQMxODYDMTg3AzE4OAMxODkDMTkwAzE5MQMxOTIDMTkzAzE5NAMxOTUDMTk2AzE5NwMxOTgDMTk5AzIwMAMyMDEDMjAyAzIwMwMyMDQDMjA1AzIwNgMyMDcDMjA4AzIwOQMyMTADMjExAzIxMgMyMTMDMjE0AzIxNQMyMTYDMjE3AzIxOAMyMTkDMjIwAzIyMQMyMjIDMjIzAzIyNAMyMjUDMjI2AzIyNwMyMjgDMjI5AzIzMAMyMzEDMjMyAzIzMwMyMzQDMjM1AzIzNgMyMzcDMjM4AzIzOQMyNDADMjQxAzI0MgMyNDMDMjQ0AzI0NQMyNDYDMjQ3AzI0OAMyNDkDMjUwAzI1MQMyNTIDMjUzAzI1NAMyNTUDMjU2AzI1NwMyNTgDMjU5AzI2MAMyNjEDMjYyAzI2MwMyNjQDMjY1AzI2NgMyNjcDMjY4AzI2OQMyNzADMjcxAzI3MgMyNzMDMjc0AzI3NQMyNzYDMjc3AzI3OAMyNzkDMjgwAzI4MQMyODIDMjgzAzI4NAMyODUDMjg2AzI4NwMyODgDMjg5AzI5MAMyOTEDMjkyAzI5MwMyOTQDMjk1AzI5NgMyOTcDMjk4AzI5OQMzMDADMzAxAzMwMgMzMDMVrwIBMQEyATMBNAE1ATYBNwE4ATkCMTACMTECMTICMTMCMTQCMTUCMTYCMTcCMTgCMTkCMjACMjECMjICMjMCMjQCMjUCMjYCMjcCMjgCMjkCMzACMzECMzICMzMCMzQCMzUCMzYCMzcCMzgCMzkCNDACNDECNDICNDMCNDQCNDUCNDYCNDcCNDgCNDkCNTACNTECNTICNTMCNTQCNTUCNTYCNTcCNTgCNTkCNjACNjECNjICNjMCNjQCNjUCNjYCNjcCNjgCNjkCNzACNzECNzICNzMCNzQCNzUCNzYCNzcCNzgCNzkCODACODECODICODMCODQCODUCODYCODcCODgCODkCOTACOTECOTICOTMCOTQCOTUCOTYCOTcCOTgCOTkDMTAwAzEwMQMxMDIDMTAzAzEwNAMxMDUDMTA2AzEwNwMxMDgDMTA5AzExMAMxMTEDMTEyAzExMwMxMTQDMTE1AzExNgMxMTcDMTE4AzExOQMxMjADMTIxAzEyMgMxMjMDMTI0AzEyNQMxMjYDMTI3AzEyOAMxMjkDMTMwAzEzMQMxMzIDMTMzAzEzNAMxMzUDMTM2AzEzNwMxMzgDMTM5AzE0MAMxNDEDMTQyAzE0MwMxNDQDMTQ1AzE0NgMxNDcDMTQ4AzE0OQMxNTADMTUxAzE1MgMxNTMDMTU0AzE1NQMxNTYDMTU3AzE1OAMxNTkDMTYwAzE2MQMxNjIDMTYzAzE2NAMxNjUDMTY2AzE2NwMxNjgDMTY5AzE3MAMxNzEDMTcyAzE3MwMxNzQDMTc1AzE3NgMxNzcDMTc4AzE3OQMxODADMTgxAzE4MgMxODMDMTg0AzE4NQMxODYDMTg3AzE4OAMxODkDMTkwAzE5MQMxOTIDMTkzAzE5NAMxOTUDMTk2AzE5NwMxOTgDMTk5AzIwMAMyMDEDMjAyAzIwMwMyMDQDMjA1AzIwNgMyMDcDMjA4AzIwOQMyMTADMjExAzIxMgMyMTMDMjE0AzIxNQMyMTYDMjE3AzIxOAMyMTkDMjIwAzIyMQMyMjIDMjIzAzIyNAMyMjUDMjI2AzIyNwMyMjgDMjI5AzIzMAMyMzEDMjMyAzIzMwMyMzQDMjM1AzIzNgMyMzcDMjM4AzIzOQMyNDADMjQxAzI0MgMyNDMDMjQ0AzI0NQMyNDYDMjQ3AzI0OAMyNDkDMjUwAzI1MQMyNTIDMjUzAzI1NAMyNTUDMjU2AzI1NwMyNTgDMjU5AzI2MAMyNjEDMjYyAzI2MwMyNjQDMjY1AzI2NgMyNjcDMjY4AzI2OQMyNzADMjcxAzI3MgMyNzMDMjc0AzI3NQMyNzYDMjc3AzI3OAMyNzkDMjgwAzI4MQMyODIDMjgzAzI4NAMyODUDMjg2AzI4NwMyODgDMjg5AzI5MAMyOTEDMjkyAzI5MwMyOTQDMjk1AzI5NgMyOTcDMjk4AzI5OQMzMDADMzAxAzMwMgMzMDMUKwOvAmdnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2RkZInE5Ycz6TSuVgyBLIiwhmqWC1Iy",
			__VIEWSTATEGENERATOR:             "CD78B523",
			ID_ucZbbzList_txtKeyWord:         "",
			ID_ucZbbzList_ucPager1_listPage:  page,
			ID_ucZbbzList_ucPager1_btnPaging: "GO",
		}
		listDoc, err := QueryCcSnList(pageListUrl, queryCcSnListFormData)
		if err != nil {
			fmt.Println(err)
			isPageListGo = false
			continue
		}
		// /html/body/form/div[3]/div[2]/div[2]/div/div[2]/table/tbody/tr[2]
		trNodes := htmlquery.Find(listDoc, `//form[@id="form1"]/div[@id="container"]/div[@class="section mobile-container-padding"]/div[@class="interaction-wrapper clear"]/div[@class="mobile-m-0"]/div[@class="smLinkList-wrapper mt14 mobile-mt-20"]/table/tbody/tr[2]/td/table/tbody/tr`)
		if len(trNodes) >= 1 {
			for _, trNode := range trNodes {
				fmt.Println("=====================开始处理数据 page = ", page, "=========================")
				titleNode := htmlquery.FindOne(trNode, `./td[2]/a`)
				title := htmlquery.InnerText(titleNode)
				title = strings.TrimSpace(title)
				title = strings.ReplaceAll(title, "-", "")
				title = strings.ReplaceAll(title, " ", "")
				title = strings.ReplaceAll(title, "|", "-")
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, "\n", "")
				title = strings.ReplaceAll(title, "\r", "")
				title = strings.ReplaceAll(title, " ", "")
				fmt.Println(title)

				codeNode := htmlquery.FindOne(trNode, `./td[3]/span`)
				code := htmlquery.InnerText(codeNode)
				code = strings.TrimSpace(code)
				code = strings.ReplaceAll(code, "/", "-")
				fmt.Println(code)

				filePath := "../www.ccsn.org.cn/" + title + "(" + code + ")" + ".pdf"
				fmt.Println(filePath)

				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				fmt.Println("=======开始下载========")

				gUidNode := htmlquery.FindOne(trNode, `./td[2]/a/@href`)
				gUidText := htmlquery.InnerText(gUidNode)
				gUidArray := strings.Split(gUidText, "Guid=")
				gUid := gUidArray[1]

				detailUrl := fmt.Sprintf("https://www.ccsn.org.cn/Zbbz/Show.aspx?Guid=%s", gUid)
				fmt.Println(detailUrl)

				showFullTextUrl := fmt.Sprintf("https://www.ccsn.org.cn/Zbbz/ShowFullText.aspx?Guid=%s", gUid)
				fmt.Println(showFullTextUrl)

				showFullTextDoc, err := QueryCcSnHtml(showFullTextUrl, detailUrl)
				//fmt.Println(htmlquery.OutputHTML(showFullTextDoc, true))
				//os.Exit(1)
				if err != nil {
					fmt.Println(err)
					continue
				}

				downloadUrlNode := htmlquery.FindOne(showFullTextDoc, `//div[@id="ID_ucShowFullText_div_media"]/a[@class="media"]/@href`)
				if downloadUrlNode == nil {
					fmt.Println("沒有下載鏈接")
					continue
				}

				downloadUrl := htmlquery.InnerText(downloadUrlNode)
				fmt.Println(downloadUrl)

				fmt.Println("=======开始下载" + title + "========")
				err = downloadCcSn(downloadUrl, detailUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "../www.ccsn.org.cn", "../upload.doc88.com/hbba.sacinfo.org.cn")
				err = copyCcSnFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")
				//DownLoadCcSnTimeSleep := 10
				DownLoadCcSnTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadCcSnTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadCcSnTimeSleep, "秒 倒计时", i, "秒===========")
				}
			}
			DownLoadCcSnPageTimeSleep := 10
			// DownLoadCcSnPageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadCcSnPageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadCcSnPageTimeSleep, "秒 倒计时", i, "秒===========")
			}
			page++
			if page > maxPage {
				isPageListGo = false
				break
			}
		}
	}
}

func QueryCcSnList(requestUrl string, queryCcSnListFormData QueryCcSnListFormData) (doc *html.Node, err error) {
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
	if CcSnEnableHttpProxy {
		client = CcSnSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("__EVENTTARGET", queryCcSnListFormData.__EVENTTARGET)
	postData.Add("__EVENTARGUMENT", queryCcSnListFormData.__EVENTARGUMENT)
	postData.Add("__VIEWSTATE", queryCcSnListFormData.__VIEWSTATE)
	postData.Add("__VIEWSTATEGENERATOR", queryCcSnListFormData.__VIEWSTATEGENERATOR)
	postData.Add("ID_ucZbbzList$txtKeyWord", queryCcSnListFormData.ID_ucZbbzList_txtKeyWord)
	postData.Add("ID_ucZbbzList$ucPager1$listPage", strconv.Itoa(queryCcSnListFormData.ID_ucZbbzList_ucPager1_listPage))
	postData.Add("ID_ucZbbzList$ucPager1$btnPaging", queryCcSnListFormData.ID_ucZbbzList_ucPager1_btnPaging)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", CcSnCookie)
	req.Header.Set("Host", "cx.www.ccsn.org.cn")
	req.Header.Set("Origin", "https://www.ccsn.org.cn")
	req.Header.Set("Referer", "https://www.ccsn.org.cn/?serviceName=bzh_api_standard")
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
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return doc, err
	}
	doc, err = decodeAndParseHTMLCcSn(string(bodyBytes))
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func decodeAndParseHTMLCcSn(gb2312Content string) (*html.Node, error) {
	// 使用GB2312解码器解码内容
	decoder := simplifiedchinese.GBK.NewDecoder() // 注意：通常GB2312在Go中对应的是GBK，而非直接使用GB2312，因为GB2312不是一个广泛支持的编码标准，而是GBK的一个子集。
	decodedContent, _, err := transform.Bytes(decoder, []byte(gb2312Content))
	if err != nil {
		return nil, err
	}
	// 将解码后的内容转换为UTF-8（通常HTML解析器需要UTF-8编码）
	utf8Content := decodedContent
	// 解析HTML
	doc, err := html.Parse(bytes.NewReader(utf8Content))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func QueryCcSnHtml(requestUrl string, referer string) (doc *html.Node, err error) {
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
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", CcSnCookie)
	req.Header.Set("Host", "www.ccsn.org.cn")
	req.Header.Set("Origin", "https://www.ccsn.org.cn")
	req.Header.Set("Referer", referer)
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
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return doc, err
	}
	doc, err = decodeAndParseHTMLCcSn(string(bodyBytes))
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func downloadCcSn(attachmentUrl string, referer string, filePath string) error {
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
	if CcSnEnableHttpProxy {
		client = CcSnSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "*/*")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", CcSnCookie)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.ccsn.org.cn")
	req.Header.Set("Referer", referer)
	//req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	//req.Header.Set("sec-ch-ua-mobile", "?0")
	//req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	//req.Header.Set("Sec-Fetch-Dest", "document")
	//req.Header.Set("Sec-Fetch-Mode", "navigate")
	//req.Header.Set("Sec-Fetch-Site", "none")
	//req.Header.Set("Sec-Fetch-User", "?1")
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

func copyCcSnFile(src, dst string) (err error) {
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
