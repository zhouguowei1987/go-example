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

var LtBzhEnableHttpProxy = false
var LtBzhHttpProxyUrl = "111.225.152.186:8089"
var LtBzhHttpProxyUrlArr = make([]string, 0)

func LtBzhHttpProxy() error {
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
					LtBzhHttpProxyUrlArr = append(LtBzhHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					LtBzhHttpProxyUrlArr = append(LtBzhHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func LtBzhSetHttpProxy() (httpclient *http.Client) {
	if LtBzhHttpProxyUrl == "" {
		if len(LtBzhHttpProxyUrlArr) <= 0 {
			err := LtBzhHttpProxy()
			if err != nil {
				LtBzhSetHttpProxy()
			}
		}
		LtBzhHttpProxyUrl = LtBzhHttpProxyUrlArr[0]
		if len(LtBzhHttpProxyUrlArr) >= 2 {
			LtBzhHttpProxyUrlArr = LtBzhHttpProxyUrlArr[1:]
		} else {
			LtBzhHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(LtBzhHttpProxyUrl)
	ProxyURL, _ := url.Parse(LtBzhHttpProxyUrl)
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

type QueryLtBzhListFormData struct {
	pageType int
	_search  bool
	nd       int64
	pageSize int
	pageNo   int
	sidx     string
	sord     string
}

var LtBzhCookie = "JSESSIONID=B0A48CAA054563683978B5DFE826187A; ur_mofcom_=\"8n3ieq4YgFJbr52CHhdIjpVP0GkG03LJfEvgxC0KT9c=\"; insert_cookie=81869604; _pk_id.23.6b23=e1b1c7e065996d2d.1756110203.; _pk_ref.23.6b23=%5B%22%22%2C%22%22%2C1756181224%2C%22http%3A%2F%2Fwww.chinajl.com.cn%2F%22%5D"

// 下载商务领域行业标准文档
// @Title 下载商务领域行业标准文档
// @Description https://ltbzh.mofcom.gov.cn/，下载商务领域行业标准文档
func main() {
	pageListUrl := "https://ltbzh.mofcom.gov.cn/ltbz/bzgl/bzglController/listbzgl"
	fmt.Println(pageListUrl)
	page := 1
	maxPage := 13
	rows := 10
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		queryLtBzhListFormData := QueryLtBzhListFormData{
			pageType: 4,
			_search:  false,
			nd:       time.Now().Unix(),
			pageSize: rows,
			pageNo:   page,
			sidx:     "",
			sord:     "asc",
		}
		queryLtBzhListResponseResult, err := QueryLtBzhList(pageListUrl, queryLtBzhListFormData)
		if err != nil {
			LtBzhHttpProxyUrl = ""
			fmt.Println(err)
		}
		for _, result := range queryLtBzhListResponseResult {
			fmt.Println("=====================开始处理数据 page = ", page, "=========================")
			code := result.KeyName
			code = strings.ReplaceAll(code, "/", "-")
			code = strings.ReplaceAll(code, "—", "-")
			fmt.Println(code)

			title := result.StanderName
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "--", "-")
			fmt.Println(title)

			filePath := "../ltbzh.mofcom.gov.cn/" + title + "(" + code + ")" + ".pdf"
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载========")

			downloadUrl := fmt.Sprintf("https://ltbzh.mofcom.gov.cn/ltbz/ltbz/front/doDownload?fileName=%s&fileRealName=%s", result.ProBzGlFile, result.ProjFileName)
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")

			err = downloadLtBzh(downloadUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "../ltbzh.mofcom.gov.cn", "../upload.doc88.com/hbba.sacinfo.org.cn")
			err = copyLtBzhFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			//DownLoadLtBzhTimeSleep := 10
			DownLoadLtBzhTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadLtBzhTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadLtBzhTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadLtBzhPageTimeSleep := 10
		// DownLoadLtBzhPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadLtBzhPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadLtBzhPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

type QueryLtBzhListResponse struct {
	Page QueryLtBzhListResponsePage `json:"page"`
}

type QueryLtBzhListResponsePage struct {
	PageNo   int                            `json:"pageNo"`
	PageSize int                            `json:"pageSize"`
	PrePage  int                            `json:"prePage"`
	Result   []QueryLtBzhListResponseResult `json:"result"`
}

type QueryLtBzhListResponseResult struct {
	StanderName  string `json:"StanderName"`
	KeyName      string `json:"keyName"`
	ProBzGlFile  string `json:"probzglFile"`
	ProjFileName string `json:"projFileName"`
}

func QueryLtBzhList(requestUrl string, queryLtBzhListFormData QueryLtBzhListFormData) (queryLtBzhListResponseResult []QueryLtBzhListResponseResult, err error) {
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
	if LtBzhEnableHttpProxy {
		client = LtBzhSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("pageType", strconv.Itoa(queryLtBzhListFormData.pageType))
	postData.Add("_search", strconv.FormatBool(queryLtBzhListFormData._search))
	postData.Add("nd", strconv.FormatInt(queryLtBzhListFormData.nd, 10))
	postData.Add("pageSize", strconv.Itoa(queryLtBzhListFormData.pageSize))
	postData.Add("pageNo", strconv.Itoa(queryLtBzhListFormData.pageNo))
	postData.Add("sidx", queryLtBzhListFormData.sidx)
	postData.Add("sord", queryLtBzhListFormData.sord)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	queryLtBzhListResponse := QueryLtBzhListResponse{}
	if err != nil {
		return queryLtBzhListResponseResult, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", LtBzhCookie)
	req.Header.Set("Host", "ltbzh.mofcom.gov.cn")
	req.Header.Set("Origin", "https://ltbzh.mofcom.gov.cn")
	req.Header.Set("Referer", "https://ltbzh.mofcom.gov.cn/ltbz/view/bzfk/listBzfk.jsp")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryLtBzhListResponseResult, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryLtBzhListResponseResult, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryLtBzhListResponseResult, err
	}
	err = json.Unmarshal(respBytes, &queryLtBzhListResponse)
	if err != nil {
		return queryLtBzhListResponseResult, err
	}
	queryLtBzhListResponseResult = queryLtBzhListResponse.Page.Result
	return queryLtBzhListResponseResult, nil
}

func downloadLtBzh(attachmentUrl string, filePath string) error {
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
	if LtBzhEnableHttpProxy {
		client = LtBzhSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", LtBzhCookie)
	req.Header.Set("Host", "ltbzh.mofcom.gov.cn")
	req.Header.Set("Referer", "https://ltbzh.mofcom.gov.cn/public/StandardPermitManager2.jsp?bzbh=&bzname=")
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

func copyLtBzhFile(src, dst string) (err error) {
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
