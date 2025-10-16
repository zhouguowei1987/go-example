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
)

// PptSpider 获取优品ppt文档
// @Title 获取优品ppt文档
// @Description https://www.ypppt.com/，将优品ppt文档入库
func main() {
	var startId = 1
	var endId = 16546
	for id := startId; id <= endId; id++ {
		err := ypPptSpider(id)
		if err != nil {
			fmt.Println(err)
		}
	}
}

var ypPptCookie = "Hm_lvt_45db753385e6d769706e10062e3d6453=1739755517,1740966035; HMACCOUNT=2CEC63D57647BCA5; __gads=ID=6693e91c7506dd79:T=1737378621:RT=1740968676:S=ALNI_MZlmhin25aPeQT8frZdlEncQkMXAg; __gpi=UID=00000ff35a510a9d:T=1737378621:RT=1740968676:S=ALNI_MZLtORzWkv1xh0aeqzM40asLeKm_g; __eoi=ID=28e7f4789c36709e:T=1737378621:RT=1740968676:S=AA-Afjb-lo12lYgccK7ZNYa-bEJn; Hm_lpvt_45db753385e6d769706e10062e3d6453=1740968762"

func ypPptSpider(id int) error {
	detailUrl := fmt.Sprintf("https://www.ypppt.com/p/d.php?aid=%d", id)
	fmt.Println(detailUrl)
	detailDoc, err := htmlquery.LoadURL(detailUrl)

	if err != nil {
		return err
	}
	// 文档名称
	titleNode := htmlquery.FindOne(detailDoc, `//div[@class="de"]/h1`)
	if titleNode == nil {
		return errors.New("下载详情页没有附件标题")
	}
	title := htmlquery.InnerText(titleNode)
	title = strings.ReplaceAll(title, " - 下载页", "")
	title = strings.TrimSpace(title)
	fmt.Println(title)
	// 过滤文件名中不含有“ppt”字样文件
	if strings.Index(strings.ToLower(title), "ppt") == -1 {
		return errors.New("过滤文件名中不含有“ppt”字样文件")
	}
	// 过滤文件名中含有“课时”字样文件
	if strings.Index(title, "课时") != -1 {
		return errors.New("过滤文件名中含有“课时”字样文件")
	}
	// 过滤文件名中含有“课件”字样文件
	if strings.Index(title, "课件") != -1 {
		return errors.New("过滤文件名中含有“课件”字样文件")
	}
	// 过滤文件名中含有“体”字样文件
	if strings.Index(title, "体") != -1 {
		return errors.New("过滤文件名中含有“体”字样文件")
	}
	// 过滤文件名中含有“图”字样文件
	if strings.Index(title, "图") != -1 {
		return errors.New("过滤文件名中含有“图”字样文件")
	}
	// 过滤文件名中含有“素材”字样文件
	if strings.Index(title, "素材") != -1 {
		return errors.New("过滤文件名中含有“素材”字样文件")
	}
	// 查看是否有下载按钮
	detailDownloadButtonNode := htmlquery.FindOne(detailDoc, `//ul[@class="down clear"]/li[1]/a`)
	if detailDownloadButtonNode == nil {
		return errors.New("详情页没有下载按钮")
	}
	// 附件下载链接
	attachUrl := htmlquery.SelectAttr(detailDownloadButtonNode, "href")
	fmt.Println(attachUrl)

	// 获取文件后缀
	downloadUrlSplitArray := strings.Split(attachUrl, ".")
	fileSuffix := downloadUrlSplitArray[len(downloadUrlSplitArray)-1]
	fileSuffixArray := []string{"pptx"}
	if !ypPptStringContains(fileSuffixArray, fileSuffix) {
		return errors.New("既不是pptx文件，跳过")
	}
	filePath := "F:\\workspace\\www.ypppt.com\\www." + fileSuffix + "_www.ypppt.com\\" + title + "." + fileSuffix
	if _, err := os.Stat(filePath); err != nil {
		fmt.Println("=======开始下载========")
		err = downloadYpPpt(attachUrl, detailUrl, filePath)
		if err != nil {
			return err
		}
		fmt.Println("=======完成下载========")
		DownLoad1PptTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoad1PptTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("id="+strconv.Itoa(id)+"===========下载", title, "成功，暂停", DownLoad1PptTimeSleep, "秒，倒计时", i, "秒===========")
		}
	}
	return nil
}
func ypPptStringContains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
func downloadYpPpt(pdfUrl string, referer string, filePath string) error {
	// 创建一个自定义的http.Transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 忽略证书验证
		},
	}
	client := &http.Client{Transport: tr}           //初始化客户端
	req, err := http.NewRequest("GET", pdfUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", ypPptCookie)
	req.Header.Set("Host", "www.ypppt.com")
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
