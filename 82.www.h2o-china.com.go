package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	H2oChinaEnableHttpProxy = false
	H2oChinaHttpProxyUrl    = "27.42.168.46:55481"
)

func H2oChinaSetHttpProxy() (httpclient http.Client) {
	ProxyURL, _ := url.Parse(H2oChinaHttpProxyUrl)
	httpclient = http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

// h2oChinaSpider 获取中国水网Pdf文档
// @Title 获取中国水网Pdf文档
// @Description https://www.h2o-china.com/，获取中国水网Pdf文档
func main() {
	page := 1
	isPageGo := true
	for isPageGo {
		listUrl := fmt.Sprintf("https://www.h2o-china.com/standard/home?ordby=dateline&sort=DESC&page=%d", page)
		fmt.Println(listUrl)
		listDoc, _ := htmlquery.LoadURL(listUrl)
		liNodes := htmlquery.Find(listDoc, `//div[@class="lists txtList"]/ul/li/em[@class="title"]/a[@class="ellip w540 i-pdf"]`)
		if len(liNodes) >= 1 {
			for _, liNode := range liNodes {
				detailUrl := htmlquery.InnerText(htmlquery.FindOne(liNode, `./@href`))
				detailUrl = "https://www.h2o-china.com" + detailUrl
				fmt.Println(detailUrl)
				detailDoc, _ := htmlquery.LoadURL(detailUrl)

				title := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="hd"]/h1`))
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, " ", "")
				fmt.Println(title)

				standardNo := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="traits"]/table/tbody/tr[3]/td[2]`))
				standardNo = strings.ReplaceAll(standardNo, "/", "-")
				standardNo = strings.ReplaceAll(standardNo, ":", "-")
				standardNo = strings.ReplaceAll(standardNo, " ", "")
				//fmt.Println(standardNo)

				downloadUrl := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="dowloads fr"]/a/@href`))
				downloadUrl = "https://www.h2o-china.com" + downloadUrl
				fmt.Println(downloadUrl)

				downloadUrlArray, err := url.Parse(downloadUrl)
				filePath := "../www.h2o-china.com/" + downloadUrlArray.Query().Get("id") + "-" + title + ".pdf"
				if len(standardNo) > 1 {
					filePath = "../www.h2o-china.com/" + downloadUrlArray.Query().Get("id") + "-" + title + "(" + standardNo + ")" + ".pdf"
				}
				//fmt.Println(filePath)

				err = downloadH2oChinaPdf(downloadUrl, filePath)
				if err != nil {
					fmt.Println(err)
				}
			}
			page++
		} else {
			isPageGo = false
			page = 1
			break
		}
	}
}

func downloadH2oChinaPdf(pdfUrl string, filePath string) error {
	// 初始化客户端
	var client http.Client
	if H2oChinaEnableHttpProxy {
		client = H2oChinaSetHttpProxy()
	}
	req, err := http.NewRequest("GET", pdfUrl, nil) //建立连接
	if err != nil {
		return err
	}
	//req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	//req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	//req.Header.Set("Cache-Control", "no-cache")
	//req.Header.Set("Connection", "keep-alive")
	//req.Header.Set("Cookie", "_gid=GA1.2.834335799.1672890802; Hm_lvt_e00731ccbc1d6c46ec1b5ace98c1390e=1672890809; backurl=%2Fstandard%2Fhome%3Fordby%3Ddateline%26sort%3DDESC%26page%3D%7Bpage%7D%3D65; Hm_lpvt_e00731ccbc1d6c46ec1b5ace98c1390e=1672897834; _gat=1; _ga_34B604LFFQ=GS1.1.1672895011.2.1.1672897834.60.0.0; _ga=GA1.1.2115154058.1672890802")
	//req.Header.Set("Host", "www.h2o-china.com")
	//req.Header.Set("Pragma", "no-cache")
	//req.Header.Set("Referer", "https://www.h2o-china.com/")
	//req.Header.Set("Upgrade-Insecure-Requests", "1")
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
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
