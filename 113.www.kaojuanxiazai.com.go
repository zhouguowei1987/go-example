package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var KaoJuanXiaZaiEnableHttpProxy = false
var KaoJuanXiaZaiHttpProxyUrl = ""
var KaoJuanXiaZaiHttpProxyUrlArr = make([]string, 0)

func KaoJuanXiaZaiHttpProxy() error {
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
					KaoJuanXiaZaiHttpProxyUrlArr = append(KaoJuanXiaZaiHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					KaoJuanXiaZaiHttpProxyUrlArr = append(KaoJuanXiaZaiHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func KaoJuanXiaZaiSetHttpProxy() (httpclient *http.Client) {
	if KaoJuanXiaZaiHttpProxyUrl == "" {
		if len(KaoJuanXiaZaiHttpProxyUrlArr) <= 0 {
			err := KaoJuanXiaZaiHttpProxy()
			if err != nil {
				KaoJuanXiaZaiSetHttpProxy()
			}
		}
		KaoJuanXiaZaiHttpProxyUrl = KaoJuanXiaZaiHttpProxyUrlArr[0]
		if len(KaoJuanXiaZaiHttpProxyUrlArr) >= 2 {
			KaoJuanXiaZaiHttpProxyUrlArr = KaoJuanXiaZaiHttpProxyUrlArr[1:]
		} else {
			KaoJuanXiaZaiHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(KaoJuanXiaZaiHttpProxyUrl)
	ProxyURL, _ := url.Parse(KaoJuanXiaZaiHttpProxyUrl)
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

type KaoJuanXiaZaiPaper struct {
	name string
	url  string
}

type KaoJuanXiaZaiSubjectsPapers struct {
	name   string
	papers []KaoJuanXiaZaiPaper
}

var kaoJuanXiaZaiSubjectsPapers = []KaoJuanXiaZaiSubjectsPapers{
	{
		name: "小学",
		papers: []KaoJuanXiaZaiPaper{
			{
				name: "小学语文",
				url:  "http://www.kaojuanxiazai.com/shijuan/1_0_0_0_0_0_p1/",
			},
			{
				name: "小学数学",
				url:  "http://www.kaojuanxiazai.com/shijuan/2_0_0_0_0_0_p1/",
			},
			{
				name: "小学英语",
				url:  "http://www.kaojuanxiazai.com/shijuan/3_0_0_0_0_0_p1/",
			},
		},
	},

	{
		name: "初中",
		papers: []KaoJuanXiaZaiPaper{
			{
				name: "初中语文",
				url:  "http://www.kaojuanxiazai.com/shijuan/4_0_0_0_0_0_p1/",
			},
			{
				name: "初中数学",
				url:  "http://www.kaojuanxiazai.com/shijuan/5_0_0_0_0_0_p1/",
			},
			{
				name: "初中英语",
				url:  "http://www.kaojuanxiazai.com/shijuan/6_0_0_0_0_0_p1/",
			},
			{
				name: "初中物理",
				url:  "http://www.kaojuanxiazai.com/shijuan/7_0_0_0_0_0_p1/",
			},
			{
				name: "初中化学",
				url:  "http://www.kaojuanxiazai.com/shijuan/8_0_0_0_0_0_p1/",
			},
			{
				name: "初中生物",
				url:  "http://www.kaojuanxiazai.com/shijuan/9_0_0_0_0_0_p1/",
			},
			{
				name: "初中政治",
				url:  "http://www.kaojuanxiazai.com/shijuan/10_0_0_0_0_0_p1/",
			},
			{
				name: "初中历史",
				url:  "http://www.kaojuanxiazai.com/shijuan/11_0_0_0_0_0_p1/",
			},
			{
				name: "初中地理",
				url:  "http://www.kaojuanxiazai.com/shijuan/12_0_0_0_0_0_p1/",
			},
		},
	},

	{
		name: "高中",
		papers: []KaoJuanXiaZaiPaper{
			{
				name: "高中语文",
				url:  "http://www.kaojuanxiazai.com/shijuan/13_0_0_0_0_0_p1/",
			},
			{
				name: "高中数学",
				url:  "http://www.kaojuanxiazai.com/shijuan/14_0_0_0_0_0_p1/",
			},
			{
				name: "高中英语",
				url:  "http://www.kaojuanxiazai.com/shijuan/15_0_0_0_0_0_p1/",
			},
			{
				name: "高中物理",
				url:  "http://www.kaojuanxiazai.com/shijuan/16_0_0_0_0_0_p1/",
			},
			{
				name: "高中化学",
				url:  "http://www.kaojuanxiazai.com/shijuan/17_0_0_0_0_0_p1/",
			},
			{
				name: "高中生物",
				url:  "http://www.kaojuanxiazai.com/shijuan/18_0_0_0_0_0_p1/",
			},
			{
				name: "高中政治",
				url:  "http://www.kaojuanxiazai.com/shijuan/19_0_0_0_0_0_p1/",
			},
			{
				name: "高中历史",
				url:  "http://www.kaojuanxiazai.com/shijuan/20_0_0_0_0_0_p1/",
			},
			{
				name: "高中地理",
				url:  "http://www.kaojuanxiazai.com/shijuan/21_0_0_0_0_0_p1/",
			},
		},
	},
}

const KaoJuanXiaZaiNextDownloadSleep = 1

// ychEduSpider 获取考卷下载试卷
// @Title 获取考卷下载试卷
// @Description http://www.kaojuanxiazai.com/，获取考卷下载试卷
func main() {
	for _, subjectsPapers := range kaoJuanXiaZaiSubjectsPapers {
		for _, paper := range subjectsPapers.papers {
			current := 1
			isPageListGo := true
			for isPageListGo {
				subjectIndexUrl := strings.ReplaceAll(paper.url, "p1", "p"+strconv.Itoa(current))
				subjectIndexDoc, err := htmlquery.LoadURL(subjectIndexUrl)
				if err != nil {
					fmt.Println(err)
					current = 1
					isPageListGo = false
					continue
				}
				liNodes := htmlquery.Find(subjectIndexDoc, `//div[@class="specification_list"]/ul/li`)
				if len(liNodes) <= 0 {
					fmt.Println(err)
					current = 1
					isPageListGo = false
					continue
				}
				for _, liNode := range liNodes {
					fmt.Println("============================================================================")
					fmt.Println("主题：", subjectsPapers.name, paper.name)
					fmt.Println("=======当前页URL", subjectIndexUrl, "========")

					viewUrl := "http://www.kaojuanxiazai.com" + htmlquery.InnerText(htmlquery.FindOne(liNode, `./div[@class="list_images_text"]/a/@href`))
					fmt.Println(viewUrl)

					viewDoc, _ := htmlquery.LoadURL(viewUrl)
					if viewDoc == nil {
						fmt.Println("获取试卷详情失败")
						continue
					}

					fileNameNode := htmlquery.FindOne(viewDoc, `//div[@class="article_title"]/div[@class="container"]/div[@class="title"]`)
					if fileNameNode == nil {
						fmt.Println("未能获取标题DIV")
						continue
					}
					fileName := htmlquery.InnerText(fileNameNode)
					fileName = strings.TrimSpace(fileName)
					fileName = strings.ReplaceAll(fileName, "<b>", "")
					fileName = strings.ReplaceAll(fileName, "</b>", "")
					fileName = strings.ReplaceAll(fileName, "/", "-")
					fileName = strings.ReplaceAll(fileName, ":", "-")
					fileName = strings.ReplaceAll(fileName, "：", "-")
					fileName = strings.ReplaceAll(fileName, "（", "(")
					fileName = strings.ReplaceAll(fileName, "）", ")")
					fileName = strings.ReplaceAll(fileName, "~", "")
					fileName = strings.ReplaceAll(fileName, ".", "")
					fileName = strings.ReplaceAll(fileName, ". ", "")
					fileName = strings.ReplaceAll(fileName, ".. ", "")
					fmt.Println(fileName)

					filePath := "../www.kaojuanxiazai.com/www.kaojuanxiazai.com/" + subjectsPapers.name + "/" + paper.name + "/" + fileName + ".doc"
					if _, err := os.Stat(filePath); err != nil {
						downLoadUrl := strings.ReplaceAll(viewUrl, "exam-", "exam/downloads/")
						fmt.Println(downLoadUrl)

						fmt.Println("=======开始下载" + strconv.Itoa(current) + "========")
						err = downloadKaoJuanXiaZai(downLoadUrl, viewUrl, filePath)
						if err != nil {
							fmt.Println(err)
							continue
						}
						fmt.Println("=======下载完成========")
						for i := 1; i <= KaoJuanXiaZaiNextDownloadSleep; i++ {
							time.Sleep(time.Second)
							fmt.Println("===========操作结束，暂停", KaoJuanXiaZaiNextDownloadSleep, "秒，倒计时", i, "秒===========")
						}
					}
				}
				current++
				isPageListGo = true
			}
		}
	}
}

func downloadKaoJuanXiaZai(attachmentUrl string, referer string, filePath string) error {
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
	if KaoJuanXiaZaiEnableHttpProxy {
		client = KaoJuanXiaZaiSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.kaojuanxiazai.com")
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
