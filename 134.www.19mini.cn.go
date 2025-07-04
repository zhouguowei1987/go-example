package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var MiNi19EnableHttpProxy = false
var MiNi19HttpProxyUrl = ""
var MiNi19HttpProxyUrlArr = make([]string, 0)

func MiNi19HttpProxy() error {
	pageMax := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// pageMax := []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
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
					MiNi19HttpProxyUrlArr = append(MiNi19HttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					MiNi19HttpProxyUrlArr = append(MiNi19HttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func MiNi19SetHttpProxy() (httpclient *http.Client) {
	if MiNi19HttpProxyUrl == "" {
		if len(MiNi19HttpProxyUrlArr) <= 0 {
			err := MiNi19HttpProxy()
			if err != nil {
				MiNi19SetHttpProxy()
			}
		}
		if len(MiNi19HttpProxyUrlArr) > 1 {
			MiNi19HttpProxyUrl = MiNi19HttpProxyUrlArr[0]
		}
		if len(MiNi19HttpProxyUrlArr) >= 2 {
			MiNi19HttpProxyUrlArr = MiNi19HttpProxyUrlArr[1:]
		} else {
			MiNi19HttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(MiNi19HttpProxyUrl)
	ProxyURL, _ := url.Parse(MiNi19HttpProxyUrl)
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

type MiNi19EducationCategory struct {
	categoryName string
	categoryUrl  string
	classId      int
	startPage    int
	wtspurl      string
}

var miNi19EducationCategory = []MiNi19EducationCategory{
	// {
	// 	categoryName: "教案",
	// 	categoryUrl:  "http://www.19mini.cn/ziyuan/jiaoan/",
	// 	classId:      29,
	// 	startPage:    1,
	// 	wtspurl:      "/ziyuan/jiaoan/",
	// },
	// {
	// 	categoryName: "试卷",
	// 	categoryUrl:  "http://www.19mini.cn/ziyuan/shijuan/",
	// 	classId:      30,
	// 	startPage:    1,
	// 	wtspurl:      "/ziyuan/shijuan/",
	// },
	{
		categoryName: "资料",
		categoryUrl:  "http://www.19mini.cn/ziyuan/ziliao/",
		classId:      33,
		startPage:    1,
		wtspurl:      "/ziyuan/ziliao/",
	},
	{
		categoryName: "训练",
		categoryUrl:  "http://www.19mini.cn/ziyuan/xunlian/",
		classId:      34,
		startPage:    1,
		wtspurl:      "/ziyuan/xunlian/",
	},
}

var MiNi19Cookie = "Hm_lvt_e82ba7292d1c4fbfbf1933dc51f62e60=1747493636,1747717671,1749290636; HMACCOUNT=1CCD0111717619C6; XLA_CI=23bb8bc755f819c1fe15ab77e57ffc56; Hm_lpvt_e82ba7292d1c4fbfbf1933dc51f62e60=1749476471; _wtspurl=wtspurl; _wtsuid=ebaf0b44-8238-4ba6-bb84-5b964a783a70; _wtscpk=26716981c1; _wtsexp=1749477828; _wtsjsk=01cc0ef384da84859bdc1d97d17189f9"

// MiNi19Spider 获取迷你语文网文档
// @Title 获取迷你语文网文档
// @Description http://www.19mini.cn/，获取迷你语文网文档
func main() {
	for _, category := range miNi19EducationCategory {
		page := category.startPage
		isPageGo := true
		MiNi19Cookie = strings.ReplaceAll(MiNi19Cookie, "_wtspurl=wtspurl", "_wtspurl="+category.wtspurl)
		for isPageGo {
			var listUrl = fmt.Sprintf(category.categoryUrl)
			if page != 1 {
				listUrl = strings.ReplaceAll(category.categoryUrl, "index.html", "") + fmt.Sprintf("index_%d.html", page)
			}
			fmt.Println(listUrl)
			listDoc, err := ListMiNi19(listUrl, category.categoryUrl)
			if err != nil {
				fmt.Println(err)
				break
			}
			divNodes := htmlquery.Find(listDoc, `//ul[@class="e2"]/li`)
			if len(divNodes) >= 1 {
				for _, divNode := range divNodes {
					fmt.Println("============================================================================")
					fmt.Println("分页：", page)
					fmt.Println("=======当前页URL", listUrl, "========")

					titleNode := htmlquery.FindOne(divNode, `./a[@class="title"]`)
					if titleNode == nil {
						fmt.Println("标题不存在")
						continue
					}
					title := htmlquery.InnerText(titleNode)
					title = strings.TrimSpace(title)
					title = strings.ReplaceAll(title, "免费", "")
					title = strings.ReplaceAll(title, "-", "")
					title = strings.ReplaceAll(title, " ", "")
					title = strings.ReplaceAll(title, "|", "-")
					fmt.Println(title)
					// 过滤文件名中含有“扫描”字样文件
					if strings.Index(title, "扫描") != -1 {
						fmt.Println("过滤文件名中含有“扫描”字样文件")
						continue
					}
					// 过滤文件名中含有“图片”字样文件
					if strings.Index(title, "图片") != -1 {
						fmt.Println("过滤文件名中含有“图片”字样文件")
						continue
					}

					detailUrl := htmlquery.InnerText(htmlquery.FindOne(divNode, `./a[@class="title"]/@href`))
					detailUrlSplitArray := strings.Split(detailUrl, "/")
					idHtml := detailUrlSplitArray[len(detailUrlSplitArray)-1]
					idStr := strings.ReplaceAll(idHtml, ".html", "")
					id, _ := strconv.Atoi(idStr)

					// 获取文档类型
					MiNi19ViewUrl := fmt.Sprintf(category.categoryUrl+"%d.html", id)
					fmt.Println(MiNi19ViewUrl)
					MiNi19ViewDoc, err := MiNi19ViewDoc(MiNi19ViewUrl, listUrl)

					// 获取文档页详情后，暂停一段时间
					MiNi19ViewDocTimeSleep := rand.Intn(3)
					for i := 1; i <= MiNi19ViewDocTimeSleep; i++ {
						time.Sleep(time.Second)
						fmt.Println("===========获取文档页详情：", title, "成功，暂停", MiNi19ViewDocTimeSleep, "秒，倒计时", i, "秒===========")
					}
					// fmt.Println(htmlquery.InnerText(MiNi19ViewDoc))
					// os.Exit(1)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fileTypeNode := htmlquery.FindOne(MiNi19ViewDoc, `//div[@class="infolist"]/span[5]`)
					if fileTypeNode == nil {
						fmt.Println("文档类型不存在")
						continue
					}
					fileType := htmlquery.InnerText(fileTypeNode)
					fmt.Println(fileType)
					if strings.Index(fileType, "doc") == -1 && strings.Index(fileType, "zip") == -1 {
						fmt.Println("文档类型不是doc或zip文档，跳过")
						continue
					}

					if strings.Index(fileType, "doc") != -1 {
						fileType = ".doc"
					} else if strings.Index(fileType, "zip") != -1 {
						fileType = ".zip"
					}

					filePath := "F:\\workspace\\www.19mini.cn\\www.19mini.cn\\" + title + fileType
					_, err = os.Stat(filePath)
					if err == nil {
						fmt.Println("文档已下载过，跳过")
						continue
					}
					MiNi19DownloadUrl := fmt.Sprintf("http://www.19mini.cn/e/DownSys/DownSoft/?classid=%d&id=%d&pathid=0", category.classId, id)
					fmt.Println(MiNi19DownloadUrl)
					MiNi19DownloadDoc, err := MiNi19DownloadDoc(MiNi19DownloadUrl, detailUrl)

					// 获取下载页详情后，暂停一段时间
					MiNi19DownloadDocTimeSleep := rand.Intn(3)
					for i := 1; i <= MiNi19DownloadDocTimeSleep; i++ {
						time.Sleep(time.Second)
						fmt.Println("===========获取下载页详情：", title, "成功，暂停", MiNi19DownloadDocTimeSleep, "秒，倒计时", i, "秒===========")
					}
					// fmt.Println(htmlquery.InnerText(MiNi19DownloadDoc))
					// os.Exit(1)
					if err != nil {
						fmt.Println(err)
						continue
					}
					// /html/body/div[2]/a
					attachmentNode := htmlquery.FindOne(MiNi19DownloadDoc, `/html/body/div[2]/a/@href`)
					if attachmentNode == nil {
						fmt.Println("没有下载链接，跳过")
						continue
					}
					attachmentUrl := "http://www.19mini.cn/e/DownSys" + strings.ReplaceAll(htmlquery.InnerText(attachmentNode), "..", "")
					fmt.Println(attachmentUrl)

					fmt.Println("=======开始下载========")
					err = downloadMiNi19(attachmentUrl, filePath, MiNi19DownloadUrl)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("=======完成下载========")
					// DownLoadMiNi19TimeSleep := rand.Intn(10)
					DownLoadMiNi19TimeSleep := 10
					for i := 1; i <= DownLoadMiNi19TimeSleep; i++ {
						time.Sleep(time.Second)
						fmt.Println("page="+strconv.Itoa(page)+"===========下载", title, "成功，暂停", DownLoadMiNi19TimeSleep, "秒，倒计时", i, "秒===========")
					}
				}
				page++
			} else {
				isPageGo = false
				page = 1
				break
			}
		}
	}
}

func ListMiNi19(requestUrl string, referer string) (doc *html.Node, err error) {
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
	if MiNi19EnableHttpProxy {
		client = MiNi19SetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}
	fmt.Println(referer)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", MiNi19Cookie)
	req.Header.Set("Referer", referer)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Mobile Safari/537.36")
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

func MiNi19DownloadDoc(requestUrl string, referer string) (doc *html.Node, err error) {
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
	if MiNi19EnableHttpProxy {
		client = MiNi19SetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", MiNi19Cookie)
	req.Header.Set("Referer", referer)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Mobile Safari/537.36")
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

func MiNi19ViewDoc(requestUrl string, referer string) (doc *html.Node, err error) {
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
	if MiNi19EnableHttpProxy {
		client = MiNi19SetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", MiNi19Cookie)
	req.Header.Set("Referer", referer)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Mobile Safari/537.36")
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

func downloadMiNi19(attachmentUrl string, filePath string, referer string) error {
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
	if MiNi19EnableHttpProxy {
		client = MiNi19SetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.19mini.cn")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", referer)
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
