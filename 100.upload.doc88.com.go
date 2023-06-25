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
	"os"
	"strconv"
)

func getKey() (uploadKey string, err error) {
	client := &http.Client{} //初始化客户端
	getKeyUrl := "https://www.doc88.com/uc/index.php?act=getkey"
	req, err := http.NewRequest("GET", getKeyUrl, nil) //建立连接
	if err != nil {
		return uploadKey, err
	}
	req.Header.Set("Accept", "*/*")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
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

func uploadFile(uploadKey string, ck string) (uploadResponseData UploadResponseData, err error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	uploadResponseData = UploadResponseData{}

	fileName := "‘红美人’柑橘设施栽培生产技术规程(DB32／T 4431-2022).pdf"
	filePath := "../finish-dbba.sacinfo.org.cn/" + fileName
	file, err := os.Open(filePath)
	if err != nil {
		return uploadResponseData, err
	}
	defer file.Close()

	err = writer.WriteField("act", "upload")
	if err != nil {
		return uploadResponseData, err
	}

	err = writer.WriteField("fileName", fileName)
	if err != nil {
		return uploadResponseData, err
	}

	formFile, err := writer.CreateFormFile("upfile", file.Name())
	if err != nil {
		return uploadResponseData, err
	}
	_, err = io.Copy(formFile, file)
	if err != nil {
		return uploadResponseData, err
	}

	uploadUrl := fmt.Sprintf("https://upload.doc88.com/u.do?v=1&uploadkey=%s&ck=%s", uploadKey, ck)
	req, err := http.NewRequest("POST", uploadUrl, body)
	if err != nil {
		return uploadResponseData, err
	}

	req.Header.Set("Accept", "*/*")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", writer.FormDataContentType())
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

func main() {
	uploadKey, err := getKey()
	if err != nil {
		return
	}
	fmt.Println("uploadKey：", uploadKey)
	ck := reverseString(uploadKey[3:8])
	if err != nil {
		return
	}
	fmt.Println("ck：", ck)
	uploadResponseData, err := uploadFile(uploadKey, ck)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", uploadResponseData)
	//fmt.Println(uploadResponseData)
}
