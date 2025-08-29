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

var LawEnableHttpProxy = false
var LawHttpProxyUrl = "111.225.152.186:8089"
var LawHttpProxyUrlArr = make([]string, 0)

func LawHttpProxy() error {
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
					LawHttpProxyUrlArr = append(LawHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					LawHttpProxyUrlArr = append(LawHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func LawSetHttpProxy() (httpclient *http.Client) {
	if LawHttpProxyUrl == "" {
		if len(LawHttpProxyUrlArr) <= 0 {
			err := LawHttpProxy()
			if err != nil {
				LawSetHttpProxy()
			}
		}
		LawHttpProxyUrl = LawHttpProxyUrlArr[0]
		if len(LawHttpProxyUrlArr) >= 2 {
			LawHttpProxyUrlArr = LawHttpProxyUrlArr[1:]
		} else {
			LawHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(LawHttpProxyUrl)
	ProxyURL, _ := url.Parse(LawHttpProxyUrl)
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

type QueryLawListRequestPayload struct {
	NONCE         string                            `json:"NONCE"`
	SIGN          string                            `json:"SIGN"`
	SIGN_TYPE     string                            `json:"SIGN_TYPE"`
	TIMESTAMP     string                            `json:"TIMESTAMP"`
	Current       int                               `json:"current"`
	FileType      int                               `json:"fileType"`
	IsLikeSearch  int                               `json:"isLikeSearch"`
	LevelCodes    []string                          `json:"levelCodes"`
	NeedHighLight bool                              `json:"needHighLight"`
	Orders        []QueryLawListRequestPayloadOrder `json:"orders"`
	SearchType    int                               `json:"searchType"`
	Size          int                               `json:"size"`
}

type QueryLawListRequestPayloadOrder struct {
	Column string `json:"column"`
	Asc    bool   `json:"asc"`
}

type LawListCategory struct {
	name           string
	url            string
	requestPayload QueryLawListRequestPayload
}

var lawListCategorys = []LawListCategory{
	{
		name: "法律法规",
		url:  "https://law.chemicalsafety.org.cn/laws/100",
		requestPayload: QueryLawListRequestPayload{
			NONCE:         "9gj7s5olrnh",
			SIGN:          "13953322E474A26679B4CC5BECC86C6893C477ECCB55C3C07E16C536EAF29F53",
			SIGN_TYPE:     "SHA256",
			TIMESTAMP:     "20250829160907",
			Current:       1,
			FileType:      0,
			IsLikeSearch:  0,
			NeedHighLight: true,
			Orders: []QueryLawListRequestPayloadOrder{
				{
					Column: "exeDate",
					Asc:    false,
				},
			},
			SearchType: 0,
			Size:       9999,
		},
	},
	{
		name: "国家标准",
		url:  "https://law.chemicalsafety.org.cn/laws/200100",
		requestPayload: QueryLawListRequestPayload{
			NONCE:         "9gj7s5olrnh",
			SIGN:          "13953322E474A26679B4CC5BECC86C6893C477ECCB55C3C07E16C536EAF29F53",
			SIGN_TYPE:     "SHA256",
			TIMESTAMP:     "20250829160907",
			Current:       1,
			FileType:      1,
			IsLikeSearch:  0,
			LevelCodes:    []string{"200100"},
			NeedHighLight: true,
			Orders: []QueryLawListRequestPayloadOrder{
				{
					Column: "exeDate",
					Asc:    false,
				},
			},
			SearchType: 0,
			Size:       9999,
		},
	},
	{
		name: "行业标准",
		url:  "https://law.chemicalsafety.org.cn/laws/200200",
		requestPayload: QueryLawListRequestPayload{
			NONCE:         "9gj7s5olrnh",
			SIGN:          "13953322E474A26679B4CC5BECC86C6893C477ECCB55C3C07E16C536EAF29F53",
			SIGN_TYPE:     "SHA256",
			TIMESTAMP:     "20250829160907",
			Current:       1,
			FileType:      1,
			IsLikeSearch:  0,
			LevelCodes:    []string{"200200"},
			NeedHighLight: true,
			Orders: []QueryLawListRequestPayloadOrder{
				{
					Column: "exeDate",
					Asc:    false,
				},
			},
			SearchType: 0,
			Size:       9999,
		},
	},
	{
		name: "地方标准",
		url:  "https://law.chemicalsafety.org.cn/laws/200300",
		requestPayload: QueryLawListRequestPayload{
			NONCE:         "9gj7s5olrnh",
			SIGN:          "13953322E474A26679B4CC5BECC86C6893C477ECCB55C3C07E16C536EAF29F53",
			SIGN_TYPE:     "SHA256",
			TIMESTAMP:     "20250829160907",
			Current:       1,
			FileType:      1,
			IsLikeSearch:  0,
			LevelCodes:    []string{"200300"},
			NeedHighLight: true,
			Orders: []QueryLawListRequestPayloadOrder{
				{
					Column: "exeDate",
					Asc:    false,
				},
			},
			SearchType: 0,
			Size:       9999,
		},
	},
	{
		name: "团体标准",
		url:  "https://law.chemicalsafety.org.cn/laws/200400",
		requestPayload: QueryLawListRequestPayload{
			NONCE:         "9gj7s5olrnh",
			SIGN:          "13953322E474A26679B4CC5BECC86C6893C477ECCB55C3C07E16C536EAF29F53",
			SIGN_TYPE:     "SHA256",
			TIMESTAMP:     "20250829160907",
			Current:       1,
			FileType:      1,
			IsLikeSearch:  0,
			LevelCodes:    []string{"200400"},
			NeedHighLight: true,
			Orders: []QueryLawListRequestPayloadOrder{
				{
					Column: "exeDate",
					Asc:    false,
				},
			},
			SearchType: 0,
			Size:       9999,
		},
	},
}

// 下载化学品安全法规标准文档
// @Title 下载化学品安全法规标准文档
// @Description https://law.chemicalsafety.org.cn/，下载化学品安全法规标准文档
func main() {
	for _, lawListCategory := range lawListCategorys {
		pageListUrl := "https://law.chemicalsafety.org.cn/api/lan/file/search"
		queryLawListRequestPayload := lawListCategory.requestPayload
		queryLawListResponseDataRecords, err := QueryLawList(pageListUrl, lawListCategory.url, queryLawListRequestPayload)
		if err != nil {
			fmt.Println(err)
			break
		}
		for id_index, law := range queryLawListResponseDataRecords {
			fmt.Println("=========开始处理数据id_index==="+strconv.Itoa(id_index)+" =====catrgory_name = ", lawListCategory.name, "==============")

			code := law.FileNo
			code = strings.ReplaceAll(code, "/", "-")
			code = strings.ReplaceAll(code, "—", "-")
			fmt.Println(code)

			title := law.Title
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "--", "-")
			fmt.Println(title)

			filePath := "../law.chemicalsafety.org.cn/" + title + "(" + code + ")" + ".pdf"
			filePath = strings.ReplaceAll(filePath, "()", "")
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载========")

			lawBaseInfoUrl := fmt.Sprintf("https://law.chemicalsafety.org.cn/api/lan/file/baseInfo/%s", law.FileId)
			fmt.Println(lawBaseInfoUrl)

			lawBaseInfoRefererUrl := fmt.Sprintf("https://law.chemicalsafety.org.cn/law/info/%s", law.FileId)
			queryLawBaseInfoResponseData, err := QueryLawBaseInfoUrl(lawBaseInfoUrl, lawBaseInfoRefererUrl)
			if err != nil {
				fmt.Println(err)
				break
			}

			downloadUrl := queryLawBaseInfoResponseData.PdfUrl
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")
			err = downloadLaw(downloadUrl, lawBaseInfoUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "../law.chemicalsafety.org.cn", "../upload.doc88.com/hbba.sacinfo.org.cn")
			err = copyLawFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			//DownLoadLawTimeSleep := 10
			DownLoadLawTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadLawTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("catrgory_name="+lawListCategory.name+",filePath="+filePath+"===========下载成功 暂停", DownLoadLawTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadLawCategoryTimeSleep := 10
		// DownLoadLawCategoryTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadLawCategoryTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("catrgory_name="+lawListCategory.name+"=========== 暂停", DownLoadLawCategoryTimeSleep, "秒 倒计时", i, "秒===========")
		}
	}
}

type QueryLawListResponse struct {
	Code int                      `json:"code"`
	Data QueryLawListResponseData `json:"data"`
	Ok   bool                     `json:"ok"`
}

type QueryLawListResponseData struct {
	Current          int                               `json:"current"`
	OptimizeCountSql bool                              `json:"optimizeCountSql"`
	Pages            int                               `json:"pages"`
	Records          []QueryLawListResponseDataRecords `json:"records"`
	SearchCount      bool                              `json:"searchCount"`
	Size             int                               `json:"size"`
	Total            int                               `json:"total"`
}

type QueryLawListResponseDataRecords struct {
	FileId string `json:"fileId"`
	FileNo string `json:"fileNo"`
	Title  string `json:"title"`
}

func QueryLawList(requestUrl string, referer string, queryLawListRequestPayload QueryLawListRequestPayload) (queryLawListResponseDataRecords []QueryLawListResponseDataRecords, err error) {
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
	if LawEnableHttpProxy {
		client = LawSetHttpProxy()
	}
	// 将数据编码为JSON格式
	queryLawListRequestPayloadJson, err := json.Marshal(queryLawListRequestPayload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 创建字符串读取器，这是http.Request需要的类型
	body := bytes.NewReader(queryLawListRequestPayloadJson)
	req, err := http.NewRequest("POST", requestUrl, body) //建立连接

	queryLawListResponse := QueryLawListResponse{}
	if err != nil {
		return queryLawListResponseDataRecords, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Host", "law.chemicalsafety.org.cn")
	req.Header.Set("Origin", "https://law.chemicalsafety.org.cn")
	req.Header.Set("Referer", referer)
	req.Header.Set("Tenant-Id", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryLawListResponseDataRecords, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryLawListResponseDataRecords, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryLawListResponseDataRecords, err
	}
	err = json.Unmarshal(respBytes, &queryLawListResponse)
	if err != nil {
		return queryLawListResponseDataRecords, err
	}
	queryLawListResponseDataRecords = queryLawListResponse.Data.Records
	return queryLawListResponseDataRecords, nil
}

type QueryLawBaseInfoResponse struct {
	Code int                          `json:"code"`
	Data QueryLawBaseInfoResponseData `json:"data"`
	Ok   bool                         `json:"ok"`
}

type QueryLawBaseInfoResponseData struct {
	FileId string `json:"fileId"`
	FileNo string `json:"fileNo"`
	Title  string `json:"title"`
	PdfUrl string `json:"pdfUrl"`
}

func QueryLawBaseInfoUrl(requestUrl string, referer string) (queryLawBaseInfoResponseData QueryLawBaseInfoResponseData, err error) {
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
	if LawEnableHttpProxy {
		client = LawSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	queryLawBaseInfoResponse := QueryLawBaseInfoResponse{}
	if err != nil {
		return queryLawBaseInfoResponseData, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Authorization", "test1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "law.chemicalsafety.org.cn")
	req.Header.Set("Origin", "https://law.chemicalsafety.org.cn")
	req.Header.Set("Referer", referer)
	req.Header.Set("Tenant-Id", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryLawBaseInfoResponseData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryLawBaseInfoResponseData, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryLawBaseInfoResponseData, err
	}
	err = json.Unmarshal(respBytes, &queryLawBaseInfoResponse)
	if err != nil {
		return queryLawBaseInfoResponseData, err
	}
	queryLawBaseInfoResponseData = queryLawBaseInfoResponse.Data
	return queryLawBaseInfoResponseData, nil
}

func downloadLaw(attachmentUrl string, referer string, filePath string) error {
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
	if LawEnableHttpProxy {
		client = LawSetHttpProxy()
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
	req.Header.Set("Host", "law.chemicalsafety.org.cn")
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

func copyLawFile(src, dst string) (err error) {
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
