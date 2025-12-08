package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var GsDfBzCookie = "CookieNid=http%3A%2F%2Fwww.gsdfbz.cn%2Ftheme%2Fdefault%2FstandardPublishDetail4830; Hm_lvt_9e72e5e781fbc2091396de8967bf1ce1=1765093641; HMACCOUNT=4E5B3419A3141A8E; JSESSIONID=8EF941D4B67C18F5025B6B83C7046E0A; Hm_lpvt_9e72e5e781fbc2091396de8967bf1ce1=1765176390"

// GsDfBzSpider 获取甘肃省地方标准文档
// @Title 获取甘肃省地方标准文档
// @Description http://www.gsdfbz.cn/，将甘肃省地方标准文档入库
func main() {
	var startId = 1
	var endId = 4830
	for id := startId; id <= endId; id++ {
		fmt.Println(id)
		detailUrl := fmt.Sprintf("http://www.gsdfbz.cn/theme/default/standardPublishDetail%d", id)
		detailDoc, err := GsDfBzDetailDoc(detailUrl)
		if err != nil {
			fmt.Println(err)
			continue
		}
		codeNode := htmlquery.FindOne(detailDoc, `//html/body/div[4]/div/div[1]/table/tbody/tr[1]/td[2]/span`)
		if codeNode == nil {
			fmt.Println("标准号不存在，跳过")
			continue
		}
		code := strings.TrimSpace(htmlquery.InnerText(codeNode))
		code = strings.ReplaceAll(code, "/", "-")
		code = strings.ReplaceAll(code, "—", "-")
		fmt.Println(code)

		titleNode := htmlquery.FindOne(detailDoc, `//html/body/div[4]/div/div[1]/table/tbody/tr[2]/td[2]`)
		if titleNode == nil {
			fmt.Println("标题不存在，跳过")
			continue
		}
		title := strings.TrimSpace(htmlquery.InnerText(titleNode))
		title = strings.ReplaceAll(title, " ", "-")
		title = strings.ReplaceAll(title, "　", "-")
		title = strings.ReplaceAll(title, "/", "-")
		title = strings.ReplaceAll(title, "--", "-")
		fmt.Println(title)

		filePath := "../www.gsdfbz.cn/" + title + "(" + code + ")" + ".pdf"
		fmt.Println(filePath)

		_, err = os.Stat(filePath)
		if err == nil {
			fmt.Println("文档已下载过，跳过")
			continue
		}

		// 下载pdf文件
		downloadUrl := fmt.Sprintf("http://www.gsdfbz.cn/getPdf/%d", id)
		fmt.Println(downloadUrl)

		fmt.Println("=======开始下载========")
		err = downloadGsDfBz(downloadUrl, detailUrl, filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//复制文件
		tempFilePath := strings.ReplaceAll(filePath, "www.gsdfbz.cn", "../upload.doc88.com/dbba.sacinfo.org.cn")
		err = copyGsDfBzFile(filePath, tempFilePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("=======完成下载========")

		// 设置倒计时
		// DownLoadTGsDfBzTimeSleep := 10
		DownLoadTGsDfBzTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadTGsDfBzTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("id="+strconv.Itoa(id)+"===========操作完成，", "暂停", DownLoadTGsDfBzTimeSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

func GsDfBzDetailDoc(url string) (doc *html.Node, err error) {
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
	req.Header.Set("Cookie", GsDfBzCookie)
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	req.Header.Set("Host", "www.gsdfbz.cn")
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
	doc, err = htmlquery.Parse(resp.Body)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func downloadGsDfBz(requestUrl string, referer string, filePath string) error {
	// 创建一个自定义的http.Transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 忽略证书验证
		},
	}
	client := &http.Client{Transport: tr}               //初始化客户端
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", GsDfBzCookie)
	req.Header.Set("Host", "www.gsdfbz.cn")
	req.Header.Set("Origin", "http://www.gsdfbz.cn")
	req.Header.Set("Referer", referer)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusCreated {
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

func copyGsDfBzFile(src, dst string) (err error) {
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
