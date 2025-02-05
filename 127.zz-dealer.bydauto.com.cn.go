package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var EditBydAutoEnableHttpProxy = false
var EditBydAutoHttpProxyUrl = "111.225.152.186:8089"
var EditBydAutoHttpProxyUrlArr = make([]string, 0)

func EditBydAutoHttpProxy() error {
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
					EditBydAutoHttpProxyUrlArr = append(EditBydAutoHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					EditBydAutoHttpProxyUrlArr = append(EditBydAutoHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func EditBydAutoSetHttpProxy() (httpclient *http.Client) {
	if EditBydAutoHttpProxyUrl == "" {
		if len(EditBydAutoHttpProxyUrlArr) <= 0 {
			err := EditBydAutoHttpProxy()
			if err != nil {
				EditBydAutoSetHttpProxy()
			}
		}
		EditBydAutoHttpProxyUrl = EditBydAutoHttpProxyUrlArr[0]
		if len(EditBydAutoHttpProxyUrlArr) >= 2 {
			EditBydAutoHttpProxyUrlArr = EditBydAutoHttpProxyUrlArr[1:]
		} else {
			EditBydAutoHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(EditBydAutoHttpProxyUrl)
	ProxyURL, _ := url.Parse(EditBydAutoHttpProxyUrl)
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
// 李奎丽
// var BydAutoEditAuthorization = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50VHlwZSI6MSwiaWQiOjY0MzIyLCJpc1N1cGVyIjpmYWxzZX0.IiINeGVqTZTqE9zHvACPX__Qu1A9YB4916lMXAumjIc"

// 高宏瑞
var BydAutoEditAuthorization = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50VHlwZSI6MSwiaWQiOjExMTA2NCwiaXNTdXBlciI6ZmFsc2V9.17zcz8xR6-cOP8OZjDUcAOoAYe2imAAKxi7vNc66PDc"
var BydAutoEditNextPageSleep = 10

type QueryEditBydAutoResponseList struct {
	Data  []QueryEditBydAutoResponseListData `json:"data"`
	File  string                             `json:"file"`
	Total int                                `json:"total"`
}

type QueryEditBydAutoResponseListData struct {
	ActivityDate   int    `json:"activityDate"`
	ActivityType   int    `json:"activityType"`
	ComeCount      int    `json:"comeCount"`
	Content        string `json:"content"`
	CustomerId     int    `json:"customerId"`
	CustomerMobile string `json:"customerMobile"`
	CustomerName   string `json:"customerName"`
	FromSource     string `json:"fromSource"`
	FromType       int    `json:"fromType"`
	IsDelay        bool   `json:"isDelay"`
	IsValid        bool   `json:"isValid"`
	Level          string `json:"level"`
	OwnerName      string `json:"ownerName"`
	SeriesName     string `json:"seriesName"`
	Source         string `json:"source"`
	SourceIdentify string `json:"sourceIdentify"`
	Status         int    `json:"status"`
}

type QueryEditBydAutoResponseGet struct {
	Info            QueryEditBydAutoResponseGetInfo            `json:"info"`
	PendingActivity QueryEditBydAutoResponseGetPendingActivity `json:"pendingActivity"`
}

type QueryEditBydAutoResponseGetInfo struct {
	CustomerId     int                                            `json:"customerId"`
	DealId         int                                            `json:"dealId"`
	IntendSerieses []QueryEditBydAutoResponseGetInfoIntendSeriese `json:"intendSerieses"`
}

type QueryEditBydAutoResponseGetInfoIntendSeriese struct {
	Id       int    `json:"id"`
	IsMaster bool   `json:"isMaster"`
	Name     string `json:"name"`
}

type QueryEditBydAutoResponseGetPendingActivity struct {
	Id int `json:"id"`
}

type QueryEditBydAutoResponseFollow struct {
	Success bool `json:"success"`
}

// ychEduSpider 编辑智蛛AI经销商系统
// @Title 编辑智蛛AI经销商系统
// @Description https://zz-dealer.bydauto.com.cn/，编辑智蛛AI经销商系统
func main() {
	pageCount := 10
	curPage := 0
	isPageListGo := true
	for isPageListGo {
		// 当前页是否处理过文档
		hasEditFlag := false
		dealerId := 826
		listRequestPayload := make(map[string]interface{})
		listRequestPayload["activityType"] = 0
		listRequestPayload["dateEnd"] = 0
		listRequestPayload["dateStart"] = 0
		listRequestPayload["dealerId"] = dealerId
		listRequestPayload["filterType"] = 0
		listRequestPayload["fromType"] = 0
		listRequestPayload["key"] = ""
		listRequestPayload["level"] = ""
		listRequestPayload["onlyTotal"] = false
		listRequestPayload["pageCount"] = pageCount
		listRequestPayload["pageStart"] = curPage
		listRequestPayload["saleIds"] = ""
		listRequestPayload["seriesIds"] = ""
		pageListUrl := "https://zz-api.bydauto.com.cn/aiApi-dealer/v1/taskRpc/list"
		fmt.Println(pageListUrl)
		queryEditBydAutoResponseList, err := QueryEditBydAutoList(pageListUrl, listRequestPayload)
		if err != nil {
			EditBydAutoHttpProxyUrl = ""
			fmt.Println(err)
			continue
		}
		//fmt.Printf("%+v", queryEditBydAutoResponseList.Data)
		//os.Exit(1)
		if len(queryEditBydAutoResponseList.Data) <= 0 {
			break
		}
		for _, queryEditBydAutoResponseListData := range queryEditBydAutoResponseList.Data {
			fmt.Printf("客户姓名：%s，客户手机号：%s\n", queryEditBydAutoResponseListData.CustomerName, queryEditBydAutoResponseListData.CustomerMobile)
			// 类型1：沟通 类型8：回访，只处理待沟通类型
			if queryEditBydAutoResponseListData.ActivityType == 1 {
				fmt.Println("====================开始处理数据================================")
				// 将待沟通的类型，处理时间随机延长2-10天
				randIntN := rand.Intn(5)
				if randIntN == 0{
				    randIntN = 1
				}
                // randIntN默认值是1-5
                randIntN = randIntN*(2*24*60*60)
				activityDate := queryEditBydAutoResponseListData.ActivityDate + randIntN + 1
				customerId := queryEditBydAutoResponseListData.CustomerId
				getUrl := "https://zz-api.bydauto.com.cn/aiApi-dealer/v2/appCustomerService/get"
				fmt.Println(getUrl)
				getRequestPayload := make(map[string]interface{})
				getRequestPayload["customerId"] = customerId
				getRequestPayload["dealerId"] = dealerId
				queryEditBydAutoResponseGet, err := QueryEditBydAutoGet(getUrl, getRequestPayload)
				//fmt.Printf("%+v\n", queryEditBydAutoResponseGet)
				if err != nil {
					EditBydAutoHttpProxyUrl = ""
					fmt.Println(err)
					continue
				}

				followUrl := "https://zz-api.bydauto.com.cn/aiApi-dealer/v2/appCustomerService/follow"
				followRequestPayload := make(map[string]interface{})
				followRequestPayload["activityId"] = queryEditBydAutoResponseGet.PendingActivity.Id
				followRequestPayload["competitionSerieses"] = ""
				followRequestPayload["completedActivityType"] = 1
				followRequestPayload["customerId"] = queryEditBydAutoResponseGet.Info.CustomerId
				followRequestPayload["dealId"] = queryEditBydAutoResponseGet.Info.DealId
				followRequestPayload["dealerId"] = dealerId
				followRequestPayload["event"] = 1
				followRequestPayload["failCity"] = ""
				followRequestPayload["failProvince"] = ""
				followRequestPayload["isSuspend"] = false
				followRequestPayload["level"] = "A"
				followRequestPayload["nextActivityType"] = 1
				followRequestPayload["nextDate"] = activityDate
				followRequestPayload["nextTestSpecId"] = nil
				followRequestPayload["quote"] = 0
				// 备注内容随机
				// 初始化随机数生成器
				rand.Seed(time.Now().UnixNano())
				// 定义一个数组
				array := []string{"已跟进", "已联系", "考虑考虑", "不着急"}
				// 获取随机索引
				index := rand.Intn(len(array))
				followRequestPayload["remark"] = array[index]
				followRequestPayload["reserveComeDate"] = 0
				var specIds []string
				for _, intendSeriese := range queryEditBydAutoResponseGet.Info.IntendSerieses {
					specIds = append(specIds, strconv.Itoa(intendSeriese.Id))
				}
				followRequestPayload["specIds"] = strings.Join(specIds, ",")
				followRequestPayload["testResult"] = 1
				followRequestPayload["testSpecId"] = 0
				//fmt.Printf("%+v\n", followRequestPayload)

				_, err = QueryEditBydAutoFollow(followUrl, followRequestPayload)
				if err != nil {
					EditBydAutoHttpProxyUrl = ""
					fmt.Println(err)
					continue
				}
				// 当前页是否处理过文档---处理过文档
				hasEditFlag = true
// 				bydAutoEditSaveTimeSleep := rand.Intn(3)
				bydAutoEditSaveTimeSleep := 15
				for i := 1; i <= bydAutoEditSaveTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page="+strconv.Itoa(curPage)+"===========更新", queryEditBydAutoResponseListData.CustomerName, "成功，暂停", bydAutoEditSaveTimeSleep, "秒，倒计时", i, "秒===========")
				}
				fmt.Println("====================处理数据完成================================")
			}
		}
		// 如果当前页没有处理过文档，则请求下一页，如果处理过文档，则继续请求当前分页
		if hasEditFlag == false {
			curPage++
			if curPage > (queryEditBydAutoResponseList.Total/pageCount)+1 {
				fmt.Println("没有更多分页了")
				isPageListGo = false
				curPage = 0
				break
			}
		}
		for i := 1; i <= BydAutoEditNextPageSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("===========翻", curPage, "页，暂停", BydAutoEditNextPageSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

func QueryEditBydAutoList(requestUrl string, listRequestPayload map[string]interface{}) (queryEditBydAutoResponseList QueryEditBydAutoResponseList, err error) {
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
	if EditBydAutoEnableHttpProxy {
		client = EditBydAutoSetHttpProxy()
	}
	payloadBytes, err := json.Marshal(listRequestPayload)
	if err != nil {
		return queryEditBydAutoResponseList, err
	}
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(payloadBytes)) //建立连接
	if err != nil {
		return queryEditBydAutoResponseList, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Authorization", BydAutoEditAuthorization)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payloadBytes)))
	req.Header.Set("Host", "zz-api.bydauto.com.cn")
	req.Header.Set("Origin", "https://zz-dealer.bydauto.com.cn")
	req.Header.Set("Referer", "https://zz-dealer.bydauto.com.cn/")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return queryEditBydAutoResponseList, err
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return queryEditBydAutoResponseList, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return queryEditBydAutoResponseList, err
	}
	err = json.Unmarshal(respBytes, &queryEditBydAutoResponseList)
	if err != nil {
		return queryEditBydAutoResponseList, err
	}
	return queryEditBydAutoResponseList, nil
}

func QueryEditBydAutoGet(requestUrl string, getRequestPayload map[string]interface{}) (queryEditBydAutoResponseGet QueryEditBydAutoResponseGet, err error) {
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
	if EditBydAutoEnableHttpProxy {
		client = EditBydAutoSetHttpProxy()
	}
	payloadBytes, err := json.Marshal(getRequestPayload)
	if err != nil {
		return queryEditBydAutoResponseGet, err
	}
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(payloadBytes)) //建立连接
	if err != nil {
		return queryEditBydAutoResponseGet, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Authorization", BydAutoEditAuthorization)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payloadBytes)))
	req.Header.Set("Host", "zz-api.bydauto.com.cn")
	req.Header.Set("Origin", "https://zz-dealer.bydauto.com.cn")
	req.Header.Set("Referer", "https://zz-dealer.bydauto.com.cn/")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return queryEditBydAutoResponseGet, err
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return queryEditBydAutoResponseGet, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return queryEditBydAutoResponseGet, err
	}
	err = json.Unmarshal(respBytes, &queryEditBydAutoResponseGet)
	if err != nil {
		return queryEditBydAutoResponseGet, err
	}
	return queryEditBydAutoResponseGet, nil
}

func QueryEditBydAutoFollow(requestUrl string, followRequestPayload map[string]interface{}) (queryEditBydAutoResponseFollow QueryEditBydAutoResponseFollow, err error) {
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
	if EditBydAutoEnableHttpProxy {
		client = EditBydAutoSetHttpProxy()
	}
	payloadBytes, err := json.Marshal(followRequestPayload)
	if err != nil {
		return queryEditBydAutoResponseFollow, err
	}
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(payloadBytes)) //建立连接
	if err != nil {
		return queryEditBydAutoResponseFollow, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Authorization", BydAutoEditAuthorization)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payloadBytes)))
	req.Header.Set("Host", "zz-api.bydauto.com.cn")
	req.Header.Set("Origin", "https://zz-dealer.bydauto.com.cn")
	req.Header.Set("Referer", "https://zz-dealer.bydauto.com.cn/")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return queryEditBydAutoResponseFollow, err
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return queryEditBydAutoResponseFollow, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return queryEditBydAutoResponseFollow, err
	}
	err = json.Unmarshal(respBytes, &queryEditBydAutoResponseFollow)
	if err != nil {
		return queryEditBydAutoResponseFollow, err
	}
	return queryEditBydAutoResponseFollow, nil
}
