package main

import (
	// "bytes"
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

var GwSoSoEnableHttpProxy = false
var GwSoSoHttpProxyUrl = "111.225.152.186:8089"
var GwSoSoHttpProxyUrlArr = make([]string, 0)

func GwSoSoHttpProxy() error {
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
					GwSoSoHttpProxyUrlArr = append(GwSoSoHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					GwSoSoHttpProxyUrlArr = append(GwSoSoHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func GwSoSoSetHttpProxy() (httpclient *http.Client) {
	if GwSoSoHttpProxyUrl == "" {
		if len(GwSoSoHttpProxyUrlArr) <= 0 {
			err := GwSoSoHttpProxy()
			if err != nil {
				GwSoSoSetHttpProxy()
			}
		}
		GwSoSoHttpProxyUrl = GwSoSoHttpProxyUrlArr[0]
		if len(GwSoSoHttpProxyUrlArr) >= 2 {
			GwSoSoHttpProxyUrlArr = GwSoSoHttpProxyUrlArr[1:]
		} else {
			GwSoSoHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(GwSoSoHttpProxyUrl)
	ProxyURL, _ := url.Parse(GwSoSoHttpProxyUrl)
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
			ResponseHeaderTimeout: time.Second * 30,
		},
	}
	return httpclient
}

type QueryGwSoSoListFormData struct {
	Start   int `json:"start"`
	Limit   int `json:"limit"`
	PageCls int `json:"pagecls"`
}

var GwSoSoCookie = "Hm_lvt_b3ecdb0e91dc1c234e7f59ad61980ec7=1773908325,1773921817,1774447409,1775015186; HMACCOUNT=1CCD0111717619C6; gongwen=3efa6df1-55f7-4919-b2f2-b2cd9a95529f; Hm_lpvt_b3ecdb0e91dc1c234e7f59ad61980ec7=1775015219"

// 下载公文搜文档
// @Title 下载公文搜文档
// @Description https://www.gwsoso.com/，下载公文搜文档
func main() {
	pageListUrl := "https://www.gwsoso.com/docs/gettop?t=0.42132661413839645"
	start := 50
	isPageListGo := true
	for isPageListGo {
		queryGwSoSoListFormData := QueryGwSoSoListFormData{
			Start:   start,
			Limit:   50,
			PageCls: 32,
		}
		queryGwSoSoListResponseData, err := QueryGwSoSoList(pageListUrl, queryGwSoSoListFormData)
		if err != nil {
			fmt.Println(err)
			isPageListGo = false
			break
		}
		if len(queryGwSoSoListResponseData) <= 0 {
			isPageListGo = false
			break
		}
		for _, data := range queryGwSoSoListResponseData {
			fmt.Println("===============开始处理数据 start = ", start, " data记录数量 = ", len(queryGwSoSoListResponseData), "==================")
			fmt.Println(data.Id)

			title := data.Title
			fmt.Println(data.Title)
			if strings.Index(data.Title, "doc") == -1 && strings.Index(data.Title, "pdf") == -1 {
				fmt.Println("文档不是doc、pdf文档，跳过")
				continue
			}
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "--", "-")
			title = strings.ReplaceAll(title, ".docx", "")
			title = strings.ReplaceAll(title, ".doc", "")
			title = strings.ReplaceAll(title, ".pdf", "")

			filePath := "D:\\workspace\\www.gwsoso.com\\www.gwsoso.com\\" + title + ".pdf"
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载========")
			downloadUrl := fmt.Sprintf("https://www.gwsoso.com/docs/pdf/%s.pdf", data.Id)
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")

			err = downloadGwSoSo(downloadUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// 查看文件大小，如果是空文件，则删除
			fileInfo, err := os.Stat(filePath)
			if err == nil && fileInfo.Size() == 0 {
				fmt.Println("空文件删除")
				err = os.Remove(filePath)
				isPageListGo = false
				break
			}
			if err != nil {
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "www.gwsoso.com\\www.gwsoso.com", "www.gwsoso.com\\temp-www.gwsoso.com")
			err = copyGwSoSoFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			//DownLoadGwSoSoTimeSleep := 10
			DownLoadGwSoSoTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadGwSoSoTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("start="+strconv.Itoa(start)+",filePath="+filePath+"===========下载成功 暂停", DownLoadGwSoSoTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}

		start = start + 50
		DownLoadGwSoSoPageTimeSleep := 10
		// DownLoadGwSoSoPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadGwSoSoPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("start="+strconv.Itoa(start)+"=========== 暂停", DownLoadGwSoSoPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
	}
}

type QueryGwSoSoListResponse struct {
	Success bool                          `json:"success"`
	Data    []QueryGwSoSoListResponseData `json:"data"`
	Total   int                           `json:"total"`
}
type QueryGwSoSoListResponseData struct {
	Id    string `json:"Id"`
	Title string `json:"Title"`
}

func QueryGwSoSoList(requestUrl string, queryGwSoSoListFormData QueryGwSoSoListFormData) (queryGwSoSoListResponseData []QueryGwSoSoListResponseData, err error) {
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
	if GwSoSoEnableHttpProxy {
		client = GwSoSoSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("start", strconv.Itoa(queryGwSoSoListFormData.Start))
	postData.Add("limit", strconv.Itoa(queryGwSoSoListFormData.Limit))
	postData.Add("pagecls", strconv.Itoa(queryGwSoSoListFormData.PageCls))
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	queryGwSoSoListResponse := QueryGwSoSoListResponse{}
	if err != nil {
		return queryGwSoSoListResponseData, err
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
	req.Header.Set("Cookie", GwSoSoCookie)
	req.Header.Set("Host", "www.gwsoso.com")
	req.Header.Set("Origin", "https://www.gwsoso.com")
	req.Header.Set("Referer", "https://www.gwsoso.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryGwSoSoListResponseData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryGwSoSoListResponseData, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryGwSoSoListResponseData, err
	}
	err = json.Unmarshal(respBytes, &queryGwSoSoListResponse)
	if err != nil {
		return queryGwSoSoListResponseData, err
	}
	queryGwSoSoListResponseData = queryGwSoSoListResponse.Data
	return queryGwSoSoListResponseData, nil
}

func downloadGwSoSo(attachmentUrl string, filePath string) error {
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
	if GwSoSoEnableHttpProxy {
		client = GwSoSoSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", GwSoSoCookie)
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("Host", "www.gwsoso.com")
	req.Header.Set("Referer", "https://www.gwsoso.com/")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
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

func copyGwSoSoFile(src, dst string) (err error) {
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

	// 创建一个文件用于保存
	fileDiv := filepath.Dir(dst)
	if _, err = os.Stat(fileDiv); err != nil {
		if os.MkdirAll(fileDiv, 0o777) != nil {
			return err
		}
	}

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
