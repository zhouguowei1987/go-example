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

var CooCoEnableHttpProxy = false
var CooCoHttpProxyUrl = ""
var CooCoHttpProxyUrlArr = make([]string, 0)

func CooCoHttpProxy() error {
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
					CooCoHttpProxyUrlArr = append(CooCoHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					CooCoHttpProxyUrlArr = append(CooCoHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func CooCoSetHttpProxy() (httpclient *http.Client) {
	if CooCoHttpProxyUrl == "" {
		if len(CooCoHttpProxyUrlArr) <= 0 {
			err := CooCoHttpProxy()
			if err != nil {
				CooCoSetHttpProxy()
			}
		}
		CooCoHttpProxyUrl = CooCoHttpProxyUrlArr[0]
		if len(CooCoHttpProxyUrlArr) >= 2 {
			CooCoHttpProxyUrlArr = CooCoHttpProxyUrlArr[1:]
		} else {
			CooCoHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(CooCoHttpProxyUrl)
	ProxyURL, _ := url.Parse(CooCoHttpProxyUrl)
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

type CooCoSubject struct {
	name   string
	papers []CooCoSubjectsPaper
}

type CooCoSubjectsPaper struct {
	name string
	url  string
}

var CooCoSubjectsPapers = []CooCoSubject{
	{
		name: "高中",
		papers: []CooCoSubjectsPaper{
			{
				name: "语文",
				url:  "http://bk.cooco.net.cn/shiti/gz/yw",
			},
			{
				name: "数学",
				url:  "http://bk.cooco.net.cn/shiti/gz/sx",
			},
			{
				name: "英语",
				url:  "http://bk.cooco.net.cn/shiti/gz/yy1",
			},
			{
				name: "物理",
				url:  "http://bk.cooco.net.cn/shiti/gz/wl",
			},
			{
				name: "化学",
				url:  "http://bk.cooco.net.cn/shiti/gz/hx",
			},
			{
				name: "地理",
				url:  "http://bk.cooco.net.cn/shiti/gz/dl",
			},
			{
				name: "历史",
				url:  "http://bk.cooco.net.cn/shiti/gz/ls",
			},
			{
				name: "生物",
				url:  "http://bk.cooco.net.cn/shiti/gz/sw",
			},
			{
				name: "政治",
				url:  "http://bk.cooco.net.cn/shiti/gz/zz",
			},
		},
	},
	{
		name: "初中",
		papers: []CooCoSubjectsPaper{
			{
				name: "语文",
				url:  "http://bk.cooco.net.cn/shiti/cz/yw",
			},
			{
				name: "数学",
				url:  "http://bk.cooco.net.cn/shiti/cz/sx",
			},
			{
				name: "英语",
				url:  "http://bk.cooco.net.cn/shiti/cz/yy1",
			},
			{
				name: "物理",
				url:  "http://bk.cooco.net.cn/shiti/cz/wl",
			},
			{
				name: "化学",
				url:  "http://bk.cooco.net.cn/shiti/cz/hx",
			},
			{
				name: "地理",
				url:  "http://bk.cooco.net.cn/shiti/cz/dl",
			},
			{
				name: "历史",
				url:  "http://bk.cooco.net.cn/shiti/cz/ls",
			},
			{
				name: "生物",
				url:  "http://bk.cooco.net.cn/shiti/cz/sw",
			},
			{
				name: "政治",
				url:  "http://bk.cooco.net.cn/shiti/cz/zz",
			},
		},
	},
	{
		name: "小学",
		papers: []CooCoSubjectsPaper{
			{
				name: "语文",
				url:  "http://bk.cooco.net.cn/shiti/xx/yw",
			},
			{
				name: "数学",
				url:  "http://bk.cooco.net.cn/shiti/xx/sx",
			},
			{
				name: "英语",
				url:  "http://bk.cooco.net.cn/shiti/xx/yy1",
			},
			{
				name: "政治",
				url:  "http://bk.cooco.net.cn/shiti/xx/zz",
			},
		},
	},
}

var CooCoNextDownloadSleep = 2

// ychEduSpider 获取备课网文档
// @Title 获取备课网文档
// @Description http://bk.cooco.net.cn/，获取备课网文档
func main() {
	for _, subject := range CooCoSubjectsPapers {
		for _, paper := range subject.papers {
			current := 1
			isPageListGo := true
			// 计算最大页数
			paperIndexUrl := paper.url
			fmt.Println(paperIndexUrl)
			paperIndexDoc, err := htmlquery.LoadURL(paperIndexUrl)
			if err != nil {
				fmt.Println(err)
				continue
			}
			paperMaxPages := 1
			paperTotalEndNode := htmlquery.FindOne(paperIndexDoc, `//div[@class="new-st-recommend"]/div[@class="w1200"]/div[@class="s-r-content"]/div[@class="s-r-stlist ywst"]/div[@id="pager-html"]/ul[@class="pagination"]/li[@class="end"]/a`)
			if paperTotalEndNode != nil {
				paperTotalEndText := htmlquery.InnerText(paperTotalEndNode)
				paperMaxPages, err = strconv.Atoi(paperTotalEndText)
				if err != nil {
					fmt.Println(err)
					continue
				}
			} else {
				paperTotalNumNodes := htmlquery.Find(paperIndexDoc, `//div[@class="new-st-recommend"]/div[@class="w1200"]/div[@class="s-r-content"]/div[@class="s-r-stlist ywst"]/div[@id="pager-html"]/ul[@class="pagination"]/li[@class="num"]`)
				if paperTotalNumNodes != nil {
					paperTotalNumLastNode := paperTotalNumNodes[len(paperTotalNumNodes)-1]
					paperTotalLastNode := htmlquery.FindOne(paperTotalNumLastNode, `./a`)
					paperTotalLastText := htmlquery.InnerText(paperTotalLastNode)
					paperMaxPages, err = strconv.Atoi(paperTotalLastText)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
			}

			for isPageListGo {
				paperListUrl := paper.url + fmt.Sprintf("?page=%d", current)
				paperListDoc, err := htmlquery.LoadURL(paperListUrl)
				if err != nil {
					fmt.Println(err)
					current = 1
					isPageListGo = false
					continue
				}
				liNodes := htmlquery.Find(paperListDoc, `//div[@class="new-st-recommend"]/div[@class="w1200"]/div[@class="s-r-content"]/div[@class="s-r-stlist ywst"]/div[@class="stlist-item clearfix"]`)
				if len(liNodes) <= 0 {
					fmt.Println(err)
					current = 1
					isPageListGo = false
					continue
				}
				for _, liNode := range liNodes {
					fmt.Println("============================================================================")
					fmt.Println("科目：", subject.name, "试卷", paper.name)
					fmt.Println("=======当前页URL", paperListUrl, "========")

					// 文档标题
					title := htmlquery.InnerText(htmlquery.FindOne(liNode, `./div[@class="fl"]/span[@class="stlist-item-name"]/@title`))
					fmt.Println(title)

					// 只下载word文档
					docNode := htmlquery.FindOne(liNode, `./div[@class="fl"]/span[@class="stlist-item-name"]/i[@class="doc"]`)
					if docNode == nil {
						fmt.Println("不是word文档")
						continue
					}
					// 文档id
					docId := htmlquery.InnerText(htmlquery.FindOne(liNode, `./div[@class="fr"]/a/@data-id`))
					fmt.Println(docId)

					err, downUrl := cooCoDownUrl(docId, paperIndexUrl)
					fmt.Println(downUrl)

					filePath := "E:\\workspace\\bk.cooco.net.cn\\bk.cooco.net.cn\\" + subject.name + "\\" + paper.name + "\\" + title + ".doc"
					_, err = os.Stat(filePath)
					if err != nil {

						fmt.Println("=======开始下载"+strconv.Itoa(current)+"-", paperMaxPages, "========")
						err = downloadCooCo(downUrl, paperIndexUrl, filePath)
						if err != nil {
							fmt.Println(err)
							continue
						}
						fmt.Println("=======下载完成========")
						for i := 1; i <= CooCoNextDownloadSleep; i++ {
							time.Sleep(time.Second)
							fmt.Println("===========操作结束，暂停", CooCoNextDownloadSleep, "秒，倒计时", i, "秒===========")
						}
					}
				}
				current++
				if current > paperMaxPages {
					fmt.Println("没有更多分页")
					break
				}
				isPageListGo = true
			}
		}
	}
}

type cooCoDownUrlReturn struct {
	Code string                 `json:"code"`
	Data cooCoDownUrlReturnData `json:"data"`
	Msg  string                 `json:"msg"`
}
type cooCoDownUrlReturnData struct {
	DocGold  int    `json:"doc_gold"`
	DownUrl  string `json:"down_url"`
	UserGold int    `json:"user_gold"`
}

func cooCoDownUrl(docId string, referer string) (err error, downUrl string) {
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
	if CooCoEnableHttpProxy {
		client = CooCoSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("doc_id", docId)
	requestUrl := "http://bk.cooco.net.cn/api/User/download"
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		return err, downUrl
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Host", "bk.cooco.net.cn")
	req.Header.Set("Origin", "http://bk.cooco.net.cn")
	req.Header.Set("Referer", referer)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return err, downUrl
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, downUrl
	}
	CooCoDownUrlReturn := &cooCoDownUrlReturn{}
	err = json.Unmarshal(respBytes, CooCoDownUrlReturn)
	if err != nil {
		return err, downUrl
	}
	if CooCoDownUrlReturn.Code != "SUCCESS" {
		return errors.New(CooCoDownUrlReturn.Msg), downUrl
	}
	downUrl = CooCoDownUrlReturn.Data.DownUrl
	return err, downUrl
}

func downloadCooCo(attachmentUrl string, referer string, filePath string) error {
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
	if CooCoEnableHttpProxy {
		client = CooCoSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "__51cke__=; _gid=GA1.2.944241540.1703144883; PHPSESSID=22d6kpnbct0v3iopp7qctbk081; __tins__21123451=%7B%22sid%22%3A%201703148480266%2C%20%22vd%22%3A%207%2C%20%22expires%22%3A%201703151933380%7D; __51laig__=27; _gat=1; _ga_34B604LFFQ=GS1.1.1703148480.2.1.1703150135.57.0.0; _ga=GA1.1.1587097358.1703144883")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.CooCo.com")
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
