package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	DbBaEnableHttpProxy = false
	DbBaHttpProxyUrl    = "111.225.152.186:8089"
)

func DbBaSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(DbBaHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type ResponseData struct {
	Current     int                   `json:"current"`
	Pages       int                   `json:"pages"`
	Records     []ResponseDataRecords `json:"records"`
	SearchCount bool                  `json:"searchCount"`
	Size        int                   `json:"size"`
	Total       int                   `json:"total"`
}
type ResponseDataRecords struct {
	ActDate    int    `json:"actDate"`
	ChName     string `json:"chName"`
	ChargeDept string `json:"chargeDept"`
	Code       string `json:"code"`
	Empty      bool   `json:"empty"`
	Industry   string `json:"industry"`
	IssueDate  int    `json:"issueDate"`
	Pk         string `json:"pk"`
	RecordDate int    `json:"recordDate"`
	RecordNo   string `json:"recordNo"`
	Status     string `json:"status"`
}

// ychEduSpider 获取地方标准文档
// @Title 获取地方标准文档
// @Description https://dbba.sacinfo.org.cn/，获取地方标准文档
func main() {
	requestUrl := "https://dbba.sacinfo.org.cn/stdQueryList"
	current := 1
	size := 50
	status := "现行"
	isPageListGo := true
	for isPageListGo {
		responseData, err := GetStdQueryList(requestUrl, current, size, status)
		if err != nil {
			fmt.Println(err)
			break
		}
		if len(responseData.Records) > 0 {
			for _, records := range responseData.Records {
				if records.Empty == false {
					fmt.Println("=======当前页为：" + strconv.Itoa(current) + "========")
					chName := strings.ReplaceAll(records.ChName, " ", "")
					chName = strings.ReplaceAll(chName, "/", "-")
					chName = strings.ReplaceAll(chName, "\n", "")
					chName = strings.ReplaceAll(chName, ":", "-")
					chName = strings.ReplaceAll(chName, "：", "-")

					code := strings.ReplaceAll(records.Code, "/", "-")
					code = strings.ReplaceAll(code, "\n", "")

					fileName := chName + "(" + code + ")"
					fmt.Println(fileName)

					downLoadUrl := fmt.Sprintf("https://dbba.sacinfo.org.cn/attachment/downloadStdFile?pk=%s", records.Pk)
					fmt.Println(downLoadUrl)

					detailUrl := fmt.Sprintf("https://dbba.sacinfo.org.cn/stdDetail/%s", records.Pk)
					fmt.Println(detailUrl)

					filePath := "../dbba.sacinfo.org.cn/" + fileName + ".pdf"
					if _, err := os.Stat(filePath); err != nil {
						fmt.Println("=======开始下载" + strconv.Itoa(current) + "========")
						err = downloadDbBa(downLoadUrl, detailUrl, filePath)
						if err != nil {
							fmt.Println(err)
							continue
						}
						fmt.Println("=======开始完成========")
					}

					// 查看文件大小，如果是空文件，则删除
					fi, err := os.Stat(filePath)
					if err == nil && fi.Size() == 0 {
						err := os.Remove(filePath)
						if err != nil {
							continue
						}
					}

					time.Sleep(time.Millisecond * 100)
				}
			}

			if current < responseData.Pages {
				current++
			} else {
				isPageListGo = false
				current = 1
				break
			}
		} else {
			isPageListGo = false
			current = 1
			break
		}
	}
}

func GetStdQueryList(requestUrl string, current int, size int, status string) (responseData ResponseData, err error) {
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
	if DbBaEnableHttpProxy {
		client = DbBaSetHttpProxy()
	}
	responseData = ResponseData{}
	postData := url.Values{}
	postData.Add("current", strconv.Itoa(current))
	postData.Add("size", strconv.Itoa(size))
	postData.Add("status", status)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return responseData, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "__51vcke__JohyjbD1C0xg9qUz=0b89b92a-6862-5964-950c-fd0424460459; __51vuft__JohyjbD1C0xg9qUz=1676033539876; mobile=15238369929; JSESSIONID=2619B18AD9F7D848ED3D0ED88206F936; __51uvsct__JohyjbD1C0xg9qUz=3; token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1SWQiOjE2MjM5ODczNDYxMDM3MDk2OTgsInVUeXBlIjoyLCJleHAiOjE2NzY2NTcyNDR9.yB-IPSls4xcQTdNKqG43q1d9hz8w3FxIYWWw0xH0shM; __vtins__JohyjbD1C0xg9qUz=%7B%22sid%22%3A%20%22647da08c-62a8-599a-995e-4a351a7137c1%22%2C%20%22vd%22%3A%207%2C%20%22stt%22%3A%20395895%2C%20%22dr%22%3A%209501%2C%20%22expires%22%3A%201676054244939%2C%20%22ct%22%3A%201676052444939%7D")
	req.Header.Set("Host", "bilianku.com")
	req.Header.Set("Origin", "https://dbba.sacinfo.org.cn")
	req.Header.Set("Referer", "https://dbba.sacinfo.org.cn/stdList")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
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
		return responseData, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseData, err
	}
	err = json.Unmarshal(respBytes, &responseData)
	if err != nil {
		return responseData, err
	}
	return responseData, nil
}

func downloadDbBa(attachmentUrl string, referer string, filePath string) error {
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
	if DbBaEnableHttpProxy {
		client = DbBaSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "dbba.sacinfo.org.cn")
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
