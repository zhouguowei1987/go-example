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
	"rsc.io/pdf"
	"strconv"
	"strings"
	"time"
)

const Doc88Cookie = "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; _qddab=3-gv2ozy.lib1y9mi; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; cdb_back[at]=0; cdb_back[n]=6; cdb_back[book_id]=0; Page_47316412869085=1; cdb_back[zip_pid]=1450463734; cdb_back[zip_pcode]=77047028061410; cdb_zip_parent_id_1450463734=4816193; current_cata_1450463734=4816193; cdb_RW_ID_385101967=33; cdb_back[pcode]=8912236959841; cdb_back[ajax]=1; cdb_back[tm]=4703; Page_Y_8912236959841=-138.1809210526316; Page_8912236959841=1; cdb_back[login]=1; cdb_back[txtloginname]=15238369929; cdb_back[txtPassword]=abcdqq123456; cdb_back[captchaCode]=1; cdb_login_if=1; cdb_uid=104598337; c_login_name=woyoceo; cdb_logined=1; cdb_pageType=2; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d232149056dc02025e76e85bd59b74cffaed92f350c3f26f274; cdb_back[image_type]=4; cdb_back[module_type]=7; Page_Y_47416414425136=-125.6546052631579; Page_47416414425136=1; cdb_back[show]=1; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_RW_ID_1649288589=1; cdb_RW_ID_1647003321=1; Page_Y_47216412770096=-138.1809210526316; Page_47216412770096=1; cdb_back[page]=1; Page_Y_29239265711315=-119.39144736842105; Page_29239265711315=1; cdb_back[member_id]=997bm2CAmkFFvqCzGqVNZZMc2cN7vRnVdPNYkH6nFaaqtkITAjYMiz4hOGITbTtTCGw; cdb_back[folder_id]=0; show_index=1; cdb_back[show_index]=1; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dall; doc88_lt=wx; cdb_tokenid=c101xkdbBa6GGY33d%2FpVTnvZ%2FfkVZMKiYOccuAjIX2Xx67fviUgGTRBUlltmCsP%2F%2BLOBu62%2FM7GRuwGN6yDKdTtSSTPLA9nn2KIlLtr%2FHDfur61BNuQyO9ML%2BphmvgqgFA; cdb_back[withdraw_name]=%E5%91%A8%E5%9B%BD%E4%BC%9F; cdb_back[withdraw_sfz]=410928198704276311; cdb_back[sharetodoc]=1; cdb_back[download]=2; cdb_RW_ID_1649701335=1; Page_Y_23773294783005=-119.39144736842105; Page_23773294783005=1; cdb_RW_ID_983733403=5; Page_1836397677547=1; cdb_back[pm_id]=1486396; cdb_back[friend_id]=0; cdb_back[inout]=all; cdb_RW_ID_1649845304=1; cdb_READED_PC_ID=%2C443443; Page_Y_46516413518071=73.63815789473684; Page_46516413518071=3; cdb_RW_ID_1649845634=1; Page_Y_67547609502610=-69.32236842105263; Page_67547609502610=3; cdb_RW_ID_1649845787=1; Page_Y_46516413518252=-119.39144736842105; Page_46516413518252=1; cdb_RW_ID_1650277378=1; Page_Y_29699231055857=-119.39144736842105; Page_29699231055857=1; cdb_RW_ID_1650277411=1; Page_49829465011799=1; cdb_RW_ID_1562360842=1; Page_Y_54759184580274=-138.1809210526316; Page_54759184580274=1; cdb_RW_ID_1650272671=1; Page_Y_29699231050259=-119.39144736842105; Page_29699231050259=1; cdb_RW_ID_1631830941=1; Page_Y_08047617518907=-138.1809210526316; Page_08047617518907=1; cdb_back[u]=1; cdb_back[t]=1; cdb_back[e_class]=upload; cdb_RW_ID_1650273287=1; Page_Y_29539238740714=108.95592105263157; Page_29539238740714=5; cdb_RW_ID_1651034378=1; Page_Y_27739239806041=-119.39144736842105; Page_27739239806041=1; cdb_back[wxcode]=69414; cdb_back[txt_amount]=52; cdb_back[checkcode]=69414; cdb_back[type]=score; cdb_RW_ID_1651034629=1; Page_Y_63347627810639=215.68750000000003; Page_63347627810639=1; cdb_RW_ID_1651368869=1; Page_Y_84559819582286=-238.7828947368421; Page_84559819582286=1; cdb_RW_ID_1651370529=1; Page_50187537124306=1; cdb_RW_ID_1243751584=2; Page_21073190753569=1; cdb_RW_ID_1243751582=2; Page_49016910286859=1; cdb_RW_ID_1651370506=1; Page_49016486027874=1; cdb_RW_ID_1468275443=5; Page_18347065342001=1; ShowSkinTip_1=1; cdb_H5R=1; showAnnotateTipIf=1; cdb_RW_ID_1450463734=9; Page_Y_77047028061410=-268.10526315789474; Page_77047028061410=1; cdb_back[size]=5; cdb_RW_ID_1651719931=1; cdb_back[uid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[m]=104598337; Page_Y_63247627479917=-119.39144736842105; Page_63247627479917=1; cdb_RW_ID_1651719964=1; cdb_back[s]=rel; cdb_RW_ID_1279647404=17; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A0lXizNXiFNXi2jMW2LEW1TH51jMW1q19or3c3gJACuvi2i3R1jsR1TM52qHU3iXiFotZBK363mHQHqPW1jphBLsV0jP5BqEV1qkSHW0Q1qFh0Tlk1m1i3iXiDutdHmtSoWlk3jfi0qPU1qk%210T0Q3gU%3D; Page_Y_14861863256545=-138.1809210526316; Page_14861863256545=1; Page_Y_84159819396687=-238.7828947368421; Page_84159819396687=1; cdb_change_message=1; cdb_msg_num=0; cdb_RW_ID_1651720213=1; cdb_back[doctype]=1; Page_Y_21573253718130=361.69736842105266; cdb_back[len]=1; Page_21573253718130=5; cdb_RW_ID_1651721004=1; Page_Y_49816486296771=0.39144736842105265; Page_49816486296771=1; cdb_back[pcid]=8371; cdb_back[p_price]=368; cdb_back[p_doc_format]=PDF; cdb_back[doc_id]=1651725607; cdb_back[p_default_points]=3; cdb_back[doccode]=1651725964; cdb_back[title]=%E7%BB%BF%E8%89%B2%E5%B7%A5%E4%B8%9A%E5%9B%AD%E5%8C%BA%E8%AF%84%E4%BB%B7%E8%A7%84%E8%8C%83%28DB21-T+3662%E2%80%942022%29; cdb_back[intro]=%E7%BB%BF%E8%89%B2%E5%B7%A5%E4%B8%9A%E5%9B%AD%E5%8C%BA%E8%AF%84%E4%BB%B7%E8%A7%84%E8%8C%83%28DB21-T+3662%E2%80%942022%29; cdb_back[p_pagecount]=26; cdb_back[doc_more_id]=1651719708%2C1651719678%2C1651719665%2C1651719658%2C1651719639%2C1651719625%2C1651719569%2C; cdb_back[id]=1651369055; cdb_back[curpage]=1; cdb_back[pid]=40629469101455; cdb_RW_ID_1651727600=1; cdb_back[srlid]=49d96fG4hHeVpzoBgbZNvd1ozrqSE0e29GyljTan+eRLUKylHAfZPLeFNpu3cDHtP4+5TVRGmbuYbCWqhcCRRavpaNQkKZq1nGwAvwLg3%2FIs; cdb_back[p_name]=%E5%9F%8E%E5%B8%82%E8%BD%A8%E9%81%93%E4%BA%A4%E9%80%9A%E7%94%B5%E5%AE%A2%E8%BD%A6%E9%A2%84%E9%AA%8C%E6%94%B6%E8%A7%84%E8%8C%83%28DB45-T+2619-2022%29; cdb_back[rel_p_id]=1450179711; cdb_back[p_id]=1651727600; Page_40629469101455=1; siftState=1; cdb_back[classify_id]=all; cdb_back[state]=all; cdb_back[menuIndex]=2; cdb_msg_time=1687708969; cdb_back[act]=upload"

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
	req.Header.Set("Cookie", Doc88Cookie)
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
func editFile(docCode string, title string, intro string, pCid int, price int, pDocFormat string) (editResponseData EditResponseData, err error) {
	client := &http.Client{}
	editResponseData = EditResponseData{}
	postData := url.Values{}
	postData.Add("doccode", docCode)
	postData.Add("title", title)
	postData.Add("intro", intro)
	postData.Add("pcid", strconv.Itoa(pCid))
	postData.Add("keyword", "")
	postData.Add("sharetodoc", "1")
	postData.Add("download", "2")
	postData.Add("p_price", strconv.Itoa(price))
	postData.Add("p_default_points", "1")
	postData.Add("p_pagecount", "")
	postData.Add("p_doc_format", pDocFormat)
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
	req.Header.Set("Cookie", Doc88Cookie)
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

type UploadChildDir struct {
	dirName string
	pCid    int
	Price   int
}

func main() {
	var uploadChildDirArr = []UploadChildDir{
		{
			dirName: "finish-www.ttbz.org.cn",
			pCid:    8370,
			Price:   468,
		},
		{
			dirName: "finish-dbba.sacinfo.org.cn",
			pCid:    8371,
			Price:   368,
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/专升本考试",
			pCid:    55789,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/初中一年级",
			pCid:    8155,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/中考",
			pCid:    8243,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/小升初",
			pCid:    8242,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/成人高考",
			pCid:    8244,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/考研",
			pCid:    8245,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/自考",
			pCid:    55789,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/高中会考",
			pCid:    8244,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/高考",
			pCid:    8244,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/初中一年级",
			pCid:    8155,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/中考",
			pCid:    8243,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/小升初",
			pCid:    8242,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/成人高考",
			pCid:    8244,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/考研",
			pCid:    8245,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/高中会考",
			pCid:    8244,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/高考",
			pCid:    8244,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/专升本考试",
			pCid:    55789,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/中考",
			pCid:    8243,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/小升初",
			pCid:    8242,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/成人高考", pCid: 8244,
			Price: 468,
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/考研",
			pCid:    8245,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/自考",
			pCid:    55789,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/高中会考", pCid: 8244,
			Price: 468,
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/高考", pCid: 8244,
			Price: 468,
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/专升本考试",
			pCid:    55789,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/中考",
			pCid:    8243,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/小升初",
			pCid:    8242,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/成人高考", pCid: 8244,
			Price: 468,
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/考研",
			pCid:    8245,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/自考",
			pCid:    55789,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/高中会考", pCid: 8244,
			Price: 468,
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/高考", pCid: 8244,
			Price: 468,
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/专升本考试",
			pCid:    55789,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/中考",
			pCid:    8243,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/小升初",
			pCid:    8242,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/成人高考", pCid: 8244,
			Price: 468,
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/考研",
			pCid:    8245,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/自考",
			pCid:    55789,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/高中会考",
			pCid:    8244,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/高考",
			pCid:    8244,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/专升本考试",
			pCid:    55789,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/中考",
			pCid:    8243,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/小升初",
			pCid:    8242,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/成人高考",
			pCid:    8244,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/考研",
			pCid:    8245,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/自考",
			pCid:    55789,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/高中会考",
			pCid:    8244,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/高考",
			pCid:    8244,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/专升本考试",
			pCid:    55789,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/中考",
			pCid:    8243,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/小升初",
			pCid:    8242,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/成人高考",
			pCid:    8244,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/考研",
			pCid:    8245,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/自考",
			pCid:    55789,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/高中会考",
			pCid:    8244,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/高考",
			pCid:    8244,
			Price:   468,
		},
		{
			dirName: "finish.topedu.ybep.com.cn/中考真题",
			pCid:    8243,
			Price:   468,
		},
		{
			dirName: "finish.topedu.ybep.com.cn/高考真题",
			pCid:    8244,
			Price:   468,
		},
		{
			dirName: "finish.www.shijuan1.com/中考试卷",
			pCid:    8243,
			Price:   468,
		},
		{
			dirName: "finish.www.shijuan1.com/高考试卷",
			pCid:    8244,
			Price:   468,
		},
		{
			dirName: "finish.www.tc260.org.cn",
			pCid:    8368,
			Price:   468,
		},
		{
			dirName: "docx.lvlin.baidu.com",
			pCid:    8131,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/二级建造师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/公务员考试",
			pCid:    8216,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/教师招聘",
			pCid:    8079,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/教师资格考试",
			pCid:    8079,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/事业单位招聘",
			pCid:    8217,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/一级建造师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/造价工程师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/证券从业资格考试",
			pCid:    8091,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/注册安全工程师考试",
			pCid:    8230,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/二级建造师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/公务员考试",
			pCid:    8216,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/教师招聘",
			pCid:    8079,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/教师资格考试",
			pCid:    8079,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/事业单位招聘",
			pCid:    8217,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/一级建造师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/造价工程师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/证券从业资格考试",
			pCid:    8091,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/注册安全工程师考试",
			pCid:    8230,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/二级建造师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/公务员考试",
			pCid:    8216,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/教师招聘",
			pCid:    8079,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/教师资格考试",
			pCid:    8079,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/事业单位招聘",
			pCid:    8217,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/消防工程师",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/一级建造师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/造价工程师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/证券从业资格考试",
			pCid:    8091,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/注册安全工程师考试",
			pCid:    8230,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/二级建造师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/公务员考试",
			pCid:    8216,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/教师招聘",
			pCid:    8079,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/教师资格考试",
			pCid:    8079,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/事业单位招聘",
			pCid:    8217,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/消防工程师",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/一级建造师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/造价工程师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/证券从业资格考试",
			pCid:    8091,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/注册安全工程师考试",
			pCid:    8230,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/二级建造师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/公务员考试",
			pCid:    8216,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/教师招聘",
			pCid:    8079,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/教师资格考试",
			pCid:    8079,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/事业单位招聘",
			pCid:    8217,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/消防工程师",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/一级建造师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/造价工程师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/证券从业资格考试",
			pCid:    8091,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/注册安全工程师考试",
			pCid:    8230,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/二级建造师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/公务员考试",
			pCid:    8216,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/教师招聘",
			pCid:    8079,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/教师资格考试",
			pCid:    8079,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/事业单位招聘",
			pCid:    8217,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/消防工程师",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/一级建造师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/造价工程师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/证券从业资格考试",
			pCid:    8091,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/注册安全工程师考试",
			pCid:    8230,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2023）.51test.net/二级建造师考试",
			pCid:    8210,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2023）.51test.net/公务员考试",
			pCid:    8216,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2023）.51test.net/教师招聘",
			pCid:    8079,
			Price:   468,
		},
		{
			dirName: "finish.tikuvip-certification（2023）.51test.net/教师资格考试",
			pCid:    8079,
			Price:   468,
		},
		{
			dirName: "docx.zzstep.com/小学/道德与法治",
			pCid:    8307,
			Price:   468,
		},
		{
			dirName: "docx.zzstep.com/小学/美术",
			pCid:    8309,
			Price:   468,
		},
		{
			dirName: "docx.zzstep.com/小学/数学",
			pCid:    8247,
			Price:   468,
		},
		{
			dirName: "docx.zzstep.com/小学/信息技术",
			pCid:    8308,
			Price:   468,
		},
		{
			dirName: "docx.zzstep.com/小学/音乐",
			pCid:    8309,
			Price:   468,
		},
		{
			dirName: "docx.zzstep.com/小学/英语",
			pCid:    8248,
			Price:   468,
		},
		{
			dirName: "docx.zzstep.com/小学/语文",
			pCid:    8246,
			Price:   468,
		},
	}
	rootPath := "../upload.doc88.com/"
	for _, childDir := range uploadChildDirArr {
		childDirPath := rootPath + childDir.dirName + "/"
		fmt.Println(childDirPath)
		files, err := ioutil.ReadDir(childDirPath)
		if err != nil {
			continue
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			fileName := file.Name()
			if fileName == ".DS_Store" {
				continue
			}
			fileExt := path.Ext(fileName)

			fmt.Println("==========开始上传==============")
			filePath := childDirPath + fileName
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
			fmt.Println("==========上传2秒后编辑文件所属类别和下载积分==============")
			time.Sleep(time.Second * 2)

			// 将上传过文件移动到"../final-upload.doc88.com/"
			finalDir := "../final-upload.doc88.com/" + childDir.dirName
			if _, err = os.Stat(finalDir); err != nil {
				if os.MkdirAll(finalDir, 0777) != nil {
					fmt.Println(err)
					break
				}
			}

			// 处理编辑文件时所需参数值
			docCode := uploadResponseData.DocCode
			title := strings.ReplaceAll(fileName, fileExt, "")
			title = strings.Replace(title, "\"", "“", 1)
			title = strings.Replace(title, "\"", "”", 1)
			intro := title
			pCid := childDir.pCid
			price := childDir.Price
			pDocFormat := ""
			switch fileExt {
			case ".pdf":
				pDocFormat = "PDF"
			case ".docx":
				pDocFormat = "DOCX"
			case ".doc":
				pDocFormat = "DOC"
			case ".ppt":
				pDocFormat = "PPT"
			case ".pptx":
				pDocFormat = "PPTX"
			}
			filePageNum := 0
			if fileExt == ".pdf" {
				// 获取PDF文件，获取总页数
				if pdfFile, err := pdf.Open(filePath); err == nil {
					filePageNum = pdfFile.NumPage()
				}
			}
			// 根据页数设置价格
			if filePageNum > 0 {
				if filePageNum > 0 && filePageNum <= 5 {
					price = 288
				} else if filePageNum > 5 && filePageNum <= 10 {
					price = 388
				} else if filePageNum > 10 && filePageNum <= 15 {
					price = 488
				} else if filePageNum > 15 && filePageNum <= 20 {
					price = 588
				} else if filePageNum > 20 && filePageNum <= 25 {
					price = 688
				} else if filePageNum > 25 && filePageNum <= 30 {
					price = 788
				} else if filePageNum > 30 && filePageNum <= 35 {
					price = 888
				} else if filePageNum > 35 && filePageNum <= 40 {
					price = 988
				} else if filePageNum > 40 && filePageNum <= 45 {
					price = 1088
				} else if filePageNum > 45 && filePageNum <= 50 {
					price = 1188
				} else {
					price = 1288
				}
			}

			// 将已上传的文件转移到指定文件夹
			fileFinal := finalDir + "/" + fileName
			err = os.Rename(filePath, fileFinal)
			if err != nil {
				fmt.Println(err)
				break
			}

			// 编辑文件所需分类和下载所需积分
			editResponseData, err := editFile(docCode, title, intro, pCid, price, pDocFormat)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println(editResponseData)

			fmt.Println("==========上传完成==============")
		}
	}
}
