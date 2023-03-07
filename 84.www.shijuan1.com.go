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
	{
		name: "语文试卷",
		category: []Grade{
			{name: "一年级", url: "https://www.shijuan1.com/a/sjyw1/"},
			{name: "二年级", url: "https://www.shijuan1.com/a/sjyw2/"},
			{name: "三年级", url: "https://www.shijuan1.com/a/sjyw3/"},
			{name: "四年级", url: "https://www.shijuan1.com/a/sjyw4/"},
			{name: "五年级", url: "https://www.shijuan1.com/a/sjyw5/"},
			{name: "六年级", url: "https://www.shijuan1.com/a/sjyw6/"},
			{name: "七年级", url: "https://www.shijuan1.com/a/sjyw7/"},
			{name: "八年级", url: "https://www.shijuan1.com/a/sjyw8/"},
			{name: "九年级", url: "https://www.shijuan1.com/a/sjyw9/"},
			{name: "中考试卷", url: "https://www.shijuan1.com/a/sjywzk/"},
			{name: "高一", url: "https://www.shijuan1.com/a/sjywg1/"},
			{name: "高二", url: "https://www.shijuan1.com/a/sjywg2/"},
			{name: "高三", url: "https://www.shijuan1.com/a/sjywg3/"},
			{name: "高考试卷", url: "https://www.shijuan1.com/a/sjywgk/"},
		},
	},
	{
		name: "数学试卷",
		category: []Grade{
			{name: "一年级", url: "https://www.shijuan1.com/a/sjsx1/"},
			{name: "二年级", url: "https://www.shijuan1.com/a/sjsx2/"},
			{name: "三年级", url: "https://www.shijuan1.com/a/sjsx3/"},
			{name: "四年级", url: "https://www.shijuan1.com/a/sjsx4/"},
			{name: "五年级", url: "https://www.shijuan1.com/a/sjsx5/"},
			{name: "六年级", url: "https://www.shijuan1.com/a/sjsx6/"},
			{name: "七年级", url: "https://www.shijuan1.com/a/sjsx7/"},
			{name: "八年级", url: "https://www.shijuan1.com/a/sjsx8/"},
			{name: "九年级", url: "https://www.shijuan1.com/a/sjsx9/"},
			{name: "中考试卷", url: "https://www.shijuan1.com/a/sjsxzk/"},
			{name: "高一", url: "https://www.shijuan1.com/a/sjsxg1/"},
			{name: "高二", url: "https://www.shijuan1.com/a/sjsxg2/"},
			{name: "高三", url: "https://www.shijuan1.com/a/sjsxg3/"},
			{name: "高考试卷", url: "https://www.shijuan1.com/a/sjsxgk/"},
		},
	},
	{
		name: "英语试卷",
		category: []Grade{
			{name: "一年级", url: "https://www.shijuan1.com/a/sjyy1/"},
			{name: "二年级", url: "https://www.shijuan1.com/a/sjyy2/"},
			{name: "三年级", url: "https://www.shijuan1.com/a/sjyy3/"},
			{name: "四年级", url: "https://www.shijuan1.com/a/sjyy4/"},
			{name: "五年级", url: "https://www.shijuan1.com/a/sjyy5/"},
			{name: "六年级", url: "https://www.shijuan1.com/a/sjyy6/"},
			{name: "七年级", url: "https://www.shijuan1.com/a/sjyy7/"},
			{name: "八年级", url: "https://www.shijuan1.com/a/sjyy8/"},
			{name: "九年级", url: "https://www.shijuan1.com/a/sjyy9/"},
			{name: "中考试卷", url: "https://www.shijuan1.com/a/sjyyzk/"},
			{name: "高一", url: "https://www.shijuan1.com/a/sjyyg1/"},
			{name: "高二", url: "https://www.shijuan1.com/a/sjyyg2/"},
			{name: "高三", url: "https://www.shijuan1.com/a/sjyyg3/"},
			{name: "高考试卷", url: "https://www.shijuan1.com/a/sjyygk/"},
		},
	},
	{
		name: "物理试卷",
		category: []Grade{
			{name: "八年级", url: "https://www.shijuan1.com/a/sjwl8/"},
			{name: "九年级", url: "https://www.shijuan1.com/a/sjwl9/"},
			{name: "中考试卷", url: "https://www.shijuan1.com/a/sjwlzk/"},
			{name: "高一", url: "https://www.shijuan1.com/a/sjwlg1/"},
			{name: "高二", url: "https://www.shijuan1.com/a/sjwlg2/"},
			{name: "高三", url: "https://www.shijuan1.com/a/sjwlg3/"},
			{name: "高考试卷", url: "https://www.shijuan1.com/a/sjwlgk/"},
		},
	},
	{
		name: "化学试卷",
		category: []Grade{
			{name: "九年级", url: "https://www.shijuan1.com/a/sjhx9/"},
			{name: "中考试卷", url: "https://www.shijuan1.com/a/sjhxzk/"},
			{name: "高一", url: "https://www.shijuan1.com/a/sjhxg1/"},
			{name: "高二", url: "https://www.shijuan1.com/a/sjhxg2/"},
			{name: "高三", url: "https://www.shijuan1.com/a/sjhxg3/"},
			{name: "高考试卷", url: "https://www.shijuan1.com/a/sjhxgk/"},
		},
	},
	{
		name: "政治试卷",
		category: []Grade{
			{name: "七年级", url: "https://www.shijuan1.com/a/sjzz7/"},
			{name: "八年级", url: "https://www.shijuan1.com/a/sjzz8/"},
			{name: "九年级", url: "https://www.shijuan1.com/a/sjzz9/"},
			{name: "中考试卷", url: "https://www.shijuan1.com/a/sjzzzk/"},
			{name: "高一", url: "https://www.shijuan1.com/a/sjzzg1/"},
			{name: "高二", url: "https://www.shijuan1.com/a/sjzzg2/"},
			{name: "高三", url: "https://www.shijuan1.com/a/sjzzg3/"},
			{name: "高考试卷", url: "https://www.shijuan1.com/a/sjzzgk/"},
		},
	},
	{
		name: "历史试卷",
		category: []Grade{
			{name: "七年级", url: "https://www.shijuan1.com/a/sjls7/"},
			{name: "八年级", url: "https://www.shijuan1.com/a/sjls8/"},
			{name: "九年级", url: "https://www.shijuan1.com/a/sjls9/"},
			{name: "中考试卷", url: "https://www.shijuan1.com/a/sjlszk/"},
			{name: "高一", url: "https://www.shijuan1.com/a/sjlsg1/"},
			{name: "高二", url: "https://www.shijuan1.com/a/sjlsg2/"},
			{name: "高三", url: "https://www.shijuan1.com/a/sjlsg3/"},
			{name: "高考试卷", url: "https://www.shijuan1.com/a/sjlsgk/"},
		},
	},
	{
		name: "地理试卷",
		category: []Grade{
			{name: "七年级", url: "https://www.shijuan1.com/a/sjdl7/"},
			{name: "八年级", url: "https://www.shijuan1.com/a/sjdl8/"},
			{name: "中考试卷", url: "https://www.shijuan1.com/a/sjdlzk/"},
			{name: "高一", url: "https://www.shijuan1.com/a/sjdlg1/"},
			{name: "高二", url: "https://www.shijuan1.com/a/sjdlg2/"},
			{name: "高三", url: "https://www.shijuan1.com/a/sjdlg3/"},
			{name: "高考试卷", url: "https://www.shijuan1.com/a/sjdlgk/"},
		},
	},
	{
		name: "生物试卷",
		category: []Grade{
			{name: "七年级", url: "https://www.shijuan1.com/a/sjsw7/"},
			{name: "八年级", url: "https://www.shijuan1.com/a/sjsw8/"},
			{name: "中考试卷", url: "https://www.shijuan1.com/a/sjswzk/"},
			{name: "高一", url: "https://www.shijuan1.com/a/sjswg1/"},
			{name: "高二", url: "https://www.shijuan1.com/a/sjswg2/"},
			{name: "高三", url: "https://www.shijuan1.com/a/sjswg3/"},
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

			var gradeId = 0
			if len(firstPaperPagesNodes) >= 3 {
				secondPageUrl := htmlquery.InnerText(htmlquery.FindOne(firstPaperPagesNodes[2], `./a/@href`))
				gradeId, _ = strconv.Atoi(strings.Split(secondPageUrl, "_")[1])
			}

			isPageListGo := true
			for isPageListGo {
				pageListUrl := fmt.Sprintf(grade.url)
				if gradeId > 0 {
					pageListUrl = fmt.Sprintf(grade.url+"list_"+strconv.Itoa(gradeId)+"_%d.html", page)
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

				if gradeId == 0 {
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
	fileFullPath := filePath + title + ".rar"
	_, err = os.Stat(fileFullPath)
	if err != nil {
		//文件不存在
		out, err := os.Create(fileFullPath)
		if err != nil {
			return err
		}
		defer out.Close()

		// 然后将响应流和文件流对接起来
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}
