package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	// "math/rand"
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

var LvPin100EnableHttpProxy = false
var LvPin100HttpProxyUrl = "111.225.152.186:8089"
var LvPin100HttpProxyUrlArr = make([]string, 0)

func LvPin100HttpProxy() error {
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
					LvPin100HttpProxyUrlArr = append(LvPin100HttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					LvPin100HttpProxyUrlArr = append(LvPin100HttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func LvPin100SetHttpProxy() (httpclient *http.Client) {
	if LvPin100HttpProxyUrl == "" {
		if len(LvPin100HttpProxyUrlArr) <= 0 {
			err := LvPin100HttpProxy()
			if err != nil {
				LvPin100SetHttpProxy()
			}
		}
		LvPin100HttpProxyUrl = LvPin100HttpProxyUrlArr[0]
		if len(LvPin100HttpProxyUrlArr) >= 2 {
			LvPin100HttpProxyUrlArr = LvPin100HttpProxyUrlArr[1:]
		} else {
			LvPin100HttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(LvPin100HttpProxyUrl)
	ProxyURL, _ := url.Parse(LvPin100HttpProxyUrl)
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

// 下载云南法网合同模板文档
// @Title 下载云南法网合同模板文档
// @Description https://yn.ai.lvpin100.com/，下载云南法网合同模板文档
func main() {
	// 第一步：获取文书模版相关类别
	// categoriesUrl := "https://front.ai.lvpin100.com/api/speed-front/doc/categories?appId=ED4E1B16C6FE657B804981B9EBC8465D&type=3&deviceType=1"
	// queryLvPin100CategoriesListResponseDataList, err := QueryLvPin100CategoriesList(categoriesUrl)
	// if err != nil {
	// 	LvPin100HttpProxyUrl = ""
	// 	fmt.Println(err)
	// }
	queryLvPin100CategoriesListResponseDataList := []QueryLvPin100CategoriesListResponseDataList{
		{
			Id:    "",
			Title: "全部",
		},
	}
	if len(queryLvPin100CategoriesListResponseDataList) <= 0 {
		fmt.Println("获取文书模版相关类别失败")
	}
	// 第二部：遍历分类，获取分类下的合同模板
	for _, category := range queryLvPin100CategoriesListResponseDataList {
		fmt.Println("=====================分类信息 id = ", category.Id, ",title = ", category.Title, "=========================")

		// 一次性获取分类下所有模板
		pageListUrl := fmt.Sprintf("https://front.ai.lvpin100.com/api/speed-front/doc/list?pageNo=1&pageSize=9999&categoryId=%s&appId=ED4E1B16C6FE657B804981B9EBC8465D&type=3&deviceType=1", category.Id)
		if len(category.Id) <= 0 {
			pageListUrl = "https://front.ai.lvpin100.com/api/speed-front/doc/list?pageNo=1&pageSize=9999&appId=ED4E1B16C6FE657B804981B9EBC8465D&type=3&deviceType=1"
		}
		fmt.Println(pageListUrl)
		queryLvPin100ListResponseDataResultList, err := QueryLvPin100List(pageListUrl)
		if err != nil {
			LvPin100HttpProxyUrl = ""
			fmt.Println(err)
		}
		for _, lvPin100 := range queryLvPin100ListResponseDataResultList {
			fmt.Println("=====================分类信息 id = ", category.Id, ",title = ", category.Title, "=========================")

			title := lvPin100.Title
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, "-", "")
			title = strings.ReplaceAll(title, " ", "")
			title = strings.ReplaceAll(title, "|", "-")
			title = strings.ReplaceAll(title, "/", "-")
			fmt.Println(title)
			fmt.Println(lvPin100.Id)

			detailUrl := fmt.Sprintf("https://yn.ai.lvpin100.com/pc/docdetail/%s", lvPin100.Id)
			fmt.Println(detailUrl)

			filePath := "../yn.ai.lvpin100.com/yn.ai.lvpin100.com/" + title
			fmt.Println(filePath)
			_, errPdf := os.Stat(filePath + ".pdf")
			_, errDoc := os.Stat(filePath + ".doc")
			_, errDocx := os.Stat(filePath + ".docx")
			_, errXls := os.Stat(filePath + ".xls")
			_, errXlsx := os.Stat(filePath + ".xlsx")
			if errPdf == nil || errDoc == nil || errDocx == nil || errXls == nil || errXlsx == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载========")

			downloadUrl := fmt.Sprintf("https://front.ai.lvpin100.com/api/speed-front/doc/download?appId=ED4E1B16C6FE657B804981B9EBC8465D&userId=87E208CA66FB40579A837F2FE6A43BB3&id=%s&deviceType=1", lvPin100.Id)
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")
			fileExtension, err := downloadLvPin100(downloadUrl, detailUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "yn.ai.lvpin100.com/yn.ai.lvpin100.com", "yn.ai.lvpin100.com/temp-yn.ai.lvpin100.com")
			err = copyLvPin100File(filePath+fileExtension, tempFilePath+fileExtension)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			// DownLoadLvPin100TimeSleep := 5
			// DownLoadLvPin100TimeSleep := rand.Intn(5)
			// for i := 1; i <= DownLoadLvPin100TimeSleep; i++ {
			// 	time.Sleep(time.Second)
			// 	fmt.Println("category="+category.Title+",filePath="+filePath+"===========下载成功 暂停", DownLoadLvPin100TimeSleep, "秒 倒计时", i, "秒===========")
			// }
		}
		DownLoadLvPin100PageTimeSleep := 5
		// DownLoadLvPin100PageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadLvPin100PageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("category="+category.Title+"=========== 暂停", DownLoadLvPin100PageTimeSleep, "秒 倒计时", i, "秒===========")
		}
	}
}

type QueryLvPin100CategoriesListResponse struct {
	Data []QueryLvPin100CategoriesListResponseDataList `json:"data"`
}

type QueryLvPin100CategoriesListResponseDataList struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func QueryLvPin100CategoriesList(requestUrl string) (queryLvPin100CategoriesListResponseDataList []QueryLvPin100CategoriesListResponseDataList, err error) {
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
	if LvPin100EnableHttpProxy {
		client = LvPin100SetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	queryLvPin100CategoriesListResponse := QueryLvPin100CategoriesListResponse{}
	if err != nil {
		return queryLvPin100CategoriesListResponseDataList, err
	}

	req.Header.Set("authority", "front.ai.lvpin100.com")
	req.Header.Set("method", "GET")
	req.Header.Set("path", "/api/speed-front/doc/categories?appId=ED4E1B16C6FE657B804981B9EBC8465D&type=3&deviceType=1")
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Host", "yn.ai.lvpin100.com")
	req.Header.Set("Origin", "https://yn.ai.lvpin100.com")
	req.Header.Set("Referer", "https://yn.ai.lvpin100.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryLvPin100CategoriesListResponseDataList, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryLvPin100CategoriesListResponseDataList, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryLvPin100CategoriesListResponseDataList, err
	}
	err = json.Unmarshal(respBytes, &queryLvPin100CategoriesListResponse)
	if err != nil {
		return queryLvPin100CategoriesListResponseDataList, err
	}
	queryLvPin100CategoriesListResponseDataList = queryLvPin100CategoriesListResponse.Data
	return queryLvPin100CategoriesListResponseDataList, nil
}

type QueryLvPin100ListResponse struct {
	Data QueryLvPin100ListResponseData `json:data`
}

type QueryLvPin100ListResponseData struct {
	Result []QueryLvPin100ListResponseDataResultList `json:result`
}

type QueryLvPin100ListResponseDataResultList struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func QueryLvPin100List(requestUrl string) (queryLvPin100ListResponseDataResultList []QueryLvPin100ListResponseDataResultList, err error) {
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
	if LvPin100EnableHttpProxy {
		client = LvPin100SetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	queryLvPin100ListResponse := QueryLvPin100ListResponse{}
	if err != nil {
		return queryLvPin100ListResponseDataResultList, err
	}

	req.Header.Set("authority", "front.ai.lvpin100.com")
	req.Header.Set("method", "GET")
	req.Header.Set("path", "/api/speed-front/doc/list?pageNo=1&pageSize=15&categoryId=F11BC3BB2ECE423ABD9999C8E3010CCF&appId=ED4E1B16C6FE657B804981B9EBC8465D&type=3&deviceType=1")
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Host", "yn.ai.lvpin100.com")
	req.Header.Set("Origin", "https://yn.ai.lvpin100.com")
	req.Header.Set("Referer", "https://yn.ai.lvpin100.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryLvPin100ListResponseDataResultList, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryLvPin100ListResponseDataResultList, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryLvPin100ListResponseDataResultList, err
	}
	err = json.Unmarshal(respBytes, &queryLvPin100ListResponse)
	if err != nil {
		return queryLvPin100ListResponseDataResultList, err
	}
	queryLvPin100ListResponseDataResultList = queryLvPin100ListResponse.Data.Result
	return queryLvPin100ListResponseDataResultList, nil
}

func downloadLvPin100(attachmentUrl string, referer string, filePath string) (fileExtension string, err error) {
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
	if LvPin100EnableHttpProxy {
		client = LvPin100SetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return fileExtension, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "yn.ai.lvpin100.com")
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
		return fileExtension, err
	}
	defer resp.Body.Close()
	// 如果访问失败 就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return fileExtension, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}

	contenttype := resp.Header.Get("Content-Type")
	fmt.Println(contenttype)
	if contenttype == "application/pdf" {
		fileExtension = ".pdf"
	} else if contenttype == "application/msword" {
		fileExtension = ".doc"
	} else if contenttype == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" {
		fileExtension = ".docx"
	} else if contenttype == "application/vnd.ms-excel" {
		fileExtension = ".xls"
	} else if contenttype == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		fileExtension = ".xlsx"
	}
	if len(fileExtension) <= 0 {
		return fileExtension, errors.New("不是pdf、doc、docx文件类型")
	}
	filePath += fileExtension

	// 创建一个文件用于保存
	fileDiv := filepath.Dir(filePath)
	if _, err = os.Stat(fileDiv); err != nil {
		if os.MkdirAll(fileDiv, 0o777) != nil {
			return fileExtension, err
		}
	}
	out, err := os.Create(filePath)
	if err != nil {
		return fileExtension, err
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fileExtension, err
	}
	return fileExtension, nil
}

func copyLvPin100File(src, dst string) (err error) {
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
