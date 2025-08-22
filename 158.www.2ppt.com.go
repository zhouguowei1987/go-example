package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

const (
	Ppt2EnableHttpProxy = false
	Ppt2HttpProxyUrl    = "111.225.152.186:8089"
)

func Ppt2SetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(Ppt2HttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type QueryPpt2DownloadUrlFormData struct {
	resourceid string
}

var Ppt2Cookie = "Hm_lvt_cd668c52b64b6c51259bb01b5a1ca278=1755330154; HMACCOUNT=1CCD0111717619C6; Hm_lpvt_cd668c52b64b6c51259bb01b5a1ca278=1755392773"

// ychEduSpider 获取爱ppt文档
// @Title 获取爱ppt文档
// @Description https://www.2ppt.com/，获取爱ppt文档
func main() {
	page := 1
	isPageListGo := true
	for isPageListGo {
		requestListUrl := fmt.Sprintf("https://www.2ppt.com/ppt/511/%d.html", page)
		referUrl := "https://www.2ppt.com/ppt/511/1.html"
		if page >= 2 {
			referUrl = fmt.Sprintf("https://www.2ppt.com/ppt/488/%d.html", page-1)
		}
		fmt.Println(requestListUrl)
		Ppt2ListDoc, err := Ppt2HtmlDoc(requestListUrl, referUrl)
		if err != nil {
			fmt.Println(err)
			break
		}
		liNodes := htmlquery.Find(Ppt2ListDoc, `//div[@class="index-cws w1400"]/div[@class="cws-list"]/ul[@class="clearfix"]/li`)
		if len(liNodes) >= 1 {
			for _, liNode := range liNodes {
				titleNode := htmlquery.FindOne(liNode, `./a/@title`)
				if titleNode == nil {
					fmt.Println("标题不存在，跳过")
					continue
				}
				title := htmlquery.InnerText(titleNode)
				title = strings.TrimSpace(title)
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, " ", "")
				fmt.Println(title)

				filePath := "F:\\workspace\\www.2ppt.com\\www.2ppt.com\\" + title + ".pptx"
				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				// 下载文档URL
				queryPpt2DownloadUrl := "https://www.2ppt.com/ppt/down.html"
				fmt.Println(queryPpt2DownloadUrl)

				buttonNode := htmlquery.FindOne(liNode, `./div[@class="cw-footer"]/a`)
				// down('ppt','oli')
				clickText := htmlquery.SelectAttr(buttonNode, "onclick")
				clickText = strings.ReplaceAll(clickText, "down('ppt','", "")
				clickText = strings.ReplaceAll(clickText, "')", "")
				resourceid := clickText
				fmt.Println(resourceid)

				queryPpt2DownloadUrlFormData := QueryPpt2DownloadUrlFormData{
					resourceid: resourceid,
				}

				queryPpt2DownloadUrlResponse, err := QueryPpt2DownloadUrl(queryPpt2DownloadUrl, queryPpt2DownloadUrlFormData, requestListUrl)
				if err != nil {
					fmt.Println(err)
					continue
				}

				downLoadUrl := queryPpt2DownloadUrlResponse.Data
				fmt.Println(downLoadUrl)
				if strings.Index(downLoadUrl, ".pptx") == -1 {
					fmt.Println("不是pptx文件，跳过")
					continue
				}

				// 开始下载
				fmt.Println("=======开始下载========")
				err = downloadPpt2(downLoadUrl, requestListUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======完成下载========")
				// 设置倒计时
				DownLoadTPpt2TimeSleep := 10
				for i := 1; i <= DownLoadTPpt2TimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page = "+strconv.Itoa(page)+"===title="+title+"===========操作完成，", "暂停", DownLoadTPpt2TimeSleep, "秒，倒计时", i, "秒===========")
				}
			}
			DownLoadPpt2PageTimeSleep := 10
			// DownLoadPpt2PageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadPpt2PageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page = "+strconv.Itoa(page)+"========= 暂停", DownLoadPpt2PageTimeSleep, "秒 倒计时", i, "秒===========")
			}
			page++
		} else {
			page = 0
			isPageListGo = false
			break
		}
	}
}

func Ppt2HtmlDoc(requestUrl string, referer string) (doc *html.Node, err error) {
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
	if Ppt2EnableHttpProxy {
		client = Ppt2SetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接
	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", Ppt2Cookie)
	req.Header.Set("Host", "www.2ppt.com")
	req.Header.Set("Origin", "https://www.2ppt.com/")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
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

type QueryPpt2DownloadUrlResponse struct {
	Code  int    `json:"code"`
	Count int    `json:"count"`
	Data  string `json:"data"`
	Msg   string `json:"msg"`
}

func QueryPpt2DownloadUrl(requestUrl string, queryPpt2DownloadUrlFormData QueryPpt2DownloadUrlFormData, referer string) (queryPpt2DownloadUrlResponse QueryPpt2DownloadUrlResponse, err error) {
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
	if Ppt2EnableHttpProxy {
		client = Ppt2SetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("resourceid", queryPpt2DownloadUrlFormData.resourceid)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		return queryPpt2DownloadUrlResponse, err
	}

	req.Header.Set("authority", "www.2ppt.com")
	req.Header.Set("method", "POST")
	req.Header.Set("path", "/ppt/down.html")
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", Ppt2Cookie)
	req.Header.Set("Origin", "https://www.2ppt.com")
	req.Header.Set("Host", "www.2ppt.com")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"124\", \"Google Chrome\";v=\"124\", \"Not-A.Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryPpt2DownloadUrlResponse, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryPpt2DownloadUrlResponse, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryPpt2DownloadUrlResponse, err
	}
	err = json.Unmarshal(respBytes, &queryPpt2DownloadUrlResponse)
	if err != nil {
		return queryPpt2DownloadUrlResponse, err
	}
	return queryPpt2DownloadUrlResponse, nil
}

func downloadPpt2(downloadUrl string, referer string, filePath string) error {
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
	if Ppt2EnableHttpProxy {
		client = Ppt2SetHttpProxy()
	}
	req, err := http.NewRequest("GET", downloadUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", Ppt2Cookie)
	req.Header.Set("Origin", "https://www.2ppt.com")
	req.Header.Set("Host", "www.2ppt.com")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"124\", \"Google Chrome\";v=\"124\", \"Not-A.Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
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
