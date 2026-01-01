package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
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
	JstEnableHttpProxy = false
	JstHttpProxyUrl    = "111.225.152.186:8089"
)

func JstSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(JstHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

var JstCookie = "wzws_sessionid=gWExNmRkYYJiY2RiYjKAMjIyLjEzNi44MC4xNjGgaJ1DGw=="

// ychEduSpider 获取吉林省工程建设地方标准文档
// @Title 获取吉林省工程建设地方标准文档
// @Description http://jst.jl.gov.cn/，获取吉林省工程建设地方标准文档
func main() {
	page := 0
	isPageListGo := true
	for isPageListGo {
		requestListUrl := "http://jst.jl.gov.cn/zwgk/dfnzh/jlsgcjsdfbzqwgk/index.html"
		referUrl := "http://jst.jl.gov.cn/zwgk/dfnzh/"
		if page > 0 {
			requestListUrl = fmt.Sprintf("http://jst.jl.gov.cn/zwgk/dfnzh/jlsgcjsdfbzqwgk/index_%d.html", page)
		}
		if page >= 2 {
			referUrl = fmt.Sprintf("http://jst.jl.gov.cn/zwgk/dfnzh/jlsgcjsdfbzqwgk/index_%d.html", page-1)
		}
		fmt.Println(requestListUrl)
		JstListDoc, err := JstHtmlDoc(requestListUrl, referUrl)
		if err != nil {
			fmt.Println(err)
			break
		}
		liNodes := htmlquery.Find(JstListDoc, `//html/body/div[3]/ul/li`)
		if len(liNodes) >= 1 {
			for _, liNode := range liNodes {
				// 中文标题
				chineseTitleNode := htmlquery.FindOne(liNode, `./a/@title`)
				if chineseTitleNode == nil {
					fmt.Println("标题不存在，跳过")
					continue
				}
				chineseTitle := htmlquery.InnerText(chineseTitleNode)
				chineseTitle = strings.TrimSpace(chineseTitle)
				chineseTitle = strings.ReplaceAll(chineseTitle, "已废止", "")
				chineseTitle = strings.ReplaceAll(chineseTitle, "/", "-")
				chineseTitle = strings.ReplaceAll(chineseTitle, "／", "-")
				chineseTitle = strings.ReplaceAll(chineseTitle, "　", "")
				chineseTitle = strings.ReplaceAll(chineseTitle, " ", "")
				chineseTitle = strings.ReplaceAll(chineseTitle, "：", ":")
				chineseTitle = strings.ReplaceAll(chineseTitle, "—", "-")
				chineseTitle = strings.ReplaceAll(chineseTitle, "－", "-")
				chineseTitle = strings.ReplaceAll(chineseTitle, "（", "(")
				chineseTitle = strings.ReplaceAll(chineseTitle, "）", ")")
				chineseTitle = strings.ReplaceAll(chineseTitle, "《", "")
				chineseTitle = strings.ReplaceAll(chineseTitle, "》", "")
				chineseTitle = strings.ReplaceAll(chineseTitle, "()", "")
				fmt.Println(chineseTitle)

				filePath := "../jst.jl.gov.cn/" + chineseTitle + ".pdf"
				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				detailUrlNode := htmlquery.FindOne(liNode, `./a/@href`)
				if detailUrlNode == nil {
					fmt.Println("没有文档详情链接，跳过")
					continue
				}
				requestListUrlArray := strings.Split(requestListUrl, "/index")
				detailUrl := requestListUrlArray[0] + strings.ReplaceAll(htmlquery.InnerText(detailUrlNode), "./", "/")
				fmt.Println(detailUrl)
				//os.Exit(1)

				JstDetailDoc, err := JstHtmlDoc(detailUrl, requestListUrl)
				//fmt.Println(htmlquery.InnerText(JstDetailDoc))
				//os.Exit(1)
				if err != nil {
					fmt.Println("获取文档详情失败，跳过")
					continue
				}

				bzDetailANode := htmlquery.FindOne(JstDetailDoc, `//html/body/div[@class="contents"]/div[@class="c_main"]/div[@id="detailCnt"]/dl[@class="annex"]/dd/span/a`)
				//fmt.Println(htmlquery.InnerText(bzDetailANode))
				//os.Exit(1)
				if bzDetailANode == nil {
					fmt.Println("没有附件链接，跳过")
					continue
				}
				bzDownloadHref := htmlquery.SelectAttr(bzDetailANode, "href")
				fmt.Println(bzDownloadHref)
				if strings.Contains(bzDownloadHref, ".pdf") == false {
					fmt.Println("附件不是pdf文件，跳过")
					continue
				}

				// 下载文档URL
				downLoadUrl := bzDownloadHref
				// 查看bzDownloadHref是否含有jst.jl.gov.cn
				if !strings.Contains(bzDownloadHref, "jst.jl.gov.cn") {
					// 不含有jst.jl.gov.cn，下载连接需要处理
					bzDetailRequestUrlBiasTIndex := strings.LastIndex(detailUrl, "/t")
					bzDownloadHref = strings.Replace(bzDownloadHref, ".", "", 1)
					downLoadUrl = detailUrl[:bzDetailRequestUrlBiasTIndex] + bzDownloadHref
				}
				fmt.Println(downLoadUrl)

				// 开始下载
				fmt.Println("=======开始下载========")
				err = downloadJst(downLoadUrl, detailUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "jst.jl.gov.cn", "temp-dbba.sacinfo.org.cn")
				err = copyJstFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======完成下载========")

				// 设置倒计时
				DownLoadTJstTimeSleep := 10
				for i := 1; i <= DownLoadTJstTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page = "+strconv.Itoa(page)+"===title="+chineseTitle+"===========操作完成，", "暂停", DownLoadTJstTimeSleep, "秒，倒计时", i, "秒===========")
				}
			}
			DownLoadJstPageTimeSleep := 10
			// DownLoadJstPageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadJstPageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page = "+strconv.Itoa(page)+"========= 暂停", DownLoadJstPageTimeSleep, "秒 倒计时", i, "秒===========")
			}
			page++
		} else {
			page = 0
			isPageListGo = false
			break
		}
	}
}

func JstHtmlDoc(requestUrl string, referer string) (doc *html.Node, err error) {
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
	if JstEnableHttpProxy {
		client = JstSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接
	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", JstCookie)
	req.Header.Set("Host", "jst.jl.gov.cn")
	req.Header.Set("Origin", "http://jst.jl.gov.cn/")
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

func downloadJst(attachmentUrl string, referer string, filePath string) error {
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
	if JstEnableHttpProxy {
		client = JstSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", JstCookie)
	req.Header.Set("Host", "jst.jl.gov.cn")
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

func copyJstFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer func(in *os.File) {
		err := in.Close()
		if err != nil {
			return
		}
	}(in)

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)

	_, err = io.Copy(out, in)
	return
}
