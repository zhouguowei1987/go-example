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
	ZJEduEnableHttpProxy = false
	ZJEduHttpProxyUrl    = "111.225.152.186:8089"
)

func ZJEduSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(ZJEduHttpProxyUrl)
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

var grades = []Grade{
	{
		name: "一年级",
		url:  "http://www.51zjedu.com/yinianji/",
	},
	{
		name: "二年级",
		url:  "http://www.51zjedu.com/ernianji/",
	},
	{
		name: "三年级",
		url:  "http://www.51zjedu.com/sannianji/",
	},
	{
		name: "四年级",
		url:  "http://www.51zjedu.com/sinianji/",
	},
	{
		name: "五年级",
		url:  "http://www.51zjedu.com/wunianji/",
	},
	{
		name: "六年级",
		url:  "http://www.51zjedu.com/liunianji/",
	},
}

// ychEduSpider 获取帝源教育文档
// @Title 获取帝源教育文档
// @Description http://www.51zjedu.com/，获取帝源教育文档
func main() {
	for _, grade := range grades {
		current := 1
		isPageListGo := true
		for isPageListGo {
			gradeIndexUrl := grade.url
			if current > 1 {
				gradeIndexUrl += fmt.Sprintf("index_%d.html", current)
			}
			gradeIndexDoc, err := htmlquery.LoadURL(gradeIndexUrl)
			if err != nil {
				fmt.Println(err)
				current = 1
				isPageListGo = false
				continue
			}
			liNodes := htmlquery.Find(gradeIndexDoc, `//div[@class="doc-list f-f14"]/ul/li`)
			if len(liNodes) <= 0 {
				fmt.Println(err)
				current = 1
				isPageListGo = false
				continue
			}
			for _, liNode := range liNodes {
				fmt.Println("============================================================================")
				fmt.Println("年级：", grade.name)
				fmt.Println("=======当前页为：" + strconv.Itoa(current) + "========")

				fileName := htmlquery.InnerText(htmlquery.FindOne(liNode, `./div[@class="doc-list-title"]/h3/a/@title`))
				fileName = strings.ReplaceAll(fileName, "/", "-")
				fileName = strings.ReplaceAll(fileName, ":", "-")
				fileName = strings.ReplaceAll(fileName, "：", "-")
				fileName = strings.ReplaceAll(fileName, "（", "(")
				fileName = strings.ReplaceAll(fileName, "）", ")")
				fmt.Println(fileName)

				filePath := "../www.51zjedu.com/www.51zjedu.com/" + grade.name + "/" + fileName
				_, errDoc := os.Stat(filePath + ".doc")
				_, errDocx := os.Stat(filePath + ".docx")
				_, errPdf := os.Stat(filePath + ".pdf")
				_, errPpt := os.Stat(filePath + ".ppt")
				if errDoc != nil && errDocx != nil && errPdf != nil && errPpt != nil {

					viewUrl := "http://www.51zjedu.com" + htmlquery.InnerText(htmlquery.FindOne(liNode, `./div[@class="doc-list-title"]/h3/a/@href`))
					fmt.Println(viewUrl)

					viewDoc, err := htmlquery.LoadURL(viewUrl)
					if err != nil {
						fmt.Println(err)
						continue
					}
					// 所需点数
					regPoints := regexp.MustCompile(`所需点数：([0-9]*)`)
					regPointsMatch := regPoints.FindAllSubmatch([]byte(htmlquery.InnerText(viewDoc)), -1)
					points, err := strconv.Atoi(string(regPointsMatch[0][1]))
					if err != nil {
						fmt.Println(err)
						continue
					}
					if points > 0 {
						fmt.Println("需要积分下载", points)
						continue
					}

					regDownloadViewUrl := regexp.MustCompile(`<a href="#ecms" onclick="window.open\('(.*?)','','width=500,height=300,resizable=yes'\);"`)
					regDownloadViewUrlMatch := regDownloadViewUrl.FindAllSubmatch([]byte(htmlquery.InnerText(viewDoc)), -1)
					downloadViewUrl := "http://www.51zjedu.com" + string(regDownloadViewUrlMatch[0][1])
					fmt.Println(downloadViewUrl)

					downloadViewDoc, err := htmlquery.LoadURL(downloadViewUrl)
					if err != nil {
						fmt.Println(err)
						continue
					}

					downloadUrlNode := htmlquery.FindOne(downloadViewDoc, `//a/@href`)
					if downloadUrlNode == nil {
						fmt.Println("没有下载链接")
						continue
					}
					downLoadUrl := strings.ReplaceAll(htmlquery.InnerText(downloadUrlNode), "../", "/")
					downLoadUrl = "http://www.51zjedu.com/e/DownSys" + downLoadUrl
					fmt.Println(downLoadUrl)

					fmt.Println("=======开始下载" + strconv.Itoa(current) + "========")
					err = downloadZJEdu(downLoadUrl, downloadViewUrl, filePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("=======开始完成========")
					time.Sleep(time.Millisecond * 200)
				}
			}
			current++
			isPageListGo = true
		}
	}
}

func downloadZJEdu(attachmentUrl string, referer string, filePath string) error {
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
	if ZJEduEnableHttpProxy {
		client = ZJEduSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.51zjedu.com")
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
	// 检查HTTP响应头中的Content-Disposition字段获取文件名和后缀
	fileName := getZJEduFileNameFromHeader(resp)
	fileExtension := filepath.Ext(fileName) // 获取文件后缀
	fileExtArr := []string{".doc", ".docx", ".pdf", ".ppt"}
	fmt.Println("文件后缀:", fileExtension)
	if !StrInArrayZJEdu(fileExtension, fileExtArr) {
		return errors.New("文件后缀：" + fileExtension + "不在下载后缀列表")
	}
	filePath += fileExtension

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

// StrInArrayZJEdu str in string list
func StrInArrayZJEdu(str string, data []string) bool {
	if len(data) > 0 {
		for _, row := range data {
			if str == row {
				return true
			}
		}
	}
	return false
}

// 从HTTP响应头中获取文件名
func getZJEduFileNameFromHeader(resp *http.Response) string {
	contentDisposition := resp.Header.Get("Content-Disposition")
	fileName := ""
	if contentDisposition != "" {
		fileName = parseZJEduFileNameFromContentDisposition(contentDisposition)
	} else {
		fileName = filepath.Base(resp.Request.URL.Path) // 默认使用URL中的文件名作为本地文件名
	}
	return fileName
}

// 从Content-Disposition字段中解析文件名
func parseZJEduFileNameFromContentDisposition(contentDisposition string) string {
	// 参考：https://tools.ietf.org/html/rfc6266#section-4.3
	// 示例：attachment; filename="example.txt" -> example.txt
	fileNameStart := len("attachment; ") + len("filename=") + 2 // 2为引号的长度
	fileNameEnd := len(contentDisposition) - 1 - len("\"")      // 最后一个双引号的位置
	fileName := contentDisposition[fileNameStart:fileNameEnd]   // 提取文件名字符串
	return fileName[1:]                                         // 去掉字符串开头的引号（如果存在）并返回结果
}
