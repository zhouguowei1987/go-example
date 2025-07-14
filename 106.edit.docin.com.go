package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
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

var DocInCookie = "last_upload_public243402665=yes; aliyungf_tc=b67c7843e29113946821791021990a9dfcb47204518972fa606e8fdfaf292f11; ifShowMsg=true; jumpIn=400; lastLoginType=weixin; doc_retrieval_welcome=1; docin_session_id=7ac77fa7-430d-43e9-8526-03bc15a2a7b1; cookie_id=CB06E60A30F00001A177E94618D0114F; time_id=2024122791421; recharge_from_type=nav-sub1; HMACCOUNT=00EDEFEA78E0441D; partner_tips=1; refererfunction=https%3A%2F%2Fwww.baidu.com%2Flink%3Furl%3DBV4-b9w4I4SW-59tL_I7aWpA9Xz-rkmA3VGPx49aXF5A2nyl6d13iDG3aOvc0Ap5Wcx1u407alaEpBIDTqO7dK3CRoTYaIaRGQtFHa-LoiW%26wd%3D%26eqid%3Df4d15be6012446b90000000667f9b823; isbaiduspider=false; Hm_lvt_6f08be44365dcdd8b6197b6770124977=1745975172; Hm_lpvt_6f08be44365dcdd8b6197b6770124977=1747032664; userchoose=170_169_174_171_175_176_177_180_181_998_000_; pbyCookieKey=1748749870334; userChoose=usertags170_169_174_171_175_176_177_180_181_998_000_; login_email=15238369929; user_password=uofLr1VXUFnZ3pLNVPDX5Kpq7lLi5mzCsPdXYATSgeY%3DH_T2mZLzKzFJFeQAZ89wgg2zIX1db6cKbXGAGFbDGqSIhC5uhBCK3LyhR7xS2mggnjj; mobilefirsttip=tip; today_first_in=1; _gid=GA1.2.91256854.1750381411; hide_home_study_banner_tips=1; s_from=direct; uaType=chrome; remindClickId=-1; saveFinProductToCookieValue=4882885303; downloadClickId=-1; booksaveClickId=-1; payReadClickId=-1; partnerLogin=-1; vip_alert_adv=-1; can_copy_alert=-1; payReadClickId_v2=-1; showFeekClickId=-1; addComdocs=-1; showShareClickId=-1; _ga=GA1.2.43085923.1672147943; _gat_gtag_UA_3158355_1=1; editOnlineloadClickId=-1; netfunction=/app/my/docin/editOne; JSESSIONID=17430A86C44A885FE6CC36415C3973E4-n2; _ga_ZYR13KTSXC=GS2.1.s1750425301$o1833$g1$t1750425470$j55$l0$h0"
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
		if currentPage >= 1 {
			pageListUrl = fmt.Sprintf("https://www.docin.com/my/upload/myUpload.do"+
				"?styleList=1&orderName=0&orderDate=0&orderVisit=0&orderStatus=0"+
				"&orderFolder=0&folderId=-1&myKeyword=&publishCount="+
				"&onlypPrivate=&totalprivatenum=0&onlypPublic=&"+
				"totalpublicnum=0&currentPage=%d&pageType=p&beginId=%d", currentPage, beginId)
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
					downPrice = 1
				} else if filePageNum > 5 && filePageNum <= 10 {
					downPrice = 2
				} else if filePageNum > 10 && filePageNum <= 15 {
					downPrice = 3
				} else if filePageNum > 15 && filePageNum <= 20 {
					downPrice = 4
				} else if filePageNum > 20 && filePageNum <= 25 {
					downPrice = 5
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
