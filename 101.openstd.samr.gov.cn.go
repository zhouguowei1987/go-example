package main

import (
	"errors"
	"fmt"
	"strings"
	"math/rand"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	OPenStdEnableHttpProxy = false
	OPenStdHttpProxyUrl    = "111.225.152.186:8089"
)

func OPenStdSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(OPenStdHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type StdCategory struct {
	StdName string
	StdUrl  string
	Page  int
	MaxPage  int
}

var OPenStdCookie = ""

// ychEduSpider 获取国家标准文档
// @Title 获取国家标准文档
// @Description https://openstd.samr.gov.cn/，获取国家标准文档
func main() {
	var StdCategories = []StdCategory{
// 		{
// 			StdName: "强制性国家标准",
// 			StdUrl:  "https://openstd.samr.gov.cn/bzgk/std/std_list_type?p.p1=1&p.p90=circulation_date&p.p91=desc",
// 			Page: 37,
// 			MaxPage: 612,
// 		},
		{
			StdName: "推荐性国家标准",
			StdUrl:  "https://openstd.samr.gov.cn/bzgk/std/std_list_type?p.p1=2&p.p90=circulation_date&p.p91=desc",
			Page: 81,
			MaxPage: 6037,
		},
// 		{
// 			StdName: "指导性技术文件",
// 			StdUrl:  "https://openstd.samr.gov.cn/bzgk/std/std_list_type?p.p1=3&p.p90=circulation_date&p.p91=desc",
// 			Page: 55,
// 			MaxPage: 83,
// 		},
	}
	for _, std := range StdCategories {
		page := std.Page
		pageSize := 10
		isPageListGo := true
		for isPageListGo {
            if page > std.MaxPage {
                isPageListGo = false
                break
            }
			pageListUrl := fmt.Sprintf(std.StdUrl+"&r=0.20175458803007884&page=%d&pageSize=%d", page, pageSize)
			pageListRefererUrl := std.StdUrl
			if page > 1{
			    pageListRefererUrl = fmt.Sprintf(std.StdUrl+"&r=0.20175458803007884&page=%d&pageSize=%d", page-1, pageSize)
			}
			fmt.Println(pageListUrl)
			// 每个分页清空cookie
			OPenStdCookie = ""
			pageListDoc, err := oPenStdHtml(pageListUrl, pageListRefererUrl)
			if err != nil {
				fmt.Println(err)
				break
			}
			trNodes := htmlquery.Find(pageListDoc, `//table[@class="table result_list table-striped table-hover"]/tbody[2]/tr`)
			if len(trNodes) > 0 {
				for _, trNode := range trNodes {
                    StdNoA := htmlquery.FindOne(trNode, `./td[2]/a`)
                    StdNo := htmlquery.InnerText(StdNoA)
                    StdNo = strings.TrimSpace(StdNo)
                    StdNo = strings.ReplaceAll(StdNo, "/", "-")
                    fmt.Println(StdNo)

                    StdNameA := htmlquery.FindOne(trNode, `./td[4]/a`)
                    StdName := htmlquery.InnerText(StdNameA)
                    StdName = strings.TrimSpace(StdName)
                    StdName = strings.ReplaceAll(StdName, " ", "")
                    StdName = strings.ReplaceAll(StdName, "　", "")
                    StdName = strings.ReplaceAll(StdName, "/", "-")
                    StdName = strings.ReplaceAll(StdName, "：", "-")
                    StdName = strings.ReplaceAll(StdName, "《", "")
                    StdName = strings.ReplaceAll(StdName, "》", "")
                    StdName = strings.ReplaceAll(StdName, "--", "-")
                    StdName = strings.ReplaceAll(StdName, "——", "-")
                    fmt.Println(StdName)

                    HCno := htmlquery.SelectAttr(StdNameA, "onclick")
                    HCno = HCno[10 : len(HCno)-3]
                    fmt.Println(HCno)

                    // 详情URL
                    detailUrl := fmt.Sprintf("https://openstd.samr.gov.cn/bzgk/std/newGbInfo?hcno=%s", HCno)
                    fmt.Println(detailUrl)
                    detailDoc, err := oPenStdHtml(detailUrl, pageListUrl)
                    if err != nil {
                        fmt.Println(err)
                        continue
                    }
                    // 查看是否有下载按钮
                    xzBtnNode := htmlquery.FindOne(detailDoc, `//html/body/div[3]/div/div/div/div/table[2]/tbody/tr[4]/td/button[@class="btn xz_btn btn-sm btn-warning app-hide"]`)
                    if xzBtnNode == nil{
                        fmt.Println("没有下载按钮，跳过")
                        continue
                    }

                    // 下载文档URL
                    downLoadUrl := fmt.Sprintf("https://openstd.samr.gov.cn/bzgk/std/viewGb?hcno=%s", HCno)
                    fmt.Println(downLoadUrl)

                    refererUrl := fmt.Sprintf("https://openstd.samr.gov.cn/bzgk/std/showGb?type=download&hcno=%s&request_locale=zh", HCno)
                    fmt.Println(refererUrl)

                    filePath := "../openstd.samr.gov.cn/" + StdName + "(" + StdNo + ")" + ".pdf"
                    _, err = os.Stat(filePath)
                    if err == nil {
                        fmt.Println("文档已下载过，跳过")
                        continue
                    }
                    fmt.Println("=======开始下载========")
                    // 下载详情页面和文件页面要同时请求，才可以下载温年
                    _, err = oPenStdHtml(refererUrl, detailUrl)
                    err = downloadOPenStd(downLoadUrl, refererUrl, filePath)
                    if err != nil {
                        fmt.Println(err)
                        continue
                    }
                    // 查看文件大小，如果是空文件，则删除
                    fileInfo, err := os.Stat(filePath)
                    if err == nil && fileInfo.Size() == 0{
                        fmt.Println("空文件删除")
                        err = os.Remove(filePath)
                    }
                    if err != nil {
                        continue
                    }
                    //复制文件
                    tempFilePath := strings.ReplaceAll(filePath, "openstd.samr.gov.cn", "temp-hbba.sacinfo.org.cn")
                    err = copyOPenStdFile(filePath, tempFilePath)
                    if err != nil {
                        fmt.Println(err)
                        continue
                    }
                    fmt.Println("=======开始完成========")
                    // 设置倒计时
                    DownLoadOPenStdTimeSleep := rand.Intn(80)
                    if DownLoadOPenStdTimeSleep <= 60{
                        DownLoadOPenStdTimeSleep = 60 + rand.Intn(30)
                    }
                    for i := 1; i <= DownLoadOPenStdTimeSleep; i++ {
                        time.Sleep(time.Second)
                        fmt.Println("page = "+strconv.Itoa(page)+"===StdName="+StdName+"===========操作完成，", "暂停", DownLoadOPenStdTimeSleep, "秒，倒计时", i, "秒===========")
                    }
                }
			}
			DownLoadOPenStdPageTimeSleep := 10
            // DownLoadOPenStdPageTimeSleep := rand.Intn(5)
            for i := 1; i <= DownLoadOPenStdPageTimeSleep; i++ {
                time.Sleep(time.Second)
                fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadOPenStdPageTimeSleep, "秒 倒计时", i, "秒===========")
            }
			page++
			if page > std.MaxPage {
                isPageListGo = false
                break
            }
		}
	}
}


func oPenStdHtml(requestUrl string, referer string) (doc *html.Node, err error) {
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
	req.Header.Set("Cookie", OPenStdCookie)
	req.Header.Set("Host", "openstd.samr.gov.cn")
	req.Header.Set("Origin", "https://openstd.samr.gov.cn")
	req.Header.Set("Referer", referer)
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
	if len(OPenStdCookie) <= 0{
	    OPenStdCookie = resp.Header.Get("Set-Cookie")
	}
	return doc, nil
}

func downloadOPenStd(attachmentUrl string, referer string, filePath string) error {
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
	if OPenStdEnableHttpProxy {
		client = OPenStdSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Range", "bytes 0-")
	fmt.Println(OPenStdCookie)
	req.Header.Set("Cookie", OPenStdCookie)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "openstd.samr.gov.cn")
// 	req.Header.Set("Origin", "https://openstd.samr.gov.cn/")
	req.Header.Set("Referer", referer)
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"148\", \"Google Chrome\";v=\"148\", \"Not/A)Brand\";v=\"99\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")
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
func copyOPenStdFile(src, dst string) (err error) {
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