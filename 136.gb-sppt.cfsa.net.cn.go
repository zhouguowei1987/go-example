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

var GbSpPtEnableHttpProxy = false
var GbSpPtHttpProxyUrl = "111.225.152.186:8089"
var GbSpPtHttpProxyUrlArr = make([]string, 0)

func GbSpPtHttpProxy() error {
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
					GbSpPtHttpProxyUrlArr = append(GbSpPtHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					GbSpPtHttpProxyUrlArr = append(GbSpPtHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func GbSpPtSetHttpProxy() (httpclient *http.Client) {
	if GbSpPtHttpProxyUrl == "" {
		if len(GbSpPtHttpProxyUrlArr) <= 0 {
			err := GbSpPtHttpProxy()
			if err != nil {
				GbSpPtSetHttpProxy()
			}
		}
		GbSpPtHttpProxyUrl = GbSpPtHttpProxyUrlArr[0]
		if len(GbSpPtHttpProxyUrlArr) >= 2 {
			GbSpPtHttpProxyUrlArr = GbSpPtHttpProxyUrlArr[1:]
		} else {
			GbSpPtHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(GbSpPtHttpProxyUrl)
	ProxyURL, _ := url.Parse(GbSpPtHttpProxyUrl)
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

type QueryGbSpPtListFormData struct {
	isLength      int
	num_tn        int
	standard_type string
	keyword       string
}

type DownloadGbSpPtFormData struct {
	task       string
	accessData string
	bzlb       string
	fact_name  string
	file_guid  string
	keyword    string
}

var GbSpPtCookie = "name=value; cookieName=cookieValue; JSESSIONID=1939457719D7E496EDAAB20F581FA15F"

// 下载食品安全国家标准数据文档
// @Title 下载食品安全国家标准数据文档
// @Description https://sppt.cfsa.net.cn:8086/db/，下载食品安全国家标准数据文档
func main() {
	pageListUrl := "https://sppt.cfsa.net.cn:8086/db?task=indexSearch"
	fmt.Println(pageListUrl)
	queryGbSpPtListFormData := QueryGbSpPtListFormData{
		isLength:      9999,
		num_tn:        2,
		standard_type: "",
		keyword:       "",
	}
	queryGbSpPtListResponse, err := QueryGbSpPtList(pageListUrl, queryGbSpPtListFormData)
	if err != nil {
		GbSpPtHttpProxyUrl = ""
		fmt.Println(err)
	}
	for id_index, gbSpPt := range queryGbSpPtListResponse {
		fmt.Println("=====================开始处理数据 id_index = ", id_index, "=========================")
		code := gbSpPt.CODE
		fmt.Println(code)

		title := gbSpPt.TITLE
		title = strings.TrimSpace(title)
		title = strings.ReplaceAll(title, "-", "")
		title = strings.ReplaceAll(title, " ", "")
		title = strings.ReplaceAll(title, "|", "-")
		fmt.Println(title)

		if len(gbSpPt.FJ) <= 0 {
			fmt.Println("数据不完整，跳过")
			continue
		}
		id_f := gbSpPt.FJ[0].ID_F
		fmt.Println(id_f)

		filePath := "../sppt.cfsa.net.cn/" + title + "(" + code + ")" + ".pdf"
		fmt.Println(filePath)

		_, err = os.Stat(filePath)
		if err == nil {
			fmt.Println("文档已下载过，跳过")
			continue
		}

		fmt.Println("=======开始下载========")

		downloadGbSpPtUrl := "https://sppt.cfsa.net.cn:8086/cfsa_aiguo"
		fmt.Println(downloadGbSpPtUrl)
		downloadGbSpPtFormData := DownloadGbSpPtFormData{
			task:       "d_p",
			accessData: "gj",
			bzlb:       "",
			fact_name:  "",
			file_guid:  id_f,
			keyword:    "",
		}
		err := downloadGbSpPt(downloadGbSpPtUrl, downloadGbSpPtFormData, filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//复制文件
		tempFilePath := strings.ReplaceAll(filePath, "../sppt.cfsa.net.cn", "../upload.doc88.com/hbba.sacinfo.org.cn")
		err = copyGbSpPtFile(filePath, tempFilePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("=======完成下载========")
		//DownLoadGbSpPtTimeSleep := 10
		DownLoadGbSpPtTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadGbSpPtTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("title="+title+"===========下载", title, "成功，暂停", DownLoadGbSpPtTimeSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

type QueryGbSpPtListResponse struct {
	CODE      string              `json:"CODE"`
	FJ        []QueryGbSpPtListFJ `json:"FJ"`
	ID        string              `json:"ID"`
	PDATE     string              `json:"PDATE"`
	SSRQ      string              `json:"SSRQ"`
	TABLENAME string              `json:"TABLENAME"`
	TITLE     string              `json:"TITLE"`
}

type QueryGbSpPtListFJ struct {
	FACT_NAME string `json:"FACT_NAME"`
	ID_F      string `json:"ID_F"`
}

func QueryGbSpPtList(requestUrl string, queryGbSpPtListFormData QueryGbSpPtListFormData) (queryGbSpPtListResponse []QueryGbSpPtListResponse, err error) {
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
	if GbSpPtEnableHttpProxy {
		client = GbSpPtSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("isLength", strconv.Itoa(queryGbSpPtListFormData.isLength))
	postData.Add("num_tn", strconv.Itoa(queryGbSpPtListFormData.num_tn))
	postData.Add("standard_type", queryGbSpPtListFormData.standard_type)
	postData.Add("keyword", queryGbSpPtListFormData.keyword)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	queryGbSpPtListResponse = []QueryGbSpPtListResponse{}
	if err != nil {
		return queryGbSpPtListResponse, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", GbSpPtCookie)
	req.Header.Set("Host", "sppt.cfsa.net.cn:8086")
	req.Header.Set("Origin", "https://sppt.cfsa.net.cn:8086")
	req.Header.Set("Referer", "https://sppt.cfsa.net.cn:8086/db")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryGbSpPtListResponse, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryGbSpPtListResponse, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryGbSpPtListResponse, err
	}
	err = json.Unmarshal(respBytes, &queryGbSpPtListResponse)
	if err != nil {
		return queryGbSpPtListResponse, err
	}
	return queryGbSpPtListResponse, nil
}

func downloadGbSpPt(requestUrl string, downloadGbSpPtFormData DownloadGbSpPtFormData, filePath string) error {
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
	if GbSpPtEnableHttpProxy {
		client = GbSpPtSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("task", downloadGbSpPtFormData.task)
	postData.Add("accessData", downloadGbSpPtFormData.accessData)
	postData.Add("bzlb", downloadGbSpPtFormData.bzlb)
	postData.Add("fact_name", downloadGbSpPtFormData.fact_name)
	postData.Add("file_guid", downloadGbSpPtFormData.file_guid)
	postData.Add("keyword", downloadGbSpPtFormData.keyword)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", GbSpPtCookie)
	req.Header.Set("Host", "sppt.cfsa.net.cn:8086")
	req.Header.Set("Origin", "https://sppt.cfsa.net.cn:8086")
	req.Header.Set("Referer", "https://sppt.cfsa.net.cn:8086/db")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
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

func copyGbSpPtFile(src, dst string) (err error) {
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
