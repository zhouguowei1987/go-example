package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
)

const (
	ShiJuan1EnableHttpProxy = false
	ShiJuan1HttpProxyUrl    = "27.42.168.46:55481"
)

func ShiJuan1SetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(ShiJuan1HttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type TestCategory struct {
	name string
	link []string
}

var AllTestCategory = []TestCategory{
	{
		name: "小学试卷",
		link: []string{
			"https://www.shijuan1.com/a/sjyw1/",
			"https://www.shijuan1.com/a/sjyw2/",
			"https://www.shijuan1.com/a/sjyw3/",
			"https://www.shijuan1.com/a/sjyw4/",
			"https://www.shijuan1.com/a/sjyw5/",
			"https://www.shijuan1.com/a/sjyw6/",

			"https://www.shijuan1.com/a/sjsx1/",
			"https://www.shijuan1.com/a/sjsx2/",
			"https://www.shijuan1.com/a/sjsx3/",
			"https://www.shijuan1.com/a/sjsx4/",
			"https://www.shijuan1.com/a/sjsx5/",
			"https://www.shijuan1.com/a/sjsx6/",

			"https://www.shijuan1.com/a/sjyy1/",
			"https://www.shijuan1.com/a/sjyy2/",
			"https://www.shijuan1.com/a/sjyy3/",
			"https://www.shijuan1.com/a/sjyy4/",
			"https://www.shijuan1.com/a/sjyy5/",
			"https://www.shijuan1.com/a/sjyy6/",
		},
	},
	{
		name: "中考试卷",
		link: []string{
			"https://www.shijuan1.com/a/sjyw7/",
			"https://www.shijuan1.com/a/sjyw8/",
			"https://www.shijuan1.com/a/sjyw9/",
			"https://www.shijuan1.com/a/sjywzk/",

			"https://www.shijuan1.com/a/sjsx7/",
			"https://www.shijuan1.com/a/sjsx8/",
			"https://www.shijuan1.com/a/sjsx9/",
			"https://www.shijuan1.com/a/sjsxzk/",

			"https://www.shijuan1.com/a/sjyy7/",
			"https://www.shijuan1.com/a/sjyy8/",
			"https://www.shijuan1.com/a/sjyy9/",
			"https://www.shijuan1.com/a/sjyyzk/",

			"https://www.shijuan1.com/a/sjwl8/",
			"https://www.shijuan1.com/a/sjwl9/",
			"https://www.shijuan1.com/a/sjwlzk/",

			"https://www.shijuan1.com/a/sjhx9/",
			"https://www.shijuan1.com/a/sjhxzk/",

			"https://www.shijuan1.com/a/sjzz7/",
			"https://www.shijuan1.com/a/sjzz8/",
			"https://www.shijuan1.com/a/sjzz9/",
			"https://www.shijuan1.com/a/sjzzzk/",

			"https://www.shijuan1.com/a/sjls7/",
			"https://www.shijuan1.com/a/sjls8/",
			"https://www.shijuan1.com/a/sjls9/",
			"https://www.shijuan1.com/a/sjlszk/",

			"https://www.shijuan1.com/a/sjdl7/",
			"https://www.shijuan1.com/a/sjdl8/",
			"https://www.shijuan1.com/a/sjdlzk/",

			"https://www.shijuan1.com/a/sjsw7/",
			"https://www.shijuan1.com/a/sjsw8/",
			"https://www.shijuan1.com/a/sjswzk/",
		},
	},
	{
		name: "高考试卷",
		link: []string{
			"https://www.shijuan1.com/a/sjywg1/",
			"https://www.shijuan1.com/a/sjywg2/",
			"https://www.shijuan1.com/a/sjywg3/",
			"https://www.shijuan1.com/a/sjywgk/",

			"https://www.shijuan1.com/a/sjsxg1/",
			"https://www.shijuan1.com/a/sjsxg2/",
			"https://www.shijuan1.com/a/sjsxg3/",
			"https://www.shijuan1.com/a/sjsxgk/",

			"https://www.shijuan1.com/a/sjyyg1/",
			"https://www.shijuan1.com/a/sjyyg2/",
			"https://www.shijuan1.com/a/sjyyg3/",
			"https://www.shijuan1.com/a/sjyygk/",

			"https://www.shijuan1.com/a/sjwlg1/",
			"https://www.shijuan1.com/a/sjwlg2/",
			"https://www.shijuan1.com/a/sjwlg3/",
			"https://www.shijuan1.com/a/sjwlgk/",

			"https://www.shijuan1.com/a/sjhxg1/",
			"https://www.shijuan1.com/a/sjhxg2/",
			"https://www.shijuan1.com/a/sjhxg3/",
			"https://www.shijuan1.com/a/sjhxgk/",

			"https://www.shijuan1.com/a/sjzzg1/",
			"https://www.shijuan1.com/a/sjzzg2/",
			"https://www.shijuan1.com/a/sjzzg3/",
			"https://www.shijuan1.com/a/sjzzgk/",

			"https://www.shijuan1.com/a/sjlsg1/",
			"https://www.shijuan1.com/a/sjlsg2/",
			"https://www.shijuan1.com/a/sjlsg3/",
			"https://www.shijuan1.com/a/sjlsgk/",

			"https://www.shijuan1.com/a/sjdlg1/",
			"https://www.shijuan1.com/a/sjdlg2/",
			"https://www.shijuan1.com/a/sjdlg3/",
			"https://www.shijuan1.com/a/sjdlgk/",

			"https://www.shijuan1.com/a/sjswg1/",
			"https://www.shijuan1.com/a/sjswg2/",
			"https://www.shijuan1.com/a/sjswg3/",
			"https://www.shijuan1.com/a/sjswgk/",
		},
	},
}

// ychEduSpider 获取第一试卷网文档
// @Title 获取第一试卷网文档
// @Description https://www.shijuan1.com/，获取第一试卷网文档
func main() {
	for _, testCategory := range AllTestCategory {
		page := 1
		for _, link := range testCategory.link {
			firstPaperDoc, _ := htmlquery.LoadURL(link)
			firstPaperPagesNodes := htmlquery.Find(firstPaperDoc, `//div[@class="dede_pages"]/ul[@class="pagelist"][1]/li`)

			var gradeId = 0
			if len(firstPaperPagesNodes) >= 3 {
				secondPageUrl := htmlquery.InnerText(htmlquery.FindOne(firstPaperPagesNodes[2], `./a/@href`))
				gradeId, _ = strconv.Atoi(strings.Split(secondPageUrl, "_")[1])
			}

			isPageListGo := true
			for isPageListGo {
				pageListUrl := fmt.Sprintf(link)
				if gradeId > 0 {
					pageListUrl = fmt.Sprintf(link+"list_"+strconv.Itoa(gradeId)+"_%d.html", page)
				}
				pageListDoc, _ := htmlquery.LoadURL(pageListUrl)
				tableTrNodes := htmlquery.Find(pageListDoc, `//div[@class="pleft"]/div[@class="listbox"]/ul[@class="c1"]/table/tbody/tr`)
				if len(tableTrNodes) >= 1 {
					for i, trNode := range tableTrNodes {
						if i == 0 {
							continue
						}
						fmt.Println("=================================================================================")
						fmt.Println(pageListUrl)

						title := htmlquery.InnerText(htmlquery.FindOne(trNode, `./td[1]`))
						title = strings.TrimSpace(title)
						title = strings.ReplaceAll(title, "/", "-")
						fmt.Println(title)

						dateText := htmlquery.InnerText(htmlquery.FindOne(trNode, `./td[6]`))
						fmt.Println(dateText)

						datePaper, _ := time.Parse("2006-01-02", dateText)
						dateStart, _ := time.Parse("2006-01-02", "2025-04-07")
						fmt.Println(dateStart)

						// 比较日期
						if datePaper.After(dateStart) == false {
							fmt.Println("日期在2025-04-07后，跳过")
							isPageListGo = false
							page = 1
							break
						}

						detailUrl := "https://www.shijuan1.com" + htmlquery.InnerText(htmlquery.FindOne(trNode, `./td[1]/a/@href`))
						detailDoc, _ := htmlquery.LoadURL(detailUrl)
						fmt.Println(detailUrl)

						filePath := "E:\\workspace\\www.shijuan1.com\\2025-04-07\\www.rar_shijuan1.com\\" + testCategory.name + "\\" + title + ".rar"
						if _, err := os.Stat(filePath); err != nil {
							downloadUrl := "https://www.shijuan1.com" + htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//ul[@class="downurllist"]/li/a/@href`))
							fmt.Println(downloadUrl)

							fmt.Println("=======开始下载" + title + "========")
							err := downloadShiJuan1(downloadUrl, detailUrl, filePath)
							if err != nil {
								fmt.Println(err)
								continue
							}
							fmt.Println("=======下载完成========")
						}
					}
					page++
				} else {
					isPageListGo = false
					page = 1
					break
				}

				if gradeId == 0 {
					isPageListGo = false
					page = 1
					break
				}
			}
		}
	}
}

func downloadShiJuan1(attachmentUrl string, referer string, filePath string) error {
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
	if ShiJuan1EnableHttpProxy {
		client = ShiJuan1SetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "__gads=ID=4fa9d5896738b30a-22c549a298d900b7:T=1675857919:RT=1675857919:S=ALNI_MaRt-4lkuhslkcGEubyXiLA8ppUFw; Hm_lvt_9400c877dfe1cf77b070ccf1be7b66af=1675857919,1676652604,1677259973,1677602207; __gpi=UID=00000bb7d488d640:T=1675857919:RT=1678152152:S=ALNI_MYmUZlhJSzKAlMesX2_56vkD0Vd_g; Hm_lpvt_9400c877dfe1cf77b070ccf1be7b66af=1678154476")
	req.Header.Set("Host", "www.shijuan1.com")
	req.Header.Set("If-Modified-Since", "Wed, 16 Nov 2022 20:09:29 GMT")
	req.Header.Set("If-None-Match", "W/\"63754379-2f74\"")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Not?A_Brand\";v=\"8\", \"Chromium\";v=\"108\", \"Google Chrome\";v=\"108\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
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
