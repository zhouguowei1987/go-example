package main

import (
	"errors"
	"fmt"
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

	"golang.org/x/net/html"

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

var HnBzwCookie = "ASP.NET_SessionId=kam0er5unlifw4ehruoltcvh"
var HnBzwPdfCookie = "ASP.NET_SessionId=04gnax45h42x5uzmleajxgzf"

// 下载湖南省标准信息公共服务平台文档
// @Title 下载湖南省标准信息公共服务平台文档
// @Description https://www.hnbzw.com/，下载湖南省标准信息公共服务平台文档
func main() {
	pageListUrl := "https://www.hnbzw.com/Standard/LocalStdSeach.aspx"
	fmt.Println(pageListUrl)
	startPage := 1
	isPageListGo := true
	for isPageListGo {
		queryHnBzwListFormData := QueryHnBzwListFormData{
			VIEWSTATE:          "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",
			VIEWSTATEGENERATOR: "86D8EADD",
			EVENTTARGET:        "AspNetPager1",
			EVENTARGUMENT:      startPage,
			EVENTVALIDATION:    "+6PWEMxXEs4EPttWoEXvTR6g707iw4p7B0EAObxFRiovJ5tj0ffJTJWlEVjkNMcmgAJvhTr1Wm96+JYyEKzAoxDdNdTds98rtFYj4KYgv37ak7v0REs5f52NTPJ+WbeVtgDwybMupZguHVp+hNXfk6zf0ujYGMhlKVS5dbKXXGP/W8S446d8icJLu46CFLzIXs4/DnnkDbngb90xutwzt6vWtkZPKkT+925EUJuwEv1vsz9pKIcogrnamS9poFG4l9sVDtJvzN0J0NFE2HqTdVmTSBPktO6yVQL8NaKLImnwfjN6z4rUGStYgmhxlUndURwE5JAab4faUMHSj0UZci7+t/k+N89uqZgDCmqopsY7Kf81Y2sDBS7nYn7UHcUYF95Q0YdAXeN/SRUKBtsifvFq0E0WrfYpwAaQVDgDFacvYT3HuG97urFFNzrd6frYnxNN+gcXeqdl+EzOlhnSi/MYh+9ga4ulUACXVl0lYCDAEeTaHybL3sxd8VaWbUT7UCJC//mB+zbwZUfOZU+ppGnPbLMd3D99nLy3yeA1gw+fllsLzt+07MGGITlCym0L2V20L1WxnH7XKHg0FYKS8HajMNw9BjaqD9D71A+KYssK36VV7lhynHw7UY2X94L3fevhYw5TLfn+xRaUN3EdBOA0XiUgBmwRk2zPkQfs6VuOjDQHAvVWOL3br08kbmDpseHaXCTiSIsA0NCA14tF5O+QcyJfLaC90Jnjmr1PBimGXvUeXP0JknZJhcTtTvSsUfMQICKYa9M5KBV1KCACRhcr7KZDczLSSO0fYM8yDGUfghXhyh0v0rj0qDlptu29i5ZPbs2iqcQTHLnU/Ny7Sve6s8Ayn0OHMtrU15wPwCcxJhJ+NVOV05ODc+kOcxtNlFd5Tii8zyLD+3isP9yLXzDvTyxJXxNklYu/4/xgG7zVJDPeQrly2rYa1TPpM2vLfanzh404TCHxWTrD+iJoFR7OU99w2plRmuWXoHzZzbrq/j8uGyQLErPrOhBxn1JmGIaOwd+6fWHkMYUMl0Yw6ZRGVix8+GJoIf1MMhBf4PUwZLgffaTAuO+Msk/MfXeKT1fNXx17fK0Xhigu9kYZ8Mb9ImdQDdQ+BhbktOVvLoGiyoH2eKXJwyP6DPVDKZnBxz5jlDLYhuze8cMiXtZe1QomfFMp9k3KSU1SwLWrHSZlXyfzUnAgdtJYevZoRGrv",
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
				tempFilePath := strings.ReplaceAll(filePath, "www.hnbzw.com", "temp-dbba.sacinfo.org.cn")
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
