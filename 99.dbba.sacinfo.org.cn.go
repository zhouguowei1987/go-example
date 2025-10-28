package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/otiai10/gosseract/v2"
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
	DbBaEnableHttpProxy = false
	DbBaHttpProxyUrl    = "111.225.152.186:8089"
)

func DbBaSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(DbBaHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type DbBaResponseData struct {
	Current     int                       `json:"current"`
	Pages       int                       `json:"pages"`
	Records     []DbBaResponseDataRecords `json:"records"`
	SearchCount bool                      `json:"searchCount"`
	Size        int                       `json:"size"`
	Total       int                       `json:"total"`
}
type DbBaResponseDataRecords struct {
	ActDate    int    `json:"actDate"`
	ChName     string `json:"chName"`
	ChargeDept string `json:"chargeDept"`
	Code       string `json:"code"`
	Empty      bool   `json:"empty"`
	Industry   string `json:"industry"`
	IssueDate  int    `json:"issueDate"`
	Pk         string `json:"pk"`
	RecordDate int    `json:"recordDate"`
	RecordNo   string `json:"recordNo"`
	Status     string `json:"status"`
}

type DbBaResponseValidateCaptcha struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

const DbBaCookie = "HMACCOUNT=487EF362690A1D5D; Hm_lvt_36f2f0446e1c2cda8410befc24743a9b=1759925510; Hm_lpvt_36f2f0446e1c2cda8410befc24743a9b=1761630427; JSESSIONID=219B2796825C6BAFE93AB153A2BCD665"

// ychEduSpider 获取地方标准文档
// @Title 获取地方标准文档
// @Description https://dbba.sacinfo.org.cn/，获取地方标准文档
func main() {
	requestUrl := "https://dbba.sacinfo.org.cn/stdQueryList"
	// 	5699
	current := 1
	maxCurrent := 751
	size := 100
	status := "现行"
	isPageListGo := true
	for isPageListGo {
		if current > maxCurrent {
			isPageListGo = false
			break
		}
		responseData, err := DbBaGetStdQueryList(requestUrl, current, size, status)
		if err != nil {
			fmt.Println(err)
			break
		}
		if len(responseData.Records) > 0 {
			for _, records := range responseData.Records {
				if records.Empty == false {
					fmt.Println("=======当前页为：" + strconv.Itoa(current) + "========")
					chName := strings.ReplaceAll(records.ChName, " ", "")
					chName = strings.ReplaceAll(chName, "/", "-")
					chName = strings.ReplaceAll(chName, "\n", "")
					chName = strings.ReplaceAll(chName, ":", "-")
					chName = strings.ReplaceAll(chName, "：", "-")

					//industry := strings.TrimSpace(records.Industry)

					code := strings.ReplaceAll(records.Code, "/", "-")
					code = strings.ReplaceAll(code, "\n", "")

					fileName := chName + "(" + code + ")"
					fmt.Println(fileName)

					filePath := "../dbba.sacinfo.org.cn/" + fileName + ".pdf"
					_, err = os.Stat(filePath)
					if err == nil {
						fmt.Println("文档已下载过，跳过")
						continue
					}

					stdDetailUrl := fmt.Sprintf("https://dbba.sacinfo.org.cn/stdDetail/%s", records.Pk)
					stdDetailDoc, err := htmlquery.LoadURL(stdDetailUrl)
					if err != nil {
						fmt.Println(err)
						continue
					}
					// 是否有查看文本按钮
					downloadButtonNode := htmlquery.FindOne(stdDetailDoc, `//div[@class="container main-body"]/div[@class="sidebar sidebar-left"]/div[@class="sidebar-tabs"]/a`)
					if downloadButtonNode == nil {
						fmt.Println("没有下载按钮跳过")
						continue
					}

					portalOnlineUrl := fmt.Sprintf("https://dbba.sacinfo.org.cn/portal/online/%s", records.Pk)
					portalOnlineDoc, err := htmlquery.LoadURL(portalOnlineUrl)
					if err != nil {
						fmt.Println(err)
						continue
					}
					// 是否有验证码窗口
					captchaModalDialogNode := htmlquery.FindOne(portalOnlineDoc, `//div[@class="container main-body"]/div[@class="row"]/div[@class="col-sm-12"]/div[@class="modal"]/div[@class="modal-dialog"]`)
					if captchaModalDialogNode == nil {
						fmt.Println("没有输入验证码窗口")
						continue
					}
				ValidateCaptchaGoTo:
					// 获取验证码图片
					// 获取当前时间的纳秒级时间戳
					nanoTimestamp := time.Now().UnixNano()
					// 将纳秒级时间戳转换为毫秒级
					millis := nanoTimestamp / 1e6 // 或者 nanoTimestamp / 1000000
					fmt.Println("当前时间的毫秒级时间戳:", millis)
					validateCodeUrl := fmt.Sprintf("https://dbba.sacinfo.org.cn/portal/validate-code?pk=%s&t=%d", records.Pk, millis)
					fmt.Println(validateCodeUrl)
					validateCodeFilePath := "./dbba-validate-code/validate-code.png"
					err = downloadValidateCodeDbBa(validateCodeUrl, validateCodeFilePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					// 获取验证码文字信息
					captcha, err := TesseractValidateCodeDbBa(validateCodeFilePath)
					captcha = strings.TrimSpace(captcha)
					fmt.Println("识别的验证码：", captcha)
					if len(captcha) != 4 {
						goto ValidateCaptchaGoTo
					}

					// 获取下载地址
					validateCaptchaReferer := fmt.Sprintf("https://dbba.sacinfo.org.cn/portal/online/%s", records.Pk)
					responseValidateCaptcha, err := validateCaptchaDbBa(captcha, records.Pk, validateCaptchaReferer)
					fmt.Println(responseValidateCaptcha, err)
					if err != nil {
						fmt.Println(err)
						continue
					}
					if responseValidateCaptcha.Code != 0 {
						fmt.Println(responseValidateCaptcha.Msg)
						goto ValidateCaptchaGoTo
					}
					downLoadUrl := fmt.Sprintf("https://dbba.sacinfo.org.cn/portal/download/%s", responseValidateCaptcha.Msg)
					fmt.Println(downLoadUrl)

					detailUrl := fmt.Sprintf("https://dbba.sacinfo.org.cn/stdDetail/%s", records.Pk)
					fmt.Println(detailUrl)

					fmt.Println("=======开始下载" + strconv.Itoa(current) + "========")
					err = downloadDbBa(downLoadUrl, detailUrl, filePath)
					if err != nil {
						fmt.Println(err)
						continue
					}

					// 查看文件大小，如果是空文件，则删除
                    fileInfo, err := os.Stat(filePath)
                    if err == nil && fileInfo.Size() == 0 {
                        fmt.Println("空文件删除")
                        err = os.Remove(filePath)
                    }
                    if err != nil {
                        continue
                    }

					//复制文件
					tempFilePath := strings.ReplaceAll(filePath, "../dbba.sacinfo.org.cn", "../upload.doc88.com/dbba.sacinfo.org.cn")
					err = DbBaCopyFile(filePath, tempFilePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("=======下载完成========")

					downloadDbBaPdfSleep := rand.Intn(5)
					for i := 1; i <= downloadDbBaPdfSleep; i++ {
						time.Sleep(time.Second)
						fmt.Println("page="+strconv.Itoa(current)+"=======chName=", chName, "成功，====== 暂停", downloadDbBaPdfSleep, "秒，倒计时", i, "秒===========")
					}
				}
			}

			DownLoadDbBaPageTimeSleep := 10
			// DownLoadDbBaPageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadDbBaPageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(current)+"====== 暂停", DownLoadDbBaPageTimeSleep, "秒 倒计时", i, "秒===========")
			}
			current++
			if current > maxCurrent {
				isPageListGo = false
				break
			}
		}
	}
}

func DbBaGetStdQueryList(requestUrl string, current int, size int, status string) (responseData DbBaResponseData, err error) {
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
	if DbBaEnableHttpProxy {
		client = DbBaSetHttpProxy()
	}
	responseData = DbBaResponseData{}
	postData := url.Values{}
	postData.Add("current", strconv.Itoa(current))
	postData.Add("size", strconv.Itoa(size))
	postData.Add("status", status)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return responseData, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "dbba.sacinfo.org.cn")
	req.Header.Set("Origin", "https://dbba.sacinfo.org.cn")
	req.Header.Set("Referer", "https://dbba.sacinfo.org.cn/stdList")
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
		return responseData, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseData, err
	}
	err = json.Unmarshal(respBytes, &responseData)
	if err != nil {
		return responseData, err
	}
	return responseData, nil
}

func downloadValidateCodeDbBa(validateCodeUrl string, filePath string) error {
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
	if DbBaEnableHttpProxy {
		client = DbBaSetHttpProxy()
	}
	req, err := http.NewRequest("GET", validateCodeUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("authority", "dbba.sacinfo.org.cn")
	req.Header.Set("method", "GET")
	path := strings.Replace(validateCodeUrl, "https://dbba.sacinfo.org.cn", "", 1)
	fmt.Println(path)
	req.Header.Set("path", path)
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "mage/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", DbBaCookie)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "dbba.sacinfo.org.cn")
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

func TesseractValidateCodeDbBa(imagePath string) (codeText string, err error) {
	// 创建Tesseract客户端
	client := gosseract.NewClient()
	defer client.Close()
	// 设置语言模型
	client.SetLanguage("eng")
	// 设置白名单字符
	client.SetWhitelist("0123456789abcdefghijklmnopqrstuvwxyz")
	// 识别图片中的文本
	err = client.SetImage(imagePath)
	if err != nil {
		return codeText, fmt.Errorf("设置图片出错: %v", err)
	}
	text, err := client.Text()
	if err != nil {
		return codeText, fmt.Errorf("识别出错: %v", err)
	}
	return text, nil
}

func validateCaptchaDbBa(captcha string, pk string, referer string) (responseValidateCaptcha DbBaResponseValidateCaptcha, err error) {
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
	if DbBaEnableHttpProxy {
		client = DbBaSetHttpProxy()
	}
	// 	fmt.Print("Enter an captcha and press enter: ")
	// 	fmt.Scanln(&captcha) // 等待用户按下回车键后继续执行
	// 	fmt.Println("You entered captcha:", captcha)
	responseValidateCaptcha = DbBaResponseValidateCaptcha{}
	requestUrl := fmt.Sprintf("https://dbba.sacinfo.org.cn/portal/validate-captcha/down?captcha=%s&pk=%s", captcha, pk)
	req, err := http.NewRequest("POST", requestUrl, nil) //建立连接
	if err != nil {
		return responseValidateCaptcha, err
	}
	req.Header.Set("authority", "dbba.sacinfo.org.cn")
	req.Header.Set("method", "POST")
	req.Header.Set("path", "/portal/validate-captcha/down")
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", DbBaCookie)
	req.Header.Set("Origin", "https://dbba.sacinfo.org.cn")
	req.Header.Set("Priority", "u=1, i")
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
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return responseValidateCaptcha, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseValidateCaptcha, err
	}
	err = json.Unmarshal(respBytes, &responseValidateCaptcha)
	if err != nil {
		return responseValidateCaptcha, err
	}
	return responseValidateCaptcha, nil
}

func downloadDbBa(attachmentUrl string, referer string, filePath string) error {
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
	if DbBaEnableHttpProxy {
		client = DbBaSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("authority", "dbba.sacinfo.org.cn")
	req.Header.Set("method", "GET")
	path := strings.Replace(attachmentUrl, "https://dbba.sacinfo.org.cn", "", 1)
	fmt.Println(path)
	req.Header.Set("path", path)
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", DbBaCookie)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "dbba.sacinfo.org.cn")
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

func DbBaCopyFile(src, dst string) (err error) {
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
