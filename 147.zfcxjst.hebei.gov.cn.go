package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type zfCxJstCategory struct {
	url     string
	name    string
	page    int
	maxPage int
}

var ZfCxJstCookie = "qspJd7aG3Y1XO=60fUHO9ocKVEtR8xeDuxT.qFr5J7jn5RfHlXxJP7md0y65aCBzcwx2z5Qr7cdOFg2rDD4j8V4_9YIAsUyvxlCGlG; qspJd7aG3Y1XP=0iPjgelfCeW5yTTCKflfL_a6u8.ZOd7MVpIK9Jl5U9N.wyFkJ1N19iCfu6zIzA7bcRwsWR3ZmDL6dXsLWun6hoLtDVve8ZHXtNssPNBeGuT_bXVU2G3sYHRA99KlnaojaaL0_JWDk2.VmgVfOniVWGOojAohcRU9c3uJiRf_gtSbXoUEkUQwC7pSjxYioNsqBSbOZ.3N2QceqDFM6y_03rCiz4LO2LIVxanQXSE2aEF6TN1KhUzy8LiCaUz_b4XYwxbAdH0oBdVKj4J3CGHZBFhKc4Z2hKKf1nVC8B8dYpt_rEtcZRX6awfK7jiY_nKu5gimQihSV_OXXCIlQQvgE_C9DDCelK5s.az3XKdaCMTWfkj7GJkp4cHtdPZgJ6l.WCWUk5_RJwfiSrNJjoWX2wq"

// zfCxJstSpider 获取河北省住房和城乡建设厅标准
// @Title 获取河北省住房和城乡建设厅标准
// @Description https://zfcxjst.hebei.gov.cn/，获取河北省住房和城乡建设厅标准
func main() {
	// 国内标准列表
	var allCategory = []zfCxJstCategory{
		{url: "https://zfcxjst.hebei.gov.cn/hbzjt/ztzl/jj/gcjsgf/fwjz", name: "房屋建筑", page: 1, maxPage: 14},  //14
		{url: "https://zfcxjst.hebei.gov.cn/hbzjt/ztzl/jj/gcjsgf/szgc", name: "市政工程", page: 1, maxPage: 7},   //7
		{url: "https://zfcxjst.hebei.gov.cn/hbzjt/ztzl/jj/gcjsgf/czjs", name: "村镇建设", page: 1, maxPage: 3},   //3
		{url: "https://zfcxjst.hebei.gov.cn/hbzjt/ztzl/jj/gcjsgf/sxjs", name: "四新技术", page: 1, maxPage: 7},   //7
		{url: "https://zfcxjst.hebei.gov.cn/hbzjt/ztzl/jj/gcjsgf/zjbz", name: "消耗量标准", page: 1, maxPage: 2}, //2
	}
	for _, category := range allCategory {
		isPageListGo := true
		for isPageListGo {
			listUrl := fmt.Sprintf("%s/index_%d.html", category.url, category.page)
			fmt.Println(listUrl)
			listDoc, err := QueryZfCxJstHtml(listUrl, category.url)
			if err != nil {
				fmt.Println(err)
				isPageListGo = false
				break
			}
			liNodes := htmlquery.Find(listDoc, `//div[@class="content"]/div[@class="guTableBox"]/div[@class="gulistTable"]/div[@class="gfbzBox"]/div[@class="biaozhunList"]/ul[@id="number"]/li`)
			if len(liNodes) >= 1 {
				for _, liNode := range liNodes {
					fmt.Println("=====================开始处理数据=========================")
					fmt.Println(category.url, category.page, category.name)

					titleNode := htmlquery.FindOne(liNode, `./div[@class="mingcheng"]/a`)
					title := htmlquery.InnerText(titleNode)
					title = strings.TrimSpace(title)
					title = strings.ReplaceAll(title, "-", "")
					title = strings.ReplaceAll(title, " ", "")
					title = strings.ReplaceAll(title, "|", "-")
					title = strings.ReplaceAll(title, "/", "-")
					title = strings.ReplaceAll(title, "\n", "")
					title = strings.ReplaceAll(title, "\r", "")
					title = strings.ReplaceAll(title, " ", "")
					fmt.Println(title)

					codeNode := htmlquery.FindOne(liNode, `./div[@class="bianhao pso"]/span`)
					code := htmlquery.InnerText(codeNode)
					code = strings.TrimSpace(code)
					code = strings.ReplaceAll(code, "/", "-")
					fmt.Println(code)

					filePath := "../zfcxjst.hebei.gov.cn/" + title + "(" + code + ")" + ".pdf"
					fmt.Println(filePath)

					_, err = os.Stat(filePath)
					if err == nil {
						fmt.Println("文档已下载过，跳过")
						continue
					}

					detailUrl := htmlquery.InnerText(htmlquery.FindOne(liNode, `./div[@class="mingcheng"]/a/@href`))
					if strings.Index(detailUrl, "zfcxjst.hebei.gov.cn") == -1 {
						detailUrl = "https://zfcxjst.hebei.gov.cn" + detailUrl
					}
					fmt.Println(detailUrl)
					detailDoc, err := QueryZfCxJstHtml(detailUrl, listUrl)
					if err != nil {
						fmt.Println("无法获取文档详情，跳过")
						continue
					}
					downNode := htmlquery.FindOne(detailDoc, `//div[@class="pc"]/div[@class="zhuye_box"]/div[@class="p_nei"]/div[@class="panel-body"]/span[@class="info_affix_file"]/a/@href`)
					if downNode == nil {
						fmt.Println("没有下载地址，跳过")
						continue
					}

					downloadUrl := "https://zfcxjst.hebei.gov.cn" + htmlquery.InnerText(downNode)
					// 只下载pdf文件
					if strings.Index(downloadUrl, ".pdf") == -1 {
						fmt.Println("不是pdf文件")
						continue
					}
					fmt.Println(downloadUrl)

					fmt.Println("=======开始下载========")
					err = downloadZfCxJstPdf(downloadUrl, filePath, detailUrl)
					if err != nil {
						fmt.Println(err)
					}
					//复制文件
					tempFilePath := strings.ReplaceAll(filePath, "../zfcxjst.hebei.gov.cn", "../upload.doc88.com/zfcxjst.hebei.gov.cn")
					err = ZfCxJstCopyFile(filePath, tempFilePath)
					if err != nil {
						fmt.Println(err)
						continue
					}

					fmt.Println("=======下载完成========")
					downloadZfCxJstPdfSleep := rand.Intn(5)
					for i := 1; i <= downloadZfCxJstPdfSleep; i++ {
						time.Sleep(time.Second)
						fmt.Println("page="+strconv.Itoa(category.page)+"=======", title, "成功，category_name="+category.name+"====== 暂停", downloadZfCxJstPdfSleep, "秒，倒计时", i, "秒===========")
					}
				}
				// DownLoadZfCxJstPageTimeSleep := 10
				DownLoadZfCxJstPageTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadZfCxJstPageTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page="+strconv.Itoa(category.page)+"====category_name="+category.name+"====== 暂停", DownLoadZfCxJstPageTimeSleep, "秒 倒计时", i, "秒===========")
				}
				category.page++
				if category.page > category.maxPage {
					isPageListGo = false
					break
				}
			}
		}
	}
}

func QueryZfCxJstHtml(requestUrl string, referer string) (doc *html.Node, err error) {
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
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", ZfCxJstCookie)
	req.Header.Set("Host", "zfcxjst.hebei.gov.cn")
	req.Header.Set("Origin", "https://zfcxjst.hebei.gov.cn")
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

func downloadZfCxJstPdf(pdfUrl string, filePath string, referer string) error {
	// 初始化客户端
	var client http.Client
	req, err := http.NewRequest("GET", pdfUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", ZfCxJstCookie)
	req.Header.Set("Host", "zfcxjst.hebei.gov.cn")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", referer)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
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

func ZfCxJstCopyFile(src, dst string) (err error) {
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
