package main

import (
	"bytes"
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

var FlkEnableHttpProxy = false
var FlkHttpProxyUrl = "111.225.152.186:8089"
var FlkHttpProxyUrlArr = make([]string, 0)

func FlkHttpProxy() error {
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
					FlkHttpProxyUrlArr = append(FlkHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					FlkHttpProxyUrlArr = append(FlkHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func FlkSetHttpProxy() (httpclient *http.Client) {
	if FlkHttpProxyUrl == "" {
		if len(FlkHttpProxyUrlArr) <= 0 {
			err := FlkHttpProxy()
			if err != nil {
				FlkSetHttpProxy()
			}
		}
		FlkHttpProxyUrl = FlkHttpProxyUrlArr[0]
		if len(FlkHttpProxyUrlArr) >= 2 {
			FlkHttpProxyUrlArr = FlkHttpProxyUrlArr[1:]
		} else {
			FlkHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(FlkHttpProxyUrl)
	ProxyURL, _ := url.Parse(FlkHttpProxyUrl)
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

type QueryFlkListRequestPayload struct {
	FlfgCodeId    []int             `json:"flfgCodeId"`
	Gbrq          []string          `json:"gbrq"`
	GbrqYear      []string          `json:"gbrqYear"`
	OrderByParam  map[string]string `json:"orderByParam"`
	PageNum       int               `json:"pageNum"`
	PageSize      int               `json:"pageSize"`
	SearchContent string            `json:"searchContent"`
	SearchRange   int               `json:"searchRange"`
	SearchType    int               `json:"searchType"`
	Sxrq          []string          `json:"sxrq"`
	Sxx           []int             `json:"sxx"`
	ZdjgCodeId    []int             `json:"zdjgCodeId"`
}

var FlkCookie = "Hm_lvt_54434aa6770b6d9fef104d146430b53b=1754290987; wzws_sessionid=gmZhYjg1ZqBowjVugDEyMy45LjI1LjEwNYFhMTZkZGE="

// 下载国家法律法规数据库文档
// @Title 下载国家法律法规数据库文档
// @Description https://flk.npc.gov.cn/，下载国家法律法规数据库文档
func main() {
	pageListUrl := "https://flk.npc.gov.cn/law-search/search/list"
	fmt.Println(pageListUrl)
	page := 46
	maxPage := 100
	rows := 100
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}

		queryFlkListRequestPayload := QueryFlkListRequestPayload{
			FlfgCodeId:    []int{},
			Gbrq:          []string{},
			GbrqYear:      []string{},
			OrderByParam:  map[string]string{"order": "-1", "sort": ""},
			PageNum:       page,
			PageSize:      rows,
			SearchContent: "",
			SearchRange:   1,
			SearchType:    2,
			Sxrq:          []string{},
			Sxx:           []int{},
			ZdjgCodeId:    []int{},
		}
		queryFlkListResponseDataStandardInfos, err := QueryFlkList(pageListUrl, queryFlkListRequestPayload)
		if err != nil {
			fmt.Println(err)
			break
		}
		for _, flk := range queryFlkListResponseDataStandardInfos {
			fmt.Println("=====================开始处理数据 page = ", page, "=========================")

			gbrq := flk.Gbrq
			fmt.Println(gbrq)

			title := flk.Title
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "--", "-")
			fmt.Println(title)

			filePath := "../flk.npc.gov.cn/flk.npc.gov.cn/" + title + "-" + flk.Flxz + "(" + gbrq + ")" + ".docx"
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载========")

			requestFlkDownloadUrl := fmt.Sprintf("https://flk.npc.gov.cn/law-search/download/pc?format=docx&bbbs=%s", flk.Bbbs)
			// fmt.Println(requestFlkDownloadUrl)
			requestFlkDownloadReferer := fmt.Sprintf("https://flk.npc.gov.cn/detail?id=%s&fileId=&type=&title=%s", flk.Bbbs, flk.Title)
			queryFlkDownloadUrlResponseData, err := QueryFlkDownloadUrl(requestFlkDownloadUrl, requestFlkDownloadReferer)
			if err != nil {
				fmt.Println(err)
				break
			}

			downloadUrl := queryFlkDownloadUrlResponseData.Url
			// fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")

			requestFlkDownloadRefererUrl := fmt.Sprintf("https://flk.npc.gov.cn/detail?id=%s&fileId=&type=&title=%s", flk.Bbbs, flk.Title)

			err = downloadFlk(downloadUrl, requestFlkDownloadRefererUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}

			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "../flk.npc.gov.cn/flk.npc.gov.cn", "../flk.npc.gov.cn/temp-flk.npc.gov.cn")
			err = copyFlkFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			//DownLoadFlkTimeSleep := 10
			DownLoadFlkTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadFlkTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadFlkTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		// DownLoadFlkPageTimeSleep := 10
		DownLoadFlkPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadFlkPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadFlkPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

type QueryFlkListResponse struct {
	Code          int                        `json:"code"`
	Rows          []QueryFlkListResponseRows `json:"rows"`
	Msg           string                     `json:"msg"`
	SearchContent string                     `json:"searchContent"`
	SearchType    int                        `json:"searchType"`
	Total         int                        `json:"total"`
}

type QueryFlkListResponseRows struct {
	Bbbs     string `json:"bbbs"`
	Title    string `json:"title"`
	ZdjgName string `json:"zdjgName"`
	Flxz     string `json:"flxz"`
	Gbrq     string `json:"gbrq"`
	Sxrq     string `json:"sxrq"`
}

func QueryFlkList(requestUrl string, queryFlkListRequestPayload QueryFlkListRequestPayload) (queryFlkListResponseRows []QueryFlkListResponseRows, err error) {
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
	if FlkEnableHttpProxy {
		client = FlkSetHttpProxy()
	}
	// 将数据编码为JSON格式
	queryFlkListRequestPayloadJson, err := json.Marshal(queryFlkListRequestPayload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 创建字符串读取器，这是http.Request需要的类型
	body := bytes.NewReader(queryFlkListRequestPayloadJson)
	req, err := http.NewRequest("POST", requestUrl, body) //建立连接

	queryFlkListResponse := QueryFlkListResponse{}
	if err != nil {
		return queryFlkListResponseRows, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Cookie", FlkCookie)
	req.Header.Set("Host", "flk.npc.gov.cn")
	req.Header.Set("Origin", "https://flk.npc.gov.cn")
	req.Header.Set("Referer", "https://flk.npc.gov.cn/search")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryFlkListResponseRows, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryFlkListResponseRows, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryFlkListResponseRows, err
	}
	err = json.Unmarshal(respBytes, &queryFlkListResponse)
	if err != nil {
		return queryFlkListResponseRows, err
	}
	queryFlkListResponseRows = queryFlkListResponse.Rows
	return queryFlkListResponseRows, nil
}

type QueryFlkDownloadUrlResponse struct {
	Code int                             `json:"code"`
	Data QueryFlkDownloadUrlResponseData `json:"data"`
	Msg  string                          `json:"msg"`
}

type QueryFlkDownloadUrlResponseData struct {
	Url   string `json:"url"`
	UrlIn string `json:"urlIn"`
}

func QueryFlkDownloadUrl(requestUrl string, referer string) (queryFlkDownloadUrlResponseData QueryFlkDownloadUrlResponseData, err error) {
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
	if FlkEnableHttpProxy {
		client = FlkSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	queryFlkDownloadUrlResponse := QueryFlkDownloadUrlResponse{}
	if err != nil {
		return queryFlkDownloadUrlResponseData, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "flk.npc.gov.cn")
	req.Header.Set("Origin", "https://flk.npc.gov.cn")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryFlkDownloadUrlResponseData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryFlkDownloadUrlResponseData, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryFlkDownloadUrlResponseData, err
	}
	err = json.Unmarshal(respBytes, &queryFlkDownloadUrlResponse)
	if err != nil {
		return queryFlkDownloadUrlResponseData, err
	}
	queryFlkDownloadUrlResponseData = queryFlkDownloadUrlResponse.Data
	return queryFlkDownloadUrlResponseData, nil
}

func downloadFlk(attachmentUrl string, referer string, filePath string) error {
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
	if FlkEnableHttpProxy {
		client = FlkSetHttpProxy()
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
	req.Header.Set("Host", "flk.npc.gov.cn")
	//req.Header.Set("Referer", referer)
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

func copyFlkFile(src, dst string) (err error) {
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
