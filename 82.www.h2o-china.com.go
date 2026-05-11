package main

import (
	"errors"
	"fmt"
	"time"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var H2oChinaCookie = "server_name_session=0710dba57d88ab91589b685360201510; backurl=%2Fstandard%2Flist%3FcolumnId%3D18%26ordby%3Ddateline%26sort%3DDESC%26page%3D%7Bpage%7D%3D"

// h2oChinaSpider 获取中国水网Pdf文档
// @Title 获取中国水网Pdf文档
// @Description https://www.h2o-china.com/，获取中国水网Pdf文档
func main() {
	page := 1
	isPageGo := true
	for isPageGo {
		listUrl := fmt.Sprintf("https://www.h2o-china.com/standard/list?columnId=18&ordby=dateline&sort=DESC&page=%d", page)
		referUrl := "https://www.h2o-china.com/standard/list"
		if page >= 2 {
			referUrl = fmt.Sprintf("https://www.h2o-china.com/standard/list?columnId=18&ordby=dateline&sort=DESC&page=%d", page-1)
		}
		listDoc, err := H2oChinaBzHtmlDoc(listUrl, referUrl)
		if err != nil {
			fmt.Println(err)
			break
		}
		liNodes := htmlquery.Find(listDoc, `//div[@class="lists txtList"]/ul/li/em[@class="title"]/a[@class="ellip w540 i-pdf"]`)
		if len(liNodes) >= 1 {
			for _, liNode := range liNodes {
				detailUrl := htmlquery.InnerText(htmlquery.FindOne(liNode, `./@href`))
				detailUrl = "https://www.h2o-china.com" + detailUrl
				fmt.Println(detailUrl)
				detailDoc, err := H2oChinaBzHtmlDoc(detailUrl, listUrl)
                if err != nil {
                    fmt.Println(err)
                    continue
                }
				title := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="hd"]/h1`))
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, " ", "")
				fmt.Println(title)

				standardNo := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="traits"]/table/tbody/tr[3]/td[2]`))
				standardNo = strings.ReplaceAll(standardNo, "/", "-")
				standardNo = strings.ReplaceAll(standardNo, ":", "-")
				standardNo = strings.ReplaceAll(standardNo, " ", "")
				//fmt.Println(standardNo)

				downloadUrl := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="dowloads fr"]/a/@href`))
				downloadUrl = "https://www.h2o-china.com" + downloadUrl
				fmt.Println(downloadUrl)

				filePath := "../www.h2o-china.com/" + title + ".pdf"
				if len(standardNo) > 1 {
					filePath = "../www.h2o-china.com/" + title + "(" + standardNo + ")" + ".pdf"
				}
				fmt.Println(filePath)
				_, err = os.Stat(filePath)
                if err == nil {
                    fmt.Println("文档已下载过，跳过")
                    continue
                }
                // 开始下载
                fmt.Println("=======开始下载========")
				err = downloadH2oChinaPdf(downloadUrl, filePath)
				if err != nil {
					fmt.Println(err)
				}
				//复制文件
                tempFilePath := strings.ReplaceAll(filePath, "www.h2o-china.com", "temp-hbba.sacinfo.org.cn")
                err = copyH2oChinaFile(filePath, tempFilePath)
                if err != nil {
                    fmt.Println(err)
                    continue
                }
                fmt.Println("=======完成下载========")

                // 设置倒计时
                DownLoadTH2oChinaTimeSleep := 10
                for i := 1; i <= DownLoadTH2oChinaTimeSleep; i++ {
                    time.Sleep(time.Second)
                    fmt.Println("===page = "+strconv.Itoa(page)+"===title="+title+"===========操作完成，", "暂停", DownLoadTH2oChinaTimeSleep, "秒，倒计时", i, "秒===========")
                }
			}
			DownLoadH2oChinaPageTimeSleep := 10
			// DownLoadH2oChinaPageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadH2oChinaPageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("===page = "+strconv.Itoa(page)+"========= 暂停", DownLoadH2oChinaPageTimeSleep, "秒 倒计时", i, "秒===========")
			}
			page++
		} else {
			isPageGo = false
			page = 1
			break
		}
	}
}

func H2oChinaBzHtmlDoc(requestUrl string, referer string) (doc *html.Node, err error) {
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
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接
	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	// 	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", H2oChinaCookie)
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	req.Header.Set("Host", "www.h2o-china.com")
	req.Header.Set("Origin", "https://www.h2o-china.com/")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return doc, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return doc, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	doc, err = htmlquery.Parse(resp.Body)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func downloadH2oChinaPdf(pdfUrl string, filePath string) error {
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
	req, err := http.NewRequest("GET", pdfUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
// 	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", H2oChinaCookie)
	req.Header.Set("Host", "www.h2o-china.com")
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
func copyH2oChinaFile(src, dst string) (err error) {
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