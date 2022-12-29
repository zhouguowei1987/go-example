package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	EnableHttpProxy = true
	HttpProxyUrl    = "27.42.168.46:55481"
)

func SetHttpProxy() (httpclient http.Client) {
	ProxyURL, _ := url.Parse(HttpProxyUrl)
	httpclient = http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

// bzkoSpider 获取标准库Pdf文档
// @Title 获取标准库Pdf文档
// @Description http://www.bzko.com/，获取标准库Pdf文档
func main() {
	var startId = 244952
	var endId = 244000
	var id = startId
	var isGoGo = true
	for isGoGo {
		fmt.Println(id)
		err := bzkoSpider(id)
		time.Sleep(time.Second * 20)
		id--
		if err != nil {
			fmt.Println(err)
			continue
		}
		if id <= endId {
			isGoGo = false
		}
	}
	//err := bzkoSpider(245019)
	//if err != nil {
	//	fmt.Println(err)
	//}
}

var cookies []*http.Cookie

func getBzko(showUrl string) (doc *html.Node, err error) {
	// 初始化客户端
	var client http.Client
	if EnableHttpProxy {
		client = SetHttpProxy()
	}
	req, err := http.NewRequest("GET", showUrl, nil) //建立连接
	if err != nil {
		return doc, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	//req.Header.Set("Cookie", "safedog-flow-item=D32C619CE1ED927EDC9E7F92A494F30A; _gid=GA1.2.252531875.1672196221; Hm_lvt_5bdc7cbf518dbe655b895f6db2b801e0=1672196226; __gads=ID=6dff11076c48e535-22132a920ed90015:T=1672196226:RT=1672196226:S=ALNI_MYnPXCTIXkHZ9f7SV3G8WrHFGvw_w; __gpi=UID=00000b98c454b02d:T=1672196226:RT=1672196226:S=ALNI_Maj4oEV-10qhFIisg1xzL1boh6siQ; IISSafeDogLGSession=E703A4230653F3462BFB2D61E2698729; _gat=1; Hm_lpvt_5bdc7cbf518dbe655b895f6db2b801e0=1672216518; _ga_34B604LFFQ=GS1.1.1672215062.3.1.1672216543.7.0.0; _ga=GA1.1.1520119252.1672196221")
	req.Header.Set("Host", "www.bzko.com")
	req.Header.Set("Pragma", "no-cache")
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
	cookies = resp.Cookies()
	doc, err = htmlquery.Parse(resp.Body)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func downloadRar(pdfUrl string, filePath string) error {
	// 初始化客户端
	var client http.Client
	if EnableHttpProxy {
		client = SetHttpProxy()
	}
	req, err := http.NewRequest("GET", pdfUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	//req.Header.Set("Cookie", "safedog-flow-item=D32C619CE1ED927EDC9E7F92A494F30A; _gid=GA1.2.252531875.1672196221; Hm_lvt_5bdc7cbf518dbe655b895f6db2b801e0=1672196226; __gads=ID=6dff11076c48e535-22132a920ed90015:T=1672196226:RT=1672196226:S=ALNI_MYnPXCTIXkHZ9f7SV3G8WrHFGvw_w; __gpi=UID=00000b98c454b02d:T=1672196226:RT=1672196226:S=ALNI_Maj4oEV-10qhFIisg1xzL1boh6siQ; IISSafeDogLGSession=E703A4230653F3462BFB2D61E2698729; _gat=1; Hm_lpvt_5bdc7cbf518dbe655b895f6db2b801e0=1672216518; _ga_34B604LFFQ=GS1.1.1672215062.3.1.1672216543.7.0.0; _ga=GA1.1.1520119252.1672196221")
	req.Header.Set("Host", "www.down.bzko.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "http://www.bzko.com/")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
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
		if os.MkdirAll(fileDiv, 0644) != nil {
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

func bzkoSpider(id int) error {
	showDownloadUrl := fmt.Sprintf("http://www.bzko.com/Common/ShowDownloadUrl.aspx?urlid=0&id=%d", id)
	showDownloadDoc, err := getBzko(showDownloadUrl)
	if err != nil {
		return err
	}
	// 文档标题
	titleNode := htmlquery.FindOne(showDownloadDoc, `//h1[@class="STYLE1"]`)
	titleText := htmlquery.InnerText(titleNode)
	titleText = strings.ReplaceAll(titleText, "/", "-")
	titleText = strings.ReplaceAll(titleText, " ", "")
	fmt.Println(titleText)

	// 文档rarUrl
	rarUrlNode := htmlquery.FindOne(showDownloadDoc, `//*[@id="content"]/table/tbody/tr/td/a/@href`)
	rarUrl := htmlquery.InnerText(rarUrlNode)
	fmt.Println(rarUrl)

	filePath := "../www.bzko.com/" + strconv.Itoa(id) + "-" + titleText + ".rar"
	err = downloadRar(rarUrl, filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
