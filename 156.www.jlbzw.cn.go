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

var JlBzwEnableHttpProxy = false
var JlBzwHttpProxyUrl = "111.225.152.186:8089"
var JlBzwHttpProxyUrlArr = make([]string, 0)

func JlBzwHttpProxy() error {
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
					JlBzwHttpProxyUrlArr = append(JlBzwHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					JlBzwHttpProxyUrlArr = append(JlBzwHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func JlBzwSetHttpProxy() (httpclient *http.Client) {
	if JlBzwHttpProxyUrl == "" {
		if len(JlBzwHttpProxyUrlArr) <= 0 {
			err := JlBzwHttpProxy()
			if err != nil {
				JlBzwSetHttpProxy()
			}
		}
		JlBzwHttpProxyUrl = JlBzwHttpProxyUrlArr[0]
		if len(JlBzwHttpProxyUrlArr) >= 2 {
			JlBzwHttpProxyUrlArr = JlBzwHttpProxyUrlArr[1:]
		} else {
			JlBzwHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(JlBzwHttpProxyUrl)
	ProxyURL, _ := url.Parse(JlBzwHttpProxyUrl)
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

type QueryJlBzwListRequestPayload struct {
	A000       string `json:"a000"`
	EndYear    string `json:"endYear"`
	Mode       string `json:"mode"`
	OrderBy    string `json:"orderBy"`
	PageNo     int    `json:"pageNo"`
	PageSize   int    `json:"pageSize"`
	SearchText string `json:"searchText"`
	StartYear  string `json:"startYear"`
}

type QueryJlBzwDownloadUrlFormData struct {
	Cond_a104       string `json:"cond.a104"`
	Cond_filename_s string `json:"cond.filename_s"`
}

// 下载吉林省地方标准文档
// @Title 下载吉林省地方标准文档
// @Description http://www.jlbzw.cn/，下载吉林省地方标准文档
func main() {
	pageListUrl := "http://www.jlbzw.cn/admin-api/portal/landmarksearch/selectStandardReadInfo"
	fmt.Println(pageListUrl)
	page := 1
	maxPage := 259
	rows := 10
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		queryJlBzwListRequestPayload := QueryJlBzwListRequestPayload{
			A000:       "",
			EndYear:    "",
			Mode:       "jilindibiao",
			OrderBy:    "",
			PageNo:     page,
			PageSize:   rows,
			SearchText: "",
			StartYear:  "",
		}
		queryJlBzwListResponseDataStandardInfos, err := QueryJlBzwList(pageListUrl, queryJlBzwListRequestPayload)
		if err != nil {
			fmt.Println(err)
			break
		}
		for _, jlBzw := range queryJlBzwListResponseDataStandardInfos {
			fmt.Println("=====================开始处理数据 page = ", page, "=========================")

			if jlBzw.Filename == "null" {
				fmt.Println("文档没有附件地址，跳过")
				continue
			}

			code := jlBzw.A100
			code = strings.ReplaceAll(code, "/", "-")
			code = strings.ReplaceAll(code, "—", "-")
			fmt.Println(code)

			title := jlBzw.A298
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "--", "-")
			fmt.Println(title)

			filePath := "../www.jlbzw.cn/" + title + "(" + code + ")" + ".pdf"
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载========")

			// 获取下载地址
			queryJlBzwDownloadUrlFormData := QueryJlBzwDownloadUrlFormData{
				Cond_a104:       jlBzw.A104,
				Cond_filename_s: jlBzw.Filename,
			}
			requestJlBzwDownloadUrl := "http://www.jlbzw.cn/admin-api/standard/preview/previewFGPDFByCond"
			queryJlBzwDownloadUrlResponse, err := QueryJlBzwDownloadUrl(requestJlBzwDownloadUrl, queryJlBzwDownloadUrlFormData)
			if err != nil {
				fmt.Println(err)
				break
			}

			downloadUrl := queryJlBzwDownloadUrlResponse.Data
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")

			requestJlBzwDownloadRefererUrl := "http://www.jlbzw.cn/localStandards"

			err = downloadJlBzw(downloadUrl, requestJlBzwDownloadRefererUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "../www.jlbzw.cn", "../upload.doc88.com/dbba.sacinfo.org.cn")
			err = copyJlBzwFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			//DownLoadJlBzwTimeSleep := 10
			DownLoadJlBzwTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadJlBzwTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadJlBzwTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadJlBzwPageTimeSleep := 10
		// DownLoadJlBzwPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadJlBzwPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadJlBzwPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

type QueryJlBzwListResponse struct {
	Code int                        `json:"code"`
	Data QueryJlBzwListResponseData `json:"data"`
	Msg  string                     `json:"msg"`
}

type QueryJlBzwListResponseData struct {
	StandardInfos     []QueryJlBzwListResponseDataStandardInfos `json:"standardInfos"`
	StandardInfos_num int                                       `json:"standardInfos_num"`
	StandardInfos_sum struct{}                                  `json:"standardInfos_sum"`
}

type QueryJlBzwListResponseDataStandardInfos struct {
	A100     string `json:"a100"`
	A104     string `json:"a104"`
	A298     string `json:"a298"`
	Filename string `json:"filename"`
}

func QueryJlBzwList(requestUrl string, queryJlBzwListRequestPayload QueryJlBzwListRequestPayload) (queryJlBzwListResponseDataStandardInfos []QueryJlBzwListResponseDataStandardInfos, err error) {
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
	if JlBzwEnableHttpProxy {
		client = JlBzwSetHttpProxy()
	}
	// 将数据编码为JSON格式
	queryJlBzwListRequestPayloadJson, err := json.Marshal(queryJlBzwListRequestPayload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 创建字符串读取器，这是http.Request需要的类型
	body := bytes.NewReader(queryJlBzwListRequestPayloadJson)
	req, err := http.NewRequest("POST", requestUrl, body) //建立连接

	queryJlBzwListResponse := QueryJlBzwListResponse{}
	if err != nil {
		return queryJlBzwListResponseDataStandardInfos, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Authorization", "test1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Host", "www.jlbzw.cn")
	req.Header.Set("Origin", "http://www.jlbzw.cn")
	req.Header.Set("Referer", "http://www.jlbzw.cn/localStandards")
	req.Header.Set("Tenant-Id", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryJlBzwListResponseDataStandardInfos, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryJlBzwListResponseDataStandardInfos, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryJlBzwListResponseDataStandardInfos, err
	}
	err = json.Unmarshal(respBytes, &queryJlBzwListResponse)
	if err != nil {
		return queryJlBzwListResponseDataStandardInfos, err
	}
	queryJlBzwListResponseDataStandardInfos = queryJlBzwListResponse.Data.StandardInfos
	return queryJlBzwListResponseDataStandardInfos, nil
}

type QueryJlBzwDownloadUrlResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

func QueryJlBzwDownloadUrl(requestUrl string, queryJlBzwDownloadUrlFormData QueryJlBzwDownloadUrlFormData) (queryJlBzwDownloadUrlResponse QueryJlBzwDownloadUrlResponse, err error) {
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
	if JlBzwEnableHttpProxy {
		client = JlBzwSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("cond.a104", queryJlBzwDownloadUrlFormData.Cond_a104)
	postData.Add("cond.filename_s", queryJlBzwDownloadUrlFormData.Cond_filename_s)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Authorization", "test1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "www.jlbzw.cn")
	req.Header.Set("Origin", "http://www.jlbzw.cn")
	req.Header.Set("Referer", "http://www.jlbzw.cn/localStandards")
	req.Header.Set("Tenant-Id", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryJlBzwDownloadUrlResponse, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryJlBzwDownloadUrlResponse, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryJlBzwDownloadUrlResponse, err
	}
	err = json.Unmarshal(respBytes, &queryJlBzwDownloadUrlResponse)
	if err != nil {
		return queryJlBzwDownloadUrlResponse, err
	}
	return queryJlBzwDownloadUrlResponse, nil
}

func downloadJlBzw(attachmentUrl string, referer string, filePath string) error {
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
	if JlBzwEnableHttpProxy {
		client = JlBzwSetHttpProxy()
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
	req.Header.Set("Host", "www.jlbzw.cn")
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

func copyJlBzwFile(src, dst string) (err error) {
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
