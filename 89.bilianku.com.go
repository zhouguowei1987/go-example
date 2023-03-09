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
	BiLianKuEnableHttpProxy = false
	BiLianKuHttpProxyUrl    = "111.225.152.186:8089"
)

func BiLianKuSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(BiLianKuHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type BiLianKuSubject struct {
	name     string
	category []BiLianKuCategory
}

type BiLianKuCategory struct {
	name string
	id   int
}

var BiLianKuAllSubject = []BiLianKuSubject{
	{
		name: "学历类",
		category: []BiLianKuCategory{
			{
				name: "自考",
				id:   11,
			},
			{
				name: "成考",
				id:   12,
			},
			{
				name: "升学考试",
				id:   13,
			},
			{
				name: "研究生考试",
				id:   14,
			},
		},
	},
	{
		name: "职业资格",
		category: []BiLianKuCategory{
			{
				name: "从业资格考试",
				id:   15,
			},
			{
				name: "新闻媒体考试",
				id:   16,
			},
			{
				name: " 演出行业考试",
				id:   17,
			},
			{
				name: "其它",
				id:   18,
			},
		},
	},
	{
		name: "公务员",
		category: []BiLianKuCategory{
			{
				name: "公务员考试",
				id:   19,
			},
			{
				name: "警察考试",
				id:   20,
			},
			{
				name: "企业事业单位考试",
				id:   21,
			},
			{
				name: "其它公务员考试",
				id:   22,
			},
		},
	},
	{
		name: "医卫类",
		category: []BiLianKuCategory{
			{
				name: "药学考试",
				id:   23,
			},
			{
				name: "医师考试",
				id:   24,
			},
			{
				name: "护理考试",
				id:   25,
			},
			{
				name: "医技职称",
				id:   26,
			},
			{
				name: "其它医卫类",
				id:   27,
			},
		},
	},
	{
		name: "建筑工程",
		category: []BiLianKuCategory{
			{
				name: "建造师",
				id:   28,
			},
			{
				name: "注册工程师",
				id:   29,
			},
			{
				name: "工程师考试",
				id:   30,
			},
			{
				name: "八大员",
				id:   31,
			},
		},
	},
	{
		name: "外语类",
		category: []BiLianKuCategory{
			{
				name: "大学英语",
				id:   32,
			},
			{
				name: "英语专业考试",
				id:   33,
			},
			{
				name: "成人英语",
				id:   34,
			},
			{
				name: "商务英语",
				id:   35,
			},
			{
				name: "雅思",
				id:   36,
			},
		},
	},
	{
		name: "外贸类",
		category: []BiLianKuCategory{
			{
				name: "货运代理",
				id:   37,
			},
			{
				name: "单证员",
				id:   38,
			},
			{
				name: "跟单员",
				id:   39,
			},
		},
	},
	{
		name: "计算机类",
		category: []BiLianKuCategory{
			{
				name: "等级考试",
				id:   40,
			},
			{
				name: "软考",
				id:   41,
			},
			{
				name: "其它",
				id:   42,
			},
		},
	},
	{
		name: "财会类",
		category: []BiLianKuCategory{
			{
				name: "会计考试",
				id:   43,
			},
			{
				name: "经济师考试",
				id:   44,
			},
			{
				name: "审计师考试",
				id:   45,
			},
			{
				name: "统计师",
				id:   46,
			},
		},
	},
	{
		name: "技能鉴定",
		category: []BiLianKuCategory{
			{
				name: "特种作业",
				id:   47,
			},
			{
				name: "安全生产",
				id:   48,
			},
			{
				name: "医疗药物",
				id:   49,
			},
			{
				name: "交通相关",
				id:   50,
			},
			{
				name: "科技类",
				id:   51,
			},
			{
				name: "食品类",
				id:   52,
			},
			{
				name: "服务行业",
				id:   53,
			},
			{
				name: "健康护理",
				id:   54,
			},
			{
				name: "其他类",
				id:   55,
			},
		},
	},
}

type apiBiLianKuResult struct {
	Count apiBiLianKuResultCount  `json:"count"`
	Data  []apiBiLianKuResultData `json:"data"`
	Multi apiBiLianKuResultMulti  `json:"multi"`
}

type apiBiLianKuResultCount struct {
	Num int `json:"num"`
}

type apiBiLianKuResultData struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	InputTime   string `json:"inputtime"`
	Time        string `json:"time"`
	Hot         string `json:"hot"`
	DownMianFei string `json:"downmianfei"`
	DowXueBi    string `json:"dowxuebi"`
}
type apiBiLianKuResultMulti struct {
	BackPage  int         `json:"BackPage"`
	NextPage  int         `json:"NextPage"`
	Page      int         `json:"Page"`
	PageCount int         `json:"PageCount"`
	PageNums  map[int]int `json:"PageNums"`
	PageSize  int         `json:"PageSize"`
	RecordNum int         `json:"RecordNum"`
}

// ychEduSpider 获取必练库文档
// @Title 获取必练库文档
// @Description https://bilianku.com/，获取必练库文档
func main() {
	for _, subject := range BiLianKuAllSubject {
		for _, category := range subject.category {
			page := 1
			apiUrl := fmt.Sprintf("http://bilianku.com/api.php?op=ajax&t=%d", time.Now().UnixNano()/1e6)
			level := 2
			sort := category.id
			for true {
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
				if BiLianKuEnableHttpProxy {
					client = BiLianKuSetHttpProxy()
				}
				postData := url.Values{}
				postData.Add("action", "kulists")
				postData.Add("t", strconv.Itoa(2))
				postData.Add("page", strconv.Itoa(page))
				postData.Add("level", strconv.Itoa(level))
				postData.Add("sort", strconv.Itoa(sort))
				postData.Add("type", "all")
				req, err := http.NewRequest("POST", apiUrl, strings.NewReader(postData.Encode())) //建立连接
				if err != nil {
					fmt.Println(err)
					page = 1
					break
				}
				req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
				req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
				req.Header.Set("Connection", "keep-alive")
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.Header.Set("Cookie", "Hm_lvt_46575bd9537c0f2a1ba086e74e274fda=1678264071; _gid=GA1.2.2071830305.1678264073; _gat=1; Leslie_1900_LoginID=N2M0NjU0MWEzZDM1NWNlNTI0MDRmZDg0OTUxY2RiMzc1OTBiNjk=; Leslie_1900_UID=M2NlNWUxODc1NzhhYzY3ZGUzN2UyYTY0ZTgxZDRiY2YzNjY3NDc=; Leslie_1900_KSUID=OWYzYmRiZGE3NGVjZjIwNzQzZTkwMTU1ZWRjZTc0NWYxNWQ1NDQ=; Leslie_1900_KSLoginID=NTdhNmRkZDAxNmYzZTMwYzU2NzdkNTYwNzBhZTBlOTg1OGYyOWM=; Hm_lpvt_46575bd9537c0f2a1ba086e74e274fda=1678324285; _ga_34B604LFFQ=GS1.1.1678324020.2.1.1678324286.53.0.0; _ga=GA1.1.790663054.1678264073")
				req.Header.Set("Host", "bilianku.com")
				req.Header.Set("Origin", "https://bilianku.com")
				req.Header.Set("Referer", "http://bilianku.com/shijuan_/"+strconv.Itoa(sort))
				req.Header.Set("Sec-Fetch-User", "?1")
				req.Header.Set("Upgrade-Insecure-Requests", "1")
				req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
				req.Header.Set("X-Requested-With", "XMLHttpRequest")
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
					page = 1
					break
				}
				apiBiLianKuResult := &apiBiLianKuResult{}
				err = json.Unmarshal(respBytes, apiBiLianKuResult)
				if err != nil {
					fmt.Println(err)
					page = 1
					break
				}
				for _, paper := range apiBiLianKuResult.Data {
					filePath := "../bilianku.com/" + subject.name + "/" + category.name + "/"
					fmt.Println(filePath)
					paperId, _ := strconv.Atoi(paper.Id)
					paperTtle := paper.Title
					fmt.Println(paperTtle)

					downloadUrl := fmt.Sprintf("http://bilianku.com/down-%d/", paperId)
					detailUrl := fmt.Sprintf("http://bilianku.com/shijuan-%d/", paperId)
					fileName := paperTtle + ".pdf"
					err = downloadBiLianKu(downloadUrl, detailUrl, filePath, fileName)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
				if apiBiLianKuResult.Multi.PageCount > page {
					page++
				} else {
					page = 1
					break
				}
			}
		}
	}
}

func downloadBiLianKu(attachmentUrl string, referer string, filePath string, fileName string) error {
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
	if BiLianKuEnableHttpProxy {
		client = BiLianKuSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("isvip", strconv.Itoa(1))
	req, err := http.NewRequest("POST", attachmentUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "Hm_lvt_46575bd9537c0f2a1ba086e74e274fda=1678264071; _gid=GA1.2.2071830305.1678264073; _gat=1; Leslie_1900_LoginID=N2M0NjU0MWEzZDM1NWNlNTI0MDRmZDg0OTUxY2RiMzc1OTBiNjk=; Leslie_1900_UID=M2NlNWUxODc1NzhhYzY3ZGUzN2UyYTY0ZTgxZDRiY2YzNjY3NDc=; Leslie_1900_KSUID=OWYzYmRiZGE3NGVjZjIwNzQzZTkwMTU1ZWRjZTc0NWYxNWQ1NDQ=; Leslie_1900_KSLoginID=NTdhNmRkZDAxNmYzZTMwYzU2NzdkNTYwNzBhZTBlOTg1OGYyOWM=; Hm_lpvt_46575bd9537c0f2a1ba086e74e274fda=1678324285; _ga_34B604LFFQ=GS1.1.1678324020.2.1.1678324286.53.0.0; _ga=GA1.1.790663054.1678264073")
	req.Header.Set("Host", "bilianku.com")
	req.Header.Set("Origin", "https://bilianku.com")
	req.Header.Set("Referer", referer)
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
