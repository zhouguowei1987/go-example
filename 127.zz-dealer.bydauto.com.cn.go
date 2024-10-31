package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
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

type QueryEditBydAutoListFormData struct {
	MenuIndex  int
	ClassifyId string
	FolderId   int
	Sort       int
	Keyword    string
	ShowIndex  int
}

type EditBydAutoResponseData struct {
	Result     string `json:"result"`
	EditTitle  string `json:"edit_title"`
	Class      string `json:"class"`
	UpdateInfo string `json:"updateinfo"`
	State      string `json:"state"`
	SaveFile   string `json:"savefile"`
	Other      string `json:"other"`
}

type EditBydAutoFormData struct {
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

//var EditDetailCookie = "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; _qddab=3-gv2ozy.lib1y9mi; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; cdb_READED_PC_ID=%2C441442; cdb_back[at]=0; cdb_back[n]=6; cdb_back[book_id]=0; cdb_back[s]=rel; Page_99439638538500=1; cdb_RW_ID_1664816554=2; Page_79399224792334=1; cdb_RW_ID_1664816550=2; Page_91061225912004=1; cdb_back[folder_id]=0; cdb_back[show_index]=1; show_index=1; cdb_back[u]=1; cdb_back[login]=1; cdb_back[txtloginname]=15238369929; cdb_back[txtPassword]=abcdqq123456; cdb_back[inout]=all; cdb_back[withdraw_name]=%E5%91%A8%E5%9B%BD%E4%BC%9F; cdb_back[withdraw_sfz]=410928198704276311; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dall; cdb_RW_ID_1674504163=1; Page_58339246386920=1; cdb_RW_ID_1674504139=1; ShowSkinTip_1=1; BG_SAVE=1; showAnnotateTipIf=1; Page_85629417657928=1; cdb_RW_ID_1674504405=1; cdb_H5R=1; Page_61399254314413=1; Page_Y_60659831344576=-92.77302631578948; cdb_RW_ID_1675722349=3; Page_60659831344576=1; cdb_RW_ID_1643689091=4; Page_Y_30287581596467=503.7927631578948; Page_30287581596467=12; cdb_back[captchaCode]=1; cdb_login_if=1; cdb_uid=104598337; c_login_name=woyoceo; cdb_logined=1; cdb_RW_ID_1676081712=1; Page_Y_97847646857473=-286.1763157894737; Page_97847646857473=2; cdb_RW_ID_1676297250=1; Page_89929414081065=1; cdb_RW_ID_1676298311=1; Page_31161262839711=1; cdb_RW_ID_1659055033=8; Page_29216483788700=1; cdb_RW_ID_1676298283=1; Page_89929414083032=1; cdb_RW_ID_1676563828=1; cdb_RW_ID_572746136=1; Page_0408320285715=1; cdb_back[m]=104598337; cdb_RW_ID_578093873=1; Page_1857576840670=1; Page_Y_31861262027989=79.46381578947368; Page_31861262027989=1; cdb_RW_ID_1676563821=1; cdb_back[uid]=b99ce806c0b55b3bdccae7bc14f8ca3e; Page_Y_69459838185249=-119.39144736842105; Page_69459838185249=1; cdb_RW_ID_1317655069=45; Page_Y_00799895233126=147.1888157894737; Page_00799895233126=1; cdb_back[sharetodoc]=1; cdb_back[download]=2; cdb_back[p_doc_format]=PDF; cdb_RW_ID_1633962014=2; Page_87387511650478=1; cdb_RW_ID_1454752637=2; Page_97316181289402=1; Page_Y_36016424242838=-119.39144736842105; cdb_RW_ID_1676757315=1; Page_Y_69559838313591=-119.39144736842105; Page_69559838313591=1; cdb_RW_ID_1676757334=1; Page_36016424282001=1; cdb_RW_ID_1653590184=1; Page_24887531364798=1; cdb_RW_ID_1676757352=1; Page_Y_97147646424123=63.80592105263158; Page_97147646424123=1; cdb_RW_ID_1676757366=1; Page_67187525232155=1; cdb_RW_ID_1530449813=2; Page_Y_54459150776295=-119.39144736842105; Page_54459150776295=1; cdb_back[id]=1; cdb_RW_ID_1410747206=1; Page_49799491545012=1; cdb_RW_ID_1455699102=40; cdb_back[doctype]=1; cdb_back[action]=data; cdb_back[p_code]=79339633255987; cdb_back[pcode]=79339633255987; cdb_back[ajax]=1; cdb_back[tm]=6681; Page_Y_79339633255987=-119.39144736842105; Page_79339633255987=1; cdb_back[t]=0; cdb_RW_ID_1676767595=2; Page_36016424242838=1; cdb_RW_ID_1633962042=3; Page_69339200527867=1; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A0tXizNXiFNXi2jMW2qvV2q3Q1TMU2LN9or3c3gJACuvi2i3R1jEW1THQ1qkV3iXiFotZBK363m0W1WpkHjFm2O3Q0Ls5HTMU1Oxi2qFh1OsQ1TvV0LsX3iXiDutdHmtSoWlk3jfi0qPU1qk%210T0Q3gU%3D; cdb_back[show_view_type]=1; cdb_back[show]=1; cdb_back[order_num]=1; cdb_back[folder_page]=0; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[member_id]=a9a7e1sN0qHKJRIQKLs1B8X7T%2FFQ83XNo7gQsSxo8i5ll8DXrqqJCk%2F6V5SU9Owq7WA; cdb_back[pm_id]=1486396; cdb_back[friend_id]=0; cdb_back[wxcode]=81763; cdb_back[txt_amount]=53; cdb_back[checkcode]=81763; cdb_back[type]=score; cdb_pageType=2; cdb_back[pcid]=8243; cdb_back[p_price]=388; cdb_back[p_default_points]=2; cdb_back[doccode]=1676981980; cdb_back[title]=2022%E5%B9%B4%E9%9D%92%E6%B5%B7%E6%B5%B7%E4%B8%9C%E4%B8%AD%E8%80%83%E6%95%B0%E5%AD%A6%E7%9C%9F%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88; cdb_back[intro]=2022%E5%B9%B4%E9%9D%92%E6%B5%B7%E6%B5%B7%E4%B8%9C%E4%B8%AD%E8%80%83%E6%95%B0%E5%AD%A6%E7%9C%9F%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88; cdb_back[p_pagecount]=8; cdb_back[doc_more_id]=1676979002%2C1676978974%2C1676978942%2C1676978926%2C1676978864%2C1676978785%2C1676978718%2C1676978188%2C1676977873%2C1676977843%2C1676977763%2C1676977703%2C1676977650%2C1676977635%2C1676977505%2C1676977341%2C1676977227%2C1676976856%2C; cdb_back[pid]=67887525691221; cdb_RW_ID_1676983773=1; cdb_back[srlid]=4ec0tRQKPZLlZvHOWM8EQ953tCVmGCoWagA8fQMHdmX4RFaHrBiTq4t5+mWiuGjv0oM4b2NabTPsPKhNux9NOhUrnF18wM%2F3+rcz%2FHFuZJC%2F; cdb_back[p_name]=2018%E5%90%89%E6%9E%97%E4%B8%AD%E8%80%83%E6%95%B0%E5%AD%A6%E7%9C%9F%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88; cdb_back[rel_p_id]=1676983773; cdb_back[p_id]=1676983773; cdb_back[page]=1; Page_Y_67887525691221=-119.39144736842105; Page_67887525691221=1; cdb_back[len]=2; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d23985013e508cdcfa515f63502ca44a80ad92f350c3f26f274; doc88_lt=wx; cdb_back[module_type]=6; cdb_back[image_type]=3; cdb_change_message=1; cdb_msg_num=0; siftState=1; cdb_back[state]=all; cdb_back[menuIndex]=2; cdb_back[classify_id]=all; cdb_msg_time=1694654372; cdb_back[act]=ajax_doc_list; cdb_back[curpage]=2"
//var EditEditCookie = "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_RW_ID_1903247601=1; Page_66416370912476=1; cdb_RW_ID_1903247574=1; Page_Y_99429852071617=-61.45723684210527; Page_99429852071617=3; cdb_RW_ID_1903247663=1; Page_Y_33273480197220=-119.39144736842105; Page_33273480197220=1; cdb_RW_ID_1903247308=1; Page_Y_66416370912075=117.82565789473685; Page_66416370912075=3; cdb_back[idcard]=410928198704276311; cdb_RW_ID_1903247654=1; Page_Y_66416370912481=251.70065789473685; Page_66416370912481=1; cdb_RW_ID_1899594632=2; Page_33373644549201=1; cdb_RW_ID_1899595345=1; cdb_READED_PC_ID=%2C441; Page_11161933030750=1; cdb_RW_ID_1900309002=1; Page_66116377073779=1; cdb_RW_ID_1899595045=5; Page_99999766363143=1; cdb_RW_ID_1900557749=1; Page_Y_33573488557794=-118.78947368421052; Page_33573488557794=1; cdb_back[inout]=all; cdb_back[wxcode]=73705; cdb_back[txt_amount]=67; cdb_back[checkcode]=73705; cdb_RW_ID_1902380443=1; Page_66916379057110=1; cdb_RW_ID_1901276810=3; cdb_back[at]=0; cdb_back[n]=6; cdb_back[doctype]=1; cdb_back[book_id]=0; Page_Y_33773483172638=-106.86513157894737; Page_33773483172638=1; cdb_H5R=1; Page_Y_99959604991344=-138.1809210526316; cdb_RW_ID_1902115626=1; Page_99959604991848=1; cdb_back[mid]=adf128e2651d9b9496bc70433c718c6f; cdb_back[member_type]=3; cdb_back[order_num]=1; cdb_back[folder_page]=0; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dall; cdb_back[login]=1; cdb_back[captchaCode]=1; cdb_RW_ID_1900557083=1; Page_Y_11061344006497=291.6282894736842; Page_11061344006497=1; cdb_RW_ID_1900557513=6; Page_99629855661692=1; cdb_RW_ID_1903524659=1; Page_33773480519254=1; cdb_RW_ID_1900557168=3; Page_99339588334921=1; cdb_RW_ID_1900557187=1; Page_77247988224754=1; cdb_RW_ID_1900557113=1; Page_33573488557330=1; cdb_RW_ID_1415594415=27; Page_30199493364493=1; cdb_RW_ID_1900307571=2; Page_66116377072826=1; cdb_RW_ID_1647369237=2; Page_28673297024107=1; cdb_RW_ID_1414552096=11; Page_69729797660584=1; cdb_RW_ID_1900557221=1; Page_Y_99629855661009=-118.78947368421052; Page_99629855661009=1; cdb_RW_ID_1903679267=1; cdb_RW_ID_1634547236=2; Page_Y_51961275056872=1028.5578947368422; Page_51961275056872=4; Page_Y_66516370423942=674.2421052631579; Page_66516370423942=4; cdb_back[pm_id]=1486396; cdb_back[friend_id]=0; cdb_RW_ID_1903882245=2; Page_Y_66316370559918=8.220394736842106; Page_66316370559918=1; cdb_RW_ID_1807498535=1; Page_Y_74487942869313=-138.1809210526316; Page_74487942869313=1; cdb_back[u]=1; show_index=1; cdb_back[folder_id]=0; cdb_back[show_index]=1; cdb_pageType=2; cdb_RW_ID_1899802279=3; Page_66916533579923=1; cdb_RW_ID_1903228797=1; Page_11061347889636=1; cdb_RW_ID_1902115722=2; Page_99959604991344=1; cdb_show_type=1; cdb_back[show]=1; cdb_RW_ID_1903886226=1; cdb_back[uid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[page]=1; Page_Y_77947981556336=16.832236842105264; Page_77947981556336=4; cdb_RW_ID_1903886218=1; Page_Y_99829852334093=-119.39144736842105; Page_99829852334093=1; cdb_RW_ID_1636908912=3; Page_67139202581597=1; cdb_RW_ID_1902584237=9; cdb_back[action]=data; cdb_back[p_code]=77187640398012; Page_Y_77187640398012=-119.39144736842105; Page_77187640398012=1; cdb_RW_ID_1902584222=5; ShowSkinTip_1=1; showAnnotateTipIf=1; Page_Y_99559604127444=-290.0625; Page_99559604127444=5; cdb_back[pcode]=99559604127444; cdb_back[ajax]=1; cdb_back[tm]=6216; cdb_back[member_id]=104598337; cdb_back[t]=0; cdb_RW_ID_1903886327=1; cdb_back[m]=0; cdb_back[s]=rel; cdb_back[id]=3; cdb_RW_ID_1632659854=2; cdb_back[len]=1; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A0VXizNXiFNXi2jMQ0LMS0j0V1qsW2LB9or3c3gJACuvi2i3R2qPT2LnW0T3Q3iXiFotZBK363mxi1TMW1T0UHqBk2LP%212qkTHutjBO0XBq3%21Hqn5HqHX3iXiDutdHmtSoWlk3jfXAv%3D%3D; Page_Y_41799280236734=-138.1809210526316; Page_41799280236734=1; Page_Y_77687641995102=-191.41776315789474; Page_77687641995102=2; cdb_back[pid]=99539580112081; cdb_RW_ID_1903886308=1; cdb_back[p_name]=%E4%B8%AD%E8%80%83%E8%AF%AD%E6%96%87%E5%8F%A4%E8%AF%97%E6%96%87%E7%B2%BE%E8%AE%B2%E5%B7%A7%E7%BB%83-%E3%80%8A%E6%B2%81%E5%9B%AD%E6%98%A5%C2%B7%E9%9B%AA%E3%80%8B; cdb_back[rel_p_id]=1903886308; cdb_back[srlid]=e54dhpGlukweNAFxNdCqThgWDAufQeMSlPO3ycJYLTntLc7K0QbkFeOAd25OjYXjkZ12pT9UmcBSsad4meR+OHKZ9NTlSsKrIAYFlv2k4a21; Page_99539580112081=1; cdb_change_message=1; cdb_msg_num=0; cdb_back[txtloginname]=15238369929; cdb_back[txtPassword]=abcdqq123456; cdb_login_if=1; cdb_uid=104598337; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d23f073f81cf508131230df2c63a1812d0ad92f350c3f26f274; c_login_name=woyoceo; cdb_logined=1; doc88_lt=wx; cdb_back[image_type]=9; cdb_back[module_type]=7; cdb_back[doc_more_id]=1903427930%2C; cdb_back[curpage]=1; siftState=1; cdb_back[classify_id]=all; cdb_back[type]=1; cdb_back[state]=all; cdb_msg_time=1701246083; cdb_back[menuIndex]=2; cdb_back[act]=getDocInfo; cdb_back[p_id]=1903886332"
//
//var EditDetailTimeSleep = 30
//var EditSaveTimeSleep = 20
//var EditNextPageSleep = 15

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

// ychEduSpider 编辑智蛛AI经销商系统
// @Title 编辑智蛛AI经销商系统
// @Description https://zz-dealer.bydauto.com.cn/，编辑智蛛AI经销商系统
func main() {
	curPage := 0
	for {
		listRequestPayload := make(map[string]interface{})
		listRequestPayload["activityType"] = 0
		listRequestPayload["dateEnd"] = 0
		listRequestPayload["dateStart"] = 0
		listRequestPayload["dealerId"] = 826
		listRequestPayload["filterType"] = 0
		listRequestPayload["fromType"] = 0
		listRequestPayload["key"] = ""
		listRequestPayload["level"] = ""
		listRequestPayload["onlyTotal"] = false
		listRequestPayload["pageCount"] = 10
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
		fmt.Printf("%+v", queryEditBydAutoResponseList.Data)
		os.Exit(1)
		//liNodes := htmlquery.Find(pageListDoc, `//div[@id="detailed"]/ul[@class="bookshow3"]/li`)
		//if len(liNodes) <= 0 {
		//	break
		//}
		//for _, liNode := range liNodes {
		//
		//	TitleNode := htmlquery.FindOne(liNode, `./div[@class="bookdoc"]/h3/a`)
		//	Title := htmlquery.SelectAttr(TitleNode, "title")
		//	fmt.Println(Title)
		//
		//	IntroNode := htmlquery.FindOne(liNode, `./div[@class="bookdoc"]/p`)
		//	Intro := htmlquery.InnerText(IntroNode)
		//
		//	PPageCountNode := htmlquery.FindOne(liNode, `./div[@class="bookimg"]/em`)
		//	PPageCount := htmlquery.InnerText(PPageCountNode)
		//	PPageCount = PPageCount[2:]
		//
		//	PPriceNode := htmlquery.FindOne(liNode, `./div[@class="bookdoc"]/ul[@class="position"]/li[6]/span[@class="jifentip"]/strong[@class="red"]`)
		//	PPrice := htmlquery.InnerText(PPriceNode)
		//
		//	filePageNum, _ := strconv.Atoi(PPageCount)
		//	PPriceNew := ""
		//	// 根据页数设置价格
		//	if filePageNum > 0 && filePageNum <= 5 {
		//		PPriceNew = "288"
		//	} else if filePageNum > 5 && filePageNum <= 10 {
		//		PPriceNew = "388"
		//	} else if filePageNum > 10 && filePageNum <= 15 {
		//		PPriceNew = "488"
		//	} else if filePageNum > 15 && filePageNum <= 20 {
		//		PPriceNew = "588"
		//	} else if filePageNum > 20 && filePageNum <= 25 {
		//		PPriceNew = "688"
		//	} else if filePageNum > 25 && filePageNum <= 30 {
		//		PPriceNew = "788"
		//	} else if filePageNum > 30 && filePageNum <= 35 {
		//		PPriceNew = "888"
		//	} else if filePageNum > 35 && filePageNum <= 40 {
		//		PPriceNew = "988"
		//	} else if filePageNum > 40 && filePageNum <= 45 {
		//		PPriceNew = "1088"
		//	} else if filePageNum > 45 && filePageNum <= 50 {
		//		PPriceNew = "1188"
		//	} else {
		//		PPriceNew = "1288"
		//	}
		//
		//	// 新旧价格一样，则跳过
		//	fmt.Println(PPrice, PPriceNew)
		//	if PPrice == PPriceNew {
		//		continue
		//	}
		//
		//	PId := htmlquery.SelectAttr(liNode, "id")
		//	PId = PId[5:]
		//
		//	for i := 1; i <= EditDetailTimeSleep; i++ {
		//		time.Sleep(time.Second)
		//		fmt.Println("page="+strconv.Itoa(curPage)+"===========获取", Title, "详情暂停", EditDetailTimeSleep, "秒，倒计时", i, "秒===========")
		//	}
		//
		//	detailUrl := "https://zz-dealer.bydauto.com.cn/uc/usr_doc_manager.php?act=getDocInfo"
		//	detailDoc, err := QueryEditBydAutoDetail(detailUrl, PId)
		//	if err != nil {
		//		EditBydAutoHttpProxyUrl = ""
		//		fmt.Println(err)
		//		continue
		//	}
		//
		//	DocCodeNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[1]/div[@class="booksedit"]/table[@class="edit-table"]/input`)
		//	DocCode := htmlquery.SelectAttr(DocCodeNode, "value")
		//
		//	PCidNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[1]/div[@class="booksedit"]/table[@class="edit-table"]/tbody/tr[3]/td[2]/div[@class="layers"]/input`)
		//	PCid := htmlquery.SelectAttr(PCidNode, "value")
		//
		//	PDocFormatNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[2]/div[@class="booksedit booksedit-bdr"]/table[@class="edit-table"]/tbody/tr[3]/td[2]/input[3]`)
		//	PDocFormat := htmlquery.SelectAttr(PDocFormatNode, "value")
		//
		//	fmt.Println("===========开始修改", Title, "价格===========")
		//	editUrl := "https://zz-dealer.bydauto.com.cn/uc/index.php"
		//	editDoc88FormData := EditBydAutoFormData{
		//		DocCode:        DocCode,
		//		Title:          Title,
		//		Intro:          Intro,
		//		PCid:           PCid,
		//		Keyword:        "",
		//		ShareToDoc:     "1",
		//		Download:       "2",
		//		PPrice:         PPriceNew,
		//		PDefaultPoints: "3",
		//		PPageCount:     PPageCount,
		//		PDocFormat:     PDocFormat,
		//		Act:            "save_info",
		//		GroupList:      "",
		//		GroupFreeList:  "",
		//	}
		//
		//	for i := 1; i <= EditSaveTimeSleep; i++ {
		//		time.Sleep(time.Second)
		//		fmt.Println("page="+strconv.Itoa(curPage)+"===========更新", Title, "成功，暂停", EditSaveTimeSleep, "秒，倒计时", i, "秒===========")
		//	}
		//
		//	_, err = EditBydAuto(editUrl, editDoc88FormData)
		//	if err != nil {
		//		EditBydAutoHttpProxyUrl = ""
		//		fmt.Println(err)
		//		continue
		//	}
		//}
		//curPage++
		//for i := 1; i <= EditNextPageSleep; i++ {
		//	time.Sleep(time.Second)
		//	fmt.Println("===========翻", curPage, "页，暂停", EditNextPageSleep, "秒，倒计时", i, "秒===========")
		//}
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
	req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50VHlwZSI6MSwiaWQiOjY0MzIyLCJpc1N1cGVyIjpmYWxzZX0.IiINeGVqTZTqE9zHvACPX__Qu1A9YB4916lMXAumjIc")
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

func QueryEditBydAutoDetail(requestUrl string, PId string) (doc *html.Node, err error) {
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
	postData := url.Values{}
	postData.Add("p_id", PId)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("Cookie", EditDetailCookie)
	req.Header.Set("Host", "zz-dealer.bydauto.com.cn")
	req.Header.Set("Origin", "https://zz-dealer.bydauto.com.cn")
	req.Header.Set("Referer", "https://zz-dealer.bydauto.com.cn/uc/doc_manager.php?act=doc_list&state=all")
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

func EditBydAuto(requestUrl string, editDoc88FormData EditBydAutoFormData) (editDoc88ResponseData EditBydAutoResponseData, err error) {
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
	editDoc88ResponseData = EditBydAutoResponseData{}
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
	//req.Header.Set("Cookie", EditEditCookie)
	req.Header.Set("Host", "zz-dealer.bydauto.com.cn")
	req.Header.Set("Origin", "https://zz-dealer.bydauto.com.cn")
	req.Header.Set("Referer", "https://zz-dealer.bydauto.com.cn/uc/doc_manager.php?act=doc_list&state=all")
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
