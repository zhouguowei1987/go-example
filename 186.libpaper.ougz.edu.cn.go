package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const (
	LibPaperEnableHttpProxy = false
	LibPaperHttpProxyUrl    = "111.225.152.186:8089"
)

func LibPaperSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(LibPaperHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type QueryLibPaperDownloadUrlFormData struct {
	resourceid string
}

var LibPaperCookie = "JSESSIONID=519CE66608776719CAF2BA698F0C9FFA"

// ychEduSpider 获取广州开放大学试卷文档
// @Title 获取广州开放大学试卷文档
// @Description https://libpaper.ougz.edu.cn/，获取广州开放大学试卷文档
func main() {
	page := 1
	isPageListGo := true
	for isPageListGo {
		requestListUrl := fmt.Sprintf("https://libpaper.ougz.edu.cn/index.jsp?wbtreeid=5488&currentnum=%d&_vt=&_vn=&_vk=&_vm=&_va=&_vs=&_vc=&_vb=&_ve=&searchScope=0", page)
		referUrl := "https://libpaper.ougz.edu.cn/"
		LibPaperListDoc, err := LibPaperHtmlDoc(requestListUrl, referUrl)
		if err != nil {
			fmt.Println(err)
			break
		}
		trNodes := htmlquery.Find(LibPaperListDoc, `//div[@class="main"]/div[@class="right"]/table[@class="pn-ltable"]/tbody[@class="pn-ltbody"]/tr`)
		if len(trNodes) >= 1 {
			for _, trNode := range trNodes {
				// 试卷代号
				codeNode := htmlquery.FindOne(trNode, `./td[@class="code"]`)
				if codeNode == nil {
					fmt.Println("试卷代号不存在，跳过")
					continue
				}
				code := htmlquery.InnerText(codeNode)
				code = strings.TrimSpace(code)
				fmt.Println(code)

				// 试卷标题
				titleNode := htmlquery.FindOne(trNode, `./td[@class="title"]`)
				if titleNode == nil {
					fmt.Println("标题不存在，跳过")
					continue
				}
				title := htmlquery.InnerText(titleNode)
				title = strings.TrimSpace(title)
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, " ", "")
				fmt.Println(title)

				// 试卷时间
				examTimeNode := htmlquery.FindOne(trNode, `./td[@class="examTime"]`)
				if examTimeNode == nil {
					fmt.Println("试卷时间不存在，跳过")
					continue
				}
				examTime := htmlquery.InnerText(examTimeNode)
				examTime = strings.TrimSpace(examTime)
				fmt.Println(examTime)

				// 试卷文件
				examFileUrlNode := htmlquery.FindOne(trNode, `./td[@class="contextPath"]/a[1]/@href`)
				if examFileUrlNode == nil {
					fmt.Println("试卷文件不存在，跳过")
					continue
				}
				examFileUrl := htmlquery.InnerText(examFileUrlNode)
				examFileUrl = strings.TrimSpace(examFileUrl)
				fmt.Println(examFileUrl)
				fileExt := filepath.Ext(examFileUrl)
				fileExt = strings.ToLower(fileExt)
				if strings.Index(fileExt, "doc") == -1 && strings.Index(fileExt, "pdf") == -1 && strings.Index(fileExt, "ppt") == -1 {
					fmt.Println("文档不是doc文档，跳过")
					continue
				}

				filePath := "E:\\workspace\\libpaper.ougz.edu.cn\\libpaper.ougz.edu.cn\\" + examTime + "国家开放大学《" + title + code + "》试题(含答案及评分标准)" + fileExt
				fmt.Println(filePath)
				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				downLoadUrl := "https://libpaper.ougz.edu.cn" + examFileUrl
				fmt.Println(downLoadUrl)

				// 开始下载
				fmt.Println("=======开始下载========")
				err = downloadLibPaper(downLoadUrl, requestListUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "libpaper.ougz.edu.cn\\libpaper.ougz.edu.cn", "libpaper.ougz.edu.cn\\temp-libpaper.ougz.edu.cn")
				err = copyLibPaperFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======完成下载========")
				// 设置倒计时
				DownLoadTLibPaperTimeSleep := 10
				for i := 1; i <= DownLoadTLibPaperTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page = "+strconv.Itoa(page)+"===title="+title+"===========操作完成，", "暂停", DownLoadTLibPaperTimeSleep, "秒，倒计时", i, "秒===========")
				}
			}
			DownLoadLibPaperPageTimeSleep := 10
			// DownLoadLibPaperPageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadLibPaperPageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page = "+strconv.Itoa(page)+"========= 暂停", DownLoadLibPaperPageTimeSleep, "秒 倒计时", i, "秒===========")
			}
			page++
		} else {
			page = 0
			isPageListGo = false
			break
		}
	}
}

func LibPaperHtmlDoc(requestUrl string, referer string) (doc *html.Node, err error) {
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
	if LibPaperEnableHttpProxy {
		client = LibPaperSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接
	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	// req.Header.Set("Cache-Control", "max-age=0")
	// req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", LibPaperCookie)
	req.Header.Set("Host", "libpaper.ougz.edu.cn")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"")
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

func downloadLibPaper(downloadUrl string, referer string, filePath string) error {
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
	if LibPaperEnableHttpProxy {
		client = LibPaperSetHttpProxy()
	}
	req, err := http.NewRequest("GET", downloadUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", LibPaperCookie)
	req.Header.Set("Origin", "https://libpaper.ougz.edu.cn")
	req.Header.Set("Host", "libpaper.ougz.edu.cn")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"124\", \"Google Chrome\";v=\"124\", \"Not-A.Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
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

func copyLibPaperFile(src, dst string) (err error) {
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
