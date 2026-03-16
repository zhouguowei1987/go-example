package main

import (
	"bytes"
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

var DocumentMeeWenEnableHttpProxy = false
var DocumentMeeWenHttpProxyUrl = "111.225.152.186:8089"
var DocumentMeeWenHttpProxyUrlArr = make([]string, 0)

func DocumentMeeWenHttpProxy() error {
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
					DocumentMeeWenHttpProxyUrlArr = append(DocumentMeeWenHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					DocumentMeeWenHttpProxyUrlArr = append(DocumentMeeWenHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func DocumentMeeWenSetHttpProxy() (httpclient *http.Client) {
	if DocumentMeeWenHttpProxyUrl == "" {
		if len(DocumentMeeWenHttpProxyUrlArr) <= 0 {
			err := DocumentMeeWenHttpProxy()
			if err != nil {
				DocumentMeeWenSetHttpProxy()
			}
		}
		DocumentMeeWenHttpProxyUrl = DocumentMeeWenHttpProxyUrlArr[0]
		if len(DocumentMeeWenHttpProxyUrlArr) >= 2 {
			DocumentMeeWenHttpProxyUrlArr = DocumentMeeWenHttpProxyUrlArr[1:]
		} else {
			DocumentMeeWenHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(DocumentMeeWenHttpProxyUrl)
	ProxyURL, _ := url.Parse(DocumentMeeWenHttpProxyUrl)
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
			ResponseHeaderTimeout: time.Second * 30,
		},
	}
	return httpclient
}

type DocumentMeeWenCategory struct {
	Name string
	Id   string
}

type QueryDocumentMeeWenListRequestPayload struct {
	Current int    `json:"current"`
	Id      string `json:"id"`
	Size    int    `json:"size"`
	Sort    int    `json:"sort"`
}
type QueryDocumentMeeWenDetailRequestPayload struct {
	Id string `json:"id"`
}

// 下载觅文普通文档
// @Title 下载觅文普通文档
// @Description https://www.meewen.com/，下载觅文普通文档
func main() {
	documentMeeWenCategory := []DocumentMeeWenCategory{
		{
			Name: "领导讲话",
			Id:   "4383e2d7-1c25-11f0-96b3-6c1ff709ec87",
		},
		{
			Name: "专题讲稿",
			Id:   "4383d175-1c25-11f0-96b3-6c1ff709ec87",
		},
		{
			Name: "遴选题库",
			Id:   "4383e29f-1c25-11f0-96b3-6c1ff709ec87",
		},
		{
			Name: "表格合同",
			Id:   "4383e241-1c25-11f0-96b3-6c1ff709ec87",
		},
		{
			Name: "公考素材",
			Id:   "4383e1f8-1c25-11f0-96b3-6c1ff709ec87",
		},
	}
	for _, category := range documentMeeWenCategory {
		fmt.Println("Name = ", category.Name, " Id = ", category.Id)
		pageListUrl := "https://www.meewen.com/meewen-portal/api/portal/document/manage/getNormalDocument"
		page := 1
		queryDocumentMeeWenListRequestPayload := QueryDocumentMeeWenListRequestPayload{
			Current: page,
			Id:      category.Id,
			Size:    1000,
			Sort:    1,
		}
		queryDocumentMeeWenListResponseDataData, err := QueryDocumentMeeWenList(pageListUrl, queryDocumentMeeWenListRequestPayload)
		if err != nil {
			fmt.Println(err)
			break
		}
		for _, data := range queryDocumentMeeWenListResponseDataData {
			fmt.Println("===============开始处理数据 page = ", page, " data记录数量 = ", len(queryDocumentMeeWenListResponseDataData), "==================")
			fmt.Println(data.Id)

			title := data.Title
			fmt.Println(data.Title)
			if strings.Index(data.Title, "doc") == -1 && strings.Index(data.Title, "pdf") == -1 {
				fmt.Println("文档不是doc、pdf文档，跳过")
				continue
			}
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "--", "-")
			title = strings.ReplaceAll(title, ".docx", "")
			title = strings.ReplaceAll(title, ".doc", "")
			title = strings.ReplaceAll(title, ".pdf", "")

			filePath := "../www.meewen.com/www.meewen.com/" + category.Name + "/" + title + ".pdf"
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载========")
			detailUrl := "https://www.meewen.com/meewen-portal/api/portal/custom/file/document/getDetail"
			queryDocumentMeeWenDetailRequestPayload := QueryDocumentMeeWenDetailRequestPayload{
				Id: data.Id,
			}
			queryDocumentMeeWenDetailResponseData, err := QueryDocumentMeeWenDetail(detailUrl, queryDocumentMeeWenDetailRequestPayload)
			fmt.Println(queryDocumentMeeWenDetailResponseData)
			if err != nil {
				fmt.Println(err)
				continue
			}

			downloadUrl := queryDocumentMeeWenDetailResponseData.PdfUrl
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")

			err = downloadDocumentMeeWen(downloadUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "www.meewen.com/www.meewen.com", "www.meewen.com/2026-03-16")
			err = copyDocumentMeeWenFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			//DownLoadDocumentMeeWenTimeSleep := 10
			DownLoadDocumentMeeWenTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadDocumentMeeWenTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadDocumentMeeWenTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadDocumentMeeWenPageTimeSleep := 10
		// DownLoadDocumentMeeWenPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadDocumentMeeWenPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadDocumentMeeWenPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
	}
}

type QueryDocumentMeeWenListResponse struct {
	Code    string                              `json:"code"`
	Data    QueryDocumentMeeWenListResponseData `json:"data"`
	Message string                              `json:"message"`
}
type QueryDocumentMeeWenListResponseData struct {
	Data []QueryDocumentMeeWenListResponseDataData `json:"data"`
}
type QueryDocumentMeeWenListResponseDataData struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func QueryDocumentMeeWenList(requestUrl string, queryDocumentMeeWenListRequestPayload QueryDocumentMeeWenListRequestPayload) (queryDocumentMeeWenListResponseDataData []QueryDocumentMeeWenListResponseDataData, err error) {
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
	if DocumentMeeWenEnableHttpProxy {
		client = DocumentMeeWenSetHttpProxy()
	}
	// 将数据编码为JSON格式
	queryDocumentMeeWenListRequestPayloadJson, err := json.Marshal(queryDocumentMeeWenListRequestPayload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 创建字符串读取器，这是http.Request需要的类型
	body := bytes.NewReader(queryDocumentMeeWenListRequestPayloadJson)
	req, err := http.NewRequest("POST", requestUrl, body) //建立连接

	queryDocumentMeeWenListResponse := QueryDocumentMeeWenListResponse{}
	if err != nil {
		return queryDocumentMeeWenListResponseDataData, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Host", "www.meewen.com")
	req.Header.Set("Origin", "https://www.meewen.com")
	req.Header.Set("Referer", "https://www.meewen.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryDocumentMeeWenListResponseDataData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryDocumentMeeWenListResponseDataData, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryDocumentMeeWenListResponseDataData, err
	}
	err = json.Unmarshal(respBytes, &queryDocumentMeeWenListResponse)
	if err != nil {
		return queryDocumentMeeWenListResponseDataData, err
	}
	queryDocumentMeeWenListResponseDataData = queryDocumentMeeWenListResponse.Data.Data
	return queryDocumentMeeWenListResponseDataData, nil
}

type QueryDocumentMeeWenDetailResponse struct {
	Code    string                                `json:"code"`
	Data    QueryDocumentMeeWenDetailResponseData `json:"data"`
	Message string                                `json:"message"`
}
type QueryDocumentMeeWenDetailResponseData struct {
	Id     string `json:"id"`
	PdfUrl string `json:"pdfurl"`
	Title  string `json:"title"`
}

func QueryDocumentMeeWenDetail(requestUrl string, queryDocumentMeeWenDetailRequestPayload QueryDocumentMeeWenDetailRequestPayload) (queryDocumentMeeWenDetailResponseData QueryDocumentMeeWenDetailResponseData, err error) {
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
	if DocumentMeeWenEnableHttpProxy {
		client = DocumentMeeWenSetHttpProxy()
	}
	// 将数据编码为JSON格式
	queryDocumentMeeWenDetailRequestPayloadJson, err := json.Marshal(queryDocumentMeeWenDetailRequestPayload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 创建字符串读取器，这是http.Request需要的类型
	body := bytes.NewReader(queryDocumentMeeWenDetailRequestPayloadJson)
	req, err := http.NewRequest("POST", requestUrl, body) //建立连接

	queryDocumentMeeWenDetailResponse := QueryDocumentMeeWenDetailResponse{}
	if err != nil {
		return queryDocumentMeeWenDetailResponseData, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Host", "www.meewen.com")
	req.Header.Set("Origin", "https://www.meewen.com")
	req.Header.Set("Referer", "https://www.meewen.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryDocumentMeeWenDetailResponseData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryDocumentMeeWenDetailResponseData, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryDocumentMeeWenDetailResponseData, err
	}
	err = json.Unmarshal(respBytes, &queryDocumentMeeWenDetailResponse)
	if err != nil {
		return queryDocumentMeeWenDetailResponseData, err
	}
	queryDocumentMeeWenDetailResponseData = queryDocumentMeeWenDetailResponse.Data
	return queryDocumentMeeWenDetailResponseData, nil
}

func downloadDocumentMeeWen(attachmentUrl string, filePath string) error {
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
	if DocumentMeeWenEnableHttpProxy {
		client = DocumentMeeWenSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("Host", "www.meewen.com")
	req.Header.Set("Referer", "https://www.meewen.com/")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
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

func copyDocumentMeeWenFile(src, dst string) (err error) {
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
