package main

import (
	// "bytes"
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

var GuoJiaJiaoYuEnableHttpProxy = false
var GuoJiaJiaoYuHttpProxyUrl = "111.225.152.186:8089"
var GuoJiaJiaoYuHttpProxyUrlArr = make([]string, 0)

func GuoJiaJiaoYuHttpProxy() error {
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
					GuoJiaJiaoYuHttpProxyUrlArr = append(GuoJiaJiaoYuHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					GuoJiaJiaoYuHttpProxyUrlArr = append(GuoJiaJiaoYuHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func GuoJiaJiaoYuSetHttpProxy() (httpclient *http.Client) {
	if GuoJiaJiaoYuHttpProxyUrl == "" {
		if len(GuoJiaJiaoYuHttpProxyUrlArr) <= 0 {
			err := GuoJiaJiaoYuHttpProxy()
			if err != nil {
				GuoJiaJiaoYuSetHttpProxy()
			}
		}
		GuoJiaJiaoYuHttpProxyUrl = GuoJiaJiaoYuHttpProxyUrlArr[0]
		if len(GuoJiaJiaoYuHttpProxyUrlArr) >= 2 {
			GuoJiaJiaoYuHttpProxyUrlArr = GuoJiaJiaoYuHttpProxyUrlArr[1:]
		} else {
			GuoJiaJiaoYuHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(GuoJiaJiaoYuHttpProxyUrl)
	ProxyURL, _ := url.Parse(GuoJiaJiaoYuHttpProxyUrl)
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

// 下载教育中国文档
// @Title 下载教育中国文档
// @Description http://www.guojiajiaoyu.com/，下载教育中国文档
func main() {
	var classNames = [...]string{"学前", "小学", "初中", "高中"}
	for _, className := range classNames {
		limit := 100
		offset := 0
		isPageListGo := true
		for isPageListGo {
			pageListUrl := fmt.Sprintf("http://www.guojiajiaoyu.com/home/works/model?sortId=&state=1&classname=%s&category=&activityClass=&keyWords=&limit=%d&offset=%d", url.QueryEscape(className), limit, offset)
			pageListReferer := fmt.Sprintf("http://www.guojiajiaoyu.com/home/category?classname=%s", className)
			queryGuoJiaJiaoYuListResponseRows, err := QueryGuoJiaJiaoYuList(pageListUrl, pageListReferer)
			if err != nil {
				fmt.Println(err)
				isPageListGo = false
				break
			}
			if len(queryGuoJiaJiaoYuListResponseRows) <= 0 {
				isPageListGo = false
				break
			}
			for _, row := range queryGuoJiaJiaoYuListResponseRows {
				fmt.Println("===============开始处理数据 offset = ", offset, " data记录数量 = ", len(queryGuoJiaJiaoYuListResponseRows), "==================")

				if len(row.AnnexUrl) <= 0 {
					fmt.Println("文档不存在，跳过")
					continue
				}

				// 查看文档后缀
				fileExt := filepath.Ext(row.AnnexUrl)
				if strings.Index(fileExt, "doc") == -1 {
					fmt.Println("文档不是doc文档，跳过")
					continue
				}
				fmt.Println(fileExt)

				if len(row.Works) <= 0 {
					fmt.Println("标题不存在，跳过")
					continue
				}

				title := row.Works
				if len(row.ActivityClass) > 0 {
					title = title + "-" + row.ActivityClass
				}
				if len(row.KeyWords) > 0 {
					title = title + "-" + row.KeyWords
				}
				if len(row.Category) > 0 {
					title = title + "-" + row.Category
				}

				title = title + fileExt
				title = strings.TrimSpace(title)
				title = strings.ReplaceAll(title, " ", "-")
				title = strings.ReplaceAll(title, "　", "-")
				title = strings.ReplaceAll(title, "：", "-")
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, "--", "-")
				fmt.Println(title)

				filePath := "../www.guojiajiaoyu.com/www.guojiajiaoyu.com/" + className + "/" + title
				fmt.Println(filePath)

				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				fmt.Println("=======开始下载========")
				downloadUrl := fmt.Sprintf("http://www.guojiajiaoyu.com%s", row.AnnexUrl)
				fmt.Println(downloadUrl)

				fmt.Println("=======开始下载" + title + "========")

				err = downloadGuoJiaJiaoYu(downloadUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				// 查看文件大小，如果是空文件，则删除
				fileInfo, err := os.Stat(filePath)
				if err == nil && fileInfo.Size() == 0 {
					fmt.Println("空文件删除")
					err = os.Remove(filePath)
					isPageListGo = false
					break
				}
				if err != nil {
					continue
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "www.guojiajiaoyu.com/www.guojiajiaoyu.com/"+className, "www.guojiajiaoyu.com/temp-www.guojiajiaoyu.com/"+className)
				err = copyGuoJiaJiaoYuFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")
				//DownLoadGuoJiaJiaoYuTimeSleep := 10
				DownLoadGuoJiaJiaoYuTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadGuoJiaJiaoYuTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("offset="+strconv.Itoa(offset)+",filePath="+filePath+"===========下载成功 暂停", DownLoadGuoJiaJiaoYuTimeSleep, "秒 倒计时", i, "秒===========")
				}
			}

			offset = offset + limit
			DownLoadGuoJiaJiaoYuPageTimeSleep := 10
			// DownLoadGuoJiaJiaoYuPageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadGuoJiaJiaoYuPageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("offset="+strconv.Itoa(offset)+"=========== 暂停", DownLoadGuoJiaJiaoYuPageTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
	}
}

type QueryGuoJiaJiaoYuListResponse struct {
	Rows  []QueryGuoJiaJiaoYuListResponseRows `json:"rows"`
	Total int                                 `json:"total"`
}
type QueryGuoJiaJiaoYuListResponseRows struct {
	ActivityClass string `json:"activityClass"`
	Category      string `json:"category"`
	KeyWords      string `json:"keyWords"`
	Works         string `json:"works"`
	AnnexUrl      string `json:"annexUrl"`
}

func QueryGuoJiaJiaoYuList(requestUrl string, referer string) (queryGuoJiaJiaoYuListResponseRows []QueryGuoJiaJiaoYuListResponseRows, err error) {
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
	if GuoJiaJiaoYuEnableHttpProxy {
		client = GuoJiaJiaoYuSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	queryGuoJiaJiaoYuListResponse := QueryGuoJiaJiaoYuListResponse{}
	if err != nil {
		return queryGuoJiaJiaoYuListResponseRows, err
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.guojiajiaoyu.com")
	req.Header.Set("Origin", "http://www.guojiajiaoyu.com")
	req.Header.Set("Referer", referer)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryGuoJiaJiaoYuListResponseRows, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryGuoJiaJiaoYuListResponseRows, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryGuoJiaJiaoYuListResponseRows, err
	}
	err = json.Unmarshal(respBytes, &queryGuoJiaJiaoYuListResponse)
	if err != nil {
		return queryGuoJiaJiaoYuListResponseRows, err
	}
	queryGuoJiaJiaoYuListResponseRows = queryGuoJiaJiaoYuListResponse.Rows
	return queryGuoJiaJiaoYuListResponseRows, nil
}

func downloadGuoJiaJiaoYu(attachmentUrl string, filePath string) error {
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
	if GuoJiaJiaoYuEnableHttpProxy {
		client = GuoJiaJiaoYuSetHttpProxy()
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
	req.Header.Set("Host", "www.guojiajiaoyu.com")
	req.Header.Set("Referer", "http://www.guojiajiaoyu.com/")
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

func copyGuoJiaJiaoYuFile(src, dst string) (err error) {
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
