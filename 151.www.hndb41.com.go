package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

var HnDb41EnableHttpProxy = false
var HnDb41HttpProxyUrl = "111.225.152.186:8089"
var HnDb41HttpProxyUrlArr = make([]string, 0)

func HnDb41HttpProxy() error {
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
					HnDb41HttpProxyUrlArr = append(HnDb41HttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					HnDb41HttpProxyUrlArr = append(HnDb41HttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func HnDb41SetHttpProxy() (httpclient *http.Client) {
	if HnDb41HttpProxyUrl == "" {
		if len(HnDb41HttpProxyUrlArr) <= 0 {
			err := HnDb41HttpProxy()
			if err != nil {
				HnDb41SetHttpProxy()
			}
		}
		HnDb41HttpProxyUrl = HnDb41HttpProxyUrlArr[0]
		if len(HnDb41HttpProxyUrlArr) >= 2 {
			HnDb41HttpProxyUrlArr = HnDb41HttpProxyUrlArr[1:]
		} else {
			HnDb41HttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(HnDb41HttpProxyUrl)
	ProxyURL, _ := url.Parse(HnDb41HttpProxyUrl)
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

type QueryHnDb41ListFormData struct {
	snumber   string
	sname     string
	pageIndex int
	sortField string
	sortOrder string
	pageSize  int
}

var HnDb41Cookie = "SESSIONID=4151517F5CEC91301C4BC202BA5BB5B0"

// 下载河南省地方标准文档
// @Title 下载河南省地方标准文档
// @Description http://www.hndb41.com/，下载河南省地方标准文档
func main() {
	pageListUrl := "http://www.hndb41.com/bzsp_getStandardPermitList.action?type=1&acc=1"
	fmt.Println(pageListUrl)
	page := 0
	maxPage := 289
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		queryHnDb41ListFormData := QueryHnDb41ListFormData{
			snumber:   "",
			sname:     "",
			pageIndex: page,
			sortField: "id",
			sortOrder: "desc",
			pageSize:  10,
		}
		queryHnDb41ListResponseData, err := QueryHnDb41List(pageListUrl, queryHnDb41ListFormData)
		if err != nil {
			HnDb41HttpProxyUrl = ""
			fmt.Println(err)
		}
		for _, data := range queryHnDb41ListResponseData {
			fmt.Println("=====================开始处理数据 page = ", page, "=========================")
			code := data.Snumber
			code = strings.ReplaceAll(code, "/", "-")
			code = strings.ReplaceAll(code, "—", "-")
			fmt.Println(code)

			title := data.Sname
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "--", "-")
			fmt.Println(title)

			filePath := "../www.hndb41.com/" + title + "(" + code + ")" + ".pdf"
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载========")

			downloadUrl := fmt.Sprintf("http://www.hndb41.com/UeditorPathfile/new/%s.pdf", strings.Replace(data.Snumber, "/", "_", 1))
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")

			err = downloadHnDb41(downloadUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "../www.hndb41.com", "../temp-www.hndb41.com")
			err = copyHnDb41File(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			//DownLoadHnDb41TimeSleep := 10
			DownLoadHnDb41TimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadHnDb41TimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadHnDb41TimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadHnDb41PageTimeSleep := 10
		// DownLoadHnDb41PageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadHnDb41PageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadHnDb41PageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

type QueryHnDb41ListResponse struct {
	Total int                           `json:"total"`
	Data  []QueryHnDb41ListResponseData `json:"data"`
}

type QueryHnDb41ListResponseData struct {
	Snumber string `json:"snumber"`
	Sname   string `json:"sname"`
	Id      int    `json:"id"`
}

func QueryHnDb41List(requestUrl string, queryHnDb41ListFormData QueryHnDb41ListFormData) (queryHnDb41ListResponseData []QueryHnDb41ListResponseData, err error) {
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
	if HnDb41EnableHttpProxy {
		client = HnDb41SetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("snumber", queryHnDb41ListFormData.snumber)
	postData.Add("sname", queryHnDb41ListFormData.sname)
	postData.Add("pageIndex", strconv.Itoa(queryHnDb41ListFormData.pageIndex))
	postData.Add("sortField", queryHnDb41ListFormData.sortField)
	postData.Add("sortOrder", queryHnDb41ListFormData.sortOrder)
	postData.Add("pageSize", strconv.Itoa(queryHnDb41ListFormData.pageSize))
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	queryHnDb41ListResponse := QueryHnDb41ListResponse{}
	if err != nil {
		return queryHnDb41ListResponseData, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", HnDb41Cookie)
	req.Header.Set("Host", "www.hndb41.com")
	req.Header.Set("Origin", "http://www.hndb41.com")
	req.Header.Set("Referer", "http://www.hndb41.com/public/StandardPermitManager2.jsp?bzbh=&bzname=")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryHnDb41ListResponseData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryHnDb41ListResponseData, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryHnDb41ListResponseData, err
	}
	err = json.Unmarshal(respBytes, &queryHnDb41ListResponse)
	if err != nil {
		return queryHnDb41ListResponseData, err
	}
	queryHnDb41ListResponseData = queryHnDb41ListResponse.Data
	return queryHnDb41ListResponseData, nil
}

func downloadHnDb41(attachmentUrl string, filePath string) error {
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
	if HnDb41EnableHttpProxy {
		client = HnDb41SetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", HnDb41Cookie)
	req.Header.Set("Host", "www.hndb41.com")
	req.Header.Set("Referer", "http://www.hndb41.com/public/StandardPermitManager2.jsp?bzbh=&bzname=")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 如果访问失败 就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}

	// 创建一个文件用于保存
	fileDiv := filepath.Dir(filePath)
	if _, err = os.Stat(fileDiv); err != nil {
		if os.MkdirAll(fileDiv, 0o777) != nil {
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

func copyHnDb41File(src, dst string) (err error) {
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
