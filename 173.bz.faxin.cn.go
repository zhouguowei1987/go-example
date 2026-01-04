package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
// 	"io/ioutil"
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
			ResponseHeaderTimeout: time.Second * 3,
		},
	}
	return httpclient
}

type QueryFaXinListFormData struct {
	sctype         string
	stdno string
	sc       string
	stdname       string
	a404         string
	a825         string
	standStatus         string
	issueDateStart         string
	issueDateEnd         string
	a205Start         string
	a205End         string
	issueDepart         string
	draftsDept         string
	reader         string
	ownerDept         string
	pageIndex         int
}

var FaXinCookie = "JSESSIONID=575479BC0D2ADCED3066678E827D6499; Hm_lvt_a317640b4aeca83b20c90d410335b70f=1766626702; HMACCOUNT=4E5B3419A3141A8E; Hm_lvt_cb8a2025f4234726e55e45f893fb7954=1766626740; Hm_lpvt_a317640b4aeca83b20c90d410335b70f=1766992871; Hm_lpvt_cb8a2025f4234726e55e45f893fb7954=1766994755"

// 下载法信标准文档
// @Title 下载法信标准文档
// @Description https://bz.faxin.cn/，下载法信标准文档
func main() {
	pageListUrl := "https://bz.faxin.cn/faxin/view/advancedsearch"
	fmt.Println(pageListUrl)
	page := 0
	maxPage := 29
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		queryFaXinListFormData := QueryFaXinListFormData{
			sctype:"",
            stdno:"",
            sc:"CN",//CN：国家标准 QT：行业标准 JJ：计量规程规范
            stdname:"",
            a404:"",
            a825:"*",
            standStatus:"",
            issueDateStart:"",
            issueDateEnd:"",
            a205Start:"",
            a205End:"",
            issueDepart:"",
            draftsDept:"",
            reader:"",
            ownerDept:"",
            pageIndex:page,
		}
		queryFaXinListDoc, err := QueryFaXinList(pageListUrl, queryFaXinListFormData)
// 		fmt.Println(htmlquery.InnerText(queryFaXinListDoc))
// 		os.Exit(1)
		if err != nil {
			fmt.Println(err)
			break
		}
		divNodes := htmlquery.Find(queryFaXinListDoc, `//html/body/div[2]/div/div[2]/div[2]/div[2]/div`)
        if len(divNodes) >= 1 {
            for _, divNode := range divNodes {
                fmt.Println("=====================开始处理数据=========================")
				codeNode := htmlquery.FindOne(divNode, `./div[1]/div/a[1]/span`)
				code := htmlquery.InnerText(codeNode)
				code = strings.ReplaceAll(code,"[国家标准]","")
				code = strings.ReplaceAll(code,"[行业标准]","")
				code = strings.ReplaceAll(code,"[计量规程规范]","")
				code = strings.ReplaceAll(code,"\n","")
				code = strings.ReplaceAll(code,"\r","")
				code = strings.ReplaceAll(code,"\t","")
				code = strings.Replace(code," ","",1)
				fmt.Println(code)

				titleNode := htmlquery.FindOne(divNode, `./div[1]/div/a[2]/span/@title`)
				title := htmlquery.InnerText(titleNode)
				title = strings.TrimSpace(title)
				title = strings.ReplaceAll(title, "-", "")
				title = strings.ReplaceAll(title, " ", "")
				title = strings.ReplaceAll(title, "|", "-")
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, "　", "-")
				fmt.Println(title)

				filePath := "../bz.faxin.cn/" + title + "(" + strings.ReplaceAll(code, "/", "-") + ")" + ".pdf"
				fmt.Println(filePath)

				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				detailUrl := fmt.Sprintf("https://bz.faxin.cn/faxin/view/online/%s",code)

				downloadUrl := fmt.Sprintf("https://bz.faxin.cn/faxin/view/haveprevpage?stdno=%s",code)
                downloadUrl = strings.ReplaceAll(downloadUrl," ","%20")
                fmt.Println(downloadUrl)

                fmt.Println("=======开始下载" + title + "========")
                err = downloadFaXin(downloadUrl, detailUrl, filePath)
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
			ResponseHeaderTimeout: time.Second * 3,
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
	req.Header.Set("Referer", "https://bz.faxin.cn/jianybz/jybzTwoGj.do?formAction=listQxtjbz")
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

func downloadFaXin(attachmentUrl string, referer string, filePath string) error {
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
	if FaXinEnableHttpProxy {
		client = FaXinSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", FaXinCookie)
	req.Header.Set("Host", "bz.faxin.cn")
	req.Header.Set("Origin", "https://bz.faxin.cn")
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
