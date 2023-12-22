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
	"regexp"
	"strconv"
	"time"
)

var TRjlSengEnableHttpProxy = false
var TRjlSengHttpProxyUrl = ""
var TRjlSengHttpProxyUrlArr = make([]string, 0)

func TRjlSengHttpProxy() error {
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
					TRjlSengHttpProxyUrlArr = append(TRjlSengHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					TRjlSengHttpProxyUrlArr = append(TRjlSengHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func TRjlSengSetHttpProxy() (httpclient *http.Client) {
	if TRjlSengHttpProxyUrl == "" {
		if len(TRjlSengHttpProxyUrlArr) <= 0 {
			err := TRjlSengHttpProxy()
			if err != nil {
				TRjlSengSetHttpProxy()
			}
		}
		TRjlSengHttpProxyUrl = TRjlSengHttpProxyUrlArr[0]
		if len(TRjlSengHttpProxyUrlArr) >= 2 {
			TRjlSengHttpProxyUrlArr = TRjlSengHttpProxyUrlArr[1:]
		} else {
			TRjlSengHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(TRjlSengHttpProxyUrl)
	ProxyURL, _ := url.Parse(TRjlSengHttpProxyUrl)
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

type TRjlSengSubject struct {
	name   string
	papers []TRjlSengSubjectsPaper
}

type TRjlSengSubjectsPaper struct {
	name string
	url  string
}

var TRjlSengSubjectsPapers = []TRjlSengSubject{
	{
		name: "教材",
		papers: []TRjlSengSubjectsPaper{
			{
				name: "初中人教版",
				url:  "https://www.trjlseng.com/crjb/",
			},
			{
				name: "初中外研版",
				url:  "https://www.trjlseng.com/cwyb/",
			},
			{
				name: "初中牛津版",
				url:  "https://www.trjlseng.com/cnjb/",
			},
			{
				name: "高中人教版",
				url:  "https://www.trjlseng.com/grjb/",
			},
			{
				name: "高中外研版",
				url:  "https://www.trjlseng.com/gwyb/",
			},
			{
				name: "高中北师版",
				url:  "https://www.trjlseng.com/gbsb/",
			},
		},
	},

	{
		name: "试题",
		papers: []TRjlSengSubjectsPaper{
			{
				name: "初一试题",
				url:  "https://www.trjlseng.com/cyst/",
			},
			{
				name: "初二试题",
				url:  "https://www.trjlseng.com/cest/",
			},
			{
				name: "初三试题",
				url:  "https://www.trjlseng.com/csst/",
			},
			{
				name: "高一试题",
				url:  "https://www.trjlseng.com/gyst/",
			},
			{
				name: "高二试题",
				url:  "https://www.trjlseng.com/gest/",
			},
			{
				name: "高三试题",
				url:  "https://www.trjlseng.com/gsst/",
			},
		},
	},

	{
		name: "中考",
		papers: []TRjlSengSubjectsPaper{
			{
				name: "中考复习",
				url:  "https://www.trjlseng.com/zkfx/",
			},
			{
				name: "中考语法",
				url:  "https://www.trjlseng.com/zkyf/",
			},
			{
				name: "中考作文",
				url:  "https://www.trjlseng.com/zkzw/",
			},
			{
				name: "中考试题",
				url:  "https://www.trjlseng.com/zkst/",
			},
		},
	},

	{
		name: "高考",
		papers: []TRjlSengSubjectsPaper{
			{
				name: "高考复习",
				url:  "https://www.trjlseng.com/gkfx/",
			},
			{
				name: "高考语法",
				url:  "https://www.trjlseng.com/gkyf/",
			},
			{
				name: "高考作文",
				url:  "https://www.trjlseng.com/gkzw/",
			},
			{
				name: "高考试题",
				url:  "https://www.trjlseng.com/gkst/",
			},
		},
	},
}

var TRjlSengNextDownloadSleep = 2

// ychEduSpider 获取中学英语网文档
// @Title 获取中学英语网文档
// @Description https://www.trjlseng.com/，获取中学英语网文档
func main() {
	for _, subject := range TRjlSengSubjectsPapers {
		for _, paper := range subject.papers {
			current := 1
			isPageListGo := true
			// 计算最大页数
			paperIndexUrl := paper.url
			paperIndexDoc, err := htmlquery.LoadURL(paperIndexUrl)
			if err != nil {
				fmt.Println(err)
				continue
			}
			paperTotalNode := htmlquery.FindOne(paperIndexDoc, `//div[@class="yzm-container"]/div[@class="yzm-content-box yzm-main-left yzm-text-list"]/div[@id="page"]/span[@class="pageinfo"]/strong`)
			paperTotalText := htmlquery.InnerText(paperTotalNode)
			pagerTotal, err := strconv.Atoi(paperTotalText)
			if err != nil {
				fmt.Println(err)
				continue
			}
			paperMaxPages := (pagerTotal / 40) + 1

			for isPageListGo {
				paperListUrl := paper.url + fmt.Sprintf("list_%d.html", current)
				paperListDoc, err := htmlquery.LoadURL(paperListUrl)
				if err != nil {
					fmt.Println(err)
					current = 1
					isPageListGo = false
					continue
				}
				liNodes := htmlquery.Find(paperListDoc, `//div[@class="yzm-container"]/div[@class="yzm-content-box yzm-main-left yzm-text-list"]/ul/li`)
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

					title := htmlquery.InnerText(htmlquery.FindOne(liNode, `./a/@title`))
					fmt.Println(title)

					viewHref := "https://www.trjlseng.com" + htmlquery.InnerText(htmlquery.FindOne(liNode, `./a/@href`))
					fmt.Println(viewHref)

					// 查看是否有附件
					viewDoc, err := htmlquery.LoadURL(viewHref)
					if err != nil {
						fmt.Println(err)
						continue
					}

					regAttachmentViewUrl := regexp.MustCompile(`<a href="/uploads/ueditor/file/(.*?)" title="`)
					regAttachmentViewUrlMatch := regAttachmentViewUrl.FindAllSubmatch([]byte(htmlquery.OutputHTML(viewDoc, true)), -1)
					if len(regAttachmentViewUrlMatch) <= 0 {
						fmt.Println("没有附件，跳过")
						continue
					}
					attachmentUrl := "https://www.trjlseng.com/uploads/ueditor/file/" + string(regAttachmentViewUrlMatch[0][1])
					fmt.Println(attachmentUrl)

					filePath := "F:\\workspace\\www.rar_trjlseng.com\\" + title + ".rar"
					_, err = os.Stat(filePath)
					if err != nil {

						fmt.Println("=======开始下载"+strconv.Itoa(current)+"-", paperMaxPages, "========")
						err = downloadTRjlSeng(attachmentUrl, viewHref, filePath)
						if err != nil {
							fmt.Println(err)
							continue
						}
						fmt.Println("=======下载完成========")
						for i := 1; i <= TRjlSengNextDownloadSleep; i++ {
							time.Sleep(time.Second)
							fmt.Println("===========操作结束，暂停", TRjlSengNextDownloadSleep, "秒，倒计时", i, "秒===========")
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

func downloadTRjlSeng(attachmentUrl string, referer string, filePath string) error {
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
	if TRjlSengEnableHttpProxy {
		client = TRjlSengSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "__51cke__=; _gid=GA1.2.944241540.1703144883; PHPSESSID=22d6kpnbct0v3iopp7qctbk081; __tins__21123451=%7B%22sid%22%3A%201703148480266%2C%20%22vd%22%3A%207%2C%20%22expires%22%3A%201703151933380%7D; __51laig__=27; _gat=1; _ga_34B604LFFQ=GS1.1.1703148480.2.1.1703150135.57.0.0; _ga=GA1.1.1587097358.1703144883")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.trjlseng.com")
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
