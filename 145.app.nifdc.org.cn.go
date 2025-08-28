package main

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	_ "os"
	"path/filepath"
	_ "path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	_ "golang.org/x/net/html"
)

var QxTjBzEnableHttpProxy = false
var QxTjBzHttpProxyUrl = "111.225.152.186:8089"
var QxTjBzHttpProxyUrlArr = make([]string, 0)

func QxTjBzHttpProxy() error {
	pageMax := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
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
					QxTjBzHttpProxyUrlArr = append(QxTjBzHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					QxTjBzHttpProxyUrlArr = append(QxTjBzHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func QxTjBzSetHttpProxy() (httpclient *http.Client) {
	if QxTjBzHttpProxyUrl == "" {
		if len(QxTjBzHttpProxyUrlArr) <= 0 {
			err := QxTjBzHttpProxy()
			if err != nil {
				QxTjBzSetHttpProxy()
			}
		}
		QxTjBzHttpProxyUrl = QxTjBzHttpProxyUrlArr[0]
		if len(QxTjBzHttpProxyUrlArr) >= 2 {
			QxTjBzHttpProxyUrlArr = QxTjBzHttpProxyUrlArr[1:]
		} else {
			QxTjBzHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(QxTjBzHttpProxyUrl)
	ProxyURL, _ := url.Parse(QxTjBzHttpProxyUrl)
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

type QueryQxTjBzListFormData struct {
	index1 int
	index2 int
}

type ViewQxTjBzFormData struct {
	type1 string
	id    string
}

var QxTjBzCookie = "JSESSIONID=20236F201F459E82C745953D125A3767; Hm_lvt_538f18182ca76c3c0acf71dbed622e93=1753434014; HMACCOUNT=1CCD0111717619C6; Hm_lpvt_538f18182ca76c3c0acf71dbed622e93=1753682828"

// 下载中国食品药品检定研究院-器械非采标推荐性行业标准文档
// @Title 下载中国食品药品检定研究院-器械非采标推荐性行业标准文档
// @Description http://app.nifdc.org.cn/jianybz/jybzTwoGj.do?formAction=listQxtjbz/，下载中国食品药品检定研究院-器械非采标推荐性行业标准文档
func main() {
	pageListUrl := "http://app.nifdc.org.cn/jianybz/jybzTwoGj.do?formAction=listQxtjbz"
	fmt.Println(pageListUrl)
	startIndex := 0
	isPageListGo := true
	for isPageListGo {
		queryQxTjBzListFormData := QueryQxTjBzListFormData{}
		if startIndex > 0 {
			queryQxTjBzListFormData.index1 = startIndex
			queryQxTjBzListFormData.index2 = startIndex - 1
		}

		queryQxTjBzListDoc, err := QueryQxTjBzList(pageListUrl, queryQxTjBzListFormData)
		if err != nil {
			fmt.Println(err)
			break
		}
		// /html/body/div[2]/div[2]/div/div/div[2]/table/tbody/tr[1]
		trNodes := htmlquery.Find(queryQxTjBzListDoc, `//html/body/div[2]/div[2]/div/div/div[2]/table/tbody/tr`)
		if len(trNodes) >= 1 {
			for _, trNode := range trNodes {
				fmt.Println("=====================开始处理数据=========================")
				codeNode := htmlquery.FindOne(trNode, `./td[3]/div/p`)
				code := htmlquery.InnerText(codeNode)
				code = strings.TrimSpace(code)
				code = strings.ReplaceAll(code, "/", "-")
				fmt.Println(code)

				titleNode := htmlquery.FindOne(trNode, `./td[2]/div/@title`)
				title := htmlquery.InnerText(titleNode)
				title = strings.TrimSpace(title)
				title = strings.ReplaceAll(title, "-", "")
				title = strings.ReplaceAll(title, " ", "")
				title = strings.ReplaceAll(title, "|", "-")
				title = strings.ReplaceAll(title, "/", "-")
				fmt.Println(title)

				filePath := "../app.nifdc.org.cn/" + title + "(" + code + ")" + ".pdf"
				fmt.Println(filePath)

				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				fmt.Println("=======开始下载========")

				buttonNode := htmlquery.FindOne(trNode, `./td[5]/div/button`)
				// ckbz('jybzTwoGj.do?formAction=viewBzpdfjs&type=qxtjbz&id=2c9048cc981213f901982c1d26bc4578')
				clickText := htmlquery.SelectAttr(buttonNode, "onclick")
				clickTextArray := strings.Split(clickText, "id=")
				id := strings.ReplaceAll(clickTextArray[1], "')", "")
				//fmt.Println(id)

				viewQxTjBzUrl := "http://app.nifdc.org.cn/jianybz/jybzTwoGj.do?formAction=viewBzpdfjs"
				//fmt.Println(viewQxTjBzUrl)
				viewQxTjBzFormData := ViewQxTjBzFormData{
					type1: "qxtjbz",
					id:    id,
				}

				viewQxTjBzDoc, err := viewQxTjBz(viewQxTjBzUrl, viewQxTjBzFormData)
				if err != nil {
					fmt.Println(err)
					continue
				}

				// /html/body/iframe
				downloadQxTjBzUrlNode := htmlquery.FindOne(viewQxTjBzDoc, `//html/body/iframe/@src`)
				downloadQxTjBzUrlText := htmlquery.InnerText(downloadQxTjBzUrlNode)
				downloadQxTjBzUrl := strings.ReplaceAll(downloadQxTjBzUrlText, "js/pdfjs2.12.313/web/viewer.html?file=", "")
				downloadQxTjBzUrl = strings.ReplaceAll(downloadQxTjBzUrl, "\\", "/")
				fmt.Println(downloadQxTjBzUrl)
				err = downloadQxTjBz(downloadQxTjBzUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "../app.nifdc.org.cn", "../upload.doc88.com/hbba.sacinfo.org.cn")
				err = QxTjBzCopyFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======完成下载========")
				//DownLoadQxTjBzTimeSleep := 10
				DownLoadQxTjBzTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadQxTjBzTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("title="+title+"===========下载", title, "成功 startIndex="+strconv.Itoa(startIndex)+"====，暂停", DownLoadQxTjBzTimeSleep, "秒，倒计时", i, "秒===========")
				}
			}
			DownLoadQxTjBzPageTimeSleep := 10
			// DownLoadQxTjBzPageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadQxTjBzPageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("startIndex="+strconv.Itoa(startIndex)+"========= 暂停", DownLoadQxTjBzPageTimeSleep, "秒 倒计时", i, "秒===========")
			}
			startIndex++
		} else {
			isPageListGo = false
			startIndex = 1
			break
		}
	}
}

func QueryQxTjBzList(requestUrl string, queryQxTjBzListFormData QueryQxTjBzListFormData) (doc *html.Node, err error) {
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
	if QxTjBzEnableHttpProxy {
		client = QxTjBzSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("index", strconv.Itoa(queryQxTjBzListFormData.index1))
	postData.Add("index", strconv.Itoa(queryQxTjBzListFormData.index2))
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", QxTjBzCookie)
	req.Header.Set("Host", "app.nifdc.org.cn")
	req.Header.Set("Origin", "http://app.nifdc.org.cn")
	req.Header.Set("Referer", "http://app.nifdc.org.cn/jianybz/jybzTwoGj.do?formAction=listQxtjbz")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
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
	doc, err = decodeAndParseHTMLQxTjBz(string(bodyBytes))
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func viewQxTjBz(requestUrl string, viewQxTjBzFormData ViewQxTjBzFormData) (doc *html.Node, err error) {
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
	if QxTjBzEnableHttpProxy {
		client = QxTjBzSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("type", viewQxTjBzFormData.type1)
	postData.Add("id", viewQxTjBzFormData.id)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		return doc, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", QxTjBzCookie)
	req.Header.Set("Host", "app.nifdc.org.cn")
	req.Header.Set("Origin", "http://app.nifdc.org.cn")
	req.Header.Set("Referer", "http://app.nifdc.org.cn/jianybz/jybzTwoGj.do?formAction=listQxtjbz")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
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
	doc, err = decodeAndParseHTMLQxTjBz(string(bodyBytes))
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func decodeAndParseHTMLQxTjBz(gb2312Content string) (*html.Node, error) {
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

func downloadQxTjBz(requestUrl string, filePath string) error {
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
	if QxTjBzEnableHttpProxy {
		client = QxTjBzSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", QxTjBzCookie)
	req.Header.Set("Host", "app.nifdc.org.cn")
	req.Header.Set("Origin", "http://app.nifdc.org.cn")
	req.Header.Set("Referer", "http://app.nifdc.org.cn/jianybz/jybzTwoGj.do?formAction=listQxtjbz")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
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

func QxTjBzCopyFile(src, dst string) (err error) {
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
