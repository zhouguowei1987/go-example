package main

import (
	"errors"
	"fmt"
	"io"

	// 	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const (
	HqWxEnableHttpProxy = false
	HqWxHttpProxyUrl    = "111.225.152.186:8089"
)

func HqWxSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(HqWxHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

var HqWxCookie = "sajssdk_2015_cross_new_user=1; Hm_lvt_1adfcf85508daf29422b179368502bb7=1781489381; HMACCOUNT=9C0CD19686802BBF; commencookie=003a8e70976e6b22e6859f579a44cc43aa577ee65e4e479bd38356516e55e8f7b9b84245d9cab5086159a6cfc2ceb3e5b501bf48b71eae7c9310f5c23c01420e8fdf91c860156c1b0a752d693608b81ce5b844d13a58a6b4b4b76b2599cad1a07e27c3d2a5b8c4d5d7afd205b3ff2b5f7a3584afcc56f3d6dfea1190a099def491be5d3ed8fd4904077b257a1791edbb43853727fa5c5af316319baff331b5e35d58963bb410392bc907ced17af8b9b4ad9b695b03398a69d74855474d5ff2155ebe865becbcffdafb10d7ba98d921a7; web_id=4223; webidCors=4223; trafficUuid=b0ae23e7-5c63-4f5b-bee2-ea8c70c74b90; _gid=GA1.2.1990149751.1781489382; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%2219ec90aacab3ca-058afdd2a11d6d8-4c657b58-2073600-19ec90aacacc9d%22%2C%22first_id%22%3A%22%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%2C%22%24latest_referrer%22%3A%22%22%7D%2C%22identities%22%3A%22eyIkaWRlbnRpdHlfY29va2llX2lkIjoiMTllYzkwYWFjYWIzY2EtMDU4YWZkZDJhMTFkNmQ4LTRjNjU3YjU4LTIwNzM2MDAtMTllYzkwYWFjYWNjOWQifQ%3D%3D%22%2C%22history_login_id%22%3A%7B%22name%22%3A%22%22%2C%22value%22%3A%22%22%7D%2C%22%24device_id%22%3A%2219ec90aacab3ca-058afdd2a11d6d8-4c657b58-2073600-19ec90aacacc9d%22%7D; _ga_MWMW10WGQW=GS2.1.s1781489382$o1$g1$t1781489948$j60$l0$h0; _ga=GA1.1.1074127803.1781489382; Hm_lpvt_1adfcf85508daf29422b179368502bb7=1781489948; lastProductId=173"

// 获取环球网校职业考试免费资料
// @Title 获取环球网校职业考试免费资料
// @Description https://m.hqwx.com/ 获取环球网校职业考试免费资料
func main() {
	// 获取分类
	hqWxCategoryUrl := "https://m.hqwx.com/ziliao/"
	hqWxCategoryRefererUrl := "https://m.hqwx.com/"
	hqWxCategoryPathUrl := "/ziliao/"
	hqWxCategoryDoc, err := QueryHqWxHtml(hqWxCategoryUrl, hqWxCategoryRefererUrl, hqWxCategoryPathUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	hqWxCategoryLiNodes := htmlquery.Find(hqWxCategoryDoc, `//div[@class="container"]/div[@class="group"]/ul/li`)
	if len(hqWxCategoryLiNodes) <= 0 {
		fmt.Println("类别数量为零")
		os.Exit(1)
	}
	for _, hqWxCategoryLiNode := range hqWxCategoryLiNodes {
		categoryNameNode := htmlquery.FindOne(hqWxCategoryLiNode, `./a`)
		categoryName := htmlquery.InnerText(categoryNameNode)
		categoryName = strings.TrimSpace(categoryName)

		categoryUrlNode := htmlquery.FindOne(hqWxCategoryLiNode, `./a/@href`)
		categoryUrl := htmlquery.InnerText(categoryUrlNode)
		categoryUrl = strings.TrimSpace(categoryUrl)

		page := 1
		isPageListGo := true
		for isPageListGo {
			requestUrl := fmt.Sprintf("https://m.hqwx.com"+categoryUrl+"?pageNo=%d", page)
			pathUrl := fmt.Sprintf(categoryUrl+"?pageNo=%d", page)
			fmt.Println(requestUrl)
			pageDoc, err := QueryHqWxHtml(requestUrl, categoryUrl, pathUrl)
			if err != nil {
				fmt.Println(err)
				page = 1
				isPageListGo = false
				continue
			}
			liNodes := htmlquery.Find(pageDoc, `//li`)
			if len(liNodes) <= 0 {
				page = 1
				isPageListGo = false
				break
			}

			for _, liNode := range liNodes {
				fmt.Println("============================================================")
				fmt.Println("================当前页列表URL", requestUrl, "=================")
				downloadHrefNode := htmlquery.FindOne(liNode, `./button[@class="download-btn dbutton"]/@data-url`)
				if downloadHrefNode == nil {
					fmt.Println("未找到下载文件节点，跳过")
					continue
				}
				downloadUrl := htmlquery.InnerText(downloadHrefNode)
				fileExt := path.Ext(downloadUrl)
				if !HqWxStrInArray(fileExt, []string{".doc", ".docx", ".pdf"}) {
					fmt.Println("文件后缀：" + fileExt + "不在下载后缀列表")
					continue
				}
				fmt.Println(downloadUrl)

				titleNode := htmlquery.FindOne(liNode, `./div[@class="info"]/div[@class="name"]/a`)
				if titleNode == nil {
					fmt.Println("未找到标题节点，跳过")
					continue
				}
				title := htmlquery.InnerText(titleNode)
				title = strings.TrimSpace(title)
				fmt.Println(title)

				filePath := "../m.hqwx.com/m.hqwx.com/" + categoryName + "/" + title
				fmt.Println(filePath)

				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				fmt.Println("=======开始下载" + title + "========")
				err = downloadHqWx(downloadUrl, requestUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "m.hqwx.com/m.hqwx.com/"+categoryName, "m.hqwx.com/temp-m.hqwx.com/"+categoryName)
				err = copyHqWxFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")
				DownLoadHqWxTimeSleep := 10
				// DownLoadHqWxTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadHqWxTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadHqWxTimeSleep, "秒 倒计时", i, "秒===========")
				}
			}
			DownLoadHqWxPageTimeSleep := 10
			// DownLoadHqWxPageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadHqWxPageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadHqWxPageTimeSleep, "秒 倒计时", i, "秒===========")
			}
			page++
		}
	}
}

func QueryHqWxHtml(requestUrl string, referer string, path string) (doc *html.Node, err error) {
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
	req.Header.Set("authority", "m.hqwx.com")
	req.Header.Set("method", "GET")
	req.Header.Set("path", path)
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", HqWxCookie)
	req.Header.Set("Host", "m.hqwx.com")
	req.Header.Set("Origin", "https://m.hqwx.com")
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

func downloadHqWx(attachmentUrl string, referer string, filePath string) error {
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
	if HqWxEnableHttpProxy {
		client = HqWxSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "m.hqwx.com")
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

func copyHqWxFile(src, dst string) (err error) {
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

	// 创建一个文件用于保存
	fileDiv := filepath.Dir(dst)
	if _, err = os.Stat(fileDiv); err != nil {
		if os.MkdirAll(fileDiv, 0o777) != nil {
			return err
		}
	}
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
func HqWxStrInArray(str string, data []string) bool {
	if len(data) > 0 {
		for _, row := range data {
			if str == row {
				return true
			}
		}
	}
	return false
}
