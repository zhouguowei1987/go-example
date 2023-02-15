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
	"time"
)

const (
	DiYiFanWenEnableHttpProxy = false
	DiYiFanWenHttpProxyUrl    = "111.225.152.186:8089"
)

func DiYiFanWenSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(DiYiFanWenHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type Category struct {
	name string
	url  string
}

var AllCategory = []Category{
	{
		name: "合同范本",
		url:  "https://www.diyifanwen.com/fanwen/hetongfanwen/",
	},
	{
		name: "工作计划",
		url:  "https://www.diyifanwen.com/fanwen/gongzuojihua/",
	},
	{
		name: "工作总结",
		url:  "https://www.diyifanwen.com/fanwen/gongzuozongjie/",
	},
	{
		name: "心得体会",
		url:  "https://www.diyifanwen.com/fanwen/xindetihui2/",
	},
	{
		name: "演讲稿",
		url:  "https://www.diyifanwen.com/yanjianggao/",
	},
}

// ychEduSpider 获取第一范文网文档
// @Title 获取第一范文网文档
// @Description https://www.diyifanwen.com/，获取第一范文网文档
func main() {
	for _, category := range AllCategory {
		childCategoryDoc, err := htmlquery.LoadURL(category.url)
		if err != nil {
			fmt.Println(err)
			break
		}
		childCategoryNodes := htmlquery.Find(childCategoryDoc, `//div[@class="childlist-data"]/dl[@class="SList"]`)
		if len(childCategoryNodes) >= 1 {
			for _, childCategoryNode := range childCategoryNodes {
				childCategoryListName := htmlquery.InnerText(htmlquery.FindOne(childCategoryNode, `./dt/a`))
				childCategoryListUrl := "https:" + htmlquery.InnerText(htmlquery.FindOne(childCategoryNode, `./dt/a/@href`))

				page := 1
				isPageListGo := true
				for isPageListGo {
					pageListUrl := fmt.Sprintf(childCategoryListUrl+"list_"+strconv.Itoa(page)+"_%d.html", page)
					if page == 1 {
						pageListUrl = childCategoryListUrl
					}
					pageListDoc, err := htmlquery.LoadURL(pageListUrl)
					if err != nil {
						fmt.Println(err)
						break
					}
					divNodes := htmlquery.Find(pageListDoc, `//div[@class="alllist-data"]/ul/li]`)
					if len(divNodes) >= 1 {
						for _, listNode := range divNodes {

							fmt.Println("=================================================================================")
							// 文档详情URL
							fileName := htmlquery.InnerText(htmlquery.FindOne(listNode, `./a`))
							fmt.Println(fileName)
							detailUrl := "https:" + htmlquery.InnerText(htmlquery.FindOne(listNode, `./a/@href`))
							fmt.Println(detailUrl)
							// 下载预览URL
							downDetailUrl := "https://s.diyifanwen.com/down/down.asp?url=" + detailUrl + "&obid=fanwen"
							// 下载文档URL
							downLoadUrl := "https://s.diyifanwen.com/down/doc.asp?id=" + detailUrl
							filePath := "../www.diyifanwen.com/" + category.name + "/" + childCategoryListName + "/"
							err = downloadDiYiFanWen(downLoadUrl, downDetailUrl, filePath, fileName+".doc")
							time.Sleep(time.Second * 15)
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
}
func downloadDiYiFanWen(attachmentUrl string, referer string, filePath string, fileName string) error {
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
	if DiYiFanWenEnableHttpProxy {
		client = DiYiFanWenSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "SetCookieTF=1; testcookie=yes; Hm_lvt_3a5e11b41af918022c823a8041a34e78=1676440526; _gid=GA1.2.898443237.1676440527; DYFWUID=1676440529075qpgonrr4i85fx4m; __bid_n=1840cffdf58e5387354207; FPTOKEN=s3SDdJ00cklGHOZzBBnYcBnfVnTCbNCvCfd+7/0HQ+Fuddgujd/9/zjvzbKggY8WQiqpHnWq2r3ta08I5vbUsboB1oORD/Um4TdsNXhBrYPyim/K8sMvwkjLfbZlniVPN7sXxPaJv0Nq6SfOtn3upRxTDZMh/ae9YR/GXDg49FTyT0dl/7CWa4kJetDFg/ysAuwJj5gWYnXVm6pMSDoHoE1EDySmHU5Z9nnI798Hog5K4v5wMyMIPGfy80bdcPTp3tmxDGXWNmgVdfAOEY1NqdbelHGHgf7tnXI3m3LSGrftuaqpBLdaR52WmMz2Jb5K8BOPgoagGWB47urF9xhoNskY0z2VEN74aV86EOfpZR50Q1NRoxtk2B1NahAA+pPpR/mc+ruLCyXVndboZTM18A==|OWQdsasbF5948RgIGgMgpz1WS6AuCWU8H8atBRS7tl4=|10|bc84415d76a9eea3cccfa5fd50cf08aa; ASPSESSIONIDCUSADTRQ=CLDJJJMBALGGPABCEPKNOCOJ; AppealCount=4; Hm_lpvt_3a5e11b41af918022c823a8041a34e78=1676447365; _ga_34B604LFFQ=GS1.1.1676440526.1.1.1676447365.56.0.0; _ga=GA1.1.1070625585.1676440527")
	req.Header.Set("Host", "s.diyifanwen.com")
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
	out, err := os.Create(filePath + fileName)
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
