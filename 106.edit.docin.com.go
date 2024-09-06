package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var EditDocInEnableHttpProxy = false
var EditDocInHttpProxyUrl = ""
var EditDocInHttpProxyUrlArr = make([]string, 0)

func EditDocInHttpProxy() error {
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
					EditDocInHttpProxyUrlArr = append(EditDocInHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					EditDocInHttpProxyUrlArr = append(EditDocInHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func EditDocInSetHttpProxy() (httpclient *http.Client) {
	if EditDocInHttpProxyUrl == "" {
		if len(EditDocInHttpProxyUrlArr) <= 0 {
			err := EditDocInHttpProxy()
			if err != nil {
				EditDocInSetHttpProxy()
			}
		}
		EditDocInHttpProxyUrl = EditDocInHttpProxyUrlArr[0]
		if len(EditDocInHttpProxyUrlArr) >= 2 {
			EditDocInHttpProxyUrlArr = EditDocInHttpProxyUrlArr[1:]
		} else {
			EditDocInHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(EditDocInHttpProxyUrl)
	ProxyURL, _ := url.Parse(EditDocInHttpProxyUrl)
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

var DocInCookie = "docin_session_id=b45ae036-1776-4855-a61f-d3441addbafb; pbyCookieKey=1719393050260; cookie_id=CACBC82BFA8000012A8D1A82E3081CC5; time_id=2024626171050; partner_tips=1; last_upload_public243402665=yes; tipMaxId=574386458; userchoose=170_169_174_171_175_176_177_180_181_998_000_; userChoose=usertags169_174_171_175_176_177_180_181_998_000_; ifShowMsg=true; _gid=GA1.2.203874693.1724818626; mobilefirsttip=tip; remindClickId=-1; login_email=15238369929; user_password=QQ7VSpIhJ4%2FLb8%2FJFb0FLKpq7lLi5mzCsPdXYATSgeY%3DH_TAOWfW55issLykbgUyREQsa3gKt8jed5YAGFbDGqSIhC5uhBCK3LyhR7xS2mggnjj; refererusertype=0; s_from=direct; uaType=chrome; netfunction=\"/my/upload/myUpload.do\"; today_first_in=1; JSESSIONID=C4B8D4541FB376B3CA3B05EA20BADD4D-n2; _ga_ZYR13KTSXC=GS1.1.1724917862.75.1.1724917918.4.0.0; _ga=GA1.2.1198245365.1719393051"
var downPrice = 5

// ychEduSpider 编辑豆丁文档
// @Title 编辑豆丁文档
// @Description https://www.docin.com/，编辑豆丁文档
func main() {
	currentPage := 1
	beginId := 0
	for {
		pageListUrl := "https://www.docin.com/my/upload/myUpload.do?onlypPublic=1&totalpublicnum=0"
		referer := "https://www.docin.com/my/upload/myUpload.do?styleList=1&orderName=0&orderDate=0&orderVisit=0&orderStatus=0&folderId=-1"
		if currentPage >= 1 {
			pageListUrl = fmt.Sprintf("https://www.docin.com/my/upload/myUpload.do?styleList=1"+
				"&orderName=0&orderDate=0&orderVisit=0&orderStatus=0&orderFolder=0"+
				"&folderId=0&myKeyword=&publishCount=&onlypPrivate=&totalprivatenum=0"+
				"&onlypPublic=&totalpublicnum=0&currentPage=%d"+
				"&pageType=n&beginId=%d", currentPage, beginId)
		}

		pageListDoc, err := QueryDocInDoc(pageListUrl, referer)
		if err != nil {
			fmt.Println(err)
			EditDocInHttpProxyUrl = ""
			continue
		}
		tbodyNodes := htmlquery.Find(pageListDoc, `//div[@class="tableWarp"]/table[@class="my-data"]/tbody`)
		if len(tbodyNodes) <= 0 {
			break
		}
		idsArr := make([]string, 0)
		for _, tbodyNode := range tbodyNodes {
			trNode := htmlquery.FindOne(tbodyNode, `./tr`)
			trId := strings.ReplaceAll(htmlquery.SelectAttr(trNode, "id"), "tr", "")

			//viewUrl := fmt.Sprintf("https://www.docin.com/p-%s.html", trId)
			//fmt.Println("访问文档详情")
			//_, err = ViewDocInDoc(viewUrl, referer)
			//if err != nil {
			//	fmt.Println(err)
			//	EditDocInHttpProxyUrl = ""
			//	continue
			//}
			//_, err := htmlquery.LoadURL(viewUrl)
			//if err != nil {
			//	continue
			//}

			idsArr = append(idsArr, trId)

			fileTitleNode := htmlquery.FindOne(tbodyNode, `./tr/td[2]/a`)
			fileTitle := htmlquery.SelectAttr(fileTitleNode, "title")
			fmt.Println(fileTitle)

			filePageNode := htmlquery.FindOne(tbodyNode, `./tr/td[4]`)
			filePage := htmlquery.InnerText(filePageNode)
			filePage = strings.TrimSpace(filePage)
			filePage = strings.ReplaceAll(filePage, "页", "")
			// 根据页数设置价格
			filePageNum, _ := strconv.Atoi(filePage)
			if filePageNum > 0 {
				if filePageNum > 0 && filePageNum <= 5 {
					downPrice = 2
				} else if filePageNum > 5 && filePageNum <= 10 {
					downPrice = 3
				} else if filePageNum > 10 && filePageNum <= 15 {
					downPrice = 4
				} else if filePageNum > 15 && filePageNum <= 20 {
					downPrice = 5
				} else if filePageNum > 20 && filePageNum <= 25 {
					downPrice = 6
				} else if filePageNum > 25 && filePageNum <= 30 {
					downPrice = 7
				} else if filePageNum > 30 && filePageNum <= 35 {
					downPrice = 8
				} else if filePageNum > 35 && filePageNum <= 40 {
					downPrice = 9
				} else if filePageNum > 40 && filePageNum <= 45 {
					downPrice = 10
				} else if filePageNum > 45 && filePageNum <= 50 {
					downPrice = 11
				} else {
					downPrice = 12
				}
			}

			// 查看文档原来价格
			filePriceNode := htmlquery.FindOne(tbodyNode, `./tr/td[5]`)
			filePrice := htmlquery.InnerText(filePriceNode)
			filePrice = strings.TrimSpace(filePrice)
			if filePrice != "免费" {
				floatFilePrice, err := strconv.ParseFloat(filePrice, 64)
				if err != nil {
					continue
				}
				originalPrice := int(floatFilePrice)
				if downPrice == originalPrice {
					continue
				}
			}

			// 开始设置价格
			fmt.Println("-----------------开始设置价格--------------------")
			editUrl := fmt.Sprintf("https://www.docin.com/app/my/docin/batchModifyPrice.do?ids=%s&down_price=%d&price_flag=0", trId, downPrice)
			_, err = QueryDocInDoc(editUrl, referer)
			if err != nil {
				fmt.Println(err)
				EditDocInHttpProxyUrl = ""
				continue
			}
			fmt.Println("-----------------开始设置价格完结--------------------")
			time.Sleep(time.Microsecond * 100)
		}
		beginId, _ = strconv.Atoi(idsArr[len(idsArr)-1])
		currentPage++
		fmt.Println(currentPage)
		referer = fmt.Sprintf("https://www.docin.com/my/upload/myUpload.do?styleList=1"+
			"&orderName=0&orderDate=0&orderVisit=0&orderStatus=0&orderFolder=0&folderId=0"+
			"&myKeyword=&publishCount=&onlypPrivate=&totalprivatenum=0&onlypPublic=1"+
			"&totalpublicnum=0&currentPage=%d&pageType=n&beginId=%d", currentPage, beginId)
	}
}

func ViewDocInDoc(requestUrl string, referer string) (doc *html.Node, err error) {
	// 初始化客户端
	var client *http.Client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*20)
				if err != nil {
					fmt.Println("dail timeout", err)
					return nil, err
				}
				return c, nil

			},
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 20,
		},
	}
	if EditDocInEnableHttpProxy {
		client = EditDocInSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.docin.com")
	req.Header.Set("Origin", "https://www.docin.com")
	req.Header.Set("Referer", referer)
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

func QueryDocInDoc(requestUrl string, referer string) (doc *html.Node, err error) {
	// 初始化客户端
	var client *http.Client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*20)
				if err != nil {
					fmt.Println("dail timeout", err)
					return nil, err
				}
				return c, nil

			},
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 20,
		},
	}
	if EditDocInEnableHttpProxy {
		client = EditDocInSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", DocInCookie)
	req.Header.Set("Host", "www.docin.com")
	req.Header.Set("Origin", "https://www.docin.com")
	req.Header.Set("Referer", referer)
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
