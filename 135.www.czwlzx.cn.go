package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var CzWlZxEnableHttpProxy = false
var CzWlZxHttpProxyUrl = ""
var CzWlZxHttpProxyUrlArr = make([]string, 0)

func CzWlZxHttpProxy() error {
	pageMax := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// pageMax := []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for _, page := range pageMax {
		freeProxyUrl := "https://www.beesproxy.com/free"
		if page > 1 {
			freeProxyUrl = fmt.Sprintf("https://www.beesproxy.com/free/page/%d", page)
		}
		beesProxyDoc, err := htmlquery.LoadURL(freeProxyUrl)
		if err != nil {
			return err
		}
		trNodes := htmlquery.Find(beesProxyDoc, `//figure[@class="wp-block-table"]/table[@class="table table-bordered bg--secondary"]/tbody/tr`)
		if len(trNodes) > 0 {
			for _, trNode := range trNodes {
				ipNode := htmlquery.FindOne(trNode, "./td[1]")
				if ipNode == nil {
					continue
				}
				ip := htmlquery.InnerText(ipNode)

				portNode := htmlquery.FindOne(trNode, "./td[2]")
				if portNode == nil {
					continue
				}
				port := htmlquery.InnerText(portNode)

				protocolNode := htmlquery.FindOne(trNode, "./td[5]")
				if protocolNode == nil {
					continue
				}
				protocol := htmlquery.InnerText(protocolNode)

				switch protocol {
				case "HTTP":
					CzWlZxHttpProxyUrlArr = append(CzWlZxHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					CzWlZxHttpProxyUrlArr = append(CzWlZxHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func CzWlZxSetHttpProxy() (httpclient *http.Client) {
	if CzWlZxHttpProxyUrl == "" {
		if len(CzWlZxHttpProxyUrlArr) <= 0 {
			err := CzWlZxHttpProxy()
			if err != nil {
				CzWlZxSetHttpProxy()
			}
		}
		if len(CzWlZxHttpProxyUrlArr) > 1 {
			CzWlZxHttpProxyUrl = CzWlZxHttpProxyUrlArr[0]
		}
		if len(CzWlZxHttpProxyUrlArr) >= 2 {
			CzWlZxHttpProxyUrlArr = CzWlZxHttpProxyUrlArr[1:]
		} else {
			CzWlZxHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(CzWlZxHttpProxyUrl)
	ProxyURL, _ := url.Parse(CzWlZxHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
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
	return httpclient
}

var CzWlZxCookie = "Hm_lvt_43bc53ae85afc8f10b75f500b7f506b6=1751030895; ASP.NET_SessionId=qukkqa55xvbqhy45qt5w4gaw; Hm_lvt_43bc53ae85afc8f10b75f500b7f506b6=1750862924,1750943997,1751263265; Hm_lpvt_43bc53ae85afc8f10b75f500b7f506b6=1751263265; HMACCOUNT=1CCD0111717619C6"

// CzWlZxSpider 获取初中物理在线文档
// @Title 获取初中物理在线文档
// @Description http://www.czwlzx.cn/，获取初中物理在线文档
func main() {
	// 154563
	var startId = 22162
	var endId = 8254
	for id := startId; id >= endId; id-- {
		err := CzWlZxSpider(id)
		if err != nil {
			fmt.Println(err)
		}
	}
	// err := CzWlZxSpider(154187)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func CzWlZxSpider(id int) error {
	detailDocUrl := fmt.Sprintf("http://www.czwlzx.cn/Item/%d.aspx", id)
	fmt.Println(detailDocUrl)
	var detailRefererUrl = "http://www.czwlzx.cn/sj/List_1219.html"
	detailDoc, err := CzWlZxDetailDoc(detailDocUrl, detailRefererUrl)
	// fmt.Println(htmlquery.InnerText(detailDoc))
	// os.Exit(1)

	if err != nil {
		return err
	}
	fmt.Println("============================================================================")
	fmt.Println("=======当前页URL", detailDocUrl, "========")

	titleNode := htmlquery.FindOne(detailDoc, `//div[@class="content"]/div[@class="ppttit"]/h2`)
	if titleNode == nil {
		return errors.New("标题不存在")
	}
	title := htmlquery.InnerText(titleNode)
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, "免费", "")
	title = strings.ReplaceAll(title, "-", "")
	title = strings.ReplaceAll(title, " ", "")
	title = strings.ReplaceAll(title, "|", "-")
	fmt.Println(title)
	// os.Exit(1)
	// 过滤文件名中含有“扫描”字样文件
	if strings.Index(title, "扫描") != -1 {
		return errors.New("过滤文件名中含有“扫描”字样文件")
	}
	// 过滤文件名中含有“图片”字样文件
	if strings.Index(title, "图片") != -1 {
		return errors.New("过滤文件名中含有“图片”字样文件")
	}
	// 过滤文件名中含有“pdf”字样文件
	if strings.Index(title, "pdf") != -1 {
		return errors.New("过滤文件名中含有“pdf”字样文件")
	}

	// 文档类型
	var fileType = ""
	fileTypeDocxNode := htmlquery.FindOne(detailDoc, `//div[@class="content"]/div[@class="ppttit"]/h2/i[@class="icon docx"]`)
	if fileTypeDocxNode != nil {
		fileType = ".docx"
	}
	fileTypeDocNode := htmlquery.FindOne(detailDoc, `//div[@class="content"]/div[@class="ppttit"]/h2/i[@class="icon doc"]`)
	if fileTypeDocNode != nil {
		fileType = ".doc"
	}

	fmt.Println(fileType)
	if len(fileType) == 0 {
		return errors.New("文档类型不是doc，跳过")
	}

	// 所需点券
	pointsNode := htmlquery.FindOne(detailDoc, `//div[@class="content"]/div[@class="maindiv"]/div[@id="shoufeitishi"]/div[@class="btn-1"]/em[@class="point1"]`)
	if pointsNode == nil {
		return errors.New("没有点券div")
	}
	pointsText := htmlquery.InnerText(pointsNode)
	// 去除空格
	pointsText = strings.Replace(pointsText, " ", "", -1)
	// 去除换行符
	pointsText = strings.Replace(pointsText, "\n", "", -1)
	// fmt.Println(pointsText)
	pointsText = strings.ReplaceAll(pointsText, "个点券", "")

	points, err := strconv.Atoi(pointsText)
	if err != nil {
		return err
	}
	fmt.Println(points)
	// os.Exit(1)
	if points > 0 {
		return errors.New("需要点券下载")
	}

	filePath := "F:\\workspace\\www.czwlzx.cn\\www.czwlzx.cn\\" + title + fileType
	_, err = os.Stat(filePath)
	if err == nil {
		return errors.New("文档已下载过，跳过")
	}
	CzWlZxDownloadUrl := fmt.Sprintf("http://www.czwlzx.cn/Common/ShowDownloadUrl.aspx?urlid=0&id=%d", id)
	fmt.Println(CzWlZxDownloadUrl)

	fmt.Println("=======开始下载========")
	err = downloadCzWlZx(CzWlZxDownloadUrl, filePath, detailDocUrl)
	if err != nil {
		return err
	}
	fmt.Println("=======完成下载========")
	// DownLoadCzWlZxTimeSleep := rand.Intn(10)
	DownLoadCzWlZxTimeSleep := 10
	for i := 1; i <= DownLoadCzWlZxTimeSleep; i++ {
		time.Sleep(time.Second)
		fmt.Println("id="+strconv.Itoa(id)+"===========下载", title, "成功，暂停", DownLoadCzWlZxTimeSleep, "秒，倒计时", i, "秒===========")
	}
	return nil
}

func CzWlZxDetailDoc(requestUrl string, referer string) (doc *html.Node, err error) {
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
	if CzWlZxEnableHttpProxy {
		client = CzWlZxSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	// req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", CzWlZxCookie)
	req.Header.Set("Host", "www.czwlzx.cn")
	req.Header.Set("Referer", referer)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
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

func downloadCzWlZx(attachmentUrl string, filePath string, referer string) error {
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
	if CzWlZxEnableHttpProxy {
		client = CzWlZxSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.czwlzx.cn")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", referer)
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
