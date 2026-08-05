package main

import (
	"bytes"
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
	_ "os"
	"path"
	"path/filepath"
	_ "path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	_ "golang.org/x/net/html"
)

var JiangYuanWangEnableHttpProxy = false
var JiangYuanWangHttpProxyUrl = "111.225.152.186:8089"
var JiangYuanWangHttpProxyUrlArr = make([]string, 0)

func JiangYuanWangHttpProxy() error {
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
					JiangYuanWangHttpProxyUrlArr = append(JiangYuanWangHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					JiangYuanWangHttpProxyUrlArr = append(JiangYuanWangHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func JiangYuanWangSetHttpProxy() (httpclient *http.Client) {
	if JiangYuanWangHttpProxyUrl == "" {
		if len(JiangYuanWangHttpProxyUrlArr) <= 0 {
			err := JiangYuanWangHttpProxy()
			if err != nil {
				JiangYuanWangSetHttpProxy()
			}
		}
		JiangYuanWangHttpProxyUrl = JiangYuanWangHttpProxyUrlArr[0]
		if len(JiangYuanWangHttpProxyUrlArr) >= 2 {
			JiangYuanWangHttpProxyUrlArr = JiangYuanWangHttpProxyUrlArr[1:]
		} else {
			JiangYuanWangHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(JiangYuanWangHttpProxyUrl)
	ProxyURL, _ := url.Parse(JiangYuanWangHttpProxyUrl)
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

type JiangYuanWangCategory struct {
	Name             string
	CateId           int
	Page             int
	MaxPage          int
	ParentClassifyId string
}

var jiangYuanWangCategory = []JiangYuanWangCategory{
	{
		Name:             "PPT",
		CateId:           5,
		Page:             1,
		MaxPage:          65,
		ParentClassifyId: "1,3,5,6",
	},
	// {
	// 	Name:             "Word",
	// 	CateId:           6,
	// 	Page:             1,
	// 	MaxPage:          60,
	// 	ParentClassifyId: "1,3,94",
	// },
	// {
	// 	Name:             "Execl",
	// 	CateId:           7,
	// 	Page:             1,
	// 	MaxPage:          35,
	// 	ParentClassifyId: "1,3",
	// },
	// {
	// 	Name:             "视频",
	// 	CateId:           8,
	// 	Page:             1,
	// 	MaxPage:          10,
	// 	ParentClassifyId: "1,3,232",
	// },
	// {
	// 	Name:             "配乐",
	// 	CateId:           9,
	// 	Page:             1,
	// 	MaxPage:          37,
	// 	ParentClassifyId: "1,5,232,237,247",
	// },
	// {
	// 	Name:             "音频",
	// 	CateId:           10,
	// 	Page:             1,
	// 	MaxPage:          88,
	// 	ParentClassifyId: "1,5,232",
	// },
}

type QueryJiangYuanWangListRequestPayload struct {
	CateId        int                                                 `json:"cate_id"`
	CategoryValue []QueryJiangYuanWangListRequestPayloadCategoryValue `json:"categoryValue"`
	Order         string                                              `json:"order"`
	Page          int                                                 `json:"page"`
	PageSize      int                                                 `json:"page_size"`
	SearchText    string                                              `json:"search_text"`
	Sort          string                                              `json:"sort"`
}

type QueryJiangYuanWangListRequestPayloadCategoryValue struct {
	ParentClassifyId int      `json:"parentClassifyId"`
	Data             []string `json:"data"`
}

// 下载江源网文档
// @Title 下载江源网文档
// @Description https://www.jiangyuanwang.com/，下载江源网文档
func main() {
	for _, category := range jiangYuanWangCategory {
		pageListUrl := "https://jywserver.jiangxiatech.com/resources/home/searchList"
		fmt.Println(pageListUrl)
		page := category.Page
		maxPage := category.MaxPage
		rows := 100
		// 处理CategoryValue字段
		var categoryValue []QueryJiangYuanWangListRequestPayloadCategoryValue
		parentClassifyIdStr := category.ParentClassifyId
		parentClassifyIdArr := strings.Split(parentClassifyIdStr, ",")
		for _, parentClassifyIdStr := range parentClassifyIdArr {
			parentClassifyId, _ := strconv.Atoi(parentClassifyIdStr)
			categoryValue = append(categoryValue, QueryJiangYuanWangListRequestPayloadCategoryValue{
				ParentClassifyId: parentClassifyId,
				Data:             []string{},
			})
		}
		isPageListGo := true
		for isPageListGo {
			if page > maxPage {
				isPageListGo = false
				break
			}
			queryJiangYuanWangListRequestPayload := QueryJiangYuanWangListRequestPayload{
				CateId:        category.CateId,
				CategoryValue: categoryValue,
				Order:         "desc",
				Page:          page,
				PageSize:      rows,
				SearchText:    "",
				Sort:          "update_time",
			}
			queryJiangYuanWangListResponseResultRows, err := QueryJiangYuanWangList(pageListUrl, queryJiangYuanWangListRequestPayload)
			if err != nil {
				fmt.Println(err)
				break
			}
			for _, jiangYuanWang := range queryJiangYuanWangListResponseResultRows {
				fmt.Println("===========Name=", category.Name, "CateId=", category.CateId, " Page=", page, "===========")

				if jiangYuanWang.Address == "" {
					fmt.Println("文档没有附件地址，跳过")
					continue
				}

				name := jiangYuanWang.Name
				name = strings.TrimSpace(name)
				name = strings.ToLower(name)
				name = strings.ReplaceAll(name, " ", "-")
				name = strings.ReplaceAll(name, "　", "-")
				name = strings.ReplaceAll(name, "/", "-")
				name = strings.ReplaceAll(name, "--", "-")
				fmt.Println(name)

				cate_name := jiangYuanWang.CateName
				cate_name = strings.TrimSpace(cate_name)
				cate_name = strings.ToLower(cate_name)
				title := name
				// 查看name字段是否含有cate_name
				if strings.Index(title, cate_name) == -1 {
					title = name + "(" + cate_name + ")"
				}

				downloadUrl := jiangYuanWang.Address
				fileExt := path.Ext(downloadUrl)

				filePath := "../www.jiangyuanwang.com/www.jiangyuanwang.com/" + cate_name + "/" + title + fileExt
				fmt.Println(filePath)

				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				fmt.Println("=======开始下载========")

				requestJiangYuanWangDownloadRefererUrl := fmt.Sprintf("https://www.jiangyuanwang.com/ClassificationTemplate/TemplateDetails/%d?classId=%d", jiangYuanWang.Id, category.CateId)

				err = downloadJiangYuanWang(downloadUrl, requestJiangYuanWangDownloadRefererUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//复制文件
				// tempFilePath := "../www.jiangyuanwang.com/temp-www.jiangyuanwang.com/" + cate_name + "/" + title + fileExt
				// err = copyJiangYuanWangFile(filePath, tempFilePath)
				// if err != nil {
				// 	fmt.Println(err)
				// 	continue
				// }
				fmt.Println("=======下载完成========")
				//DownLoadJiangYuanWangTimeSleep := 10
				DownLoadJiangYuanWangTimeSleep := rand.Intn(10)
				for i := 1; i <= DownLoadJiangYuanWangTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("===========title=", title, "===========下载成功 暂停", DownLoadJiangYuanWangTimeSleep, "秒 倒计时", i, "秒===========")
				}
			}
			page++
			if page > maxPage {
				isPageListGo = false
				break
			}
			DownLoadJiangYuanWangPageTimeSleep := 10
			// DownLoadJiangYuanWangPageTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadJiangYuanWangPageTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("===========cate_name = ", category.Name, " page = "+strconv.Itoa(page)+"===== 暂停", DownLoadJiangYuanWangPageTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
	}
}

type QueryJiangYuanWangListResponse struct {
	Code   int                                  `json:"code"`
	Result QueryJiangYuanWangListResponseResult `json:"result"`
	Msg    string                               `json:"msg"`
}

type QueryJiangYuanWangListResponseResult struct {
	count int                                        `json:"count"`
	Rows  []QueryJiangYuanWangListResponseResultRows `json:"rows"`
}

type QueryJiangYuanWangListResponseResultRows struct {
	Address  string `json:"address"`
	CateName string `json:"cate_name"`
	Name     string `json:"name"`
	Id       int    `json:"id"`
}

func QueryJiangYuanWangList(requestUrl string, queryJiangYuanWangListRequestPayload QueryJiangYuanWangListRequestPayload) (queryJiangYuanWangListResponseResultRows []QueryJiangYuanWangListResponseResultRows, err error) {
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
	if JiangYuanWangEnableHttpProxy {
		client = JiangYuanWangSetHttpProxy()
	}
	// 将数据编码为JSON格式
	queryJiangYuanWangListRequestPayloadJson, err := json.Marshal(queryJiangYuanWangListRequestPayload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 创建字符串读取器，这是http.Request需要的类型
	body := bytes.NewReader(queryJiangYuanWangListRequestPayloadJson)
	req, err := http.NewRequest("POST", requestUrl, body) //建立连接

	queryJiangYuanWangListResponse := QueryJiangYuanWangListResponse{}
	if err != nil {
		return queryJiangYuanWangListResponseResultRows, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Authorization", "test1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Host", "www.jiangyuanwang.com")
	req.Header.Set("Origin", "https://www.jiangyuanwang.com")
	req.Header.Set("Referer", "https://www.jiangyuanwang.com/")
	req.Header.Set("Tenant-Id", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryJiangYuanWangListResponseResultRows, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryJiangYuanWangListResponseResultRows, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryJiangYuanWangListResponseResultRows, err
	}
	err = json.Unmarshal(respBytes, &queryJiangYuanWangListResponse)
	if err != nil {
		return queryJiangYuanWangListResponseResultRows, err
	}
	queryJiangYuanWangListResponseResultRows = queryJiangYuanWangListResponse.Result.Rows
	return queryJiangYuanWangListResponseResultRows, nil
}

func downloadJiangYuanWang(attachmentUrl string, referer string, filePath string) error {
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
	if JiangYuanWangEnableHttpProxy {
		client = JiangYuanWangSetHttpProxy()
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
	req.Header.Set("Host", "resources.jiangxiatech.com")
	req.Header.Set("Referer", referer)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
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

func copyJiangYuanWangFile(src, dst string) (err error) {
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
