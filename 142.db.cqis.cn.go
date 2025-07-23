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

var CqIsEnableHttpProxy = false
var CqIsHttpProxyUrl = "111.225.152.186:8089"
var CqIsHttpProxyUrlArr = make([]string, 0)

func CqIsHttpProxy() error {
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
					CqIsHttpProxyUrlArr = append(CqIsHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					CqIsHttpProxyUrlArr = append(CqIsHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func CqIsSetHttpProxy() (httpclient *http.Client) {
	if CqIsHttpProxyUrl == "" {
		if len(CqIsHttpProxyUrlArr) <= 0 {
			err := CqIsHttpProxy()
			if err != nil {
				CqIsSetHttpProxy()
			}
		}
		CqIsHttpProxyUrl = CqIsHttpProxyUrlArr[0]
		if len(CqIsHttpProxyUrlArr) >= 2 {
			CqIsHttpProxyUrlArr = CqIsHttpProxyUrlArr[1:]
		} else {
			CqIsHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(CqIsHttpProxyUrl)
	ProxyURL, _ := url.Parse(CqIsHttpProxyUrl)
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

type QueryCqIsListFormData struct {
	page      int
	rows      int
	queryJson QueryJson
}

type QueryJson struct {
	FinishCode             string `json:"FinishCode"`
	SourceName             string `json:"SourceName"`
	BZName                 string `json:"BZName"`
	ClassificationIndustry string `json:"ClassificationIndustry"`
	StandardCategory       string `json:"StandardCategory"`
	OtherCategory          string `json:"OtherCategory"`
}

var CqIsCookie = "__RequestVerificationToken=igF-mWhmVgvlECbsH8Lx3A7Kx9Zmv2zCM98q8cltHjNFPD2KOqFWd2jyRHc05oWowSVaK1oR1ksEYz6nG6c54Wms22jwhnO-YnbguZ9NeoQ1; userbrowse=4f074005-2472-97f0-f88e-0ad97d4b09a5"

// 下载重庆市地方标准文档
// @Title 下载重庆市地方标准文档
// @Description http://db.cqis.cn/LocalStandard/Index/，下载重庆市地方标准文档
func main() {
	pageListUrl := "http://db.cqis.cn/LocalStandard/GetDB50ShareList"
	fmt.Println(pageListUrl)
	page := 1
	maxPage := 136
	rows := 15
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		queryCqIsListFormData := QueryCqIsListFormData{
			page: page,
			rows: rows,
			queryJson: QueryJson{
				FinishCode:             "",
				SourceName:             "重庆",
				BZName:                 "",
				ClassificationIndustry: "",
				StandardCategory:       "",
				OtherCategory:          "",
			},
		}
		queryCqIsListResponseRows, err := QueryCqIsList(pageListUrl, queryCqIsListFormData)
		if err != nil {
			CqIsHttpProxyUrl = ""
			fmt.Println(err)
		}
		for _, cqIs := range queryCqIsListResponseRows {
			fmt.Println("=====================开始处理数据 page = ", page, "=========================")
			code := cqIs.Project_FinishCODE
			code = strings.ReplaceAll(code, "/", "-")
			code = strings.ReplaceAll(code, "—", "-")
			fmt.Println(code)

			title := cqIs.BZName
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, "-", "")
			title = strings.ReplaceAll(title, " ", "")
			title = strings.ReplaceAll(title, "|", "-")
			fmt.Println(title)

			filePath := "../db.cqis.cn/" + title + "(" + code + ")" + ".pdf"
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载========")

			downloadUrl := fmt.Sprintf("http://db.cqis.cn/LocalStandard/PdfShowPage?keyValue=%s", cqIs.BZProjectPublicId)
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")

			detailUrl := fmt.Sprintf("http://db.cqis.cn/LocalStandard/ShowPdf?keyValue=%s", cqIs.Project_FinishCODE)
			err = downloadCqIs(downloadUrl, detailUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "../db.cqis.cn", "../upload.doc88.com/db.cqis.cn")
			err = copyCqIsFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			//DownLoadCqIsTimeSleep := 10
			DownLoadCqIsTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadCqIsTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadCqIsTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadCqIsPageTimeSleep := 10
		// DownLoadCqIsPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadCqIsPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadCqIsPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

type QueryCqIsListResponse struct {
	Rows    []QueryCqIsListResponseRows `json:"rows"`
	Page    int                         `json:"page"`
	Records int                         `json:"records"`
	Total   int                         `json:"total"`
}

type QueryCqIsListResponseRows struct {
	BZProjectPublicId  string `json:"BZProjectPublicId"`
	BZName             string `json:"BZName"`
	Project_FinishCODE string `json:"Project_FinishCODE"`
}

func QueryCqIsList(requestUrl string, queryCqIsListFormData QueryCqIsListFormData) (queryCqIsListResponseRows []QueryCqIsListResponseRows, err error) {
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
	if CqIsEnableHttpProxy {
		client = CqIsSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("page", strconv.Itoa(queryCqIsListFormData.page))
	postData.Add("rows", strconv.Itoa(queryCqIsListFormData.rows))
	queryJson, err := json.Marshal(queryCqIsListFormData.queryJson)
	postData.Add("queryJson", string(queryJson))
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	queryCqIsListResponse := QueryCqIsListResponse{}
	if err != nil {
		return queryCqIsListResponseRows, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", CqIsCookie)
	req.Header.Set("Host", "db.cqis.cn")
	req.Header.Set("Origin", "http://db.cqis.cn")
	req.Header.Set("Referer", "http://db.cqis.cn/LocalStandard/Index")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryCqIsListResponseRows, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryCqIsListResponseRows, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryCqIsListResponseRows, err
	}
	err = json.Unmarshal(respBytes, &queryCqIsListResponse)
	if err != nil {
		return queryCqIsListResponseRows, err
	}
	queryCqIsListResponseRows = queryCqIsListResponse.Rows
	return queryCqIsListResponseRows, nil
}

func downloadCqIs(attachmentUrl string, referer string, filePath string) error {
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
	if CqIsEnableHttpProxy {
		client = CqIsSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", CqIsCookie)
	req.Header.Set("Host", "db.cqis.cn")
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

func copyCqIsFile(src, dst string) (err error) {
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
