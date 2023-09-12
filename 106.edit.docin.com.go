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

const (
	EditDocInEnableHttpProxy = false
	EditDocInHttpProxyUrl    = "111.225.152.186:8089"
)

func EditDocInSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(EditDocInHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

var DocInCookie = "docin_session_id=fd912190-3bfb-46b6-874e-e80e03a13e6a; cookie_id=CA075B2FFC700001BB661788D6701A55; time_id=2022102510372; partner_tips=1; last_upload_public243402665=yes; __bid_n=18504630ab46cccb5c4207; FEID=v10-b58297500cfb5dd9acd684d3a79fe0046c214fae; __xaf_fpstarttimer__=1672884568247; __xaf_thstime__=1672884568292; __xaf_fptokentimer__=1672884568303; page_length=100; visitTopicIds=285258; FPTOKEN=kvb7/hh60j7cFgDOQawZiYvvkPY/qB48fxed8C/fl+P2R42S57Os2TYTN+VR6NmHcA1nWF/Zzlu+7TY7/HtU73imcFBpnSP39G/w34CSqNw2RManN9ZWM8Rd53fSy4NLTbovm4+LUM81V3W7SHzatAkSv+h00+noHqnIfLZikrRsHYCNvzkQCjEQ/dCCZbYyRhByiFLA/Z8H5yonq8PxCwHsklqAaHvp6Qpf2rRUMG+XO/yBsrGI/S35CJvnLiqAwS7rmeC81UJCMTKO8QV55/E7b3uTfd3ag7iJetYaONGBKi1RHyqebvtjkchSNruW7MYQ/XA9CYUKNQoU1w4YRgjvCWsaLWYjgjQKwEb8kOGIkQOv0tIJlrNL4gB7tc2+FNjHOlN4Bv4Psr3hA+gnLA==|WAHGPrsYPHcHXFUCs4FqFpoq7UHi1svYa41yPZFyRyw=|10|72d0d163d5eb20094432489b1e57cea9; aliyungf_tc=b81c45404af697ae66c345ab4e4631abd0abfb134ababa4519e41f5fd12fb63d; login_email=15238369929; user_password=bGKLZV42BM55cslJQjQYVqpq7lLi5mzCsPdXYATSgeY%3DH_TQFZ9qedzOno0nj9FG3FRIEIH8Vz6PIWcAGFbDGqSIhC5uhBCK3LyhR7xS2mggnjj; tipMaxId=0; ifShowMsg=true; jumpIn=401; isbaiduspider=false; recharge_from_type=nav-sub1; netfunction=\"/my/upload/myUpload.do\"; mobilefirsttip=tip; today_first_in=1; _gid=GA1.2.1089910833.1694394686; s_from=direct; uaType=chrome; JSESSIONID=BBA18BAF2FE03930D8AD3C4E0349F2F5-n1; remindClickId=-1; _gat_gtag_UA_3158355_1=1; _gat=1; _ga_34B604LFFQ=GS1.1.1694397751.127.1.1694399062.57.0.0; _ga=GA1.1.1839309374.1666665423; _ga_ZYR13KTSXC=GS1.1.1694397751.66.1.1694399086.33.0.0"
var downPrice = 5

// ychEduSpider 编辑豆丁文档
// @Title 编辑豆丁文档
// @Description https://www.docin.com/，编辑豆丁文档
func main() {
	currentPage := 1
	beginId := 0
	for {
		pageListUrl := "https://www.docin.com/my/upload/myUpload.do?onlypPublic=1&totalpublicnum=0"
		referer := "https://www.docin.com/my/upload/myUpload.do?onlypPublic=1&totalpublicnum=0"
		if currentPage > 1 {
			pageListUrl = fmt.Sprintf("https://www.docin.com/my/upload/myUpload.do?styleList=1"+
				"&orderName=0&orderDate=0&orderVisit=0&orderStatus=0&orderFolder=0&folderId=0"+
				"&myKeyword=&publishCount=&onlypPrivate=&totalprivatenum=0&onlypPublic=1"+
				"&totalpublicnum=0&currentPage=%d&pageType=n&beginId=%d", currentPage, beginId)
		}

		fmt.Println(pageListUrl)
		pageListDoc, err := QueryDocInDoc(pageListUrl, referer)
		if err != nil {
			fmt.Println(err)
			break
		}
		tbodyNodes := htmlquery.Find(pageListDoc, `//div[@class="tableWarp"]/table[@class="my-data"]/tbody`)
		if len(tbodyNodes) <= 0 {
			break
		}
		idsArr := make([]string, 0)
		for _, tbodyNode := range tbodyNodes {
			trNode := htmlquery.FindOne(tbodyNode, `./tr`)
			trId := strings.ReplaceAll(htmlquery.SelectAttr(trNode, "id"), "tr", "")
			idsArr = append(idsArr, trId)

			filePageNode := htmlquery.FindOne(tbodyNode, `./tr/td[4]`)
			filePage := htmlquery.InnerText(filePageNode)
			filePage = strings.TrimSpace(filePage)
			filePage = strings.ReplaceAll(filePage, "页", "")
			// 根据页数设置价格
			filePageNum, _ := strconv.Atoi(filePage)
			if filePageNum > 0 {
				if filePageNum > 0 && filePageNum <= 8 {
					downPrice = 2
				} else if filePageNum > 8 && filePageNum <= 18 {
					downPrice = 3
				} else if filePageNum > 18 && filePageNum <= 28 {
					downPrice = 4
				} else if filePageNum > 28 && filePageNum <= 38 {
					downPrice = 5
				} else if filePageNum > 38 && filePageNum <= 48 {
					downPrice = 6
				} else if filePageNum > 48 && filePageNum <= 58 {
					downPrice = 7
				} else {
					downPrice = 8
				}
			}
			fmt.Println("-----------------开始设置价格--------------------")

			fileTitleNode := htmlquery.FindOne(tbodyNode, `./tr/td[2]/a`)
			fileTitle := htmlquery.SelectAttr(fileTitleNode, "title")
			fmt.Println(fileTitle)

			// 开始设置价格
			editUrl := fmt.Sprintf("https://www.docin.com/app/my/docin/batchModifyPrice.do?ids=%s&down_price=%d&price_flag=1", trId, downPrice)
			_, err := QueryDocInDoc(editUrl, referer)
			if err != nil {
				fmt.Println(err)
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
