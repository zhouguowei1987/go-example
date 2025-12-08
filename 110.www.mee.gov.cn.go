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

type MeeCategory struct {
	name string
	url  string
}

var meecategory = []MeeCategory{
	{name: "水环境质量标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/shjbh/shjzlbz/"},
	{name: "水污染物排放标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/shjbh/swrwpfbz/"},
	{name: "相关标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/shjbh/xgbzh/"},
	{name: "大气环境质量标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/dqhjbh/dqhjzlbz/"},
	{name: "大气固定源污染物排放标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/dqhjbh/dqgdwrywrwpfbz/"},
	{name: "大气移动源污染物排放标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/dqhjbh/dqydywrwpfbz/"},
	{name: "相关标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/dqhjbh/xgbz/"},
	{name: "声环境质量标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/wlhj/shjzlbz/"},
	{name: "环境噪声排放标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/wlhj/hjzspfbz/"},
	{name: "土壤环境保护", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/trhj/"},
	{name: "固体废物污染控制标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/gthw/gtfwwrkzbz/"},
	{name: "危险废物鉴别方法标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/gthw/wxfwjbffbz/"},
	{name: "其他相关标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/gthw/qtxgbz/"},
	{name: "电磁辐射标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/hxxhj/dcfsbz/"},
	{name: "放射性环境标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/hxxhj/fsxhjbz/"},
	{name: "相关监测方法标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/hxxhj/xgjcffbz/"},
	{name: "相关标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/hxxhj/xgbz/"},
	{name: "生态环境保护", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/stzl/"},
	{name: "环境影响评价", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/hp/"},
	{name: "排污许可", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/pwxk/"},
	{name: "清洁生产标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/other/qjscbz/"},
	{name: "环境影响评价技术导则", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/other/pjjsdz/"},
	{name: "环保验收技术规范", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/other/hbysjsgf/"},
	{name: "环境标志产品标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/other/hjbz/"},
	{name: "环保产品技术要求", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/other/hbcpjsyq/"},
	{name: "环境保护工程技术规范", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/other/hjbhgc/"},
	{name: "环境保护信息标准", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/other/xxbz/"},
	{name: "其他", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/other/qt/"},
	{name: "污染防治技术政策", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/wrfzjszc/"},
	{name: "可行技术指南", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/kxxjszn/"},
	{name: "环境监测方法标准及监测规范", url: "https://www.mee.gov.cn/ywgz/fgbz/bz/bzwb/jcffbz/"},
}

var MeeCookie = "wdcid=3fa507778241451a; viewsid=3484113bb5b0423ca247eef0f4d6a1d4; ariauseGraymode=false; Hm_lvt_0f50400dd25408cef4f1afb556ccb34f=1763092858,1764820005; Hm_lpvt_0f50400dd25408cef4f1afb556ccb34f=1764820005; HMACCOUNT=4E5B3419A3141A8E; arialoadData=true; wdlast=1764820006; wdses=2c58020ac2fbd186"

// ychEduSpider 获取生态环境标准文档
// @Title 获取生态环境标准文档
// @Description https://www.mee.gov.cn/，获取生态环境标准文档
func main() {
	for _, category := range meecategory {
		fmt.Println(category.name, category.url)
		page := 0
		isPageListGo := true
		for isPageListGo {
			requestListUrl := category.url + "index.shtml"
			referUrl := category.url
			if page > 0 {
				requestListUrl = fmt.Sprintf(category.url+"index_%d.shtml", page)
			}
			if page >= 2 {
				referUrl = fmt.Sprintf(category.url+"index_%d.shtml", page-1)
			}
			fmt.Println(requestListUrl)
			meeBzListDoc, err := MeeBzHtmlDoc(requestListUrl, referUrl)
			if err != nil {
				fmt.Println(err)
				break
			}
			liNodes := htmlquery.Find(meeBzListDoc, `//html/body/div[4]/div[2]/div/div[2]/ul[2]/li`)
			if len(liNodes) >= 1 {
				for _, liNode := range liNodes {
					detailUrlNode := htmlquery.FindOne(liNode, `./a/@href`)
					if detailUrlNode == nil {
						fmt.Println("没有文档详情链接，跳过")
						continue
					}
					detailUrl := category.url + htmlquery.InnerText(detailUrlNode)
					fmt.Println(detailUrl)
					//os.Exit(1)

					meeDetailDoc, err := MeeBzHtmlDoc(detailUrl, requestListUrl)
					if err != nil {
						fmt.Println("获取文档详情失败，跳过")
						continue
					}

					bzDetailANodes := htmlquery.Find(meeDetailDoc, `//div[@class="neiright_Content"]/div[@class="neiright_JPZ_GK_CP"]/a`)
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
									bzDetailRequestUrlBiasTIndex := strings.LastIndex(detailUrl, "/t")
									bzDownloadHref = strings.Replace(bzDownloadHref, ".", "", 1)
									downLoadUrl = detailUrl[:bzDetailRequestUrlBiasTIndex] + bzDownloadHref
								}
								fmt.Println(downLoadUrl)

								filePath := "../www.mee.gov.cn/" + chineseTitle + ".pdf"
								_, err = os.Stat(filePath)
								if err == nil {
									fmt.Println("文档已下载过，跳过")
									continue
								}

								// 开始下载
								fmt.Println("=======开始下载========")
								err = downloadMee(downLoadUrl, detailUrl, filePath)
								if err != nil {
									fmt.Println(err)
									continue
								}
								//复制文件
								tempFilePath := strings.ReplaceAll(filePath, "../www.mee.gov.cn", "../upload.doc88.com/www.mee.gov.cn")
								err = copyMeeFile(filePath, tempFilePath)
								if err != nil {
									fmt.Println(err)
									continue
								}
								fmt.Println("=======完成下载========")

								// 设置倒计时
								DownLoadTMeeTimeSleep := 10
								for i := 1; i <= DownLoadTMeeTimeSleep; i++ {
									time.Sleep(time.Second)
									fmt.Println("category_name = "+category.name+",===page = "+strconv.Itoa(page)+"===title="+chineseTitle+"===========操作完成，", "暂停", DownLoadTMeeTimeSleep, "秒，倒计时", i, "秒===========")
								}
							}
						}
					}
				}
				DownLoadMeePageTimeSleep := 10
				// DownLoadMeePageTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadMeePageTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("category_name = "+category.name+",===page = "+strconv.Itoa(page)+"========= 暂停", DownLoadMeePageTimeSleep, "秒 倒计时", i, "秒===========")
				}
				page++
			} else {
				page = 0
				isPageListGo = false
				break
			}
		}
	}
}

func MeeBzHtmlDoc(requestUrl string, referer string) (doc *html.Node, err error) {
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
	req.Header.Set("Cookie", MeeCookie)
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
	req.Header.Set("Cookie", MeeCookie)
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

func copyMeeFile(src, dst string) (err error) {
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
