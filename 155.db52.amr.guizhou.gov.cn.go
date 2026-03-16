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

var Db52EnableHttpProxy = false
var Db52HttpProxyUrl = "111.225.152.186:8089"
var Db52HttpProxyUrlArr = make([]string, 0)

func Db52HttpProxy() error {
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
					Db52HttpProxyUrlArr = append(Db52HttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					Db52HttpProxyUrlArr = append(Db52HttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func Db52SetHttpProxy() (httpclient *http.Client) {
	if Db52HttpProxyUrl == "" {
		if len(Db52HttpProxyUrlArr) <= 0 {
			err := Db52HttpProxy()
			if err != nil {
				Db52SetHttpProxy()
			}
		}
		Db52HttpProxyUrl = Db52HttpProxyUrlArr[0]
		if len(Db52HttpProxyUrlArr) >= 2 {
			Db52HttpProxyUrlArr = Db52HttpProxyUrlArr[1:]
		} else {
			Db52HttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(Db52HttpProxyUrl)
	ProxyURL, _ := url.Parse(Db52HttpProxyUrl)
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

type QueryDb52ListRequestPayload struct {
	Areaid      string `json:"areaid"`
	Flag      string `json:"flag"`
	IndustryType      string `json:"industryType"`
	IsAsc      string `json:"isAsc"`
	Name      string `json:"name"`
	Number      string `json:"number"`
	Param      QueryDb52ListRequestPayloadParam `json:"param"`
}
type QueryDb52ListRequestPayloadParam struct {
    PageNum int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

var Db52Cookie = "_trs_uv=mkki6w5m_4177_fod; MM_mq4qQammP3BA4=N2RiNTQxZTQ4NjU3MmRlYcl8P7aViJqnDKgbLcWNUuSpeKxX22X-8MelvKbKqGae11imuv2px1WHiCnacAv1iGUI5tLN9mg_fcmsALzP88Q"

// 下载贵州省地方标准文档
// @Title 下载贵州省地方标准文档
// @Description http://db52.amr.guizhou.gov.cn/，下载贵州省地方标准文档
func main() {
	page := 1
	maxPage := 177
	pageSize := 10
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		pageListUrl := fmt.Sprintf("https://db52.amr.guizhou.gov.cn/v2-service/std-back/client/standardpermit/page?pageNum=%d&pageSize=%d",page,pageSize)
	    fmt.Println(pageListUrl)
		queryDb52ListRequestPayload := QueryDb52ListRequestPayload{
			Areaid:      "",
			Flag:  "",
			IndustryType:   "",
			IsAsc:     "desc",
			Name: "",
			Number:  "",
			Param: QueryDb52ListRequestPayloadParam{
			    PageNum: page,
			    PageSize:pageSize,
			},
		}
		queryDb52ListResponseRows, err := QueryDb52List(pageListUrl, queryDb52ListRequestPayload)
// 		fmt.Println(queryDb52ListResponseRows)
		if err != nil {
			fmt.Println(err)
			break
		}
		for _, data := range queryDb52ListResponseRows {
			fmt.Println("=====================开始处理数据 page = ", page, "=========================")

			code := data.Snumber
			code = strings.ReplaceAll(code, "/", "-")
			code = strings.ReplaceAll(code, "—", "-")
			fmt.Println(code)

			title := data.Sname
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "--", "-")
			fmt.Println(title)

			filePath := "../db52.amr.guizhou.gov.cn/" + title + "(" + code + ")" + ".pdf"
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载========")

			downloadUrl := fmt.Sprintf("https://db52.amr.guizhou.gov.cn/v2-service/std-back/client/common/download/byPath?id=%d&disposition=attachment&name=%s", data.Id,url.QueryEscape(data.Sname))
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")

			err = downloadDb52(downloadUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "db52.amr.guizhou.gov.cn", "temp-dbba.sacinfo.org.cn")
			err = copyDb52File(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			//DownLoadDb52TimeSleep := 10
			DownLoadDb52TimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadDb52TimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadDb52TimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadDb52PageTimeSleep := 10
		// DownLoadDb52PageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadDb52PageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadDb52PageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

type QueryDb52ListResponse struct {
	Code      int                         `json:"code"`
	Rows      []QueryDb52ListResponseRows `json:"rows"`
	Msg       string                      `json:"msg"`
	Total     int                         `json:"total"`
}

type QueryDb52ListResponseRows struct {
    Id  int `json:"id"`
	Sname    string `json:"sname"`
	Snumber  string `json:"snumber"`
}

func QueryDb52List(requestUrl string, queryDb52ListRequestPayload QueryDb52ListRequestPayload) (queryDb52ListResponseRows []QueryDb52ListResponseRows, err error) {
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
	if Db52EnableHttpProxy {
		client = Db52SetHttpProxy()
	}
	// 将数据编码为JSON格式
	queryDb52ListRequestPayloadJson, err := json.Marshal(queryDb52ListRequestPayload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 创建字符串读取器，这是http.Request需要的类型
	body := bytes.NewReader(queryDb52ListRequestPayloadJson)
	req, err := http.NewRequest("POST", requestUrl, body) //建立连接

	queryDb52ListResponse := QueryDb52ListResponse{}
	if err != nil {
		return queryDb52ListResponseRows, err
	}

    req.Header.Set("authority", "db52.amr.guizhou.gov.cn")
	req.Header.Set("method", "POST")
	path := strings.Replace(requestUrl, "https://db52.amr.guizhou.gov.cn", "", 1)
	fmt.Println(path)
	req.Header.Set("path", path)
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Cookie", Db52Cookie)
	req.Header.Set("Host", "db52.amr.guizhou.gov.cn")
	req.Header.Set("Origin", "http://db52.amr.guizhou.gov.cn")
	req.Header.Set("Referer", "https://db52.amr.guizhou.gov.cn/v2/localStandard")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryDb52ListResponseRows, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryDb52ListResponseRows, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryDb52ListResponseRows, err
	}
	err = json.Unmarshal(respBytes, &queryDb52ListResponse)
	if err != nil {
		return queryDb52ListResponseRows, err
	}
	queryDb52ListResponseRows = queryDb52ListResponse.Rows
	return queryDb52ListResponseRows, nil
}

func downloadDb52(attachmentUrl string, filePath string) error {
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
	if Db52EnableHttpProxy {
		client = Db52SetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("authority", "db52.amr.guizhou.gov.cn")
	req.Header.Set("method", "GET")
	path := strings.Replace(attachmentUrl, "https://db52.amr.guizhou.gov.cn", "", 1)
	fmt.Println(path)
	req.Header.Set("path", path)
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", Db52Cookie)
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("Host", "db52.amr.guizhou.gov.cn")
	req.Header.Set("Referer", "https://db52.amr.guizhou.gov.cn/v2/localStandard")
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

func copyDb52File(src, dst string) (err error) {
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
