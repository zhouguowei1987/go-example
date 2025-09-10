package main

import (
	"crypto/tls"
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
	"math/rand"
)

var TbzCookie = "__jsluid_s=b7d35d18c6c44705ce234044421b8f67; Hm_lvt_8c446e9fafe752e4975210bc30d7ab9d=1752074026,1752916674,1753086531,1754059674; HMACCOUNT=1CCD0111717619C6; ASP.NET_SessionId=smepcq3yb0e4tp525qc2521i; Hm_lpvt_8c446e9fafe752e4975210bc30d7ab9d=1754298876"

// TbzSpider 获取全国团体标准信息平台Pdf文档
// @Title 获取全国团体标准信息平台Pdf文档
// @Description https://www.ttbz.org.cn/，将全国团体标准信息平台Pdf文档入库
func main() {
//     144758
	var startId = 145900
	var endId = 145989
	for id := startId; id <= endId; id++ {
		fmt.Println(id)
		pdfsUrl := fmt.Sprintf("https://www.ttbz.org.cn/Pdfs/Index/?ftype=st&pms=%d", id)
		pdfsDoc, err := TbzPdfsDoc(pdfsUrl)
		if err != nil {
			fmt.Println(err)
			continue
		}
		iframeSrcNode := htmlquery.FindOne(pdfsDoc, `//iframe[@id="myiframe"]/@src`)
		if iframeSrcNode == nil{
		    fmt.Println("页面不存在")
			continue
		}
		iframeSrc := htmlquery.InnerText(iframeSrcNode)
		fmt.Println(iframeSrc)

		// 下载pdf文件
		pdfUrl := strings.ReplaceAll(iframeSrc, "/Home/PdfView?file=", "https://www.ttbz.org.cn")
		// 移除rnd参数
		pdfUrl = strings.Split(pdfUrl, "&")[0]
		fmt.Println(pdfUrl)

		fmt.Println("=======开始下载========")
		filePath, err := downloadTbzPdf(pdfUrl, id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// 查看文件大小，如果是空文件，则删除
		fileInfo, err := os.Stat(filePath)
		if err == nil && fileInfo.Size() == 0 {
			fmt.Println("空文件删除")
			err = os.Remove(filePath)
		}
		if err != nil {
            continue
        }
		//复制文件
		tempFilePath := strings.ReplaceAll(filePath, "www.ttbz.org.cn", "temp-www.ttbz.org.cn")
		err = copyTbzFile(filePath, tempFilePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("=======完成下载========")

		// 设置倒计时
		// DownLoadTTbzTimeSleep := 10
		DownLoadTTbzTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadTTbzTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("id="+strconv.Itoa(id)+"===========操作完成，", "暂停", DownLoadTTbzTimeSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

func TbzPdfsDoc(url string) (doc *html.Node, err error) {
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
	req.Header.Set("Cookie", TbzCookie)
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	req.Header.Set("Host", "www.ttbz.org.cn")
	req.Header.Set("Pragma", "no-cache")
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

func downloadTbzPdf(pdfUrl string, pdfId int) (filePath string, err error) {
	// 创建一个自定义的http.Transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 忽略证书验证
		},
	}
	client := &http.Client{Transport: tr}           //初始化客户端                     //初始化客户端
	req, err := http.NewRequest("GET", pdfUrl, nil) //建立连接
	if err != nil {
		return filePath, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.ttbz.org.cn")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return filePath, err
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return filePath, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}

	contentDisposition := resp.Header.Get("Content-Disposition")
	contentDispositionUnescape, err := url.QueryUnescape(contentDisposition)
	if err != nil {
		return filePath, err
	}
	if len(contentDispositionUnescape) <= 0 {
		return filePath, err
	}
	contentDispositionUnescape = strings.ReplaceAll(contentDispositionUnescape, "inline;FileName=", "")
	contentDispositionUnescape = strings.TrimSpace(contentDispositionUnescape)
	contentDispositionUnescape = strings.Replace(contentDispositionUnescape, "_", "-", 1)
	contentDispositionUnescapeArray := strings.Split(contentDispositionUnescape, "_")
	code := contentDispositionUnescapeArray[0]
	code = strings.ReplaceAll(code, "/", "-")
	fmt.Println(code)

	titleArray := contentDispositionUnescapeArray[1:]
	title := strings.Join(titleArray,"-")
	title = strings.ReplaceAll(title, " ", "")
	title = strings.ReplaceAll(title, "　", "")
	title = strings.ReplaceAll(title, "/", "-")
	title = strings.ReplaceAll(title, "--", "-")
	title = strings.ReplaceAll(title, ".pdf", "")
	fmt.Println(title)
	filePath = "../www.ttbz.org.cn/" + strconv.Itoa(pdfId) + "-" + title + "(" + code + ").pdf"
	fmt.Println(filePath)

	_, err = os.Stat(filePath)
    if err == nil {
        return filePath, errors.New("文档已下载过，跳过")
    }

	// 创建一个文件用于保存
	fileDiv := filepath.Dir(filePath)
	if _, err = os.Stat(fileDiv); err != nil {
		if os.MkdirAll(fileDiv, 0644) != nil {
			return filePath, err
		}
	}
	out, err := os.Create(filePath)
	if err != nil {
		return filePath, err
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return filePath, err
	}
	return filePath, nil
}

func copyTbzFile(src, dst string) (err error) {
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
