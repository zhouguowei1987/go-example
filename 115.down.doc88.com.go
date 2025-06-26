package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var DownDoc88EnableHttpProxy = false
var DownDoc88HttpProxyUrl = "111.225.152.186:8089"
var DownDoc88HttpProxyUrlArr = make([]string, 0)

func DownDoc88HttpProxy() error {
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
					DownDoc88HttpProxyUrlArr = append(DownDoc88HttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					DownDoc88HttpProxyUrlArr = append(DownDoc88HttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func DownDoc88SetHttpProxy() (httpclient *http.Client) {
	if DownDoc88HttpProxyUrl == "" {
		if len(DownDoc88HttpProxyUrlArr) <= 0 {
			err := DownDoc88HttpProxy()
			if err != nil {
				DownDoc88SetHttpProxy()
			}
		}
		DownDoc88HttpProxyUrl = DownDoc88HttpProxyUrlArr[0]
		if len(DownDoc88HttpProxyUrlArr) >= 2 {
			DownDoc88HttpProxyUrlArr = DownDoc88HttpProxyUrlArr[1:]
		} else {
			DownDoc88HttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(DownDoc88HttpProxyUrl)
	ProxyURL, _ := url.Parse(DownDoc88HttpProxyUrl)
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

type QueryDownDoc88ListFormData struct {
	MenuIndex  int
	ClassifyId string
	FolderId   int
	Sort       int
	Keyword    string
	ShowIndex  int
}

var DownListCookie = "new_user_task_degree=100; cdb_sys_sid=usrcbi0lk9bcuoaok1d5h70ke0; cdb_RW_ID_1392951280=2340; cdb_RW_ID_1672846620=3839; cdb_RW_ID_2008636090=1; cdb_RW_ID_2008629885=1; cdb_RW_ID_335703776=5; cdb_RW_ID_228722088=7; cdb_READED_PC_ID=%2C447; cdb_RW_ID_2008798462=1587; cdb_RW_ID_1675317234=950; cdb_RW_ID_2008457230=686; cdb_RW_ID_2005292790=19; cdb_RW_ID_2036761478=1; cdb_RW_ID_2035886221=1; cdb_RW_ID_2037028201=2; cdb_RW_ID_1443403331=620; cdb_RW_ID_2043350748=2; cdb_RW_ID_2050714108=1; cdb_RW_ID_2038752397=334; ShowSkinTip_1=1; showAnnotateTipIf=1; Page_13743815423194=1; PHPSESSID=usrcbi0lk9bcuoaok1d5h70ke0; doc88_lt=doc88; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d23d9581971a85c09bd546a550a976fe5c6d92f350c3f26f274; cdb_login_if=1; cdb_uid=104598337; cdb_tokenid=d9c3JT8Nn4biuHjlgRiRhsxY7%2BiGJVwjubocIDlhJA9UZ2c7CgVEhg%2Bz1qAOKRoFf1u9DkC7%2Fy9t8ou31UQ4%2BPaYW%2F9thhQbhPp%2FOozoNUBb; c_login_name=woyoceo; cdb_logined=1; show_index=1; cdb_change_message=1; cdb_msg_num=0; siftState=0; cdb_pageType=2; cdb_msg_time=1750770483"
var DownLoadUrlCookie = "cdb_sys_sid=usrcbi0lk9bcuoaok1d5h70ke0; cdb_RW_ID_1392951280=2340; cdb_RW_ID_1672846620=3839; cdb_RW_ID_2008636090=1; cdb_RW_ID_2008629885=1; cdb_RW_ID_335703776=5; cdb_RW_ID_228722088=7; cdb_READED_PC_ID=%2C447; cdb_RW_ID_2008798462=1587; cdb_RW_ID_1675317234=950; cdb_RW_ID_2008457230=686; cdb_RW_ID_2005292790=19; cdb_RW_ID_2036761478=1; cdb_RW_ID_2035886221=1; cdb_RW_ID_2037028201=2; cdb_RW_ID_1443403331=620; cdb_RW_ID_2043350748=2; cdb_RW_ID_2050714108=1; cdb_RW_ID_2038752397=334; ShowSkinTip_1=1; showAnnotateTipIf=1; Page_13743815423194=1; PHPSESSID=usrcbi0lk9bcuoaok1d5h70ke0; doc88_lt=doc88; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d23d9581971a85c09bd546a550a976fe5c6d92f350c3f26f274; cdb_login_if=1; cdb_uid=104598337; cdb_tokenid=d9c3JT8Nn4biuHjlgRiRhsxY7%2BiGJVwjubocIDlhJA9UZ2c7CgVEhg%2Bz1qAOKRoFf1u9DkC7%2Fy9t8ou31UQ4%2BPaYW%2F9thhQbhPp%2FOozoNUBb; c_login_name=woyoceo; cdb_logined=1; show_index=1; cdb_change_message=1; cdb_msg_num=0; siftState=0; cdb_pageType=2; cdb_msg_time=1750770483"

var DownNextPageSleep = 10
var TodayCurrentDownLoadCount = 0
var TodayMaxDownLoadCount = 150000

// ychEduSpider 下载道客巴巴文档
// @Title 下载道客巴巴文档
// @Description https://www.doc88.com/，下载道客巴巴文档
func main() {
	curPage := 29
	isPageListGo := true
	for isPageListGo {
		pageListUrl := fmt.Sprintf("https://www.doc88.com/uc/doc_manager.php?act=ajax_doc_list&curpage=%d", curPage)
		fmt.Println(pageListUrl)
		queryDownDoc88ListFormData := QueryDownDoc88ListFormData{
			MenuIndex:  4,
			ClassifyId: "all",
			FolderId:   0,
			Sort:       2,
			Keyword:    "",
			ShowIndex:  1,
		}
		pageListDoc, err := QueryDownDoc88List(pageListUrl, queryDownDoc88ListFormData)
		if err != nil {
			DownDoc88HttpProxyUrl = ""
			fmt.Println(err)
			continue
		}
		liNodes := htmlquery.Find(pageListDoc, `//div[@id="detailed"]/ul[@class="bookshow3"]/li`)
		if len(liNodes) <= 0 {
			break
		}
		for _, liNode := range liNodes {

			TitleNode := htmlquery.FindOne(liNode, `./div[@class="bookdoc"]/h3/a`)
			Title := htmlquery.SelectAttr(TitleNode, "title")
			fmt.Println(Title)

			CateLogNameNode := htmlquery.FindOne(liNode, `./div[@class="bookdoc"]/div[@class="posttime"]/span/span[@class="catelog_name"]`)
			if CateLogNameNode == nil {
				continue
			}
			CateLogName := htmlquery.InnerText(CateLogNameNode)
			if CateLogName != "实用文书>标准规范>行业标准" {
				continue
			}

			fileDownLoadCountNode := htmlquery.FindOne(liNode, `./div[@class="bookdoc"]/ul[@class="position"]/li[8]/span[@class="red"]`)
			if fileDownLoadCountNode == nil {
				continue
			}
			fileDownLoadCountStr := htmlquery.InnerText(fileDownLoadCountNode)
			fileDownLoadCount, _ := strconv.Atoi(fileDownLoadCountStr)
			if fileDownLoadCount <= 0 {
				continue
			}
			fmt.Println("=====下载次数", fileDownLoadCount, "=====")

			filePath := "../down.ttbz.org.cn/" + Title + ".pdf"
			if _, err := os.Stat(filePath); err != nil {

				PidStr := htmlquery.SelectAttr(liNode, "id")
				PidStr = strings.ReplaceAll(PidStr, "bdoc_", "")
				Pid, _ := strconv.Atoi(PidStr)

				requestQueryDownDoc88DownLoadUrl := fmt.Sprintf("https://www.doc88.com/doc.php?act=download&pid=%d", Pid)
				QueryDownDoc88DownLoadUrlDoc, err := QueryDownDoc88DownLoadUrl(requestQueryDownDoc88DownLoadUrl)
				if err != nil {
					DownDoc88HttpProxyUrl = ""
					fmt.Println(err)
					continue
				}

				refererUrl := "https://www.doc88.com/uc/doc_manager.php?act=doc_list&state=myshare"

				downloadUrl := htmlquery.InnerText(QueryDownDoc88DownLoadUrlDoc)
				fmt.Println(downloadUrl)

				fmt.Println("=======开始下载" + Title + "========")
				err = DownLoadDoc88(downloadUrl, refererUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")

				// DownLoadDoc88TimeSleep := rand.Intn(10)
				DownLoadDoc88TimeSleep := 10
				for i := 1; i <= DownLoadDoc88TimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("curPage="+strconv.Itoa(curPage)+"===========下载", Title, "成功，暂停", DownLoadDoc88TimeSleep, "秒，倒计时", i, "秒===========")
				}

				TodayCurrentDownLoadCount++
				if TodayCurrentDownLoadCount >= TodayMaxDownLoadCount {
					fmt.Println("=======今天已下载", TodayCurrentDownLoadCount, "个文档，停止下载========")
					isPageListGo = false
					break
				}
			}
		}
		curPage++
		for i := 1; i <= DownNextPageSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("===========翻", curPage, "页，暂停", DownNextPageSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

func QueryDownDoc88List(requestUrl string, queryDownDoc88ListFormData QueryDownDoc88ListFormData) (doc *html.Node, err error) {
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
	if DownDoc88EnableHttpProxy {
		client = DownDoc88SetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("menuIndex", strconv.Itoa(queryDownDoc88ListFormData.MenuIndex))
	postData.Add("classify_id", queryDownDoc88ListFormData.ClassifyId)
	postData.Add("folder_id", strconv.Itoa(queryDownDoc88ListFormData.FolderId))
	postData.Add("sort", strconv.Itoa(queryDownDoc88ListFormData.Sort))
	postData.Add("keyword", queryDownDoc88ListFormData.Keyword)
	postData.Add("show_index", strconv.Itoa(queryDownDoc88ListFormData.ShowIndex))
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", DownListCookie)
	req.Header.Set("Host", "www.doc88.com")
	req.Header.Set("Origin", "https://www.doc88.com")
	req.Header.Set("Referer", "https://www.doc88.com/uc/doc_manager.php?act=doc_list&state=myshare")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return doc, err
	}
	defer resp.Body.Close()
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

func QueryDownDoc88DownLoadUrl(requestUrl string) (doc *html.Node, err error) {
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
	if DownDoc88EnableHttpProxy {
		client = DownDoc88SetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", DownLoadUrlCookie)
	req.Header.Set("Host", "www.doc88.com")
	req.Header.Set("Origin", "https://www.doc88.com")
	req.Header.Set("Referer", "https://www.doc88.com/uc/doc_manager.php?act=doc_list&state=all")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return doc, err
	}
	defer resp.Body.Close()
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

func DownLoadDoc88(attachmentUrl string, referer string, filePath string) error {
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
	if DownDoc88EnableHttpProxy {
		client = DownDoc88SetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.doc88.com")
	req.Header.Set("Origin", "https://www.doc88.com")
	req.Header.Set("Referer", referer)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
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
