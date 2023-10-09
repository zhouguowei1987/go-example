package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	LvLinEnableHttpProxy = false
	LvLinHttpProxyUrl    = "111.225.152.186:8089"
)

func LvLinSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(LvLinHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type LawContractListResponse struct {
	ErrMsg string                      `json:"errmsg"`
	ErrNo  int                         `json:"errno"`
	Data   LawContractListResponseData `json:"data"`
	LogId  int                         `json:"logid"`
}
type LawContractListResponseData struct {
	ContractList []LawContractListResponseDataContractList `json:"contractList"`
	Total        int                                       `json:"total"`
}

type LawContractListResponseDataContractList struct {
	Title string `json:"title"`
	Cmd   string `json:"cmd"`
}

// ychEduSpider 获取律临合同文书
// @Title 获取律临合同文书
// @Description https://lvlin.baidu.com/，获取律临合同文书
func main() {
	pn := 1
	rn := 50
	//精选合同：1 企业合同：3
	lawContractType := 1
	isPageListGo := true
	for isPageListGo {
		requestUrl := fmt.Sprintf("https://lvlin.baidu.com/pc/zdQuestionApi/question/api/lawcontractlist?clientType=pc&law_category1=&law_category2=&type=%d&pn=%d&rn=%d", lawContractType, pn, rn)
		lawContractListResponse, err := GetLawContractList(requestUrl)
		if err != nil {
			fmt.Println(err)
			break
		}
		if lawContractListResponse.ErrNo == 0 {
			for _, contract := range lawContractListResponse.Data.ContractList {
				fmt.Println("=======当前页为：" + strconv.Itoa(pn) + "========")
				title := contract.Title
				fmt.Println(title)

				detailUrl := contract.Cmd
				fmt.Println(detailUrl)

				detailDoc, err := htmlquery.LoadURL(detailUrl)
				if err != nil {
					fmt.Println(err)
					break
				}

				downloadButtonNode := htmlquery.FindOne(detailDoc, `//div[@class="content-box"]/div[@class="content-left"]/div[@class="bottom"]/div[@class="btn-box"]/a[2]/@href`)
				if downloadButtonNode == nil {
					fmt.Println("没有下载按钮")
					continue
				}

				downLoadUrl := htmlquery.InnerText(downloadButtonNode)
				fmt.Println(downLoadUrl)

				filePath := "../lvlin.baidu.com/" + title + ".docx"
				if _, err := os.Stat(filePath); err != nil {
					fmt.Println("=======开始下载" + strconv.Itoa(pn) + "========")
					err = downloadLvLin(downLoadUrl, detailUrl, filePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("=======开始完成========")
				}

				// 查看文件大小，如果是空文件，则删除
				fi, err := os.Stat(filePath)
				if err == nil && fi.Size() == 0 {
					err := os.Remove(filePath)
					if err != nil {
						continue
					}
				}

				time.Sleep(time.Millisecond * 100)
			}

			if pn < lawContractListResponse.Data.Total/rn {
				pn++
			} else {
				isPageListGo = false
				pn = 1
				break
			}
		} else {
			isPageListGo = false
			pn = 1
			break
		}
	}
}

func GetLawContractList(requestUrl string) (lawContractListResponse LawContractListResponse, err error) {
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
	if LvLinEnableHttpProxy {
		client = LvLinSetHttpProxy()
	}
	lawContractListResponse = LawContractListResponse{}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return lawContractListResponse, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "lvlin.baidu.com")
	req.Header.Set("Origin", "https://lvlin.baidu.com/")
	req.Header.Set("Referer", "https://lvlin.baidu.com/pc/channel")
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
		return lawContractListResponse, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return lawContractListResponse, err
	}
	err = json.Unmarshal(respBytes, &lawContractListResponse)
	if err != nil {
		return lawContractListResponse, err
	}
	return lawContractListResponse, nil
}

func downloadLvLin(attachmentUrl string, referer string, filePath string) error {
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
	if LvLinEnableHttpProxy {
		client = LvLinSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "lvlin.baidu.com")
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
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}

	// 创建一个文件用于保存
	fileDiv := filepath.Dir(filePath)
	if _, err = os.Stat(fileDiv); err != nil {
		if os.MkdirAll(fileDiv, 0777) != nil {
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
