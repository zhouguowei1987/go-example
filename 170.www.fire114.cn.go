package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var Fire114Cookie = "PHPSESSID=sgeiq3i8ut557f35o4mnur1p83; UM_distinctid=199997a7ccd1070-0a26c6e829008b-26001d51-1fa400-199997a7cce577; sessionid=sgeiq3i8ut557f35o4mnur1p83; CNZZDATA1978074=cnzz_eid%3D843923194-1759216565-https%253A%252F%252Fnew.fire114.cn%252F%26ntime%3D1759992324"

// Fire114Spider 获取消防百事通文档
// @Title 获取消防百事通文档
// @Description http://www.fire114.cn/，将消防百事通文档入库
func main() {
	// 139565
	var startId = 140337
	// 19966
	var endId = 140000
	for id := startId; id >= endId; id-- {
		detailUrl := fmt.Sprintf("https://www.fire114.cn/islibd/%d.html", id)
		fmt.Println(detailUrl)
		refererUrl := fmt.Sprintf("https://www.fire114.cn/islibd/%d.html", id-1)
		fmt.Println(refererUrl)
		detailDoc, err := Fire114Doc(detailUrl, refererUrl)
		if err != nil {
			fmt.Println(err)
			continue
		}
		titleNode := htmlquery.FindOne(detailDoc, `//div[@class="mainWrap"]/div[@class="containWrap row"]/div[@class="col-9"]/div[@class="titleWrap"]/div/span[@class="title"]`)
		if titleNode == nil {
			fmt.Println("没有标题节点")
			continue
		}
		title := htmlquery.InnerText(titleNode)

		title = strings.ReplaceAll(title, "《", "")
		title = strings.ReplaceAll(title, "》", "")
		title = strings.ReplaceAll(title, "【", "")
		title = strings.ReplaceAll(title, "】", "")
		title = strings.ReplaceAll(title, "()", "")
		title = strings.ReplaceAll(title, "（)", "")
		title = strings.ReplaceAll(title, "‘", "")
		title = strings.TrimSpace(title)
		if strings.Index(title, ".pdf") == -1 && strings.Index(title, ".doc") == -1 {
			fmt.Println("不是pdf、doc文档")
			continue
		}
		if strings.Index(title, "检验报告") != -1 {
			fmt.Println("标题含有‘检验报告’字样，跳过")
			continue
		}
		fmt.Println(title)

		filePath := "../www.fire114.cn/" + title
		fmt.Println(filePath)
		_, err = os.Stat(filePath)
		if err == nil {
			fmt.Println("文档已下载过，跳过")
			continue
		}
		fmt.Println("=======开始下载========")
		iframeSrcNode := htmlquery.FindOne(detailDoc, `//div[@class="mainWrap"]/div[@class="containWrap row"]/div[@class="col-9"]/iframe/@src`)
		if iframeSrcNode == nil {
			fmt.Println("没有下载链接节点")
			continue
		}
		downLoadUrl := htmlquery.InnerText(iframeSrcNode)
		if strings.Index(downLoadUrl, "https://new.fire114.cn") != -1 {
			downLoadUrl = strings.ReplaceAll(downLoadUrl, "/Cms/widget/pdfjs/web/viewer-2.html?cwgpdfsrcurl=", "")
		} else if strings.Index(downLoadUrl, "https://oss.fire114.cn") != -1 {
			downLoadUrl = strings.ReplaceAll(downLoadUrl, "https://oss.fire114.cn", "https://new.fire114.cn/uploads")
		} else if strings.Index(downLoadUrl, "https://view.officeapps.live.com") != -1 {
			downLoadUrl = strings.ReplaceAll(downLoadUrl, "https://view.officeapps.live.com/op/view.aspx?src=", "https://")
		}
		fmt.Println(downLoadUrl)
		err = downFire114(downLoadUrl, filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("=======完成下载========")

		//复制文件
		tempFilePath := strings.ReplaceAll(filePath, "../www.fire114.cn", "../temp-www.fire114.cn")
		err = copyFire114File(filePath, tempFilePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// 设置倒计时
		// DownLoadFire114TimeSleep := 10
		DownLoadFire114TimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadFire114TimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("id="+strconv.Itoa(id)+"===========操作完成，", "暂停", DownLoadFire114TimeSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

func Fire114Doc(requestUrl string, refererUrl string) (doc *html.Node, err error) {
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
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", Fire114Cookie)
	req.Header.Set("Host", "www.fire114.cn")
	req.Header.Set("Referer", refererUrl)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
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

func downFire114(pdfUrl string, filePath string) error {
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
	req, err := http.NewRequest("GET", pdfUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", Fire114Cookie)
	req.Header.Set("Host", "www.fire114.cn")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
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
func copyFire114File(src, dst string) (err error) {
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
