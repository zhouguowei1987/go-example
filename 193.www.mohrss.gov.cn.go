package main

import (
	"errors"
	"fmt"
	"io"

	// 	"math/rand"
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
	MoHrSsEnableHttpProxy = false
	MoHrSsHttpProxyUrl    = "111.225.152.186:8089"
)

func MoHrSsSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(MoHrSsHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

var MoHrSsCookie = "path=/; __tst_status=878897575#; EO_Bot_Ssid=2891841536; Hm_lvt_64e46e3f389bd47c0981fa5e4b9f2405=1778485750,1778570141,1778652195; HMACCOUNT=9C0CD19686802BBF; arialoadData=false; ariauseGraymode=false; Hm_lpvt_64e46e3f389bd47c0981fa5e4b9f2405=1778652312"

// 获取中华人民共和国人力资源和社会保障部标准
// @Title 获取中华人民共和国人力资源和社会保障部标准
// @Description https://www.mohrss.gov.cn/ 获取中华人民共和国人力资源和社会保障部标准
func main() {
	maxPage := 3
	page := 1
	isPageListGo := true
	for isPageListGo {
		requestUrl := "https://www.mohrss.gov.cn/SYrlzyhshbzb/zwgk/ghcw/bz/index.html"
		refererUrl := "https://www.mohrss.gov.cn/SYrlzyhshbzb/zwgk/ghcw/bz/index.html"
		if page >= 2 {
			requestUrl = fmt.Sprintf("https://www.mohrss.gov.cn/SYrlzyhshbzb/zwgk/ghcw/bz/index_%d.html", page-1)
		}
		if page >= 3 {
			refererUrl = fmt.Sprintf("https://www.mohrss.gov.cn/SYrlzyhshbzb/zwgk/ghcw/bz/index_%d.html", page-2)
		}
		fmt.Println(requestUrl)
		pageDoc, err := QueryMoHrSsHtml(requestUrl, refererUrl)
		if err != nil {
			fmt.Println(err)
			isPageListGo = false
			continue
		}
		if err != nil {
			fmt.Println(err)
		}
		liNodes := htmlquery.Find(pageDoc, `//html/body/div[1]/div/div/div[2]/div[1]/div[2]/div`)
		if len(liNodes) <= 0 {
			isPageListGo = false
			break
		}

		for _, liNode := range liNodes {

			titleNode := htmlquery.FindOne(liNode, `//a`)
			if titleNode == nil {
				fmt.Println("未找到标题节点，跳过")
				continue
			}
			title := strings.TrimSpace(htmlquery.InnerText(titleNode))
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "／", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "：", "-")
			title = strings.ReplaceAll(title, "—", "-")
			title = strings.ReplaceAll(title, "--", "-")
			title = strings.ReplaceAll(title, ".pdf", "")
			title = strings.ReplaceAll(title, "（", "(")
			title = strings.ReplaceAll(title, "）", ")")
			title = strings.ReplaceAll(title, "《", "")
			title = strings.ReplaceAll(title, "》", "")
			fmt.Println(title)
			if strings.Index(title, "LD") == -1 {
				fmt.Println("不含‘LD’字样，跳过")
				continue
			}

			filePath := "../www.mohrss.gov.cn/" + title + ".pdf"
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			aHrefNode := htmlquery.FindOne(liNode, `//a/@href`)
			if aHrefNode == nil {
				fmt.Println("未找到链接节点，跳过")
				continue
			}
			detailUrl := htmlquery.InnerText(aHrefNode)
			requestUrlArray := strings.Split(requestUrl, "/")
			detailUrl = strings.Join(requestUrlArray[:len(requestUrlArray)-1], "/") + "/" + detailUrl
			fmt.Println(detailUrl)

			detailDoc, err := QueryMoHrSsHtml(detailUrl, requestUrl)
			if err != nil {
				fmt.Println(err)
				continue
			}

			detailEditorNode := htmlquery.FindOne(detailDoc, `//div[@class="TRS_Editor"]`)
			if detailEditorNode == nil {
				fmt.Println("未找到‘TRS_Editor’文件节点，跳过")
				continue
			}

			detailDownloadHrefNode := htmlquery.FindOne(detailEditorNode, `//a/@href`)
			if detailDownloadHrefNode == nil {
				fmt.Println("未找到下载文件节点，跳过")
				continue
			}
			downloadUrl := htmlquery.InnerText(detailDownloadHrefNode)
			downloadUrl = strings.ReplaceAll(downloadUrl, "./", "")
			detailUrlArray := strings.Split(detailUrl, "/")
			downloadUrl = strings.Join(detailUrlArray[:len(detailUrlArray)-1], "/") + "/" + downloadUrl
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")
			err = downloadMoHrSs(downloadUrl, detailUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "www.mohrss.gov.cn", "temp-hbba.sacinfo.org.cn")
			err = copyMoHrSsFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			DownLoadMoHrSsTimeSleep := 10
			// DownLoadMoHrSsTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadMoHrSsTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadMoHrSsTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadMoHrSsPageTimeSleep := 10
		// DownLoadMoHrSsPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadMoHrSsPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadMoHrSsPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

func QueryMoHrSsHtml(requestUrl string, referer string) (doc *html.Node, err error) {
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
	req.Header.Set("Cookie", MoHrSsCookie)
	req.Header.Set("Host", "www.mohrss.gov.cn")
	req.Header.Set("Origin", "https://www.mohrss.gov.cn")
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

func downloadMoHrSs(attachmentUrl string, referer string, filePath string) error {
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
	if MoHrSsEnableHttpProxy {
		client = MoHrSsSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.mohrss.gov.cn")
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

func copyMoHrSsFile(src, dst string) (err error) {
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
