package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
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
	for _, treeListInitData := range tiKuVip51TestTreeListInitData {
		fmt.Println(treeListInitData.Name)
		for _, childTreeListInitData := range treeListInitData.Children {
			// 请求下级文件夹
			tiKuVip51TestOneLevelTreeListData, err := treeList(childTreeListInitData.Path)
			if err != nil {
				fmt.Println(err)
				continue
			}

			for _, oneLevelTreeListData := range tiKuVip51TestOneLevelTreeListData {
				if oneLevelTreeListData.IsParent {
					tiKuVip51TestTwoLevelTreeListData, err := treeList(oneLevelTreeListData.Path)
					if err != nil {
						fmt.Println(err)
						continue
					}
					for _, twoLevelTreeListData := range tiKuVip51TestTwoLevelTreeListData {
						if twoLevelTreeListData.IsParent {
							// 请求下级文件夹
							tiKuVip51TestThreeLevelTreeListData, err := treeList(twoLevelTreeListData.Path)
							if err != nil {
								fmt.Println(err)
								continue
							}
							for _, threeLevelTreeListData := range tiKuVip51TestThreeLevelTreeListData {
								if !threeLevelTreeListData.IsParent {
									// 请求文件
									tiKuVip51TestThreeLevelTreeListDataPathListDataFileList, err := PathList(threeLevelTreeListData.Path)
									if err != nil {
										fmt.Println(err)
										continue
									}
									tiKuVip51TestDownloadUrl(tiKuVip51TestThreeLevelTreeListDataPathListDataFileList)
								}
							}
						} else {
							// 请求文件
							tiKuVip51TestTwoLevelTreeListDataPathListDataFileList, err := PathList(twoLevelTreeListData.Path)
							if err != nil {
								fmt.Println(err)
								continue
							}
							// 下载文件
							tiKuVip51TestDownloadUrl(tiKuVip51TestTwoLevelTreeListDataPathListDataFileList)
						}
					}
				} else {
					// 请求文件
					tiKuVip51TestOneLevelTreeListDataPathListDataFileList, err := PathList(oneLevelTreeListData.Path)
					if err != nil {
						fmt.Println(err)
						continue
					}
					// 下载文件
					tiKuVip51TestDownloadUrl(tiKuVip51TestOneLevelTreeListDataPathListDataFileList)
				}
			}
		}
	}
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

type TiKuVip51TestTreeListResult struct {
	Code    bool                        `json:"code"`
	UseTime float64                     `json:"use_time"`
	Data    []TiKuVip51TestTreeListData `json:"data"`
}

type TiKuVip51TestTreeListData struct {
	IsParent    bool   `json:"isParent"`
	Name        string `json:"name"`
	ATime       int64  `json:"atime"`
	CTime       int64  `json:"ctime"`
	IsReadable  int    `json:"isReadable"`
	IsWriteable int    `json:"isWriteable"`
	Mode        string `json:"mode"`
	MTime       int64  `json:"mtime"`
	Path        string `json:"path"`
	Type        string `json:"type"`
}

func treeList(path string) (tiKuVip51TestTreeListData []TiKuVip51TestTreeListData, err error) {
	apiUrl := fmt.Sprintf("https://tikuvip.51test.net/index.php?share/treeList&app=folder&user=100&sid=BzcEWh8C")
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
	postData := make(map[string]string)
	postData["path"] = path
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k, v := range postData {
		w.WriteField(k, v)
	}
	w.Close()

	req, err := http.NewRequest("POST", apiUrl, body) //建立连接
	if err != nil {
		fmt.Println(err)
		return tiKuVip51TestTreeListData, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("cookie", "__yjs_duid=1_9906f2c8d7c017db33d48b6a18ccd56b1675842585518; __bid_n=18630008879d6e0aee4207; Hm_lvt_f4ae163e87a012d4ab5106f993decb4c=1678772717; HOST=https://tikuvip.51test.net/; APP_HOST=https://tikuvip.51test.net/; kodUserLanguage=zh-CN; KOD_SESSION_SSO=n8v73kudnk59513468tvhg87tb; KOD_SESSION_ID_8e194=4oqj9kpf2pv95a06ob9jtgn9t9; Hm_lvt_c3d24798f142b815b7069d109e892e45=1681193783; Hm_lpvt_c3d24798f142b815b7069d109e892e45=1681193791; FPTOKEN=UE8f6wb8MckxSxA89FIzTZ9nPQu2jojdRQf4VC8sc1QD/+6ogqCPSaSNzQEejyrOERvpDkcNLTobKcgtrh4HBHntkNvvf/elpIuqc/oTjkvNrkyQSRfyPlZ6jm9dYo3c/48EVuuYIExPhgwAdN3uYxIbH7T3h+s+F/RMB9b7hk8HDyBAJqiZIcMKwKqYgiPDRL8unhH+0FqFBIoJADdvPMAxfThWnvolGsCCpU+jZUPoXZmBCWUf88+amY4wvGNiBbcoIZYh1tZfd2Hd+AF+HWz0RsBdmMc0FHvDBx/mxmPIeZrFww3nE7PP185jycCBL2D4vAVHyVvvBI8R5nCUrC/zG8ya3XsTSI0LHdQsN/Dg8J5Fvjjlmfb//2zgb6jqF9AZvf0CZ/3R5YRO2kkdKw==|wfSc6dBQePRjKadw5Z+bnD3Mhzt+CJbcPPgSMmC9Td4=|10|acd5f2507c1333048f3938686dc63972; Hm_lpvt_f4ae163e87a012d4ab5106f993decb4c=1681195808; kodVersionCheck=check-at-1681195868")
	req.Header.Set("origin", "https://tikuvip.51test.net")
	req.Header.Set("referer", "https://tikuvip.51test.net/?share/folder&user=100&sid=BzcEWh8C&uid=8034602&uip=222.70.7.91&downloaddate=2023-04-11&token=00fded6d31b3f0bf47f14d3251bc120c")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return tiKuVip51TestTreeListData, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return tiKuVip51TestTreeListData, err
	}
	tiKuVip51TestTreeListResult := &TiKuVip51TestTreeListResult{}
	err = json.Unmarshal(respBytes, tiKuVip51TestTreeListResult)
	if err != nil {
		fmt.Println(err)
		return tiKuVip51TestTreeListData, err
	}

	tiKuVip51TestTreeListData = tiKuVip51TestTreeListResult.Data
	return tiKuVip51TestTreeListData, nil
}

type TiKuVip51TestPathListResult struct {
	Code    bool                      `json:"code"`
	UseTime float64                   `json:"use_time"`
	Data    TiKuVip51TestPathListData `json:"data"`
}

type TiKuVip51TestPathListData struct {
	FileList []TiKuVip51TestPathListDataFileList `json:"fileList"`
}

type TiKuVip51TestPathListDataFileList struct {
	IsParent    bool   `json:"isParent"`
	Ext         string `json:"ext"`
	Name        string `json:"name"`
	ATime       int64  `json:"atime"`
	CTime       int64  `json:"ctime"`
	IsReadable  int    `json:"isReadable"`
	IsWriteable int    `json:"isWriteable"`
	Mode        string `json:"mode"`
	MTime       int64  `json:"mtime"`
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	Type        string `json:"type"`
}

func PathList(path string) (tiKuVip51TestPathListDataFileList []TiKuVip51TestPathListDataFileList, err error) {
	apiUrl := fmt.Sprintf("https://tikuvip.51test.net/index.php?share/pathList&user=100&sid=BzcEWh8C&path=%s", path)
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
		return tiKuVip51TestPathListDataFileList, err
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
		return tiKuVip51TestPathListDataFileList, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return tiKuVip51TestPathListDataFileList, err
	}
	tiKuVip51TestPathListResult := &TiKuVip51TestPathListResult{}
	err = json.Unmarshal(respBytes, tiKuVip51TestPathListResult)
	if err != nil {
		fmt.Println(err)
		return tiKuVip51TestPathListDataFileList, err
	}

	tiKuVip51TestPathListDataFileList = tiKuVip51TestPathListResult.Data.FileList
	return tiKuVip51TestPathListDataFileList, nil
}

var downloadNumber = 0
var sleepSecond = 30

func tiKuVip51TestDownloadUrl(tiKuVip51TestPathListDataFileList []TiKuVip51TestPathListDataFileList) {
	for _, pathListDataFile := range tiKuVip51TestPathListDataFileList {
		pathArray := strings.Split(pathListDataFile.Path, "/")
		handlePath := make([]string, 0, len(pathArray)-2)
		for i, v := range pathArray {
			if i != 0 && i != len(pathArray)-1 {
				handlePath = append(handlePath, v)
			}
		}
		// 只保留一级目录
		filePath := "../tikuvip.51test.net/" + handlePath[0] + "/"
		fileName := pathListDataFile.Name
		downloadUrl := fmt.Sprintf("https://tikuvip.51test.net/index.php?pluginApp/to/officeLive/&path={userShare}:100/真题题库%s", pathListDataFile.Path)

		fmt.Println("=======================")
		fmt.Println(downloadUrl)
		fmt.Println(filePath)
		fmt.Println(fileName)

		if _, err := os.Stat(filePath + fileName); err != nil {
			err := downloadTiKuVip51Test(downloadUrl, filePath, fileName)
			downloadNumber++
			if err != nil {
				fmt.Println(err)
				continue
			}
			if downloadNumber >= 15 {
				fmt.Sprintf("=======下载15个文件，暂停%d秒=======", sleepSecond)
				time.Sleep(time.Second * time.Duration(sleepSecond))
				downloadNumber = 0
			} else {
				time.Sleep(time.Second * 1)
			}
		}
	}
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
	req.Header.Set("cookie", "__yjs_duid=1_9906f2c8d7c017db33d48b6a18ccd56b1675842585518; __bid_n=18630008879d6e0aee4207; Hm_lvt_f4ae163e87a012d4ab5106f993decb4c=1678772717; HOST=https://tikuvip.51test.net/; APP_HOST=https://tikuvip.51test.net/; kodUserLanguage=zh-CN; KOD_SESSION_SSO=n8v73kudnk59513468tvhg87tb; KOD_SESSION_ID_8e194=4oqj9kpf2pv95a06ob9jtgn9t9; Hm_lvt_c3d24798f142b815b7069d109e892e45=1681193783; Hm_lpvt_c3d24798f142b815b7069d109e892e45=1681193791; FPTOKEN=UE8f6wb8MckxSxA89FIzTZ9nPQu2jojdRQf4VC8sc1QD/+6ogqCPSaSNzQEejyrOERvpDkcNLTobKcgtrh4HBHntkNvvf/elpIuqc/oTjkvNrkyQSRfyPlZ6jm9dYo3c/48EVuuYIExPhgwAdN3uYxIbH7T3h+s+F/RMB9b7hk8HDyBAJqiZIcMKwKqYgiPDRL8unhH+0FqFBIoJADdvPMAxfThWnvolGsCCpU+jZUPoXZmBCWUf88+amY4wvGNiBbcoIZYh1tZfd2Hd+AF+HWz0RsBdmMc0FHvDBx/mxmPIeZrFww3nE7PP185jycCBL2D4vAVHyVvvBI8R5nCUrC/zG8ya3XsTSI0LHdQsN/Dg8J5Fvjjlmfb//2zgb6jqF9AZvf0CZ/3R5YRO2kkdKw==|wfSc6dBQePRjKadw5Z+bnD3Mhzt+CJbcPPgSMmC9Td4=|10|acd5f2507c1333048f3938686dc63972; Hm_lpvt_f4ae163e87a012d4ab5106f993decb4c=1681195808; kodVersionCheck=check-at-1681195868")
	req.Header.Set("Host", "tikuvip.51test.net")
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
