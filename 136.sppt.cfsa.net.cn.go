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

var SpPtCfSaEnableHttpProxy = false
var SpPtCfSaHttpProxyUrl = "111.225.152.186:8089"
var SpPtCfSaHttpProxyUrlArr = make([]string, 0)

func SpPtCfSaHttpProxy() error {
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
					SpPtCfSaHttpProxyUrlArr = append(SpPtCfSaHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					SpPtCfSaHttpProxyUrlArr = append(SpPtCfSaHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func SpPtCfSaSetHttpProxy() (httpclient *http.Client) {
	if SpPtCfSaHttpProxyUrl == "" {
		if len(SpPtCfSaHttpProxyUrlArr) <= 0 {
			err := SpPtCfSaHttpProxy()
			if err != nil {
				SpPtCfSaSetHttpProxy()
			}
		}
		SpPtCfSaHttpProxyUrl = SpPtCfSaHttpProxyUrlArr[0]
		if len(SpPtCfSaHttpProxyUrlArr) >= 2 {
			SpPtCfSaHttpProxyUrlArr = SpPtCfSaHttpProxyUrlArr[1:]
		} else {
			SpPtCfSaHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(SpPtCfSaHttpProxyUrl)
	ProxyURL, _ := url.Parse(SpPtCfSaHttpProxyUrl)
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

type QuerySpPtCfSaListFormData struct {
	isLength      int
	num_tn        int
	standard_type string
	keyword       string
}

type DownloadSpPtCfSaFormData struct {
	task       string
	accessData string
	bzlb       string
	fact_name  string
	file_guid  string
	keyword    string
}

var SpPtCfSaCookie = "name=value; cookieName=cookieValue; JSESSIONID=1939457719D7E496EDAAB20F581FA15F"

// 下载食品安全国家标准数据文档
// @Title 下载食品安全国家标准数据文档
// @Description https://sppt.cfsa.net.cn:8086/db/，下载食品安全国家标准数据文档
func main() {
	pageListUrl := "https://sppt.cfsa.net.cn:8086/db?task=indexSearch"
	fmt.Println(pageListUrl)
	querySpPtCfSaListFormData := QuerySpPtCfSaListFormData{
		isLength:      9999,
		num_tn:        2,
		standard_type: "",
		keyword:       "",
	}
	querySpPtCfSaListResponse, err := QuerySpPtCfSaList(pageListUrl, querySpPtCfSaListFormData)
	if err != nil {
		SpPtCfSaHttpProxyUrl = ""
		fmt.Println(err)
	}
	for id_index, spPtCfSa := range querySpPtCfSaListResponse {
		fmt.Println("=====================开始处理数据 id_index = ", id_index, "=========================")
		code := spPtCfSa.CODE
		fmt.Println(code)

		title := spPtCfSa.TITLE
		title = strings.TrimSpace(title)
		title = strings.ReplaceAll(title, "-", "")
		title = strings.ReplaceAll(title, " ", "")
		title = strings.ReplaceAll(title, "|", "-")
		fmt.Println(title)

		if len(spPtCfSa.FJ) < 0 {
			fmt.Println("数据不完整，跳过")
			continue
		}
		id_f := spPtCfSa.FJ[0].ID_F
		fmt.Println(id_f)

		filePath := "../sppt.cfsa.net.cn/" + title + "(" + code + ")" + ".pdf"
		fmt.Println(filePath)

		_, err = os.Stat(filePath)
		if err == nil {
			fmt.Println("文档已下载过，跳过")
			continue
		}

		fmt.Println("=======开始下载========")

		downloadSpPtCfSaUrl := "https://sppt.cfsa.net.cn:8086/cfsa_aiguo"
		fmt.Println(downloadSpPtCfSaUrl)
		downloadSpPtCfSaFormData := DownloadSpPtCfSaFormData{
			task:       "d_p",
			accessData: "gj",
			bzlb:       "",
			fact_name:  "",
			file_guid:  id_f,
			keyword:    "",
		}
		err := downloadSpPtCfSa(downloadSpPtCfSaUrl, downloadSpPtCfSaFormData, filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("=======完成下载========")
		//DownLoadSpPtCfSaTimeSleep := 10
		DownLoadSpPtCfSaTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadSpPtCfSaTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("title="+title+"===========下载", title, "成功，暂停", DownLoadSpPtCfSaTimeSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

type QuerySpPtCfSaListResponse struct {
	CODE      string                `json:"CODE"`
	FJ        []QuerySpPtCfSaListFJ `json:"FJ"`
	ID        string                `json:"ID"`
	PDATE     string                `json:"PDATE"`
	SSRQ      string                `json:"SSRQ"`
	TABLENAME string                `json:"TABLENAME"`
	TITLE     string                `json:"TITLE"`
}

type QuerySpPtCfSaListFJ struct {
	FACT_NAME string `json:"FACT_NAME"`
	ID_F      string `json:"ID_F"`
}

func QuerySpPtCfSaList(requestUrl string, querySpPtCfSaListFormData QuerySpPtCfSaListFormData) (querySpPtCfSaListResponse []QuerySpPtCfSaListResponse, err error) {
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
	if SpPtCfSaEnableHttpProxy {
		client = SpPtCfSaSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("isLength", strconv.Itoa(querySpPtCfSaListFormData.isLength))
	postData.Add("num_tn", strconv.Itoa(querySpPtCfSaListFormData.num_tn))
	postData.Add("standard_type", querySpPtCfSaListFormData.standard_type)
	postData.Add("keyword", querySpPtCfSaListFormData.keyword)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	querySpPtCfSaListResponse = []QuerySpPtCfSaListResponse{}
	if err != nil {
		return querySpPtCfSaListResponse, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", SpPtCfSaCookie)
	req.Header.Set("Host", "sppt.cfsa.net.cn:8086")
	req.Header.Set("Origin", "https://sppt.cfsa.net.cn:8086")
	req.Header.Set("Referer", "https://sppt.cfsa.net.cn:8086/db")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return querySpPtCfSaListResponse, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return querySpPtCfSaListResponse, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return querySpPtCfSaListResponse, err
	}
	err = json.Unmarshal(respBytes, &querySpPtCfSaListResponse)
	if err != nil {
		return querySpPtCfSaListResponse, err
	}
	return querySpPtCfSaListResponse, nil
}

func downloadSpPtCfSa(requestUrl string, downloadSpPtCfSaFormData DownloadSpPtCfSaFormData, filePath string) error {
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
	if SpPtCfSaEnableHttpProxy {
		client = SpPtCfSaSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("task", downloadSpPtCfSaFormData.task)
	postData.Add("accessData", downloadSpPtCfSaFormData.accessData)
	postData.Add("bzlb", downloadSpPtCfSaFormData.bzlb)
	postData.Add("fact_name", downloadSpPtCfSaFormData.fact_name)
	postData.Add("file_guid", downloadSpPtCfSaFormData.file_guid)
	postData.Add("keyword", downloadSpPtCfSaFormData.keyword)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", SpPtCfSaCookie)
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
