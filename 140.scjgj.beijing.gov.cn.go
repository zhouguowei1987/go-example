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

var ScJgjEnableHttpProxy = false
var ScJgjHttpProxyUrl = "111.225.152.186:8089"
var ScJgjHttpProxyUrlArr = make([]string, 0)

func ScJgjHttpProxy() error {
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
					ScJgjHttpProxyUrlArr = append(ScJgjHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					ScJgjHttpProxyUrlArr = append(ScJgjHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func ScJgjSetHttpProxy() (httpclient *http.Client) {
	if ScJgjHttpProxyUrl == "" {
		if len(ScJgjHttpProxyUrlArr) <= 0 {
			err := ScJgjHttpProxy()
			if err != nil {
				ScJgjSetHttpProxy()
			}
		}
		ScJgjHttpProxyUrl = ScJgjHttpProxyUrlArr[0]
		if len(ScJgjHttpProxyUrlArr) >= 2 {
			ScJgjHttpProxyUrlArr = ScJgjHttpProxyUrlArr[1:]
		} else {
			ScJgjHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(ScJgjHttpProxyUrl)
	ProxyURL, _ := url.Parse(ScJgjHttpProxyUrl)
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

type QueryScJgjListFormData struct {
	page         int
	standardCode string
	namecn       string
	status       int
}

// 下载北京市市场监督局文档
// @Title 下载北京市市场监督局文档
// @Description https://scjgj.beijing.gov.cn/cxfw/，下载北京市市场监督局文档
func main() {
	pageListUrl := "https://cx.scjgj.beijing.gov.cn/shiyao/nosession/showall/bzh_api_standard"
	fmt.Println(pageListUrl)
	page := 1
	maxPage := 234
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		queryScJgjListFormData := QueryScJgjListFormData{
			page:         page,
			standardCode: "",
			namecn:       "",
			status:       0,
		}
		queryScJgjListResponseContent, err := QueryScJgjList(pageListUrl, queryScJgjListFormData)
		if err != nil {
			ScJgjHttpProxyUrl = ""
			fmt.Println(err)
		}
		for _, scJgj := range queryScJgjListResponseContent {
			fmt.Println("=====================开始处理数据 page = ", page, "=========================")
			code := scJgj.StandardCode
			code = strings.ReplaceAll(code, "/", "-")
			code = strings.ReplaceAll(code, "—", "-")
			fmt.Println(code)

			title := scJgj.NameCn
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, "-", "")
			title = strings.ReplaceAll(title, " ", "")
			title = strings.ReplaceAll(title, "|", "-")
			fmt.Println(title)

			detailUrl := fmt.Sprintf("https://scjgj.beijing.gov.cn/cxfw/201911/t20191118_513430.html?name=bzh_api_standard&id=%d", scJgj.Id)
			fmt.Println(detailUrl)

			filePath := "../scjgj.beijing.gov.cn/" + title + "(" + code + ")" + ".pdf"
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载========")

			downloadUrl := scJgj.StaUrl
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")
			err = downloadScJgj(downloadUrl, detailUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "../scjgj.beijing.gov.cn", "../upload.doc88.com/scjgj.beijing.gov.cn")
			err = copyScJgjFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			//DownLoadScJgjTimeSleep := 10
			DownLoadScJgjTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadScJgjTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadScJgjTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadScJgjPageTimeSleep := 10
		// DownLoadScJgjPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadScJgjPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadScJgjPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

type QueryScJgjListResponse struct {
	Content          []QueryScJgjListResponseContent `json:"content"`
	Number           int                             `json:"number"`
	NumberOfElements int                             `json:"numberOfElements"`
	QuerySql         string                          `json:"querySql"`
	Size             int                             `json:"size"`
	TotalElements    int                             `json:"totalElements"`
	TotalPages       int                             `json:"totalPages"`
}

type QueryScJgjListResponseContent struct {
	Id           int    `json:"id"`
	NameCn       string `json:"namecn"`
	StaUrl       string `json:"sta_url"`
	StandardCode string `json:"standardCode"`
}

func QueryScJgjList(requestUrl string, queryScJgjListFormData QueryScJgjListFormData) (queryScJgjListResponseContent []QueryScJgjListResponseContent, err error) {
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
	if ScJgjEnableHttpProxy {
		client = ScJgjSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("page", strconv.Itoa(queryScJgjListFormData.page))
	postData.Add("standardCode", queryScJgjListFormData.standardCode)
	postData.Add("namecn", queryScJgjListFormData.namecn)
	postData.Add("status", strconv.Itoa(queryScJgjListFormData.status))
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	queryScJgjListResponse := QueryScJgjListResponse{}
	if err != nil {
		return queryScJgjListResponseContent, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Host", "cx.scjgj.beijing.gov.cn")
	req.Header.Set("Origin", "https://scjgj.beijing.gov.cn")
	req.Header.Set("Referer", "https://scjgj.beijing.gov.cn/cxfw/?serviceName=bzh_api_standard")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryScJgjListResponseContent, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryScJgjListResponseContent, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryScJgjListResponseContent, err
	}
	err = json.Unmarshal(respBytes, &queryScJgjListResponse)
	if err != nil {
		return queryScJgjListResponseContent, err
	}
	queryScJgjListResponseContent = queryScJgjListResponse.Content
	return queryScJgjListResponseContent, nil
}

func downloadScJgj(attachmentUrl string, referer string, filePath string) error {
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
	if ScJgjEnableHttpProxy {
		client = ScJgjSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "cx.scjgj.beijing.gov.cn")
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

func copyScJgjFile(src, dst string) (err error) {
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
