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
	"regexp"
	"strconv"
	"strings"
	"time"
)

var YuWenChaZiDianEnableHttpProxy = false
var YuWenChaZiDianHttpProxyUrl = ""
var YuWenChaZiDianHttpProxyUrlArr = make([]string, 0)

func YuWenChaZiDianHttpProxy() error {
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
					YuWenChaZiDianHttpProxyUrlArr = append(YuWenChaZiDianHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					YuWenChaZiDianHttpProxyUrlArr = append(YuWenChaZiDianHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func YuWenChaZiDianSetHttpProxy() (httpclient *http.Client) {
	if YuWenChaZiDianHttpProxyUrl == "" {
		if len(YuWenChaZiDianHttpProxyUrlArr) <= 0 {
			err := YuWenChaZiDianHttpProxy()
			if err != nil {
				YuWenChaZiDianSetHttpProxy()
			}
		}
		YuWenChaZiDianHttpProxyUrl = YuWenChaZiDianHttpProxyUrlArr[0]
		if len(YuWenChaZiDianHttpProxyUrlArr) >= 2 {
			YuWenChaZiDianHttpProxyUrlArr = YuWenChaZiDianHttpProxyUrlArr[1:]
		} else {
			YuWenChaZiDianHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(YuWenChaZiDianHttpProxyUrl)
	ProxyURL, _ := url.Parse(YuWenChaZiDianHttpProxyUrl)
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

//var YuWenChaZiDianNextDownloadSleep = 2

// ychEduSpider 获取查字典语文网试卷
// @Title 获取查字典语文网试卷
// @Description https://yuwen.chazidian.com/，获取查字典语文网试卷
func main() {
	page := 950
	isPageListGo := true

	for isPageListGo {
		yuWenChaZiDianResponse, err := YuWenChaZiDianGetCatePageDataApi(page)
		if err != nil {
			fmt.Println(err)
			page = 1
			isPageListGo = false
			break
		}
		if len(yuWenChaZiDianResponse.Data) <= 0 {
			fmt.Println(err)
			page = 1
			isPageListGo = false
			break
		}
		for _, data := range yuWenChaZiDianResponse.Data {
			fmt.Println("============================================================================")
			fmt.Println("=======当前页URL", page, "========")

			title := data.Title
			fmt.Println(title)
			if strings.Contains(title, "图片") || strings.Contains(title, "扫描") {
				fmt.Println("标题含有图片，扫描字样")
				continue
			}

			viewHref := "https://yuwen.chazidian.com/shijuan" + strconv.Itoa(data.Id)
			fmt.Println(viewHref)

			// 查看是否有附件
			viewDoc, err := htmlquery.LoadURL(viewHref)
			if err != nil {
				fmt.Println(err)
				continue
			}

			regAttachmentViewUrl := regexp.MustCompile(`<a href="//yuwen.chazidian.com/uploadfile/(.*?)" style="color: brown">立即下载</a>`)
			regAttachmentViewUrlMatch := regAttachmentViewUrl.FindAllSubmatch([]byte(htmlquery.OutputHTML(viewDoc, true)), -1)
			if len(regAttachmentViewUrlMatch) <= 0 {
				fmt.Println("没有附件，跳过")
				continue
			}
			attachmentUrl := "https://yuwen.chazidian.com/uploadfile/" + string(regAttachmentViewUrlMatch[0][1])
			fmt.Println(attachmentUrl)
			fileExtIndex := strings.LastIndex(attachmentUrl, ".")
			fileExt := attachmentUrl[fileExtIndex:]
			if !YuWenChaZiDianStrInArray(fileExt, []string{".doc", ".docx", ".rar"}) {
				fmt.Println("文件后缀：" + fileExt + "不在下载后缀列表")
				continue
			}

			filePath := "E:\\workspace\\yuwen.chazidian.com\\yuwen.rar_chazidian.com\\" + title + fileExt
			_, err = os.Stat(filePath)
			if err != nil {
				fmt.Println("=======开始下载" + strconv.Itoa(page) + "========")
				err = downloadYuWenChaZiDian(attachmentUrl, viewHref, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")
				//for i := 1; i <= YuWenChaZiDianNextDownloadSleep; i++ {
				//	time.Sleep(time.Second)
				//	fmt.Println("===========操作结束，暂停", YuWenChaZiDianNextDownloadSleep, "秒，倒计时", i, "秒===========")
				//}
			}
		}
		page++
		isPageListGo = true
	}
}

type YuWenChaZiDianResponse struct {
	Code int                          `json:"code"`
	Data []YuWenChaZiDianResponseData `json:"Data"`
}
type YuWenChaZiDianResponseData struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Date        string `json:"date"`
}

func YuWenChaZiDianGetCatePageDataApi(page int) (yuWenChaZiDianResponse YuWenChaZiDianResponse, err error) {
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
	if YuWenChaZiDianEnableHttpProxy {
		client = YuWenChaZiDianSetHttpProxy()
	}
	yuWenChaZiDianResponse = YuWenChaZiDianResponse{}
	postData := url.Values{}
	postData.Add("table", "ol_shijuan")
	postData.Add("page", strconv.Itoa(page))
	postData.Add("catid", "7, 8, 9, 10, 11, 12, 13, 14, 16")
	postData.Add("kewenid", "")
	postData.Add("nianji", "")
	postData.Add("banben", "")
	postData.Add("ce", "")
	req, err := http.NewRequest("POST", "https://yuwen.chazidian.com/index/api/getCatePageDataApi", strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return yuWenChaZiDianResponse, err
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", "Hm_lvt_392e83603a0def58379f6aa1f9e6a93b=1723435257; HMACCOUNT=2CEC63D57647BCA5; Hm_lvt_1b8f22a621ad677920c7dfdb50ececf1=1723435257; PHPSESSID=73d9a4951fdff72d542b29eb9bf473e8; Hm_lvt_9ac05ff87f40912e7310348f7565b387=1723435303; Hm_lpvt_1b8f22a621ad677920c7dfdb50ececf1=1723435659; Hm_lpvt_392e83603a0def58379f6aa1f9e6a93b=1723435659; Hm_lpvt_9ac05ff87f40912e7310348f7565b387=1723435660")
	req.Header.Set("Host", "yuwen.chazidian.com")
	req.Header.Set("Origin", "https://yuwen.chazidian.com")
	req.Header.Set("Referer", "https://yuwen.chazidian.com/shijuan/")
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
		return yuWenChaZiDianResponse, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	respString := string(respBytes)
	respString = strings.ReplaceAll(respString, "\\\"", "\"")
	respString = strings.Trim(respString, "\"")
	respString, err = strconv.Unquote(strings.Replace(strconv.Quote(respString), `\\u`, `\u`, -1))
	respString = strings.ReplaceAll(respString, "\\", "")
	if err != nil {
		return yuWenChaZiDianResponse, err
	}
	err = json.Unmarshal([]byte(respString), &yuWenChaZiDianResponse)
	if err != nil {
		return yuWenChaZiDianResponse, err
	}
	return yuWenChaZiDianResponse, nil
}

func downloadYuWenChaZiDian(attachmentUrl string, referer string, filePath string) error {
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
	if YuWenChaZiDianEnableHttpProxy {
		client = YuWenChaZiDianSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "yuwen.chazidian.com")
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

// YuWenChaZiDianStrInArray str in string list
func YuWenChaZiDianStrInArray(str string, data []string) bool {
	if len(data) > 0 {
		for _, row := range data {
			if str == row {
				return true
			}
		}
	}
	return false
}
