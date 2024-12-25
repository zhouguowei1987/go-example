package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"math/rand"
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
	TopEduEnableHttpProxy = false
	TopEduHttpProxyUrl    = "111.225.152.186:8089"
)

func TopEduSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(TopEduHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type TopEduSubject struct {
	name string
	url  string
	tp   string
}

var AllTopEduSubject = []TopEduSubject{
	{
		name: "中考真题",
		url:  "http://topedu.ybep.com.cn/new_x/project/really_test.php?tp=z",
		tp:   "z",
	},
	{
		name: "高考真题",
		url:  "http://topedu.ybep.com.cn/new_x/project/really_test.php?tp=g",
		tp:   "g",
	},
}

// ychEduSpider 获取鼎尖资源网文档
// @Title 获取鼎尖资源网文档
// @Description http://topedu.ybep.com.cn/，获取鼎尖资源网文档
func main() {
	for _, subject := range AllTopEduSubject {
		page := 1
		totalPage := 1
		for page <= totalPage {
			pageListUrl := subject.url
			pageListUrl = fmt.Sprintf(subject.url+"&page=%d", page)
			fmt.Println(pageListUrl)

			pageListDoc, err := htmlquery.LoadURL(pageListUrl)
			if err != nil {
				fmt.Println(err)
				break
			}
			if page == 1 {
				// 获取总页数
				pageNodes := htmlquery.Find(pageListDoc, `//ul[@class="pagination"]/li`)
				totalPageNodeXpath := "//ul[@class=\"pagination\"]/li[" + strconv.Itoa(len(pageNodes)) + "]"
				totalPageNode := htmlquery.InnerText(htmlquery.FindOne(pageListDoc, totalPageNodeXpath))
				totalPage, _ = strconv.Atoi(totalPageNode)
			}

			dlNodes := htmlquery.Find(pageListDoc, `//ul[@class="cbottoms"]/li`)
			if len(dlNodes) >= 1 {
				for _, dlNode := range dlNodes {

					fmt.Println("=================================================================================")
					// 文档详情URL
					fileName := htmlquery.InnerText(htmlquery.FindOne(dlNode, `./span/em/a`))
					fileName = strings.ReplaceAll(fileName, "_", "")
					fileName = strings.ReplaceAll(fileName, " ", "")
					fmt.Println(fileName)
					if !strings.Contains(fileName, "doc") {
						continue
					}

					detailUrl := htmlquery.InnerText(htmlquery.FindOne(dlNode, `./span/em/a/@href`))
					detailUrl = "http://topedu.ybep.com.cn/new_x/project/" + detailUrl
					fmt.Println(detailUrl)
					//解析url 并保证没有错误
					u, err := url.Parse(detailUrl)
					if err != nil {
						fmt.Println(err)
						continue
					}
					urlParam, err := url.ParseQuery(u.RawQuery)
					if err != nil {
						fmt.Println(err)
						continue
					}
					attachmentUrl := fmt.Sprintf("http://topedu.ybep.com.cn/new_x/project/clouddown.php?pg=really&id=%s&tp=%s&opt=0&n=%s", urlParam.Get("id"), subject.tp, urlParam.Get("n"))
					fmt.Println(attachmentUrl)

					filePath := "../topedu.ybep.com.cn/" + subject.name + "/" + fileName
					if _, err := os.Stat(filePath); err != nil {
						fmt.Println("=======开始下载========")
						err = downloadTopEdu(attachmentUrl, detailUrl, filePath)
						if err != nil {
							fmt.Println(err)
							continue
						}
						fmt.Println("=======开始完成========")
						DownLoadTopEduTimeSleep := rand.Intn(5)
						for i := 1; i <= DownLoadTopEduTimeSleep; i++ {
							time.Sleep(time.Second)
							fmt.Println("page="+strconv.Itoa(page)+"===========下载", fileName, "成功，暂停", DownLoadTopEduTimeSleep, "秒，倒计时", i, "秒===========")
						}
					}
				}
				page++
			} else {
				page = 1
				break
			}
		}
	}
}

func downloadTopEdu(attachmentUrl string, referer string, filePath string) error {
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
	if TopEduEnableHttpProxy {
		client = TopEduSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "__51vcke__JnUuVbB1pNBuOI9J=7a1ac260-d305-5b5d-b8ff-d2650e75d779; __51vuft__JnUuVbB1pNBuOI9J=1716949675665; PHPSESSID=cvbi2rkq2ua8a3l3gcqfehp11l; __51uvsct__JnUuVbB1pNBuOI9J=3; __vtins__JnUuVbB1pNBuOI9J=%7B%22sid%22%3A%20%226025d880-da35-5d23-8302-7bafb55f6519%22%2C%20%22vd%22%3A%2036%2C%20%22stt%22%3A%20495709%2C%20%22dr%22%3A%2044058%2C%20%22expires%22%3A%201735098428534%2C%20%22ct%22%3A%201735096628534%7D")
	req.Header.Set("Host", "opedu.ybep.com.cn")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", referer)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
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
