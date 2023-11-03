package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"io/ioutil"
	"math/rand"
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
	ZZStepEnableHttpProxy = true
	ZZStepHttpProxyUrl    = "http://115.29.148.215:8999"
)

func ZZStepSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(ZZStepHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type ZZStepSubject struct {
	name string
	url  string
}

var subjects = []ZZStepSubject{
	{
		name: "试卷",
		url:  "http://www2.zzstep.com/front/paper/index.html",
	},
	{
		name: "中考",
		url:  "http://www2.zzstep.com/front/beikao/index.html",
	},
}

var NextDownloadSleep = 60

var randStringLength = 8

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 当前账号已下载文档数量
var eachUsernameDownloadCurrentCount = 0

// 每个账号最大下载数量
var eachUsernameDownloadMaxCount = 80
var password = "123456"
var refer = "http://www.zzstep.com/"
var ZZStepCookie = "looyu_id=a8650c60217050342a55f236559725ac_20002564%3A9; _99_mon=%5B0%2C0%2C0%5D; PHPSESSID=0fsf2nbmid269b03lvhl9btp63; Hm_lvt_e22ffa3ceb40ebd886ecaaa1c24eb75d=1698997450; Hm_lvt_5d656a0b54535a39b1d88c84eff1725b=1698997450; looyu_20002564=v%3A17c0cfff633630c6ec83f487d6852626%2Cref%3A%2Cr%3A%2Cmon%3A//m6817.talk99.cn/monitor%2Cp0%3Ahttp%253A//www2.zzstep.com/; zzstep_front_user=%00l%0BoWfTd%5CcWl%07a%073_lR0%060%5C%3AV2R%3C%07%7BRs%03eW0%01%3A%04%26%0CvQ3Qc%0D%2B%05%3DWl%04i%01iQ1%04bQoR8%00j%0BoWmT%27%5CiWd%07k%07%21_2Ri%062%5CkVhRg%07mRe%03%7DW8%01s%04%3E%0C%3BQ%60Q%25%0Dn%05hWw%042%01%1FQb%04%12Q%3BR%22%00g%0B%2CWlTl%5CiW~%07%24%07p_9Rr%06%3F%5CaVkRc%07%22R%3B%03%2CW9%018%04%3E%0C%21Q%3DQi%0D%7B%05gWG%040%01%1FQl%04rQmRs%00f%0BnWfTn%5CqW0%07%3E%07d_5Rn%063%5CyVhRc%07wR%22%03dWp%01%3A%047%0C9QxQ~%0Dl%05uW%27%04h%01%20; zzstep_front_user=%00l%0BoWfTd%5CcWl%07a%073_lR0%060%5C%3AV2R%3C%07%7BRs%03eW0%01%3A%04%26%0CvQ3Qc%0D%2B%05%3DWl%04i%01iQ1%04bQoR8%00j%0BoWmT%27%5CiWd%07k%07%21_2Ri%062%5CkVhRg%07mRe%03%7DW8%01s%04%3E%0C%3BQ%60Q%25%0Dn%05hWw%042%01%1FQb%04%12Q%3BR%22%00g%0B%2CWlTl%5CiW~%07%24%07p_9Rr%06%3F%5CaVkRc%07%22R%3B%03%2CW9%018%04%3E%0C%21Q%3DQi%0D%7B%05gWG%040%01%1FQl%04rQmRs%00f%0BnWfTn%5CqW0%07%3E%07d_5Rn%063%5CyVhRc%07wR%22%03dWp%01%3A%047%0C9QxQ~%0Dl%05uW%27%04h%01%20; cdb_cookietime=2592000; cdb_compound=e198GOozSUa22Tp2vut%2Bv%2Fjhe9JBqI%2FoaeOHl7SbE32on%2BK%2F4pb3x2cSgFq0z8QGKyCmWSjrLsZ4mz4%2BquyPvEVBRYJb2pNXKpVsjO3S; cdb_auth=geE%2BCuBgAdjMwSwW1uW12KzHCtilCj46KpazC%2FFEzWwC5AI5nt%2FruO556%2FvzOHW6HA; Hm_lpvt_5d656a0b54535a39b1d88c84eff1725b=1698998778; Hm_lpvt_e22ffa3ceb40ebd886ecaaa1c24eb75d=1698998778"

// ychEduSpider 获取中国教育出版网文档
// @Title 获取中国教育出版网文档
// @Description http://www2.zzstep.com/，获取中国教育出版网文档
func main() {
	//// 首先注册登陆新账号
	//err := ZZStepRegisterLoginUsername()
	//if err != nil {
	//	return
	//}
	for _, subject := range subjects {
		current := 30
		isPageListGo := true
		for isPageListGo {
			subjectIndexUrl := subject.url
			if current > 1 {
				subjectIndexUrl += fmt.Sprintf("?studysection=204&subject=29&page=%d", current)
			}
			subjectIndexDoc, err := htmlquery.LoadURL(subjectIndexUrl)
			if err != nil {
				fmt.Println(err)
				current = 1
				isPageListGo = false
				continue
			}
			liNodes := htmlquery.Find(subjectIndexDoc, `//div[@class="zy-list fn-mt20"]/ul[@class="reslist"]/li[@class="fn-pt20 fn-pb20"]`)
			if len(liNodes) <= 0 {
				fmt.Println(err)
				current = 1
				isPageListGo = false
				continue
			}
			for _, liNode := range liNodes {
				fmt.Println("============================================================================")
				fmt.Println("主题：", subject.name)
				fmt.Println("=======当前页为：" + strconv.Itoa(current) + "========")

				// 所需智币
				pointsNode := htmlquery.FindOne(liNode, `./div[@class="btn-item fn-left"]/div[@class="money fn-pt10"]`)
				if pointsNode == nil {
					fmt.Println("没有智币div")
					continue
				}
				pointsText := htmlquery.InnerText(pointsNode)
				fmt.Println(pointsText)
				pointsText = strings.ReplaceAll(pointsText, "智币", "")

				points, err := strconv.Atoi(pointsText)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if points > 0 {
					fmt.Println("需要智币下载", points)
					continue
				}

				// 当前文件类型
				fileExtTextNode := htmlquery.FindOne(liNode, `./div[@class="filetype fn-pl10 fn-left"]/img/@src`)
				if fileExtTextNode == nil {
					fmt.Println("没有文件类型div")
					continue
				}
				fileExtText := htmlquery.InnerText(fileExtTextNode)
				fileExtText = strings.ReplaceAll(fileExtText, "/public/front/images/", "")
				if fileExtText != "typeicon-word.png" {
					fmt.Println(fileExtText, "不在下载后缀列表")
					continue
				}

				fileName := htmlquery.InnerText(htmlquery.FindOne(liNode, `./div[@class="zy-box fn-left"]/div[@class="subject-t"]/a`))
				fileName = strings.TrimSpace(fileName)
				fileName = strings.ReplaceAll(fileName, "/", "-")
				fileName = strings.ReplaceAll(fileName, ":", "-")
				fileName = strings.ReplaceAll(fileName, "：", "-")
				fileName = strings.ReplaceAll(fileName, "（", "(")
				fileName = strings.ReplaceAll(fileName, "）", ")")
				fmt.Println(fileName)

				filePath := "../www2.zzstep.com/www2.zzstep.com/" + subject.name + "/" + fileName
				_, errDoc := os.Stat(filePath + ".doc")
				_, errDocx := os.Stat(filePath + ".docx")
				if errDoc != nil && errDocx != nil {
					viewUrl := "http://www2.zzstep.com" + htmlquery.InnerText(htmlquery.FindOne(liNode, `./div[@class="zy-box fn-left"]/div[@class="subject-t"]/a/@href`))
					fmt.Println(viewUrl)

					downLoadUrl := strings.ReplaceAll(viewUrl, "index", "download")
					fmt.Println(downLoadUrl)

					fmt.Println("=======开始下载" + strconv.Itoa(current) + "========")
					err = downloadZZStep(downLoadUrl, viewUrl, filePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("=======下载完成========")
					for i := 1; i <= NextDownloadSleep; i++ {
						time.Sleep(time.Second)
						fmt.Println("===========操作结束，暂停", NextDownloadSleep, "秒，倒计时", i, "秒===========")
					}
					if eachUsernameDownloadCurrentCount++; eachUsernameDownloadCurrentCount >= eachUsernameDownloadMaxCount {
						//注册登陆新账号
						err = ZZStepRegisterLoginUsername()
						if err != nil {
							return
						}
					}
				}
			}
			current++
			isPageListGo = true
		}
	}
}

func ZZStepRegisterLoginUsername() error {
	// 注册新账号
	rand.Seed(time.Now().UnixNano()) // 设置随机种子
	// 生成长度为randStringLength的随机字符串
	username := randStringBytes(randStringLength)
	fmt.Println(username)
	err := ZZStepRegisterRandUsername(username, password, password, refer)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 登陆
	err = ZZStepLoginUsername(username, password, refer)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

type ZZStepRegisterRandUsernameResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func ZZStepRegisterRandUsername(username string, password string, password2 string, refer string) error {
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
	if ZZStepEnableHttpProxy {
		client = ZZStepSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("username", username)
	postData.Add("password", password)
	postData.Add("password2", password2)
	postData.Add("refer", refer)
	requestUrl := "http://www2.zzstep.com/front/regist/registbyusername.html"
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Host", "www2.zzstep.com")
	req.Header.Set("Origin", "http://www2.zzstep.com")
	req.Header.Set("Referer", "http://www2.zzstep.com/front/regist/index.html?refer=http://www.zzstep.com/")
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

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	zZStepRegisterRandAccountResp := ZZStepRegisterRandUsernameResp{}
	err = json.Unmarshal(respBytes, &zZStepRegisterRandAccountResp)
	if err != nil {
		return err
	}

	if zZStepRegisterRandAccountResp.Code != 1 {
		return errors.New(zZStepRegisterRandAccountResp.Msg)
	}
	return nil
}

type ZZStepLoginUsernameResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func ZZStepLoginUsername(username string, passwordu string, refer string) error {
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
	if ZZStepEnableHttpProxy {
		client = ZZStepSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("username", username)
	postData.Add("passwordu", passwordu)
	postData.Add("type", "username")
	postData.Add("refer", refer)
	requestUrl := "http://www2.zzstep.com/front/login/dologin.html"
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Host", "www2.zzstep.com")
	req.Header.Set("Origin", "http://www2.zzstep.com")
	req.Header.Set("Referer", "http://www2.zzstep.com/front/login/index.html?refer=http://www2.zzstep.com/")
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

	respBytes, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBytes))
	if err != nil {
		return err
	}

	zZStepLoginUsernameResp := ZZStepLoginUsernameResp{}
	err = json.Unmarshal(respBytes, &zZStepLoginUsernameResp)
	if err != nil {
		return err
	}

	if zZStepLoginUsernameResp.Code != 1 {
		return errors.New(zZStepLoginUsernameResp.Msg)
	}

	// 重新设置cookie
	ZZStepCookie = resp.Header.Get("Set-Cookie")
	return nil
}

func downloadZZStep(attachmentUrl string, referer string, filePath string) error {
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
	if ZZStepEnableHttpProxy {
		client = ZZStepSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", ZZStepCookie)
	req.Header.Set("Host", "www2.zzstep.com")
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
	// 检查HTTP响应头中的Content-Disposition字段获取文件名和后缀
	fileName := getZZStepFileNameFromHeader(resp)
	fileExtension := filepath.Ext(fileName) // 获取文件后缀
	fileExtArr := []string{".doc", ".docx"}
	fmt.Println("文件后缀:", fileExtension)
	if !StrInArrayZZStep(fileExtension, fileExtArr) {
		return errors.New("文件后缀：" + fileExtension + "不在下载后缀列表")
	}
	filePath += fileExtension
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

// 生成随机字符串
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// StrInArrayZZStep str in string list
func StrInArrayZZStep(str string, data []string) bool {
	if len(data) > 0 {
		for _, row := range data {
			if str == row {
				return true
			}
		}
	}
	return false
}

// 从HTTP响应头中获取文件名
func getZZStepFileNameFromHeader(resp *http.Response) string {
	contentDisposition := resp.Header.Get("Content-Disposition")
	fileName := ""
	if contentDisposition != "" {
		fileName = parseZZStepFileNameFromContentDisposition(contentDisposition)
	} else {
		fileName = filepath.Base(resp.Request.URL.Path) // 默认使用URL中的文件名作为本地文件名
	}
	return fileName
}

// 从Content-Disposition字段中解析文件名
func parseZZStepFileNameFromContentDisposition(contentDisposition string) string {
	// 参考：https://tools.ietf.org/html/rfc6266#section-4.3
	// 示例：attachment; filename="example.txt" -> example.txt
	fileNameStart := len("attachment; ") + len("filename=") + 1
	fileNameEnd := len(contentDisposition) - 1
	fileName := contentDisposition[fileNameStart:fileNameEnd] // 提取文件名字符串
	return fileName[:]                                        // 去掉字符串开头的引号（如果存在）并返回结果
}
