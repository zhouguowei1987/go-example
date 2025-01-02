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
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	Ppt818EnableHttpProxy = false
	Ppt818HttpProxyUrl    = "111.225.152.186:8089"
)

func Ppt818SetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(Ppt818HttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type Ppt818Subject struct {
	name string
	url  string
}

var AllPpt818Subject = []Ppt818Subject{
	//{
	//	name: "excel模板",
	//	url:  "http://www.ppt818.com/list_1_quanbu/",
	//},
	{
		name: "ppt模板",
		url:  "http://www.ppt818.com/list_2_quanbu/",
	},
	//{
	//	name: "word模板",
	//	url:  "http://www.ppt818.com/list_3_quanbu/",
	//},
}

// ychEduSpider 获取pc6文档
// @Title 获取pc6文档
// @Description http://www.ppt818.com/，获取pc6文档
func main() {
	for _, subject := range AllPpt818Subject {
		pageListUrl := subject.url
		fmt.Println(pageListUrl)

		page := 1
		isPageListGo := true
		for isPageListGo {
			if page > 1 {
				pageListUrl = fmt.Sprintf(subject.url+"updated_at_%d.html", page)
			}
			fmt.Println(pageListUrl)

			pageListDoc, err := htmlquery.LoadURL(pageListUrl)
			if err != nil {
				fmt.Println(err)
				break
			}

			dlNodes := htmlquery.Find(pageListDoc, `//div[@class="clearfix mb-w"]/a]`)
			if len(dlNodes) >= 1 {
				for _, dlNode := range dlNodes {

					fmt.Println("=================================================================================")
					// 文档详情URL
					fileName := htmlquery.InnerText(htmlquery.FindOne(dlNode, `./div[@class="abs bb mb-b-des"]/div[@class="ell f14 cfff mbbd-inner"]`))
					fmt.Println(fileName)

					// 跳过文件名中含有“课件”字样文件
					if strings.Index(fileName, "课件") != -1 {
						fmt.Println("跳过文件名中含有“课件”字样文件")
						continue
					}

					// 跳过文件名中含有“课时”字样文件
					if strings.Index(fileName, "课时") != -1 {
						fmt.Println("跳过文件名中含有“课时”字样文件")
						continue
					}

					// 跳过文件名中不含有“PPT模板”字样文件
					if strings.Index(strings.ToLower(fileName), "ppt") == -1 {
						fmt.Println("跳过文件名中不含有“ppt”字样文件")
						continue
					}

					detailUrl := htmlquery.InnerText(htmlquery.FindOne(dlNode, `./@href`))
					detailUrl = "http://www.ppt818.com" + detailUrl
					fmt.Println(detailUrl)

					detailDoc, _ := htmlquery.LoadURL(detailUrl)

					// 下载文档URL
					downLoadUrl := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//a[@class="db f16 cfff tc mt30 download-btn downloadCount"]/@href`))
					fmt.Println(downLoadUrl)

					// 文件格式
					attachmentFormat := strings.Split(downLoadUrl, ".")

					if In(attachmentFormat[len(attachmentFormat)-1], []string{"ppt", "pptx"}) == false {
						fmt.Println("不是ppt文件")
						continue
					}
					filePath := "../www.ppt818.com/www.ppt818.com/" + subject.name + "/"
					fileName = fileName + "." + attachmentFormat[len(attachmentFormat)-1]
					if _, err = os.Stat(filePath + fileName); err != nil {
						fmt.Println("=======开始下载" + fileName + "========")
						err = downloadPpt818(downLoadUrl, detailUrl, filePath, fileName)
						if err != nil {
							fmt.Println(err)
							continue
						}
						fmt.Println("=======下载完成========")
						DownLoadPPT818TimeSleep := rand.Intn(5)
						for i := 1; i <= DownLoadPPT818TimeSleep; i++ {
							time.Sleep(time.Second)
							fmt.Println("page="+strconv.Itoa(page)+"===========下载", fileName, "成功，暂停", DownLoadPPT818TimeSleep, "秒，倒计时", i, "秒===========")
						}
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
func downloadPpt818(attachmentUrl string, referer string, filePath string, fileName string) error {
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
	if Ppt818EnableHttpProxy {
		client = Ppt818SetHttpProxy()
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
	req.Header.Set("Host", "http://www.ppt818.com/")
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

func In(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}
