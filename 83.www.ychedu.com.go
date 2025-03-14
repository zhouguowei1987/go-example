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
	YchEduEnableHttpProxy = false
	YchEduHttpProxyUrl    = "27.42.168.46:55481"
)

func YchEduSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(YchEduHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type AdultEducationCategory struct {
	categoryName string
	categoryUrl  string
}

var adultEducationCategory = []AdultEducationCategory{
	//{categoryName: "中考试题", categoryUrl: "http://yw.ychedu.com/ZCZT/Index.html"},
	//{categoryName: "中考试题", categoryUrl: "http://yw.ychedu.com/ZCML/Index.html"},
	//{categoryName: "中考试题", categoryUrl: "http://sx.ychedu.com/ZCZT/Index.html"},
	//{categoryName: "中考试题", categoryUrl: "http://sx.ychedu.com/ZCML/Index.html"},
	//{categoryName: "中考试题", categoryUrl: "http://yy.ychedu.com/ZCZT/Index.html"},
	//{categoryName: "中考试题", categoryUrl: "http://yy.ychedu.com/ZCML/Index.html"},
	//{categoryName: "中考试题", categoryUrl: "http://wl.ychedu.com/ZCZT/Index.html"},
	//{categoryName: "中考试题", categoryUrl: "http://wl.ychedu.com/ZCML/Index.html"},
	//{categoryName: "中考试题", categoryUrl: "http://hx.ychedu.com/ZCZT/Index.html"},
	//{categoryName: "中考试题", categoryUrl: "http://hx.ychedu.com/ZCML/Index.html"},
	//{categoryName: "中考试题", categoryUrl: "http://zz.ychedu.com/ZCZT/Index.html"},
	//{categoryName: "中考试题", categoryUrl: "http://zz.ychedu.com/ZCML/Index.html"},
	//{categoryName: "中考试题", categoryUrl: "http://ls.ychedu.com/ZCZT/Index.html"},
	//{categoryName: "中考试题", categoryUrl: "http://ls.ychedu.com/ZCML/Index.html"},
	//{categoryName: "中考试题", categoryUrl: "http://qt.ychedu.com/ZCZT/Index.html"},
	//{categoryName: "中考试题", categoryUrl: "http://qt.ychedu.com/ZCML/Index.html"},
	{categoryName: "高考试题", categoryUrl: "http://yw.ychedu.com/GKZT/Index.html"},
	{categoryName: "高考试题", categoryUrl: "http://yw.ychedu.com/GKML/Index.html"},
	{categoryName: "高考试题", categoryUrl: "http://sx.ychedu.com/GKZT/Index.html"},
	{categoryName: "高考试题", categoryUrl: "http://sx.ychedu.com/GKML/Index.html"},
	{categoryName: "高考试题", categoryUrl: "http://yy.ychedu.com/GKZT/Index.html"},
	{categoryName: "高考试题", categoryUrl: "http://yy.ychedu.com/GKML/Index.html"},
	{categoryName: "高考试题", categoryUrl: "http://wl.ychedu.com/GKZT/Index.html"},
	{categoryName: "高考试题", categoryUrl: "http://wl.ychedu.com/GKML/Index.html"},
	{categoryName: "高考试题", categoryUrl: "http://hx.ychedu.com/GKZT/Index.html"},
	{categoryName: "高考试题", categoryUrl: "http://hx.ychedu.com/GKML/Index.html"},
	{categoryName: "高考试题", categoryUrl: "http://zz.ychedu.com/GKZT/Index.html"},
	{categoryName: "高考试题", categoryUrl: "http://zz.ychedu.com/GKML/Index.html"},
	{categoryName: "高考试题", categoryUrl: "http://ls.ychedu.com/GKZT/Index.html"},
	{categoryName: "高考试题", categoryUrl: "http://ls.ychedu.com/GKML/Index.html"},
	{categoryName: "高考试题", categoryUrl: "http://qt.ychedu.com/GKZT/Index.html"},
	{categoryName: "高考试题", categoryUrl: "http://qt.ychedu.com/GKML/Index.html"},
}

// ychEduSpider 获取宜城教育文档
// @Title 获取宜城教育文档
// @Description http://www.ychedu.com/，获取宜城教育文档
func main() {
	for _, category := range adultEducationCategory {
		page := 0
		maxPage := 0
		isPageGo := true
		for isPageGo {
			var listUrl = fmt.Sprintf(category.categoryUrl)
			if page != 0 {
				listUrl = strings.ReplaceAll(category.categoryUrl, "Index.html", "") + fmt.Sprintf("List_%d.html", page)
			}
			fmt.Println(listUrl)
			// 获取最大页面
			listDoc, err := htmlquery.LoadURL(listUrl)
			if err != nil {
				fmt.Println(err)
				break
			}
			if maxPage == 0 {
				countNode := htmlquery.FindOne(listDoc, `//div[@class="showpage"]/b`)
				if countNode == nil {
					fmt.Println("文档总数不存在")
					break
				}
				countInt, _ := strconv.Atoi(htmlquery.InnerText(countNode))
				maxPage = countInt/(27) + 1
				//page = maxPage / 2
			}
			divNodes := htmlquery.Find(listDoc, `//div[@class="bk21"]/div[@align="center"][1]/div`)
			if len(divNodes) >= 1 {
				for _, divNode := range divNodes {
					detailUrl := htmlquery.InnerText(htmlquery.FindOne(divNode, `./ul[@id="soft_lb1"]/div/li/a/@href`))
					detailDoc, _ := htmlquery.LoadURL(detailUrl)
					fmt.Println(detailUrl)

					titleNode := htmlquery.FindOne(divNode, `./ul[@id="soft_lb1"]/div/li/a`)
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

					ychEduDownloadUrlNode := htmlquery.FindOne(detailDoc, `//div[@class="nr10down"]/a/@href`)
					if ychEduDownloadUrlNode == nil {
						fmt.Println("下载链接不存在")
						continue
					}
					ychEduDownloadUrl := htmlquery.InnerText(ychEduDownloadUrlNode)
					fmt.Println(ychEduDownloadUrl)
					zipFilePath := "F:\\workspace\\www.ychedu.com\\www.ychedu.com\\" + category.categoryName + "\\" + title + ".zip"
					rarFilePath := "F:\\workspace\\www.ychedu.com\\www.ychedu.com\\" + category.categoryName + "\\" + title + ".rar"
					docFilePath := "F:\\workspace\\www.ychedu.com\\www.ychedu.com\\" + category.categoryName + "\\" + title + ".doc"
					docxFilePath := "F:\\workspace\\www.ychedu.com\\www.ychedu.com\\" + category.categoryName + "\\" + title + ".docx"
					_, zipErr := os.Stat(zipFilePath)
					_, rarErr := os.Stat(rarFilePath)
					_, docErr := os.Stat(docFilePath)
					_, docxErr := os.Stat(docxFilePath)
					if zipErr != nil && rarErr != nil && docErr != nil && docxErr != nil {
						fmt.Println("=======开始下载========")
						filePath := "F:\\workspace\\www.ychedu.com\\www.ychedu.com\\" + category.categoryName
						err := downloadYchEdu(ychEduDownloadUrl, filePath, title)
						if err != nil {
							fmt.Println(err)
							// 创建一个空文件，防止重复下载访问
							filePath = filePath + "\\" + title + ".rar"
							fmt.Println(filePath)
							fmt.Println("=======开始创建空文件========")
							_, err := os.Create(filePath)
							if err != nil {
								fmt.Println(err)
								fmt.Println("创建空文件失败")
							}
							fmt.Println("=======完成创建空文件========")
						}
						fmt.Println("=======完成下载========")
						DownLoadYchEduTimeSleep := rand.Intn(10)
						for i := 1; i <= DownLoadYchEduTimeSleep; i++ {
							time.Sleep(time.Second)
							fmt.Println("page="+strconv.Itoa(page)+"===========下载", title, "成功，暂停", DownLoadYchEduTimeSleep, "秒，倒计时", i, "秒===========")
						}
					}
				}
				page++
				if page > maxPage {
					isPageGo = false
					page = 0
					fmt.Println("没有更多分页")
					break
				}
			} else {
				isPageGo = false
				page = 0
				break
			}
		}
	}
}

func ychEduStringContains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func downloadYchEdu(attachmentUrl string, filePath string, title string) error {
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
	if YchEduEnableHttpProxy {
		client = YchEduSetHttpProxy()
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
	req.Header.Set("Host", "www.ychedu.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "http://www.ychedu.com/")
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

	var suffix string
	contentType := resp.Header.Get("Content-Type")
	fmt.Println(contentType)
	switch contentType {
	case "application/x-zip-compressed":
		suffix = "zip"
		break
	case "application/octet-stream":
		suffix = "rar"
		break
	case "application/msword":
		suffix = "doc"
		break
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		suffix = "docx"
		break
	default:
		return nil
	}
	fileSuffixArray := []string{"zip", "rar", "doc", "docx"}
	if !ychEduStringContains(fileSuffixArray, suffix) {
		return errors.New("既不是zip文件，也不是rar文件，跳过")
	}
	// 创建一个文件用于保存
	filePath = filePath + "\\" + title + "." + suffix
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
