package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
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
	YchEduEnableHttpProxy = false
	YchEduHttpProxyUrl    = "27.42.168.46:55481"
)

func YchEduSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(YchEduHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type AdultEducationCategory struct {
	categoryName string
	categoryUrl  string
}

var adultEducationCategory = []AdultEducationCategory{
	//{categoryName: "医药类资源", categoryUrl: "http://www.ychedu.com/CRJY/YYST/"},
	//{categoryName: "母婴育儿", categoryUrl: "http://www.ychedu.com/CRJY/myye/"},
	{categoryName: "中医中药", categoryUrl: "http://www.ychedu.com/CRJY/zyzy/"},
	{categoryName: "汽车驾驶技巧", categoryUrl: "http://www.ychedu.com/CRJY/qcjsjq/"},
	{categoryName: "中国名著小说100部", categoryUrl: "http://www.ychedu.com/CRJY/zgxs/"},
	{categoryName: "世界名著100名", categoryUrl: "http://www.ychedu.com/CRJY/sjmz/"},
	//{categoryName: "招生就业移民留学", categoryUrl: "http://www.ychedu.com/CRJY/QTST/"},
	{categoryName: "意境最美诗词", categoryUrl: "http://www.ychedu.com/CRJY/shici/"},
	//{categoryName: "财会类试题", categoryUrl: "http://www.ychedu.com/CRJY/CKST/"},
	//{categoryName: "公务员试题", categoryUrl: "http://www.ychedu.com/CRJY/GWYST/"},
	//{categoryName: "职业资格考试题", categoryUrl: "http://www.ychedu.com/CRJY/ZYZG/"},
	//{categoryName: "建筑类试题", categoryUrl: "http://www.ychedu.com/CRJY/GZLST/"},
	//{categoryName: "考研试题", categoryUrl: "http://www.ychedu.com/CRJY/KYST/"},
	//{categoryName: "考博试题", categoryUrl: "http://www.ychedu.com/CRJY/KBST/"},
	//{categoryName: "自考试题", categoryUrl: "http://www.ychedu.com/CRJY/ZKST/"},
	//{categoryName: "英语四级试题", categoryUrl: "http://www.ychedu.com/CRJY/YYSJ/"},
	//{categoryName: "英语六级试题", categoryUrl: "http://www.ychedu.com/CRJY/YYLJ/"},
	//{categoryName: "计算机试题", categoryUrl: "http://www.ychedu.com/CRJY/JSJST/"},
	//{categoryName: "司法考试试题", categoryUrl: "http://www.ychedu.com/CRJY/SFKS/"},
}

// ychEduSpider 获取宜城教育文档
// @Title 获取宜城教育文档
// @Description http://www.ychedu.com/，获取宜城教育文档
func main() {
	for _, category := range adultEducationCategory {
		detail := make(map[int]string)
		page := 0
		isPageGo := true
		for isPageGo {
			func() {
				defer func() {
					if p := recover(); p != nil {
						fmt.Printf("panic recover! p: %#v\n", p)
						detail[page] = "panic"
					}
				}()
				var listUrl = fmt.Sprintf(category.categoryUrl)
				if page != 0 {
					listUrl = fmt.Sprintf(category.categoryUrl+"List_%d.html", page)
				}
				fmt.Println(listUrl)
				listDoc, _ := htmlquery.LoadURL(listUrl)
				divNodes := htmlquery.Find(listDoc, `//div[@class="bk21"]/div[@align="center"][1]/div`)
				if len(divNodes) >= 1 {
					for _, divNode := range divNodes {
						detailUrl := htmlquery.InnerText(htmlquery.FindOne(divNode, `./ul[@id="soft_lb1"]/div/li/a/@href`))
						detailDoc, _ := htmlquery.LoadURL(detailUrl)
						fmt.Println(detailUrl)

						//div[@class="nr10y"]/table/tbody/tr[1]/td/p[1]/strong/font
						title := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="nr10y"]/table/tbody/tr[1]/td/p[1]`))
						//titleArray := strings.Split(title, "-")
						//title = titleArray[0]
						title = strings.ReplaceAll(title, " ", "")
						fmt.Println(title)

						downloadUrl := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="nr10down"]/a/@href`))
						fmt.Println(downloadUrl)

						//downloadUrlArray, err := url.Parse(downloadUrl)
						//softID := downloadUrlArray.Query().Get("SoftID")
						filePath := "../www.ychedu.com/" + category.categoryName + "/"
						fmt.Println(filePath)

						err := downloadYchEdu(downloadUrl, filePath, title)
						if err != nil {
							fmt.Println(err)
						}
					}
					page++
				} else {
					isPageGo = false
					page = 0
				}
			}()
		}
	}
}

func downloadYchEdu(attachmentUrl string, filePath string, title string) error {
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
	if YchEduEnableHttpProxy {
		client = YchEduSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.ychedu.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "http://www.ychedu.com/")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}

	var suffix string
	contentType := resp.Header.Get("Content-Type")
	switch contentType {
	case "application/msword":
		// doc
		suffix = ".doc"
		break
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		// docx
		suffix = ".docx"
		break
		//case "application/vnd.ms-powerpoint":
		//	// ppt
		//	suffix = ".ppt"
		//	break
		//case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
		//	// pptx
		//	suffix = ".pptx"
		break
	default:
		return nil
	}
	// 创建一个文件用于保存
	fileDiv := filepath.Dir(filePath)
	if _, err = os.Stat(fileDiv); err != nil {
		if os.MkdirAll(fileDiv, 0777) != nil {
			return err
		}
	}
	out, err := os.Create(filePath + title + suffix)
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
