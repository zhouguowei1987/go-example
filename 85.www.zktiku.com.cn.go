package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
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
	ZkTiKuEnableHttpProxy = false
	ZkTiKuHttpProxyUrl    = "27.42.168.46:55481"
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
		url:  "https://www.zktiku.com.cn/document/list?subjectId=1550777663498985474",
	},
	{
		name: "数学",
		url:  "https://www.zktiku.com.cn/document/list?subjectId=1550777680674660353",
	},
	{
		name: "英语",
		url:  "https://www.zktiku.com.cn/document/list?subjectId=1550777696118087682",
	},
	{
		name: "物理",
		url:  "https://www.zktiku.com.cn/document/list?subjectId=1550777712421347330",
	},
	{
		name: "化学",
		url:  "https://www.zktiku.com.cn/document/list?subjectId=1550777727697002498",
	},
	{
		name: "生物学",
		url:  "https://www.zktiku.com.cn/document/list?subjectId=1550777745950613506",
	},
	{
		name: "历史",
		url:  "https://www.zktiku.com.cn/document/list?subjectId=1550777926037250049",
	},
	{
		name: "思想政治",
		url:  "https://www.zktiku.com.cn/document/list?subjectId=1550777827135561729",
	},
	{
		name: "地理",
		url:  "https://www.zktiku.com.cn/document/list?subjectId=1550777951320514561",
	},
}

// ychEduSpider 获取名校教研文档
// @Title 获取名校教研文档
// @Description https://www.zktiku.com.cn/，获取名校教研文档
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

					detailUrl := "https://www.zktiku.com.cn" + htmlquery.InnerText(htmlquery.FindOne(listNode, `./a/@href`))
					fmt.Println(detailUrl)
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

							downloadUrl := "https://www.zktiku.com.cn" + htmlquery.InnerText(htmlquery.FindOne(fileNode, `./div[@class="kemu-info-item-r"]/a/@href`))
							downloadUrl = strings.ReplaceAll(downloadUrl, "preview?fileDetailId", "downloadexec?id")
							fmt.Println(downloadUrl)

							filePath := "../www.zktiku.com.cn/" + subject.name + "/"
							fileName := fileTile + suffix
							err := downloadZkTiKu(downloadUrl, filePath, fileName)
							if err != nil {
								fmt.Println(err)
								continue
							}
						}
					}
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
	client := &http.Client{}                     //初始化客户端
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

func downloadZkTiKu(attachmentUrl string, filePath string, fileName string) error {
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
