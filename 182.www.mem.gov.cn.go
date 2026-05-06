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
	MemEnableHttpProxy = false
	MemHttpProxyUrl    = "111.225.152.186:8089"
)

func MemSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(MemHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

var MemCookie = "Hm_lvt_7c3492d683dc7a90fd44bf8bfd57e50c=1778038119; HMACCOUNT=4E5B3419A3141A8E; Hm_lpvt_7c3492d683dc7a90fd44bf8bfd57e50c=1778038265"

// ychEduSpider 获取应急管理部标准文档
// @Title 获取应急管理部标准文档
// @Description https://www.mem.gov.cn/，获取应急管理部标准文档
func main() {
    page := 1
    isPageListGo := true
    for isPageListGo {
        requestListUrl := "https://www.mem.gov.cn/fw/flfgbz/bz/bzwb/index.shtml"
        referUrl := "https://www.mem.gov.cn/fw/flfgbz/bz/bzwb/index.shtml"
        if page > 0 {
            requestListUrl = fmt.Sprintf("https://www.mem.gov.cn/fw/flfgbz/bz/bzwb/index_%d.shtml", page)
        }
        if page >= 2 {
            referUrl = fmt.Sprintf("https://www.mem.gov.cn/fw/flfgbz/bz/bzwb/index_%d.shtml", page-1)
        }
        fmt.Println(requestListUrl)
        fmt.Println(referUrl)
        memBzListDoc, err := MemBzHtmlDoc(requestListUrl, referUrl)
        if err != nil {
            fmt.Println(err)
            break
        }
        liNodes := htmlquery.Find(memBzListDoc, `//li/a[@class="newttle"]`)
        if len(liNodes) >= 1 {
            for _, liNode := range liNodes {
                fmt.Println("=====================开始处理列表-分割线==========================")

                fmt.Println("=======page = " + strconv.Itoa(page) + "=========")
                downloadHrefNode := htmlquery.FindOne(liNode, `./@href`)
                downLoadUrl := htmlquery.InnerText(downloadHrefNode)
                if strings.Contains(downLoadUrl, ".pdf") {
                    // 查看downloadHref是否含有www.mem.gov.cn
                    if !strings.Contains(downLoadUrl, "www.mem.gov.cn") {
                        // 不含有www.mem.gov.cn，下载连接需要处理
                        bzDetailRequestUrlBiasTIndex := strings.LastIndex(requestListUrl, "/")
                        downLoadUrl = strings.Replace(downLoadUrl, ".", "", 1)
                        downLoadUrl = requestListUrl[:bzDetailRequestUrlBiasTIndex] + downLoadUrl
                    }
                    // 中文标题
                    chineseTitle := htmlquery.InnerText(liNode)
                    chineseTitle = chineseTitle[:len(chineseTitle) - 16]
                    chineseTitle = strings.TrimSpace(chineseTitle)
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
                    fmt.Println(chineseTitle)

                    filePath := "../www.mem.gov.cn/" + chineseTitle + ".pdf"
                    _, err = os.Stat(filePath)
                    if err == nil {
                        fmt.Println("文档已下载过，跳过")
                        continue
                    }

                    // 开始下载
                    fmt.Println("=======开始下载========")
                    err = downloadMem(downLoadUrl, requestListUrl, filePath)
                    if err != nil {
                        fmt.Println(err)
                        continue
                    }
                    //复制文件
                    tempFilePath := strings.ReplaceAll(filePath, "www.mem.gov.cn", "temp-hbba.sacinfo.org.cn")
                    err = copyMemFile(filePath, tempFilePath)
                    if err != nil {
                        fmt.Println(err)
                        continue
                    }
                    fmt.Println("=======完成下载========")

                    // 设置倒计时
                    DownLoadTMemTimeSleep := 10
                    for i := 1; i <= DownLoadTMemTimeSleep; i++ {
                        time.Sleep(time.Second)
                        fmt.Println("===page = "+strconv.Itoa(page)+"===title="+chineseTitle+"===========操作完成，", "暂停", DownLoadTMemTimeSleep, "秒，倒计时", i, "秒===========")
                    }
                }
            }
            DownLoadMemPageTimeSleep := 10
            // DownLoadMemPageTimeSleep := rand.Intn(5)
            for i := 1; i <= DownLoadMemPageTimeSleep; i++ {
                time.Sleep(time.Second)
                fmt.Println("===page = "+strconv.Itoa(page)+"========= 暂停", DownLoadMemPageTimeSleep, "秒 倒计时", i, "秒===========")
            }
            page++
        } else {
            page = 0
            isPageListGo = false
            break
        }
    }
}

func MemBzHtmlDoc(requestUrl string, referer string) (doc *html.Node, err error) {
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
	if MemEnableHttpProxy {
		client = MemSetHttpProxy()
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
	req.Header.Set("Cookie", MemCookie)
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	req.Header.Set("Host", "www.mem.gov.cn")
	req.Header.Set("Origin", "https://www.mem.gov.cn/")
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

func downloadMem(attachmentUrl string, referer string, filePath string) error {
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
	if MemEnableHttpProxy {
		client = MemSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", MemCookie)
	req.Header.Set("Host", "www.mem.gov.cn")
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

func copyMemFile(src, dst string) (err error) {
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
