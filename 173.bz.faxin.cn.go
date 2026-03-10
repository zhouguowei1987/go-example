package main

import (
	"errors"
	"fmt"
	"io"
	"golang.org/x/net/html"
// 	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
)

var FaXinEnableHttpProxy = false
var FaXinHttpProxyUrl = "111.225.152.186:8089"
var FaXinHttpProxyUrlArr = make([]string, 0)

func FaXinHttpProxy() error {
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
					FaXinHttpProxyUrlArr = append(FaXinHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					FaXinHttpProxyUrlArr = append(FaXinHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func FaXinSetHttpProxy() (httpclient *http.Client) {
	if FaXinHttpProxyUrl == "" {
		if len(FaXinHttpProxyUrlArr) <= 0 {
			err := FaXinHttpProxy()
			if err != nil {
				FaXinSetHttpProxy()
			}
		}
		FaXinHttpProxyUrl = FaXinHttpProxyUrlArr[0]
		if len(FaXinHttpProxyUrlArr) >= 2 {
			FaXinHttpProxyUrlArr = FaXinHttpProxyUrlArr[1:]
		} else {
			FaXinHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(FaXinHttpProxyUrl)
	ProxyURL, _ := url.Parse(FaXinHttpProxyUrl)
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

type QueryFaXinListFormData struct {
	sctype         string
	stdno          string
	sc             string
	stdname        string
	a404           string
	a825           string
	standStatus    string
	issueDateStart string
	issueDateEnd   string
	a205Start      string
	a205End        string
	issueDepart    string
	draftsDept     string
	reader         string
	ownerDept      string
	pageIndex      int
}
type QueryFaXinStdOnlineFormData struct {
	a100       string
	saleagtid  string
	username   string
	readdevice string
	encstr     string
}

var FaXinCookie = "JSESSIONID=221443C9BCC3F5ABFB08FEDAA4F80249; Hm_lvt_cb8a2025f4234726e55e45f893fb7954=1771997810; HMACCOUNT=4E5B3419A3141A8E; Hm_lvt_a317640b4aeca83b20c90d410335b70f=1771997824; HMACCOUNT=4E5B3419A3141A8E; Hm_lpvt_a317640b4aeca83b20c90d410335b70f=1773107689; lawapp_web=9B62EB060E96A41C08BFB4BFBB88DD5BF2BB3818921B383451F0D19BE409F11B4AE47A8408591BD48C8D7DDE66CC1A9A92C784A019D702252E95B404BEFF7E3D44970A51F90ECF18B63C050347503388DB057CC550BFC652175F8363CD3F4B178192F485E4C0C4ED85C3AB5DB15129F7518997575FB568B08088D2D9642173F328A30857EFCDA4CB7AF400784D51BFC8DEE90BB2; Hm_lpvt_cb8a2025f4234726e55e45f893fb7954=1773112420"
var FaXinEncStrCookie = "JSESSIONID=221443C9BCC3F5ABFB08FEDAA4F80249; Hm_lvt_cb8a2025f4234726e55e45f893fb7954=1771997810; HMACCOUNT=4E5B3419A3141A8E; Hm_lvt_a317640b4aeca83b20c90d410335b70f=1771997824; HMACCOUNT=4E5B3419A3141A8E; Hm_lpvt_a317640b4aeca83b20c90d410335b70f=1773107689; lawapp_web=9B62EB060E96A41C08BFB4BFBB88DD5BF2BB3818921B383451F0D19BE409F11B4AE47A8408591BD48C8D7DDE66CC1A9A92C784A019D702252E95B404BEFF7E3D44970A51F90ECF18B63C050347503388DB057CC550BFC652175F8363CD3F4B178192F485E4C0C4ED85C3AB5DB15129F7518997575FB568B08088D2D9642173F328A30857EFCDA4CB7AF400784D51BFC8DEE90BB2; JSESSIONID=2116BE012B79BB7E567CF8045F19A765; Hm_lpvt_cb8a2025f4234726e55e45f893fb7954=1773119144"
var SpcStdOnlineCookie = "JSESSIONID=73AED56B4D4427C678D2491D782ADA7C; Hm_lpvt_6d75523a84ebfd663067173dd3baab34=1773119308"
var SpcDownloadCookie = "Qs_lvt_503365=1759975608%2C1763093101%2C1766475585%2C1768268241%2C1769395188; Qs_pv_503365=2169124988614570200%2C1976145068643755300%2C2957571777807069000%2C3500804953458690000%2C4312190314712403000; Hm_lvt_6d75523a84ebfd663067173dd3baab34=1771997845; HMACCOUNT=4E5B3419A3141A8E; Hm_lvt_b5bdb81ba9543a3b3567778ab86df74b=1773053945; Hm_lvt_f1e3be80a525d0bb9e111ca6cbaa2457=1773053945; Hm_lpvt_b5bdb81ba9543a3b3567778ab86df74b=1773114537; Hm_lpvt_f1e3be80a525d0bb9e111ca6cbaa2457=1773114537; JSESSIONID=AEA0BFF5E999FDA4B9D5AEECF0588C52; Hm_lpvt_6d75523a84ebfd663067173dd3baab34=1773119230"

// 下载法信标准文档
// @Title 下载法信标准文档
// @Description https://bz.faxin.cn/，下载法信标准文档
func main() {
	pageListUrl := "https://bz.faxin.cn/faxin/view/advancedsearch"
	fmt.Println(pageListUrl)
	page := 0
	maxPage := 30 //对不起，目前网站只支持下翻30页
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		queryFaXinListFormData := QueryFaXinListFormData{
			sctype:         "",
			stdno:          "",
			sc:             "CN", //CN：国家标准 QT：行业标准 JJ：计量规程规范
			stdname:        "",
			a404:           "",
			a825:           "*",
			standStatus:    "",
			issueDateStart: "",
			issueDateEnd:   "",
			a205Start:      "",
			a205End:        "",
			issueDepart:    "",
			draftsDept:     "",
			reader:         "",
			ownerDept:      "",
			pageIndex:      page,
		}
		queryFaXinListDoc, err := QueryFaXinList(pageListUrl, queryFaXinListFormData)
		if err != nil {
			fmt.Println(err)
			break
		}
		divNodes := htmlquery.Find(queryFaXinListDoc, `//html/body/div[2]/div/div[2]/div[2]/div[2]/div[@class="search-list"]`)
		fmt.Println("=======一共有====", len(divNodes), "====记录=======")
		if len(divNodes) >= 1 {
			for _, divNode := range divNodes {
				fmt.Println("=====================开始处理数据=========================")
				titleNode := htmlquery.FindOne(divNode, `./div[1]/div/a[2]/span/@title`)
				title := htmlquery.InnerText(titleNode)
				title = strings.TrimSpace(title)
				title = strings.ReplaceAll(title, "-", "")
				title = strings.ReplaceAll(title, " ", "")
				title = strings.ReplaceAll(title, "|", "-")
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, "　", "-")
				fmt.Println(title)

				codeNode := htmlquery.FindOne(divNode, `./div[1]/div/a[1]/span`)
				code := htmlquery.InnerText(codeNode)
				code = strings.ReplaceAll(code, "[国家标准]", "")
				code = strings.ReplaceAll(code, "[行业标准]", "")
				code = strings.ReplaceAll(code, "[计量规程规范]", "")
				code = strings.ReplaceAll(code, "\n", "")
				code = strings.ReplaceAll(code, "\r", "")
				code = strings.ReplaceAll(code, "\t", "")
				code = strings.Replace(code, " ", "", 1)
				fmt.Println(code)
// 				标准号中含有“/”，推荐性标准不提供在线阅读，跳过
                if strings.Index(code,"/") != -1{
                    fmt.Println("推荐性标准不提供在线阅读，跳过")
                    continue
                }

				filePath := "../bz.faxin.cn/" + title + "(" + strings.ReplaceAll(code, "/", "-") + ")" + ".pdf"
				fmt.Println(filePath)

				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				// 获取encstr
				faXinEncStrUrl := fmt.Sprintf("https://bz.faxin.cn/faxin/stdlib/getencstr?a100=%s", strings.ReplaceAll(code, " ", "%20"))
				fmt.Println(faXinEncStrUrl)
				faXinEncStrReferer := fmt.Sprintf("https://bz.faxin.cn/faxin/view/online/%s/?", strings.ReplaceAll(code, " ", "%20"))
				faXinEncStrDoc, err := getFaXinEncStr(faXinEncStrUrl, faXinEncStrReferer)
				if err != nil {
					fmt.Println(err)
					continue
				}
				faXinEncStr := htmlquery.InnerText(faXinEncStrDoc)
				fmt.Println(faXinEncStr)
				if len(faXinEncStr) < 32 {
					fmt.Println("获取EncStr失败")
					break
				}
				// 获取在线阅读页面内容
				queryFaXinStdOnlineFormData := QueryFaXinStdOnlineFormData{
					a100:       code,
					saleagtid:  "0114",
					username:   "",
					readdevice: "1",
					encstr:     faXinEncStr,
				}
				stdOnlineUrl := "https://www.spc.org.cn/gb168/agtvip/stdonline"
				queryFaXinStdOnlineDoc, err := QueryFaXinStdOnline(stdOnlineUrl,queryFaXinStdOnlineFormData)
				if err != nil {
// 				[{"errorCode":"804010016","value":"因版权限制，此标准不提供在线阅读"}]
					fmt.Println(err)
					continue
				}

				// 获取在线阅读enc
				regFaXinEnc := regexp.MustCompile(`var enc = "(.*?)";`)
				regFaXinEncMatch := regFaXinEnc.FindAllSubmatch([]byte(htmlquery.InnerText(queryFaXinStdOnlineDoc)), -1)
				if len(regFaXinEncMatch) == 0 {
					fmt.Println("获取enc失败")
					break
				}
				faXinEnc := string(regFaXinEncMatch[0][1])
				fmt.Println("enc==", faXinEnc)
				downloadMyFoxit := faXinEnc

				// 获取在线阅读rc
				regFaXinRc := regexp.MustCompile(`var rc = "(.*?)";`)
				regFaXinRcMatch := regFaXinRc.FindAllSubmatch([]byte(htmlquery.InnerText(queryFaXinStdOnlineDoc)), -1)
				if len(regFaXinRcMatch) == 0 {
					fmt.Println("获取rc失败")
					break
				}
				faXinRc := string(regFaXinRcMatch[0][1])
				fmt.Println("rc==", faXinRc)

				downloadReferer := "https://www.spc.org.cn/gb168/agtvip/stdonline"
				downloadUrl := fmt.Sprintf("https://www.spc.org.cn/stdlib/onlinereading?token=%s&type=", faXinRc)
				fmt.Println(downloadUrl)

				fmt.Println("=======开始下载" + title + "========")
				err = downloadFaXin(downloadUrl, downloadReferer, downloadMyFoxit, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "bz.faxin.cn", "temp-hbba.sacinfo.org.cn")
				err = copyFaXinFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")
				//DownLoadFaXinTimeSleep := 10
				DownLoadFaXinTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadFaXinTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadFaXinTimeSleep, "秒 倒计时", i, "秒===========")
				}
			}
		}
		DownLoadFaXinPageTimeSleep := 10
		// DownLoadFaXinPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadFaXinPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadFaXinPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

func QueryFaXinList(requestUrl string, queryFaXinListFormData QueryFaXinListFormData) (doc *html.Node, err error) {
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
	if FaXinEnableHttpProxy {
		client = FaXinSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("stdno", queryFaXinListFormData.stdno)
	postData.Add("sc", queryFaXinListFormData.sc)
	postData.Add("stdname", queryFaXinListFormData.stdname)
	postData.Add("a404", queryFaXinListFormData.a404)
	postData.Add("a825", queryFaXinListFormData.a825)
	postData.Add("standStatus", queryFaXinListFormData.standStatus)
	postData.Add("issueDateStart", queryFaXinListFormData.issueDateStart)
	postData.Add("issueDateEnd", queryFaXinListFormData.issueDateEnd)
	postData.Add("a205Start", queryFaXinListFormData.a205Start)
	postData.Add("a205End", queryFaXinListFormData.a205End)
	postData.Add("issueDepart", queryFaXinListFormData.issueDepart)
	postData.Add("draftsDept", queryFaXinListFormData.draftsDept)
	postData.Add("reader", queryFaXinListFormData.reader)
	postData.Add("ownerDept", queryFaXinListFormData.ownerDept)
	postData.Add("pageIndex", strconv.Itoa(queryFaXinListFormData.pageIndex))
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", FaXinCookie)
	req.Header.Set("Host", "bz.faxin.cn")
	req.Header.Set("Origin", "https://bz.faxin.cn")
	req.Header.Set("Referer", "https://bz.faxin.cn/faxin/view/advancedsearch")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return doc, err
	}
	defer resp.Body.Close()
//     body, err := ioutil.ReadAll(resp.Body)
//     if err != nil {
//     fmt.Println(err)
//         return
//     }
//     fmt.Println(string(body))
//     os.Exit(1)
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

func getFaXinEncStr(requestUrl string, refererUrl string) (doc *html.Node, err error) {
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
	if FaXinEnableHttpProxy {
		client = FaXinSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", FaXinEncStrCookie)
	req.Header.Set("Host", "bz.faxin.cn")
	req.Header.Set("Origin", "https://bz.faxin.cn")
	req.Header.Set("Referer", refererUrl)
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return doc, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
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

func QueryFaXinStdOnline(requestUrl string, queryFaXinStdOnlineFormData QueryFaXinStdOnlineFormData) (doc *html.Node, err error) {
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
	if FaXinEnableHttpProxy {
		client = FaXinSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("a100", queryFaXinStdOnlineFormData.a100)
	postData.Add("saleagtid", queryFaXinStdOnlineFormData.saleagtid)
	postData.Add("username", queryFaXinStdOnlineFormData.username)
	postData.Add("readdevice", queryFaXinStdOnlineFormData.readdevice)
	postData.Add("encstr", queryFaXinStdOnlineFormData.encstr)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", SpcStdOnlineCookie)
	req.Header.Set("Host", "www.spc.org.cn")
	req.Header.Set("Origin", "https://bz.faxin.cn")
	req.Header.Set("Referer", "https://bz.faxin.cn/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return doc, err
	}
	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
//    if err != nil {
//       fmt.Println(err)
//       return
//    }
//    fmt.Println(string(body))
//    os.Exit(1)

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

func downloadFaXin(attachmentUrl string, referer string, myFoxit string, filePath string) error {
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
	if FaXinEnableHttpProxy {
		client = FaXinSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "*/*")
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", SpcDownloadCookie)
	req.Header.Set("Host", "www.spc.org.cn")
	req.Header.Set("Myfoxit", myFoxit)
	req.Header.Set("Origin", "https://www.spc.org.cn")
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
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 如果访问失败 就打印当前状态码
// 	if resp.StatusCode != http.StatusOK {
// 		return errors.New("http status :" + strconv.Itoa(resp.StatusCode))
// 	}

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

func copyFaXinFile(src, dst string) (err error) {
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
