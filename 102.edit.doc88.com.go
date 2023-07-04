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

var ListCookie = "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; _qddab=3-gv2ozy.lib1y9mi; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; Page_Y_01473236654992=-122.80263157894737; Page_01473236654992=1; cdb_RW_ID_1638301172=1; Page_10259852509934=1; cdb_RW_ID_1653079317=1; Page_Y_27916480723062=-119.39144736842105; Page_27916480723062=1; cdb_RW_ID_1445399964=8; Page_69716118033341=1; cdb_RW_ID_1366298896=178; Page_23773022146642=1; cdb_RW_ID_1653081930=1; cdb_RW_ID_1452176401=1; Page_97339637942689=1; Page_Y_27916480756307=-82.98684210526316; Page_27916480756307=8; cdb_RW_ID_1653081818=1; cdb_RW_ID_1433712132=1; Page_60299488590980=1; Page_Y_64861207491919=-119.39144736842105; Page_64861207491919=1; show_index=1; cdb_back[inout]=all; Page_Y_30459815029852=-122.80263157894737; cdb_RW_ID_1653081628=1; Page_48739230819271=1; cdb_RW_ID_1653081672=2; Page_15029462539410=1; cdb_RW_ID_1653081638=2; Page_30459815029852=1; cdb_RW_ID_1653080853=1; Page_48347621858521=1; cdb_back[sharetodoc]=1; cdb_back[download]=2; cdb_back[p_doc_format]=PDF; cdb_back[p_price]=388; cdb_back[doc_id]=1653317081; cdb_back[p_default_points]=3; cdb_back[doccode]=1653317081; cdb_back[title]=%E7%AC%AC%E4%B8%89%E6%96%B9%E7%8E%AF%E4%BF%9D%E6%9C%8D%E5%8A%A1%E6%9C%BA%E6%9E%84%E6%9C%8D%E5%8A%A1%E8%A7%84%E8%8C%83%28DB1311-T+024-2022%29; cdb_back[intro]=%E7%AC%AC%E4%B8%89%E6%96%B9%E7%8E%AF%E4%BF%9D%E6%9C%8D%E5%8A%A1%E6%9C%BA%E6%9E%84%E6%9C%8D%E5%8A%A1%E8%A7%84%E8%8C%83%28DB1311-T+024-2022%29; cdb_back[p_pagecount]=13; Page_Y_15229462298588=-119.39144736842105; cdb_RW_ID_1653319089=1; Page_Y_64761207713493=-271.2730263157895; Page_64761207713493=6; cdb_RW_ID_1653319099=2; Page_15229462298588=1; cdb_RW_ID_1429305629=1; Page_Y_97987806143506=-138.1809210526316; Page_97987806143506=1; cdb_back[txtloginname]=15238369929; cdb_back[txtPassword]=abcdqq123456; cdb_back[captchaCode]=1; cdb_login_if=1; cdb_uid=104598337; c_login_name=woyoceo; cdb_logined=1; cdb_RW_ID_1518666299=1; cdb_back[action]=data; cdb_back[p_code]=20929693444088; Page_Y_20929693444088=-119.39144736842105; Page_20929693444088=1; cdb_RW_ID_1653318012=1; Page_Y_30559815592094=-119.39144736842105; Page_30559815592094=1; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[order_num]=1; cdb_back[folder_page]=0; ShowSkinTip_1=1; cdb_H5R=1; showAnnotateTipIf=1; Page_Y_15729462686450=-119.39144736842105; cdb_back[s]=rel; cdb_RW_ID_1653595533=1; Page_78973250545500=1; cdb_RW_ID_1434151083=4; Page_30716101686750=1; cdb_RW_ID_1435331049=1; Page_14461570771453=1; cdb_RW_ID_1653595163=1; cdb_RW_ID_1452176313=1; Page_70387830725171=1; Page_Y_15729462686942=2.7401315789473686; Page_15729462686942=1; cdb_RW_ID_35698275=1; cdb_back[pcode]=496983267053; cdb_back[ajax]=1; cdb_back[tm]=4868; cdb_back[uid]=60c577fa1c7172c3a0bc7253bf638a89; cdb_back[page]=1; Page_Y_496983267053=-138.1809210526316; Page_496983267053=1; cdb_RW_ID_1653595134=1; cdb_READED_PC_ID=%2C443440443443; Page_51499238363984=1; cdb_RW_ID_1449904571=1; Page_99029778857619=1; cdb_back[m]=104598337; cdb_back[len]=1; cdb_back[id]=1; cdb_RW_ID_1648574648=1; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A0tXizNXiFNXi2jMW2LnU0L3V0TE50jN9or3c3gJACuvi2i3R1jsT1qkV0TvV3iXiFotZBK363j3THuHVBqPSHuBj2q1l2qxi2LnQ0OMTHTBl0u3U0LFm3iXiDutdHmtSoWlk3jfi0qPU1qk%210T0Q3gU%3D; Page_57087589328589=1; cdb_change_message=1; cdb_msg_num=0; cdb_RW_ID_1458169830=10; Page_99216185643507=1; cdb_back[type]=score; cdb_RW_ID_1652297412=1; Page_Y_49516489932169=-119.39144736842105; Page_49516489932169=1; cdb_RW_ID_1653315303=1; Page_48039230093080=1; siftState=1; cdb_back[classify_id]=all; cdb_RW_ID_1653595345=2; Page_Y_78973250545095=-119.39144736842105; Page_78973250545095=1; cdb_RW_ID_1653595602=2; Page_15729462686450=1; cdb_back[member_id]=f397lfHZI49aK0A71mo3+N5EMJlgiRMzMCU2KaYdzCm1YI3YFfpg8ukXs+th6w7zNVE; cdb_back[show]=2; cdb_back[curpage]=5; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dmyshare; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d231511c82e759a92389de181255884092dd92f350c3f26f274; doc88_lt=wx; cdb_back[module_type]=7; cdb_back[t]=1; cdb_back[image_type]=2; cdb_back[state]=all; cdb_back[menuIndex]=2; cdb_back[pcid]=563; cdb_RW_ID_1653081057=1; cdb_back[n]=6; cdb_back[doctype]=1; cdb_back[book_id]=0; cdb_back[at]=0; Page_48347621857824=1; cdb_back[p]=2023%2F07%2F02%2F27916480757539.xml.gz; cdb_back[pid]=27916480757539; cdb_RW_ID_1653080892=1; cdb_back[p_name]=2023%E5%B9%B4%E5%9B%9B%E5%B7%9D%E6%88%90%E9%83%BD%E4%B8%AD%E8%80%83%E8%AF%AD%E6%96%87%E8%AF%95%E9%A2%98%28%E5%90%AB%E7%AD%94%E6%A1%88%29; cdb_back[rel_p_id]=1653080892; cdb_back[p_id]=1653080892; cdb_back[srlid]=8b24ILemgJftutNFw75m70ZeBTlrwseA+UNP3EnjgaOtajQQ7WrCtmPBMUjDW1i+1jazKc4lpP%2F%2FYnZ%2FWochrCOb3+T8Brtbv%2FZE+Zulndxo; Page_27916480757539=1; cdb_back[u]=1; cdb_msg_time=1688440007; cdb_back[act]=ajax_doc_folder"
var DetailCookie = "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; _qddab=3-gv2ozy.lib1y9mi; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; Page_Y_30459815029852=-122.80263157894737; cdb_RW_ID_1653081628=1; Page_48739230819271=1; cdb_RW_ID_1653081672=2; Page_15029462539410=1; cdb_RW_ID_1653081638=2; Page_30459815029852=1; cdb_RW_ID_1653080853=1; Page_48347621858521=1; cdb_back[sharetodoc]=1; cdb_back[download]=2; cdb_back[p_doc_format]=PDF; cdb_back[doc_id]=1653317081; Page_Y_15229462298588=-119.39144736842105; cdb_RW_ID_1653319089=1; Page_Y_64761207713493=-271.2730263157895; Page_64761207713493=6; cdb_RW_ID_1653319099=2; Page_15229462298588=1; cdb_RW_ID_1429305629=1; Page_Y_97987806143506=-138.1809210526316; Page_97987806143506=1; cdb_back[txtloginname]=15238369929; cdb_back[txtPassword]=abcdqq123456; cdb_back[captchaCode]=1; cdb_login_if=1; cdb_uid=104598337; c_login_name=woyoceo; cdb_logined=1; cdb_RW_ID_1518666299=1; cdb_back[action]=data; cdb_back[p_code]=20929693444088; Page_Y_20929693444088=-119.39144736842105; Page_20929693444088=1; cdb_RW_ID_1653318012=1; Page_Y_30559815592094=-119.39144736842105; Page_30559815592094=1; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[order_num]=1; cdb_back[folder_page]=0; ShowSkinTip_1=1; cdb_H5R=1; showAnnotateTipIf=1; Page_Y_15729462686450=-119.39144736842105; cdb_RW_ID_1434151083=4; Page_30716101686750=1; cdb_RW_ID_1435331049=1; Page_14461570771453=1; cdb_RW_ID_1653595163=1; cdb_RW_ID_1452176313=1; Page_70387830725171=1; Page_Y_15729462686942=2.7401315789473686; Page_15729462686942=1; cdb_RW_ID_35698275=1; Page_Y_496983267053=-138.1809210526316; Page_496983267053=1; cdb_RW_ID_1653595134=1; cdb_READED_PC_ID=%2C443440443443; Page_51499238363984=1; cdb_RW_ID_1449904571=1; Page_99029778857619=1; cdb_back[m]=104598337; cdb_RW_ID_1648574648=1; Page_57087589328589=1; cdb_change_message=1; cdb_msg_num=0; cdb_RW_ID_1458169830=10; Page_99216185643507=1; cdb_RW_ID_1652297412=1; Page_Y_49516489932169=-119.39144736842105; Page_49516489932169=1; cdb_RW_ID_1653315303=1; Page_48039230093080=1; siftState=1; cdb_back[classify_id]=all; cdb_RW_ID_1653595345=2; Page_Y_78973250545095=-119.39144736842105; Page_78973250545095=1; cdb_back[show]=2; cdb_RW_ID_1653081057=1; cdb_back[n]=6; cdb_back[book_id]=0; cdb_back[at]=0; Page_48347621857824=1; cdb_RW_ID_1653080892=1; Page_27916480757539=1; cdb_back[u]=1; cdb_back[folder_id]=0; cdb_back[show_index]=1; cdb_back[pcid]=8371; cdb_back[p_default_points]=2; cdb_RW_ID_1653595602=3; cdb_back[doctype]=10; Page_15729462686450=1; cdb_search_format=; cdb_back[ct]=0; cdb_back[h]=1; cdb_RW_ID_1651724888=1; Page_Y_40629469107333=221.16776315789474; Page_40629469107333=3; cdb_RW_ID_1653081607=1; cdb_RW_ID_1651369048=1; Page_27039239025861=2; cdb_RW_ID_1651369046=1; Page_Y_84559819586078=-238.7828947368421; Page_84559819586078=1; cdb_RW_ID_1649285080=1; Page_Y_49429478036535=91.99013157894737; Page_49429478036535=1; cdb_RW_ID_1646989512=1; Page_24661252393018=1; cdb_RW_ID_1646932079=1; Page_24661252378463=1; cdb_RW_ID_1653595533=2; Page_78973250545500=1; searchType=0; cdb_back[p]=%2F2023%2F06%2F22%2F49916486753927.xml.gz; cdb_back[curpage]=1; cdb_back[tabnum]=0; cdb_back[statusval]=%5C%27%5C%27; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dall; cdb_back[login]=1; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d231511c82e759a9238cf9b9586bf103e47d92f350c3f26f274; doc88_lt=wx; cdb_back[image_type]=3; cdb_back[module_type]=5; cdb_back[t]=1; show_index=1; cdb_pageType=2; cdb_RW_ID_1653594222=1; cdb_back[pid]=20159743167507; cdb_back[s]=rel; cdb_back[id]=1; cdb_back[p_name]=DB50%2FT+1269-2022+%E7%95%9C%E7%A6%BD%E7%B2%AA%E6%B1%A1%E8%B5%84%E6%BA%90%E5%8C%96%E5%88%A9%E7%94%A8+%E6%9C%AF%E8%AF%AD; cdb_VIEW_DOC_ID=%2C1427594304; cdb_RW_ID_1427594304=4; cdb_back[srlid]=a692XDmjgKzZsJnEfIzBVQFdDF1LlU97MxG8Duf5JYuMvjvXyk8N4+y3+++4S+DRn1SwEfG8MxmKYXztnSwKZesX1364hgDWgFq6GvGfZCQ; Page_20159743167507=1; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A0tXizNXiFNXi2jMW2LnU1LnV0qsX0jl9or3c3gJACuvi2i3R1jsT1qkU0j3S3iXiFotZBK363jsV0OsV2uHQBuHS0WsV1uBmBq1j0T0T1uHR0u0R0Lli3iXiDutdHmtSoWlk3jfi0qPU1qk%210T0Q3gU%3D; cdb_back[pcode]=15729462687000; cdb_back[ajax]=1; cdb_back[tm]=6730; cdb_back[uid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[page]=1; cdb_back[rel_p_id]=1653594222; cdb_back[member_id]=104598337; cdb_back[len]=3; Page_Y_15729462687000=-119.39144736842105; Page_15729462687000=1; cdb_back[inout]=all; cdb_back[type]=1; cdb_back[state]=all; cdb_back[menuIndex]=2; cdb_msg_time=1688454311; cdb_back[p_id]=1653595572; cdb_back[doccode]=1653595572; cdb_back[title]=%E8%87%AA%E7%84%B6%E7%81%BE%E5%AE%B3%E5%8F%97%E6%8D%9F%E7%AB%B9%E6%9E%97%E6%81%A2%E5%A4%8D%E6%8A%80%E6%9C%AF%E8%A7%84%E7%A8%8B%28DB33-T+2505-2022%29; cdb_back[intro]=%E8%87%AA%E7%84%B6%E7%81%BE%E5%AE%B3%E5%8F%97%E6%8D%9F%E7%AB%B9%E6%9E%97%E6%81%A2%E5%A4%8D%E6%8A%80%E6%9C%AF%E8%A7%84%E7%A8%8B%28DB33-T+2505-2022%29; cdb_back[p_price]=388; cdb_back[p_pagecount]=9; cdb_back[act]=save_info"
var EditCookie = "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; _qddab=3-gv2ozy.lib1y9mi; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; Page_Y_30459815029852=-122.80263157894737; cdb_RW_ID_1653081628=1; Page_48739230819271=1; cdb_RW_ID_1653081672=2; Page_15029462539410=1; cdb_RW_ID_1653081638=2; Page_30459815029852=1; cdb_RW_ID_1653080853=1; Page_48347621858521=1; cdb_back[sharetodoc]=1; cdb_back[download]=2; cdb_back[p_doc_format]=PDF; cdb_back[doc_id]=1653317081; Page_Y_15229462298588=-119.39144736842105; cdb_RW_ID_1653319089=1; Page_Y_64761207713493=-271.2730263157895; Page_64761207713493=6; cdb_RW_ID_1653319099=2; Page_15229462298588=1; cdb_RW_ID_1429305629=1; Page_Y_97987806143506=-138.1809210526316; Page_97987806143506=1; cdb_back[txtloginname]=15238369929; cdb_back[txtPassword]=abcdqq123456; cdb_back[captchaCode]=1; cdb_login_if=1; cdb_uid=104598337; c_login_name=woyoceo; cdb_logined=1; cdb_RW_ID_1518666299=1; cdb_back[action]=data; cdb_back[p_code]=20929693444088; Page_Y_20929693444088=-119.39144736842105; Page_20929693444088=1; cdb_RW_ID_1653318012=1; Page_Y_30559815592094=-119.39144736842105; Page_30559815592094=1; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[order_num]=1; cdb_back[folder_page]=0; ShowSkinTip_1=1; cdb_H5R=1; showAnnotateTipIf=1; Page_Y_15729462686450=-119.39144736842105; cdb_RW_ID_1434151083=4; Page_30716101686750=1; cdb_RW_ID_1435331049=1; Page_14461570771453=1; cdb_RW_ID_1653595163=1; cdb_RW_ID_1452176313=1; Page_70387830725171=1; Page_Y_15729462686942=2.7401315789473686; Page_15729462686942=1; cdb_RW_ID_35698275=1; Page_Y_496983267053=-138.1809210526316; Page_496983267053=1; cdb_RW_ID_1653595134=1; cdb_READED_PC_ID=%2C443440443443; Page_51499238363984=1; cdb_RW_ID_1449904571=1; Page_99029778857619=1; cdb_back[m]=104598337; cdb_RW_ID_1648574648=1; Page_57087589328589=1; cdb_change_message=1; cdb_msg_num=0; cdb_RW_ID_1458169830=10; Page_99216185643507=1; cdb_RW_ID_1652297412=1; Page_Y_49516489932169=-119.39144736842105; Page_49516489932169=1; cdb_RW_ID_1653315303=1; Page_48039230093080=1; siftState=1; cdb_back[classify_id]=all; cdb_RW_ID_1653595345=2; Page_Y_78973250545095=-119.39144736842105; Page_78973250545095=1; cdb_back[show]=2; cdb_RW_ID_1653081057=1; cdb_back[n]=6; cdb_back[book_id]=0; cdb_back[at]=0; Page_48347621857824=1; cdb_RW_ID_1653080892=1; Page_27916480757539=1; cdb_back[u]=1; cdb_back[folder_id]=0; cdb_back[show_index]=1; cdb_back[doccode]=1653595602; cdb_back[title]=%E8%87%AA%E7%84%B6%E7%94%9F%E6%80%81%E7%A9%BA%E9%97%B4%E5%88%86%E7%B1%BB%E6%8C%87%E5%8D%97%28DB61-T+1604-2022%29; cdb_back[intro]=%E8%87%AA%E7%84%B6%E7%94%9F%E6%80%81%E7%A9%BA%E9%97%B4%E5%88%86%E7%B1%BB%E6%8C%87%E5%8D%97%28DB61-T+1604-2022%29; cdb_back[pcid]=8371; cdb_back[p_price]=288; cdb_back[p_default_points]=2; cdb_back[p_pagecount]=7; cdb_RW_ID_1653595602=3; cdb_back[doctype]=10; Page_15729462686450=1; cdb_search_format=; cdb_back[ct]=0; cdb_back[h]=1; cdb_RW_ID_1651724888=1; Page_Y_40629469107333=221.16776315789474; Page_40629469107333=3; cdb_RW_ID_1653081607=1; cdb_RW_ID_1651369048=1; Page_27039239025861=2; cdb_RW_ID_1651369046=1; Page_Y_84559819586078=-238.7828947368421; Page_84559819586078=1; cdb_RW_ID_1649285080=1; Page_Y_49429478036535=91.99013157894737; Page_49429478036535=1; cdb_RW_ID_1646989512=1; Page_24661252393018=1; cdb_RW_ID_1646932079=1; Page_24661252378463=1; cdb_RW_ID_1653595533=2; Page_78973250545500=1; searchType=0; cdb_back[p]=%2F2023%2F06%2F22%2F49916486753927.xml.gz; cdb_back[curpage]=1; cdb_back[tabnum]=0; cdb_back[statusval]=%5C%27%5C%27; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dall; cdb_back[login]=1; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d231511c82e759a9238cf9b9586bf103e47d92f350c3f26f274; doc88_lt=wx; cdb_back[image_type]=3; cdb_back[module_type]=5; cdb_back[t]=1; show_index=1; cdb_pageType=2; cdb_RW_ID_1653594222=1; cdb_back[pid]=20159743167507; cdb_back[s]=rel; cdb_back[id]=1; cdb_back[p_name]=DB50%2FT+1269-2022+%E7%95%9C%E7%A6%BD%E7%B2%AA%E6%B1%A1%E8%B5%84%E6%BA%90%E5%8C%96%E5%88%A9%E7%94%A8+%E6%9C%AF%E8%AF%AD; cdb_VIEW_DOC_ID=%2C1427594304; cdb_RW_ID_1427594304=4; cdb_back[srlid]=a692XDmjgKzZsJnEfIzBVQFdDF1LlU97MxG8Duf5JYuMvjvXyk8N4+y3+++4S+DRn1SwEfG8MxmKYXztnSwKZesX1364hgDWgFq6GvGfZCQ; Page_20159743167507=1; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A0tXizNXiFNXi2jMW2LnU1LnV0qsX0jl9or3c3gJACuvi2i3R1jsT1qkU0j3S3iXiFotZBK363jsV0OsV2uHQBuHS0WsV1uBmBq1j0T0T1uHR0u0R0Lli3iXiDutdHmtSoWlk3jfi0qPU1qk%210T0Q3gU%3D; cdb_back[pcode]=15729462687000; cdb_back[ajax]=1; cdb_back[tm]=6730; cdb_back[uid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[page]=1; cdb_back[rel_p_id]=1653594222; cdb_back[member_id]=104598337; cdb_back[len]=3; Page_Y_15729462687000=-119.39144736842105; Page_15729462687000=1; cdb_back[inout]=all; cdb_back[type]=1; cdb_back[state]=all; cdb_back[menuIndex]=2; cdb_msg_time=1688454311; cdb_back[act]=getDocInfo; cdb_back[p_id]=1653595572"

var EditCount = 1

// ychEduSpider 编辑道客巴巴文档
// @Title 编辑道客巴巴文档
// @Description https://www.doc88.com/，编辑道客巴巴文档
func main() {
	curPage := 220
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

			PPageCountNode := htmlquery.FindOne(liNode, `./div[@class="bookimg"]/em`)
			PPageCount := htmlquery.InnerText(PPageCountNode)
			PPageCount = PPageCount[2:]

			PPriceNode := htmlquery.FindOne(liNode, `./div[@class="bookdoc"]/ul[@class="position"]/li[6]/span[@class="jifentip"]/strong[@class="red"]`)
			PPrice := htmlquery.InnerText(PPriceNode)

			filePageNum, _ := strconv.Atoi(PPageCount)
			PPriceNew := ""
			// 根据页数设置价格
			if filePageNum > 0 && filePageNum <= 8 {
				PPriceNew = "288"
			} else if filePageNum > 8 && filePageNum <= 18 {
				PPriceNew = "388"
			} else if filePageNum > 18 && filePageNum <= 28 {
				PPriceNew = "488"
			} else if filePageNum > 28 && filePageNum <= 38 {
				PPriceNew = "588"
			} else if filePageNum > 38 && filePageNum <= 48 {
				PPriceNew = "688"
			} else if filePageNum > 48 && filePageNum <= 58 {
				PPriceNew = "788"
			} else {
				PPriceNew = "888"
			}

			// 新旧价格一样，则跳过
			fmt.Println(PPrice, PPriceNew)
			if PPrice == PPriceNew {
				continue
			}
			fmt.Println("===========开始修改价格=============", EditCount)

			PId := htmlquery.SelectAttr(liNode, "id")
			PId = PId[5:]

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

			PDocFormatNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[2]/div[@class="booksedit booksedit-bdr"]/table[@class="edit-table"]/tbody/tr[3]/td[2]/input[3]`)
			PDocFormat := htmlquery.SelectAttr(PDocFormatNode, "value")

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
			editDoc88ResponseData, err := EditDoc88(editUrl, editDoc88FormData)
			if err != nil {
				fmt.Println(err)
				break
			}
			EditCount++
			fmt.Println(editDoc88ResponseData)
			if EditCount > 3 {
				EditCount = 1
				fmt.Println("==========更新数量超过3，暂停120秒==========")
				time.Sleep(time.Second * 120)
			} else {
				fmt.Println("==========更新成功，暂停25秒==========")
				time.Sleep(time.Second * 25)
			}
		}
		curPage++
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
