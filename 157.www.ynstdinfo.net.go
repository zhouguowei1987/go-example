package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	_ "os"
	"path/filepath"
	_ "path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	_ "golang.org/x/net/html"
)

var YnStdInfoEnableHttpProxy = false
var YnStdInfoHttpProxyUrl = "111.225.152.186:8089"
var YnStdInfoHttpProxyUrlArr = make([]string, 0)

func YnStdInfoHttpProxy() error {
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
					YnStdInfoHttpProxyUrlArr = append(YnStdInfoHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					YnStdInfoHttpProxyUrlArr = append(YnStdInfoHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func YnStdInfoSetHttpProxy() (httpclient *http.Client) {
	if YnStdInfoHttpProxyUrl == "" {
		if len(YnStdInfoHttpProxyUrlArr) <= 0 {
			err := YnStdInfoHttpProxy()
			if err != nil {
				YnStdInfoSetHttpProxy()
			}
		}
		YnStdInfoHttpProxyUrl = YnStdInfoHttpProxyUrlArr[0]
		if len(YnStdInfoHttpProxyUrlArr) >= 2 {
			YnStdInfoHttpProxyUrlArr = YnStdInfoHttpProxyUrlArr[1:]
		} else {
			YnStdInfoHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(YnStdInfoHttpProxyUrl)
	ProxyURL, _ := url.Parse(YnStdInfoHttpProxyUrl)
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

type QueryYnStdInfoListFormData struct {
	number          string
	stdName         string
	stdOrg          string
	ics             string
	ccs             string
	status          string
	statusA         bool
	statusACheckbox bool
	statusW         bool
	statusN         bool
	statusNCheckbox bool
	page            int
}

var YnStdInfoCookie = "COOKIE_SUPPORT=true; GUEST_LANGUAGE_ID=zh_CN; UM_distinctid=198ac2e12d011bb-05be4ef07d0d1b-26001d51-1fa400-198ac2e12d1165b; JSESSIONID=1828AB0C971FFA7EF08390C0EDCCDFE0; CNZZDATA1264399389=324404656-1755235357-https%253A%252F%252Fwww.baidu.com%252F%7C1755248893; LFR_SESSION_STATE_20158=1755248892839"

// 下载云南省地方标准文档
// @Title 下载云南省地方标准文档
// @Description https://www.ynstdinfo.net/，下载云南省地方标准文档
func main() {
	pageListUrl := "https://www.ynstdinfo.net/web/guest/db53?p_auth=wWHFjHNS&p_p_id=simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD&p_p_lifecycle=1&p_p_state=normal&p_p_mode=view&p_p_col_id=column-1&p_p_col_count=2&_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_javax.portlet.action=search&_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_mvcPath=%2Fhtml%2Fsimplesearch%2Fview.jsp&_disable_auto_search=true"
	fmt.Println(pageListUrl)
	page := 1
	maxPage := 152
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}
		queryYnStdInfoListFormData := QueryYnStdInfoListFormData{
			number:          "",
			stdName:         "",
			stdOrg:          "",
			ics:             "",
			ccs:             "",
			status:          "A W N",
			statusA:         true,
			statusACheckbox: true,
			statusW:         false,
			statusN:         true,
			statusNCheckbox: true,
			page:            page,
		}
		queryYnStdInfoListDoc, err := QueryYnStdInfoList(pageListUrl, queryYnStdInfoListFormData)
		if err != nil {
			fmt.Println(err)
			break
		}

		liNodes := htmlquery.Find(queryYnStdInfoListDoc, `//html/body/div/div/div/div/div/div/div[1]/section/div/div/div/form/div[2]/ul/li`)
		if len(liNodes) >= 1 {
			for _, liNode := range liNodes {
				fmt.Println("=====================开始处理数据 page = ", page, "=========================")

				codeNode := htmlquery.FindOne(liNode, `./table/tbody/tr[1]/td[1]/a[1]`)
				if codeNode == nil {
					fmt.Println("标准号不存在，跳过")
					continue
				}
				code := strings.TrimSpace(htmlquery.InnerText(codeNode))
				code = strings.ReplaceAll(code, "/", "-")
				code = strings.ReplaceAll(code, "—", "-")
				fmt.Println(code)

				titleNode := htmlquery.FindOne(liNode, `./table/tbody/tr[2]/td`)
				if titleNode == nil {
					fmt.Println("标题不存在，跳过")
					continue
				}
				title := strings.TrimSpace(htmlquery.InnerText(titleNode))
				title = strings.ReplaceAll(title, " ", "-")
				title = strings.ReplaceAll(title, "　", "-")
				title = strings.ReplaceAll(title, "/", "-")
				title = strings.ReplaceAll(title, "--", "-")
				fmt.Println(title)

				filePath := "../www.ynstdinfo.net/" + title + "(" + code + ").pdf"
				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				downloadUrlNode := htmlquery.FindOne(liNode, `./table/tbody/tr[1]/td[1]/a[2]/@href`)
				if downloadUrlNode == nil {
					fmt.Println("详情链接地址不存在，跳过")
					continue
				}
				downloadUrl := strings.TrimSpace(htmlquery.InnerText(downloadUrlNode))
				fmt.Println(downloadUrl)

				fmt.Println("=======开始下载" + title + "========")

				err = downloadYnStdInfo(downloadUrl, pageListUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				os.Exit(1)
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "../www.ynstdinfo.net", "../upload.doc88.com/dbba.sacinfo.org.cn")
				err = copyYnStdInfoFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")
				//DownLoadYnStdInfoTimeSleep := 10
				DownLoadYnStdInfoTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadYnStdInfoTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadYnStdInfoTimeSleep, "秒 倒计时", i, "秒===========")
				}
			}
		}
		DownLoadYnStdInfoPageTimeSleep := 10
		// DownLoadYnStdInfoPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadYnStdInfoPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadYnStdInfoPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

func QueryYnStdInfoList(requestUrl string, queryYnStdInfoListFormData QueryYnStdInfoListFormData) (doc *html.Node, err error) {
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
	if YnStdInfoEnableHttpProxy {
		client = YnStdInfoSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_number", queryYnStdInfoListFormData.number)
	postData.Add("_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_stdName", queryYnStdInfoListFormData.stdName)
	postData.Add("_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_stdOrg", queryYnStdInfoListFormData.stdOrg)
	postData.Add("_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_ics", queryYnStdInfoListFormData.ics)
	postData.Add("_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_ccs", queryYnStdInfoListFormData.ccs)
	postData.Add("_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_status", queryYnStdInfoListFormData.status)
	postData.Add("_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_statusA", strconv.FormatBool(queryYnStdInfoListFormData.statusA))
	postData.Add("_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_statusACheckbox", strconv.FormatBool(queryYnStdInfoListFormData.statusACheckbox))
	postData.Add("_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_statusW", strconv.FormatBool(queryYnStdInfoListFormData.statusW))
	postData.Add("_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_statusN", strconv.FormatBool(queryYnStdInfoListFormData.statusN))
	postData.Add("_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_statusNCheckbox", strconv.FormatBool(queryYnStdInfoListFormData.statusNCheckbox))
	postData.Add("_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_page", strconv.Itoa(queryYnStdInfoListFormData.page))
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("authority", "www.ynstdinfo.net")
	req.Header.Set("method", "POST")
	req.Header.Set("path", "/web/guest/db53?p_auth=9NnvQ604&p_p_id=simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD&p_p_lifecycle=1&p_p_state=normal&p_p_mode=view&p_p_col_id=column-1&p_p_col_count=2&_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_javax.portlet.action=search&_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_mvcPath=%2Fhtml%2Fsimplesearch%2Fview.jsp&_disable_auto_search=true")
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", YnStdInfoCookie)
	req.Header.Set("Host", "www.ynstdinfo.net")
	req.Header.Set("Origin", "https://www.ynstdinfo.net")
	//req.Header.Set("Priority", "\nu=0, i")
	req.Header.Set("Referer", "https://www.ynstdinfo.net/web/guest/db53?p_auth=9NnvQ604&p_p_id=simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD&p_p_lifecycle=1&p_p_state=normal&p_p_mode=view&p_p_col_id=column-1&p_p_col_count=2&_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_javax.portlet.action=search&_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_mvcPath=%2Fhtml%2Fsimplesearch%2Fview.jsp&_disable_auto_search=true")
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
		return doc, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return doc, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	doc, err = htmlquery.Parse(resp.Body)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func downloadYnStdInfo(attachmentUrl string, referer string, filePath string) error {
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
	if YnStdInfoEnableHttpProxy {
		client = YnStdInfoSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("authority", "www.ynstdinfo.net")
	req.Header.Set("method", "POST")
	req.Header.Set("path", "/web/guest/db53?p_auth=9NnvQ604&p_p_id=simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD&p_p_lifecycle=1&p_p_state=normal&p_p_mode=view&p_p_col_id=column-1&p_p_col_count=2&_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_javax.portlet.action=search&_simplesearch_WAR_bhystdportlet_INSTANCE_OQOoH9AprmUD_mvcPath=%2Fhtml%2Fsimplesearch%2Fview.jsp&_disable_auto_search=true")
	req.Header.Set("scheme", "https")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.ynstdinfo.net")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	os.Exit(1)

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

func copyYnStdInfoFile(src, dst string) (err error) {
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
