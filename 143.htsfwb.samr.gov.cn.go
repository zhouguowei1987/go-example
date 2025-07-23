package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	_ "os"
	"path/filepath"
	_ "path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	_ "golang.org/x/net/html"
)

var HtSfWbEnableHttpProxy = false
var HtSfWbHttpProxyUrl = "111.225.152.186:8089"
var HtSfWbHttpProxyUrlArr = make([]string, 0)

func HtSfWbHttpProxy() error {
	pageMax := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
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
					HtSfWbHttpProxyUrlArr = append(HtSfWbHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					HtSfWbHttpProxyUrlArr = append(HtSfWbHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func HtSfWbSetHttpProxy() (httpclient *http.Client) {
	if HtSfWbHttpProxyUrl == "" {
		if len(HtSfWbHttpProxyUrlArr) <= 0 {
			err := HtSfWbHttpProxy()
			if err != nil {
				HtSfWbSetHttpProxy()
			}
		}
		HtSfWbHttpProxyUrl = HtSfWbHttpProxyUrlArr[0]
		if len(HtSfWbHttpProxyUrlArr) >= 2 {
			HtSfWbHttpProxyUrlArr = HtSfWbHttpProxyUrlArr[1:]
		} else {
			HtSfWbHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(HtSfWbHttpProxyUrl)
	ProxyURL, _ := url.Parse(HtSfWbHttpProxyUrl)
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

var HtSfWbCookie = "Hm_lvt_54db9897e5a65f7a7b00359d86015d8d=1752905161,1753168348; HMACCOUNT=1CCD0111717619C6; Hm_lpvt_54db9897e5a65f7a7b00359d86015d8d=1753254437; __jsluid_s=e2a3e160f8044ebe5131f1716a90b079; samr=isopen"

// 下载地方合同示范文本
// @Title 下载地方合同示范文本
// @Description https://htsfwb.samr.gov.cn/Local/，下载地方合同示范文本
func main() {
	page := 1
	maxPage := 50
	// true:地方合同范本 false:部委合同范本
	locHtSfWb := false
	suffixHtSfWb := ".docx"
	// 1:docx 2:pdf
	typeHtSfWb := 1
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		pageListUrl := fmt.Sprintf("https://htsfwb.samr.gov.cn/api/content/SearchTemplates?loc=%t&p=%d&key=", locHtSfWb, page)
		fmt.Println(pageListUrl)
		queryHtSfWbListResponseData, err := QueryHtSfWbList(pageListUrl)
		if err != nil {
			HtSfWbHttpProxyUrl = ""
			fmt.Println(err)
			break
		}
		if len(queryHtSfWbListResponseData) <= 0 {
			fmt.Println("没有更多数据，停止")
			break
		}
		for _, htSfWb := range queryHtSfWbListResponseData {
			fmt.Println("=====================开始处理数据 page = ", page, "=========================")

			title := htSfWb.Title
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "--", "-")
			fmt.Println(title)

			filePath := "../htsfwb.samr.gov.cn/" + title + suffixHtSfWb
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载========")

			downloadUrl := fmt.Sprintf("https://htsfwb.samr.gov.cn/api/File/DownTemplate?id=%s&type=%d", htSfWb.Id, typeHtSfWb)
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")

			detailUrl := fmt.Sprintf("https://htsfwb.samr.gov.cn/View?id=%s", htSfWb.Id)
			fmt.Println(detailUrl)

			err = downloadHtSfWb(downloadUrl, detailUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "../htsfwb.samr.gov.cn", "../upload.doc88.com/htsfwb.samr.gov.cn")
			err = copyHtSfWbFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			//DownLoadHtSfWbTimeSleep := 10
			DownLoadHtSfWbTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadHtSfWbTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadHtSfWbTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadHtSfWbPageTimeSleep := 10
		// DownLoadHtSfWbPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadHtSfWbPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadHtSfWbPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

type QueryHtSfWbListResponse struct {
	Data      []QueryHtSfWbListResponseData `json:"Data"`
	Page      int                           `json:"Page"`
	Total     int                           `json:"Total"`
	TotalPage int                           `json:"TotalPage"`
}

type QueryHtSfWbListResponseData struct {
	Id    string `json:"Id"`
	Title string `json:"Title"`
}

func QueryHtSfWbList(requestUrl string) (queryHtSfWbListResponseData []QueryHtSfWbListResponseData, err error) {
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
	if HtSfWbEnableHttpProxy {
		client = HtSfWbSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	queryHtSfWbListResponse := QueryHtSfWbListResponse{}
	if err != nil {
		return queryHtSfWbListResponseData, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", HtSfWbCookie)
	req.Header.Set("Host", "htsfwb.samr.gov.cn")
	req.Header.Set("Origin", "https://htsfwb.samr.gov.cn")
	req.Header.Set("Referer", "https://htsfwb.samr.gov.cn/Local")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryHtSfWbListResponseData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryHtSfWbListResponseData, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryHtSfWbListResponseData, err
	}
	err = json.Unmarshal(respBytes, &queryHtSfWbListResponse)
	if err != nil {
		return queryHtSfWbListResponseData, err
	}
	queryHtSfWbListResponseData = queryHtSfWbListResponse.Data
	return queryHtSfWbListResponseData, nil
}

func downloadHtSfWb(attachmentUrl string, referer string, filePath string) error {
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
	if HtSfWbEnableHttpProxy {
		client = HtSfWbSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", HtSfWbCookie)
	req.Header.Set("Host", "htsfwb.samr.gov.cn")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
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

func copyHtSfWbFile(src, dst string) (err error) {
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
