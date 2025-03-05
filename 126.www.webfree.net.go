package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
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
)

const (
	WebFreeEnableHttpProxy = false
	WebFreeHttpProxyUrl    = "111.225.152.186:8089"
)

func WebFreeSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(WebFreeHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type WebFree struct {
	name string
	url  string
}

var webfrees = []WebFree{
	{
		name: "国家标准",
		url:  "https://www.webfree.net/downloads/gb",
	},
	//{
	//	name: "行业标准",
	//	url:  "https://www.webfree.net/hangye-biaozhun",
	//},
	//{
	//	name: "地方标准",
	//	url:  "https://www.webfree.net/difang-biaozhun",
	//},
	//{
	//	name: "书籍图集",
	//	url:  "https://www.webfree.net/downloads/book-and-drawings",
	//},
}

type DownLoadWebFreeFormData struct {
	__wpdm_ID string
	dataType  string
	execute   string
	action    string
	password  string
}

type DownLoadWebFreeResponse struct {
	DownLoadUrl string `json:"downloadurl"`
	Success     bool   `json:"success"`
}

// ychEduSpider 协筑资源标准文档
// @Title 协筑资源标准文档
// @Description https://www.webfree.net/，协筑资源标准文档
func main() {
	for _, webfree := range webfrees {
		current := 343
		minCurrent := 1
		isPageListGo := true
		for isPageListGo {
			if current < minCurrent {
				isPageListGo = false
				break
			}
			webfreeIndexUrl := webfree.url
			if current > 1 {
				webfreeIndexUrl += fmt.Sprintf("/?cp=%d", current)
			}
			webfreeIndexDoc, err := htmlquery.LoadURL(webfreeIndexUrl)
			if err != nil {
				fmt.Println(err)
				break
			}
			liNodes := htmlquery.Find(webfreeIndexDoc, `//div[@id="content_wpdm_package_1"]/div[@class="row"]/div[@class="col-lg-12 col-md-12 col-12"]`)
			if len(liNodes) <= 0 {
				isPageListGo = false
				fmt.Println("没有数据暂停")
				break
			}
			for _, liNode := range liNodes {
				fmt.Println("============================================================================")
				fmt.Println("标准类别：", webfree.name)
				fmt.Println("=======当前页为：" + strconv.Itoa(current) + "========")

				fileANode := htmlquery.FindOne(liNode, `./div[@class="entry entry-cpt"]/div[@class="entry-container"]/div[@class="entry-head"]/h2[@class="entry-title"]/a`)
				fileName := htmlquery.InnerText(fileANode)
				fileName = strings.TrimSpace(fileName)
				fileName = strings.ReplaceAll(fileName, "/", "-")
				fileName = strings.ReplaceAll(fileName, " ", "")
				fmt.Println(fileName)

				filePath := "E:\\workspace\\www.webfree.net\\www.webfree.net/" + webfree.name + "/" + fileName + ".pdf"
				_, errPdf := os.Stat(filePath)
				if errPdf != nil {

					// 标准文件id
					fileId := htmlquery.InnerText(htmlquery.FindOne(fileANode, `./@href`))
					fileId = strings.ReplaceAll(fileId, "https://www.webfree.net/download/", "")
					fmt.Println(fileId)

					downLoadWebFreeFormData := DownLoadWebFreeFormData{
						__wpdm_ID: fileId,
						dataType:  "json",
						execute:   "wpdm_getlink",
						action:    "wpdm_ajax_call",
						password:  "Webfree.net",
					}

					downLoadWebFreeResponse, err := DownLoadWebFreeUrl(downLoadWebFreeFormData)
					if err != nil {
						fmt.Println(err)
						continue
					}
					if downLoadWebFreeResponse.Success != true {
						fmt.Println(err)
						continue
					}

					downLoadUrl := downLoadWebFreeResponse.DownLoadUrl
					fmt.Println(downLoadUrl)

					fmt.Println("=======开始下载" + strconv.Itoa(current) + "========")
					err = downloadWebFree(downLoadUrl, filePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("=======开始完成========")
					time.Sleep(time.Millisecond * 200)
				}
			}
			if current > minCurrent {
				current--
			} else {
				isPageListGo = false
				break
			}
		}
	}
}

func DownLoadWebFreeUrl(downLoadWebFreeFormData DownLoadWebFreeFormData) (downLoadWebFreeResponse DownLoadWebFreeResponse, err error) {
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
	if WebFreeEnableHttpProxy {
		client = WebFreeSetHttpProxy()
	}
	downLoadWebFreeResponse = DownLoadWebFreeResponse{}
	postData := url.Values{}
	postData.Add("__wpdm_ID", downLoadWebFreeFormData.__wpdm_ID)
	postData.Add("dataType", downLoadWebFreeFormData.dataType)
	postData.Add("execute", downLoadWebFreeFormData.execute)
	postData.Add("action", downLoadWebFreeFormData.action)
	postData.Add("password", downLoadWebFreeFormData.password)
	req, err := http.NewRequest("POST", "https://www.webfree.net/wp-json/wpdm/validate-password", strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return downLoadWebFreeResponse, err
	}

	req.Header.Set("authority", "www.webfree.net")
	req.Header.Set("method", "POST")
	req.Header.Set("path", "/wp-json/wpdm/validate-password")
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", "__wpdm_client=3e8ac23582935cd8c0ed700cd0f84b21; PHPSESSID=458e78ddae9c5dbe00147d3fe1531906")
	req.Header.Set("Host", "www.webfree.net")
	req.Header.Set("Origin", "https://www.webfree.net")
	req.Header.Set("Priority", "u=1, i")
	referer := fmt.Sprintf("https://www.webfree.net/?__wpdmlo=%s&REFERRER=https://www.webfree.net/download/%d", downLoadWebFreeFormData.__wpdm_ID, downLoadWebFreeFormData.__wpdm_ID)
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"124\", \"Google Chrome\";v=\"124\", \"Not-A.Brand\";v=\"99\"")
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
		return downLoadWebFreeResponse, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return downLoadWebFreeResponse, err
	}
	err = json.Unmarshal(respBytes, &downLoadWebFreeResponse)
	if err != nil {
		return downLoadWebFreeResponse, err
	}
	return downLoadWebFreeResponse, nil
}

func downloadWebFree(attachmentUrl string, filePath string) error {
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
	if WebFreeEnableHttpProxy {
		client = WebFreeSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "__51cke__=; _gid=GA1.2.944241540.1703144883; PHPSESSID=22d6kpnbct0v3iopp7qctbk081; __tins__21123451=%7B%22sid%22%3A%201703148480266%2C%20%22vd%22%3A%207%2C%20%22expires%22%3A%201703151933380%7D; __51laig__=27; _gat=1; _ga_34B604LFFQ=GS1.1.1703148480.2.1.1703150135.57.0.0; _ga=GA1.1.1587097358.1703144883")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.webfree.net")
	req.Header.Set("Referer", "https://www.webfree.net")
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
