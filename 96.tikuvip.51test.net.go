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

var tiKuVipCookie = "__yjs_duid=1_9906f2c8d7c017db33d48b6a18ccd56b1675842585518; Hm_lvt_c3d24798f142b815b7069d109e892e45=1685671841; Hm_lpvt_c3d24798f142b815b7069d109e892e45=1685671841; __bid_n=18630008879d6e0aee4207; FPTOKEN=MGPTS+u3bngAqtwpnIUwMibRc1guZDiPU0y/2Ry3RnRwNpndpk8N5snuF0RDdLI6D7w8f5NXPlGBl14aAisNEV6A8iqHUQqR1EwEnrGn0sY+vFM7G6zIHTnvcpseJEolkFkvl7mQTyW5oghJsJns52laUwjY4Xwb3Ag54dTqGcjOgCTAY5JrqEIcPNQZE1j1+Pzwqu1X2jS2S8l2NQ9aTFwgIZQVTqzAY6AOna96FkqU/CvWAfQ12YXMJz9RncwRRP03odS/0u8JwxVAG+TlkJzeTdCimks1nzllXuTPQ0paE9/OS02kJEYFQRjxCcjeUJ3u98zwxWreQbiXuC11Bdxf5UNnooQHzp6zCEn9joxYb+QG6ftQ9R6pJjrjqthOLK8pNMcXuQY0SJakvg0Nwg==|MJv40hOMa4u5wYnIcPFTUVWo07Na1ehPiHJz0gCwQOk=|10|62703f260015d5c13013d6e6360b5596; Hm_lvt_f4ae163e87a012d4ab5106f993decb4c=1688434770; Hm_lvt_f0e5d9e2cfc9b0e39766d14a6cbd7c33=1688526341; Hm_lpvt_f0e5d9e2cfc9b0e39766d14a6cbd7c33=1688526341; Hm_lpvt_f4ae163e87a012d4ab5106f993decb4c=1689610425; KOD_SESSION_ID_8e194=dmr4nju6okomkthdpp9gg9m51c; PHPSESSID=26sqogif7l94h6gvitvhs745jh"

// ychEduSpider 获取无忧考试网真题
// @Title 获取无忧考试网真题
// @Description https://tikuvip.51test.net/，获取无忧考试网真题
func main() {
	saveCategory := map[string]bool{
		"自考":      true,
		"专升本考试":   true,
		"一级建造师考试": true,
		"小升初":     true,
		"考研":      true,
		"高中会考":    true,
		"高考":      true,
		"二级建造师考试": true,
		"成人高考":    true,
		"中考":      true,
	}
	tiKuVip51TestTreeListInitData, err := treeListInit()
	if err != nil {
		fmt.Println(err)
	}
	for _, treeListInitData := range tiKuVip51TestTreeListInitData {
		fmt.Println(treeListInitData.Name)
		for _, childTreeListInitData := range treeListInitData.Children {
			fmt.Println(childTreeListInitData.Name)
			if _, ok := saveCategory[childTreeListInitData.Name]; !ok {
				continue
			}
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
	//req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", tiKuVipCookie)
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
	//req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Cookie", tiKuVipCookie)
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
	//req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", tiKuVipCookie)
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

var maxDownloadNumber = 20
var downloadNumber = 0
var sleepSecond = 30

var tiKuVip51TestSaveYear = []string{"2023"}

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
		filePath := "../tikuvip（2023）.51test.net/" + handlePath[0] + "/"
		fileName := pathListDataFile.Name
		fileName = strings.Trim(fileName, " ")
		for _, year := range tiKuVip51TestSaveYear {
			if strings.Contains(fileName, year) {
				// 开始下载
				attachmentUrl := fmt.Sprintf("https://tikuvip.51test.net/index.php?pluginApp/to/officeLive/&path={userShare}:100/真题题库%s", pathListDataFile.Path)

				fmt.Println("=======================")
				fmt.Println(attachmentUrl)

				downloadDocUrl, err := downloadTiKuVip51TestUrl(attachmentUrl)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(downloadDocUrl)

				if _, err := os.Stat(filePath + fileName); err != nil {
					fmt.Println("=======开始下载========")
					fmt.Println(filePath)
					fmt.Println(fileName)
					err := downloadTiKuVip51Test(downloadDocUrl, filePath, fileName)
					downloadNumber++
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("=======开始完成========")
					if downloadNumber >= maxDownloadNumber {
						fmt.Printf("=======下载%d个文件，暂停%d秒=======\n", maxDownloadNumber, sleepSecond)
						time.Sleep(time.Second * time.Duration(sleepSecond))
						downloadNumber = 0
					} else {
						time.Sleep(time.Second * 1)
					}
				}
				break
			}
		}
	}
}

func downloadTiKuVip51TestUrl(attachmentUrl string) (downloadUrl string, err error) {
	// 初始化客户端
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return downloadUrl, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "zip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", tiKuVipCookie)
	req.Header.Set("Host", "tikuvip.51test.net")
	req.Header.Set("Origin", "https://tikuvip.51test.net")
	req.Header.Set("Referer", "https://tikuvip.51test.net/")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return downloadUrl, err
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode == http.StatusOK {
		downloadUrl = attachmentUrl
	} else if resp.StatusCode == http.StatusFound {
		downloadUrl = resp.Header.Get("Location")
	}
	return downloadUrl, nil
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

	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "zip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", tiKuVipCookie)
	req.Header.Set("Host", "tikuvip.51test.net")
	req.Header.Set("Origin", "https://tikuvip.51test.net")
	req.Header.Set("Referer", "https://tikuvip.51test.net/")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
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
