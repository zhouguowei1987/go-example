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
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	SjFgEnableHttpProxy = false
	SjFgHttpProxyUrl    = "111.225.152.186:8089"
)

func SjFgSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(SjFgHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type SjFgCategory struct {
	validLevel string
	validLevelChild string
	departDuty int
}

var sjFgCategory = []SjFgCategory{
	{
		validLevel: "法律",
		validLevelChild: "以市场监管部门为主要执行部门",
		departDuty: 4001,
	},
	{
		validLevel: "法律",
		validLevelChild: "其他涉及市场监管部门职能",
		departDuty: 4002,
	},
	{
		validLevel: "法律",
		validLevelChild: "行政执法普遍适用",
		departDuty: 4003,
	},
	{
		validLevel: "行政法规",
		validLevelChild: "以市场监管部门为主要执行部门",
		departDuty: 4001,
	},
	{
		validLevel: "行政法规",
		validLevelChild: "其他涉及市场监管部门职能",
		departDuty: 4002,
	},
	{
		validLevel: "行政法规",
		validLevelChild: "行政执法普遍适用",
		departDuty: 4003,
	},
	{
		validLevel: "规章",
		validLevelChild: "总局公布的规章",
		departDuty: 5001,
	},
	{
		validLevel: "规章",
		validLevelChild: "总局参与的联合规章",
		departDuty: 5003,
	},
}

type SjFgListFormData struct {
	pageNo      int
	pageSize    int
	searchScope string
	searchType  string
	timeValid   string
	lawType     string
	validLevel  string
	departDuty  int
	valid       int
	lawName     string
	startTime   string
	pubTime     string
}

type QueryLawByLawIdRequestPayload struct {
	Id interface{} `json:"id"`
}

var SjFgCookie = "HMACCOUNT=1CCD0111717619C6; Hm_lvt_54db9897e5a65f7a7b00359d86015d8d=1756176872,1756314720; Hm_lpvt_54db9897e5a65f7a7b00359d86015d8d=1757428720; __jsluid_s=06e6e8d62d686727e4c088a24e481331"

// ychEduSpider 市场监管法律法规规章数据库文档
// @Title 市场监管法律法规规章数据库文档
// @Description https://sjfg.samr.gov.cn/，市场监管法律法规规章数据库文档
func main() {
	for _, sjfg := range sjFgCategory {
		pageNo := 1
		pageSize := 10
		isPageListGo := true
		requestUrl := "https://sjfg.samr.gov.cn/law/law_search/getLawStore.do"
		for isPageListGo {
			sjFgListFormData := SjFgListFormData{
				pageNo:      pageNo,
				pageSize:    pageSize,
				searchScope: "标题",
				searchType:  "模糊查询",
				timeValid:   "",
				lawType:     "",
				validLevel:  sjfg.validLevel,
				departDuty:  sjfg.departDuty,
				valid:       1,
				lawName:     "",
				startTime:   "",
				pubTime:     "",
			}
			fmt.Println(sjFgListFormData)
			sjFgListResponsePage, err := SjFgList(requestUrl, sjFgListFormData)
			if err != nil {
				fmt.Println(err)
				break
			}
			maxPageNo := (sjFgListResponsePage.Count / pageSize) + 1
			if pageNo >= maxPageNo {
				isPageListGo = false
				break
			}
			for _, result := range sjFgListResponsePage.Result {
				fmt.Println("============================================================================")
				fmt.Println("=======当前页为：" + strconv.Itoa(pageNo) + "========类别=" + sjfg.validLevel + "-" + sjfg.validLevelChild)

				lawId := result[0]
				fmt.Println(lawId)
				queryLawByLawIdRequestPayload := QueryLawByLawIdRequestPayload{
					Id: lawId,
				}

				queryLawByLawIdUrl := "https://sjfg.samr.gov.cn/law/law_search/queryLawByLawId.do"
				queryLawByLawIdReferer := fmt.Sprintf("https://sjfg.samr.gov.cn/law/pageInfo/law_search_new.law_details?lawId=%s&label=1", lawId)
				queryLawByLawIdResponse, err := QueryLawByLawId(queryLawByLawIdUrl, queryLawByLawIdReferer, queryLawByLawIdRequestPayload)
				if err != nil {
					fmt.Println(err)
					continue
				}
				title := strings.TrimSpace(queryLawByLawIdResponse.Page.LawName)
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, " ", "")
				fmt.Println(title)

				// 下载docx文档
				docName := queryLawByLawIdResponse.Page.FileUrl
				fmt.Println(docName)
				if len(queryLawByLawIdResponse.Page.FileUrl) > 0 {
					docFilePath := "../sjfg.samr.gov.cn/" + title + "." + strings.Split(docName, ".")[1]
					fmt.Println(docFilePath)
					_, err = os.Stat(docFilePath)
					if err == nil {
						fmt.Println("pdf文档已下载过，跳过")
						continue
					}
					fmt.Println("=======开始下载doc文件========")
					downLoadUrl := fmt.Sprintf("https://sjfg.samr.gov.cn/law/file%s", docName)
					err = downloadSjFg(downLoadUrl, docFilePath)
					if err != nil {
						fmt.Println(err)
						continue
					}

					//复制文件
					tempFilePath := strings.ReplaceAll(docFilePath, "../sjfg.samr.gov.cn", "../upload.doc88.com/sjfg.samr.gov.cn")
					err = SjFgCopyFile(docFilePath, tempFilePath)
					if err != nil {
						fmt.Println(err)
						continue
					}

					fmt.Println("=======下载完成========")

					downloadSjFgDocSleep := rand.Intn(5)
					for i := 1; i <= downloadSjFgDocSleep; i++ {
						time.Sleep(time.Second)
						fmt.Println("page="+strconv.Itoa(pageNo)+"=======title=", title, "成功，====== 暂停", downloadSjFgDocSleep, "秒，倒计时", i, "秒===========")
					}
				}

				// 下载pdf文档
				pdfName := queryLawByLawIdResponse.Page.FilePath
				fmt.Println(pdfName)
				if len(queryLawByLawIdResponse.Page.FilePath) > 0 {
					pdfFilePath := "../sjfg.samr.gov.cn/" + title + "." + strings.Split(pdfName, ".")[1]
					fmt.Println(pdfFilePath)
					_, err = os.Stat(pdfFilePath)
					if err == nil {
						fmt.Println("pdf文档已下载过，跳过")
						continue
					}
					fmt.Println("=======开始下载pdf文件========")
					downLoadUrl := fmt.Sprintf("https://sjfg.samr.gov.cn/law/file%s", pdfName)
					err = downloadSjFg(downLoadUrl, pdfFilePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					//复制文件
					tempFilePath := strings.ReplaceAll(pdfFilePath, "../sjfg.samr.gov.cn", "../upload.doc88.com/sjfg.samr.gov.cn")
					err = SjFgCopyFile(pdfFilePath, tempFilePath)
					if err != nil {
						fmt.Println(err)
						continue
					}

					fmt.Println("=======下载完成========")

					downloadSjFgPdfSleep := rand.Intn(5)
					for i := 1; i <= downloadSjFgPdfSleep; i++ {
						time.Sleep(time.Second)
						fmt.Println("page="+strconv.Itoa(pageNo)+"=======title=", title, "成功，====== 暂停", downloadSjFgPdfSleep, "秒，倒计时", i, "秒===========")
					}
				}
			}
			// DownLoadSjFgPageTimeSleep := 10
			DownLoadSjFgPageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadSjFgPageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(pageNo)+"======标准类别=" + sjfg.validLevel + "-" + sjfg.validLevelChild+"====== 暂停", DownLoadSjFgPageTimeSleep, "秒 倒计时", i, "秒===========")
			}
			pageNo++
			if pageNo > maxPageNo {
				isPageListGo = false
				break
			}
		}
	}
}

type SjFgListResponse struct {
	AllCount int                  `json:"allCount"`
	Page     SjFgListResponsePage `json:"page"`
}
type SjFgListResponsePage struct {
	Count    int             `json:"count"`
	PageNo   int             `json:"pageNo"`
	PageSize int             `json:"pageSize"`
	Result   [][]interface{} `json:"result"`
}

func SjFgList(requestUrl string, sjFgListFormData SjFgListFormData) (sjFgListResponsePage SjFgListResponsePage, err error) {
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
	if SjFgEnableHttpProxy {
		client = SjFgSetHttpProxy()
	}
	sjFgListResponse := SjFgListResponse{}
	postData := url.Values{}
	postData.Add("pageNo", strconv.Itoa(sjFgListFormData.pageNo))
	postData.Add("pageSize", strconv.Itoa(sjFgListFormData.pageSize))
	postData.Add("searchScope", sjFgListFormData.searchScope)
	postData.Add("searchType", sjFgListFormData.searchType)
	postData.Add("timeValid", sjFgListFormData.timeValid)
	postData.Add("lawType", sjFgListFormData.lawType)
	postData.Add("validLevel", sjFgListFormData.validLevel)
	postData.Add("departDuty", strconv.Itoa(sjFgListFormData.departDuty))
	postData.Add("valid", strconv.Itoa(sjFgListFormData.valid))
	postData.Add("lawName", sjFgListFormData.lawName)
	postData.Add("startTime", sjFgListFormData.startTime)
	postData.Add("pubTime", sjFgListFormData.pubTime)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return sjFgListResponsePage, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", SjFgCookie)
	req.Header.Set("Host", "sjfg.samr.gov.cn")
	req.Header.Set("Origin", "https://sjfg.samr.gov.cn")
	req.Header.Set("Referer", "https://sjfg.samr.gov.cn/law/pageInfo/law_search_new.law?key=%E8%A7%84%E7%AB%A0")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"124\", \"Google Chrome\";v=\"124\", \"Not-A.Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return sjFgListResponsePage, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(respBytes))
	//os.Exit(1)
	if err != nil {
		return sjFgListResponsePage, err
	}
	err = json.Unmarshal(respBytes, &sjFgListResponse)
	if err != nil {
		return sjFgListResponsePage, err
	}
	sjFgListResponsePage = sjFgListResponse.Page
	return sjFgListResponsePage, nil
}

type QueryLawByLawIdResponse struct {
	Page QueryLawByLawIdResponsePage `json:"page"`
}
type QueryLawByLawIdResponsePage struct {
	FilePath string `json:"filePath"`
	FileUrl  string `json:"fileUrl"`
	LawName  string `json:"lawName"`
}

func QueryLawByLawId(requestUrl string, referer string, queryLawByLawIdRequestPayload QueryLawByLawIdRequestPayload) (queryLawByLawIdResponse QueryLawByLawIdResponse, err error) {
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
	if SjFgEnableHttpProxy {
		client = SjFgSetHttpProxy()
	}
	// 将数据编码为JSON格式
	queryLawByLawIdRequestPayloadJson, err := json.Marshal(queryLawByLawIdRequestPayload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 创建字符串读取器，这是http.Request需要的类型
	body := bytes.NewReader(queryLawByLawIdRequestPayloadJson)
	req, err := http.NewRequest("POST", requestUrl, body) //建立连接

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Cookie", SjFgCookie)
	req.Header.Set("Host", "sjfg.samr.gov.cn")
	req.Header.Set("Origin", "http://sjfg.samr.gov.cn")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"124\", \"Google Chrome\";v=\"124\", \"Not-A.Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryLawByLawIdResponse, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryLawByLawIdResponse, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryLawByLawIdResponse, err
	}
	err = json.Unmarshal(respBytes, &queryLawByLawIdResponse)
	if err != nil {
		return queryLawByLawIdResponse, err
	}
	return queryLawByLawIdResponse, nil
}

func downloadSjFg(attachmentUrl string, filePath string) error {
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
	if SjFgEnableHttpProxy {
		client = SjFgSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", SjFgCookie)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "sjfg.samr.gov.cn")
	req.Header.Set("Referer", "https://sjfg.samr.gov.cn")
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

func SjFgCopyFile(src, dst string) (err error) {
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
