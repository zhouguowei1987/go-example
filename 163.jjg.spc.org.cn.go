package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	_ "os"
	"path/filepath"
	_ "path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	_ "golang.org/x/net/html"
)

var JjgEnableHttpProxy = false
var JjgHttpProxyUrl = "111.225.152.186:8089"
var JjgHttpProxyUrlArr = make([]string, 0)

func JjgHttpProxy() error {
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
					JjgHttpProxyUrlArr = append(JjgHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					JjgHttpProxyUrlArr = append(JjgHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func JjgSetHttpProxy() (httpclient *http.Client) {
	if JjgHttpProxyUrl == "" {
		if len(JjgHttpProxyUrlArr) <= 0 {
			err := JjgHttpProxy()
			if err != nil {
				JjgSetHttpProxy()
			}
		}
		JjgHttpProxyUrl = JjgHttpProxyUrlArr[0]
		if len(JjgHttpProxyUrlArr) >= 2 {
			JjgHttpProxyUrlArr = JjgHttpProxyUrlArr[1:]
		} else {
			JjgHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(JjgHttpProxyUrl)
	ProxyURL, _ := url.Parse(JjgHttpProxyUrl)
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

type QueryJjgListFormData struct {
	pageindex int
	statusxtb string
	statusgc  string
	statusxp  string
	statusqt  string
	text      string
}

type ViewJjgFormData struct {
	a100       string
	ismobile   string
	standclass string
}

var JjgCookie = "JSESSIONID=4E69177646AB19631E4D6EB0158BA234; Hm_lvt_f1e3be80a525d0bb9e111ca6cbaa2457=1758592011; Hm_lvt_283fc9e84b8a5a7401b75b0f774e2120=1766475585; HMACCOUNT=4E5B3419A3141A8E; Qs_lvt_503365=1756359161%2C1758243669%2C1759975608%2C1763093101%2C1766475585; closeclick=closeclick; Qs_pv_503365=1535695938568337700%2C1146327973359442300%2C3509218405693096000%2C8388552161518928%2C3956291511354711000; Hm_lpvt_283fc9e84b8a5a7401b75b0f774e2120=1766476559"

// 下载国家计量技术规范文档
// @Title 下载国家计量技术规范文档
// @Description https://jjg.spc.org.cn/，下载国家计量技术规范文档
func main() {
	pageListUrl := "https://jjg.spc.org.cn/resmea/view/search"
	fmt.Println(pageListUrl)
	page := 0
	maxPage := 265
	isPageListGo := true
	for isPageListGo {
		if page > maxPage {
			isPageListGo = false
			break
		}

		queryJjgListFormData := QueryJjgListFormData{
			pageindex: page,
			statusxtb: "检定系统表",
			statusgc:  "检定规程",
			statusxp:  "型评大纲",
			statusqt:  "其他计量技术规范",
			text:      "",
		}
		queryJjgListDoc, err := QueryJjgList(pageListUrl, queryJjgListFormData)
		if err != nil {
			fmt.Println(err)
			break
		}
		trNodes := htmlquery.Find(queryJjgListDoc, `//html/body/div[2]/div[1]/div[3]/div[2]/div/table/tbody/tr`)
		if len(trNodes) >= 2 {
			for _, trNode := range trNodes {
				fmt.Println("=====================开始处理数据 page = ", page, "=========================")
				codeNode := htmlquery.FindOne(trNode, `./td[2]/span`)
				if codeNode == nil {
					fmt.Println("标准号不存在，跳过")
					continue
				}
				code := strings.TrimSpace(htmlquery.InnerText(codeNode))
				code = strings.ReplaceAll(code, "/", "-")
				code = strings.ReplaceAll(code, "—", "-")
				fmt.Println(code)

				titleNode := htmlquery.FindOne(trNode, `./td[4]`)
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

				filePath := "../jjg.spc.org.cn/" + title + "(" + code + ")" + ".pdf"
				fmt.Println(filePath)

				_, err = os.Stat(filePath)
				if err == nil {
					fmt.Println("文档已下载过，跳过")
					continue
				}

				fmt.Println("=======开始下载========")

				viewJjgUrl := "https://jjg.spc.org.cn/resmea/view/stdonline"
				viewJjgFormData := ViewJjgFormData{
					a100:       code,
					ismobile:   "",
					standclass: "",
				}
				viewJjgDoc, err := ViewJjg(viewJjgUrl, viewJjgFormData)
				if err != nil {
					fmt.Println(err)
					continue
				}

				regEnc := regexp.MustCompile("var enc = \"(.*?)\";")
				findEnc := regEnc.Find([]byte(htmlquery.InnerText(viewJjgDoc)))
				encStr := string(findEnc)
				enc := strings.ReplaceAll(encStr, "var enc = \"", "")
				enc = strings.ReplaceAll(enc, "\";", "")
				fmt.Println(enc)

				regRc := regexp.MustCompile("var rc = \"(.*?)\";")
				findRc := regRc.Find([]byte(htmlquery.InnerText(viewJjgDoc)))
				rcStr := string(findRc)
				rc := strings.ReplaceAll(rcStr, "var rc = \"", "")
				rc = strings.ReplaceAll(rc, "\";", "")
				fmt.Println(rc)

				downloadJjgUrl := fmt.Sprintf("https://jjg.spc.org.cn/resmea/view/onlinereading?token=%s", rc)
				fmt.Println(downloadJjgUrl)
				err = downloadJjg(downloadJjgUrl, enc, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "jjg.spc.org.cn", "temp-hbba.sacinfo.org.cn")
				err = copyJjgFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======完成下载========")
				//DownLoadJjgTimeSleep := 10
				DownLoadJjgTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadJjgTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("title="+title+"===========下载", title, "成功，暂停", DownLoadJjgTimeSleep, "秒，倒计时", i, "秒===========")
				}
			}
		}
		DownLoadJjgPageTimeSleep := 10
		// DownLoadJjgPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadJjgPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadJjgPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

func QueryJjgList(requestUrl string, queryJjgListFormData QueryJjgListFormData) (doc *html.Node, err error) {
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
	if JjgEnableHttpProxy {
		client = JjgSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("pageindex", strconv.Itoa(queryJjgListFormData.pageindex))
	postData.Add("statusxtb", queryJjgListFormData.statusxtb)
	postData.Add("statusgc", queryJjgListFormData.statusgc)
	postData.Add("statusxp", queryJjgListFormData.statusxp)
	postData.Add("statusqt", queryJjgListFormData.statusqt)
	postData.Add("text", queryJjgListFormData.text)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", JjgCookie)
	req.Header.Set("Host", "jjg.spc.org.cn")
	req.Header.Set("Origin", "https://jjg.spc.org.cn")
	req.Header.Set("Referer", "https://jjg.spc.org.cn/resmea/view/search")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
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

func ViewJjg(requestUrl string, viewJjgFormData ViewJjgFormData) (doc *html.Node, err error) {
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
	if JjgEnableHttpProxy {
		client = JjgSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("a100", viewJjgFormData.a100)
	postData.Add("ismobile", viewJjgFormData.ismobile)
	postData.Add("standclass", viewJjgFormData.standclass)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", JjgCookie)
	req.Header.Set("Host", "jjg.spc.org.cn")
	req.Header.Set("Origin", "https://jjg.spc.org.cn")
	req.Header.Set("Referer", fmt.Sprintf("ttps://jjg.spc.org.cn/resmea/standard/%s", viewJjgFormData.a100))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
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

func downloadJjg(requestUrl string, enc string, filePath string) error {
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
	if JjgEnableHttpProxy {
		client = JjgSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", JjgCookie)
	req.Header.Set("Host", "jjg.spc.org.cn")
	req.Header.Set("Myfoxit", enc)
	req.Header.Set("Origin", "https://jjg.spc.org.cn")
	req.Header.Set("Referer", "https://jjg.spc.org.cn/resmea/view/stdonline")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusCreated {
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

func copyJjgFile(src, dst string) (err error) {
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
