package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
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

var Bz100EnableHttpProxy = false
var Bz100HttpProxyUrl = "111.225.152.186:8089"
var Bz100HttpProxyUrlArr = make([]string, 0)

func Bz100HttpProxy() error {
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
					Bz100HttpProxyUrlArr = append(Bz100HttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					Bz100HttpProxyUrlArr = append(Bz100HttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func Bz100SetHttpProxy() (httpclient *http.Client) {
	if Bz100HttpProxyUrl == "" {
		if len(Bz100HttpProxyUrlArr) <= 0 {
			err := Bz100HttpProxy()
			if err != nil {
				Bz100SetHttpProxy()
			}
		}
		Bz100HttpProxyUrl = Bz100HttpProxyUrlArr[0]
		if len(Bz100HttpProxyUrlArr) >= 2 {
			Bz100HttpProxyUrlArr = Bz100HttpProxyUrlArr[1:]
		} else {
			Bz100HttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(Bz100HttpProxyUrl)
	ProxyURL, _ := url.Parse(Bz100HttpProxyUrl)
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

type QueryBz100ListFormData struct {
	searchBtnTJ      string
	a200             string
	pager_pageNumber int
	pager_orderBy    string
	pager_orderType  string
}

type Bz100DownloadFormData struct {
	id     string
	sid    string
	ouidop string
	dbuuid string
}

var Bz100Cookie = "JSESSIONID=F1B22DE16B52CBC0A1A151DA191B8D10.z"

// 下载山东省地方标准
// @Title 下载山东省地方标准
// @Description https://www.bz100.cn/，下载山东省地方标准
func main() {
	pageListUrl := "https://www.bz100.cn/member/standard/standard!getfreedb.action"
	fmt.Println(pageListUrl)
	startPage := 1
	isPageListGo := true
	for isPageListGo {
		queryBz100ListFormData := QueryBz100ListFormData{
			searchBtnTJ:      "no",
			a200:             "现行",
			pager_pageNumber: startPage,
			pager_orderBy:    "t_order",
			pager_orderType:  "asc",
		}

		queryBz100ListDoc, err := QueryBz100List(pageListUrl, queryBz100ListFormData)
		if err != nil {
			fmt.Println(err)
			break
		}

		trNodes := htmlquery.Find(queryBz100ListDoc, `//div[@class="search_wrap"]/div[@class="search_con"]/div[@class="con_list"]/table[@class="list_xq"]/tbody[@id="goaler"]/tr`)
		if len(trNodes) >= 1 {
			for _, trNode := range trNodes {
				fmt.Println("=====================开始处理数据=========================")
				codeNode := htmlquery.FindOne(trNode, `./td[2]/a/b/font`)
				code := htmlquery.InnerText(codeNode)
				code = strings.TrimSpace(code)
				code = strings.ReplaceAll(code, "/", "-")
				fmt.Println(code)

				titleNode := htmlquery.FindOne(trNode, `./td[3]`)
				title := htmlquery.InnerText(titleNode)
				title = strings.TrimSpace(title)
				title = strings.ReplaceAll(title, "-", "")
				title = strings.ReplaceAll(title, " ", "")
				title = strings.ReplaceAll(title, "|", "-")
				fmt.Println(title)

				filePath := "../www.bz100.cn/" + title + "(" + code + ")" + ".pdf"
				fmt.Println(filePath)

				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				fmt.Println("=======开始下载========")

				buttonNode := htmlquery.FindOne(trNode, `./td[7]/a`)
				// qwyulansub('5F8F007452E043CBBD465CC36A2805C6','402881b436e7105a0136e80033430002','DB37','8a81a0a55ac6a975015ac6c383f60004')
				clickText := htmlquery.SelectAttr(buttonNode, "onclick")
				clickText = strings.ReplaceAll(clickText, "qwyulansub(", "")
				clickText = strings.ReplaceAll(clickText, ")", "")
				clickText = strings.ReplaceAll(clickText, "'", "")
				clickTextArray := strings.Split(clickText, ",")
				id := clickTextArray[0]
				sid := clickTextArray[2]
				ouidop := clickTextArray[1]
				dbuuid := clickTextArray[3]

				bz100DownloadUrl := "https://www.bz100.cn/member/standard/standard!qypreviewdb.action"
				bz100DownloadFormData := Bz100DownloadFormData{
					id:     id,
					sid:    sid,
					ouidop: ouidop,
					dbuuid: dbuuid,
				}
				fmt.Println(bz100DownloadFormData)

				err = downloadBz100(bz100DownloadUrl, bz100DownloadFormData, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "../www.bz100.cn", "../upload.doc88.com/www.bz100.cn")
				err = Bz100CopyFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======完成下载========")
				//DownLoadBz100TimeSleep := 10
				DownLoadBz100TimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadBz100TimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("title="+title+"===========下载", title, "成功 startPage="+strconv.Itoa(startPage)+"====，暂停", DownLoadBz100TimeSleep, "秒，倒计时", i, "秒===========")
				}
			}
			DownLoadBz100PageTimeSleep := 10
			// DownLoadBz100PageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadBz100PageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("startPage="+strconv.Itoa(startPage)+"========= 暂停", DownLoadBz100PageTimeSleep, "秒 倒计时", i, "秒===========")
			}
			startPage++
		} else {
			isPageListGo = false
			startPage = 1
			break
		}
	}
}

func QueryBz100List(requestUrl string, queryBz100ListFormData QueryBz100ListFormData) (doc *html.Node, err error) {
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
	if Bz100EnableHttpProxy {
		client = Bz100SetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("searchBtnTJ", queryBz100ListFormData.searchBtnTJ)
	postData.Add("a200", queryBz100ListFormData.a200)
	postData.Add("pager_pageNumber", strconv.Itoa(queryBz100ListFormData.pager_pageNumber))
	postData.Add("pager_orderBy", queryBz100ListFormData.pager_orderBy)
	postData.Add("pager_orderType", queryBz100ListFormData.pager_orderType)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", Bz100Cookie)
	req.Header.Set("Host", "www.bz100.cn")
	req.Header.Set("Origin", "http://www.bz100.cn")
	req.Header.Set("Referer", "https://www.bz100.cn/member/standard/standard!getfreedb.action")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return doc, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
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

func downloadBz100(requestUrl string, bz100DownloadFormData Bz100DownloadFormData, filePath string) error {
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
	if Bz100EnableHttpProxy {
		client = Bz100SetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("id", bz100DownloadFormData.id)
	postData.Add("sid", bz100DownloadFormData.sid)
	postData.Add("ouidop", bz100DownloadFormData.ouidop)
	postData.Add("dbuuid", bz100DownloadFormData.dbuuid)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", Bz100Cookie)
	req.Header.Set("Host", "www.bz100.cn")
	req.Header.Set("Origin", "http://www.bz100.cn")
	req.Header.Set("Referer", "https://www.bz100.cn")
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

func Bz100CopyFile(src, dst string) (err error) {
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
