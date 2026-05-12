package main

import (
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
	"golang.org/x/net/html"
)

const (
	ShanXiEnableHttpProxy = false
	ShanXiHttpProxyUrl    = "111.225.152.186:8089"
)

func ShanXiSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(ShanXiHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

var ShanXiCookie = "trs_uv=mp25ydum_5381_1nfl; _trs_ua_s_1=mp28edg1_5381_61si"

// ychEduSpider 获取山西省市场监督管理局标准文档
// @Title 获取山西省市场监督管理局标准文档
// @Description https://scjgj.shanxi.gov.cn/，获取山西省市场监督管理局标准文档
func main() {
	page := 0
	isPageListGo := true
	for isPageListGo {
		requestListUrl := "https://scjgj.shanxi.gov.cn/gzzt/bzgl/dfbzcx/index.shtml"
		referUrl := "https://scjgj.shanxi.gov.cn/gzzt/bzgl/dfbzcx/index.shtml"
		if page >= 1 {
			requestListUrl = fmt.Sprintf("https://scjgj.shanxi.gov.cn/gzzt/bzgl/dfbzcx/index_%d.shtml", page)
		}
		if page >= 2 {
			referUrl = fmt.Sprintf("https://scjgj.shanxi.gov.cn/gzzt/bzgl/dfbzcx/index_%d.shtml", page-1)
		}
		fmt.Println(requestListUrl)
		fmt.Println(referUrl)
		saacBzListDoc, err := ShanXiBzHtmlDoc(requestListUrl, referUrl)
		if err != nil {
			fmt.Println(err)
			break
		}
		liNodes := htmlquery.Find(saacBzListDoc, `//html/body/div/div[2]/div/div[2]/ul/ul/li`)
		if len(liNodes) >= 1 {
			for _, liNode := range liNodes {
				fmt.Println("=====================开始处理列表-分割线==========================")

				fmt.Println("=======page = " + strconv.Itoa(page) + "=========")

				// 中文标题
				aHrefNode := htmlquery.FindOne(liNode, `./a/span[@class="li_con"]`)
				title := htmlquery.InnerText(aHrefNode)
				title = strings.TrimSpace(title)
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, "／", "-")
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, "　", "-")
				title = strings.ReplaceAll(title, " ", "-")
				title = strings.ReplaceAll(title, "：", "-")
				title = strings.ReplaceAll(title, "--", "-")
				title = strings.ReplaceAll(title, "—", "-")
				title = strings.ReplaceAll(title, "• ", "")
				title = strings.ReplaceAll(title, " ", "")
				title = strings.ReplaceAll(title, ".pdf", "")
				title = strings.ReplaceAll(title, "（", "(")
				title = strings.ReplaceAll(title, "）", ")")
				title = strings.ReplaceAll(title, "《", "")
				title = strings.ReplaceAll(title, "》", "")
				fmt.Println(title)

				filePath := "../scjgj.shanxi.gov.cn/" + title + ".pdf"
				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				downloadHrefNode := htmlquery.FindOne(liNode, `./a/@href`)
				downLoadUrl := htmlquery.InnerText(downloadHrefNode)
				if strings.Contains(downLoadUrl, ".pdf") {
					// 查看downloadHref是否含有scjgj.shanxi.gov.cn
					if !strings.Contains(downLoadUrl, "scjgj.shanxi.gov.cn") {
						// 不含有scjgj.shanxi.gov.cn，下载连接需要处理
						// 查看有几个“..”
						count := strings.Count(downLoadUrl, "..")
						if count > 0 {
							requestListUrlArray := strings.Split(requestListUrl, "/")
							requestListUrlArray = requestListUrlArray[:len(requestListUrlArray)-(count+1)]
							requestListUrlDownLoadUrl := strings.Join(requestListUrlArray, "/")
							downLoadUrl = strings.Replace(downLoadUrl, "../", "", count)
							downLoadUrl = requestListUrlDownLoadUrl + "/" + downLoadUrl
						} else {
							bzDetailRequestUrlBiasTIndex := strings.LastIndex(requestListUrl, "/")
							downLoadUrl = strings.Replace(downLoadUrl, ".", "", 1)
							downLoadUrl = requestListUrl[:bzDetailRequestUrlBiasTIndex] + downLoadUrl
						}
					}
					fmt.Println(downLoadUrl)
					// 开始下载
					fmt.Println("=======开始下载========")
					err = downloadShanXi(downLoadUrl, requestListUrl, filePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					//复制文件
					tempFilePath := strings.ReplaceAll(filePath, "scjgj.shanxi.gov.cn", "temp-hbba.sacinfo.org.cn")
					err = copyShanXiFile(filePath, tempFilePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("=======完成下载========")

					// 设置倒计时
					DownLoadTShanXiTimeSleep := 10
					for i := 1; i <= DownLoadTShanXiTimeSleep; i++ {
						time.Sleep(time.Second)
						fmt.Println("===page = "+strconv.Itoa(page)+"===title="+title+"===========操作完成，", "暂停", DownLoadTShanXiTimeSleep, "秒，倒计时", i, "秒===========")
					}
				}
			}
			DownLoadShanXiPageTimeSleep := 10
			// DownLoadShanXiPageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadShanXiPageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("===page = "+strconv.Itoa(page)+"========= 暂停", DownLoadShanXiPageTimeSleep, "秒 倒计时", i, "秒===========")
			}
			page++
		} else {
			page = 0
			isPageListGo = false
			break
		}
	}
}

func ShanXiBzHtmlDoc(requestUrl string, referer string) (doc *html.Node, err error) {
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
	if ShanXiEnableHttpProxy {
		client = ShanXiSetHttpProxy()
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
	req.Header.Set("Cookie", ShanXiCookie)
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	req.Header.Set("Host", "scjgj.shanxi.gov.cn")
	req.Header.Set("Origin", "https://scjgj.shanxi.gov.cn/")
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

func downloadShanXi(attachmentUrl string, referer string, filePath string) error {
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
	if ShanXiEnableHttpProxy {
		client = ShanXiSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", ShanXiCookie)
	req.Header.Set("Host", "scjgj.shanxi.gov.cn")
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

func copyShanXiFile(src, dst string) (err error) {
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
