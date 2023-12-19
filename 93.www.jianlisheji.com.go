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
	"strings"
	"time"
)

const (
	JianLiSheJiEnableHttpProxy = false
	JianLiSheJiHttpProxyUrl    = "111.225.152.186:8089"
)

func JianLiSheJiSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(JianLiSheJiHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type JianLiSheJiSubject struct {
	name string
	url  string
}

var AllJianLiSheJiSubject = []JianLiSheJiSubject{
	{
		name: "中文简历",
		url:  "https://www.jianlisheji.com/jianli/jianlimuban/",
	},
	//{
	//	name: "英文简历",
	//	url:  "https://www.jianlisheji.com/jianli/yingwenjianli/",
	//},
	//{
	//	name: "表格简历",
	//	url:  "https://www.jianlisheji.com/jianli/biaogejianli/",
	//},
	{
		name: "小升初简历",
		url:  "https://www.jianlisheji.com/jianli/xiaoshengchu/",
	},
}

const JianLiSheJiCurrentAccountDownloadCountMAx = 10

var JianLiSheJiCurrentAccountId = 900000
var JianLiSheJiCookie = fmt.Sprintf("Hm_lvt_935dcd404e08577ddce430adb43b2cc9=1680165201; _gid=GA1.2.1939606718.1680165202; user_type=free; vip_expire_min=0; Hm_lpvt_935dcd404e08577ddce430adb43b2cc9=1680166294; JSESSIONID=598908E6CCA80F94C38EA17D7BEFF728; userid=%d; id_enpt=BiYMiYhrhb5tQJLJH8c4uQ==; avatar=https://www.jianlisheji.com/public/pc/common/default_head.png; mobile_flag=false; wechat_flag=true; login_token=dc0b2df7b35145ecea97e74c802ec4a9; _ga=GA1.1.1112728547.1680165202; _ga_34B604LFFQ=GS1.1.1680229977.2.1.1680233253.60.0.0", JianLiSheJiCurrentAccountId)
var JianLiSheJiCurrentAccountDownloadCount = 0
var JianLiSheJiNextDownloadSleep = 2

// ychEduSpider 获取简历设计文档
// @Title 获取简历设计文档
// @Description https://www.jianlisheji.com/，获取简历设计文档
func main() {
	for _, subject := range AllJianLiSheJiSubject {
		page := 1
		isPageListGo := true
		for isPageListGo {
			pageListUrl := fmt.Sprintf(subject.url+"?pageNumber=%d", page)
			fmt.Println(pageListUrl)

			pageListDoc, err := htmlquery.LoadURL(pageListUrl)
			if err != nil {
				fmt.Println(err)
				break
			}

			dlNodes := htmlquery.Find(pageListDoc, `//div[@class="word_item_con"]`)
			if len(dlNodes) >= 1 {
				for _, dlNode := range dlNodes {

					fmt.Println("=================================================================================")
					// 文档详情URL
					detailUrl := htmlquery.InnerText(htmlquery.FindOne(dlNode, `./dl[@class="word_item"]/dt/span/a/@href`))
					fmt.Println(detailUrl)

					detailDoc, _ := htmlquery.LoadURL(detailUrl)
					fileName := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="word_right"]/div[@class="inner"]/div[@class="title"]`))
					fileName = strings.ReplaceAll(fileName, "免费下载", "")
					fmt.Println(fileName)

					// 格式
					fileFormat := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="word_right"]/div[@class="inner"]/div[@class="detail"]/ul/li[1]/span[2]`))

					//编号
					code := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="word_right"]/div[@class="inner"]/div[@class="detail"]/ul/li[3]/span[2]`))

					vipCheckUrl := fmt.Sprintf("https://www.jianlisheji.com/vip/check/?v=%d&type=word_download&code=%s&rid=", time.Now().UnixMilli(), code)
					vipCheckReturn, err := vipCheckJianLiSheJi(vipCheckUrl, detailUrl)
					if err != nil {
						fmt.Println(err)
						// 当前账号id加一
						JianLiSheJiCurrentAccountId++
						JianLiSheJiCurrentAccountDownloadCount = 0
						continue
					}
					fmt.Println("========当前账户ID,", JianLiSheJiCurrentAccountId, "============")
					fmt.Println("========当前账户已下载,", JianLiSheJiCurrentAccountDownloadCount, "个文档============")
					// 下载文档URL
					downLoadUrl := fmt.Sprintf("https://www.jianlisheji.com/download/vip_download_word/?code=%s&keyid=%s&time=%s&encrypt=%s", vipCheckReturn.code, vipCheckReturn.keyid, vipCheckReturn.time, vipCheckReturn.encrypt)

					filePath := "F:\\workspace\\www.jianlisheji.com\\" + subject.name + "\\"
					fileName = fileName + "." + fileFormat
					if _, err := os.Stat(filePath + fileName); err != nil {
						fmt.Println("=======开始下载========")
						err = downloadJianLiSheJi(downLoadUrl, detailUrl, filePath, fileName)
						if err != nil {
							fmt.Println(err)
							// 当前账号id加一
							JianLiSheJiCurrentAccountId++
							JianLiSheJiCurrentAccountDownloadCount = 0
							continue
						}
						fmt.Println("=======下载完成========")
						for i := 1; i <= JianLiSheJiNextDownloadSleep; i++ {
							time.Sleep(time.Second)
							fmt.Println("===========操作结束，暂停", JianLiSheJiNextDownloadSleep, "秒，倒计时", i, "秒===========")
						}
						if JianLiSheJiCurrentAccountDownloadCount++; JianLiSheJiCurrentAccountDownloadCount >= JianLiSheJiCurrentAccountDownloadCountMAx {
							// 当前账号id加一
							JianLiSheJiCurrentAccountId++
							JianLiSheJiCurrentAccountDownloadCount = 0
							continue
						}
					}
				}
				page++
			} else {
				isPageListGo = false
				page = 1
				break
			}
		}
	}
}

type vipCheckResult struct {
	State int    `json:"state"`
	Msg   string `json:"msg"`
}

type vipCheckReturn struct {
	code    string
	keyid   string
	time    string
	encrypt string
}

func vipCheckJianLiSheJi(vipCheckUrl string, referer string) (vipCheckReturn vipCheckReturn, err error) {
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
	if JianLiSheJiEnableHttpProxy {
		client = JianLiSheJiSetHttpProxy()
	}
	postData := url.Values{}
	req, err := http.NewRequest("GET", vipCheckUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return vipCheckReturn, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", JianLiSheJiCookie)
	req.Header.Set("Host", "www.jianlisheji.com")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return vipCheckReturn, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return vipCheckReturn, err
	}
	vipCheckResult := &vipCheckResult{}
	err = json.Unmarshal(respBytes, vipCheckResult)
	if err != nil {
		fmt.Println(111)
		return vipCheckReturn, err
	}
	if vipCheckResult.State != 2 {
		return vipCheckReturn, errors.New(vipCheckResult.Msg)
	}
	var items map[string]interface{}
	_ = json.Unmarshal([]byte(vipCheckResult.Msg), &items)

	vipCheckReturn.code = items["code"].(string)
	vipCheckReturn.encrypt = items["encrypt"].(string)
	vipCheckReturn.keyid = items["keyid"].(string)
	vipCheckReturn.time = items["time"].(string)

	return vipCheckReturn, nil
}

func downloadJianLiSheJi(attachmentUrl string, referer string, filePath string, fileName string) error {
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
	if JianLiSheJiEnableHttpProxy {
		client = JianLiSheJiSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "https://www.jianlisheji.com")
	req.Header.Set("Referer", referer)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
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
	out, err := os.Create(filePath + fileName)
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
