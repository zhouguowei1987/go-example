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
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const (
	CvmaEnableHttpProxy = false
	CvmaHttpProxyUrl    = "111.225.152.186:8089"
)

func CvmaSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(CvmaHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

var CvmaCookie = "server_name_session=f4a44218d8a35a36bc8a7185ffa36b37; PHPSESSID=70vt2ul33djvpvg3im4n96oe0h"

// ychEduSpider 获取中国兽医协会标准文档
// @Title 获取中国兽医协会标准文档
// @Description https://www.cvma.org.cn/，获取中国兽医协会标准文档
func main() {
	page := 1
	isPageListGo := true
	for isPageListGo {
		requestListUrl := fmt.Sprintf("https://www.cvma.org.cn/api/search-biaozhunku.php?siteid=10000&type=1&page=2&pagesize=10&period=&bznum=&name=&unit=&schedule=已发布&page=%d&pagesize=100", page)
		// fmt.Println(requestListUrl)
		cvmaBzListResponse, err := GetCvmaBzList(requestListUrl)
		if err != nil {
			fmt.Println(err)
			break
		}
		if len(cvmaBzListResponse.Data) >= 1 {
			for _, cvmaBz := range cvmaBzListResponse.Data {
				fmt.Println("=====================开始处理列表-分割线==========================")

				fmt.Println("=======page = " + strconv.Itoa(page) + "=========")

				// 标题
				name := cvmaBz.Name
				name = strings.TrimSpace(name)
				name = strings.ReplaceAll(name, "/", "-")
				name = strings.ReplaceAll(name, "：", ":")
				name = strings.ReplaceAll(name, "—", "-")
				name = strings.ReplaceAll(name, "－", "-")
				name = strings.ReplaceAll(name, "—", "-")
				name = strings.ReplaceAll(name, "（", "(")
				name = strings.ReplaceAll(name, "）", ")")
				name = strings.ReplaceAll(name, "《", "")
				name = strings.ReplaceAll(name, "》", "")
				fmt.Println(name)
				// 标准号
				bznum := cvmaBz.Bznum
				bznum = strings.TrimSpace(bznum)
				bznum = strings.ReplaceAll(bznum, "/", "-")
				fmt.Println(bznum)

				filePath := "../www.cvma.org.cn/" + name + "(" + bznum + ")" + ".pdf"
				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}
				// 标准详情地址
				regDetailUrl := regexp.MustCompile(`<a target=\'_blank\' href=\'(.*?)\'><span class=\'el-link--inner\'>预览</span>`)
				regDetailUrlMatch := regDetailUrl.FindAllSubmatch([]byte(cvmaBz.Annex), -1)
				if len(regDetailUrlMatch) <= 0 {
					fmt.Println("未找到详情地址，跳过")
					continue
				}
				detailUrl := "https://www.cvma.org.cn" + string(regDetailUrlMatch[0][1])
				fmt.Println(detailUrl)

				detailDoc, err := CvmaBzHtmlDoc(detailUrl, "https://www.cvma.org.cn/6788/index.html")
				if err != nil {
					fmt.Println(err)
					continue
				}
				// 下载地址
				regDownloadUrl := regexp.MustCompile(`<input type="hidden" id="pdf" value="(.*?)"/>`)
				regDownloadUrlMatch := regDownloadUrl.FindAllSubmatch([]byte(htmlquery.OutputHTML(detailDoc, true)), -1)
				if len(regDownloadUrlMatch) <= 0 {
					fmt.Println("未找下载地址，跳过")
					continue
				}
				downLoadUrl := string(regDownloadUrlMatch[0][1])
				fmt.Println(downLoadUrl)
				// 开始下载
				fmt.Println("=======开始下载========")
				err = downloadCvma(downLoadUrl, requestListUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "www.cvma.org.cn", "temp-www.cvma.org.cn")
				err = copyCvmaFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======完成下载========")

				// 设置倒计时
				DownLoadTCvmaTimeSleep := 10
				for i := 1; i <= DownLoadTCvmaTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("===page = "+strconv.Itoa(page)+"===name="+name+"===========操作完成，", "暂停", DownLoadTCvmaTimeSleep, "秒，倒计时", i, "秒===========")
				}
			}
			DownLoadCvmaPageTimeSleep := 10
			// DownLoadCvmaPageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadCvmaPageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("===page = "+strconv.Itoa(page)+"========= 暂停", DownLoadCvmaPageTimeSleep, "秒 倒计时", i, "秒===========")
			}
			page++
		} else {
			page = 0
			isPageListGo = false
			break
		}
	}
}

type CvmaBzListResponse struct {
	Count string                   `json:"count"`
	Code  int                      `json:"code"`
	Data  []CvmaBzListResponseData `json:"data"`
	Msg   string                   `json:"msg"`
}

type CvmaBzListResponseData struct {
	Annex string `json:"annex"`
	Bznum string `json:"bznum"`
	Name  string `json:"name"`
}

func GetCvmaBzList(requestUrl string) (cvmaBzListResponse CvmaBzListResponse, err error) {
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
	if CvmaEnableHttpProxy {
		client = CvmaSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return cvmaBzListResponse, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.cvma.org.cn")
	req.Header.Set("Origin", "https://www.cvma.org.cn/")
	req.Header.Set("Referer", "https://www.cvma.org.cn/6788/index.html")
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
		return cvmaBzListResponse, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cvmaBzListResponse, err
	}
	err = json.Unmarshal(respBytes, &cvmaBzListResponse)
	if err != nil {
		return cvmaBzListResponse, err
	}
	return cvmaBzListResponse, nil
}

func CvmaBzHtmlDoc(requestUrl string, referer string) (doc *html.Node, err error) {
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
	if CvmaEnableHttpProxy {
		client = CvmaSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接
	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	// 	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", CvmaCookie)
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	req.Header.Set("Host", "www.cvma.org.cn")
	req.Header.Set("Origin", "https://www.cvma.org.cn/")
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

func downloadCvma(attachmentUrl string, referer string, filePath string) error {
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
	if CvmaEnableHttpProxy {
		client = CvmaSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", CvmaCookie)
	req.Header.Set("Host", "www.cvma.org.cn")
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

func copyCvmaFile(src, dst string) (err error) {
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
