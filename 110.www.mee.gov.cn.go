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
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	MeeEnableHttpProxy = false
	MeeHttpProxyUrl    = "111.225.152.186:8089"
)

func MeeSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(MeeHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

// ychEduSpider 获取生态环境标准文档
// @Title 获取生态环境标准文档
// @Description https://www.mee.gov.cn/，获取生态环境标准文档
func main() {
	// 第一步获取所有大分类
	categoryIndexRequestUrl := "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/"
	categoryIndexDoc, err := htmlquery.LoadURL(categoryIndexRequestUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	stBzCBaseTabUlLiNodes := htmlquery.Find(categoryIndexDoc, `//ul[@class="stbzCBaseTabUl"]/li`)
	if len(stBzCBaseTabUlLiNodes) <= 0 {
		fmt.Println("没有大分类")
		os.Exit(1)
	}
	for _, liNode := range stBzCBaseTabUlLiNodes {
		smallCategoryRequestUrls := make([]string, 0)
		categoryRequestUrlNode := htmlquery.FindOne(liNode, `./a/@href`)
		categoryRequestUrl := htmlquery.InnerText(categoryRequestUrlNode)
		categoryRequestUrl = strings.ReplaceAll(categoryRequestUrl, ".", "")
		categoryRequestUrl = strings.ReplaceAll(categoryRequestUrl, "/", "")
		categoryRequestUrl = "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/" + categoryRequestUrl + "/"
		categoryEachDoc, err := htmlquery.LoadURL(categoryRequestUrl)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// 第二步查看是否有二级分类，如果有二级分类的话页面会跳转
		regWindowLocationHref := regexp.MustCompile(`window.location.href = "./(.*?)/"`)
		regWindowLocationHrefMatch := regWindowLocationHref.FindAllSubmatch([]byte(htmlquery.InnerText(categoryEachDoc)), -1)
		if len(regWindowLocationHrefMatch) == 0 {
			// 没有二级分类
			smallCategoryRequestUrls = append(smallCategoryRequestUrls, categoryRequestUrl)
		} else {
			// 有二级分类
			preRequestUrl := categoryRequestUrl
			categoryRequestUrl = categoryRequestUrl + string(regWindowLocationHrefMatch[0][1]) + "/"
			categoryEachDoc, err := htmlquery.LoadURL(categoryRequestUrl)
			if err != nil {
				fmt.Println(err)
				continue
			}
			stBzInnerNavLiNodes := htmlquery.Find(categoryEachDoc, `//div[@class="bgtList"]/ul[@class="stBzinnerNav"]/li`)
			if len(stBzInnerNavLiNodes) <= 0 {
				fmt.Println("没有二级分类")
				continue
			}
			for _, liNode := range stBzInnerNavLiNodes {
				liANode := htmlquery.FindOne(liNode, `./a/@href`)
				liAHref := htmlquery.InnerText(liANode)
				liAHref = strings.ReplaceAll(liAHref, ".", "")
				liAHref = strings.ReplaceAll(liAHref, "/", "")
				if len(liAHref) == 0 {
					smallCategoryRequestUrls = append(smallCategoryRequestUrls, categoryRequestUrl)
				} else {
					smallCategoryRequestUrls = append(smallCategoryRequestUrls, preRequestUrl+liAHref+"/")
				}
			}
		}

		if len(smallCategoryRequestUrls) > 0 {
			for _, smallCategoryUrl := range smallCategoryRequestUrls {
				smallCategoryIndexPageUrl := smallCategoryUrl
				smallCategoryIndexPageDoc, err := htmlquery.LoadURL(smallCategoryIndexPageUrl)
				if err != nil {
					fmt.Println(err)
					continue
				}
				// 获取最大分页
				smallCategoryMaxPage := 0
				regCountPage := regexp.MustCompile(`var countPage = ([0-9]*)//共多少页`)
				regCountPageMatch := regCountPage.FindAllSubmatch([]byte(htmlquery.InnerText(smallCategoryIndexPageDoc)), -1)
				if len(regCountPageMatch) > 0 {
					smallCategoryMaxPageStr := string(regCountPageMatch[0][1])
					smallCategoryMaxPage, _ = strconv.Atoi(smallCategoryMaxPageStr)
				}
				for smallCategoryPage := 0; smallCategoryPage < smallCategoryMaxPage; smallCategoryPage++ {
					smallCategoryPageText := ""
					referer := smallCategoryUrl
					if smallCategoryPage > 0 {
						smallCategoryPageText = "_" + strconv.Itoa(smallCategoryPage)
						referer = smallCategoryPageText
					}
					smallCategoryRequestUrl := fmt.Sprintf(smallCategoryUrl+"index%s.shtml", smallCategoryPageText)
					smallCategoryDoc, err := MeeBzList(smallCategoryRequestUrl, referer)
					if err != nil {
						fmt.Println(err)
						continue
					}

					bgtListLiNodes := htmlquery.Find(smallCategoryDoc, `//div[@class="bgtList"]/ul[@class="zzjgGrzyCUl"]/li`)
					if len(bgtListLiNodes) > 0 {
						for _, bgtListLiNode := range bgtListLiNodes {
							fmt.Println("=======================================================")
							liANode := htmlquery.FindOne(bgtListLiNode, `./a/@href`)
							liAHref := htmlquery.InnerText(liANode)
							bzDetailRequestUrl := ""
							if strings.Contains(liAHref, "../") {
								smallCategoryRequestUrlBiasIndex := strings.LastIndex(smallCategoryRequestUrl, "/index")
								smallCategoryRequestUrlTemp := smallCategoryRequestUrl[:smallCategoryRequestUrlBiasIndex]

								biasCount := strings.Count(liAHref, "../")
								smallCategoryRequestUrlTempArr := strings.Split(smallCategoryRequestUrlTemp, "/")
								smallCategoryRequestUrlTempArr = smallCategoryRequestUrlTempArr[:len(smallCategoryRequestUrlTempArr)-biasCount]
								smallCategoryRequestUrlTemp = strings.Join(smallCategoryRequestUrlTempArr, "/") + "/"

								bzDetailRequestUrl = smallCategoryRequestUrlTemp + strings.ReplaceAll(liAHref, "../", "")
							} else {
								liAHref = strings.Replace(liAHref, ".", "", 1)
								smallCategoryRequestUrlIndexShtmlIndex := strings.LastIndex(smallCategoryRequestUrl, "/index")
								bzDetailRequestUrl = smallCategoryRequestUrl[:smallCategoryRequestUrlIndexShtmlIndex] + liAHref
							}

							fmt.Println(smallCategoryRequestUrl)
							fmt.Println(bzDetailRequestUrl)

							bzDetailDoc, err := htmlquery.LoadURL(bzDetailRequestUrl)
							if err != nil {
								fmt.Println(err)
								continue
							}
							bzDetailANodes := htmlquery.Find(bzDetailDoc, `//div[@class="neiright_Content"]/div[@class="neiright_JPZ_GK_CP"]//a`)
							if len(bzDetailANodes) > 0 {
								for _, bzDetailANode := range bzDetailANodes {
									bzDownloadHrefNode := htmlquery.FindOne(bzDetailANode, `./@href`)
									bzDownloadHref := htmlquery.InnerText(bzDownloadHrefNode)
									fmt.Println(bzDownloadHref)
									if strings.Contains(bzDownloadHref, ".pdf") {
										// 中文标题
										chineseTitle := htmlquery.InnerText(bzDetailANode)
										chineseTitle = strings.TrimSpace(chineseTitle)
										chineseTitle = strings.ReplaceAll(chineseTitle, "/", "-")
										chineseTitle = strings.ReplaceAll(chineseTitle, "／", "-")
										chineseTitle = strings.ReplaceAll(chineseTitle, "　", "")
										chineseTitle = strings.ReplaceAll(chineseTitle, " ", "")
										chineseTitle = strings.ReplaceAll(chineseTitle, "：", ":")
										chineseTitle = strings.ReplaceAll(chineseTitle, "—", "-")
										chineseTitle = strings.ReplaceAll(chineseTitle, "－", "-")
										chineseTitle = strings.ReplaceAll(chineseTitle, "（", "(")
										chineseTitle = strings.ReplaceAll(chineseTitle, "）", ")")
										chineseTitle = strings.ReplaceAll(chineseTitle, "《", "")
										chineseTitle = strings.ReplaceAll(chineseTitle, "》", "")
										fmt.Println(chineseTitle)

										// 下载文档URL
										downLoadUrl := bzDownloadHref
										// 查看bzDownloadHref是否含有www.mee.gov.cn
										if !strings.Contains(bzDownloadHref, "www.mee.gov.cn") {
											// 不含有www.mee.gov.cn，下载连接需要处理
											bzDetailRequestUrlBiasTIndex := strings.LastIndex(bzDetailRequestUrl, "/t")
											bzDownloadHref = strings.Replace(bzDownloadHref, ".", "", 1)
											downLoadUrl = bzDetailRequestUrl[:bzDetailRequestUrlBiasTIndex] + bzDownloadHref
										}
										fmt.Println(downLoadUrl)

										filePath := "../www.mee.gov.cn/" + chineseTitle + ".pdf"
										if _, err := os.Stat(filePath); err != nil {
											// 开始下载
											fmt.Println("=======开始下载========")
											err = downloadMee(downLoadUrl, bzDetailRequestUrl, filePath)
											if err != nil {
												fmt.Println(err)
												continue
											}
											fmt.Println("=======开始完成========")
										}
										time.Sleep(time.Millisecond * 100)
									}
								}
							}
							fmt.Println("=======================================================")
						}
					}
				}
			}
		}
	}
}

func MeeBzList(requestUrl string, referer string) (doc *html.Node, err error) {
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
	if MeeEnableHttpProxy {
		client = MeeSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接
	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.mee.gov.cn")
	req.Header.Set("Origin", "https://www.mee.gov.cn/")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
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

func downloadMee(attachmentUrl string, referer string, filePath string) error {
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
	if MeeEnableHttpProxy {
		client = MeeSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.mee.gov.cn")
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
