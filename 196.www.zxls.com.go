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
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const (
	ZxLsEnableHttpProxy = false
	ZxLsHttpProxyUrl    = "111.225.152.186:8089"
)

func ZxLsSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(ZxLsHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

var ZxLsCookie = "Hm_lvt_c546156b33a73aaa69021f8a527d9e26=1780378702,1780551837,1780623585; HMACCOUNT=9C0CD19686802BBF; ASP.NET_SessionId=e5n1ox454px4b1fmlc0ae545; last_heartbeat_time=1780626778891; Hm_lpvt_c546156b33a73aaa69021f8a527d9e26=1780626816"

// 获取中学历史教学园地试卷
// @Title 获取中学历史教学园地试卷
// @Description https://www.zxls.com/ 获取中学历史教学园地试卷
func main() {
	maxPage := 601
	page := 1
	isPageListGo := true
	for isPageListGo {
		requestUrl := fmt.Sprintf("https://www.zxls.com/generation/gzst/List_5884_%d.html", page)
		refererUrl := "https://www.zxls.com/generation/gzst/List_5884_1.html"
		pathUrl := fmt.Sprintf("/generation/gzst/List_5884_%d.html", page)
		if page >= 2 {
			refererUrl = fmt.Sprintf("https://www.zxls.com/generation/gzst/List_5884_%d.html", page-1)
		}
		fmt.Println(requestUrl)
		pageDoc, err := QueryZxLsHtml(requestUrl, refererUrl, pathUrl)
		if err != nil {
			fmt.Println(err)
			isPageListGo = false
			continue
		}
		if err != nil {
			fmt.Println(err)
		}
		liNodes := htmlquery.Find(pageDoc, `//html/body/div[2]/div[2]/div[2]/div[5]/div[4]/ul[@id="infolist"]/li[@class="list-vh-item"]`)
		if len(liNodes) <= 0 {
			isPageListGo = false
			break
		}

		for _, liNode := range liNodes {

			// 查看是否免费
			pointsNode := htmlquery.FindOne(liNode, `./span[@class="point mx-point"]/@data-point`)
			if pointsNode == nil {
				fmt.Println("没有点span")
				continue
			}
			pointsText := htmlquery.InnerText(pointsNode)

			points, err := strconv.Atoi(pointsText)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if points > 0 {
				fmt.Println("需要点下载", points)
				continue
			}

			detailHrefNode := htmlquery.FindOne(liNode, `./a/@href`)
			if detailHrefNode == nil {
				fmt.Println("未找到链接节点，跳过")
				continue
			}
			detailHref := strings.TrimSpace(htmlquery.InnerText(detailHrefNode))
			detailHref = strings.TrimSpace(detailHref)

			detailHrefArray := strings.Split(detailHref, "/")
			downloadId := strings.ReplaceAll(detailHrefArray[len(detailHrefArray)-1], ".html", "")

			showDownloadUrl := fmt.Sprintf("https://www.zxls.com/Common/ShowDownloadUrl.aspx?id=%s", downloadId)
			showDownloadPathUrl := fmt.Sprintf("/Common/ShowDownloadUrl.aspx?id=%s", downloadId)
			fmt.Println(showDownloadUrl)
			showDownloadDoc, err := QueryZxLsHtml(showDownloadUrl, "https://www.zxls.com"+detailHref, showDownloadPathUrl)
			if err != nil {
				fmt.Println(err)
				continue
			}
			titleNode := htmlquery.FindOne(showDownloadDoc, `//html/body/div[3]/div[2]/div[1]/div[2]/div[5]/h2`)
			if titleNode == nil {
				fmt.Println("未找到标题节点，跳过")
				continue
			}
			title := strings.TrimSpace(htmlquery.InnerText(titleNode))
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, "当前资料-", "")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "／", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "：", "-")
			title = strings.ReplaceAll(title, "—", "-")
			title = strings.ReplaceAll(title, "--", "-")
			title = strings.ReplaceAll(title, "（", "(")
			title = strings.ReplaceAll(title, "）", ")")
			title = strings.ReplaceAll(title, "《", "")
			title = strings.ReplaceAll(title, "》", "")
			fmt.Println(title)

			downloadHrefNode := htmlquery.FindOne(showDownloadDoc, `//html/body/div[3]/div[2]/div[1]/div[2]/div[6]/form/div[3]/dl/dd[1]/div[1]/div[2]/div[1]/span[1]/a/@href`)
			if downloadHrefNode == nil {
				fmt.Println("未找到下载文件节点，跳过")
				// continue
			}
			downloadUrl := htmlquery.InnerText(downloadHrefNode)
			fileExt := path.Ext(downloadUrl)
			downloadUrl = "https://www.zxls.com" + downloadUrl
			fmt.Println(downloadUrl)

			filePath := "../www.zxls.com/www.zxls.com/" + title + fileExt
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载" + title + "========")
			err = downloadZxLs(downloadUrl, showDownloadUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "www.zxls.com/www.zxls.com", "www.zxls.com/temp-www.zxls.com")
			err = copyZxLsFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			DownLoadZxLsTimeSleep := 10
			// DownLoadZxLsTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadZxLsTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadZxLsTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadZxLsPageTimeSleep := 10
		// DownLoadZxLsPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadZxLsPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadZxLsPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

func QueryZxLsHtml(requestUrl string, referer string, path string) (doc *html.Node, err error) {
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
	req.Header.Set("authority", "www.zxls.com")
	req.Header.Set("method", "GET")
	req.Header.Set("path", path)
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", ZxLsCookie)
	req.Header.Set("Host", "www.zxls.com")
	req.Header.Set("Origin", "https://www.zxls.com")
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
	return doc, nil
}

func downloadZxLs(attachmentUrl string, referer string, filePath string) error {
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
	if ZxLsEnableHttpProxy {
		client = ZxLsSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.zxls.com")
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

func copyZxLsFile(src, dst string) (err error) {
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
