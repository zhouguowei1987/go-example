package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var EditDoc88EnableHttpProxy = false
var EditDoc88HttpProxyUrl = "111.225.152.186:8089"
var EditDoc88HttpProxyUrlArr = make([]string, 0)

func EditDoc88HttpProxy() error {
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
					EditDoc88HttpProxyUrlArr = append(EditDoc88HttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					EditDoc88HttpProxyUrlArr = append(EditDoc88HttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func EditDoc88SetHttpProxy() (httpclient *http.Client) {
	if EditDoc88HttpProxyUrl == "" {
		if len(EditDoc88HttpProxyUrlArr) <= 0 {
			err := EditDoc88HttpProxy()
			if err != nil {
				EditDoc88SetHttpProxy()
			}
		}
		EditDoc88HttpProxyUrl = EditDoc88HttpProxyUrlArr[0]
		if len(EditDoc88HttpProxyUrlArr) >= 2 {
			EditDoc88HttpProxyUrlArr = EditDoc88HttpProxyUrlArr[1:]
		} else {
			EditDoc88HttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(EditDoc88HttpProxyUrl)
	ProxyURL, _ := url.Parse(EditDoc88HttpProxyUrl)
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

type QueryEditDoc88ListFormData struct {
	MenuIndex  int
	ClassifyId string
	FolderId   int
	Sort       int
	Keyword    string
	ShowIndex  int
}

type EditDoc88ResponseData struct {
	Result     string `json:"result"`
	EditTitle  string `json:"edit_title"`
	Class      string `json:"class"`
	UpdateInfo string `json:"updateinfo"`
	State      string `json:"state"`
	SaveFile   string `json:"savefile"`
	Other      string `json:"other"`
}

type EditDoc88FormData struct {
	DocCode        string
	Title          string
	Intro          string
	PCid           string
	Keyword        string
	ShareToDoc     string
	Download       string
	PPrice         string
	PDefaultPoints string
	PPageCount     string
	PDocFormat     string
	Act            string
	GroupList      string
	GroupFreeList  string
}

var EditListCookie = "new_user_task_degree=100; cdb_sys_sid=usrcbi0lk9bcuoaok1d5h70ke0; cdb_RW_ID_1392951280=2340; cdb_RW_ID_1672846620=3839; cdb_RW_ID_2008636090=1; cdb_RW_ID_2008629885=1; cdb_RW_ID_335703776=5; cdb_RW_ID_228722088=7; cdb_RW_ID_2008798462=1587; cdb_RW_ID_1675317234=950; cdb_RW_ID_2008457230=686; cdb_RW_ID_2005292790=19; cdb_RW_ID_2036761478=1; cdb_RW_ID_2035886221=1; cdb_RW_ID_2037028201=2; cdb_RW_ID_1443403331=620; cdb_RW_ID_2043350748=2; cdb_RW_ID_2050714108=1; cdb_RW_ID_2038752397=334; ShowSkinTip_1=1; showAnnotateTipIf=1; Page_13743815423194=1; doc88_lt=doc88; cdb_login_if=1; cdb_uid=104598337; cdb_tokenid=d9c3JT8Nn4biuHjlgRiRhsxY7%2BiGJVwjubocIDlhJA9UZ2c7CgVEhg%2Bz1qAOKRoFf1u9DkC7%2Fy9t8ou31UQ4%2BPaYW%2F9thhQbhPp%2FOozoNUBb; c_login_name=woyoceo; cdb_logined=1; cdb_RW_ID_1453569777=9; Page_84861507023666=1; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d23cc6089fd85ee433394d2a19ca101ca39d92f350c3f26f274; cdb_change_message=1; cdb_msg_num=0; cdb_RW_ID_2064386957=1; Page_40320547234861=1; cdb_RW_ID_2058901312=1; Page_Y_49220563859290=-93.23711340206185; Page_49220563859290=1; cdb_RW_ID_2061952146=1; Page_Y_50780457630785=-93.23711340206185; Page_50780457630785=1; cdb_RW_ID_2063767762=1; cdb_RW_ID_783946052=59; Page_Y_21271820727721=687.5515463917526; Page_21271820727721=4; Page_Y_1436697352408=-80.54432989690721; Page_1436697352408=1; cdb_RW_ID_2063766612=1; cdb_READED_PC_ID=%2C440447; Page_Y_27237820422297=-401.7783505154639; Page_27237820422297=3; cdb_RW_ID_84150571=51; cdb_VIEW_DOC_ID=%2C783946052%2C84150571; Page_776450728247=1; cdb_RW_ID_2063765341=1; Page_Y_84854085381579=283.3917525773196; Page_84854085381579=4; cdb_RW_ID_2063764256=1; Page_Y_27237820426732=222.05154639175257; Page_27237820426732=6; PHPSESSID=usrcbi0lk9bcuoaok1d5h70ke0; siftState=1; show_index=1; cdb_msg_time=1751115070"
var EditDetailCookie = "cdb_sys_sid=usrcbi0lk9bcuoaok1d5h70ke0; cdb_RW_ID_1392951280=2340; cdb_RW_ID_1672846620=3839; cdb_RW_ID_2008636090=1; cdb_RW_ID_2008629885=1; cdb_RW_ID_335703776=5; cdb_RW_ID_228722088=7; cdb_RW_ID_2008798462=1587; cdb_RW_ID_1675317234=950; cdb_RW_ID_2008457230=686; cdb_RW_ID_2005292790=19; cdb_RW_ID_2036761478=1; cdb_RW_ID_2035886221=1; cdb_RW_ID_2037028201=2; cdb_RW_ID_1443403331=620; cdb_RW_ID_2043350748=2; cdb_RW_ID_2050714108=1; cdb_RW_ID_2038752397=334; ShowSkinTip_1=1; showAnnotateTipIf=1; Page_13743815423194=1; doc88_lt=doc88; cdb_login_if=1; cdb_uid=104598337; cdb_tokenid=d9c3JT8Nn4biuHjlgRiRhsxY7%2BiGJVwjubocIDlhJA9UZ2c7CgVEhg%2Bz1qAOKRoFf1u9DkC7%2Fy9t8ou31UQ4%2BPaYW%2F9thhQbhPp%2FOozoNUBb; c_login_name=woyoceo; cdb_logined=1; cdb_RW_ID_1453569777=9; Page_84861507023666=1; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d23cc6089fd85ee433394d2a19ca101ca39d92f350c3f26f274; cdb_change_message=1; cdb_msg_num=0; cdb_RW_ID_2064386957=1; Page_40320547234861=1; cdb_RW_ID_2058901312=1; Page_Y_49220563859290=-93.23711340206185; Page_49220563859290=1; cdb_RW_ID_2061952146=1; Page_Y_50780457630785=-93.23711340206185; Page_50780457630785=1; cdb_RW_ID_2063767762=1; cdb_RW_ID_783946052=59; Page_Y_21271820727721=687.5515463917526; Page_21271820727721=4; Page_Y_1436697352408=-80.54432989690721; Page_1436697352408=1; cdb_RW_ID_2063766612=1; cdb_READED_PC_ID=%2C440447; Page_Y_27237820422297=-401.7783505154639; Page_27237820422297=3; cdb_RW_ID_84150571=51; cdb_VIEW_DOC_ID=%2C783946052%2C84150571; Page_776450728247=1; cdb_RW_ID_2063765341=1; Page_Y_84854085381579=283.3917525773196; Page_84854085381579=4; cdb_RW_ID_2063764256=1; Page_Y_27237820426732=222.05154639175257; Page_27237820426732=6; PHPSESSID=usrcbi0lk9bcuoaok1d5h70ke0; siftState=1; show_index=1; cdb_msg_time=1751115070"
var EditEditCookie = "cdb_sys_sid=usrcbi0lk9bcuoaok1d5h70ke0; cdb_RW_ID_1392951280=2340; cdb_RW_ID_1672846620=3839; cdb_RW_ID_2008636090=1; cdb_RW_ID_2008629885=1; cdb_RW_ID_335703776=5; cdb_RW_ID_228722088=7; cdb_RW_ID_2008798462=1587; cdb_RW_ID_1675317234=950; cdb_RW_ID_2008457230=686; cdb_RW_ID_2005292790=19; cdb_RW_ID_2036761478=1; cdb_RW_ID_2035886221=1; cdb_RW_ID_2037028201=2; cdb_RW_ID_1443403331=620; cdb_RW_ID_2043350748=2; cdb_RW_ID_2050714108=1; cdb_RW_ID_2038752397=334; ShowSkinTip_1=1; showAnnotateTipIf=1; Page_13743815423194=1; doc88_lt=doc88; cdb_login_if=1; cdb_uid=104598337; cdb_tokenid=d9c3JT8Nn4biuHjlgRiRhsxY7%2BiGJVwjubocIDlhJA9UZ2c7CgVEhg%2Bz1qAOKRoFf1u9DkC7%2Fy9t8ou31UQ4%2BPaYW%2F9thhQbhPp%2FOozoNUBb; c_login_name=woyoceo; cdb_logined=1; cdb_RW_ID_1453569777=9; Page_84861507023666=1; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d23cc6089fd85ee433394d2a19ca101ca39d92f350c3f26f274; cdb_change_message=1; cdb_msg_num=0; cdb_RW_ID_2064386957=1; Page_40320547234861=1; cdb_RW_ID_2058901312=1; Page_Y_49220563859290=-93.23711340206185; Page_49220563859290=1; cdb_RW_ID_2061952146=1; Page_Y_50780457630785=-93.23711340206185; Page_50780457630785=1; cdb_RW_ID_2063767762=1; cdb_RW_ID_783946052=59; Page_Y_21271820727721=687.5515463917526; Page_21271820727721=4; Page_Y_1436697352408=-80.54432989690721; Page_1436697352408=1; cdb_RW_ID_2063766612=1; cdb_READED_PC_ID=%2C440447; Page_Y_27237820422297=-401.7783505154639; Page_27237820422297=3; cdb_RW_ID_84150571=51; cdb_VIEW_DOC_ID=%2C783946052%2C84150571; Page_776450728247=1; cdb_RW_ID_2063765341=1; Page_Y_84854085381579=283.3917525773196; Page_84854085381579=4; cdb_RW_ID_2063764256=1; Page_Y_27237820426732=222.05154639175257; Page_27237820426732=6; PHPSESSID=usrcbi0lk9bcuoaok1d5h70ke0; siftState=1; show_index=1; cdb_msg_time=1751115070"

var EditDetailTimeSleep = 10
var EditSaveTimeSleep = 10
var EditNextPageSleep = 15

// ychEduSpider 编辑道客巴巴文档
// @Title 编辑道客巴巴文档
// @Description https://www.doc88.com/，编辑道客巴巴文档
func main() {
	curPage := 2586

	for {
		pageListUrl := fmt.Sprintf("https://www.doc88.com/uc/doc_manager.php?act=ajax_doc_list&curpage=%d", curPage)
		fmt.Println(pageListUrl)
		queryEditDoc88ListFormData := QueryEditDoc88ListFormData{
			MenuIndex:  4,
			ClassifyId: "all",
			FolderId:   0,
			Sort:       2,
			Keyword:    "",
			ShowIndex:  1,
		}
		pageListDoc, err := QueryEditDoc88List(pageListUrl, queryEditDoc88ListFormData)
		if err != nil {
			EditDoc88HttpProxyUrl = ""
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

			IntroNode := htmlquery.FindOne(liNode, `./div[@class="bookdoc"]/p`)
			Intro := htmlquery.InnerText(IntroNode)

			PPageCountNode := htmlquery.FindOne(liNode, `./div[@class="bookimg"]/em`)
			PPageCount := htmlquery.InnerText(PPageCountNode)
			PPageCount = PPageCount[2:]

			PPriceNode := htmlquery.FindOne(liNode, `./div[@class="bookdoc"]/ul[@class="position"]/li[6]/span[@class="jifentip"]/strong[@class="red"]`)
			PPrice := htmlquery.InnerText(PPriceNode)

			filePageNum, _ := strconv.Atoi(PPageCount)
			PPriceNew := ""
			// 根据页数设置价格
			if filePageNum > 0 && filePageNum <= 5 {
				PPriceNew = "288"
			} else if filePageNum > 5 && filePageNum <= 10 {
				PPriceNew = "388"
			} else if filePageNum > 10 && filePageNum <= 15 {
				PPriceNew = "488"
			} else if filePageNum > 15 && filePageNum <= 20 {
				PPriceNew = "588"
			} else if filePageNum > 20 && filePageNum <= 25 {
				PPriceNew = "688"
			} else if filePageNum > 25 && filePageNum <= 30 {
				PPriceNew = "788"
			} else if filePageNum > 30 && filePageNum <= 35 {
				PPriceNew = "888"
			} else {
				PPriceNew = "988"
			}

			CateLogNameNode := htmlquery.FindOne(liNode, `./div[@class="bookdoc"]/div[@class="posttime"]/span/span[@class="catelog_name"]`)
			if CateLogNameNode == nil {
				continue
			}
			CateLogName := htmlquery.InnerText(CateLogNameNode)
			if CateLogName != "实用文书>标准规范>国家标准" && CateLogName != "实用文书>标准规范>行业标准" && CateLogName != "实用文书>标准规范>地方标准" && CateLogName != "实用文书>标准规范>企业标准" {
				PPriceNew = "200"
			}

			// 新旧价格一样，则跳过
			fmt.Println(PPrice, PPriceNew)
			if PPrice == PPriceNew {
				continue
			}

			PId := htmlquery.SelectAttr(liNode, "id")
			PId = PId[5:]

			for i := 1; i <= EditDetailTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(curPage)+"===========获取", Title, "详情暂停", EditDetailTimeSleep, "秒，倒计时", i, "秒===========")
			}

			detailUrl := "https://www.doc88.com/uc/usr_doc_manager.php?act=getDocInfo"
			detailDoc, err := QueryEditDoc88Detail(detailUrl, PId)
			if err != nil {
				EditDoc88HttpProxyUrl = ""
				fmt.Println(err)
				continue
			}

			DocCodeNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[1]/div[@class="booksedit"]/table[@class="edit-table"]/input`)
			DocCode := htmlquery.SelectAttr(DocCodeNode, "value")

			PCidNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[1]/div[@class="booksedit"]/table[@class="edit-table"]/tbody/tr[3]/td[2]/div[@class="layers"]/input`)
			PCid := htmlquery.SelectAttr(PCidNode, "value")

			PDocFormatNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[2]/div[@class="booksedit booksedit-bdr"]/table[@class="edit-table"]/tbody/tr[3]/td[2]/input[3]`)
			PDocFormat := htmlquery.SelectAttr(PDocFormatNode, "value")

			fmt.Println("===========开始修改", Title, "价格===========")
			editUrl := "https://www.doc88.com/uc/index.php"
			editDoc88FormData := EditDoc88FormData{
				DocCode:        DocCode,
				Title:          Title,
				Intro:          Intro,
				PCid:           PCid,
				Keyword:        "",
				ShareToDoc:     "1",
				Download:       "2",
				PPrice:         PPriceNew,
				PDefaultPoints: "3",
				PPageCount:     PPageCount,
				PDocFormat:     PDocFormat,
				Act:            "save_info",
				GroupList:      "",
				GroupFreeList:  "",
			}

			for i := 1; i <= EditSaveTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(curPage)+"===========更新", Title, "成功，暂停", EditSaveTimeSleep, "秒，倒计时", i, "秒===========")
			}

			_, err = EditDoc88(editUrl, editDoc88FormData)
			if err != nil {
				EditDoc88HttpProxyUrl = ""
				fmt.Println(err)
				continue
			}
		}
		curPage++
		for i := 1; i <= EditNextPageSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("===========翻", curPage, "页，暂停", EditNextPageSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

func QueryEditDoc88List(requestUrl string, queryEditDoc88ListFormData QueryEditDoc88ListFormData) (doc *html.Node, err error) {
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
	if EditDoc88EnableHttpProxy {
		client = EditDoc88SetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("menuIndex", strconv.Itoa(queryEditDoc88ListFormData.MenuIndex))
	postData.Add("classify_id", queryEditDoc88ListFormData.ClassifyId)
	postData.Add("folder_id", strconv.Itoa(queryEditDoc88ListFormData.FolderId))
	postData.Add("sort", strconv.Itoa(queryEditDoc88ListFormData.Sort))
	postData.Add("keyword", queryEditDoc88ListFormData.Keyword)
	postData.Add("show_index", strconv.Itoa(queryEditDoc88ListFormData.ShowIndex))
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", EditListCookie)
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

func QueryEditDoc88Detail(requestUrl string, PId string) (doc *html.Node, err error) {
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
	if EditDoc88EnableHttpProxy {
		client = EditDoc88SetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("p_id", PId)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", EditDetailCookie)
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

func EditDoc88(requestUrl string, editDoc88FormData EditDoc88FormData) (editDoc88ResponseData EditDoc88ResponseData, err error) {
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
	if EditDoc88EnableHttpProxy {
		client = EditDoc88SetHttpProxy()
	}
	editDoc88ResponseData = EditDoc88ResponseData{}
	postData := url.Values{}
	postData.Add("doccode", editDoc88FormData.DocCode)
	postData.Add("title", editDoc88FormData.Title)
	postData.Add("intro", editDoc88FormData.Intro)
	postData.Add("pcid", editDoc88FormData.PCid)
	postData.Add("keyword", editDoc88FormData.Keyword)
	postData.Add("sharetodoc", editDoc88FormData.ShareToDoc)
	postData.Add("download", editDoc88FormData.Download)
	postData.Add("p_price", editDoc88FormData.PPrice)
	postData.Add("p_default_points", editDoc88FormData.PDefaultPoints)
	postData.Add("p_pagecount", editDoc88FormData.PPageCount)
	postData.Add("p_doc_format", editDoc88FormData.PDocFormat)
	postData.Add("act", "save_info")
	postData.Add("group_list", editDoc88FormData.GroupList)
	postData.Add("group_free_list", editDoc88FormData.GroupFreeList)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return editDoc88ResponseData, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", EditEditCookie)
	req.Header.Set("Host", "www.doc88.com")
	req.Header.Set("Origin", "https://www.doc88.com")
	req.Header.Set("Referer", "https://www.doc88.com/uc/doc_manager.php?act=doc_list&state=all")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return editDoc88ResponseData, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return editDoc88ResponseData, err
	}
	err = json.Unmarshal(respBytes, &editDoc88ResponseData)
	if err != nil {
		return editDoc88ResponseData, err
	}
	return editDoc88ResponseData, nil
}
