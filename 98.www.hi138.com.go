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
	Hi138EnableHttpProxy = false
	Hi138HttpProxyUrl    = "111.225.152.186:8089"
)

func Hi138SetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(Hi138HttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

// ychEduSpider 获取免费论文下载中心文档
// @Title 获取免费论文下载中心文档
// @Description http://www.hi138.com/，获取免费论文下载中心文档
func main() {
	page := 1
	isPageListGo := true
	indexDoc, err := htmlquery.LoadURL("http://www.hi138.com")
	if err != nil {
		fmt.Println(err)
	}
	mainNodes := htmlquery.Find(indexDoc, `//div[@class="copyright"]/div[@class="main"]/div[@class="bottom"]/div[@class="main"]`)
	if len(mainNodes) <= 0 {
		return
	}

	for _, mainNode := range mainNodes {
		dlNodes := htmlquery.Find(mainNode, `./div[@class="fr"]/dl`)
		if len(dlNodes) <= 0 {
			break
		}
		for _, dlNode := range dlNodes {
			bigSubject := htmlquery.FindOne(dlNode, `./dt/a[@class="big"]`)
			bigSubjectName := htmlquery.InnerText(bigSubject)
			smallSubjects := htmlquery.Find(dlNode, `./dt/a[not(@class="big")]`)
			if len(smallSubjects) <= 0 {
				break
			}
			for _, smallSubject := range smallSubjects {
				smallSubjectName := htmlquery.InnerText(smallSubject)
				smallSubjectUrl := htmlquery.InnerText(htmlquery.FindOne(smallSubject, `./@href`))
				for isPageListGo {
					smallSubjectListUrl := fmt.Sprintf("http://www.hi138.com"+smallSubjectUrl+"%d/", page)
					fmt.Println(smallSubjectListUrl)

					pageListDoc, err := htmlquery.LoadURL(smallSubjectListUrl)
					if err != nil {
						fmt.Println(err)
						break
					}

					liNodes := htmlquery.Find(pageListDoc, `//div[@class="bleft"]/ul[@class="list list_b"]/li`)
					if len(liNodes) >= 1 {
						for _, liNode := range liNodes {
							// 文档详情URL
							detailUrl := htmlquery.InnerText(htmlquery.FindOne(liNode, `./a/@href`))
							detailUrlSplit := strings.Split(detailUrl, "/")
							fileIdString := strings.ReplaceAll(detailUrlSplit[len(detailUrlSplit)-1], ".asp", "")
							fileId, _ := strconv.Atoi(fileIdString)
							fileName := htmlquery.InnerText(htmlquery.FindOne(liNode, `./a`))
							fileName = strings.ReplaceAll(fileName, "/", "-")
							fileName = strings.ReplaceAll(fileName, ".", "")
							fileName = strings.ReplaceAll(fileName, " ", "")
							fmt.Println(fileName)
							// 下载文档URL
							downLoadUrl := fmt.Sprintf("http://down.hi138.com/downloadfile.asp?id=%d", fileId)
							fmt.Println(downLoadUrl)

							filePath := "../www.hi138.com/" + bigSubjectName + "/" + smallSubjectName + "/" + fileName + ".docx"
							if _, err := os.Stat(filePath); err != nil {
								fmt.Println("=======开始下载========")
								err = downloadHi138(downLoadUrl, detailUrl, filePath)
								if err != nil {
									fmt.Println(err)
									continue
								}
								fmt.Println("=======开始完成========")
							}
							time.Sleep(time.Second * 1)
							break
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
func downloadHi138(attachmentUrl string, referer string, filePath string) error {
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
	if Hi138EnableHttpProxy {
		client = Hi138SetHttpProxy()
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
	req.Header.Set("Host", "http://www.hi138.com")
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
