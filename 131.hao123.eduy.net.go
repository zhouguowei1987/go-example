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
	EduYEnableHttpProxy = false
	EduYHttpProxyUrl    = "27.42.168.46:55481"
)

func EduYSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(EduYHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type EduYEducationCategory struct {
	categoryName string
	categoryUrl  string
	classId      int
}

var eduYEducationCategory = []EduYEducationCategory{
	//{
	//	categoryName: "小学",
	//	categoryUrl:  "http://hao123.eduy.net/jiaoxueziyuan/shiti/xiaoxue/index.html",
	//	classId:      29,
	//},
	//{
	//	categoryName: "初中",
	//	categoryUrl:  "http://hao123.eduy.net/jiaoxueziyuan/shiti/chuzhong/index.html",
	//	classId:      30,
	//},
	{
		categoryName: "高中",
		categoryUrl:  "http://hao123.eduy.net/jiaoxueziyuan/shiti/gaozhong/index.html",
		classId:      31,
	},
}

// eduYSpider 获取阳光数学网文档
// @Title 获取阳光数学网文档
// @Description http://hao123.eduy.net/，获取阳光数学网文档
func main() {
	for _, category := range eduYEducationCategory {
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
			// /html/body/table[4]/tbody/tr/td[1]/table[2]/tbody/tr/td/table[1]/tbody/tr/td/table/tbody/tr[1]
			divNodes := htmlquery.Find(listDoc, `//html/body/table[4]/tbody/tr/td[1]/table[2]/tbody/tr/td/table[1]/tbody/tr/td/table/tbody/tr`)
			if len(divNodes) >= 1 {
				for _, divNode := range divNodes {
					// 第一个td中含有img标签的才是要提取的内容
					imgNode := htmlquery.FindOne(divNode, `./td[1]/img`)
					if imgNode == nil {
						fmt.Println("不是要提取的内容，跳过")
						continue
					}

					titleNode := htmlquery.FindOne(divNode, `./td[1]/b/a`)
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

					detailUrl := htmlquery.InnerText(htmlquery.FindOne(divNode, `./td[1]/b/a/@href`))
					detailUrlSplitArray := strings.Split(detailUrl, "/")
					idHtml := detailUrlSplitArray[len(detailUrlSplitArray)-1]
					idStr := strings.ReplaceAll(idHtml, ".html", "")
					id, _ := strconv.Atoi(idStr)

					eduYDownloadUrl := fmt.Sprintf("http://hao123.eduy.net/e/DownSys/DownSoft/?classid=%d&id=%d&pathid=0", category.classId, id)
					eduYDownloadDoc, err := htmlquery.LoadURL(eduYDownloadUrl)
					if err != nil {
						fmt.Println(err)
						continue
					}
					// /html/body/center/table/tbody/tr/td/table/tbody/tr[8]/td[2]/a
					attachmentNode := htmlquery.FindOne(eduYDownloadDoc, `/html/body/center/table/tbody/tr/td/table/tbody/tr[8]/td[2]/a/@href`)
					if attachmentNode == nil {
						fmt.Println("没有下载链接，跳过")
						continue
					}
					attachmentUrl := "http://hao123.eduy.net/e/DownSys" + strings.ReplaceAll(htmlquery.InnerText(attachmentNode), "..", "")
					fmt.Println(attachmentUrl)
					filePath := "F:\\workspace\\hao123.eduy.net\\hao123.eduy.net\\" + category.categoryName + "\\" + title + ".rar"
					_, err = os.Stat(filePath)
					if err != nil {
						fmt.Println("=======开始下载========")
						err := downloadEduY(attachmentUrl, filePath, eduYDownloadUrl)
						if err != nil {
							fmt.Println(err)
							continue
						}
						fmt.Println("=======完成下载========")
						DownLoadEduYTimeSleep := rand.Intn(10)
						for i := 1; i <= DownLoadEduYTimeSleep; i++ {
							time.Sleep(time.Second)
							fmt.Println("page="+strconv.Itoa(page)+"===========下载", title, "成功，暂停", DownLoadEduYTimeSleep, "秒，倒计时", i, "秒===========")
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

func downloadEduY(attachmentUrl string, filePath string, referer string) error {
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
	if EduYEnableHttpProxy {
		client = EduYSetHttpProxy()
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
	req.Header.Set("Host", "hao123.eduy.net")
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
