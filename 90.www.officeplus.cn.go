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
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	OfficePlusEnableHttpProxy = false
	OfficePlusHttpProxyUrl    = "111.225.152.186:8089"
)

func OfficePlusSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(OfficePlusHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type OfficePlusCategory struct {
	name     string
	keywords string
}

var OfficePlusAllCategory = []OfficePlusCategory{
	{
		name:     "Word模板",
		keywords: "word-content",
	},
	//{
	//	name:     "Excel模板",
	//	keywords: "excel-content",
	//},
}

type apiOfficePlusListResult struct {
	PageCount      int                            `json:"pageCount"`
	Items          []apiOfficePlusListResultItems `json:"items"`
	TotalItemCount int                            `json:"totalItemCount"`
}

type apiOfficePlusListResultItems struct {
	Id       string `json:"id"`
	FileName string `json:"fileName"`
	Title    string `json:"title"`
}

var officePlusAuthorization = "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJXaWNVc2VySWQiOiIwZDZmNmRkYy05ZDJiLTQ2ZTUtYTFhNS1kNGRkYTZjYWFkODAiLCJPZmZpY2VQbHVzVXNlcklkIjoiNzI4OWU4NzYtNjY4ZS01YTEzLTRiMTctM2EwYTM2ODk4ODg1Iiwic3ViIjoid2ljYXV0aEBvZmZpY2VwbHVzLmNuIiwiYXVkIjoid2ljYXV0aEBvZmZpY2VwbHVzLmNuIiwibmJmIjoxNzQxMTg3MzU4LCJleHAiOjE3NDEyNzM3NTgsImlhdCI6MTc0MTE4NzM1OCwiaXNzIjoid2ljYXV0aEBvZmZpY2VwbHVzLmNuIn0.mrYB82pkorBq4pVVTZsUYwJbqYjv3cZDohmz9VJRDZtQxatSZdNfTUWc6D5eUFm_ClRWV2Qcnpp5-TimnsQG5cvuXScsBuQ8YaxjuTkyxYHAqELVH5GFv6v0YEyF7eLIb0VtaQNy1PMjA0XqblQR0rI_1Y0QTN3jRIegrjJO9wGt2oe05AW-zhB_mB5Nxd3HklW8RX8Y37lQVYOGSA-gdpBLcWBhohzZF2vRfB5SYmRLeBHl4Hi6bQLWPAWmz8UEBW1Wz0Ma9f8OyYtF34VV6X56jCgAGjuMKZYuxu4RnXe9ApLNLX0fbog8gL1-tBdAWkHxh41j5MykJ5VyIneF1w"

// ychEduSpider 获取office-plus模板文档
// @Title 获取office-plus模板文档
// @Description https://www.officeplus.cn/，获取office-plus模板文档
func main() {
	for _, category := range OfficePlusAllCategory {
		page := 0
		for true {
			apiUrl := fmt.Sprintf("https://api.officeplus.cn/api/website/v2.1/contents/%s/search?orderBy=Total&pageIndex=%d&pageSize=30&paymentType=0&l2Method=1", category.keywords, page)
			fmt.Println(apiUrl)
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
			if OfficePlusEnableHttpProxy {
				client = OfficePlusSetHttpProxy()
			}
			postData := url.Values{}
			req, err := http.NewRequest("GET", apiUrl, strings.NewReader(postData.Encode())) //建立连接
			if err != nil {
				fmt.Println(err)
				page = 0
				break
			}
			req.Header.Set("Accept", "application/json, text/plain, */*")
			req.Header.Set("accept-encoding", "gzip, deflate, br")
			req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
			req.Header.Set("Origin", "https://www.officeplus.cn")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
			resp, err := client.Do(req) //拿到返回的内容
			if err != nil {
				fmt.Println(err)
				page = 1
				break
			}
			defer resp.Body.Close()
			respBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				page = 0
				break
			}
			apiOfficePlusListResult := &apiOfficePlusListResult{}
			err = json.Unmarshal(respBytes, apiOfficePlusListResult)
			if err != nil {
				fmt.Println(err)
				page = 0
				break
			}
			fmt.Println("======================================================")
			fmt.Println(category.name, category.keywords, page)
			for _, item := range apiOfficePlusListResult.Items {
				itemId := item.Id
				itemTitle := item.Title
				itemTitle = strings.ReplaceAll(itemTitle, "/", "-")
				itemTitle = strings.ReplaceAll(itemTitle, " ", "")
				fmt.Println(itemId, itemTitle)

				attachUrl, err := downloadUrl(itemId)
				if err != nil {
					fmt.Println(err)
					continue
				}
				filePath := "../www.officeplus.cn/" + category.name + "/" + itemTitle + "." + strings.Split(item.FileName, ".")[1]
				_, err = os.Stat(filePath)
				if err != nil {
					fmt.Println("=======开始下载========")
					err = downloadOfficePlus(attachUrl, filePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("=======完成下载========")
					DownLoadYchEduTimeSleep := rand.Intn(5)
					for i := 1; i <= DownLoadYchEduTimeSleep; i++ {
						time.Sleep(time.Second)
						fmt.Println("page="+strconv.Itoa(page)+"===========下载", itemTitle, "成功，暂停", DownLoadYchEduTimeSleep, "秒，倒计时", i, "秒===========")
					}
				}
			}
			if apiOfficePlusListResult.PageCount > page {
				page++
			} else {
				page = 0
				break
			}
		}
	}
}

type downloadUrlResult struct {
	DownloadUrl string `json:"downloadUrl"`
	UserId      string `json:"userId"`
}

func downloadUrl(id string) (fileUrl string, err error) {
	apiUrl := fmt.Sprintf("https://api.officeplus.cn/api/website/v2.1/download/%s/download-url", id)
	fmt.Println(apiUrl)
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
	if OfficePlusEnableHttpProxy {
		client = OfficePlusSetHttpProxy()
	}
	postData := url.Values{}
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	//req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("authorization", officePlusAuthorization)
	req.Header.Set("Origin", "https://www.officeplus.cn")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	downloadUrlResult := &downloadUrlResult{}
	err = json.Unmarshal(respBytes, downloadUrlResult)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fileUrl = downloadUrlResult.DownloadUrl
	fmt.Println(fileUrl)
	return fileUrl, nil
}

func downloadOfficePlus(attachmentUrl string, filePath string) error {
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
	if OfficePlusEnableHttpProxy {
		client = OfficePlusSetHttpProxy()
	}
	postData := url.Values{}
	req, err := http.NewRequest("GET", attachmentUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "Hm_lvt_c1b8f67dbb5e69873155161621f66842=1679986802; _gid=GA1.2.989925163.1679986804; optoken=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJXaWNVc2VySWQiOiIwZDZmNmRkYy05ZDJiLTQ2ZTUtYTFhNS1kNGRkYTZjYWFkODAiLCJXaWNVc2VyTmFtZSI6IldJQzhRQUlBMUNHIiwiV3hOaWNrTmFtZSI6Iktv5Y2XIiwiV3hIZWFkSW1hZ2VVcmwiOiJodHRwczovL3RoaXJkd3gucWxvZ28uY24vbW1vcGVuL3ZpXzMyL1EwajRUd0dUZlRMS2xHb2ZpYVB5bndKMnhkY2NqS3RqeUFQRWNLc3M4RWFrdXR3dDBFQXpWNjdUMFViZnd6YWhGVUdDQjAzVDFIcWFqWW9ZVjNjSVJPQS8xMzIiLCJXeFVuaW9uSWQiOiJvekdfZDZwckxyMVpTd3pEb25kWU9XOUtTblBJIiwiV3hPcGVuSWQiOiJvMjZBSzY4b25BOGFuUWRZTlR6X1ZucVRZb0VRIiwic3ViIjoid2ljYXV0aEBvZmZpY2VwbHVzLmNuIiwiYXVkIjoid2ljYXV0aEBvZmZpY2VwbHVzLmNuIiwibmJmIjoxNjc5OTg2ODE4LCJleHAiOjE2ODA1OTE2MTgsImlhdCI6MTY3OTk4NjgxOCwiaXNzIjoid2ljYXV0aEBvZmZpY2VwbHVzLmNuIn0.C2SQIlQgcaLAWGmDyx7DjZb0kNOhd0PEv1DJ66DB0Z2xnqo6D1B5CmFz-JudSYSDCYSpHrWwZAOtEIIuX8Qg2zgzxwYGmz_eTLLtO5SYdi-2QLymcA9s5wqd3NuBaO9fDvmJb5BqZ6uuAVPpB0jOBS_STjPTqVoxmB7SHmprNH4JeY5nUYT9a0y9mYBChNDwtPg0a59-BIME_-N8P5SpcsZelsqC8my7dbZMbiiI3ESVmIqgDM302z3t-RilgJ9oZAWlkDQ-Qg7npX5QrlNspd50F34mfm8dsnTcsC9h5K2gLBbg2zneBaDkR_sOdNmG1E5QIYTjkplLr11HEJemIA; Hm_lpvt_c1b8f67dbb5e69873155161621f66842=1679988896; _ga_34B604LFFQ=GS1.1.1679986804.1.1.1679988897.3.0.0; _ga=GA1.1.345890739.1679986804")
	req.Header.Set("Host", "content-prod.officeplus.cn")
	req.Header.Set("Origin", "https://www.officeplus.cn")
	req.Header.Set("Referer", "https://www.officeplus.cn/")
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
