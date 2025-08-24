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
	_ "os"
	"path/filepath"
	_ "path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/antchfx/htmlquery"
	_ "golang.org/x/net/html"
)

var BjJsEnableHttpProxy = false
var BjJsHttpProxyUrl = "111.225.152.186:8089"
var BjJsHttpProxyUrlArr = make([]string, 0)

func BjJsHttpProxy() error {
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
					BjJsHttpProxyUrlArr = append(BjJsHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					BjJsHttpProxyUrlArr = append(BjJsHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func BjJsSetHttpProxy() (httpclient *http.Client) {
	if BjJsHttpProxyUrl == "" {
		if len(BjJsHttpProxyUrlArr) <= 0 {
			err := BjJsHttpProxy()
			if err != nil {
				BjJsSetHttpProxy()
			}
		}
		BjJsHttpProxyUrl = BjJsHttpProxyUrlArr[0]
		if len(BjJsHttpProxyUrlArr) >= 2 {
			BjJsHttpProxyUrlArr = BjJsHttpProxyUrlArr[1:]
		} else {
			BjJsHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(BjJsHttpProxyUrl)
	ProxyURL, _ := url.Parse(BjJsHttpProxyUrl)
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

type QueryBjJsListFormData struct {
	filter_LIKE_name         string
	filter_LIKE_standard_num string
	currentPage              int
	pageSize                 int
	recordCount              int
	OrderByField             string
	OrderByDesc              string
}

var BjJsCookie = "JSESSIONID=BB2E476658D59EB2D06117C807C5445B; history={title:\"ä¸œåŸŽåŒºé\u009D©æ–°é‡Œ26ã€\u008128å\u008F·ä½\u008Få®…é¡¹ç›®å»ºè®¾æ–¹æ¡ˆ\",id:\"53613709\",url:\"/bjjs/sy/zwfw66/cxzx-sjzt/gcjsl/xmjsfa/dcqxmjsfags/53613709/index.shtml\",date:\"2025-08-22\"}"

// 下载北京市住房和城乡建设委员会文档
// @Title 下载北京市住房和城乡建设委员会文档
// @Description http://bjjs.zjw.beijing.gov.cn/，下载北京市住房和城乡建设委员会文档
func main() {
	pageListUrl := "http://bjjs.zjw.beijing.gov.cn/eportal/ui?pageId=53613181"
	fmt.Println(pageListUrl)
	page := 1
	maxPage := 27
	rows := 10
	totalRecord := 262
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		queryBjJsListFormData := QueryBjJsListFormData{
			filter_LIKE_name:         "",
			filter_LIKE_standard_num: "",
			currentPage:              page,
			pageSize:                 rows,
			recordCount:              totalRecord,
			OrderByField:             "",
			OrderByDesc:              "",
		}
		queryBjJsListDoc, err := QueryBjJsList(pageListUrl, queryBjJsListFormData)
		if err != nil {
			fmt.Println(err)
			break
		}

		trNodes := htmlquery.Find(queryBjJsListDoc, `//div[@class="xxgk_content"]/div[@class="xxjs column"]/div[@class="portlet"]/div[2]/form/div[@class="bzcx_"]/div[1]/table/tbody/tr`)
		if len(trNodes) >= 2 {
			for trIndex, trNode := range trNodes {
				if trIndex == 0 {
					fmt.Println("表头数据，跳过")
					continue
				}
				fmt.Println("=====================开始处理数据 page = ", page, "=========================")

				codeNode := htmlquery.FindOne(trNode, `./td[1]`)
				if codeNode == nil {
					fmt.Println("标准号不存在，跳过")
					continue
				}
				code := strings.TrimSpace(htmlquery.InnerText(codeNode))
				code = strings.ReplaceAll(code, "/", "-")
				code = strings.ReplaceAll(code, "—", "-")
				fmt.Println(code)

				titleNode := htmlquery.FindOne(trNode, `./td[2]/a/@title`)
				if titleNode == nil {
					fmt.Println("标题不存在，跳过")
					continue
				}
				title := strings.TrimSpace(htmlquery.InnerText(titleNode))
				title = strings.ReplaceAll(title, " ", "-")
				title = strings.ReplaceAll(title, "　", "-")
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, "--", "-")
				fmt.Println(title)

				filePath := "../bjjs.zjw.beijing.gov.cn/" + title + "(" + code + ").pdf"
				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				viewUrlNode := htmlquery.FindOne(trNode, `./td[2]/a/@href`)
				if viewUrlNode == nil {
					fmt.Println("预览链接不存在，跳过")
					continue
				}
				viewUrl := htmlquery.InnerText(viewUrlNode)
				viewDoc, err := htmlquery.LoadURL(viewUrl)
				if err != nil {
					fmt.Println(err)
					continue
				}

				reg := regexp.MustCompile("var path2 = \"(.*?)\";")
				path2 := reg.Find([]byte(htmlquery.InnerText(viewDoc)))
				path2Str := string(path2)
				path2StrHandle := strings.ReplaceAll(path2Str, "var path2 = \"", "")
				path2StrHandle = strings.ReplaceAll(path2StrHandle, "\";", "")

				downloadUrl := "http://111.206.112.79:8080/fileTest/" + path2StrHandle

				fmt.Println("=======开始下载" + title + "========")

				err = downloadBjJs(downloadUrl, pageListUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "../bjjs.zjw.beijing.gov.cn", "../upload.doc88.com/dbba.sacinfo.org.cn")
				err = copyBjJsFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")
				//DownLoadBjJsTimeSleep := 10
				DownLoadBjJsTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadBjJsTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadBjJsTimeSleep, "秒 倒计时", i, "秒===========")
				}
			}
		}
		DownLoadBjJsPageTimeSleep := 10
		// DownLoadBjJsPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadBjJsPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadBjJsPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

func QueryBjJsList(requestUrl string, queryBjJsListFormData QueryBjJsListFormData) (doc *html.Node, err error) {
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
	if BjJsEnableHttpProxy {
		client = BjJsSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("filter_LIKE_name", queryBjJsListFormData.filter_LIKE_name)
	postData.Add("filter_LIKE_standard_num", queryBjJsListFormData.filter_LIKE_standard_num)
	postData.Add("currentPage", strconv.Itoa(queryBjJsListFormData.currentPage))
	postData.Add("pageSize", strconv.Itoa(queryBjJsListFormData.pageSize))
	postData.Add("recordCount", strconv.Itoa(queryBjJsListFormData.recordCount))
	postData.Add("OrderByField", queryBjJsListFormData.OrderByField)
	postData.Add("OrderByDesc", queryBjJsListFormData.OrderByDesc)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", BjJsCookie)
	req.Header.Set("Host", "bjjs.zjw.beijing.gov.cn")
	req.Header.Set("Origin", "http://bjjs.zjw.beijing.gov.cn")
	req.Header.Set("Referer", "http://bjjs.zjw.beijing.gov.cn/eportal/ui?pageId=53613181")
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

func downloadBjJs(attachmentUrl string, referer string, filePath string) error {
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
	if BjJsEnableHttpProxy {
		client = BjJsSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "bjjs.zjw.beijing.gov.cn")
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

func copyBjJsFile(src, dst string) (err error) {
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
