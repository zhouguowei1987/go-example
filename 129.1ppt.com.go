package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// PptSpider 获取第一ppt文档
// @Title 获取第一ppt文档
// @Description https://1ppt.com/，将第一ppt文档入库
func main() {
	//103045
	var startId = 129682
	var endId = 129683
	goCh := make(chan int, endId-startId)
	for id := startId; id <= endId; id++ {
		// 设置下载倒计时
		//DownLoadPptTimeSleep := rand.Intn(20)
		DownLoadPptTimeSleep := 10
		for i := 1; i <= DownLoadPptTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("id="+strconv.Itoa(id)+"===========操作完成，", "暂停", DownLoadPptTimeSleep, "秒，倒计时", i, "秒===========")
		}
		go func(id int) {
			err := pptSpider(id)
			if err != nil {
				fmt.Println(err)
			}
			goCh <- id
		}(id)
		fmt.Println(<-goCh)
	}
	//pptSpider(82515)
}

func pptSpider(id int) error {
	detailUrl := fmt.Sprintf("https://www.1ppt.com/article/%d.html", id)
	detailDoc, err := htmlquery.LoadURL(detailUrl)
	if err != nil {
		return err
	}
	// 查看是否有下载按钮
	detailDownloadButtonNode := htmlquery.FindOne(detailDoc, `//ul[@class="downurllist"]/li[1]`)
	if detailDownloadButtonNode == nil {
		return errors.New("详情页没有下载按钮")
	}
	downloadUrl := fmt.Sprintf("https://www.1ppt.com/plus/download.php?open=0&aid=%d&cid=3", id)
	downloadDetailDoc, err := getPptDownloadDoc(downloadUrl, detailUrl)
	if err != nil {
		return err
	}
	// 查看是否有下载按钮
	downloadButtonNode := htmlquery.FindOne(downloadDetailDoc, `//ul[@class="downloadlist"]/li[@class"c1"]/a`)
	if downloadButtonNode == nil {
		return errors.New("下载详情页没有下载按钮")
	}

	// 文档名称
	titleNode := htmlquery.FindOne(downloadDetailDoc, `//dl[@class="downloadpage"]/dt/h1//a`)
	if titleNode == nil {
		return errors.New("下载详情页没有附件标题")
	}
	title := htmlquery.InnerText(titleNode)
	fmt.Println(title)

	// 附件下载链接
	attachUrl := htmlquery.SelectAttr(downloadButtonNode, "href")
	fmt.Println(attachUrl)

	// 获取文件后缀
	downloadUrlSplitArray := strings.Split(attachUrl, ".")
	fileSuffix := downloadUrlSplitArray[len(downloadUrlSplitArray)-1]
	if fileSuffix != "zip" {
		return errors.New("不是zip文件，跳过")
	}

	filePath := "../1ppt.com/" + title + ".zip"
	if _, err := os.Stat(filePath); err != nil {
		fmt.Println("=======开始下载========")
		err = downloadPpt(attachUrl, filePath)
		if err != nil {
			return err
		}
		fmt.Println("=======完成下载========")
	}
	return nil
}

func getPptDownloadDoc(url string, reffer string) (doc *html.Node, err error) {
	client := &http.Client{}                     //初始化客户端
	req, err := http.NewRequest("GET", url, nil) //建立连接
	if err != nil {
		return doc, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "mizToken=202501191741290.5355677329169450.001375657287244314; Hm_lvt_087ceb5ea69d10fb5bbb6bc49c209fa2=1737279673; HMACCOUNT=2CEC63D57647BCA5; acw_tc=2760826617373331730588303edf6930cea2e1031f6c2dc10a9d79b45ba631; Hm_lpvt_087ceb5ea69d10fb5bbb6bc49c209fa2=1737333516")
	req.Header.Set("Host", "www.1ppt.com")
	req.Header.Set("Reffer", reffer)
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

func downloadPpt(pdfUrl string, filePath string) error {
	client := &http.Client{}                        //初始化客户端
	req, err := http.NewRequest("GET", pdfUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "mizToken=202501191741290.5355677329169450.001375657287244314; Hm_lvt_087ceb5ea69d10fb5bbb6bc49c209fa2=1737279673; HMACCOUNT=2CEC63D57647BCA5; acw_tc=2760826617373331730588303edf6930cea2e1031f6c2dc10a9d79b45ba631; Hm_lpvt_087ceb5ea69d10fb5bbb6bc49c209fa2=1737333516")
	req.Header.Set("Host", "1ppt.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
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
