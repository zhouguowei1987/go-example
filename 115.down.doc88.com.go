package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

var DownListCookie = "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; _qddab=3-gv2ozy.lib1y9mi; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; cdb_READED_PC_ID=%2C441442; cdb_back[at]=0; cdb_back[n]=6; cdb_back[book_id]=0; cdb_back[s]=rel; Page_99439638538500=1; cdb_RW_ID_1664816554=2; Page_79399224792334=1; cdb_RW_ID_1664816550=2; Page_91061225912004=1; cdb_back[folder_id]=0; cdb_back[show_index]=1; show_index=1; cdb_back[u]=1; cdb_back[login]=1; cdb_back[txtloginname]=15238369929; cdb_back[txtPassword]=abcdqq123456; cdb_back[inout]=all; cdb_back[withdraw_name]=%E5%91%A8%E5%9B%BD%E4%BC%9F; cdb_back[withdraw_sfz]=410928198704276311; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dall; cdb_RW_ID_1674504163=1; Page_58339246386920=1; cdb_RW_ID_1674504139=1; ShowSkinTip_1=1; BG_SAVE=1; showAnnotateTipIf=1; Page_85629417657928=1; cdb_RW_ID_1674504405=1; cdb_H5R=1; Page_61399254314413=1; Page_Y_60659831344576=-92.77302631578948; cdb_RW_ID_1675722349=3; Page_60659831344576=1; cdb_RW_ID_1643689091=4; Page_Y_30287581596467=503.7927631578948; Page_30287581596467=12; cdb_back[captchaCode]=1; cdb_login_if=1; cdb_uid=104598337; c_login_name=woyoceo; cdb_logined=1; cdb_RW_ID_1676081712=1; Page_Y_97847646857473=-286.1763157894737; Page_97847646857473=2; cdb_RW_ID_1676297250=1; Page_89929414081065=1; cdb_RW_ID_1676298311=1; Page_31161262839711=1; cdb_RW_ID_1659055033=8; Page_29216483788700=1; cdb_RW_ID_1676298283=1; Page_89929414083032=1; cdb_RW_ID_1676563828=1; cdb_RW_ID_572746136=1; Page_0408320285715=1; cdb_back[m]=104598337; cdb_RW_ID_578093873=1; Page_1857576840670=1; Page_Y_31861262027989=79.46381578947368; Page_31861262027989=1; cdb_RW_ID_1676563821=1; cdb_back[uid]=b99ce806c0b55b3bdccae7bc14f8ca3e; Page_Y_69459838185249=-119.39144736842105; Page_69459838185249=1; cdb_RW_ID_1317655069=45; Page_Y_00799895233126=147.1888157894737; Page_00799895233126=1; cdb_back[sharetodoc]=1; cdb_back[download]=2; cdb_back[p_doc_format]=PDF; cdb_RW_ID_1633962014=2; Page_87387511650478=1; cdb_RW_ID_1454752637=2; Page_97316181289402=1; Page_Y_36016424242838=-119.39144736842105; cdb_RW_ID_1676757315=1; Page_Y_69559838313591=-119.39144736842105; Page_69559838313591=1; cdb_RW_ID_1676757334=1; Page_36016424282001=1; cdb_RW_ID_1653590184=1; Page_24887531364798=1; cdb_RW_ID_1676757352=1; Page_Y_97147646424123=63.80592105263158; Page_97147646424123=1; cdb_RW_ID_1676757366=1; Page_67187525232155=1; cdb_RW_ID_1530449813=2; Page_Y_54459150776295=-119.39144736842105; Page_54459150776295=1; cdb_back[id]=1; cdb_RW_ID_1410747206=1; Page_49799491545012=1; cdb_RW_ID_1455699102=40; cdb_back[doctype]=1; cdb_back[action]=data; cdb_back[p_code]=79339633255987; cdb_back[pcode]=79339633255987; cdb_back[ajax]=1; cdb_back[tm]=6681; Page_Y_79339633255987=-119.39144736842105; Page_79339633255987=1; cdb_back[t]=0; cdb_RW_ID_1676767595=2; Page_36016424242838=1; cdb_RW_ID_1633962042=3; Page_69339200527867=1; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A0tXizNXiFNXi2jMW2qvV2q3Q1TMU2LN9or3c3gJACuvi2i3R1jEW1THQ1qkV3iXiFotZBK363m0W1WpkHjFm2O3Q0Ls5HTMU1Oxi2qFh1OsQ1TvV0LsX3iXiDutdHmtSoWlk3jfi0qPU1qk%210T0Q3gU%3D; cdb_back[show_view_type]=1; cdb_back[show]=1; cdb_back[order_num]=1; cdb_back[folder_page]=0; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[member_id]=a9a7e1sN0qHKJRIQKLs1B8X7T%2FFQ83XNo7gQsSxo8i5ll8DXrqqJCk%2F6V5SU9Owq7WA; cdb_back[pm_id]=1486396; cdb_back[friend_id]=0; cdb_back[wxcode]=81763; cdb_back[txt_amount]=53; cdb_back[checkcode]=81763; cdb_back[type]=score; cdb_pageType=2; cdb_back[pcid]=8243; cdb_back[curpage]=1; cdb_back[p_price]=388; cdb_back[p_default_points]=2; cdb_back[doccode]=1676981980; cdb_back[title]=2022%E5%B9%B4%E9%9D%92%E6%B5%B7%E6%B5%B7%E4%B8%9C%E4%B8%AD%E8%80%83%E6%95%B0%E5%AD%A6%E7%9C%9F%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88; cdb_back[intro]=2022%E5%B9%B4%E9%9D%92%E6%B5%B7%E6%B5%B7%E4%B8%9C%E4%B8%AD%E8%80%83%E6%95%B0%E5%AD%A6%E7%9C%9F%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88; cdb_back[p_pagecount]=8; cdb_back[doc_more_id]=1676979002%2C1676978974%2C1676978942%2C1676978926%2C1676978864%2C1676978785%2C1676978718%2C1676978188%2C1676977873%2C1676977843%2C1676977763%2C1676977703%2C1676977650%2C1676977635%2C1676977505%2C1676977341%2C1676977227%2C1676976856%2C; cdb_back[pid]=67887525691221; cdb_RW_ID_1676983773=1; cdb_back[srlid]=4ec0tRQKPZLlZvHOWM8EQ953tCVmGCoWagA8fQMHdmX4RFaHrBiTq4t5+mWiuGjv0oM4b2NabTPsPKhNux9NOhUrnF18wM%2F3+rcz%2FHFuZJC%2F; cdb_back[p_name]=2018%E5%90%89%E6%9E%97%E4%B8%AD%E8%80%83%E6%95%B0%E5%AD%A6%E7%9C%9F%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88; cdb_back[rel_p_id]=1676983773; cdb_back[p_id]=1676983773; cdb_back[page]=1; Page_Y_67887525691221=-119.39144736842105; Page_67887525691221=1; cdb_back[len]=2; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d23985013e508cdcfa515f63502ca44a80ad92f350c3f26f274; doc88_lt=wx; cdb_back[module_type]=6; cdb_back[image_type]=3; cdb_change_message=1; cdb_msg_num=0; siftState=1; cdb_back[state]=all; cdb_back[menuIndex]=2; cdb_back[classify_id]=all; cdb_msg_time=1694654372; cdb_back[act]=ajax_doc_list"
var DownLoadUrlCookie = "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_back[u]=1; show_index=1; cdb_READED_PC_ID=%2C; cdb_back[idcard]=410928198704276311; cdb_back[withdraw_name]=%E5%91%A8%E5%9B%BD%E4%BC%9F; cdb_back[withdraw_sfz]=410928198704276311; cdb_back[phone]=15238369929; cdb_RW_ID_1902409532=1; Page_Y_33173481984501=-119.39144736842105; Page_33173481984501=1; cdb_back[inout]=all; cdb_RW_ID_1904183284=1; Page_Y_11461345197895=-119.39144736842105; Page_11461345197895=1; cdb_RW_ID_1902409455=12; Page_99739587685633=1; cdb_RW_ID_1904503551=1; Page_Y_73747980281227=-119.39144736842105; Page_73747980281227=1; cdb_RW_ID_1904503542=1; Page_Y_69616371870819=-119.39144736842105; Page_69616371870819=1; cdb_RW_ID_1904503509=1; cdb_back[at]=0; cdb_back[n]=6; cdb_back[doctype]=1; cdb_back[book_id]=0; Page_Y_90999614318316=107.64802631578948; Page_90999614318316=2; cdb_RW_ID_1904503498=1; cdb_H5R=1; Page_Y_31373489580946=-119.39144736842105; Page_31373489580946=1; cdb_RW_ID_1904503457=1; Page_Y_73747980281024=360.5230263157895; Page_73747980281024=3; cdb_RW_ID_1880214367=16; Page_18561994815726=1; Page_Y_88261506310890=535.891447368421; cdb_back[complaint_report_id]=413163; cdb_RW_ID_1457915285=18; Page_88261506310890=1; cdb_back[p_code]=88261506310890; cdb_back[showIndex]=1; cdb_back[login]=1; cdb_back[txtPassword]=abcdqq123456; cdb_RW_ID_1905260859=1; Page_Y_97639583728135=244.00065789473683; Page_97639583728135=3; cdb_RW_ID_1905264292=1; Page_Y_18561340825838=-21.529605263157897; Page_18561340825838=3; cdb_RW_ID_1905264195=1; cdb_RW_ID_1667043285=2; Page_79799225148073=1; Page_Y_69116378941638=321.3782894736842; Page_69116378941638=1; cdb_back[wx_code]=5835; cdb_pageType=2; cdb_back[p_default_points]=3; cdb_RW_ID_1905574302=1; Page_94159601137504=1; cdb_RW_ID_1905574317=1; Page_Y_90629856617291=-119.39144736842105; Page_90629856617291=1; cdb_back[pm_id]=1486396; cdb_back[friend_id]=0; cdb_back[wxcode]=53996; cdb_back[txt_amount]=99; cdb_back[checkcode]=53996; cdb_back[uid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[page]=1; Page_Y_90299613536964=-304.1546052631579; Page_90299613536964=2; cdb_RW_ID_1905759091=1; cdb_RW_ID_1654358688=3; Page_Y_64661205709299=-0.5868421052631578; Page_64661205709299=1; cdb_back[member_type]=3; cdb_RW_ID_1631852444=8; Page_98973203651999=1; Page_Y_90429856168589=-238.7828947368421; Page_90429856168589=1; cdb_back[folder_id]=0; cdb_back[show_index]=1; cdb_RW_ID_1904815874=1; ShowSkinTip_1=1; showAnnotateTipIf=1; Page_73347980572540=1; cdb_back[captchaCode]=1; cdb_RW_ID_1904815835=1; Page_Y_90029857396326=121.74013157894738; Page_90029857396326=2; cdb_RW_ID_1415211533=137; Page_Y_86516168966800=216.4703947368421; Page_86516168966800=2; cdb_RW_ID_1905759608=1; cdb_back[id]=1; cdb_RW_ID_1806144666=1; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A0tXizNXiFNXi2jMQ0LM%210TM%210jPT0jF9or3c3gJACuvi2i3R2qPV1Ts51jP%213iXiFotZBK363m1l2O050jHR1Tv%21BLPW0uHQ0Wpm2qEW0WHX1jkUHqpm3iXiDutdHmtSoWlk3jfi0qPU1qk%210T0Q3gU%3D; Page_99499712944222=1; Page_Y_31273485754286=-119.39144736842105; Page_31273485754286=1; cdb_RW_ID_1904182776=2; Page_99839586917442=1; cdb_RW_ID_1905759061=1; Page_Y_90299613536129=-119.39144736842105; Page_90299613536129=1; cdb_RW_ID_1905757620=1; Page_Y_69416378282497=-119.39144736842105; Page_69416378282497=1; cdb_RW_ID_1905759605=1; Page_18261340603240=1; cdb_RW_ID_1905759550=1; cdb_back[ajax]=1; Page_Y_73647982429228=-119.39144736842105; Page_73647982429228=1; cdb_RW_ID_1905759515=1; cdb_back[pcode]=73647982429272; cdb_back[tm]=5256; Page_Y_73647982429272=-119.39144736842105; Page_73647982429272=1; cdb_back[t]=0; cdb_back[txtloginname]=15238369929; cdb_login_if=1; cdb_uid=104598337; c_login_name=woyoceo; cdb_logined=1; cdb_back[doccode]=1905985765; cdb_back[title]=%E9%92%A2%E6%A1%A5%E7%BB%93%E6%9E%84%E8%B4%A8%E9%87%8F%E8%AF%84%E4%BC%B0%E6%95%B0%E6%8D%AE%E7%AE%A1%E7%90%86%E7%B3%BB%E7%BB%9F%28T-CASMES+218%E2%80%942023%29; cdb_back[intro]=%E9%92%A2%E6%A1%A5%E7%BB%93%E6%9E%84%E8%B4%A8%E9%87%8F%E8%AF%84%E4%BC%B0%E6%95%B0%E6%8D%AE%E7%AE%A1%E7%90%86%E7%B3%BB%E7%BB%9F%28T-CASMES+218%E2%80%942023%29; cdb_back[pcid]=8370; cdb_back[sharetodoc]=1; cdb_back[download]=2; cdb_back[p_price]=488; cdb_back[p_pagecount]=13; cdb_back[p_doc_format]=PDF; Page_Y_31773485464074=156.57894736842107; cdb_RW_ID_1905989027=1; Page_Y_31773485464817=-119.39144736842105; Page_31773485464817=1; cdb_back[doc_more_id]=1905989027%2C; siftState=1; cdb_RW_ID_1905989372=1; cdb_back[m]=104598337; cdb_back[len]=3; Page_Y_70287643696120=-45.799342105263165; Page_70287643696120=2; cdb_back[classify_id]=all; cdb_change_message=1; cdb_msg_num=0; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dall; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d230b015e3dd0c066c9b9112d0717f149f1d92f350c3f26f274; doc88_lt=wx; cdb_back[module_type]=7; cdb_back[image_type]=4; cdb_back[type]=score; cdb_back[p_name]=%E5%9B%9B%E5%B7%9D%E7%9C%81%E8%B5%84%E9%98%B3%E5%B8%82%E9%9B%81%E6%B1%9F%E5%8C%BA%E4%B9%9D%E5%B9%B4%E7%BA%A7%28%E4%B8%8A%29%E6%9C%9F%E6%9C%AB%E7%89%A9%E7%90%86%E8%AF%95%E5%8D%B7%28%E8%A7%A3%E6%9E%90%E7%89%88%29; cdb_back[rel_p_id]=1905989379; cdb_back[p_id]=1905989379; cdb_back[pid]=31773485464074; cdb_RW_ID_1905989379=3; cdb_back[srlid]=43004MveWY2TDUOb3aB2uWjvzFt9eJ5vsmi9yM2XczbGwuFYBdEhVvyi64HNipHTvRSnrNaKBynWkl4L2Es41gLgUOd59yvvusRnOkmqaXAI; Page_31773485464074=1; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[member_id]=3642hn3bVJTpRlis0CbLQ23iD5f7jAO2oOEhnFTIRs5TzLhbZ74AKWBzUQdaI6at1BI; cdb_back[show]=1; cdb_back[order_num]=1; cdb_back[folder_page]=0; cdb_back[state]=myshare; cdb_back[menuIndex]=4; cdb_msg_time=1701956410; cdb_back[act]=ajax_doc_list; cdb_back[curpage]=4"

var DownNextPageSleep = 10
var TodayCurrentDownLoadCount = 0
var TodayMaxDownLoadCount = 20

// ychEduSpider 下载道客巴巴文档
// @Title 下载道客巴巴文档
// @Description https://www.doc88.com/，下载道客巴巴文档
func main() {
	curPage := 510
	isPageListGo := true
	for isPageListGo {
		pageListUrl := fmt.Sprintf("https://www.doc88.com/uc/doc_manager.php?act=ajax_doc_list&curpage=%d", curPage)
		fmt.Println(pageListUrl)
		queryDownDoc88ListFormData := QueryDownDoc88ListFormData{
			MenuIndex:  4,
			ClassifyId: "all",
			FolderId:   0,
			Sort:       1,
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

			filePath := "../down.doc88.com/" + Title + ".pdf"
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
