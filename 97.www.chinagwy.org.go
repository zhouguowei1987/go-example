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
	ChinaGwyEnableHttpProxy = false
	ChinaGwyHttpProxyUrl    = "111.225.152.186:8089"
)

func ChinaGwySetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(ChinaGwyHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

// ychEduSpider 获取公考资讯网文档
// @Title 获取公考资讯网文档
// @Description https://www.chinagwy.org/，获取公考资讯网文档
func main() {
	maxPage := 50
	page := 1
	isPageListGo := true
	for isPageListGo {
		pageListUrl := fmt.Sprintf("https://www.chinagwy.org/html/stzx/7_%d.html", page)
		fmt.Println(pageListUrl)

		pageListDoc, err := htmlquery.LoadURL(pageListUrl)
		if err != nil {
			fmt.Println(err)
			break
		}

		dlNodes := htmlquery.Find(pageListDoc, `//div[@class="con"]/ul[@class="list01"]/li[not(@class="line")]`)
		if len(dlNodes) >= 1 {
			for _, dlNode := range dlNodes {

				fmt.Println("=================================================================================")
				// 文档详情URL
				fileName := htmlquery.InnerText(htmlquery.FindOne(dlNode, `./a[2]`))
				fmt.Println(fileName)

				detailUrl := htmlquery.InnerText(htmlquery.FindOne(dlNode, `./a[2]/@href`))
				fmt.Println(detailUrl)

				detailDoc, _ := htmlquery.LoadURL(detailUrl)
				detailDocText := htmlquery.OutputHTML(detailDoc, true)

				reg := regexp.MustCompile(`<a href="http://www.chinagwy.org/files/(.*?).pdf" target="_blank">(.*?).pdf</a>`)
				regFindStingMatch := reg.FindStringSubmatch(detailDocText)

				if len(regFindStingMatch) != 3 {
					continue
				}
				aPdfFileName := regFindStingMatch[1]
				if strings.Index(regFindStingMatch[2], "答案") == -1 {
					continue
				}

				// 下载文档URL
				downLoadUrl := fmt.Sprintf("http://www.chinagwy.org/files/%s.pdf", aPdfFileName)
				fmt.Println(downLoadUrl)

				// 文件格式
				attachmentFormat := strings.Split(downLoadUrl, ".")
				filePath := "../www.chinagwy.org/" + fileName + "(含答案)." + attachmentFormat[len(attachmentFormat)-1]
				if _, err := os.Stat(filePath); err != nil {
					fmt.Println("=======开始下载========")
					err = downloadChinaGwy(downLoadUrl, detailUrl, filePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("=======开始完成========")
				}
				time.Sleep(time.Second * 1)
			}
			page++
			if page > maxPage {
				isPageListGo = false
				break
			}
		} else {
			isPageListGo = false
			break
		}
	}
}
func downloadChinaGwy(attachmentUrl string, referer string, filePath string) error {
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
	if ChinaGwyEnableHttpProxy {
		client = ChinaGwySetHttpProxy()
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
	req.Header.Set("Host", "https://www.chinagwy.org")
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
