package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	// 	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/otiai10/gosseract/v2"
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

type HdBaResponseData struct {
	Current     int                       `json:"current"`
	Pages       int                       `json:"pages"`
	Records     []HdBaResponseDataRecords `json:"records"`
	SearchCount bool                      `json:"searchCount"`
	Size        int                       `json:"size"`
	Total       int                       `json:"total"`
}
type HdBaResponseDataRecords struct {
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

type HdBaResponseValidateCaptcha struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

const DbBaCookie = "HMACCOUNT=487EF362690A1D5D; Hm_lvt_36f2f0446e1c2cda8410befc24743a9b=1780896225; Hm_lpvt_36f2f0446e1c2cda8410befc24743a9b=1782603684; JSESSIONID=AC388871DFDCAE74D574B58558A8889E"

// ychEduSpider 获取改版-改版-地方标准文档
// @Title 获取改版-地方标准文档
// @Description https://dbba.sacinfo.org.cn/，获取改版-地方标准文档
func main() {
	requestUrl := "https://dbba.sacinfo.org.cn/stdQueryList"
	// 	5699
	current := 1
	maxCurrent := 867
	size := 100
	status := "现行"
	isPageListGo := true
	for isPageListGo {
		if current > maxCurrent {
			isPageListGo = false
			break
		}
		responseData, err := DbBaGetStdQueryList(requestUrl, current, size, status)
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

					//industry := strings.TrimSpace(records.Industry)

					code := strings.ReplaceAll(records.Code, "/", "-")
					code = strings.ReplaceAll(code, "\n", "")

					fileName := chName + "(" + code + ")"
					fmt.Println(fileName)

					filePath := "../dbba.sacinfo.org.cn/" + fileName + ".pdf"
					if _, err := os.Stat(filePath); err != nil {

						stdDetailUrl := fmt.Sprintf("https://dbba.sacinfo.org.cn/stdDetail/%s", records.Pk)
						stdOnlineUrl := fmt.Sprintf("https://dbba.sacinfo.org.cn/portal/online/%s", records.Pk)
						stdOnlinePathUrl := fmt.Sprintf("/portal/online/%s", records.Pk)
						stdOnlineDoc, err := DbBaGetStdQueryHtml(stdOnlineUrl, stdDetailUrl, stdOnlinePathUrl)
						if err != nil {
							fmt.Println(err)
							continue
						}
						// 是否有下载按钮
						downloadButtonNode := htmlquery.FindOne(stdOnlineDoc, `//div[@class="container main-body"]/div[@class="row"]/div[@class="col-sm-12"]/div/div[3]/a[@class="btn btn-warning"]`)
						if downloadButtonNode == nil {
							fmt.Println("没有下载按钮跳过")
							continue
						}
						downLoadUrl := fmt.Sprintf("https://dbba.sacinfo.org.cn/portal/download/%s", records.Pk)
						fmt.Println(downLoadUrl)

						fmt.Println("=======开始下载" + strconv.Itoa(current) + "========")
						err = downloadDbBa(downLoadUrl, stdOnlineUrl, filePath)
						if err != nil {
							fmt.Println(err)
							continue
						}

						// 查看文件大小，如果是空文件，则删除
						fileInfo, err := os.Stat(filePath)
						if err == nil && fileInfo.Size() == 0 {
							fmt.Println("空文件删除")
							err = os.Remove(filePath)
						}
						if err != nil {
							continue
						}

						//复制文件
						tempFilePath := strings.ReplaceAll(filePath, "dbba.sacinfo.org.cn", "temp-dbba.sacinfo.org.cn")
						err = DbBaCopyFile(filePath, tempFilePath)
						if err != nil {
							fmt.Println(err)
							continue
						}
						fmt.Println("=======下载完成========")

						// downloadDbBaPdfSleep := rand.Intn(5)
						downloadDbBaPdfSleep := 10
						for i := 1; i <= downloadDbBaPdfSleep; i++ {
							time.Sleep(time.Second)
							fmt.Println("page="+strconv.Itoa(current)+"=======chName=", chName, "成功，====== 暂停", downloadDbBaPdfSleep, "秒，倒计时", i, "秒===========")
						}
					}
				}
			}

			DownLoadDbBaPageTimeSleep := 10
			// DownLoadDbBaPageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadDbBaPageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(current)+"====== 暂停", DownLoadDbBaPageTimeSleep, "秒 倒计时", i, "秒===========")
			}
			current++
			if current > maxCurrent {
				isPageListGo = false
				break
			}
		}
	}
}

func DbBaGetStdQueryList(requestUrl string, current int, size int, status string) (responseData HdBaResponseData, err error) {
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
	if DbBaEnableHttpProxy {
		client = DbBaSetHttpProxy()
	}
	responseData = HdBaResponseData{}
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
	req.Header.Set("Host", "dbba.sacinfo.org.cn")
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

func DbBaGetStdQueryHtml(requestUrl string, referer string, path string) (doc *html.Node, err error) {
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
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}
	req.Header.Set("authority", "dbba.sacinfo.org.cn")
	req.Header.Set("method", "GET")
	req.Header.Set("path", path)
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", DbBaCookie)
	req.Header.Set("Host", "dbba.sacinfo.org.cn")
	req.Header.Set("Origin", "https://dbba.sacinfo.org.cn")
	req.Header.Set("Referer", referer)
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
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
			ResponseHeaderTimeout: time.Second * 30,
		},
	}
	if DbBaEnableHttpProxy {
		client = DbBaSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("authority", "dbba.sacinfo.org.cn")
	req.Header.Set("method", "GET")
	path := strings.Replace(attachmentUrl, "https://dbba.sacinfo.org.cn", "", 1)
	fmt.Println(path)
	req.Header.Set("path", path)
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", DbBaCookie)
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

func DbBaCopyFile(src, dst string) (err error) {
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
