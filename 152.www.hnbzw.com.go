package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
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

var HnBzwEnableHttpProxy = false
var HnBzwHttpProxyUrl = "111.225.152.186:8089"
var HnBzwHttpProxyUrlArr = make([]string, 0)

func HnBzwHttpProxy() error {
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
					HnBzwHttpProxyUrlArr = append(HnBzwHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					HnBzwHttpProxyUrlArr = append(HnBzwHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func HnBzwSetHttpProxy() (httpclient *http.Client) {
	if HnBzwHttpProxyUrl == "" {
		if len(HnBzwHttpProxyUrlArr) <= 0 {
			err := HnBzwHttpProxy()
			if err != nil {
				HnBzwSetHttpProxy()
			}
		}
		HnBzwHttpProxyUrl = HnBzwHttpProxyUrlArr[0]
		if len(HnBzwHttpProxyUrlArr) >= 2 {
			HnBzwHttpProxyUrlArr = HnBzwHttpProxyUrlArr[1:]
		} else {
			HnBzwHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(HnBzwHttpProxyUrl)
	ProxyURL, _ := url.Parse(HnBzwHttpProxyUrl)
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

type QueryHnBzwListFormData struct {
	VIEWSTATE          string
	VIEWSTATEGENERATOR string
	EVENTTARGET        string
	EVENTARGUMENT      int
	EVENTVALIDATION    string
	txtKey             string
	txtNo              string
	drpState           int
	drpYear            int
	txtDrafter         string
	txtDrafterMan      string
}

var HnBzwCookie = "ASP.NET_SessionId=s3ubburyp1cfifdu2sixxmnu"
var HnBzwPdfCookie = "ASP.NET_SessionId=04gnax45h42x5uzmleajxgzf"

// 下载湖南省标准信息公共服务平台文档
// @Title 下载湖南省标准信息公共服务平台文档
// @Description https://www.hnbzw.com/，下载湖南省标准信息公共服务平台文档
func main() {
	pageListUrl := "https://www.hnbzw.com/Standard/LocalStdSeach.aspx"
	fmt.Println(pageListUrl)
	startPage := 310
	isPageListGo := true
	for isPageListGo {
		queryHnBzwListFormData := QueryHnBzwListFormData{
			VIEWSTATE:          "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",
			VIEWSTATEGENERATOR: "86D8EADD",
			EVENTTARGET:        "AspNetPager1",
			EVENTARGUMENT:      startPage,
			EVENTVALIDATION:    "uW3LBifzXocVEoKJ4edSNYBiqsAAVrEFXapm31UVleVItp+2H5wtQFnJFk4/GoN4HP8PC1vgE0ShkMCdqifiF8eS57bzFIXYvfM/ZV+faqKpx40DEW5ALkw/Gh9lT5i6LZq4Obq4yyoaQ1z974Yn9OltOUJwqjHfHm8ggd37tTm0useKDya6/0p3EO8vepxpVA2w1wz/QyepaIE75NH60SBq9qqAotxLFjPlIMBWGn9LSUM4ji9hvvDJG+JvFYqk8Pid8bbbRYBD3Ez/foKl7YuTcRCv9i7zOO+Jdce8PpHvZnzlNSCfVBakGr3VF9MqOvl3KSEhk/fl7kQM+x3I/cojfUglVpYXJrMkH/rdvg+odbbe3734iWfF4bhvYrV5uyn3c+97p3dvra5cGW/kxdRtM5bHDT/QlMN5JEoNusbVgHUD2lcG0/m0fdCw6Sj3SuKK+82x5o/v5Jad0tn1+mbRv/Ta8y+XaUSiwHJYc9IoJKRBQ3brBhc/EAFaEaDYGrU83EcMXzhcSgJhRomePwh+nUwuiDxMAOm5uJfGlGA8wGW0a97IMtfVFVTWYZrxVls6KrWmktz1BWHPldHVuYR62X1rlwjQql9XBIDuGGuqKHXsGYj88X15C8eKgg2QUISTv0KtsouDX082KCCsa3QPrQpH4JwilzyCuI+wkQkBNbI22GSHp5mbzJWUAz1u8QrrEvjBRtfpsEUc1jBXNlrMa+6XW33RhIyr7Qd6PoDOYBhhmAUkKohiH/xCJ/GTKuDKmbj61Jrg00yejjlHaDdz9taFuTk5ZlX9FGyEbBD4jFdGw8Ohe6awPORfKu7sbx/yoYFrsjzPkbyJROvncrReFydv3m9ij6sqrXyAxHtN65IXmIfo6WAaafpHqwZ39ho+NdOjf4/LiHkZmcLkWrI5AYqRHB01pRNYf/d/wTo3/7b5eKATSWHbmdBaVnDktA6L3TB9sPVVRBXnl0dlIDOUIjWd7HCg9MUQBsbQwKxxgdqvvZ3ai15ICxYpZRgprBkUyuBNlMeXoEDNCPqp3bK6EEEWulSVhPinxygNKiz7i4T4JBwCnp6U3+Io+uDN9cXxSC+IaBb6We0cFx9vPQ+uCQlkPxxrt/yhZYq7l3s/Ch+AhfdZXj48DBa453yUggw64VwFFoVoPxaOZb/s6HkcLOgVvM8AkXfS/ULlvwX5DPV02OSVxN/IK3ScDude",
			txtKey:             "",
			txtNo:              "",
			drpState:           0,
			drpYear:            0,
			txtDrafter:         "",
			txtDrafterMan:      "",
		}
		queryHnBzwListDoc, err := QueryHnBzwList(pageListUrl, queryHnBzwListFormData)
		if err != nil {
			fmt.Println(err)
			break
		}
		ulNodes := htmlquery.Find(queryHnBzwListDoc, `//form[@id="form1"]/div[@class="gj-cx-b"]/div[@class="con-gj"]/div[@class="gj-cx"]/div[@class="gj-zd"]/div[@class="gj-zd"]/div[@class="gjlb"]/ul[@class="list-mc"]`)
		if len(ulNodes) > 0 {
			for _, ulNode := range ulNodes {
				fmt.Println("=====================开始处理数据 page = ", startPage, "=========================")

				codeNode := htmlquery.FindOne(ulNode, `./li[2]/span[1]`)
				code := htmlquery.InnerText(codeNode)
				code = strings.ReplaceAll(code, "/", "-")
				code = strings.ReplaceAll(code, "—", "-")
				fmt.Println(code)

				titleNode := htmlquery.FindOne(ulNode, `./li[3]`)
				title := htmlquery.InnerText(titleNode)
				title = strings.TrimSpace(title)
				title = strings.ReplaceAll(title, " ", "-")
				title = strings.ReplaceAll(title, "　", "-")
				title = strings.ReplaceAll(title, "中文名：", "")
				title = strings.ReplaceAll(title, ":", "-")
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, "--", "-")
				fmt.Println(title)

				filePath := "../www.hnbzw.com/" + title + "(" + code + ")" + ".pdf"
				fmt.Println(filePath)

				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				// 查看是否有下载按钮
				buttonNode := htmlquery.FindOne(ulNode, `./li[1]/span[@class="lisy-xq"]/a[1]`)
				hrefText := htmlquery.SelectAttr(buttonNode, "href")
				if len(hrefText) <= 0 {
					fmt.Println("没有预览按钮，跳过")
					continue
				}
				previewHref := "https://www.hnbzw.com" + hrefText
				fmt.Println(previewHref)

				previewDoc, err := previewHnBzwDoc(previewHref)
				if err != nil {
					fmt.Println("获取文档详情失败，跳过")
					continue
				}
				// /html/body/form/div[4]/div/div/div/ul/table/tbody/tr[2]/td/input
				downloadButtonNode := htmlquery.FindOne(previewDoc, `//form[@id="form1"]/div[@class="gj-cx-b"]/div[@class="con-gj"]/div[@class="gj-cx-goods"]/div[@class="gj-zd"]/ul[@class="gj-bts"]/table/tbody/tr[2]/td/input`)
				if downloadButtonNode == nil {
					fmt.Println("没有下载按钮，跳过")
					continue
				}
				fmt.Println("=======开始下载========")

				// /html/body/form/input
				downloadNode := htmlquery.FindOne(previewDoc, `//html/body/form/input/@value`)
				if downloadNode == nil {
					fmt.Println("获取下载链接失败，跳过")
					continue
				}
				downloadUrl := htmlquery.InnerText(downloadNode)
				fmt.Println(downloadUrl)

				fmt.Println("=======开始下载" + title + "========")

				err = downloadHnBzw(downloadUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "../www.hnbzw.com", "../temp-www.hnbzw.com")
				err = copyHnBzwFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")
				//DownLoadHnBzwTimeSleep := 10
				DownLoadHnBzwTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadHnBzwTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("startPage="+strconv.Itoa(startPage)+",filePath="+filePath+"===========下载成功 暂停", DownLoadHnBzwTimeSleep, "秒 倒计时", i, "秒===========")
				}
			}
			DownLoadHnBzwPageTimeSleep := 10
			// DownLoadHnBzwPageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadHnBzwPageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("startPage="+strconv.Itoa(startPage)+"=========== 暂停", DownLoadHnBzwPageTimeSleep, "秒 倒计时", i, "秒===========")
			}
			startPage++
		} else {
			isPageListGo = false
			startPage = 1
			break
		}
	}
}

func QueryHnBzwList(requestUrl string, queryHnBzwListFormData QueryHnBzwListFormData) (doc *html.Node, err error) {
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
	if HnBzwEnableHttpProxy {
		client = HnBzwSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("__VIEWSTATE", queryHnBzwListFormData.VIEWSTATE)
	postData.Add("__VIEWSTATEGENERATOR", queryHnBzwListFormData.VIEWSTATEGENERATOR)
	postData.Add("__EVENTTARGET", queryHnBzwListFormData.EVENTTARGET)
	postData.Add("__EVENTARGUMENT", strconv.Itoa(queryHnBzwListFormData.EVENTARGUMENT))
	postData.Add("__EVENTVALIDATION", queryHnBzwListFormData.EVENTVALIDATION)
	postData.Add("txtKey", queryHnBzwListFormData.txtKey)
	postData.Add("txtNo", queryHnBzwListFormData.txtNo)
	postData.Add("drpState", strconv.Itoa(queryHnBzwListFormData.drpState))
	postData.Add("drpYear", strconv.Itoa(queryHnBzwListFormData.drpYear))
	postData.Add("txtDrafter", queryHnBzwListFormData.txtDrafter)
	postData.Add("txtDrafterMan", queryHnBzwListFormData.txtDrafterMan)

	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", HnBzwCookie)
	req.Header.Set("Host", "www.hnbzw.com")
	req.Header.Set("Origin", "https://www.hnbzw.com")
	req.Header.Set("Referer", "https://www.hnbzw.com/Standard/StdSearch.aspx")
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

func previewHnBzwDoc(requestUrl string) (doc *html.Node, err error) {
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
	if HnBzwEnableHttpProxy {
		client = HnBzwSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", HnBzwCookie)
	req.Header.Set("Host", "www.hnbzw.com")
	req.Header.Set("Origin", "https://www.hnbzw.com")
	req.Header.Set("Referer", "https://www.hnbzw.com/Standard/StdSearch.aspx")
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

func downloadHnBzw(attachmentUrl string, filePath string) error {
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
	if HnBzwEnableHttpProxy {
		client = HnBzwSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "*/*")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", HnBzwPdfCookie)
	req.Header.Set("Host", "pdf.hnbzw.com")
	req.Header.Set("Referer", "https://www.hnbzw.com/")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"124\", \"Google Chrome\";v=\"124\", \"Not-A.Brand\";v=\"99\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
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

func copyHnBzwFile(src, dst string) (err error) {
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
