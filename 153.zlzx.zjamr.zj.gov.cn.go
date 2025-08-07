package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	_ "os"
	"path/filepath"
	_ "path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	_ "golang.org/x/net/html"
)

var ZlZxEnableHttpProxy = false
var ZlZxHttpProxyUrl = "111.225.152.186:8089"
var ZlZxHttpProxyUrlArr = make([]string, 0)

func ZlZxHttpProxy() error {
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
					ZlZxHttpProxyUrlArr = append(ZlZxHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					ZlZxHttpProxyUrlArr = append(ZlZxHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func ZlZxSetHttpProxy() (httpclient *http.Client) {
	if ZlZxHttpProxyUrl == "" {
		if len(ZlZxHttpProxyUrlArr) <= 0 {
			err := ZlZxHttpProxy()
			if err != nil {
				ZlZxSetHttpProxy()
			}
		}
		ZlZxHttpProxyUrl = ZlZxHttpProxyUrlArr[0]
		if len(ZlZxHttpProxyUrlArr) >= 2 {
			ZlZxHttpProxyUrlArr = ZlZxHttpProxyUrlArr[1:]
		} else {
			ZlZxHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(ZlZxHttpProxyUrl)
	ProxyURL, _ := url.Parse(ZlZxHttpProxyUrl)
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
			ResponseHeaderTimeout: time.Second * 3,
		},
	}
	return httpclient
}

var ZlZxCookie = "node_id=nginx_1; cna=MYMPIfeeTQACAf////+VQMz4; arialoadData=false; zh_choose_undefined=s; cssstyle=1; zwlogBaseinfo=eyJsbF91c2VyaWQiOiIiLCJsb2dfc3RhdHVzIjoi5pyq55m75b2VIiwidXNlclR5cGUiOiJndWVzdCIsImJpel9zZXNzaW9uX2lkIjoiMjFjODlmY2ZiNDQ4NDJjMTgxMDhiODAwMmQ5MDdjMjEiLCJzaXRlX2lkIjoxLCJwYWdlX21vZGUiOiLluLjop4TmqKHlvI8ifQ==; _d_id=6cbe0b89f633a6ccfdaae4c793e662"

// 下载浙江标准在线文档
// @Title 下载浙江标准在线文档
// @Description https://zlzx.zjamr.zj.gov.cn/LocalStandard/Index/，下载浙江标准在线文档
func main() {
	page := 1
	maxPage := 3206
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		pageListUrl := fmt.Sprintf("https://zlzx.zjamr.zj.gov.cn/bzzx/public/news/list/BZBP/ALL/%d.html", page)
		fmt.Println(pageListUrl)

		referUrl := "https://zlzx.zjamr.zj.gov.cn/bzzx/public/news/list/BZBP/ALL/1.html"
		if page > 1 {
			referUrl = fmt.Sprintf("https://zlzx.zjamr.zj.gov.cn/bzzx/public/news/list/BZBP/ALL/%d.html", page-1)
		}

		queryZlZxListDoc, err := QueryZlZxList(pageListUrl, referUrl)
		if err != nil {
			fmt.Println(err)
			break
		}
		// /html/body/div[2]/div[2]/ul[1]/a[1]
		aNodes := htmlquery.Find(queryZlZxListDoc, `//html/body[@class="zybj"]/div[@class="xwlbs-cc"]/div[@class="xwlbs"]/ul[@class="tzgg_bj"]/a`)
		if len(aNodes) > 0 {
			for _, aNode := range aNodes {
				fmt.Println("=====================开始处理数据 page = ", page, "=========================")
				// /html/body/div[2]/div[2]/ul[1]/a[1]/li/span[1]
				codeTitleNode := htmlquery.FindOne(aNode, `./li/span[1]`)
				codeTitle := htmlquery.InnerText(codeTitleNode)
				codeTitleArray := strings.Split(codeTitle, "  |  ")
				code := strings.ReplaceAll(codeTitleArray[0], "/", "-")
				code = strings.ReplaceAll(code, "—", "-")
				fmt.Println(code)

				title := strings.TrimSpace(codeTitleArray[1])
				title = strings.ReplaceAll(title, " ", "-")
				title = strings.ReplaceAll(title, "　", "-")
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, "--", "-")
				fmt.Println(title)

				filePath := "../zlzx.zjamr.zj.gov.cn/" + title + "(" + code + ")" + ".pdf"
				fmt.Println(filePath)

				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				viewZlZxHrefNode := htmlquery.FindOne(aNode, `./@href`)
				viewZlZxHref := htmlquery.InnerText(viewZlZxHrefNode)
				viewZlZxHref = "https://zlzx.zjamr.zj.gov.cn" + viewZlZxHref
				fmt.Println(viewZlZxHref)

				viewDoc, err := viewZlZxDoc(viewZlZxHref)
				if err != nil {
					fmt.Println("获取文档详情失败，跳过")
					continue
				}

				// 查看是否有下载按钮
				downloadButtonNode := htmlquery.FindOne(viewDoc, `//html/body[@class="zybj"]/div[@class="xwxqy"]/div[@class="gj-cx-c"]/div[@class="gj-zd"]/ul[@class="gj-bt"]/table/tbody/tr[2]/td[2]/a`)
				if downloadButtonNode == nil {
					fmt.Println("没有下载按钮，跳过")
					continue
				}

				// dbPdfView(&quot;https:\/\/zlzx.zjamr.zj.gov.cn\/bzzx\/rest\/redirect\/files\/localFile\/2025-07-09\/e2901fcd233d5674d1b8d0331df4d89a.pdf&quot;,&quot;db164bf3b31f40e0a05f0d8808c69cc3&quot;,&quot;DB33\/T 1436-2025&quot;)
				clickText := htmlquery.SelectAttr(downloadButtonNode, "onclick")
				clickText = strings.ReplaceAll(clickText, "dbPdfView(", "")
				clickText = strings.ReplaceAll(clickText, ")", "")
				clickText = strings.ReplaceAll(clickText, "\"", "")
				clickText = strings.ReplaceAll(clickText, "\\", "")
				clickTextArray := strings.Split(clickText, ",")
				fmt.Println("=======开始下载========")

				downloadUrl := clickTextArray[0]
				fmt.Println(downloadUrl)
				fmt.Println("=======开始下载" + title + "========")

				err = downloadZlZx(downloadUrl, viewZlZxHref, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "../zlzx.zjamr.zj.gov.cn", "../upload.doc88.com/temp-zlzx.zjamr.zj.gov.cn")
				err = copyZlZxFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")
				//DownLoadZlZxTimeSleep := 10
				DownLoadZlZxTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadZlZxTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadZlZxTimeSleep, "秒 倒计时", i, "秒===========")
				}
			}
			DownLoadZlZxPageTimeSleep := 10
			// DownLoadZlZxPageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadZlZxPageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadZlZxPageTimeSleep, "秒 倒计时", i, "秒===========")
			}
			page++
			if page > maxPage {
				isPageListGo = false
				break
			}
		}
	}
}

func QueryZlZxList(requestUrl string, referer string) (doc *html.Node, err error) {
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
	if ZlZxEnableHttpProxy {
		client = ZlZxSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", ZlZxCookie)
	req.Header.Set("Host", "zlzx.zjamr.zj.gov.cn")
	req.Header.Set("Origin", "https://zlzx.zjamr.zj.gov.cn")
	req.Header.Set("Referer", referer)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
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

func viewZlZxDoc(requestUrl string) (doc *html.Node, err error) {
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
	if ZlZxEnableHttpProxy {
		client = ZlZxSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", ZlZxCookie)
	req.Header.Set("Host", "www.hnbzw.com")
	req.Header.Set("Origin", "https://www.hnbzw.com")
	req.Header.Set("Referer", "https://www.hnbzw.com/Standard/StdSearch.aspx")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
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

func downloadZlZx(attachmentUrl string, referer string, filePath string) error {
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
	if ZlZxEnableHttpProxy {
		client = ZlZxSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", ZlZxCookie)
	req.Header.Set("Host", "zlzx.zjamr.zj.gov.cn")
	req.Header.Set("Referer", referer)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 如果访问失败 就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}

	// 创建一个文件用于保存
	fileDiv := filepath.Dir(filePath)
	if _, err = os.Stat(fileDiv); err != nil {
		if os.MkdirAll(fileDiv, 0o777) != nil {
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

func copyZlZxFile(src, dst string) (err error) {
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
