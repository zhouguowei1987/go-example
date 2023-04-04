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
	"strings"
	"time"
)

const (
	XiaZai51EnableHttpProxy = false
	XiaZai51HttpProxyUrl    = "111.225.152.186:8089"
)

func XiaZai51SetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(XiaZai51HttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type XiaZai51Subject struct {
	name string
	url  string
}

var AllXiaZai51Subject = []XiaZai51Subject{
	{
		name: "合同模板",
		url:  "http://www.51xiazai.cn/sort/614/",
	},
}

// ychEduSpider 获取云奥科技文档
// @Title 获取云奥科技文档
// @Description http://www.51xiazai.cn/，获取云奥科技文档
func main() {
	for _, subject := range AllXiaZai51Subject {
		page := 1
		isPageListGo := true
		for isPageListGo {
			pageListUrl := subject.url
			if page > 1 {
				pageListUrl = fmt.Sprintf(subject.url+"%d/", page)
			}
			fmt.Println(pageListUrl)

			pageListDoc, err := htmlquery.LoadURL(pageListUrl)
			if err != nil {
				fmt.Println(err)
				break
			}

			dlNodes := htmlquery.Find(pageListDoc, `//div[@class="listCont"]/ul[@id="soft_list"]/li`)
			if len(dlNodes) >= 1 {
				for _, dlNode := range dlNodes {

					fmt.Println("=================================================================================")
					// 文档详情URL
					fileName := htmlquery.InnerText(htmlquery.FindOne(dlNode, `./div[@class="info"]/h4/a`))
					fmt.Println(fileName)

					detailUrl := htmlquery.InnerText(htmlquery.FindOne(dlNode, `./div[@class="info"]/h4/a/@href`))
					fmt.Println(detailUrl)

					detailDoc, _ := htmlquery.LoadURL(detailUrl)

					reg := regexp.MustCompile(`var downurlstr="(.*?)"`)
					//将所有null替换为空字符串
					detailDocText := htmlquery.InnerText(detailDoc)
					regFindStingMatch := reg.FindStringSubmatch(detailDocText)

					// 下载文档URL
					downLoadUrl := regFindStingMatch[1]
					fmt.Println(downLoadUrl)

					// 文件格式
					attachmentFormat := strings.Split(downLoadUrl, ".")

					filePath := "../www.51xiazai.cn/" + subject.name + "/"
					err = downloadXiaZai51(downLoadUrl, detailUrl, filePath, fileName+"."+attachmentFormat[len(attachmentFormat)-1])
					if err != nil {
						fmt.Println(err)
						continue
					}
					time.Sleep(time.Second * 1)
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
func downloadXiaZai51(attachmentUrl string, referer string, filePath string, fileName string) error {
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
	if XiaZai51EnableHttpProxy {
		client = XiaZai51SetHttpProxy()
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
	req.Header.Set("Host", "https://softforspeed.51xiazai.cn")
	req.Header.Set("Referer", referer)
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
