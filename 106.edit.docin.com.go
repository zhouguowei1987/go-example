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

var DocInCookie = "cookie_id=CA1BC7B9D32000011AEA1CB081C0DA70; time_id=20221227213222; partner_tips=1; __bid_n=18553c8c217b5d5e2e4207; FEID=v10-9bb248a4c21a53f72760ecda6234cbecf70a7381; __xaf_fpstarttimer__=1672147945190; __xaf_thstime__=1672147945346; __xaf_fptokentimer__=1672147945861; last_upload_public243402665=yes; indexnoticeupdatetag31526379=unshow; __root_domain_v=.docin.com; _qddaz=QD.151176617720177; _ga_ZYR13KTSXC=deleted; _ga_ZYR13KTSXC=deleted; FPTOKEN=yJ0+QSXweZ4Cqq0iuE4OoxgWMakB+ymq7HPZB8+AFbd3gLJfGXg+uhwC0PoTY2vuB9fF5/2qsuUHnMT2yBvsmKwBeb2Es9MY/cERDMA9Eo0Q4NnpM+qitBzZ5FgbqMdP5jaPkSY3peXMVlupZpjbAVsSoBjCQ5h+OaSsGsZHh5XMvNOM2sUd+BUqUvfHYTZcEf1zoMXQMbPlOFUOa1qcC5h6YlPH6Q3uN6f67bocIhZijom17xgVRl5ISjwvfhBdEIDZWVbcvcDY6+VU8WZfDOkd/p66Bm9Sz/OeodH8SiuMetE/mcgTgF5KiLFLh3yS8JVVgndmKGV3Yppl9eVrOxYrlQPZjf0rOfxGNdetjSbCEVs+HX/Usks8sfKMRey9ZSjQ0XKZdAhRoUoQwCk4BA==|FOHf58AIodtuneILuIhr9eSK/3MZMJ4ih8ikwzMvr5Y=|10|f132030d89062ab6a03a37cb0d2b63f1; pbyCookieKey=1692338365244; _gid=GA1.2.1560629737.1701515273; userchoose=104_174_171_102_175_176_169_177_178_179_180_181_998_000_; userChoose=usertags104_174_171_102_175_176_169_177_178_179_180_181_998_000_; tipMaxId=563451380; login_email=15238369929; user_password=KV6Pt8FSiB82Vxxi25QqE6pq7lLi5mzCsPdXYATSgeY%3DH_T4n74cojitZ4iWF2fxZGg7U0c79BT4bEwAGFbDGqSIhC5uhBCK3LyhR7xS2mggnjj; visitTopicIds=\"279237,285253,285264,284991,285122\"; mobilefirsttip=tip; today_first_in=1; aliyungf_tc=4a42abfa3ed4fa0922018ab2aa741da3315330cec5bf606e5f90ad248b88f29f; ifShowMsg=true; docin_session_id=b7443cb8-debe-44ce-b172-b1dbcf3099e7; s_from=direct; uaType=chrome; netfunction=\"/my/upload/myUpload.do\"; remindClickId=-1; jumpIn=400; saveFinProductToCookieValue=4571846455; JSESSIONID=8F3DA4005105E58975F7058E748DF75C-n2; _ga=GA1.2.43085923.1672147943; _gat_gtag_UA_3158355_1=1; _ga_ZYR13KTSXC=GS1.1.1703679016.773.1.1703679537.55.0.0"
var downPrice = 5

// ychEduSpider 编辑豆丁文档
// @Title 编辑豆丁文档
// @Description https://www.docin.com/，编辑豆丁文档
func main() {
	currentPage := 1
	beginId := 0
	for {
		pageListUrl := "https://www.docin.com/my/upload/myUpload.do?styleList=1&orderName=0&orderDate=0&orderVisit=0&orderStatus=0&folderId=-1"
		referer := "https://www.docin.com/my/upload/myUpload.do?styleList=1&orderName=0&orderDate=0&orderVisit=0&orderStatus=0&folderId=-1"
		if currentPage > 1 {
			pageListUrl = fmt.Sprintf("https://www.docin.com/my/upload/myUpload.do?styleList=1"+
				"&orderName=0&orderDate=0&orderVisit=0&orderStatus=0"+
				"&orderFolder=0&folderId=-1&myKeyword=&publishCount=&onlypPrivate="+
				"&totalprivatenum=0&onlypPublic=&totalpublicnum=0"+
				"&currentPage=%d&pageType=n&beginId=%d", currentPage, beginId)
		}

		fmt.Println(pageListUrl)
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
