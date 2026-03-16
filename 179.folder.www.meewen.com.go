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

var FolderMeeWenEnableHttpProxy = false
var FolderMeeWenHttpProxyUrl = "111.225.152.186:8089"
var FolderMeeWenHttpProxyUrlArr = make([]string, 0)

func FolderMeeWenHttpProxy() error {
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
					FolderMeeWenHttpProxyUrlArr = append(FolderMeeWenHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					FolderMeeWenHttpProxyUrlArr = append(FolderMeeWenHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func FolderMeeWenSetHttpProxy() (httpclient *http.Client) {
	if FolderMeeWenHttpProxyUrl == "" {
		if len(FolderMeeWenHttpProxyUrlArr) <= 0 {
			err := FolderMeeWenHttpProxy()
			if err != nil {
				FolderMeeWenSetHttpProxy()
			}
		}
		FolderMeeWenHttpProxyUrl = FolderMeeWenHttpProxyUrlArr[0]
		if len(FolderMeeWenHttpProxyUrlArr) >= 2 {
			FolderMeeWenHttpProxyUrlArr = FolderMeeWenHttpProxyUrlArr[1:]
		} else {
			FolderMeeWenHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(FolderMeeWenHttpProxyUrl)
	ProxyURL, _ := url.Parse(FolderMeeWenHttpProxyUrl)
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

type FolderMeeWenCategory struct {
	Name string
	Id   string
}

type NormalFolderMeeWenListRequestPayload struct {
	Current int    `json:"current"`
	Id      string `json:"id"`
	Size    int    `json:"size"`
	Sort    int    `json:"sort"`
}

type QueryFolderMeeWenListRequestPayload struct {
	Current int    `json:"current"`
	Id      string `json:"id"`
	Size    int    `json:"size"`
	Path    string `json:"path"`
}
type QueryFolderMeeWenDetailRequestPayload struct {
	Id string `json:"id"`
}

// 下载觅文文件夹
// @Title 下载觅文文件夹
// @Description https://www.meewen.com/，下载觅文文件夹
func main() {
	folderMeeWenCategory := []FolderMeeWenCategory{
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
	for _, category := range folderMeeWenCategory {
		fmt.Println("Name = ", category.Name, " Id = ", category.Id)
		// 查询文件夹
		normalFolderListUrl := "https://www.meewen.com/meewen-portal/api/portal/document/manage/getNormalFolder"
		normalFolderMeeWenListRequestPayload := NormalFolderMeeWenListRequestPayload{
			Current: 1,
			Id:      category.Id,
			Size:    1000,
			Sort:    1,
		}
		normalFolderMeeWenListResponseData, err := NormalFolderMeeWenList(normalFolderListUrl, normalFolderMeeWenListRequestPayload)
		if err != nil {
			fmt.Println(err)
			break
		}
		for _, folder := range normalFolderMeeWenListResponseData {
			fmt.Println("Name = ", category.Name, "===", folder.Title, " Id = ", category.Id)
			pageListUrl := "https://www.meewen.com/meewen-portal/api/portal/document/manage/getPackageFile"
			page := 1
			queryFolderMeeWenListRequestPayload := QueryFolderMeeWenListRequestPayload{
				Current: page,
				Id:      folder.Id,
				Size:    1000,
				Path:    "",
			}
			queryFolderMeeWenListResponseDataData, err := QueryFolderMeeWenList(pageListUrl, queryFolderMeeWenListRequestPayload)
			if err != nil {
				fmt.Println(err)
				break
			}
			for _, data := range queryFolderMeeWenListResponseDataData {
				fmt.Println("===============开始处理数据 page = ", page, " data记录数量 = ", len(queryFolderMeeWenListResponseDataData), "==================")
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
				title = strings.ReplaceAll(title, ".docx", "")
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
				queryFolderMeeWenDetailRequestPayload := QueryFolderMeeWenDetailRequestPayload{
					Id: data.Id,
				}
				queryFolderMeeWenDetailResponseData, err := QueryFolderMeeWenDetail(detailUrl, queryFolderMeeWenDetailRequestPayload)
				fmt.Println(queryFolderMeeWenDetailResponseData)
				if err != nil {
					fmt.Println(err)
					continue
				}

				downloadUrl := queryFolderMeeWenDetailResponseData.PdfUrl
				fmt.Println(downloadUrl)

				fmt.Println("=======开始下载" + title + "========")

				err = downloadFolderMeeWen(downloadUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "www.meewen.com/www.meewen.com", "www.meewen.com/2026-03-16")
				err = copyFolderMeeWenFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")
				//DownLoadFolderMeeWenTimeSleep := 10
				DownLoadFolderMeeWenTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadFolderMeeWenTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadFolderMeeWenTimeSleep, "秒 倒计时", i, "秒===========")
				}
			}
			DownLoadFolderMeeWenPageTimeSleep := 10
			// DownLoadFolderMeeWenPageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadFolderMeeWenPageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadFolderMeeWenPageTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
	}
}

type NormalFolderMeeWenListResponse struct {
	Code    string                               `json:"code"`
	Data    []NormalFolderMeeWenListResponseData `json:"data"`
	Message string                               `json:"message"`
}
type NormalFolderMeeWenListResponseData struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func NormalFolderMeeWenList(requestUrl string, normalFolderMeeWenListRequestPayload NormalFolderMeeWenListRequestPayload) (normalFolderMeeWenListResponseData []NormalFolderMeeWenListResponseData, err error) {
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
	if FolderMeeWenEnableHttpProxy {
		client = FolderMeeWenSetHttpProxy()
	}
	// 将数据编码为JSON格式
	normalFolderMeeWenListRequestPayloadJson, err := json.Marshal(normalFolderMeeWenListRequestPayload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 创建字符串读取器，这是http.Request需要的类型
	body := bytes.NewReader(normalFolderMeeWenListRequestPayloadJson)
	req, err := http.NewRequest("POST", requestUrl, body) //建立连接

	normalFolderMeeWenListResponse := NormalFolderMeeWenListResponse{}
	if err != nil {
		return normalFolderMeeWenListResponseData, err
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
		return normalFolderMeeWenListResponseData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return normalFolderMeeWenListResponseData, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return normalFolderMeeWenListResponseData, err
	}
	err = json.Unmarshal(respBytes, &normalFolderMeeWenListResponse)
	if err != nil {
		return normalFolderMeeWenListResponseData, err
	}
	normalFolderMeeWenListResponseData = normalFolderMeeWenListResponse.Data
	return normalFolderMeeWenListResponseData, nil
}

type QueryFolderMeeWenListResponse struct {
	Code    string                            `json:"code"`
	Data    QueryFolderMeeWenListResponseData `json:"data"`
	Message string                            `json:"message"`
}
type QueryFolderMeeWenListResponseData struct {
	Data []QueryFolderMeeWenListResponseDataData `json:"data"`
}
type QueryFolderMeeWenListResponseDataData struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func QueryFolderMeeWenList(requestUrl string, queryFolderMeeWenListRequestPayload QueryFolderMeeWenListRequestPayload) (queryFolderMeeWenListResponseDataData []QueryFolderMeeWenListResponseDataData, err error) {
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
	if FolderMeeWenEnableHttpProxy {
		client = FolderMeeWenSetHttpProxy()
	}
	// 将数据编码为JSON格式
	queryFolderMeeWenListRequestPayloadJson, err := json.Marshal(queryFolderMeeWenListRequestPayload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 创建字符串读取器，这是http.Request需要的类型
	body := bytes.NewReader(queryFolderMeeWenListRequestPayloadJson)
	req, err := http.NewRequest("POST", requestUrl, body) //建立连接

	queryFolderMeeWenListResponse := QueryFolderMeeWenListResponse{}
	if err != nil {
		return queryFolderMeeWenListResponseDataData, err
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
		return queryFolderMeeWenListResponseDataData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryFolderMeeWenListResponseDataData, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryFolderMeeWenListResponseDataData, err
	}
	err = json.Unmarshal(respBytes, &queryFolderMeeWenListResponse)
	if err != nil {
		return queryFolderMeeWenListResponseDataData, err
	}
	queryFolderMeeWenListResponseDataData = queryFolderMeeWenListResponse.Data.Data
	return queryFolderMeeWenListResponseDataData, nil
}

type QueryFolderMeeWenDetailResponse struct {
	Code    string                              `json:"code"`
	Data    QueryFolderMeeWenDetailResponseData `json:"data"`
	Message string                              `json:"message"`
}
type QueryFolderMeeWenDetailResponseData struct {
	Id     string `json:"id"`
	PdfUrl string `json:"pdfurl"`
	Title  string `json:"title"`
}

func QueryFolderMeeWenDetail(requestUrl string, queryFolderMeeWenDetailRequestPayload QueryFolderMeeWenDetailRequestPayload) (queryFolderMeeWenDetailResponseData QueryFolderMeeWenDetailResponseData, err error) {
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
	if FolderMeeWenEnableHttpProxy {
		client = FolderMeeWenSetHttpProxy()
	}
	// 将数据编码为JSON格式
	queryFolderMeeWenDetailRequestPayloadJson, err := json.Marshal(queryFolderMeeWenDetailRequestPayload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 创建字符串读取器，这是http.Request需要的类型
	body := bytes.NewReader(queryFolderMeeWenDetailRequestPayloadJson)
	req, err := http.NewRequest("POST", requestUrl, body) //建立连接

	queryFolderMeeWenDetailResponse := QueryFolderMeeWenDetailResponse{}
	if err != nil {
		return queryFolderMeeWenDetailResponseData, err
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
		return queryFolderMeeWenDetailResponseData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryFolderMeeWenDetailResponseData, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryFolderMeeWenDetailResponseData, err
	}
	err = json.Unmarshal(respBytes, &queryFolderMeeWenDetailResponse)
	if err != nil {
		return queryFolderMeeWenDetailResponseData, err
	}
	queryFolderMeeWenDetailResponseData = queryFolderMeeWenDetailResponse.Data
	return queryFolderMeeWenDetailResponseData, nil
}

func downloadFolderMeeWen(attachmentUrl string, filePath string) error {
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
	if FolderMeeWenEnableHttpProxy {
		client = FolderMeeWenSetHttpProxy()
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
	req.Header.Set("Sec-Fetch-Dest", "folder")
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

func copyFolderMeeWenFile(src, dst string) (err error) {
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
