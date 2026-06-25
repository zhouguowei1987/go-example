package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	OPenStdEnableHttpProxy = false
	OPenStdHttpProxyUrl    = "111.225.152.186:8089"
)

func OPenStdSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(OPenStdHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type StdCategory struct {
	StdName string
	StdUrl  string
}

// ychEduSpider 获取国家标准文档
// @Title 获取国家标准文档
// @Description https://openstd.samr.gov.cn/，获取国家标准文档
func main() {
	var StdCategories = []StdCategory{
		{
			StdName: "强制性国家标准",
			StdUrl:  "https://openstd.samr.gov.cn/bzgk/gb/std_list_type?p.p1=1&p.p90=circulation_date&p.p91=desc",
		},
		{
			StdName: "推荐性国家标准",
			StdUrl:  "https://openstd.samr.gov.cn/bzgk/gb/std_list_type?p.p1=2&p.p90=circulation_date&p.p91=desc",
		},
		{
			StdName: "指导性技术文件",
			StdUrl:  "https://openstd.samr.gov.cn/bzgk/gb/std_list_type?p.p1=3&p.p90=circulation_date&p.p91=desc",
		},
	}
	for _, std := range StdCategories {
		page := 1
		pageSize := 10
		for {
			pageListUrl := fmt.Sprintf(std.StdUrl+"&r=0.20175458803007884&page=%d&pageSize=%d", page, pageSize)
			fmt.Println(pageListUrl)
			pageListDoc, err := htmlquery.LoadURL(pageListUrl)
			if err != nil {
				page = 1
				fmt.Println(err)
				break
			}
			trNodes := htmlquery.Find(pageListDoc, `//table[@class="table result_list table-striped table-hover"]/tbody[2]/tr`)
			if len(trNodes) <= 0 {
				page = 1
				break
			}
			for _, trNode := range trNodes {
				StdNoA := htmlquery.FindOne(trNode, `./td[2]/a`)
				StdNo := htmlquery.InnerText(StdNoA)
				fmt.Println(StdNo)

				StdNameA := htmlquery.FindOne(trNode, `./td[4]/a`)
				StdName := htmlquery.InnerText(StdNameA)
				fmt.Println(StdName)

				HCno := htmlquery.SelectAttr(StdNameA, "onclick")
				HCno = HCno[10 : len(HCno)-3]
				fmt.Println(HCno)

				// 详情URL
				detailUrl := fmt.Sprintf("https://openstd.samr.gov.cn/bzgk/std/newGbInfo?hcno=%s", HCno)
				fmt.Println(detailUrl)

				// 下载文档URL
				downLoadUrl := fmt.Sprintf("https://openstd.samr.gov.cn/bzgk/std/viewGb?hcno=%s", HCno)
				fmt.Println(downLoadUrl)

				refererUrl := fmt.Sprintf("https://openstd.samr.gov.cn/bzgk/std/showGb?type=download&hcno=%s&request_locale=zh", HCno)
				fmt.Println(refererUrl)

				filePath := "../openstd.samr.gov.cn/" + StdName + "(" + StdNo + ")" + ".pdf"
				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}
				fmt.Println("=======开始下载========")
                err = downloadOPenStd(downLoadUrl, refererUrl, filePath)
                if err != nil {
                    fmt.Println(err)
                    continue
                }
                fmt.Println("=======开始完成========")
				// 设置倒计时
				DownLoadOPenStdTimeSleep := 10
				for i := 1; i <= DownLoadOPenStdTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page = "+strconv.Itoa(page)+"===StdName="+StdName+"===========操作完成，", "暂停", DownLoadOPenStdTimeSleep, "秒，倒计时", i, "秒===========")
				}
			}
			page++
		}
	}
}

func downloadOPenStd(attachmentUrl string, referer string, filePath string) error {
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
	if OPenStdEnableHttpProxy {
		client = OPenStdSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
// 	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
// 	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	// 获取当前时间的纳秒级时间戳
    nowTimestamp := time.Now().Unix()
	OPenStdCookie := fmt.Sprintf("JSESSIONID=9C35F92116F63703FEF2CC9F0C54D32D; Hm_lvt_54db9897e5a65f7a7b00359d86015d8d=1777604946; HMACCOUNT=4E5B3419A3141A8E; Hm_lvt_50758913e6f0dfc9deacbfebce3637e4=1781876650; Hm_lpvt_50758913e6f0dfc9deacbfebce3637e4=%d", nowTimestamp)
	fmt.Println(OPenStdCookie)
	req.Header.Set("Cookie", OPenStdCookie)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "openstd.samr.gov.cn")
// 	req.Header.Set("Origin", "https://openstd.samr.gov.cn/")
	req.Header.Set("Referer", referer)
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"148\", \"Google Chrome\";v=\"148\", \"Not/A)Brand\";v=\"99\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println(resp.Header.Get("Content-Disposition"))
	fmt.Println(resp.Header.Get("Content-Length"))
	fmt.Println(resp.Header.Get("Content-Type"))
// 	os.Exit(1)
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
