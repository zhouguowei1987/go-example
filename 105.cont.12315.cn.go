package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"io/ioutil"
	"mime"
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
	ContEnableHttpProxy = false
	ContHttpProxyUrl    = "111.225.152.186:8089"
)

func ContSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(ContHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type ContResponse struct {
	Data      []ContResponseData `json:"Data"`
	Page      int                `json:"Page"`
	Total     int                `json:"Total"`
	TotalPage int                `json:"TotalPage"`
}
type ContResponseData struct {
	Id    string `json:"Id"`
	Title string `json:"Title"`
}

// ychEduSpider 获取部委合同示范文本文档
// @Title 获取部委合同示范文本文档
// @Description https://cont.12315.cn/，获取部委合同示范文本文档
func main() {
	current := 1
	isPageListGo := true
	for isPageListGo {
		requestUrl := fmt.Sprintf("https://cont.12315.cn/api/content/SearchTemplates?loc=false&p=%d&key=", current)
		contResponse, err := GetSearchTemplates(requestUrl)
		if err != nil {
			fmt.Println(err)
			break
		}
		if current > contResponse.TotalPage {
			isPageListGo = false
			break
		}
		if len(contResponse.Data) > 0 {
			for _, data := range contResponse.Data {
				fmt.Println("=======当前页为：" + strconv.Itoa(current) + "========")
				id := data.Id

				viewUrl := fmt.Sprintf("https://cont.12315.cn/View?id=%s", id)
				fmt.Println(viewUrl)

				viewDoc, err := htmlquery.LoadURL(viewUrl)
				if err != nil {
					fmt.Println(err)
					break
				}

				releaseNumberNode := htmlquery.FindOne(viewDoc, `//div[@class="samr-view-info"]/div[3]/div[@class="info-content"]`)
				if releaseNumberNode == nil {
					continue
				}
				releaseNumber := htmlquery.InnerText(releaseNumberNode)
				releaseNumber = strings.ReplaceAll(releaseNumber, "/", "-")
				releaseNumber = strings.ReplaceAll(releaseNumber, "\n", "")

				title := strings.ReplaceAll(data.Title, " ", "")
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, "\n", "")
				title = strings.ReplaceAll(title, ":", "-")
				title = strings.ReplaceAll(title, "：", "-")

				fileName := title + "(" + releaseNumber + ")"
				fmt.Println(fileName)

				downLoadUrl := fmt.Sprintf("https://cont.12315.cn/api/File/DownTemplate?id=%s&type=1", id)
				fmt.Println(downLoadUrl)

				filePath := "../cont.12315.cn/" + fileName + ".docx"
				if _, err := os.Stat(filePath); err != nil {
					fmt.Println("=======开始下载" + strconv.Itoa(current) + "========")
					err = downloadCont(downLoadUrl, viewUrl, filePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("=======开始完成========")
					time.Sleep(time.Millisecond * 200)
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

			if current < contResponse.TotalPage {
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

func GetSearchTemplates(requestUrl string) (contResponse ContResponse, err error) {
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
	if ContEnableHttpProxy {
		client = ContSetHttpProxy()
	}
	contResponse = ContResponse{}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return contResponse, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "f2b414daf483410388b9bd9460d8c62d=WyIzODQzNTE0OTgxIl0")
	req.Header.Set("Host", "cont.12315.cn")
	req.Header.Set("Origin", "https://cont.12315.cn")
	req.Header.Set("Referer", "https://cont.12315.cn/National")
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
		return contResponse, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return contResponse, err
	}
	err = json.Unmarshal(respBytes, &contResponse)
	if err != nil {
		return contResponse, err
	}
	return contResponse, nil
}

func downloadCont(attachmentUrl string, referer string, filePath string) error {
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
	if ContEnableHttpProxy {
		client = ContSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "cont.12315.cn")
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
	// 从响应头获取文件名
	fileName := getFilenameFromHeader(resp.Header)
	if fileName != "" {
		fileExt := strings.Split(fileName, ".")
		filePath = strings.ReplaceAll(filePath, "docx", fileExt[1])
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

// 从响应头获取文件名
func getFilenameFromHeader(header http.Header) string {
	contentDisposition := header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		// 如果获取不到文件名，则使用默认文件名
		return ""
	}
	return params["filename"]
}
