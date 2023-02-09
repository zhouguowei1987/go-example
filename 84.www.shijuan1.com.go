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

type Grade struct {
	name string
	url  string
}

type TestCategory struct {
	name     string
	category []Grade
}

var AllTestCategory = []TestCategory{
	//{
	//	name: "语文试卷",
	//	category: []Grade{
	//		{name: "中考试卷", url: "https://www.shijuan1.com/a/sjywzk/"},
	//		{name: "高考试卷", url: "https://www.shijuan1.com/a/sjywgk/"},
	//	},
	//},
	//{
	//	name: "数学试卷",
	//	category: []Grade{
	//		{name: "中考试卷", url: "https://www.shijuan1.com/a/sjsxzk/"},
	//		{name: "高考试卷", url: "https://www.shijuan1.com/a/sjsxgk/"},
	//	},
	//},
	//{
	//	name: "英语试卷",
	//	category: []Grade{
	//		{name: "中考试卷", url: "https://www.shijuan1.com/a/sjyyzk/"},
	//		{name: "高考试卷", url: "https://www.shijuan1.com/a/sjyygk/"},
	//	},
	//},
	{
		name: "物理试卷",
		category: []Grade{
			{name: "中考试卷", url: "https://www.shijuan1.com/a/sjwlzk/"},
			{name: "高考试卷", url: "https://www.shijuan1.com/a/sjwlgk/"},
		},
	},
	{
		name: "化学试卷",
		category: []Grade{
			{name: "中考试卷", url: "https://www.shijuan1.com/a/sjhxzk/"},
			{name: "高考试卷", url: "https://www.shijuan1.com/a/sjhxgk/"},
		},
	},
	{
		name: "政治试卷",
		category: []Grade{
			{name: "中考试卷", url: "https://www.shijuan1.com/a/sjzzzk/"},
			{name: "高考试卷", url: "https://www.shijuan1.com/a/sjzzgk/"},
		},
	},
	{
		name: "历史试卷",
		category: []Grade{
			{name: "中考试卷", url: "https://www.shijuan1.com/a/sjlszk/"},
			{name: "高考试卷", url: "https://www.shijuan1.com/a/sjlsgk/"},
		},
	},
	{
		name: "地理试卷",
		category: []Grade{
			{name: "中考试卷", url: "https://www.shijuan1.com/a/sjdlzk/"},
			{name: "高考试卷", url: "https://www.shijuan1.com/a/sjdlgk/"},
		},
	},
	{
		name: "生物试卷",
		category: []Grade{
			{name: "中考试卷", url: "https://www.shijuan1.com/a/sjswzk/"},
			{name: "高考试卷", url: "https://www.shijuan1.com/a/sjswgk/"},
		},
	},
}

// ychEduSpider 获取第一试卷网文档
// @Title 获取第一试卷网文档
// @Description https://www.shijian1.com/，获取第一试卷网文档
func main() {
	for _, testCategory := range AllTestCategory {
		page := 1
		for _, grade := range testCategory.category {
			firstPaperDoc, _ := htmlquery.LoadURL(grade.url)
			firstPaperPagesNodes := htmlquery.Find(firstPaperDoc, `//div[@class="dede_pages"]/ul[@class="pagelist"][1]/li`)
			secondPageUrl := htmlquery.InnerText(htmlquery.FindOne(firstPaperPagesNodes[2], `./a/@href`))
			gradeId := strings.Split(secondPageUrl, "_")[1]

			isPageListGo := true
			for isPageListGo {
				pageListUrl := fmt.Sprintf(grade.url+"list_"+gradeId+"_%d.html", page)
				pageListDoc, _ := htmlquery.LoadURL(pageListUrl)
				tableTrNodes := htmlquery.Find(pageListDoc, `//div[@class="pleft"]/div[@class="listbox"]/ul[@class="c1"]/table/tbody/tr`)
				if len(tableTrNodes) >= 1 {
					for i, trNode := range tableTrNodes {
						if i == 0 {
							continue
						}
						fmt.Println("=================================================================================")
						fmt.Println(pageListUrl)

						detailUrl := "https://www.shijuan1.com" + htmlquery.InnerText(htmlquery.FindOne(trNode, `./td[1]/a/@href`))
						detailDoc, _ := htmlquery.LoadURL(detailUrl)
						fmt.Println(detailUrl)

						title := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="pleft"]/div[@class="viewbox"]/div[@class="title"]/h2`))
						title = strings.ReplaceAll(title, "/", "-")
						title = strings.ReplaceAll(title, " ", "")
						fmt.Println(title)

						updateDate := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="pleft"]/div[@class="viewbox"]/div[@class="infolist"]/span[6]`))
						yearMonthDay := strings.Split(updateDate, "-")
						if year, _ := strconv.Atoi(yearMonthDay[0]); year < 2020 {
							isPageListGo = false
							page = 1
							break
						}

						downloadUrl := "https://www.shijuan1.com" + htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//ul[@class="downurllist"]/li/a/@href`))
						fmt.Println(downloadUrl)

						filePath := "../www.shijuan1.com/" + testCategory.name + "/" + grade.name + "/"

						err := downloadShiJuan1(downloadUrl, detailUrl, filePath, title)
						if err != nil {
							fmt.Println(err)
							continue
						}
					}
					page++
				} else {
					isPageListGo = false
					page = 1
					break
				}
			}
		}
	}
}

func downloadShiJuan1(attachmentUrl string, referer string, filePath string, title string) error {
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
	req.Header.Set("Cookie", "Hm_lvt_9400c877dfe1cf77b070ccf1be7b66af=1675857919; __gads=ID=4fa9d5896738b30a-22c549a298d900b7:T=1675857919:RT=1675857919:S=ALNI_MaRt-4lkuhslkcGEubyXiLA8ppUFw; __gpi=UID=00000bb7d488d640:T=1675857919:RT=1675857919:S=ALNI_MYmUZlhJSzKAlMesX2_56vkD0Vd_g; Hm_lpvt_9400c877dfe1cf77b070ccf1be7b66af=1675865606")
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
	out, err := os.Create(filePath + title + ".rar")
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
