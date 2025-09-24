package main

import (
    "errors"
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"math/rand"
	"io"
)

var OsTaEnableHttpProxy = false
var OsTaHttpProxyUrl = ""
var OsTaHttpProxyUrlArr = make([]string, 0)

func OsTaHttpProxy() error {
	pageMax := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	//pageMax := []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
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
					OsTaHttpProxyUrlArr = append(OsTaHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					OsTaHttpProxyUrlArr = append(OsTaHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func OsTaSetHttpProxy() (httpclient *http.Client) {
	if OsTaHttpProxyUrl == "" {
		if len(OsTaHttpProxyUrlArr) <= 0 {
			err := OsTaHttpProxy()
			if err != nil {
				OsTaSetHttpProxy()
			}
		}
		OsTaHttpProxyUrl = OsTaHttpProxyUrlArr[0]
		if len(OsTaHttpProxyUrlArr) >= 2 {
			OsTaHttpProxyUrlArr = OsTaHttpProxyUrlArr[1:]
		} else {
			OsTaHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(OsTaHttpProxyUrl)
	ProxyURL, _ := url.Parse(OsTaHttpProxyUrl)
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

var OsTaCookie = "Hm_lvt_e85984af56dd04582a569a53719e397f=1757738790,1758696058; HMACCOUNT=1CCD0111717619C6; Hm_lpvt_e85984af56dd04582a569a53719e397f=1758696361"

//var OsTaNextDownloadSleep = 2

// ychEduSpider 获取国家职业技能标准
// @Title 获取国家职业技能标准
// @Description http://www.osta.org.cn/，获取国家职业技能标准
func main() {
	page := 1
	isPageListGo := true
	for isPageListGo {
		requestUrl := fmt.Sprintf("https://www.osta.org.cn/api/public/skillStandardList?pageSize=20&pageNum=%d", page)
		osTaListResponseBodyList, err := OsTaList(requestUrl)
		if err != nil {
			fmt.Println(err)
			page = 1
			isPageListGo = false
			break
		}
		if len(osTaListResponseBodyList) <= 0 {
			fmt.Println(err)
			page = 1
			isPageListGo = false
			break
		}
		for _, row := range osTaListResponseBodyList {
			fmt.Println("=====================开始处理数据 page = ", page, "=========================")

			filePath := "../www.osta.org.cn/" + row.Name + "-"+row.IssueNumber+"（" + row.Code + "）.pdf"
			fmt.Println(filePath)
            _, err = os.Stat(filePath)
            if err == nil {
                fmt.Println("文档已下载过，跳过")
                continue
            }
            downLoadUrl := "https://www.osta.org.cn/api/sys/downloadFile/decrypt?fileName=" + row.StandardInfo
            fmt.Println(downLoadUrl)

            fmt.Println("=======开始下载========")
            err = downloadOsTa(downLoadUrl, filePath)
            if err != nil {
                fmt.Println(err)
                continue
            }
            fmt.Println("=======完成下载========")

            //复制文件
            tempFilePath := strings.ReplaceAll(filePath, "../www.osta.org.cn", "../upload.doc88.com/www.osta.org.cn")
            err = copyOsTaFile(filePath, tempFilePath)
            if err != nil {
                fmt.Println(err)
                continue
            }

            // 设置倒计时
            // DownLoadOsTaTimeSleep := 10
            DownLoadOsTaTimeSleep := rand.Intn(5)
            for i := 1; i <= DownLoadOsTaTimeSleep; i++ {
                time.Sleep(time.Second)
                fmt.Println("name="+row.Name+"===========操作完成，", "暂停", DownLoadOsTaTimeSleep, "秒，倒计时", i, "秒===========")
            }
		}
		page++
		isPageListGo = true
	}
}

type OsTaListResponse struct {
	Code  int                        `json:"code"`
	Msg   string                     `json:"msg"`
	Body  OsTaListResponseBody `json:"body"`
}
type OsTaListResponseBody struct {
	List         []OsTaListResponseBodyList    `json:"list"`
}
type OsTaListResponseBodyList struct {
	Id         int    `json:"id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	IssueNumber   string `json:"issueNumber"`
	StandardInfo   string `json:"standardInfo"`
}

func OsTaList(requestUrl string) (osTaListResponseBodyList []OsTaListResponseBodyList, err error) {
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
	if OsTaEnableHttpProxy {
		client = OsTaSetHttpProxy()
	}
	osTaListResponse := OsTaListResponse{}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return osTaListResponseBodyList, err
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", OsTaCookie)
	req.Header.Set("Host", "www.osta.org.cn")
	req.Header.Set("Origin", "http://www.osta.org.cn")
	req.Header.Set("Referer", "https://www.osta.org.cn/skillStandard")
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
		return osTaListResponseBodyList, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return osTaListResponseBodyList, err
	}
	err = json.Unmarshal(respBytes, &osTaListResponse)
	if err != nil {
		return osTaListResponseBodyList, err
	}
    osTaListResponseBodyList = osTaListResponse.Body.List
	return osTaListResponseBodyList, nil
}

func downloadOsTa(pdfUrl string, filePath string) error {
	client := &http.Client{}                        //初始化客户端
	req, err := http.NewRequest("GET", pdfUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.osta.org.cn")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
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
		if os.MkdirAll(fileDiv, 0644) != nil {
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

func copyOsTaFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer func(in *os.File) {
		err := in.Close()
		if err != nil {
			return
		}
	}(in)

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)

	_, err = io.Copy(out, in)
	return
}

