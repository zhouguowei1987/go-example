package main

import (
	"errors"
	"fmt"
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

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const (
	CsiscHbEnableHttpProxy = false
	CsiscHbHttpProxyUrl    = "111.225.152.186:8089"
)

func CsiscHbSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(CsiscHbHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

var CsiscHbCookie = "acw_tc=2760774217786375310827330ebed6ebff555898c30504aa788bd2254d5237"

// 获取资本市场标准网-行业标准
// @Title 获取资本市场标准网-行业标准
// @Description https://www.csisc.cn/ 获取资本市场标准网-行业标准
func main() {
	maxPage := 79
	page := 1
	isPageListGo := true
	for isPageListGo {
		requestUrl := "https://www.csisc.cn/zbscbzw/c100176/yfb_hb_list/code_0.shtml"
		refererUrl := "https://www.csisc.cn/zbscbzw/c100176/yfb_hb_list/code_0.shtml"
		if page >= 2 {
			requestUrl = fmt.Sprintf("https://www.csisc.cn/zbscbzw/c100176/yfb_hb_list/code_0_%d.shtml", page)
		}
		if page >= 3 {
			refererUrl = fmt.Sprintf("https://www.csisc.cn/zbscbzw/c100176/yfb_hb_list/code_0_%d.shtml", page-1)
		}
		fmt.Println(requestUrl)
		pageDoc, err := QueryCsiscHbHtml(requestUrl, refererUrl)
		if err != nil {
			fmt.Println(err)
			isPageListGo = false
			continue
		}
		if err != nil {
			fmt.Println(err)
		}
		liNodes := htmlquery.Find(pageDoc, `//li[@class="item"]`)
		if len(liNodes) <= 0 {
			isPageListGo = false
			break
		}

		for _, liNode := range liNodes {

			aHrefNode := htmlquery.FindOne(liNode, `./a/@href`)
			if aHrefNode == nil {
				continue
			}
			detailUrl := htmlquery.InnerText(aHrefNode)
			detailUrl = "https://www.csisc.cn" + detailUrl
			fmt.Println(detailUrl)

			detailDoc, err := QueryCsiscHbHtml(detailUrl, requestUrl)
			if err != nil {
				fmt.Println(err)
				continue
			}

			titleNode := htmlquery.FindOne(detailDoc, `//div[@class="w980 mb"]/div[@class="mainbox clearfix"]/div[@class="innerbox clearfix"]/div[@class="maincontent singlecontent article"]/div[@class="inbox"]/div[@class="page_content"]/div[@class="article-content"]/p[2]/font`)
			if titleNode == nil {
				fmt.Println("未找到标题节点，跳过")
				continue
			}
			title := strings.TrimSpace(htmlquery.InnerText(titleNode))
			title = strings.ReplaceAll(title, "标准中文名称：", "")
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "--", "-")
			fmt.Println(title)

			codeNode := htmlquery.FindOne(detailDoc, `//div[@class="w980 mb"]/div[@class="mainbox clearfix"]/div[@class="innerbox clearfix"]/div[@class="maincontent singlecontent article"]/div[@class="inbox"]/div[@class="page_content"]/div[@class="article-content"]/p[1]/font`)
			if codeNode == nil {
				fmt.Println("未找到标准号节点，跳过")
				continue
			}
			code := strings.TrimSpace(htmlquery.InnerText(codeNode))
			code = strings.ReplaceAll(code, "标准编号：", "")
			code = strings.ReplaceAll(code, "/", "-")
			code = strings.ReplaceAll(code, "—", "-")
			fmt.Println(code)

			filePath := "../www.csisc.cn/" + title + "(" + code + ")" + ".pdf"
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			detailDownloadHrefNode := htmlquery.FindOne(detailDoc, `//div[@class="w980 mb"]/div[@class="mainbox clearfix"]/div[@class="innerbox clearfix"]/div[@class="maincontent singlecontent article"]/div[@class="inbox"]/div[@class="page_content"]/div[@class="article-content"]/font/p[2]/font/a/@href`)
			if detailDownloadHrefNode == nil {
				fmt.Println("未找到标准号节点，跳过")
				continue
			}
			downloadUrl := htmlquery.InnerText(detailDownloadHrefNode)
			detailUrlArray := strings.Split(detailUrl, "/")
			downloadUrl = strings.Join(detailUrlArray[:len(detailUrlArray)-1], "/") + "/" + downloadUrl
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")
			err = downloadCsiscHb(downloadUrl, detailUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "www.csisc.cn", "temp-hbba.sacinfo.org.cn")
			err = copyCsiscHbFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			//DownLoadCsiscHbTimeSleep := 10
			DownLoadCsiscHbTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadCsiscHbTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadCsiscHbTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadCsiscHbPageTimeSleep := 10
		// DownLoadCsiscHbPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadCsiscHbPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadCsiscHbPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

func QueryCsiscHbHtml(requestUrl string, referer string) (doc *html.Node, err error) {
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
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", CsiscHbCookie)
	req.Header.Set("Host", "www.csisc.cn")
	req.Header.Set("Origin", "https://www.csisc.cn")
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

func downloadCsiscHb(attachmentUrl string, referer string, filePath string) error {
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
	if CsiscHbEnableHttpProxy {
		client = CsiscHbSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.csisc.cn")
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

func copyCsiscHbFile(src, dst string) (err error) {
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
