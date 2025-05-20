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
	MiNi19EnableHttpProxy = false
	MiNi19HttpProxyUrl    = "27.42.168.46:55481"
)

func MiNi19SetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(MiNi19HttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type MiNi19EducationCategory struct {
	categoryName string
	categoryUrl  string
	classId      int
}

var miNi19EducationCategory = []MiNi19EducationCategory{
	{
		categoryName: "试卷",
		categoryUrl:  "http://www.19mini.cn/ziyuan/shijuan/",
		classId:      30,
	},
}

// MiNi19Spider 获取迷你语文网文档
// @Title 获取迷你语文网文档
// @Description http://www.19mini.cn/，获取迷你语文网文档
func main() {
	for _, category := range miNi19EducationCategory {
		page := 1
		isPageGo := true
		for isPageGo {
			var listUrl = fmt.Sprintf(category.categoryUrl)
			if page != 1 {
				listUrl = strings.ReplaceAll(category.categoryUrl, "index.html", "") + fmt.Sprintf("index_%d.html", page)
			}
			fmt.Println(listUrl)
			listDoc, err := htmlquery.LoadURL(listUrl)
			if err != nil {
				fmt.Println(err)
				break
			}
			divNodes := htmlquery.Find(listDoc, `//ul[@class="e2"]/li`)
			if len(divNodes) >= 1 {
				for _, divNode := range divNodes {
					titleNode := htmlquery.FindOne(divNode, `./a[@class="title"]`)
					if titleNode == nil {
						fmt.Println("标题不存在")
						continue
					}
					title := htmlquery.InnerText(titleNode)
					title = strings.TrimSpace(title)
					title = strings.ReplaceAll(title, "免费", "")
					title = strings.ReplaceAll(title, "-", "")
					title = strings.ReplaceAll(title, " ", "")
					title = strings.ReplaceAll(title, "|", "-")
					fmt.Println(title)

					detailUrl := htmlquery.InnerText(htmlquery.FindOne(divNode, `./a[@class="title"]/@href`))
					detailUrlSplitArray := strings.Split(detailUrl, "/")
					idHtml := detailUrlSplitArray[len(detailUrlSplitArray)-1]
					idStr := strings.ReplaceAll(idHtml, ".html", "")
					id, _ := strconv.Atoi(idStr)

					MiNi19DownloadUrl := fmt.Sprintf("http://www.19mini.cn/e/DownSys/DownSoft/?classid=%d&id=%d&pathid=0", category.classId, id)
					fmt.Println(MiNi19DownloadUrl)
					MiNi19DownloadDoc, err := htmlquery.LoadURL(MiNi19DownloadUrl)
					if err != nil {
						fmt.Println(err)
						continue
					}
					// /html/body/div[2]/a
					attachmentNode := htmlquery.FindOne(MiNi19DownloadDoc, `/html/body/div[2]/a/@href`)
					if attachmentNode == nil {
						fmt.Println("没有下载链接，跳过")
						continue
					}
					attachmentUrl := "http://www.19mini.cn/e/DownSys" + strings.ReplaceAll(htmlquery.InnerText(attachmentNode), "..", "")
					fmt.Println(attachmentUrl)

					// 获取文档类型
					MiNi19ViewUrl := fmt.Sprintf("http://www.19mini.cn/ziyuan/shijuan/%d.html", id)
					fmt.Println(MiNi19ViewUrl)
					MiNi19ViewDoc, err := htmlquery.LoadURL(MiNi19ViewUrl)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fileTypeNode := htmlquery.FindOne(MiNi19ViewDoc, `//div[@class="infolist"]/span[5]`)
					if fileTypeNode == nil {
						fmt.Println("文档类型不存在")
						continue
					}
					fileType := htmlquery.InnerText(fileTypeNode)
					fmt.Println(fileType)
					if strings.Index(fileType, "doc") == -1 {
						fmt.Println("文档类型不是doc文档，跳过")
						continue
					} else {
						fileType = ".doc"
					}

					filePath := "F:\\workspace\\www.19mini.cn\\www.19mini.cn\\" + category.categoryName + "\\" + title + fileType
					_, err = os.Stat(filePath)
					if err != nil {
						fmt.Println("=======开始下载========")
						err := downloadMiNi19(attachmentUrl, filePath, MiNi19DownloadUrl)
						if err != nil {
							fmt.Println(err)
							continue
						}
						fmt.Println("=======完成下载========")
						DownLoadMiNi19TimeSleep := rand.Intn(10)
						for i := 1; i <= DownLoadMiNi19TimeSleep; i++ {
							time.Sleep(time.Second)
							fmt.Println("page="+strconv.Itoa(page)+"===========下载", title, "成功，暂停", DownLoadMiNi19TimeSleep, "秒，倒计时", i, "秒===========")
						}
					}
				}
				page++
			} else {
				isPageGo = false
				page = 1
				break
			}
		}
	}
}

func downloadMiNi19(attachmentUrl string, filePath string, referer string) error {
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
	if MiNi19EnableHttpProxy {
		client = MiNi19SetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.19mini.cn")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", referer)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
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
