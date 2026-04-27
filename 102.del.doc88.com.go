package main

import (
// 	"encoding/json"
	"errors"
	"fmt"
// 	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type QueryDelDoc88ListFormData struct {
	MenuIndex  int
	ClassifyId string
	FolderId   int
	Sort       int
	Keyword    string
	ShowIndex  int
}

type DelDoc88FormData struct {
	Id        string
}

// Doc88Cookie 15238369929
var DelListCookie = "PHPSESSID=6nh9k3a70toarhknllurprqi75; cdb_sys_sid=6nh9k3a70toarhknllurprqi75; cdb_RW_ID_2105657455=1; ShowSkinTip_1=1; Page_Y_77468140206500=-154.36776315789473; Page_77468140206500=1; cdb_RW_ID_315077174=29; Page_1502296511917=1; cdb_RW_ID_2105657438=1; Page_Y_00871385257906=156.78947368421052; Page_00871385257906=3; cdb_RW_ID_2105744054=1; Page_88990913544134=1; cdb_RW_ID_2105265700=2; Page_10380743053244=1; cdb_RW_ID_2105744043=1; Page_Y_11780743288481=-144.4736842105263; Page_11780743288481=1; cdb_RW_ID_2105743986=1; Page_Y_22920956172834=12.134868421052632; Page_22920956172834=6; cdb_RW_ID_2026327809=157; Page_Y_18568482786943=-119.39144736842105; Page_18568482786943=1; cdb_RW_ID_2105918627=1; Page_Y_68737983591274=546.0690789473684; Page_68737983591274=3; cdb_RW_ID_2105920627=1; Page_17919678397492=1; cdb_RW_ID_2106009624=1; Page_Y_68037982885276=13.700657894736842; Page_68037982885276=1; cdb_RW_ID_2106132922=2; Page_Y_75620954920800=-52.84539473684211; Page_75620954920800=4; cdb_RW_ID_2106132918=1; Page_Y_41390912980697=-113.12828947368422; Page_41390912980697=1; cdb_RW_ID_2093562922=13; Page_77380461350600=1; cdb_RW_ID_2023400579=3; cdb_RW_ID_2106132873=2; Page_Y_98571382301670=-119.39144736842105; Page_98571382301670=1; Page_Y_67319790177823=-119.39144736842105; Page_67319790177823=1; cdb_RW_ID_2105923789=1; Page_Y_17919678390253=802.858552631579; Page_17919678390253=2; cdb_RW_ID_2051011410=11; Page_Y_86519786766167=-119.39144736842105; Page_86519786766167=1; cdb_RW_ID_2106859899=1; Page_Y_07143786529599=228.2138157894737; Page_07143786529599=3; cdb_RW_ID_2106859898=1; Page_Y_93071382654646=-119.39144736842105; Page_93071382654646=1; cdb_RW_ID_2106859894=1; Page_79220954368387=1; cdb_RW_ID_2105767240=1; showAnnotateTipIf=1; Page_77168140626854=1; cdb_RW_ID_2106862809=1; Page_Y_07143786563589=-144.4736842105263; Page_07143786563589=1; cdb_RW_ID_2106868358=1; Page_87180745959139=1; cdb_RW_ID_1856434724=2; Page_36616584101291=1; cdb_RW_ID_2106868326=1; Page_Y_79220954343204=212.21052631578948; Page_79220954343204=20; cdb_RW_ID_2106868289=1; Page_Y_87180745959096=-94.73684210526316; Page_87180745959096=24; cdb_RW_ID_2106868211=1; cdb_RW_ID_1469485837=1; Page_18947069052514=1; Page_Y_07143786565377=154.73684210526318; Page_07143786565377=6; cdb_RW_ID_2106868195=1; Page_Y_87180745959763=-241.26315789473685; Page_87180745959763=23; cdb_RW_ID_2106131694=1; Page_Y_98571382303249=-119.39144736842105; Page_98571382303249=1; cdb_RW_ID_2080095563=12; cdb_H5R=1; Page_Y_37219757738840=307.28618421052636; Page_37219757738840=1; cdb_RW_ID_2106079090=1; Page_Y_68637982845858=-119.39144736842105; Page_68637982845858=1; cdb_RW_ID_2056545947=160; Page_68643826202904=1; cdb_RW_ID_2107439286=1; Page_79654903756428=1; cdb_RW_ID_2099399935=1; Page_Y_69619733033308=-270.4901315789474; Page_69619733033308=13; Page_Y_77643895975374=-119.39144736842105; cdb_RW_ID_2098918217=8; Page_77643895975374=1; cdb_RW_ID_2098918393=30; Page_Y_77580469679161=-119.39144736842105; Page_77580469679161=1; cdb_RW_ID_2099720859=2; Page_70880466204936=1; cdb_RW_ID_2107904089=1; cdb_RW_ID_2107904087=1; Page_Y_70620951857531=356.60855263157896; Page_70620951857531=11; cdb_RW_ID_2107904086=1; Page_Y_58068146345492=369.91776315789474; Page_58068146345492=4; cdb_RW_ID_2107904961=1; Page_74154903607689=1; cdb_search_format=; searchType=0; cdb_RW_ID_2101884506=3; Page_Y_81868141995042=-196.11513157894737; Page_81868141995042=2; cdb_RW_ID_2065106162=2; Page_Y_15920546954940=-71.63486842105263; Page_15920546954940=2; Page_Y_91771386838646=502.2269736842106; cdb_RW_ID_2105239433=1; Page_Y_03443785878542=176.54276315789474; Page_03443785878542=1; cdb_READED_PC_ID=%2C730730730447440593443448440441440443; cdb_RW_ID_1908563572=2; Page_Y_70487649351320=-119.39144736842105; Page_70487649351320=1; cdb_login_if=1; cdb_uid=104598337; c_login_name=woyoceo; cdb_logined=1; cdb_change_message=1; cdb_msg_num=0; doc88_lt=wx; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d233dd02f2a58c6b2b5841b62253b0c1d47d92f350c3f26f274; cdb_RW_ID_2107957090=1; Page_91271387457848=1; cdb_RW_ID_2108010816=1; Page_58668149414912=1; cdb_RW_ID_1987519680=21; cdb_VIEW_DOC_ID=%2C1987519680; Page_Y_36116352863457=361.3059210526316; Page_36116352863457=2; cdb_RW_ID_2108010801=1; Page_Y_19219675767576=-119.39144736842105; Page_19219675767576=1; siftState=1; show_index=1; cdb_RW_ID_2108010898=2; Page_91771386838646=1; cdb_pageType=2; cdb_msg_time=1777300748"
var DelDelCookie = "PHPSESSID=6nh9k3a70toarhknllurprqi75; cdb_sys_sid=6nh9k3a70toarhknllurprqi75; cdb_RW_ID_2105657455=1; ShowSkinTip_1=1; Page_Y_77468140206500=-154.36776315789473; Page_77468140206500=1; cdb_RW_ID_315077174=29; Page_1502296511917=1; cdb_RW_ID_2105657438=1; Page_Y_00871385257906=156.78947368421052; Page_00871385257906=3; cdb_RW_ID_2105744054=1; Page_88990913544134=1; cdb_RW_ID_2105265700=2; Page_10380743053244=1; cdb_RW_ID_2105744043=1; Page_Y_11780743288481=-144.4736842105263; Page_11780743288481=1; cdb_RW_ID_2105743986=1; Page_Y_22920956172834=12.134868421052632; Page_22920956172834=6; cdb_RW_ID_2026327809=157; Page_Y_18568482786943=-119.39144736842105; Page_18568482786943=1; cdb_RW_ID_2105918627=1; Page_Y_68737983591274=546.0690789473684; Page_68737983591274=3; cdb_RW_ID_2105920627=1; Page_17919678397492=1; cdb_RW_ID_2106009624=1; Page_Y_68037982885276=13.700657894736842; Page_68037982885276=1; cdb_RW_ID_2106132922=2; Page_Y_75620954920800=-52.84539473684211; Page_75620954920800=4; cdb_RW_ID_2106132918=1; Page_Y_41390912980697=-113.12828947368422; Page_41390912980697=1; cdb_RW_ID_2093562922=13; Page_77380461350600=1; cdb_RW_ID_2023400579=3; cdb_RW_ID_2106132873=2; Page_Y_98571382301670=-119.39144736842105; Page_98571382301670=1; Page_Y_67319790177823=-119.39144736842105; Page_67319790177823=1; cdb_RW_ID_2105923789=1; Page_Y_17919678390253=802.858552631579; Page_17919678390253=2; cdb_RW_ID_2051011410=11; Page_Y_86519786766167=-119.39144736842105; Page_86519786766167=1; cdb_RW_ID_2106859899=1; Page_Y_07143786529599=228.2138157894737; Page_07143786529599=3; cdb_RW_ID_2106859898=1; Page_Y_93071382654646=-119.39144736842105; Page_93071382654646=1; cdb_RW_ID_2106859894=1; Page_79220954368387=1; cdb_RW_ID_2105767240=1; showAnnotateTipIf=1; Page_77168140626854=1; cdb_RW_ID_2106862809=1; Page_Y_07143786563589=-144.4736842105263; Page_07143786563589=1; cdb_RW_ID_2106868358=1; Page_87180745959139=1; cdb_RW_ID_1856434724=2; Page_36616584101291=1; cdb_RW_ID_2106868326=1; Page_Y_79220954343204=212.21052631578948; Page_79220954343204=20; cdb_RW_ID_2106868289=1; Page_Y_87180745959096=-94.73684210526316; Page_87180745959096=24; cdb_RW_ID_2106868211=1; cdb_RW_ID_1469485837=1; Page_18947069052514=1; Page_Y_07143786565377=154.73684210526318; Page_07143786565377=6; cdb_RW_ID_2106868195=1; Page_Y_87180745959763=-241.26315789473685; Page_87180745959763=23; cdb_RW_ID_2106131694=1; Page_Y_98571382303249=-119.39144736842105; Page_98571382303249=1; cdb_RW_ID_2080095563=12; cdb_H5R=1; Page_Y_37219757738840=307.28618421052636; Page_37219757738840=1; cdb_RW_ID_2106079090=1; Page_Y_68637982845858=-119.39144736842105; Page_68637982845858=1; cdb_RW_ID_2056545947=160; Page_68643826202904=1; cdb_RW_ID_2107439286=1; Page_79654903756428=1; cdb_RW_ID_2099399935=1; Page_Y_69619733033308=-270.4901315789474; Page_69619733033308=13; Page_Y_77643895975374=-119.39144736842105; cdb_RW_ID_2098918217=8; Page_77643895975374=1; cdb_RW_ID_2098918393=30; Page_Y_77580469679161=-119.39144736842105; Page_77580469679161=1; cdb_RW_ID_2099720859=2; Page_70880466204936=1; cdb_RW_ID_2107904089=1; cdb_RW_ID_2107904087=1; Page_Y_70620951857531=356.60855263157896; Page_70620951857531=11; cdb_RW_ID_2107904086=1; Page_Y_58068146345492=369.91776315789474; Page_58068146345492=4; cdb_RW_ID_2107904961=1; Page_74154903607689=1; cdb_search_format=; searchType=0; cdb_RW_ID_2101884506=3; Page_Y_81868141995042=-196.11513157894737; Page_81868141995042=2; cdb_RW_ID_2065106162=2; Page_Y_15920546954940=-71.63486842105263; Page_15920546954940=2; Page_Y_91771386838646=502.2269736842106; cdb_RW_ID_2105239433=1; Page_Y_03443785878542=176.54276315789474; Page_03443785878542=1; cdb_READED_PC_ID=%2C730730730447440593443448440441440443; cdb_RW_ID_1908563572=2; Page_Y_70487649351320=-119.39144736842105; Page_70487649351320=1; cdb_login_if=1; cdb_uid=104598337; c_login_name=woyoceo; cdb_logined=1; cdb_change_message=1; cdb_msg_num=0; doc88_lt=wx; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d233dd02f2a58c6b2b5841b62253b0c1d47d92f350c3f26f274; cdb_RW_ID_2107957090=1; Page_91271387457848=1; cdb_RW_ID_2108010816=1; Page_58668149414912=1; cdb_RW_ID_1987519680=21; cdb_VIEW_DOC_ID=%2C1987519680; Page_Y_36116352863457=361.3059210526316; Page_36116352863457=2; cdb_RW_ID_2108010801=1; Page_Y_19219675767576=-119.39144736842105; Page_19219675767576=1; siftState=1; show_index=1; cdb_RW_ID_2108010898=2; Page_91771386838646=1; cdb_pageType=2; cdb_msg_time=1777300748"

var DelSaveTimeSleep = 5
var DelNextPageSleep = 10

// ychEduSpider 删除道客巴巴文档
// @Title 删除道客巴巴文档
// @Description https://www.doc88.com/，删除道客巴巴文档
func main() {
	curPage := 1
	for {
		pageListUrl := fmt.Sprintf("https://www.doc88.com/uc/doc_manager.php?act=ajax_doc_list&curpage=%d", curPage)
		fmt.Println(pageListUrl)
		queryDelDoc88ListFormData := QueryDelDoc88ListFormData{
			MenuIndex:  4,
			ClassifyId: "all",
			FolderId:   0,
			Sort:       1,
			Keyword:    "",
			ShowIndex:  1,
		}
		pageListDoc, err := QueryDelDoc88List(pageListUrl, queryDelDoc88ListFormData)
		if err != nil {
			fmt.Println(err)
			continue
		}
		liNodes := htmlquery.Find(pageListDoc, `//div[@id="detailed"]/ul[@class="bookshow3"]/li`)
		if len(liNodes) <= 0 {
			break
		}
		for _, liNode := range liNodes {

			fmt.Println("=======一共有====", len(liNodes), "====文档=======")
			CateLogNameNode := htmlquery.FindOne(liNode, `./div[@class="bookdoc"]/div[@class="posttime"]/span/span[@class="catelog_name"]`)
			if CateLogNameNode == nil {
				continue
			}
			CateLogName := htmlquery.InnerText(CateLogNameNode)
			if strings.Index(CateLogName, "标准") == -1 {
				continue
			}

            TitleNode := htmlquery.FindOne(liNode, `./div[@class="bookdoc"]/h3/a`)
			Title := htmlquery.SelectAttr(TitleNode, "title")
			fmt.Println(Title)
			if strings.Index(Title, "T-") == -1 {
				continue
			}

			PId := htmlquery.SelectAttr(liNode, "id")
			PId = PId[5:]

			fmt.Println("===========开始删除", Title, "===========")
			delUrl := "https://www.doc88.com/uc/index.php?act=del_doc"
			delDoc88FormData := DelDoc88FormData{
				Id:        PId,
			}

			err = DelDoc88(delUrl, delDoc88FormData)
			if err != nil {
				fmt.Println(err)
				continue
			}
			for i := 1; i <= DelSaveTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("page="+strconv.Itoa(curPage)+"===========更新", Title, "成功，暂停", DelSaveTimeSleep, "秒，倒计时", i, "秒===========")
			}
		}
		curPage++
		for i := 1; i <= DelNextPageSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("===========翻", curPage, "页，暂停", DelNextPageSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

func QueryDelDoc88List(requestUrl string, queryDelDoc88ListFormData QueryDelDoc88ListFormData) (doc *html.Node, err error) {
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
	postData := url.Values{}
	postData.Add("menuIndex", strconv.Itoa(queryDelDoc88ListFormData.MenuIndex))
	postData.Add("classify_id", queryDelDoc88ListFormData.ClassifyId)
	postData.Add("folder_id", strconv.Itoa(queryDelDoc88ListFormData.FolderId))
	postData.Add("sort", strconv.Itoa(queryDelDoc88ListFormData.Sort))
	postData.Add("keyword", queryDelDoc88ListFormData.Keyword)
	postData.Add("show_index", strconv.Itoa(queryDelDoc88ListFormData.ShowIndex))
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", DelListCookie)
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

func DelDoc88(requestUrl string, delDoc88FormData DelDoc88FormData) error{
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
	postData := url.Values{}
	postData.Add("id", delDoc88FormData.Id)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", DelDelCookie)
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
		return err
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	return nil
}
