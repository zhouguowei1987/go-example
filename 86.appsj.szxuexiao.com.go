package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	SZXueXiaoEnableHttpProxy = false
	SZXueXiaoHttpProxyUrl    = "111.225.152.186:8089"
)

func SZXueXiaoSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(SZXueXiaoHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type Grade86 struct {
	name string
	url  string
}

var AllGrade = []Grade86{
	{
		name: "一年级",
		url:  "https://appsj.szxuexiao.com/yinianji/index.html",
	},
	{
		name: "二年级",
		url:  "https://appsj.szxuexiao.com/ernianji/index.html",
	},
	{
		name: "三年级",
		url:  "https://appsj.szxuexiao.com/sannianji/index.html",
	},
	{
		name: "四年级",
		url:  "https://appsj.szxuexiao.com/sinianji/index.html",
	},
	{
		name: "五年级",
		url:  "https://appsj.szxuexiao.com/wunianji/index.html",
	},
	{
		name: "六年级",
		url:  "https://appsj.szxuexiao.com/liunianji/index.html",
	},
	{
		name: "语文",
		url:  "https://appsj.szxuexiao.com/yuwen/index.html",
	},
	{
		name: "数学",
		url:  "https://appsj.szxuexiao.com/shuxue/index.html",
	},
	{
		name: "英语",
		url:  "https://appsj.szxuexiao.com/yingyu/index.html",
	},
	{
		name: "科学",
		url:  "https://appsj.szxuexiao.com/kexue/index.html",
	},
}

// ychEduSpider 获取名校教研文档
// @Title 获取名校教研文档
// @Description https://appsj.szxuexiao.com/，获取名校教研文档
func main() {
	for _, grade := range AllGrade {
		page := 1
		isPageListGo := true
		for isPageListGo {

			pageListUrl := strings.Replace(grade.url, "index", "index_"+strconv.Itoa(page), -1)
			if page == 1 {
				pageListUrl = grade.url
			}
			pageListDoc, err := htmlquery.LoadURL(pageListUrl)
			if err != nil {
				fmt.Println(err)
				break
			}
			divNodes := htmlquery.Find(pageListDoc, `//div[@class="list-group list-group-flush"]/a[@class="list-group-item"]`)
			if len(divNodes) >= 1 {
				for _, listNode := range divNodes {

					detailUrl := "https://appsj.szxuexiao.com/" + htmlquery.SelectAttr(listNode, "href")
					detailDoc, err := htmlquery.LoadURL(detailUrl)
					if err != nil {
						continue
					}

					// 查看是否有《点击下载文档》连接
					downloadPan := htmlquery.FindOne(detailDoc, `//div[@class="entry-content"]/a[@class="btn btn-primary btn-lg active"]/@href`)
					if downloadPan == nil {
						isPageListGo = false
						page = 1
						break
					}
					downloadPanUrl := htmlquery.InnerText(downloadPan)
					downloadPwd := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="entry-content"]/font`))

					baiduPanDownloadUrl := downloadPanUrl + "?pwd=" + downloadPwd
					_, err = url.ParseRequestURI(baiduPanDownloadUrl)
					if err == nil {
						fmt.Println(baiduPanDownloadUrl)
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
