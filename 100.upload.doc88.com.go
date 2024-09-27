package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"math/rand"
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

// Doc88Cookie 15238369929
//const Doc88Cookie = "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_back[u]=1; cdb_back[t]=0; cdb_back[folder_id]=0; cdb_RW_ID_1911744937=1; cdb_back[uid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[m]=104598337; cdb_back[s]=rel; cdb_RW_ID_1632382023=4; cdb_RW_ID_1632382037=2; cdb_RW_ID_1911770319=2; cdb_back[doc_more_id]=1911863340%2C1911863304%2C1911863302%2C; cdb_RW_ID_1911863322=1; cdb_RW_ID_1911863300=1; cdb_RW_ID_1911852207=1; cdb_back[len]=1; cdb_RW_ID_1911863297=1; cdb_H5R=1; cdb_RW_ID_1911852206=1; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A0tXizNXiFNXi2jMQ0LvU0THR2qsS0Tx9or3c3gJACuvi2i3R2qMR2LHT0jkQ3iXiFotZBK363jsWBLxkBLE%211LnRBuMRBjJi2qhhHjhj0Lv%211TJl1qHQ3iXiDutdHmtSoWlk3jfi0qPU1qk%210T0Q3gU%3D; cdb_back[wxcode]=72100; cdb_back[txt_amount]=58; cdb_back[checkcode]=72100; cdb_back[p_price]=388; cdb_back[p_default_points]=2; cdb_back[id]=1911954892; cdb_RW_ID_1911954981=1; cdb_back[curpage]=1; cdb_back[classify_id]=all; cdb_back[type]=1; cdb_back[doccode]=1912168242; cdb_back[title]=%E5%BA%93%E5%B0%94%E5%8B%92%E9%A6%99%E6%A2%A8%E7%9C%81%E5%8A%9B%E5%AF%86%E6%A4%8D%E5%9B%AD%E6%97%A0%E4%BA%BA%E6%9C%BA%E8%BE%85%E5%8A%A9%E6%8E%88%E7%B2%89%E6%8A%80%E6%9C%AF%E8%A7%84%E7%A8%8B%28T-SHZSAQS+00241%E2%80%942024%29; cdb_back[intro]=%E5%BA%93%E5%B0%94%E5%8B%92%E9%A6%99%E6%A2%A8%E7%9C%81%E5%8A%9B%E5%AF%86%E6%A4%8D%E5%9B%AD%E6%97%A0%E4%BA%BA%E6%9C%BA%E8%BE%85%E5%8A%A9%E6%8E%88%E7%B2%89%E6%8A%80%E6%9C%AF%E8%A7%84%E7%A8%8B%28T-SHZSAQS+00241%E2%80%942024%29; cdb_back[p_pagecount]=8; cdb_back[state]=all; cdb_back[menuIndex]=2; cdb_back[pid]=74287670759082; cdb_RW_ID_1912168247=1; cdb_back[srlid]=b9f6mq9y7gPJ2hk40dqCH%2FgnRMmDglJwoznONTW58zP+BdE5CALigXcxWr3teIWPxkg0I7xemECrLGtohS4i6L1Va+2CCqvjFGlJtRdJqMgE; cdb_back[p_name]=%E6%96%B0%E7%96%86%E5%8D%97%E7%96%86%E5%86%AC%E5%B0%8F%E9%BA%A6%E5%B9%B2%E6%92%AD%E6%B9%BF%E5%87%BA%E9%AB%98%E4%BA%A7%E6%A0%BD%E5%9F%B9%E6%8A%80%E6%9C%AF%E8%A7%84%E7%A8%8B%28T-SHZSAQS+00234%E2%80%942024%29; cdb_back[rel_p_id]=1912168247; cdb_back[p_id]=1912168247; cdb_back[pcode]=74287670759082; cdb_back[_]=1704623575425; cdb_back[act]=get_new_user_task_degree; cdb_RW_ID_1634511839=17; cdb_RW_ID_1409870694=3; cdb_RW_ID_1467682874=6; cdb_RW_ID_1912369859=1; cdb_RW_ID_1912451073=1; cdb_RW_ID_1912290011=2; cdb_RW_ID_1912481477=1; cdb_logined=1; cdb_RW_ID_1912828853=1; cdb_RW_ID_1912878693=1; cdb_RW_ID_1062694187=40; cdb_RW_ID_1912878687=1; cdb_RW_ID_1912876498=1; cdb_RW_ID_1911422879=1; cdb_RW_ID_1910919563=2; show_index=1; cdb_RW_ID_1913002043=1; cdb_RW_ID_1457915724=42; Page_11273957435719=1; cdb_RW_ID_1916137836=1; Page_Y_94859698953258=-119.39144736842105; Page_94859698953258=1; cdb_RW_ID_1916292836=1; Page_90599692060782=1; cdb_RW_ID_1916500170=1; Page_Y_90799692311951=-119.39144736842105; Page_90799692311951=1; cdb_RW_ID_1916500980=1; Page_Y_70987675344694=-119.39144736842105; Page_70987675344694=1; cdb_RW_ID_1916500962=1; Page_Y_90799692311620=-122.80263157894737; Page_90799692311620=1; cdb_RW_ID_1916681511=2; Page_Y_69316364456866=337.81907894736844; Page_69316364456866=2; cdb_RW_ID_1916684147=1; Page_94659698827973=1; cdb_RW_ID_1916684176=1; Page_Y_69316364451624=-119.39144736842105; Page_69316364451624=1; cdb_RW_ID_1916684157=1; Page_Y_18361312295106=-119.39144736842105; Page_18361312295106=1; cdb_RW_ID_1916684107=1; Page_Y_73947976650784=-119.39144736842105; Page_73947976650784=1; cdb_RW_ID_1916684067=1; Page_Y_90699692274125=29.35855263157895; Page_90699692274125=3; cdb_RW_ID_1916981067=1; Page_Y_98199692679125=-119.39144736842105; Page_98199692679125=1; c_login_name=woyoceo; cdb_RW_ID_1916984890=1; cdb_READED_PC_ID=%2C; Page_Y_60716364351537=-119.39144736842105; Page_60716364351537=1; cdb_login_if=1; cdb_uid=104598337; cdb_RW_ID_1917205942=1; Page_Y_71787672043680=-119.39144736842105; Page_71787672043680=1; cdb_RW_ID_1917205923=1; Page_Y_60616362978390=568.7730263157895; Page_60616362978390=3; cdb_pageType=2; cdb_RW_ID_1916982030=1; Page_92529894830525=1; Page_Y_05029891159633=-213.70065789473685; cdb_RW_ID_1917701612=1; Page_Y_84861316641218=278.1611842105263; Page_84861316641218=1; cdb_RW_ID_1917701596=1; Page_Y_05029891159684=657.9177631578948; Page_05029891159684=4; cdb_RW_ID_1917701588=2; Page_05029891159633=3; cdb_RW_ID_1435851882=1; Page_Y_90559751219224=514.3618421052632; Page_90559751219224=10; doc88_lt=wx; cdb_tokenid=d3190%2BtXmQYDNmn5AZL9Q5UeWsavTx8ECVqc4kOmJZ2jiLCyC7Bc0ujSYz%2Fkq2jqOqVVK3MzGPKx783xNEQ8BRKW3CWOord1sK2JHn2pAMiMpYOWUEY3SscXPn2F17dYYg; cdb_RW_ID_1918007638=1; Page_Y_97116365772405=295.54276315789474; Page_97116365772405=2; cdb_RW_ID_1918007629=1; Page_Y_05729893551408=253.26644736842107; Page_05729893551408=7; cdb_RW_ID_1918008543=2; ShowSkinTip_1=1; showAnnotateTipIf=1; Page_Y_97116365775810=-238.38157894736844; Page_97116365775810=1; cdb_RW_ID_1918008545=1; Page_Y_18973436886595=5.088815789473685; Page_18973436886595=3; cdb_change_message=1; cdb_msg_num=0; cdb_RW_ID_1918815527=1; Page_37847975572234=1; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928feeda5c6169a4ebe4cd6d0a7e57ad73f20ed2a9a8be13a295adb2f27dcfc55127add2c2984ea49c396d186f19a3268e3e0624004b34b9699812f33eba660401bbe21a5fa05e706bb956dfac0eef035b2ed92f350c3f26f274; siftState=1; cdb_msg_time=1707639443"

// Doc88Cookie 15803889687
const Doc88Cookie = "cdb_sys_sid=vc0gvec15bo7h8rjq84llgd1h4; PHPSESSID=vc0gvec15bo7h8rjq84llgd1h4; cdb_change_message=1; cdb_msg_num=0; cdb_msg_time=1705637593; show_index=1; siftState=1; cdb_token=5176691bb4a2b7d6f48bef23f1c8c6bd6976e1452443c283d33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86c3480eba7246bba18f516cec2b411ba4e4de4720219466f830686a5776f5804dd60f82c713aca24f; cdb_tokenid=f4ec%2F1tvyaOMSK9UxFUgZJuEPiN3OwbisWAvNCi61AWsiQ9S50A2U2B66txBgWSEJVeGmWi8SkcMC8XLLbajoTPgd8Uy4U42fhoxFX2Ee%2FE; c_login_name=1056391268; cdb_login_if=1; cdb_system_sign=1; cdb_back[act]=get_intro; cdb_back[p]=2021%2F10%2F19%2F74687192717605.xml.gz; cdb_back[order]=0; cdb_back[pcid]=593; cdb_back[pid]=1905989379; cdb_back[curpage]=1; cdb_back[folder_id]=0; cdb_back[folder_page]=0; cdb_back[member_id]=6d4d209o34hyMZNLDyvpApTMemXbDwsEwc1QjGjbHoA%2FIZY1sHFLGSCJT7afUd0Ifg4; cdb_back[order_num]=1; cdb_back[show]=1; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[classify_id]=all; cdb_back[menuIndex]=2; cdb_back[show_index]=1; cdb_back[state]=all; cdb_back[t]=0; cdb_back[u]=1; cdb_back[log_state]=0; cdb_back[log_type_page]=all; cdb_back[captchaCode]=1; cdb_back[checkagreement]=on; cdb_back[txtPassword2]=abcdqq123456; cdb_back[txtPassword]=abcdqq123456; cdb_back[txtemail]=963094343%40qq.com; cdb_back[txtloginname]=woyoceo2023; cdb_back[username]=woyoceo2023; cdb_back[book_id]=3367789; cdb_back[classify]=8140; cdb_back[code]=66202dc48d76ef18; cdb_back[at]=0; cdb_back[doctype]=1; cdb_back[n]=6; cdb_back[p_id]=1401760345; cdb_back[p_name]=%E5%A4%A7%E5%AD%A6%E7%AA%81%E5%8F%91%E5%85%AC%E5%85%B1%E4%BA%8B%E4%BB%B6%E5%BA%94%E6%80%A5%E9%A2%84%E6%A1%88; cdb_back[rel_p_id]=1401760345; cdb_back[srlid]=128129RcKWcCMH%2F7zX0d+LXHyarhH5ko9rgGbBtxe8Y3NPDEfpvwVaPApe8esAntuluwmBC6fIe5K8RzK+HWKTRHm56Y4NVPn40Iskf0L84; cdb_READED_PC_ID=%2C; cdb_RW_ID_1401760345=3989; cdb_back[image_type]=9; cdb_back[login]=1; cdb_back[module_type]=7; cdb_back[email]=963094343%40qq.com; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dall"

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
			Price:   388,
		},
		{
			dirName: "finish-dbba.sacinfo.org.cn",
			pCid:    8371,
			Price:   388,
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/高考",
			pCid:    8244,
			Price:   388,
		},
		{
			dirName: "finish.topedu.ybep.com.cn/中考真题",
			pCid:    8243,
			Price:   388,
		},
		{
			dirName: "finish.topedu.ybep.com.cn/高考真题",
			pCid:    8244,
			Price:   388,
		},
		{
			dirName: "www.docx_shijuan1.com/中考试卷",
			pCid:    8243,
			Price:   388,
		},
		{
			dirName: "www.docx_shijuan1.com/高考试卷",
			pCid:    8244,
			Price:   388,
		},
		{
			dirName: "www.docx_shijuan1.com/小学试卷",
			pCid:    8242,
			Price:   388,
		},
		{
			dirName: "docx.51zjedu.com/一年级",
			pCid:    8242,
			Price:   388,
		},
		{
			dirName: "docx.51zjedu.com/二年级",
			pCid:    8242,
			Price:   388,
		},
		{
			dirName: "docx.51zjedu.com/三年级",
			pCid:    8242,
			Price:   388,
		},
		{
			dirName: "docx.51zjedu.com/四年级",
			pCid:    8242,
			Price:   388,
		},
		{
			dirName: "docx.51zjedu.com/五年级",
			pCid:    8242,
			Price:   388,
		},
		{
			dirName: "docx.51zjedu.com/六年级",
			pCid:    8242,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/小学/道德与法治",
			pCid:    8307,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/小学/美术",
			pCid:    8309,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/小学/数学",
			pCid:    8247,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/小学/信息技术",
			pCid:    8308,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/小学/音乐",
			pCid:    8309,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/小学/英语",
			pCid:    8248,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/小学/语文",
			pCid:    8246,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/初中/数学",
			pCid:    8156,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/初中/英语",
			pCid:    8157,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/初中/语文",
			pCid:    8155,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/初中/物理",
			pCid:    8158,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/初中/化学",
			pCid:    8160,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/初中/音乐",
			pCid:    55918,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/初中/道德与法治",
			pCid:    8163,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/初中/美术",
			pCid:    55918,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/初中/生物",
			pCid:    8162,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/初中/地理",
			pCid:    8159,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/初中/历史",
			pCid:    8161,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/高中/地理",
			pCid:    8256,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/高中/化学",
			pCid:    8183,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/高中/历史",
			pCid:    8257,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/高中/生物",
			pCid:    8186,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/高中/数学",
			pCid:    8180,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/高中/物理",
			pCid:    8182,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/高中/英语",
			pCid:    8181,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/高中/语文",
			pCid:    8179,
			Price:   388,
		},
		{
			dirName: "docx.zzstep.com/高中/政治",
			pCid:    8158,
			Price:   388,
		},
		{
			dirName: "docx.gzenxx.com",
			pCid:    55919,
			Price:   388,
		},
		{
			dirName: "docx.trjlseng.com",
			pCid:    8157,
			Price:   388,
		},
		{
			dirName: "yuwen.docx_chazidian.com",
			pCid:    55919,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/小学/政治",
			pCid:    8307,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/小学/数学",
			pCid:    8247,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/小学/英语",
			pCid:    8248,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/小学/语文",
			pCid:    8246,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/数学",
			pCid:    8156,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/英语",
			pCid:    8157,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/语文",
			pCid:    8155,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/物理",
			pCid:    8158,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/化学",
			pCid:    8160,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/政治",
			pCid:    8163,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/生物",
			pCid:    8162,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/地理",
			pCid:    8159,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/历史",
			pCid:    8161,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/地理",
			pCid:    8256,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/化学",
			pCid:    8183,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/历史",
			pCid:    8257,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/生物",
			pCid:    8186,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/数学",
			pCid:    8180,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/物理",
			pCid:    8182,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/英语",
			pCid:    8181,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/语文",
			pCid:    8179,
			Price:   388,
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/政治",
			pCid:    8158,
			Price:   388,
		},
		{
			dirName: "47.108.163.154",
			pCid:    8368,
			Price:   388,
		},
		{
			dirName: "www.webfree.net/国家标准",
			pCid:    8368,
			Price:   388,
		},
		{
			dirName: "www.webfree.net/行业标准",
			pCid:    8370,
			Price:   388,
		},
		{
			dirName: "www.webfree.net/地方标准",
			pCid:    8371,
			Price:   388,
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
		shuffleArray(files)
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

func shuffleArray(arr []fs.FileInfo) {
	rand.Seed(time.Now().UnixNano())
	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}
