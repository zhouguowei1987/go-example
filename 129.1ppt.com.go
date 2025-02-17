package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"math/rand"
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
	var startId = 3
	var endId = 130463
	for id := startId; id <= endId; id++ {
		err := pptSpider(id)
		if err != nil {
			fmt.Println(err)
		}
	}
	//pptSpider(130283)
}

func pptSpider(id int) error {
	detailUrl := fmt.Sprintf("https://www.1ppt.com/article/%d.html", id)
	fmt.Println(detailUrl)
	detailDoc, err := htmlquery.LoadURL(detailUrl)

	if err != nil {
		return err
	}
	// 查看是否有下载按钮
	detailDownloadButtonNode := htmlquery.FindOne(detailDoc, `//ul[@class="downurllist"]/li[1]`)
	if detailDownloadButtonNode == nil {
		return errors.New("详情页没有下载按钮")
	}

	downloadDetailUrl := fmt.Sprintf("https://www.1ppt.com/plus/download.php?open=0&aid=%d&cid=3", id)
	fmt.Println(downloadDetailUrl)
	downloadDetailDoc, err := getPptDownloadDetailDoc(downloadDetailUrl, detailUrl)
	if err != nil {
		return err
	}
	// 文档名称
	titleNode := htmlquery.FindOne(downloadDetailDoc, `//dl[@class="downloadpage"]/dt/h1/a`)
	if titleNode == nil {
		return errors.New("下载详情页没有附件标题")
	}
	title := htmlquery.InnerText(titleNode)
	fmt.Println(title)
	// 过滤文件名中含有“图”字样文件
	if strings.Index(title, "图") != -1 {
		return errors.New("过滤文件名中含有“图”字样文件")
	}
	// 过滤文件名中含有“张”字样文件
	if strings.Index(title, "张") != -1 {
		return errors.New("过滤文件名中含有“张”字样文件")
	}
	// 过滤文件名中含有“套”字样文件
	if strings.Index(title, "套") != -1 {
		return errors.New("过滤文件名中含有“套”字样文件")
	}
	// 过滤文件名中含有“个”字样文件
	if strings.Index(title, "个") != -1 {
		return errors.New("过滤文件名中含有“个”字样文件")
	}
	// 过滤文件名中含有“页”字样文件
	if strings.Index(title, "页") != -1 {
		return errors.New("过滤文件名中含有“页”字样文件")
	}
	// 过滤文件名中含有“年”字样文件
	if strings.Index(title, "年") != -1 {
		return errors.New("过滤文件名中含有“年”字样文件")
	}
	// 过滤文件名中含有“素材”字样文件
	if strings.Index(title, "素材") != -1 {
		return errors.New("过滤文件名中含有“素材”字样文件")
	}

	// 查看是否有下载按钮
	downloadButtonNode := htmlquery.FindOne(downloadDetailDoc, `//ul[@class="downloadlist"]/li[@class="c1"]/a`)
	if downloadButtonNode == nil {
		return errors.New("下载详情页没有下载按钮")
	}
	// 附件下载链接
	attachUrl := htmlquery.SelectAttr(downloadButtonNode, "href")
	fmt.Println(attachUrl)

	// 获取文件后缀
	downloadUrlSplitArray := strings.Split(attachUrl, ".")
	fileSuffix := downloadUrlSplitArray[len(downloadUrlSplitArray)-1]
	fileSuffixArray := []string{"zip", "rar"}
	if !stringContains(fileSuffixArray, fileSuffix) {
		return errors.New("既不是zip文件，也不是rar文件，跳过")
	}
	filePath := "F:\\workspace\\www.1ppt.com\\www." + fileSuffix + "_1ppt.com/" + title + "." + fileSuffix
	if _, err := os.Stat(filePath); err != nil {
		fmt.Println("=======开始下载========")
		err = downloadPpt(attachUrl, downloadDetailUrl, filePath)
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

func stringContains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func getPptDownloadDetailDoc(url string, referer string) (doc *html.Node, err error) {
	client := &http.Client{}                     //初始化客户端
	req, err := http.NewRequest("GET", url, nil) //建立连接
	if err != nil {
		return doc, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "mizToken=202501200928550.359382022694436640.9092217902997255; HMACCOUNT=00EDEFEA78E0441D; acw_tc=1a0c655917391955004167396e0045853fbe2d10f2f74d83faf9ee6a1f2601; Hm_lvt_087ceb5ea69d10fb5bbb6bc49c209fa2=1737336287,1739195501; Hm_lpvt_087ceb5ea69d10fb5bbb6bc49c209fa2=1739196761")
	req.Header.Set("Host", "www.1ppt.com")
	req.Header.Set("Referer", referer)
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
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
	doc, err = decodeAndParseHTML(string(bodyBytes))
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func decodeAndParseHTML(gb2312Content string) (*html.Node, error) {
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

func downloadPpt(pdfUrl string, referer string, filePath string) error {
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
