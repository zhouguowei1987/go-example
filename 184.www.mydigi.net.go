package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var MyDiGiCookie = "ASPSESSIONIDSQSRBSRQ=MHEDBGHBGANCNJKMHDNPAHDF; __utma=136034221.338647901.1778124082.1778124082.1778124082.1; __utmb=136034221; __utmc=136034221; __utmz=136034221.1778124082.1.1.utmccn=(direct)|utmcsr=(direct)|utmcmd=(none); newasp%5Fnet="

// MyDiGiSpider 获取说明书之家文档
// @Title 获取说明书之家文档
// @Description http://www.mydigi.net/，将说明书之家文档入库
func main() {
	var startId = 8338
	var endId = 40628
	for id := startId; id <= endId; id++ {
		showUrl := fmt.Sprintf("http://www.mydigi.net/soft/show.asp?id=%d", id)
		fmt.Println(showUrl)
		showDoc, err := MyDiGiHtmlDoc(showUrl)
		if err != nil {
			fmt.Println(err)
			continue
		}

		titleNode := htmlquery.FindOne(showDoc, `//html/body/table[7]/tbody/tr[1]/td/table/tbody/tr/td[2]`)
		if titleNode == nil {
			fmt.Println("标题不存在，跳过")
			continue
		}
		title := strings.TrimSpace(htmlquery.InnerText(titleNode))
		title = strings.TrimSpace(title)
		title = strings.ReplaceAll(title, "您当前的位置：说明书之家 -> ", "")
		titleArray := strings.Split(title, " -> ")
		title = titleArray[1] + "(" + titleArray[0] + ")" + "-" + titleArray[2]
		title = strings.ReplaceAll(title, "/", "-")
		title = strings.ReplaceAll(title, "／", "-")
		title = strings.ReplaceAll(title, "/", "-")
		title = strings.ReplaceAll(title, "　", "-")
		title = strings.ReplaceAll(title, " ", "-")
		title = strings.ReplaceAll(title, "：", ":")
		title = strings.ReplaceAll(title, "—", "-")
		title = strings.ReplaceAll(title, "－", "-")
		title = strings.ReplaceAll(title, "（", "(")
		title = strings.ReplaceAll(title, "）", ")")
		title = strings.ReplaceAll(title, "《", "")
		title = strings.ReplaceAll(title, "》", "")
		fmt.Println(title)

		filePath := "../www.mydigi.net/" + title + ".pdf"
		fmt.Println(filePath)

		_, err = os.Stat(filePath)
		if err == nil {
			fmt.Println("文档已下载过，跳过")
			continue
		}

		manualDownUrl := fmt.Sprintf("http://www.mydigi.net/soft/manualdown.asp?softid=%d", id)
		fmt.Println(showUrl)
		manualDownDoc, err := MyDiGiHtmlDoc(manualDownUrl)
		if err != nil {
			fmt.Println(err)
			continue
		}
		downloadViewHrefNode := htmlquery.FindOne(manualDownDoc, `//html/body/table[8]/tbody/tr/td[3]/table/tbody/tr[9]/td/div/a[1]/@href`)
		if downloadViewHrefNode == nil {
			fmt.Println("获取下载链接页面不存在，跳过")
			continue
		}
		downloadViewHref := strings.TrimSpace(htmlquery.InnerText(downloadViewHrefNode))
		downloadViewHref = strings.TrimSpace(downloadViewHref)
		downloadViewUrl := "http://www.mydigi.net/soft/" + downloadViewHref
		fmt.Println(downloadViewUrl)
		downloadViewDoc, err := MyDiGiHtmlDoc(downloadViewUrl)
		if err != nil {
			fmt.Println(err)
			continue
		}
		downloadUrlNode := htmlquery.FindOne(downloadViewDoc, `//html/body/div/table[3]/tbody/tr[3]/td/a/@href`)
		if downloadUrlNode == nil {
			fmt.Println("文档下载链接不存在，跳过")
			continue
		}
		downloadUrl := strings.TrimSpace(htmlquery.InnerText(downloadUrlNode))
		fmt.Println(downloadUrl)

		fmt.Println("=======开始下载========")
		err = downloadMyDiGi(downloadUrl, downloadViewUrl, filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//复制文件
		tempFilePath := strings.ReplaceAll(filePath, "www.mydigi.net", "upload.doc88.com/www.mydigi.net")
		err = copyMyDiGiFile(filePath, tempFilePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("=======完成下载========")

		// 设置倒计时
		// DownLoadTMyDiGiTimeSleep := 10
		DownLoadTMyDiGiTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadTMyDiGiTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("id="+strconv.Itoa(id)+"===========操作完成，", "暂停", DownLoadTMyDiGiTimeSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

func MyDiGiHtmlDoc(url string) (doc *html.Node, err error) {
	// 创建一个自定义的http.Transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 忽略证书验证
		},
	}
	client := &http.Client{Transport: tr}        //初始化客户端
	req, err := http.NewRequest("GET", url, nil) //建立连接
	if err != nil {
		return doc, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", MyDiGiCookie)
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	req.Header.Set("Host", "www.mydigi.net")
	req.Header.Set("Referer", "http://www.mydigi.net/sat/standard/standardlist/0")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return doc, err
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return doc, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return doc, err
	}
	doc, err = decodeAndParseHTMLMyDiGi(string(bodyBytes))
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func decodeAndParseHTMLMyDiGi(gb2312Content string) (*html.Node, error) {
	// 使用GB2312解码器解码内容
	decoder := simplifiedchinese.GBK.NewDecoder() // 注意：通常GB2312在Go中对应的是GBK，而非直接使用GB2312，因为GB2312不是一个广泛支持的编码标准，而是GBK的一个子集。
	decodedContent, _, err := transform.Bytes(decoder, []byte(gb2312Content))
	if err != nil {
		return nil, err
	}
	// 将解码后的内容转换为UTF-8（通常HTML解析器需要UTF-8编码）
	utf8Content := decodedContent
	// 解析HTML
	doc, err := html.Parse(bytes.NewReader(utf8Content))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func downloadMyDiGi(requestUrl string, referer string, filePath string) error {
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
	} //初始化客户端
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", MyDiGiCookie)
	req.Header.Set("Host", "www.mydigi.net")
	req.Header.Set("Origin", "http://www.mydigi.net")
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

func copyMyDiGiFile(src, dst string) (err error) {
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
