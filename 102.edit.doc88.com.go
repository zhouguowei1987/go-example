package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	EditDoc88EnableHttpProxy = false
	EditDoc88HttpProxyUrl    = "111.225.152.186:8089"
)

func EditDoc88SetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(EditDoc88HttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
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

var ListCookie = "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; _qddab=3-gv2ozy.lib1y9mi; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; cdb_READED_PC_ID=%2C441442; cdb_back[at]=0; cdb_back[n]=6; cdb_back[book_id]=0; cdb_back[s]=rel; Page_99439638538500=1; cdb_RW_ID_1664816554=2; Page_79399224792334=1; cdb_RW_ID_1664816550=2; Page_91061225912004=1; cdb_back[folder_id]=0; cdb_back[show_index]=1; show_index=1; cdb_back[u]=1; cdb_back[login]=1; cdb_back[txtloginname]=15238369929; cdb_back[txtPassword]=abcdqq123456; cdb_back[inout]=all; cdb_back[withdraw_name]=%E5%91%A8%E5%9B%BD%E4%BC%9F; cdb_back[withdraw_sfz]=410928198704276311; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dall; cdb_RW_ID_1674504163=1; Page_58339246386920=1; cdb_RW_ID_1674504139=1; ShowSkinTip_1=1; BG_SAVE=1; showAnnotateTipIf=1; Page_85629417657928=1; cdb_RW_ID_1674504405=1; cdb_H5R=1; Page_61399254314413=1; Page_Y_60659831344576=-92.77302631578948; cdb_RW_ID_1675722349=3; Page_60659831344576=1; cdb_RW_ID_1643689091=4; Page_Y_30287581596467=503.7927631578948; Page_30287581596467=12; cdb_back[captchaCode]=1; cdb_login_if=1; cdb_uid=104598337; c_login_name=woyoceo; cdb_logined=1; cdb_RW_ID_1676081712=1; Page_Y_97847646857473=-286.1763157894737; Page_97847646857473=2; cdb_RW_ID_1676297250=1; Page_89929414081065=1; cdb_RW_ID_1676298311=1; Page_31161262839711=1; cdb_RW_ID_1659055033=8; Page_29216483788700=1; cdb_RW_ID_1676298283=1; Page_89929414083032=1; cdb_RW_ID_1676563828=1; cdb_RW_ID_572746136=1; Page_0408320285715=1; cdb_back[m]=104598337; cdb_RW_ID_578093873=1; Page_1857576840670=1; Page_Y_31861262027989=79.46381578947368; Page_31861262027989=1; cdb_RW_ID_1676563821=1; cdb_back[uid]=b99ce806c0b55b3bdccae7bc14f8ca3e; Page_Y_69459838185249=-119.39144736842105; Page_69459838185249=1; cdb_RW_ID_1317655069=45; Page_Y_00799895233126=147.1888157894737; Page_00799895233126=1; cdb_back[sharetodoc]=1; cdb_back[download]=2; cdb_back[p_doc_format]=PDF; cdb_RW_ID_1633962014=2; Page_87387511650478=1; cdb_RW_ID_1454752637=2; Page_97316181289402=1; Page_Y_36016424242838=-119.39144736842105; cdb_RW_ID_1676757315=1; Page_Y_69559838313591=-119.39144736842105; Page_69559838313591=1; cdb_RW_ID_1676757334=1; Page_36016424282001=1; cdb_RW_ID_1653590184=1; Page_24887531364798=1; cdb_RW_ID_1676757352=1; Page_Y_97147646424123=63.80592105263158; Page_97147646424123=1; cdb_RW_ID_1676757366=1; Page_67187525232155=1; cdb_RW_ID_1530449813=2; Page_Y_54459150776295=-119.39144736842105; Page_54459150776295=1; cdb_back[id]=1; cdb_RW_ID_1410747206=1; Page_49799491545012=1; cdb_RW_ID_1455699102=40; cdb_back[doctype]=1; cdb_back[action]=data; cdb_back[p_code]=79339633255987; cdb_back[pcode]=79339633255987; cdb_back[ajax]=1; cdb_back[tm]=6681; Page_Y_79339633255987=-119.39144736842105; Page_79339633255987=1; cdb_back[t]=0; cdb_RW_ID_1676767595=2; Page_36016424242838=1; cdb_RW_ID_1633962042=3; Page_69339200527867=1; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A0tXizNXiFNXi2jMW2qvV2q3Q1TMU2LN9or3c3gJACuvi2i3R1jEW1THQ1qkV3iXiFotZBK363m0W1WpkHjFm2O3Q0Ls5HTMU1Oxi2qFh1OsQ1TvV0LsX3iXiDutdHmtSoWlk3jfi0qPU1qk%210T0Q3gU%3D; cdb_back[show_view_type]=1; cdb_back[show]=1; cdb_back[order_num]=1; cdb_back[folder_page]=0; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[member_id]=a9a7e1sN0qHKJRIQKLs1B8X7T%2FFQ83XNo7gQsSxo8i5ll8DXrqqJCk%2F6V5SU9Owq7WA; cdb_back[pm_id]=1486396; cdb_back[friend_id]=0; cdb_back[wxcode]=81763; cdb_back[txt_amount]=53; cdb_back[checkcode]=81763; cdb_back[type]=score; cdb_pageType=2; cdb_back[pcid]=8243; cdb_back[curpage]=1; cdb_back[p_price]=388; cdb_back[p_default_points]=2; cdb_back[doccode]=1676981980; cdb_back[title]=2022%E5%B9%B4%E9%9D%92%E6%B5%B7%E6%B5%B7%E4%B8%9C%E4%B8%AD%E8%80%83%E6%95%B0%E5%AD%A6%E7%9C%9F%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88; cdb_back[intro]=2022%E5%B9%B4%E9%9D%92%E6%B5%B7%E6%B5%B7%E4%B8%9C%E4%B8%AD%E8%80%83%E6%95%B0%E5%AD%A6%E7%9C%9F%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88; cdb_back[p_pagecount]=8; cdb_back[doc_more_id]=1676979002%2C1676978974%2C1676978942%2C1676978926%2C1676978864%2C1676978785%2C1676978718%2C1676978188%2C1676977873%2C1676977843%2C1676977763%2C1676977703%2C1676977650%2C1676977635%2C1676977505%2C1676977341%2C1676977227%2C1676976856%2C; cdb_back[pid]=67887525691221; cdb_RW_ID_1676983773=1; cdb_back[srlid]=4ec0tRQKPZLlZvHOWM8EQ953tCVmGCoWagA8fQMHdmX4RFaHrBiTq4t5+mWiuGjv0oM4b2NabTPsPKhNux9NOhUrnF18wM%2F3+rcz%2FHFuZJC%2F; cdb_back[p_name]=2018%E5%90%89%E6%9E%97%E4%B8%AD%E8%80%83%E6%95%B0%E5%AD%A6%E7%9C%9F%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88; cdb_back[rel_p_id]=1676983773; cdb_back[p_id]=1676983773; cdb_back[page]=1; Page_Y_67887525691221=-119.39144736842105; Page_67887525691221=1; cdb_back[len]=2; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d23985013e508cdcfa515f63502ca44a80ad92f350c3f26f274; doc88_lt=wx; cdb_back[module_type]=6; cdb_back[image_type]=3; cdb_change_message=1; cdb_msg_num=0; siftState=1; cdb_back[state]=all; cdb_back[menuIndex]=2; cdb_back[classify_id]=all; cdb_msg_time=1694654372; cdb_back[act]=ajax_doc_list"
var DetailCookie = "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; _qddab=3-gv2ozy.lib1y9mi; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; cdb_READED_PC_ID=%2C441442; cdb_back[at]=0; cdb_back[n]=6; cdb_back[book_id]=0; cdb_back[s]=rel; Page_99439638538500=1; cdb_RW_ID_1664816554=2; Page_79399224792334=1; cdb_RW_ID_1664816550=2; Page_91061225912004=1; cdb_back[folder_id]=0; cdb_back[show_index]=1; show_index=1; cdb_back[u]=1; cdb_back[login]=1; cdb_back[txtloginname]=15238369929; cdb_back[txtPassword]=abcdqq123456; cdb_back[inout]=all; cdb_back[withdraw_name]=%E5%91%A8%E5%9B%BD%E4%BC%9F; cdb_back[withdraw_sfz]=410928198704276311; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dall; cdb_RW_ID_1674504163=1; Page_58339246386920=1; cdb_RW_ID_1674504139=1; ShowSkinTip_1=1; BG_SAVE=1; showAnnotateTipIf=1; Page_85629417657928=1; cdb_RW_ID_1674504405=1; cdb_H5R=1; Page_61399254314413=1; Page_Y_60659831344576=-92.77302631578948; cdb_RW_ID_1675722349=3; Page_60659831344576=1; cdb_RW_ID_1643689091=4; Page_Y_30287581596467=503.7927631578948; Page_30287581596467=12; cdb_back[captchaCode]=1; cdb_login_if=1; cdb_uid=104598337; c_login_name=woyoceo; cdb_logined=1; cdb_RW_ID_1676081712=1; Page_Y_97847646857473=-286.1763157894737; Page_97847646857473=2; cdb_RW_ID_1676297250=1; Page_89929414081065=1; cdb_RW_ID_1676298311=1; Page_31161262839711=1; cdb_RW_ID_1659055033=8; Page_29216483788700=1; cdb_RW_ID_1676298283=1; Page_89929414083032=1; cdb_RW_ID_1676563828=1; cdb_RW_ID_572746136=1; Page_0408320285715=1; cdb_back[m]=104598337; cdb_RW_ID_578093873=1; Page_1857576840670=1; Page_Y_31861262027989=79.46381578947368; Page_31861262027989=1; cdb_RW_ID_1676563821=1; cdb_back[uid]=b99ce806c0b55b3bdccae7bc14f8ca3e; Page_Y_69459838185249=-119.39144736842105; Page_69459838185249=1; cdb_RW_ID_1317655069=45; Page_Y_00799895233126=147.1888157894737; Page_00799895233126=1; cdb_back[sharetodoc]=1; cdb_back[download]=2; cdb_back[p_doc_format]=PDF; cdb_RW_ID_1633962014=2; Page_87387511650478=1; cdb_RW_ID_1454752637=2; Page_97316181289402=1; Page_Y_36016424242838=-119.39144736842105; cdb_RW_ID_1676757315=1; Page_Y_69559838313591=-119.39144736842105; Page_69559838313591=1; cdb_RW_ID_1676757334=1; Page_36016424282001=1; cdb_RW_ID_1653590184=1; Page_24887531364798=1; cdb_RW_ID_1676757352=1; Page_Y_97147646424123=63.80592105263158; Page_97147646424123=1; cdb_RW_ID_1676757366=1; Page_67187525232155=1; cdb_RW_ID_1530449813=2; Page_Y_54459150776295=-119.39144736842105; Page_54459150776295=1; cdb_back[id]=1; cdb_RW_ID_1410747206=1; Page_49799491545012=1; cdb_RW_ID_1455699102=40; cdb_back[doctype]=1; cdb_back[action]=data; cdb_back[p_code]=79339633255987; cdb_back[pcode]=79339633255987; cdb_back[ajax]=1; cdb_back[tm]=6681; Page_Y_79339633255987=-119.39144736842105; Page_79339633255987=1; cdb_back[t]=0; cdb_RW_ID_1676767595=2; Page_36016424242838=1; cdb_RW_ID_1633962042=3; Page_69339200527867=1; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A0tXizNXiFNXi2jMW2qvV2q3Q1TMU2LN9or3c3gJACuvi2i3R1jEW1THQ1qkV3iXiFotZBK363m0W1WpkHjFm2O3Q0Ls5HTMU1Oxi2qFh1OsQ1TvV0LsX3iXiDutdHmtSoWlk3jfi0qPU1qk%210T0Q3gU%3D; cdb_back[show_view_type]=1; cdb_back[show]=1; cdb_back[order_num]=1; cdb_back[folder_page]=0; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[member_id]=a9a7e1sN0qHKJRIQKLs1B8X7T%2FFQ83XNo7gQsSxo8i5ll8DXrqqJCk%2F6V5SU9Owq7WA; cdb_back[pm_id]=1486396; cdb_back[friend_id]=0; cdb_back[wxcode]=81763; cdb_back[txt_amount]=53; cdb_back[checkcode]=81763; cdb_back[type]=score; cdb_pageType=2; cdb_back[pcid]=8243; cdb_back[p_price]=388; cdb_back[p_default_points]=2; cdb_back[doccode]=1676981980; cdb_back[title]=2022%E5%B9%B4%E9%9D%92%E6%B5%B7%E6%B5%B7%E4%B8%9C%E4%B8%AD%E8%80%83%E6%95%B0%E5%AD%A6%E7%9C%9F%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88; cdb_back[intro]=2022%E5%B9%B4%E9%9D%92%E6%B5%B7%E6%B5%B7%E4%B8%9C%E4%B8%AD%E8%80%83%E6%95%B0%E5%AD%A6%E7%9C%9F%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88; cdb_back[p_pagecount]=8; cdb_back[doc_more_id]=1676979002%2C1676978974%2C1676978942%2C1676978926%2C1676978864%2C1676978785%2C1676978718%2C1676978188%2C1676977873%2C1676977843%2C1676977763%2C1676977703%2C1676977650%2C1676977635%2C1676977505%2C1676977341%2C1676977227%2C1676976856%2C; cdb_back[pid]=67887525691221; cdb_RW_ID_1676983773=1; cdb_back[srlid]=4ec0tRQKPZLlZvHOWM8EQ953tCVmGCoWagA8fQMHdmX4RFaHrBiTq4t5+mWiuGjv0oM4b2NabTPsPKhNux9NOhUrnF18wM%2F3+rcz%2FHFuZJC%2F; cdb_back[p_name]=2018%E5%90%89%E6%9E%97%E4%B8%AD%E8%80%83%E6%95%B0%E5%AD%A6%E7%9C%9F%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88; cdb_back[rel_p_id]=1676983773; cdb_back[p_id]=1676983773; cdb_back[page]=1; Page_Y_67887525691221=-119.39144736842105; Page_67887525691221=1; cdb_back[len]=2; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d23985013e508cdcfa515f63502ca44a80ad92f350c3f26f274; doc88_lt=wx; cdb_back[module_type]=6; cdb_back[image_type]=3; cdb_change_message=1; cdb_msg_num=0; siftState=1; cdb_back[state]=all; cdb_back[menuIndex]=2; cdb_back[classify_id]=all; cdb_msg_time=1694654372; cdb_back[act]=ajax_doc_list; cdb_back[curpage]=2"
var EditCookie = "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; _qddab=3-gv2ozy.lib1y9mi; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; cdb_READED_PC_ID=%2C441442; cdb_back[at]=0; cdb_back[n]=6; cdb_back[book_id]=0; cdb_back[s]=rel; Page_99439638538500=1; cdb_RW_ID_1664816554=2; Page_79399224792334=1; cdb_RW_ID_1664816550=2; Page_91061225912004=1; cdb_back[folder_id]=0; cdb_back[show_index]=1; show_index=1; cdb_back[u]=1; cdb_back[login]=1; cdb_back[txtloginname]=15238369929; cdb_back[txtPassword]=abcdqq123456; cdb_back[inout]=all; cdb_back[withdraw_name]=%E5%91%A8%E5%9B%BD%E4%BC%9F; cdb_back[withdraw_sfz]=410928198704276311; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dall; cdb_RW_ID_1674504163=1; Page_58339246386920=1; cdb_RW_ID_1674504139=1; ShowSkinTip_1=1; BG_SAVE=1; showAnnotateTipIf=1; Page_85629417657928=1; cdb_RW_ID_1674504405=1; cdb_H5R=1; Page_61399254314413=1; Page_Y_60659831344576=-92.77302631578948; cdb_RW_ID_1675722349=3; Page_60659831344576=1; cdb_RW_ID_1643689091=4; Page_Y_30287581596467=503.7927631578948; Page_30287581596467=12; cdb_back[captchaCode]=1; cdb_login_if=1; cdb_uid=104598337; c_login_name=woyoceo; cdb_logined=1; cdb_RW_ID_1676081712=1; Page_Y_97847646857473=-286.1763157894737; Page_97847646857473=2; cdb_RW_ID_1676297250=1; Page_89929414081065=1; cdb_RW_ID_1676298311=1; Page_31161262839711=1; cdb_RW_ID_1659055033=8; Page_29216483788700=1; cdb_RW_ID_1676298283=1; Page_89929414083032=1; cdb_RW_ID_1676563828=1; cdb_RW_ID_572746136=1; Page_0408320285715=1; cdb_back[m]=104598337; cdb_RW_ID_578093873=1; Page_1857576840670=1; Page_Y_31861262027989=79.46381578947368; Page_31861262027989=1; cdb_RW_ID_1676563821=1; cdb_back[uid]=b99ce806c0b55b3bdccae7bc14f8ca3e; Page_Y_69459838185249=-119.39144736842105; Page_69459838185249=1; cdb_RW_ID_1317655069=45; Page_Y_00799895233126=147.1888157894737; Page_00799895233126=1; cdb_back[sharetodoc]=1; cdb_back[download]=2; cdb_back[p_doc_format]=PDF; cdb_RW_ID_1633962014=2; Page_87387511650478=1; cdb_RW_ID_1454752637=2; Page_97316181289402=1; Page_Y_36016424242838=-119.39144736842105; cdb_RW_ID_1676757315=1; Page_Y_69559838313591=-119.39144736842105; Page_69559838313591=1; cdb_RW_ID_1676757334=1; Page_36016424282001=1; cdb_RW_ID_1653590184=1; Page_24887531364798=1; cdb_RW_ID_1676757352=1; Page_Y_97147646424123=63.80592105263158; Page_97147646424123=1; cdb_RW_ID_1676757366=1; Page_67187525232155=1; cdb_RW_ID_1530449813=2; Page_Y_54459150776295=-119.39144736842105; Page_54459150776295=1; cdb_back[id]=1; cdb_RW_ID_1410747206=1; Page_49799491545012=1; cdb_RW_ID_1455699102=40; cdb_back[doctype]=1; cdb_back[action]=data; cdb_back[p_code]=79339633255987; cdb_back[pcode]=79339633255987; cdb_back[ajax]=1; cdb_back[tm]=6681; Page_Y_79339633255987=-119.39144736842105; Page_79339633255987=1; cdb_back[t]=0; cdb_RW_ID_1676767595=2; Page_36016424242838=1; cdb_RW_ID_1633962042=3; Page_69339200527867=1; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A0tXizNXiFNXi2jMW2qvV2q3Q1TMU2LN9or3c3gJACuvi2i3R1jEW1THQ1qkV3iXiFotZBK363m0W1WpkHjFm2O3Q0Ls5HTMU1Oxi2qFh1OsQ1TvV0LsX3iXiDutdHmtSoWlk3jfi0qPU1qk%210T0Q3gU%3D; cdb_back[show_view_type]=1; cdb_back[show]=1; cdb_back[order_num]=1; cdb_back[folder_page]=0; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[member_id]=a9a7e1sN0qHKJRIQKLs1B8X7T%2FFQ83XNo7gQsSxo8i5ll8DXrqqJCk%2F6V5SU9Owq7WA; cdb_back[pm_id]=1486396; cdb_back[friend_id]=0; cdb_back[wxcode]=81763; cdb_back[txt_amount]=53; cdb_back[checkcode]=81763; cdb_back[type]=score; cdb_pageType=2; cdb_back[pcid]=8243; cdb_back[p_price]=388; cdb_back[p_default_points]=2; cdb_back[doccode]=1676981980; cdb_back[title]=2022%E5%B9%B4%E9%9D%92%E6%B5%B7%E6%B5%B7%E4%B8%9C%E4%B8%AD%E8%80%83%E6%95%B0%E5%AD%A6%E7%9C%9F%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88; cdb_back[intro]=2022%E5%B9%B4%E9%9D%92%E6%B5%B7%E6%B5%B7%E4%B8%9C%E4%B8%AD%E8%80%83%E6%95%B0%E5%AD%A6%E7%9C%9F%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88; cdb_back[p_pagecount]=8; cdb_back[doc_more_id]=1676979002%2C1676978974%2C1676978942%2C1676978926%2C1676978864%2C1676978785%2C1676978718%2C1676978188%2C1676977873%2C1676977843%2C1676977763%2C1676977703%2C1676977650%2C1676977635%2C1676977505%2C1676977341%2C1676977227%2C1676976856%2C; cdb_back[pid]=67887525691221; cdb_RW_ID_1676983773=1; cdb_back[srlid]=4ec0tRQKPZLlZvHOWM8EQ953tCVmGCoWagA8fQMHdmX4RFaHrBiTq4t5+mWiuGjv0oM4b2NabTPsPKhNux9NOhUrnF18wM%2F3+rcz%2FHFuZJC%2F; cdb_back[p_name]=2018%E5%90%89%E6%9E%97%E4%B8%AD%E8%80%83%E6%95%B0%E5%AD%A6%E7%9C%9F%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88; cdb_back[rel_p_id]=1676983773; cdb_back[page]=1; Page_Y_67887525691221=-119.39144736842105; Page_67887525691221=1; cdb_back[len]=2; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d23985013e508cdcfa515f63502ca44a80ad92f350c3f26f274; doc88_lt=wx; cdb_back[module_type]=6; cdb_back[image_type]=3; cdb_change_message=1; cdb_msg_num=0; siftState=1; cdb_back[state]=all; cdb_back[menuIndex]=2; cdb_back[classify_id]=all; cdb_msg_time=1694654372; cdb_back[curpage]=2; cdb_back[act]=getDocInfo; cdb_back[p_id]=1676983246"

var EditCount = 1

// ychEduSpider 编辑道客巴巴文档
// @Title 编辑道客巴巴文档
// @Description https://www.doc88.com/，编辑道客巴巴文档
func main() {
	curPage := 351
	for {
		pageListUrl := fmt.Sprintf("https://www.doc88.com/uc/doc_manager.php?act=ajax_doc_list&curpage=%d", curPage)
		fmt.Println(pageListUrl)
		queryEditDoc88ListFormData := QueryEditDoc88ListFormData{
			MenuIndex:  2,
			ClassifyId: "all",
			FolderId:   0,
			Sort:       1,
			Keyword:    "",
			ShowIndex:  1,
		}
		pageListDoc, err := QueryEditDoc88List(pageListUrl, queryEditDoc88ListFormData)
		if err != nil {
			fmt.Println(err)
			break
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

			CatalogNameNode := htmlquery.FindOne(liNode, `./div[@class="bookdoc"]/div[@class="posttime"]/span/span[@class="catelog_name"]`)
			CatalogName := htmlquery.InnerText(CatalogNameNode)
			if !strings.Contains(CatalogName, "标准规范") {
				fmt.Println("===========不是团体标准，跳过===========")
				continue
			}

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
			} else if filePageNum > 35 && filePageNum <= 40 {
				PPriceNew = "988"
			} else if filePageNum > 40 && filePageNum <= 45 {
				PPriceNew = "1088"
			} else if filePageNum > 45 && filePageNum <= 50 {
				PPriceNew = "1188"
			} else {
				PPriceNew = "1288"
			}

			// 新旧价格一样，则跳过
			fmt.Println(PPrice, PPriceNew)
			if PPrice == PPriceNew {
				continue
			}

			PId := htmlquery.SelectAttr(liNode, "id")
			PId = PId[5:]

			for i := 1; i <= 35; i++ {
				time.Sleep(time.Second)
				fmt.Println("===========获取", Title, "详情暂停35秒，倒计时", i, "秒===========")
			}

			detailUrl := "https://www.doc88.com/uc/usr_doc_manager.php?act=getDocInfo"
			detailDoc, err := QueryEditDoc88Detail(detailUrl, PId)
			if err != nil {
				fmt.Println(err)
				break
			}
			DocCodeNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[1]/div[@class="booksedit"]/table[@class="edit-table"]/input`)
			DocCode := htmlquery.SelectAttr(DocCodeNode, "value")

			PCidNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[1]/div[@class="booksedit"]/table[@class="edit-table"]/tbody/tr[3]/td[2]/div[@class="layers"]/input`)
			PCid := htmlquery.SelectAttr(PCidNode, "value")
			if PCid != "8370" {
				fmt.Println("===========不是团体标准，跳过===========")
				continue
			}

			PDocFormatNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[2]/div[@class="booksedit booksedit-bdr"]/table[@class="edit-table"]/tbody/tr[3]/td[2]/input[3]`)
			PDocFormat := htmlquery.SelectAttr(PDocFormatNode, "value")

			fmt.Println("===========开始修改", Title, "价格===========", EditCount)
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
			_, err = EditDoc88(editUrl, editDoc88FormData)
			if err != nil {
				fmt.Println(err)
				break
			}
			EditCount++
			if EditCount > 3 {
				EditCount = 1
				for i := 1; i <= 80; i++ {
					time.Sleep(time.Second)
					fmt.Println("===========更新数量超过3，暂停80秒，倒计时", i, "秒===========")
				}
			} else {
				for i := 1; i <= 45; i++ {
					time.Sleep(time.Second)
					fmt.Println("===========更新", Title, "成功，暂停45秒，倒计时", i, "秒===========")
				}
			}
		}
		curPage++
		for i := 1; i <= 20; i++ {
			time.Sleep(time.Second)
			fmt.Println("===========翻", curPage, "页，暂停20秒，倒计时", i, "秒===========")
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
	req.Header.Set("Cookie", ListCookie)
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
	req.Header.Set("Cookie", DetailCookie)
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
	req.Header.Set("Cookie", EditCookie)
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
