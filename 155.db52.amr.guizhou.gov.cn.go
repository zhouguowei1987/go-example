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
			ResponseHeaderTimeout: time.Second * 3,
		},
	}
	return httpclient
}

type QueryDb52ListFormData struct {
	flag      string
	areaCode  string
	stdName   string
	stdNo     string
	pageIndex int
	pageSize  int
}

var Db52Cookie = "MM_mq4qQammP3BA4=YWQ0ZDViZDJjNjU1ZDM3NWL_PeyHzJuOhwD3DeuuFPagVwvsL_hlX5xSBMX33CpDzk4E-yCuM-FAzJQt3iaB4kftLzMlhQu_iNEMY43S09s"

// 下载贵州省地方标准文档
// @Title 下载贵州省地方标准文档
// @Description http://db52.amr.guizhou.gov.cn/，下载贵州省地方标准文档
func main() {
	pageListUrl := "https://db52.amr.guizhou.gov.cn/yongjie-gzstd-api/std/std/lib/listLib"
	fmt.Println(pageListUrl)
	page := 0
	maxPage := 205
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		queryDb52ListFormData := QueryDb52ListFormData{
			flag:      "",
			areaCode:  "",
			stdName:   "",
			stdNo:     "",
			pageIndex: page,
			pageSize:  10,
		}
		queryDb52ListResponseData, err := QueryDb52List(pageListUrl, queryDb52ListFormData)
		if err != nil {
			fmt.Println(err)
			break
		}
		for _, data := range queryDb52ListResponseData {
			fmt.Println("=====================开始处理数据 page = ", page, "=========================")

			if len(data.FilePath) <= 0 {
				fmt.Println("文档没有附件地址，跳过")
				continue
			}

			code := data.StdNo
			code = strings.ReplaceAll(code, "/", "-")
			code = strings.ReplaceAll(code, "—", "-")
			fmt.Println(code)

			title := data.StdName
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "--", "-")
			fmt.Println(title)
			if strings.Index(title, "测试项目") != -1 {
				fmt.Println("测试文档，跳过")
				continue
			}

			filePath := "../db52.amr.guizhou.gov.cn/" + title + "(" + code + ")" + ".pdf"
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载========")

			downloadUrl := fmt.Sprintf("https://db52.amr.guizhou.gov.cn/yongjie-gzstd-api/%s", data.FilePath)
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
	Data      []QueryDb52ListResponseData `json:"data"`
	Msg       string                      `json:"msg"`
	PageIndex int                         `json:"pageIndex"`
	PageSize  int                         `json:"pageSize"`
	Total     int                         `json:"total"`
	TotalPage int                         `json:"totalPage"`
}

type QueryDb52ListResponseData struct {
	StdNo    string `json:"stdNo"`
	StdName  string `json:"stdName"`
	FilePath string `json:"filePath"`
}

func QueryDb52List(requestUrl string, queryDb52ListFormData QueryDb52ListFormData) (queryDb52ListResponseData []QueryDb52ListResponseData, err error) {
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
	if Db52EnableHttpProxy {
		client = Db52SetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("flag", queryDb52ListFormData.flag)
	postData.Add("areaCode", queryDb52ListFormData.areaCode)
	postData.Add("stdName", queryDb52ListFormData.stdName)
	postData.Add("stdNo", queryDb52ListFormData.stdNo)
	postData.Add("pageIndex", strconv.Itoa(queryDb52ListFormData.pageIndex))
	postData.Add("pageSize", strconv.Itoa(queryDb52ListFormData.pageSize))
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	queryDb52ListResponse := QueryDb52ListResponse{}
	if err != nil {
		return queryDb52ListResponseData, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", Db52Cookie)
	req.Header.Set("Host", "db52.amr.guizhou.gov.cn")
	req.Header.Set("Origin", "http://db52.amr.guizhou.gov.cn")
	req.Header.Set("Referer", "https://db52.amr.guizhou.gov.cn/portal/localStandardsSearch?searchName=%E5%9C%B0%E6%96%B9%E6%A0%87%E5%87%86%E6%9F%A5%E8%AF%A2")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryDb52ListResponseData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryDb52ListResponseData, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryDb52ListResponseData, err
	}
	err = json.Unmarshal(respBytes, &queryDb52ListResponse)
	if err != nil {
		return queryDb52ListResponseData, err
	}
	queryDb52ListResponseData = queryDb52ListResponse.Data
	return queryDb52ListResponseData, nil
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
			ResponseHeaderTimeout: time.Second * 3,
		},
	}
	if Db52EnableHttpProxy {
		client = Db52SetHttpProxy()
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
	req.Header.Set("Cookie", Db52Cookie)
	req.Header.Set("Host", "db52.amr.guizhou.gov.cn")
	req.Header.Set("Referer", "https://db52.amr.guizhou.gov.cn/portal/localStandardsSearch?searchName=%E5%9C%B0%E6%96%B9%E6%A0%87%E5%87%86%E6%9F%A5%E8%AF%A2")
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
