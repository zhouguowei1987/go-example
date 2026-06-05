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

var KsBzEnableHttpProxy = false
var KsBzHttpProxyUrl = "111.225.152.186:8089"
var KsBzHttpProxyUrlArr = make([]string, 0)

func KsBzHttpProxy() error {
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
					KsBzHttpProxyUrlArr = append(KsBzHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					KsBzHttpProxyUrlArr = append(KsBzHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func KsBzSetHttpProxy() (httpclient *http.Client) {
	if KsBzHttpProxyUrl == "" {
		if len(KsBzHttpProxyUrlArr) <= 0 {
			err := KsBzHttpProxy()
			if err != nil {
				KsBzSetHttpProxy()
			}
		}
		KsBzHttpProxyUrl = KsBzHttpProxyUrlArr[0]
		if len(KsBzHttpProxyUrlArr) >= 2 {
			KsBzHttpProxyUrlArr = KsBzHttpProxyUrlArr[1:]
		} else {
			KsBzHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(KsBzHttpProxyUrl)
	ProxyURL, _ := url.Parse(KsBzHttpProxyUrl)
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

var KsBzListCookie = "JSESSIONID=KlmXV3CuDTNAvxWO-pzNe1fDvMqTAxFpbNKE0M-0"

// 下载矿山安全标准文档
// @Title 下载矿山安全标准文档
// @Description https://ksbz.chinamine-safety.gov.cn/，下载矿山安全标准文档
func main() {
	page := 1
	maxPage := 7
	count := 100
	isPageListGo := true
	for isPageListGo {
		timestamp := time.Now().Unix()
		pageListUrl := fmt.Sprintf("https://ksbz.chinamine-safety.gov.cn/revisionSystem/standardFile/getStandardFilePage?_t=%d&fileType=&page=%d&size=%d&releaseTime=&title=&standardStatus=", timestamp, page, count)
		fmt.Println(pageListUrl)
		queryKsBzListResponseResultLists, err := QueryKsBzList(pageListUrl)
		if err != nil {
			fmt.Println(err)
			break
		}
		for _, ksBz := range queryKsBzListResponseResultLists {
			fmt.Println("=========开始处理数据==============")

			code := ksBz.StandardNumber
			code = strings.ReplaceAll(code, "/", "-")
			code = strings.ReplaceAll(code, "—", "-")
			fmt.Println(code)

			title := ksBz.Title
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "--", "-")
			fmt.Println(title)

			filePath := "../ksbz.chinamine-safety.gov.cn/" + title + "(" + code + ")" + ".pdf"
			filePath = strings.ReplaceAll(filePath, "()", "")
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			detailUrl := fmt.Sprintf("https://ksbz.chinamine-safety.gov.cn/revisionSystem/standardFile/queryById?_t=%d&id=%s&isLike=2", timestamp, ksBz.Id)
			fmt.Println(detailUrl)
			queryKsBzDetailResponseResultStandardAnnexList, err := QueryKsBzDetail(detailUrl)
			if err != nil {
				fmt.Println(err)
				break
			}
            // 使用QueryEscape进行URL编码
            chineseTitle := url.QueryEscape(queryKsBzDetailResponseResultStandardAnnexList[0].FileName)
			downloadUrl := fmt.Sprintf("https://ksbz.chinamine-safety.gov.cn/revisionSystem/legalFileUpload/downloadFile?src=%s&chineseTitle=%s", queryKsBzDetailResponseResultStandardAnnexList[0].Src, chineseTitle)
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")

			err = downloadKsBz(downloadUrl, detailUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "ksbz.chinamine-safety.gov.cn", "temp-ksbz.chinamine-safety.gov.cn")
			err = copyKsBzFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			DownLoadKsBzTimeSleep := 10
			// DownLoadKsBzTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadKsBzTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(page)+"====== 暂停,filePath="+filePath+"===========下载成功 暂停", DownLoadKsBzTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadKsBzCategoryTimeSleep := 10
		// DownLoadKsBzCategoryTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadKsBzCategoryTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"====== 暂停", DownLoadKsBzCategoryTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

type QueryKsBzListResponse struct {
	Code      int                         `json:"code"`
	Success   bool                        `json:"success"`
	Result    QueryKsBzListResponseResult `json:"result"`
	TimeStamp int                         `json:"timestamp"`
}
type QueryKsBzListResponseResult struct {
	Lists []QueryKsBzListResponseResultLists `json:"lists"`
	Count int                                `json:"count"`
}

type QueryKsBzListResponseResultLists struct {
	Id             string `json:"id"`
	StandardNumber string `json:"standard_number"`
	Title          string `json:"title"`
}

func QueryKsBzList(requestUrl string) (queryKsBzListResponseResultLists []QueryKsBzListResponseResultLists, err error) {
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
	if KsBzEnableHttpProxy {
		client = KsBzSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	queryKsBzListResponse := QueryKsBzListResponse{}
	if err != nil {
		return queryKsBzListResponseResultLists, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Cookie", KsBzListCookie)
	req.Header.Set("Host", "ksbz.chinamine-safety.gov.cn")
	req.Header.Set("Origin", "https://ksbz.chinamine-safety.gov.cn")
	req.Header.Set("Referer", "https://ksbz.chinamine-safety.gov.cn/mine/standard")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryKsBzListResponseResultLists, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryKsBzListResponseResultLists, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryKsBzListResponseResultLists, err
	}
	err = json.Unmarshal(respBytes, &queryKsBzListResponse)
	if err != nil {
		fmt.Println(err)
		return queryKsBzListResponseResultLists, err
	}
	queryKsBzListResponseResultLists = queryKsBzListResponse.Result.Lists
	return queryKsBzListResponseResultLists, nil
}

type QueryKsBzDetailResponse struct {
	Code      int                           `json:"code"`
	Success   bool                          `json:"success"`
	Result    QueryKsBzDetailResponseResult `json:"result"`
	TimeStamp int                           `json:"timestamp"`
}
type QueryKsBzDetailResponseResult struct {
	StandardAnnexList []QueryKsBzDetailResponseResultStandardAnnexList `json:"standardAnnexList"`
}

type QueryKsBzDetailResponseResultStandardAnnexList struct {
	FileName string `json:"fileName"`
	Src      string `json:"src"`
}

func QueryKsBzDetail(requestUrl string) (queryKsBzDetailResponseResultStandardAnnexList []QueryKsBzDetailResponseResultStandardAnnexList, err error) {
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
	if KsBzEnableHttpProxy {
		client = KsBzSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	queryKsBzDetailResponse := QueryKsBzDetailResponse{}
	if err != nil {
		return queryKsBzDetailResponseResultStandardAnnexList, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Cookie", KsBzListCookie)
	req.Header.Set("Host", "ksbz.chinamine-safety.gov.cn")
	req.Header.Set("Origin", "https://ksbz.chinamine-safety.gov.cn")
	req.Header.Set("Referer", "https://ksbz.chinamine-safety.gov.cn/mine/standard")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryKsBzDetailResponseResultStandardAnnexList, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryKsBzDetailResponseResultStandardAnnexList, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryKsBzDetailResponseResultStandardAnnexList, err
	}
	err = json.Unmarshal(respBytes, &queryKsBzDetailResponse)
	if err != nil {
		fmt.Println(err)
		return queryKsBzDetailResponseResultStandardAnnexList, err
	}
	queryKsBzDetailResponseResultStandardAnnexList = queryKsBzDetailResponse.Result.StandardAnnexList
	return queryKsBzDetailResponseResultStandardAnnexList, nil
}

func downloadKsBz(attachmentUrl string, referer string, filePath string) error {
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
	if KsBzEnableHttpProxy {
		client = KsBzSetHttpProxy()
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
	req.Header.Set("Cookie", KsBzListCookie)
	req.Header.Set("Host", "ksbz.chinamine-safety.gov.cn")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
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

func copyKsBzFile(src, dst string) (err error) {
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
