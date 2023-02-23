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
	ChaZiDianEnableHttpProxy = false
	ChaZiDianHttpProxyUrl    = "111.225.152.186:8089"
)

func ChaZiDianSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(ChaZiDianHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type ChaZiDianSubject struct {
	name string
	url  string
}

var AllChaZiDianSubject = []ChaZiDianSubject{
	//{
	//	name: "生物试题",
	//	url:  "https://shengwu.chazidian.com/shiti/",
	//},
	//{
	//	name: "物理试题",
	//	url:  "https://wuli.chazidian.com/shiti/",
	//},
	//{
	//	name: "化学试题",
	//	url:  "https://huaxue.chazidian.com/shiti/",
	//},
	//{
	//	name: "政治试题",
	//	url:  "https://zhengzhi.chazidian.com/shiti/",
	//},
	//{
	//	name: "历史试题",
	//	url:  "https://lishi.chazidian.com/shiti/",
	//},
	{
		name: "地理试题",
		url:  "https://dili.chazidian.com/shiti/",
	},
}

// ychEduSpider 获取查字典文档
// @Title 获取查字典文档
// @Description https://www.chazidian.com/，获取查字典文档
func main() {
	for _, subject := range AllChaZiDianSubject {
		page := 1
		isPageListGo := true
		for isPageListGo {
			pageListUrl := fmt.Sprintf(subject.url+"?page=%d", page)
			pageListDoc, err := htmlquery.LoadURL(pageListUrl)
			if err != nil {
				fmt.Println(err)
				break
			}
			dlNodes := htmlquery.Find(pageListDoc, `//div[@class="zyxz-m"]/div[@class="m-shl"]/dl]`)
			if len(dlNodes) >= 1 {
				for _, dlNode := range dlNodes {

					fmt.Println("=================================================================================")
					// 文档详情URL
					fileName := htmlquery.InnerText(htmlquery.FindOne(dlNode, `./dt/a`))
					fmt.Println(fileName)

					detailUrl := htmlquery.InnerText(htmlquery.FindOne(dlNode, `./dt/a/@href`))
					fmt.Println(detailUrl)

					// 发布时间
					uploadDate := htmlquery.InnerText(htmlquery.FindOne(dlNode, `./dd[1]`))
					uploadDate = strings.ReplaceAll(uploadDate, "发布时间：", "")
					fmt.Println(uploadDate)

					yearMonthDay := strings.Split(uploadDate, "-")
					if year, _ := strconv.Atoi(yearMonthDay[0]); year < 2019 {
						isPageListGo = false
						page = 1
						break
					}

					detailDoc, _ := htmlquery.LoadURL(detailUrl)
					// 文件格式
					attachmentFormat := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="jxff-l"]/dl[@class="wz-zyxq"]/dd[2]/span`))
					attachmentFormat = strings.ReplaceAll(attachmentFormat, "格式：", "")
					fmt.Println(attachmentFormat)

					// 下载文档URL
					downLoadUrl := strings.ReplaceAll(subject.url, "/shiti/", "") + htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="jxff-l"]/dl[@class="wz-zyxq"]/dd[@class="zy-djxz"]/a/@href`))
					//downLoadUrl := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="jxff-l"]/dl[@class="wz-zyxq"]/dd[@class="zy-djxz"]/a/@href`))
					fmt.Println(downLoadUrl)
					filePath := "../www.chazidian.com/" + subject.name + "/"
					err = downloadChaZiDian(downLoadUrl, filePath, fileName+"."+attachmentFormat)
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
func downloadChaZiDian(attachmentUrl string, filePath string, fileName string) error {
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
	if ChaZiDianEnableHttpProxy {
		client = ChaZiDianSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

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
