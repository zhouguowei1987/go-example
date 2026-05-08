package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"strings"
	"time"

	"github.com/antchfx/htmlquery"
)

var GetTiKuEnableHttpProxy = false
var GetTiKuHttpProxyUrl = ""
var GetTiKuHttpProxyUrlArr = make([]string, 0)

func GetTiKuHttpProxy() error {
	pageMax := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, page := range pageMax {
		freeProxyUrl := "https://www.beesproxy.com/free"
		if page > 1 {
			freeProxyUrl = fmt.Sprintf("https://www.beesproxy.com/free/page/%d", page)
		}
		beesProxyDoc, err := htmlquery.LoadURL(freeProxyUrl)
		if err != nil {
			return err
		}
		trNodes := htmlquery.Find(beesProxyDoc, `//figure[@class="wp-block-table"]/table[@class="table table-bordered bg--secondary"]/tbody/tr`)
		if len(trNodes) > 0 {
			for _, trNode := range trNodes {
				ipNode := htmlquery.FindOne(trNode, "./td[1]")
				if ipNode == nil {
					continue
				}
				ip := htmlquery.InnerText(ipNode)

				portNode := htmlquery.FindOne(trNode, "./td[2]")
				if portNode == nil {
					continue
				}
				port := htmlquery.InnerText(portNode)

				protocolNode := htmlquery.FindOne(trNode, "./td[5]")
				if protocolNode == nil {
					continue
				}
				protocol := htmlquery.InnerText(protocolNode)

				switch protocol {
				case "HTTP":
					GetTiKuHttpProxyUrlArr = append(GetTiKuHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					GetTiKuHttpProxyUrlArr = append(GetTiKuHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func GetTiKuSetHttpProxy() (httpclient *http.Client) {
	if GetTiKuHttpProxyUrl == "" {
		if len(GetTiKuHttpProxyUrlArr) <= 0 {
			err := GetTiKuHttpProxy()
			if err != nil {
				GetTiKuSetHttpProxy()
			}
		}
		GetTiKuHttpProxyUrl = GetTiKuHttpProxyUrlArr[0]
		if len(GetTiKuHttpProxyUrlArr) >= 2 {
			GetTiKuHttpProxyUrlArr = GetTiKuHttpProxyUrlArr[1:]
		} else {
			GetTiKuHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(GetTiKuHttpProxyUrl)
	ProxyURL, _ := url.Parse(GetTiKuHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*3)
				if err != nil {
					fmt.Println("dail timeout", err)
					return nil, err
				}
				return c, nil

			},
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 30,
		},
	}
	return httpclient
}

type GetTiKuSubjects struct {
	name string
	id   int
}

type GetTiKuSubjectGrades struct {
	name string
	id   int
}

var getTiKuSubjects = []GetTiKuSubjects{
	// {
	// 	name: "语文",
	// 	id:   1,
	// },
	// {
	// 	name: "数学",
	// 	id:   2,
	// },
	{
		name: "英语",
		id:   3,
	},
}
var getTiKuSubjectGrades = []GetTiKuSubjectGrades{
	{
		name: "一年级上",
		id:   111,
	},
	{
		name: "一年级下",
		id:   112,
	},
	{
		name: "二年级上",
		id:   121,
	},
	{
		name: "二年级下",
		id:   122,
	},
	{
		name: "三年级上",
		id:   131,
	},
	{
		name: "三年级下",
		id:   132,
	},
	{
		name: "四年级上",
		id:   141,
	},
	{
		name: "四年级下",
		id:   142,
	},
	{
		name: "五年级上",
		id:   151,
	},
	{
		name: "五年级下",
		id:   152,
	},
	{
		name: "六年级上",
		id:   161,
	},
	{
		name: "六年级下",
		id:   162,
	},
}

type GetTiKuPaper struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type GetTiKuPaperDetail struct {
	Content string                      `json:"content"`
	Grade   string                      `json:"grade"`
	PubTime int                         `json:"pub_time"`
	Subject string                      `json:"subject"`
	DownArr []GetTiKuPaperDetailDownArr `json:"down_arr"`
	Title   string                      `json:"title"`
}
type GetTiKuPaperDetailDownArr struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

const GetTiKuNextDownloadSleep = 5

// ychEduSpider 获取考卷下载试卷
// @Title 获取考卷下载试卷
// @Description http://gettiku.com/，获取考卷下载试卷
func main() {
	for _, subject := range getTiKuSubjects {
		for _, grade := range getTiKuSubjectGrades {
			current := 1
			isPageListGo := true
			for isPageListGo {
				subjectIndexUrl := fmt.Sprintf("https://gettiku.com/api/article/posts?grade=%d&subject=%d&tag=1&page=%d", grade.id, subject.id, current)
				fmt.Println(subjectIndexUrl)
				subjectIndexDoc, err := htmlquery.LoadURL(subjectIndexUrl)
				subjectIndexStr := htmlquery.InnerText(subjectIndexDoc)
				var getTiKuPapers []GetTiKuPaper
				err = json.Unmarshal([]byte(subjectIndexStr), &getTiKuPapers)
				if err != nil {
					fmt.Println(err)
					current = 1
					isPageListGo = false
					continue
				}
				if len(getTiKuPapers) <= 0 {
					fmt.Println(err)
					current = 1
					isPageListGo = false
					continue
				}
				for _, paper := range getTiKuPapers {
					fmt.Println("============================================================================")
					fmt.Println("主题：", subject.name, grade.name, paper.Title)
					if strings.Index(paper.Title, "试") == -1 {
						fmt.Println("不含有'试'字，跳过")
						continue
					}
					fmt.Println("=======当前页URL", subjectIndexUrl, "========")

					viewUrl := fmt.Sprintf("https://gettiku.com/api/article/q?id=%d", paper.Id)
					fmt.Println(viewUrl)

					viewDoc, _ := htmlquery.LoadURL(viewUrl)
					if viewDoc == nil {
						fmt.Println("获取试卷详情失败")
						continue
					}
					viewDocStr := htmlquery.InnerText(viewDoc)
					var getTiKuPaperDetail GetTiKuPaperDetail
					err = json.Unmarshal([]byte(viewDocStr), &getTiKuPaperDetail)

					if err != nil {
						fmt.Println(err)
						current = 1
						isPageListGo = false
						continue
					}
					fileName := getTiKuPaperDetail.DownArr[0].Name
					fileName = strings.TrimSpace(fileName)
					if strings.Index(fileName, "含答案") == -1 {
						fileNameArray := strings.Split(fileName, ".")
						fileName = fileNameArray[0] + "(含答案)" + "." + fileNameArray[1]
					}
					fileName = grade.name + "(" + subject.name + ")" + fileName
					fileName = strings.ReplaceAll(fileName, "<b>", "")
					fileName = strings.ReplaceAll(fileName, "</b>", "")
					fileName = strings.ReplaceAll(fileName, "/", "-")
					fileName = strings.ReplaceAll(fileName, ":", "-")
					fileName = strings.ReplaceAll(fileName, " ", "_")
					fileName = strings.ReplaceAll(fileName, "：", "-")
					fileName = strings.ReplaceAll(fileName, "（", "(")
					fileName = strings.ReplaceAll(fileName, "）", ")")
					fileName = strings.ReplaceAll(fileName, "word，", "")
					fileName = strings.ReplaceAll(fileName, "word版，", "")
					fileName = strings.ReplaceAll(fileName, "(含答案)(含答案)", "(含答案)")
					fmt.Println(fileName)

					filePath := "../gettiku.com/gettiku.com/" + subject.name + "/" + fileName
					fmt.Println(filePath)
					_, err = os.Stat(filePath)
					if err == nil {
						fmt.Println("文档已下载过，跳过")
						continue
					}
					downLoadUrl := getTiKuPaperDetail.DownArr[0].Url
					fmt.Println(downLoadUrl)

					fmt.Println("=======开始下载" + strconv.Itoa(current) + "========")
					err = downloadGetTiKu(downLoadUrl, viewUrl, filePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					//复制文件
					tempFilePath := strings.ReplaceAll(filePath, "gettiku.com/gettiku.com", "gettiku.com/temp-gettiku.com")
					err = copyGetTiKuFile(filePath, tempFilePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("=======下载完成========")
					for i := 1; i <= GetTiKuNextDownloadSleep; i++ {
						time.Sleep(time.Second)
						fmt.Println("===========操作结束，暂停", GetTiKuNextDownloadSleep, "秒，倒计时", i, "秒===========")
					}
				}
				current++
				isPageListGo = true
			}
		}
	}
}

func downloadGetTiKu(attachmentUrl string, referer string, filePath string) error {
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
			ResponseHeaderTimeout: time.Second * 30,
		},
	}
	if GetTiKuEnableHttpProxy {
		client = GetTiKuSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "gettiku.com")
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

func copyGetTiKuFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(in *os.File) {
		err := in.Close()
		if err != nil {
			return
		}
	}(in)

	// 创建一个文件用于保存
	fileDiv := filepath.Dir(dst)
	if _, err = os.Stat(fileDiv); err != nil {
		if os.MkdirAll(fileDiv, 0o777) != nil {
			return err
		}
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)

	_, err = io.Copy(out, in)
	return nil
}
