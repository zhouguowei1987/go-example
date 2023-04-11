package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
	TiKuVip51TestEnableHttpProxy = false
	TiKuVip51TestHttpProxyUrl    = "111.225.152.186:8089"
)

func TiKuVip51TestSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(TiKuVip51TestHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

// ychEduSpider 获取无忧考试网真题
// @Title 获取无忧考试网真题
// @Description https://tikuvip.51test.net/，获取无忧考试网真题
func main() {
	tiKuVip51TestTreeListInitData, err := treeListInit()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", tiKuVip51TestTreeListInitData)
}

type TiKuVip51TestTreeListInitResult struct {
	Code    bool                            `json:"code"`
	UseTime float64                         `json:"use_time"`
	Data    []TiKuVip51TestTreeListInitData `json:"data"`
}

type TiKuVip51TestTreeListInitData struct {
	IsParent bool                                    `json:"isParent"`
	Name     string                                  `json:"name"`
	Open     bool                                    `json:"open"`
	Path     string                                  `json:"path"`
	Type     string                                  `json:"type"`
	Children []TiKuVip51TestTreeListInitDataChildren `json:"children"`
}

type TiKuVip51TestTreeListInitDataChildren struct {
	ATime       int64  `json:"atime"`
	CTime       int64  `json:"ctime"`
	IsParent    bool   `json:"isParent"`
	IsReadable  int    `json:"isReadable"`
	IsWriteable int    `json:"isWriteable"`
	Mode        string `json:"mode"`
	MTime       int64  `json:"mtime"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
}

func treeListInit() (tiKuVip51TestTreeListInitData []TiKuVip51TestTreeListInitData, err error) {
	apiUrl := fmt.Sprintf("https://tikuvip.51test.net/index.php?share/treeList&app=folder&user=100&sid=BzcEWh8C&type=init")
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
	if TiKuVip51TestEnableHttpProxy {
		client = TiKuVip51TestSetHttpProxy()
	}
	postData := url.Values{}
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		fmt.Println(err)
		return tiKuVip51TestTreeListInitData, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("cookie", "__yjs_duid=1_9906f2c8d7c017db33d48b6a18ccd56b1675842585518; __bid_n=18630008879d6e0aee4207; Hm_lvt_f4ae163e87a012d4ab5106f993decb4c=1678772717; HOST=https://tikuvip.51test.net/; APP_HOST=https://tikuvip.51test.net/; kodUserLanguage=zh-CN; KOD_SESSION_SSO=n8v73kudnk59513468tvhg87tb; KOD_SESSION_ID_8e194=4oqj9kpf2pv95a06ob9jtgn9t9; Hm_lvt_c3d24798f142b815b7069d109e892e45=1681193783; Hm_lpvt_c3d24798f142b815b7069d109e892e45=1681193791; FPTOKEN=UE8f6wb8MckxSxA89FIzTZ9nPQu2jojdRQf4VC8sc1QD/+6ogqCPSaSNzQEejyrOERvpDkcNLTobKcgtrh4HBHntkNvvf/elpIuqc/oTjkvNrkyQSRfyPlZ6jm9dYo3c/48EVuuYIExPhgwAdN3uYxIbH7T3h+s+F/RMB9b7hk8HDyBAJqiZIcMKwKqYgiPDRL8unhH+0FqFBIoJADdvPMAxfThWnvolGsCCpU+jZUPoXZmBCWUf88+amY4wvGNiBbcoIZYh1tZfd2Hd+AF+HWz0RsBdmMc0FHvDBx/mxmPIeZrFww3nE7PP185jycCBL2D4vAVHyVvvBI8R5nCUrC/zG8ya3XsTSI0LHdQsN/Dg8J5Fvjjlmfb//2zgb6jqF9AZvf0CZ/3R5YRO2kkdKw==|wfSc6dBQePRjKadw5Z+bnD3Mhzt+CJbcPPgSMmC9Td4=|10|acd5f2507c1333048f3938686dc63972; Hm_lpvt_f4ae163e87a012d4ab5106f993decb4c=1681195808; kodVersionCheck=check-at-1681195868")
	req.Header.Set("referer", "https://tikuvip.51test.net/?share/folder&user=100&sid=BzcEWh8C&uid=8034602&uip=222.70.7.91&downloaddate=2023-04-11&token=00fded6d31b3f0bf47f14d3251bc120c")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return tiKuVip51TestTreeListInitData, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return tiKuVip51TestTreeListInitData, err
	}
	tiKuVip51TestTreeListInitResult := &TiKuVip51TestTreeListInitResult{}
	err = json.Unmarshal(respBytes, tiKuVip51TestTreeListInitResult)
	if err != nil {
		fmt.Println(err)
		return tiKuVip51TestTreeListInitData, err
	}

	tiKuVip51TestTreeListInitData = tiKuVip51TestTreeListInitResult.Data
	return tiKuVip51TestTreeListInitData, nil
}

type TiKuVip51TestDownloadUrlResult struct {
	DownloadUrl    string `json:"downloadUrl"`
	IsFreeDownload bool   `json:"isFreeDownload"`
	UserId         string `json:"userId"`
}

func TiKuVip51TestDownloadUrl(id string) (fileUrl string, err error) {
	apiUrl := fmt.Sprintf("https://api.officeplus.cn/api/v2.0/web/download/%s/download-url", id)
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
	if TiKuVip51TestEnableHttpProxy {
		client = TiKuVip51TestSetHttpProxy()
	}
	postData := url.Values{}
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJXaWNVc2VySWQiOiIwZDZmNmRkYy05ZDJiLTQ2ZTUtYTFhNS1kNGRkYTZjYWFkODAiLCJXaWNVc2VyTmFtZSI6IldJQzhRQUlBMUNHIiwiV3hOaWNrTmFtZSI6Iktv5Y2XIiwiV3hIZWFkSW1hZ2VVcmwiOiJodHRwczovL3RoaXJkd3gucWxvZ28uY24vbW1vcGVuL3ZpXzMyL1EwajRUd0dUZlRMS2xHb2ZpYVB5bndKMnhkY2NqS3RqeUFQRWNLc3M4RWFrdXR3dDBFQXpWNjdUMFViZnd6YWhGVUdDQjAzVDFIcWFqWW9ZVjNjSVJPQS8xMzIiLCJXeFVuaW9uSWQiOiJvekdfZDZwckxyMVpTd3pEb25kWU9XOUtTblBJIiwiV3hPcGVuSWQiOiJvMjZBSzY4b25BOGFuUWRZTlR6X1ZucVRZb0VRIiwic3ViIjoid2ljYXV0aEBvZmZpY2VwbHVzLmNuIiwiYXVkIjoid2ljYXV0aEBvZmZpY2VwbHVzLmNuIiwibmJmIjoxNjc5OTg2ODE4LCJleHAiOjE2ODA1OTE2MTgsImlhdCI6MTY3OTk4NjgxOCwiaXNzIjoid2ljYXV0aEBvZmZpY2VwbHVzLmNuIn0.C2SQIlQgcaLAWGmDyx7DjZb0kNOhd0PEv1DJ66DB0Z2xnqo6D1B5CmFz-JudSYSDCYSpHrWwZAOtEIIuX8Qg2zgzxwYGmz_eTLLtO5SYdi-2QLymcA9s5wqd3NuBaO9fDvmJb5BqZ6uuAVPpB0jOBS_STjPTqVoxmB7SHmprNH4JeY5nUYT9a0y9mYBChNDwtPg0a59-BIME_-N8P5SpcsZelsqC8my7dbZMbiiI3ESVmIqgDM302z3t-RilgJ9oZAWlkDQ-Qg7npX5QrlNspd50F34mfm8dsnTcsC9h5K2gLBbg2zneBaDkR_sOdNmG1E5QIYTjkplLr11HEJemIA")
	req.Header.Set("Origin", "https://tikuvip.51test.net")
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
	TiKuVip51TestDownloadUrlResult := &TiKuVip51TestDownloadUrlResult{}
	err = json.Unmarshal(respBytes, TiKuVip51TestDownloadUrlResult)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fileUrl = TiKuVip51TestDownloadUrlResult.DownloadUrl
	fmt.Println(fileUrl)
	return fileUrl, nil
}

func downloadTiKuVip51Test(attachmentUrl string, filePath string, fileName string) error {
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
	if TiKuVip51TestEnableHttpProxy {
		client = TiKuVip51TestSetHttpProxy()
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
	req.Header.Set("Origin", "https://tikuvip.51test.net")
	req.Header.Set("Referer", "https://tikuvip.51test.net/")
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
