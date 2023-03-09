package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	ZkTiKuEnableHttpProxy = false
	ZkTiKuHttpProxyUrl    = "111.225.152.186:8089"
)

func ZkTiKuSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(ZkTiKuHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type Subject struct {
	name string
	url  string
}

var AllSubject = []Subject{
	{
		name: "语文",
		url:  "https://bilianku.com/document/list?subjectId=1550777663498985474",
	},
	{
		name: "数学",
		url:  "https://bilianku.com/document/list?subjectId=1550777680674660353",
	},
	{
		name: "英语",
		url:  "https://bilianku.com/document/list?subjectId=1550777696118087682",
	},
	{
		name: "物理",
		url:  "https://bilianku.com/document/list?subjectId=1550777712421347330",
	},
	{
		name: "化学",
		url:  "https://bilianku.com/document/list?subjectId=1550777727697002498",
	},
	{
		name: "生物学",
		url:  "https://bilianku.com/document/list?subjectId=1550777745950613506",
	},
	{
		name: "历史",
		url:  "https://bilianku.com/document/list?subjectId=1550777926037250049",
	},
	{
		name: "思想政治",
		url:  "https://bilianku.com/document/list?subjectId=1550777827135561729",
	},
	{
		name: "地理",
		url:  "https://bilianku.com/document/list?subjectId=1550777951320514561",
	},
}

// ychEduSpider 获取名校教研文档
// @Title 获取名校教研文档
// @Description https://bilianku.com/，获取名校教研文档
func main() {
	for _, subject := range AllSubject {
		page := 1
		indexSubjectDoc, err := getZkTiKu(subject.url)
		if err != nil {
			fmt.Println(err)
			break
		}
		indexSubjectPagesNodes := htmlquery.Find(indexSubjectDoc, `//div[@class="kemu-c"]/div[@class="kemu-c-list"]/div[@class="layui-box layui-laypage layui-laypage-default"]/a`)

		var maxPageIndex = 0
		if len(indexSubjectPagesNodes) >= 3 {
			maxPageIndex, _ = strconv.Atoi(htmlquery.InnerText(indexSubjectPagesNodes[len(indexSubjectPagesNodes)-2]))
		}

		isPageListGo := true
		for isPageListGo {
			// 科目最后一页，停止
			if page > maxPageIndex {
				break
			}

			pageListUrl := fmt.Sprintf(subject.url+"&pageIndex=%d", page)
			pageListDoc, err := getZkTiKu(pageListUrl)
			if err != nil {
				fmt.Println(err)
				break
			}
			divNodes := htmlquery.Find(pageListDoc, `//div[@class="kemu-c"]/div[@class="kemu-c-list"]/div[@class="kemu-c-item"]`)
			if len(divNodes) >= 1 {
				for _, listNode := range divNodes {

					fmt.Println("=================================================================================")
					fmt.Println(pageListUrl)

					detailUrl := "https://bilianku.com" + htmlquery.InnerText(htmlquery.FindOne(listNode, `./a/@href`))
					detailDoc, err := getZkTiKu(detailUrl)
					if err != nil {
						fmt.Println(err)
						continue
					}

					// 下载文件列表
					fileNodes := htmlquery.Find(detailDoc, `//div[@class="kemu-info-list"]/div[@class="kemu-info-item"]`)
					if len(fileNodes) >= 1 {
						for _, fileNode := range fileNodes {

							// 文件类型
							suffix := ""
							imgSrc := htmlquery.InnerText(htmlquery.FindOne(fileNode, `./div[@class="kemu-info-item-l"]/img/@src`))
							if strings.Index(imgSrc, "pdf") > -1 {
								suffix = ".pdf"
							}
							if strings.Index(imgSrc, "docx") > -1 {
								suffix = ".docx"
							}

							fileTile := htmlquery.InnerText(htmlquery.FindOne(fileNode, `./div[@class="kemu-info-item-c"]/div[@class="kemu-info-item-c-t"]`))
							fileTile = strings.ReplaceAll(fileTile, "/", "-")
							fileTile = strings.ReplaceAll(fileTile, " ", "")
							fileTile = strings.ReplaceAll(fileTile, "\n", "")
							fileTile = strings.ReplaceAll(fileTile, "\r", "")
							fmt.Println(fileTile)

							onClickText := htmlquery.InnerText(htmlquery.FindOne(fileNode, `./div[@class="kemu-info-item-r"]/a[@class="kemu-info-item-r-i bgc3c8c93"]/@onclick`))
							reg := regexp.MustCompile(`[0-9]+`)
							fileDetailId, _ := strconv.Atoi(reg.FindAllString(onClickText, -1)[0])

							downloadUrl := "https://bilianku.com/document/download"
							fmt.Println(downloadUrl)
							err = downloadZkTiKu(downloadUrl, detailUrl, fileDetailId)
							if err != nil {
								fmt.Println(err)
								continue
							}

							downloadExecUrl := fmt.Sprintf("https://bilianku.com/document/downloadexec?id=%d", fileDetailId)
							fmt.Println(downloadExecUrl)
							filePath := "../bilianku.com/" + subject.name + "/"
							fileName := fileTile + suffix
							err = downloadExecZkTiKu(downloadExecUrl, detailUrl, filePath, fileName)
							if err != nil {
								fmt.Println(err)
								continue
							}
							time.Sleep(time.Second * 3)
						}
					}
					time.Sleep(time.Second * 3)
				}
				page++
			} else {
				isPageListGo = false
				page = 1
				break
			}
		}
	}
}

func getZkTiKu(url string) (doc *html.Node, err error) {
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
	if ZkTiKuEnableHttpProxy {
		client = ZkTiKuSetHttpProxy()
	}
	req, err := http.NewRequest("GET", url, nil) //建立连接
	if err != nil {
		return doc, err
	}
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return doc, err
	}
	defer resp.Body.Close()
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

type downloadResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func downloadZkTiKu(downloadUrl string, referer string, fileDetailId int) error {
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
	if ZkTiKuEnableHttpProxy {
		client = ZkTiKuSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("id", strconv.Itoa(fileDetailId))
	req, err := http.NewRequest("POST", downloadUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "__51vcke__JohyjbD1C0xg9qUz=0b89b92a-6862-5964-950c-fd0424460459; __51vuft__JohyjbD1C0xg9qUz=1676033539876; mobile=15238369929; JSESSIONID=2619B18AD9F7D848ED3D0ED88206F936; __51uvsct__JohyjbD1C0xg9qUz=3; token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1SWQiOjE2MjM5ODczNDYxMDM3MDk2OTgsInVUeXBlIjoyLCJleHAiOjE2NzY2NTcyNDR9.yB-IPSls4xcQTdNKqG43q1d9hz8w3FxIYWWw0xH0shM; __vtins__JohyjbD1C0xg9qUz=%7B%22sid%22%3A%20%22647da08c-62a8-599a-995e-4a351a7137c1%22%2C%20%22vd%22%3A%207%2C%20%22stt%22%3A%20395895%2C%20%22dr%22%3A%209501%2C%20%22expires%22%3A%201676054244939%2C%20%22ct%22%3A%201676052444939%7D")
	req.Header.Set("Host", "bilianku.com")
	req.Header.Set("Origin", "https://bilianku.com")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("token", "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1SWQiOjE2MjM5ODczNDYxMDM3MDk2OTgsInVUeXBlIjoyLCJleHAiOjE2NzY2NTk2MTJ9.QkOLuwo8vwfNne0ClPbDYgtqZphDr4vuuFM9Ny9yb0I")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	downloadResult := &downloadResult{}
	err = json.Unmarshal(respBytes, downloadResult)
	if err != nil {
		fmt.Println(111)
		return err
	}
	if downloadResult.Code != 1 {
		return errors.New(downloadResult.Msg)
	}
	return nil
}

func downloadExecZkTiKu(attachmentUrl string, referer string, filePath string, fileName string) error {
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
	if ZkTiKuEnableHttpProxy {
		client = ZkTiKuSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "__51vcke__JohyjbD1C0xg9qUz=0b89b92a-6862-5964-950c-fd0424460459; __51vuft__JohyjbD1C0xg9qUz=1676033539876; mobile=15238369929; JSESSIONID=2619B18AD9F7D848ED3D0ED88206F936; __51uvsct__JohyjbD1C0xg9qUz=3; token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1SWQiOjE2MjM5ODczNDYxMDM3MDk2OTgsInVUeXBlIjoyLCJleHAiOjE2NzY2NTcyNDR9.yB-IPSls4xcQTdNKqG43q1d9hz8w3FxIYWWw0xH0shM; __vtins__JohyjbD1C0xg9qUz=%7B%22sid%22%3A%20%22647da08c-62a8-599a-995e-4a351a7137c1%22%2C%20%22vd%22%3A%207%2C%20%22stt%22%3A%20395895%2C%20%22dr%22%3A%209501%2C%20%22expires%22%3A%201676054244939%2C%20%22ct%22%3A%201676052444939%7D")
	req.Header.Set("Host", "bilianku.com")
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
	out, err := os.Create(filePath + fileName)
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
