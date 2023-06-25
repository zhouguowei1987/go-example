package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func getKey() (uploadKey string, err error) {
	client := &http.Client{} //初始化客户端
	getKeyUrl := "https://www.doc88.com/uc/index.php?act=getkey"
	req, err := http.NewRequest("GET", getKeyUrl, nil) //建立连接
	if err != nil {
		return uploadKey, err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; _qddab=3-gv2ozy.lib1y9mi; cdb_RW_ID_1471810323=1; Page_29129719395202=1; Page_Y_57187589391399=34.838815789473685; cdb_search_format=; cdb_back[ct]=0; cdb_back[h]=1; cdb_RW_ID_1632548210=1; Page_98673201596138=1; cdb_RW_ID_1543056311=1; Page_13147201826177=1; cdb_RW_ID_1069339956=65; Page_Y_8179126886632=-138.1809210526316; Page_8179126886632=1; cdb_RW_ID_1646641328=1; Page_Y_28239262269071=-138.1809210526316; Page_28239262269071=1; cdb_RW_ID_1646641322=2; Page_Y_47416414416099=-544.3631578947369; Page_47416414416099=1; cdb_RW_ID_1647380718=1; Page_Y_54687582194279=342.5164473684211; Page_54687582194279=10; cdb_back[txtAutologin]=1; cdb_RW_ID_1648583377=2; Page_29039261310044=1; cdb_back[pm_id]=1486396; cdb_back[friend_id]=0; cdb_RW_ID_1644494744=34; Page_Y_80716411131211=-119.39144736842105; Page_80716411131211=1; cdb_RW_ID_1648583588=1; Page_57187589391399=1; cdb_back[order_num]=1; cdb_back[folder_page]=0; cdb_RW_ID_1648583239=1; Page_21761259097873=1; cdb_back[e_class]=upload; cdb_RW_ID_1644494937=3; Page_62529477787821=1; cdb_RW_ID_1647738331=1; Page_Y_49529471123229=-119.39144736842105; Page_49529471123229=1; cdb_RW_ID_1458169749=9; Page_33447025769409=1; cdb_back[showIndex]=1; cdb_RW_ID_1644494741=1; cdb_back[p_code]=21847600090407; Page_21847600090407=1; cdb_RW_ID_1468275443=1; Page_Y_18347065342001=-119.39144736842105; Page_18347065342001=1; cdb_RW_ID_1446555443=62; Page_Y_94259778111775=90.42434210526316; Page_94259778111775=4; cdb_back[doc_more_id]=1649056298%2C1649056297%2C1649056295%2C1649056294%2C1649056291%2C1649056290%2C1649056289%2C1649056288%2C1649056287%2C1649056286%2C1649056284%2C1649056282%2C; cdb_RW_ID_1649056792=1; Page_29399246132560=1; cdb_RW_ID_1648385391=2; Page_Y_29739261013059=-138.1809210526316; Page_29739261013059=1; cdb_RW_ID_1647578794=1; ShowSkinTip_1=1; cdb_H5R=1; showAnnotateTipIf=1; Page_45829471613187=1; cdb_RW_ID_1649056801=1; Page_49629478564359=1; cdb_RW_ID_1647738329=1; Page_Y_29839264401075=-119.39144736842105; Page_29839264401075=1; cdb_RW_ID_1052791534=13; Page_Y_6621789236801=-276.3618421052632; Page_6621789236801=1; cdb_RW_ID_1052532622=18; Page_9962560620400=1; cdb_RW_ID_1471810433=1; Page_09439649198600=1; cdb_RW_ID_1072730880=5; Page_0763847408118=1; cdb_back[p]=%2F2023%2F06%2F11%2F23373296801778.xml.gz; cdb_back[type]=1; cdb_back[size]=404; cdb_RW_ID_1649285204=1; Page_57587586093048=1; cdb_RW_ID_1646480878=1; Page_54387585894929=1; cdb_back[show_view_type]=1; cdb_RW_ID_1649285182=1; Page_Y_46416413958659=-119.39144736842105; Page_46416413958659=1; cdb_back[s]=rel; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; cdb_RW_ID_1647512358=1; cdb_back[at]=0; cdb_back[n]=6; cdb_back[book_id]=0; Page_47316412869085=1; cdb_back[zip_pid]=1450463734; cdb_back[zip_pcode]=77047028061410; Page_Y_77047028061410=-203.68421052631578; cdb_zip_parent_id_1450463734=4816193; current_cata_1450463734=4816193; cdb_RW_ID_1450463734=7; Page_77047028061410=1; cdb_RW_ID_385101967=33; cdb_back[pcode]=8912236959841; cdb_back[ajax]=1; cdb_back[tm]=4703; Page_Y_8912236959841=-138.1809210526316; Page_8912236959841=1; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dmyprivate; cdb_back[login]=1; cdb_back[txtloginname]=15238369929; cdb_back[txtPassword]=abcdqq123456; cdb_back[captchaCode]=1; cdb_login_if=1; cdb_uid=104598337; c_login_name=woyoceo; cdb_logined=1; cdb_change_message=1; cdb_msg_num=0; cdb_pageType=2; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d232149056dc02025e76e85bd59b74cffaed92f350c3f26f274; doc88_lt=wx; cdb_back[image_type]=4; cdb_back[module_type]=7; cdb_back[curpage]=1; cdb_READED_PC_ID=%2C443; Page_Y_47416414425136=-125.6546052631579; Page_47416414425136=1; cdb_back[show]=1; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_RW_ID_1649288589=1; cdb_back[pid]=47216412770096; cdb_back[id]=4; cdb_RW_ID_1647003321=1; cdb_back[doctype]=10; cdb_back[p_name]=DB13T+5694-2023%E7%8C%AA%E5%A1%9E%E5%86%85%E5%8D%A1%E7%97%85%E6%AF%92%E7%97%85%E9%98%B2%E6%8E%A7%E6%8A%80%E6%9C%AF%E8%A7%84%E7%A8%8B; cdb_back[p_id]=1647003321; cdb_back[srlid]=b6edxyVv6TsGTwDHvUBb+jUJdmt1JKdYCo0oa+7HadEveIeXJGgxi2b3BqUc6scu1x4gMGAZQqpzsHX1P707V1OliEHjIExnfy%2FL9EpmCh8; cdb_VIEW_DOC_ID=%2C1647003321; cdb_back[m]=104598337; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A1NXizNXiFNXi2jMW2LH%211Tn51jPV0jh9or3c3gJACuvi2i3R1jv50jn%211qn53iXiFotZBK363jtk0WH%212uNj0qHX1LH%21BOH%211qM%212uxl2usT0LNjHW0Q3iXiDutdHmtSoWlk3jfi0qPU1qk%210T0Q3gU%3D; Page_Y_47216412770096=-138.1809210526316; Page_47216412770096=1; cdb_back[uid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[page]=1; cdb_back[rel_p_id]=1647001853; Page_Y_29239265711315=-119.39144736842105; cdb_back[len]=2; Page_29239265711315=1; cdb_back[member_id]=997bm2CAmkFFvqCzGqVNZZMc2cN7vRnVdPNYkH6nFaaqtkITAjYMiz4hOGITbTtTCGw; cdb_back[folder_id]=0; cdb_back[state]=myshare; siftState=1; show_index=1; cdb_back[u]=1; cdb_back[t]=1; cdb_back[menuIndex]=4; cdb_back[classify_id]=all; cdb_back[show_index]=1; cdb_msg_time=1686879171; cdb_back[act]=getkey")
	req.Header.Set("Host", "www.doc88.com")
	req.Header.Set("Referer", "https://www.doc88.com/uc/index.php?act=upload")
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
		return uploadKey, err
	}
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return uploadKey, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return uploadKey, err
	}
	uploadKey = string(s)
	return uploadKey, nil
}

func reverseString(s string) string {
	b := []byte(s)
	n := len(b)
	for i := 0; i < n/2; i++ {
		b[i], b[n-i-1] = b[n-i-1], b[i]
	}
	return string(b)

}

type UploadResponseData struct {
	Result  string `json:"result"`
	Message string `json:"message"`
	DocCode string `json:"doccode"`
	PId     string `json:"p_id"`
	PCode   string `json:"p_code"`
}

// 上传文件
func uploadFile(filePath string, uploadKey string, ck string) (uploadResponseData UploadResponseData, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	uploadResponseData = UploadResponseData{}

	file, err := os.Open(filePath)
	if err != nil {
		return uploadResponseData, err
	}
	defer file.Close()

	fileWriter, err := bodyWriter.CreateFormFile("upfile", filepath.Base(file.Name()))
	if err != nil {
		return uploadResponseData, err
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return uploadResponseData, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	err = bodyWriter.WriteField("act", "upload")
	if err != nil {
		return uploadResponseData, err
	}

	err = bodyWriter.WriteField("fileName", filepath.Base(file.Name()))
	if err != nil {
		return uploadResponseData, err
	}

	uploadUrl := fmt.Sprintf("https://upload.doc88.com/u.do?v=1&uploadkey=%s&ck=%s", uploadKey, ck)
	req, err := http.NewRequest("POST", uploadUrl, bodyBuf)
	if err != nil {
		return uploadResponseData, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Host", "upload.doc88.com")
	req.Header.Set("Origin", "https://www.doc88.com")
	req.Header.Set("Referer", "https://www.doc88.com/")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return uploadResponseData, err
	}

	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBytes, &uploadResponseData)
	if err != nil {
		return uploadResponseData, err
	}
	return uploadResponseData, nil
}

type EditResponseData struct {
	Result     string `json:"result"`
	EditTitle  string `json:"edit_title"`
	Class      string `json:"class"`
	UpdateInfo string `json:"updateinfo"`
	State      string `json:"state"`
	SaveFile   string `json:"savefile"`
	Other      string `json:"other"`
}

// 编辑文件所属分类和下载所需积分
func editFile(doccode string, title string, intro string, pcid int, price int) (editResponseData EditResponseData, err error) {
	client := &http.Client{}
	editResponseData = EditResponseData{}
	postData := url.Values{}
	postData.Add("doccode", doccode)
	postData.Add("title", title)
	postData.Add("intro", intro)
	postData.Add("pcid", strconv.Itoa(pcid))
	postData.Add("keyword", "")
	postData.Add("sharetodoc", "1")
	postData.Add("download", "2")
	postData.Add("p_price", strconv.Itoa(price))
	postData.Add("p_default_points", "1")
	postData.Add("p_pagecount", "")
	postData.Add("p_doc_format", "PDF")
	postData.Add("act", "save_info")
	postData.Add("group_list", "")
	postData.Add("group_free_list", "")
	requestUrl := "https://www.doc88.com/uc/index.php"
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return editResponseData, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; _qddab=3-gv2ozy.lib1y9mi; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; cdb_back[at]=0; cdb_back[n]=6; cdb_back[book_id]=0; Page_47316412869085=1; cdb_back[zip_pid]=1450463734; cdb_back[zip_pcode]=77047028061410; cdb_zip_parent_id_1450463734=4816193; current_cata_1450463734=4816193; cdb_RW_ID_385101967=33; cdb_back[pcode]=8912236959841; cdb_back[ajax]=1; cdb_back[tm]=4703; Page_Y_8912236959841=-138.1809210526316; Page_8912236959841=1; cdb_back[login]=1; cdb_back[txtloginname]=15238369929; cdb_back[txtPassword]=abcdqq123456; cdb_back[captchaCode]=1; cdb_login_if=1; cdb_uid=104598337; c_login_name=woyoceo; cdb_logined=1; cdb_pageType=2; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d232149056dc02025e76e85bd59b74cffaed92f350c3f26f274; cdb_back[image_type]=4; cdb_back[module_type]=7; Page_Y_47416414425136=-125.6546052631579; Page_47416414425136=1; cdb_back[show]=1; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_RW_ID_1649288589=1; cdb_RW_ID_1647003321=1; Page_Y_47216412770096=-138.1809210526316; Page_47216412770096=1; cdb_back[page]=1; Page_Y_29239265711315=-119.39144736842105; Page_29239265711315=1; cdb_back[member_id]=997bm2CAmkFFvqCzGqVNZZMc2cN7vRnVdPNYkH6nFaaqtkITAjYMiz4hOGITbTtTCGw; cdb_back[folder_id]=0; show_index=1; cdb_back[show_index]=1; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dall; doc88_lt=wx; cdb_tokenid=c101xkdbBa6GGY33d%2FpVTnvZ%2FfkVZMKiYOccuAjIX2Xx67fviUgGTRBUlltmCsP%2F%2BLOBu62%2FM7GRuwGN6yDKdTtSSTPLA9nn2KIlLtr%2FHDfur61BNuQyO9ML%2BphmvgqgFA; cdb_back[withdraw_name]=%E5%91%A8%E5%9B%BD%E4%BC%9F; cdb_back[withdraw_sfz]=410928198704276311; cdb_back[sharetodoc]=1; cdb_back[download]=2; cdb_back[p_price]=368; cdb_back[p_doc_format]=PDF; cdb_RW_ID_1649701335=1; Page_Y_23773294783005=-119.39144736842105; Page_23773294783005=1; cdb_RW_ID_983733403=5; Page_1836397677547=1; cdb_back[pm_id]=1486396; cdb_back[friend_id]=0; cdb_back[inout]=all; cdb_RW_ID_1649845304=1; cdb_READED_PC_ID=%2C443443; Page_Y_46516413518071=73.63815789473684; Page_46516413518071=3; cdb_back[pcid]=8371; cdb_RW_ID_1649845634=1; Page_Y_67547609502610=-69.32236842105263; Page_67547609502610=3; cdb_RW_ID_1649845787=1; Page_Y_46516413518252=-119.39144736842105; Page_46516413518252=1; cdb_RW_ID_1650277378=1; Page_Y_29699231055857=-119.39144736842105; Page_29699231055857=1; cdb_RW_ID_1650277411=1; Page_49829465011799=1; cdb_RW_ID_1562360842=1; Page_Y_54759184580274=-138.1809210526316; Page_54759184580274=1; cdb_RW_ID_1650272671=1; Page_Y_29699231050259=-119.39144736842105; Page_29699231050259=1; cdb_RW_ID_1631830941=1; Page_Y_08047617518907=-138.1809210526316; Page_08047617518907=1; cdb_back[u]=1; cdb_back[t]=1; cdb_back[e_class]=upload; cdb_RW_ID_1650273287=1; cdb_back[doctype]=1; Page_Y_29539238740714=108.95592105263157; Page_29539238740714=5; cdb_RW_ID_1651034378=1; Page_Y_27739239806041=-119.39144736842105; Page_27739239806041=1; cdb_back[wxcode]=69414; cdb_back[txt_amount]=52; cdb_back[checkcode]=69414; cdb_back[type]=score; cdb_RW_ID_1651034629=1; Page_Y_63347627810639=215.68750000000003; Page_63347627810639=1; cdb_back[doc_id]=1651368735; cdb_back[doc_more_id]=1651034595%2C1651034592%2C1651034589%2C1651034582%2C1651034578%2C1651034575%2C1651034572%2C1651034567%2C1651034565%2C1651034563%2C1651034559%2C1651034554%2C1651034552%2C1651034542%2C1651034540%2C; cdb_RW_ID_1651368869=1; Page_Y_84559819582286=-238.7828947368421; Page_84559819582286=1; cdb_back[curpage]=1; cdb_RW_ID_1651370529=1; Page_50187537124306=1; cdb_RW_ID_1243751584=2; Page_21073190753569=1; cdb_back[len]=1; cdb_RW_ID_1243751582=2; Page_49016910286859=1; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A0lXizNXiFNXi2jMW2LEW1qM52LEV1TN9or3c3gJACuvi2i3R1jsR0TEX1q353iXiFotZBK363jtj2Lll0uBl1Lxh1Opl0uvW0OHSHjkX0ThlHqkW0OvU3iXiDutdHmtSoWlk3jfi0qPU1qk%210T0Q3gU%3D; siftState=1; cdb_back[classify_id]=all; cdb_RW_ID_1651370506=1; Page_49016486027874=1; cdb_back[p_default_points]=1; cdb_RW_ID_1468275443=5; Page_18347065342001=1; cdb_msg_num=0; cdb_change_message=1; cdb_back[pid]=77047028061410; ShowSkinTip_1=1; cdb_H5R=1; showAnnotateTipIf=1; cdb_back[p_name]=%E9%81%93%E5%AE%A2%E5%B7%B4%E5%B7%B4%E6%89%B9%E9%87%8F%E4%B8%8A%E4%BC%A02.2.zip; cdb_back[rel_p_id]=1450463734; cdb_back[uid]=c6d5e05d8dc638aa41ac6258c990722f; cdb_back[m]=91184039; cdb_back[id]=99759710785351; cdb_RW_ID_1450463734=9; cdb_back[srlid]=5218PoDd+R6vOwaPFTayEnd8+wO2CoRciFo3O3IQePkbmKh7rPlm55VTh0hfkm9rqICxOgoxcJxQcaiH7auKiNNKB06DTnYyuMtYVF%2FASoU; Page_Y_77047028061410=-268.10526315789474; Page_77047028061410=1; cdb_back[doccode]=1651651349; cdb_back[title]=%E3%80%8A%E9%AB%98%E9%BB%8E%E8%B4%A1%E5%B1%B1%E7%8C%AA%E4%BF%9D%E7%A7%8D%E6%8A%80%E6%9C%AF%E8%A7%84%E7%A8%8B%E3%80%8B%28DB5333-T+28%E2%80%942022%29; cdb_back[intro]=%E3%80%8A%E9%AB%98%E9%BB%8E%E8%B4%A1%E5%B1%B1%E7%8C%AA%E4%BF%9D%E7%A7%8D%E6%8A%80%E6%9C%AF%E8%A7%84%E7%A8%8B%E3%80%8B%28DB5333-T+28%E2%80%942022%29; cdb_back[size]=5; cdb_back[p_id]=1651707794; cdb_back[state]=myprivate; cdb_back[menuIndex]=3; cdb_msg_time=1687673004; cdb_back[act]=getkey")
	req.Header.Set("Host", "www.doc88.com")
	req.Header.Set("Origin", "https://www.doc88.com")
	req.Header.Set("Referer", "https://www.doc88.com/uc/index.php?act=upload")
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
		return editResponseData, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return editResponseData, err
	}
	err = json.Unmarshal(respBytes, &editResponseData)
	if err != nil {
		return editResponseData, err
	}
	return editResponseData, nil
}
func main() {
	rootPath := "../finish-dbba.sacinfo.org.cn/"
	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		fileExt := path.Ext(fileName)
		if fileExt != ".pdf" {
			continue
		}
		fmt.Println("==========开始上传==============")
		filePath := rootPath + fileName
		fmt.Println(filePath)
		uploadKey, err := getKey()
		if err != nil {
			fmt.Println(err)
			continue
		}
		ck := reverseString(uploadKey[3:8])
		if err != nil {
			fmt.Println(err)
			continue
		}
		uploadResponseData, err := uploadFile(filePath, uploadKey, ck)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(uploadResponseData)

		if result, _ := strconv.Atoi(uploadResponseData.Result); result != 0 {
			fmt.Println(uploadResponseData.Message)
			break
		}
		fmt.Println("==========上传5秒后编辑文件所属类别和下载积分==============")
		time.Sleep(time.Second * 5)
		// 编辑文件所需分类和下载所需积分
		doccode := uploadResponseData.DocCode
		title := strings.ReplaceAll(fileName, fileExt, "")
		intro := title
		// 地方标准分类8371，团体标准分类8370
		pcid := 8370
		if strings.Contains(fileName, "DB") {
			pcid = 8371
		}

		price := 368
		editResponseData, err := editFile(doccode, title, intro, pcid, price)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(editResponseData)

		// 将上传过文件移动到"../final-dbba.sacinfo.org.cn/"
		finalDir := "../final-dbba.sacinfo.org.cn"
		if _, err = os.Stat(finalDir); err != nil {
			if os.MkdirAll(finalDir, 0777) != nil {
				fmt.Println(err)
				break
			}
		}
		fileFinal := finalDir + "/" + fileName
		err = os.Rename(filePath, fileFinal)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("==========上传完成==============")
	}
}
