package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
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
)

// SessionId 15238369929
var SessionId = "upl1c9poffaqearfmdclkl0ru5"
var Token = "3cabba46df1b1ec810b2a7a65b828e01"
var Cookie = "PHPSESSID=upl1c9poffaqearfmdclkl0ru5; a_5001214301011114=1; a_8110060017006044=1; a_7103123105006051=1; a_5342222004010022=1; a_8104054004006045=1; a_8103126004006045=1; a_6013130131010050=1; a_6025040104010102=1; a_6021214104010102=1; a_5004303242011123=1; a_7131160111004152=1; a_6010131022002203=1; a_8035044027005111=1; a_7121003036004131=1; a_7004120132006053=1; a_6211042023001240=1; a_8035101016006033=1; home_51596649=1; a_8066031117004074=1; home_6907122=3; a_5224202203004314=1; home_4248920=2; a_5133202000000000=1; home_162=1; a_7165030124006060=1; a_7045043145006062=1; a_6143002243010113=1; a_5303324322011220=1; a_7042030005004126=1; a_5222013144004022=1; Hm_lvt_27fe35f56bdde9c16e63129d678cd236=1718439297; a_5232333324011040=1; a_8021021024006075=1; 30d8fb61e609cac11=1718448470%2C1; a_8017126106004101=1; Hm_lpvt_27fe35f56bdde9c16e63129d678cd236=1718449050; a_7054021151006126=1; HMACCOUNT=487EF362690A1D5D; a_5344222244011240=1; a_6120010043010204=1; a_6241044143010205=1; a_8141030077006115=1; a_7165163120006140=1; a_8134101077006115=1; a_8014065011006117=1; a_6234145141010211=1; a_8005031037002052=1; a_5144131204011310=1; a_8112117055006121=1; a_6153234152004141=1; a_8006034074006122=1; a_5024232124010331=1; a_8004115074006122=1; home_46465572=13; a_7044043062006146=1; a_6105132212010224=1; CLIENT_SYS_UN_ID=3rvhcmcfS1q3WmPXIfihAg==; 23180e48f45d3ffdf83d11798c3cf2df=1732012068%2C1; a_7015041021010003=1; a_6132002022011003=1; a_8106135077007005=1; return_url=http%3A%2F%2Fmax.book118.com%2Fuser_center_v1%2Fdoc%2Findex%2Findex.html; a_7015163064010013=1; a_7142164035005012=1; a_8011102054003133=1; detail_show_similar=0; a_5021041141012020=1; 0ab7d33081eeb=2022%E5%B9%B4%E5%9B%9B%E5%B7%9D%E7%9C%81%E8%87%AA%E8%B4%A1%E5%B8%82%E4%B8%AD%E8%80%83%E5%8C%96%E5%AD%A6%E8%AF%95%E5%8D%B7%EF%BC%88%E5%90%AB%E8%A7%A3%E6%9E%90%E7%89%88%EF%BC%89; input_search_logs=a%3A1%3A%7Bi%3A0%3Ba%3A2%3A%7Bs%3A8%3A%22keywords%22%3Bs%3A61%3A%222022%E5%B9%B4%E5%9B%9B%E5%B7%9D%E7%9C%81%E8%87%AA%E8%B4%A1%E5%B8%82%E4%B8%AD%E8%80%83%E5%8C%96%E5%AD%A6%E8%AF%95%E5%8D%B7%EF%BC%88%E5%90%AB%E8%A7%A3%E6%9E%90%E7%89%88%EF%BC%89%22%3Bs%3A4%3A%22time%22%3Bi%3A1735272270%3B%7D%7D; Hm_lvt_90eba6da4a0cf4f7f07614e930242936=1735272273; Hm_lpvt_90eba6da4a0cf4f7f07614e930242936=1735272279; a_6045050241011024=1; a_7043140166010022=1; a_7101032025010022=1; a_5143112204012041=1; Hm_lvt_01a0a5632981ad913df7ee8d0d145f4c=1738324767; 5a9a221b83986f79ee93b689251380af=1738324795%2C2; Hm_lpvt_01a0a5632981ad913df7ee8d0d145f4c=1738324796; s_m=174645614%3D%3Esimilar%7C%7C%7C377633538%3D%3Epreview_re%7C%7C%7C908837254%3D%3Esimilar%7C%7C%7C1841443071%3D%3Esimilar%7C%7C%7C1978635343%3D%3Esimilar%7C%7C%7C1982161289%3D%3Esimilar%7C%7C%7Ccdh%3D%3Ec865f4f0%7C%7C%7C-425000343%3D%3Esimilar%7C%7C%7C-897180436%3D%3Esimilar%7C%7C%7C-780484034%3D%3Esimilar%7C%7C%7C-1911702904%3D%3Esimilar%7C%7C%7C-982241833%3D%3Esimilar%7C%7C%7C-660567806%3D%3Esearch_filter_item_href%7C%7C%7C-25029623%3D%3Esecondcate_doclist_item_href; a_8010036121007026=1; home_50939525=1; d6b93d63cc960c878126=1738324831%2C1; Hm_lvt_af8c54428f2dd7308990f5dd456fae6d=1738324833; Hm_lpvt_af8c54428f2dd7308990f5dd456fae6d=1738324833; a_8113022016007031=1; a_6012003212011043=1; a_6041102003004001=1; ef7656dc08a0f1cf4c78acb87d97a1b9=1741931420%2C1; 94ca48fd8a42333b_login_passwordlogin=1742706951%2C1; Hm_lvt_ed4f006fba260fb55ee1dfcb3e754e1c=1741930141,1744418820; Hm_lpvt_ed4f006fba260fb55ee1dfcb3e754e1c=1744418820; a_7101164052010003=1; s_v=cdh%3D%3Ec865f4f0%7C%7C%7Cvid%3D%3E1695983286868757343%7C%7C%7Cfsts%3D%3E1695983286%7C%7C%7Cdsfs%3D%3E561%7C%7C%7Cnps%3D%3E143; s_rfd=cdh%3D%3Ec865f4f0%7C%7C%7Ctrd%3D%3Ewww.baidu.com%7C%7C%7Cftrd%3D%3Ebaidu.com; s_s=cdh%3D%3Ec865f4f0%7C%7C%7Clast_req%3D%3E1744418820%7C%7C%7Csid%3D%3E1744418820838268523%7C%7C%7Cdsps%3D%3E29; Hm_lvt_f32e81852cb54f29133561587adb93c1=1746182314; a_5120024002012200=1; DOC_INFORM=read; 94ca48fd8a42333b_code_getgraphcode=1747901554%2C1; max_u_token=7f9e225b8297ddc92f901482cb99705d; operation_user_center=1; c4da14928424747de8b677208095de01=1748102531%2C2; Hm_lpvt_f32e81852cb54f29133561587adb93c1=1748102547"

// 金币上传 MoldType:0 CoinScoreType:0
// 积分上传  MoldType:4 CoinScoreType:4

var MoldType = "0"
var CoinScoreType = "0"

type VerifyUploadDocumentResponse struct {
	Code    string                           `json:"code"`
	Data    VerifyUploadDocumentResponseData `json:"data"`
	Message string                           `json:"message"`
}

type VerifyUploadDocumentResponseData struct {
	IsAllowUpload string `json:"isAllowUpload"`
	Reason        string `json:"reason"`
}

func VerifyUploadDocument(title string, format string, price string, md5 string) (isAllowUpload bool, err error) {
	client := &http.Client{} //初始化客户端
	postData := url.Values{}
	postData.Add("mold_type", MoldType)
	postData.Add("type", CoinScoreType)
	postData.Add("session_id", SessionId)
	postData.Add("title", title)
	postData.Add("format", format)
	postData.Add("systemCategory", "0")
	postData.Add("folder", "0")
	postData.Add("price", price)
	switch MoldType {
	case strconv.Itoa(0):
		// 金币上传
		postData.Add("readPrice", "0")
		postData.Add("reeReadPage", "0")
		break
	case strconv.Itoa(4):
		// 积分上传
		break
	}
	postData.Add("contentMD5", md5)
	requestUrl := "https://max.book118.com/user_center_v1/upload/Api/verifyUploadDocument.html"
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	req.Header.Set("Accept", "*/*")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", Cookie)
	req.Header.Set("Host", "max.book118.com")
	req.Header.Set("Origin", "https://max.book118.com")
	referer := "https://max.book118.com/user_center_v1/upload/Upload/ordinary.html"
	switch MoldType {
	case strconv.Itoa(0):
		// 金币上传
		referer = "https://max.book118.com/user_center_v1/upload/Upload/ordinary.html"
		break
	case strconv.Itoa(4):
		// 积分上传
		referer = "https://max.book118.com/user_center_v1/home/reward/index.html"
		break
	}
	req.Header.Set("Referer", referer)
	req.Header.Set("Sec-Ch-Ua", "\"Not/A)Brand\";v=\"99\", \"Google Chrome\";v=\"115\", \"Chromium\";v=\"115\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return false, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	respBytes, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(respBytes))
	if err != nil {
		return false, err
	}
	verifyUploadDocumentResponse := VerifyUploadDocumentResponse{}
	err = json.Unmarshal(respBytes, &verifyUploadDocumentResponse)
	if err != nil {
		return false, err
	}
	if verifyUploadDocumentResponse.Data.IsAllowUpload != "1" {
		return false, errors.New(verifyUploadDocumentResponse.Data.Reason)
	}
	return true, nil
}

type GetDocCateResponse struct {
	Code    int32                  `json:"code"`
	Data    GetDocCateResponseData `json:"data"`
	Message string                 `json:"message"`
}

type GetDocCateResponseData struct {
	CateId   string `json:"cate_id"`
	CateName string `json:"cate_name"`
}

func GetDocCate(title string) (systemCategory GetDocCateResponseData, err error) {
	client := &http.Client{} //初始化客户端
	postData := url.Values{}
	postData.Add("title", title)
	requestUrl := "https://max.book118.com/user_center_v1/upload/Api/getDocCate.html"
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	req.Header.Set("Accept", "*/*")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", Cookie)
	req.Header.Set("Host", "max.book118.com")
	req.Header.Set("Origin", "https://max.book118.com")
	req.Header.Set("Referer", "https://max.book118.com/user_center_v1/home/reward/index.html")
	req.Header.Set("Sec-Ch-Ua", "\"Not/A)Brand\";v=\"99\", \"Google Chrome\";v=\"115\", \"Chromium\";v=\"115\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	systemCategory = GetDocCateResponseData{}
	if err != nil {
		return systemCategory, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return systemCategory, err
	}
	getDocCateResponse := GetDocCateResponse{}
	err = json.Unmarshal(respBytes, &getDocCateResponse)
	if err != nil {
		return systemCategory, err
	}
	systemCategory.CateId = getDocCateResponse.Data.CateId
	systemCategory.CateName = getDocCateResponse.Data.CateName
	return systemCategory, nil
}

type Book118UploadResponse struct {
	Code    string                    `json:"code"`
	Data    Book118UploadResponseData `json:"data"`
	Message string                    `json:"message"`
}

type Book118UploadResponseData struct {
	Aid                  string `json:"aid"`
	AuditScore           int32  `json:"audit_score"`
	NextUploadScore      int32  `json:"next_upload_score"`
	RemainNumber         int32  `json:"remainNumber"`
	UploadRewardAllScore int32  `json:"upload_reward_all_score"`
	UploadRewardScore    int32  `json:"upload_reward_score"`
	UseNumber            int32  `json:"useNumber"`
}

// Book18Upload 上传文件
func Book18Upload(filePath string, id string, md5 string, title string, systemCategory string, price string) (uploadResponse Book118UploadResponse, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	file, err := os.Open(filePath)
	if err != nil {
		return uploadResponse, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return uploadResponse, err
	}

	// 获取文件大小（字节数）
	fileSize := fileInfo.Size()

	// 获取文件修改时间
	modTime := fileInfo.ModTime()
	formattedTime := modTime.Format("Mon Jan 02 2006 15:04:05 GMT-0700 (MST)")

	fileWriter, err := bodyWriter.CreateFormFile("single", filepath.Base(file.Name()))
	if err != nil {
		return uploadResponse, err
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return uploadResponse, err
	}

	contentType := bodyWriter.FormDataContentType()
	err = bodyWriter.Close()
	if err != nil {
		return uploadResponse, err
	}

	err = bodyWriter.WriteField("mold_type", MoldType)
	if err != nil {
		return uploadResponse, err
	}
	err = bodyWriter.WriteField("type", CoinScoreType)
	if err != nil {
		return uploadResponse, err
	}
	err = bodyWriter.WriteField("session_id", SessionId)
	if err != nil {
		return uploadResponse, err
	}
	err = bodyWriter.WriteField("token", Token)
	if err != nil {
		return uploadResponse, err
	}
	err = bodyWriter.WriteField("uploadKeyword", "0")
	if err != nil {
		return uploadResponse, err
	}
	err = bodyWriter.WriteField("id", "WU_FILE_"+id)
	if err != nil {
		return uploadResponse, err
	}
	err = bodyWriter.WriteField("name", file.Name())
	if err != nil {
		return uploadResponse, err
	}
	err = bodyWriter.WriteField("lastModifiedDate", formattedTime)
	if err != nil {
		return uploadResponse, err
	}

	err = bodyWriter.WriteField("size", strconv.Itoa(int(fileSize)))
	if err != nil {
		return uploadResponse, err
	}
	err = bodyWriter.WriteField("md5", md5)
	if err != nil {
		return uploadResponse, err
	}
	err = bodyWriter.WriteField("title", title)
	if err != nil {
		return uploadResponse, err
	}
	err = bodyWriter.WriteField("systemCategory", systemCategory)
	if err != nil {
		return uploadResponse, err
	}
	err = bodyWriter.WriteField("price", price)
	if err != nil {
		return uploadResponse, err
	}

	uploadUrl := "https://upfile9.book118.com/upload/single/upload"
	req, err := http.NewRequest("POST", uploadUrl, bodyBuf)
	if err != nil {
		return uploadResponse, err
	}

	req.Header.Set("Accept", "*/*")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Cookie", Cookie)
	req.Header.Set("Host", "max.book118.com")
	req.Header.Set("Origin", "https://max.book118.com")
	req.Header.Set("Referer", "https://max.book118.com/user_center_v1/home/reward/index.html")
	req.Header.Set("Sec-Ch-Ua", "\"Not/A)Brand\";v=\"99\", \"Google Chrome\";v=\"115\", \"Chromium\";v=\"115\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return uploadResponse, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	respBytes, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBytes))
	//os.Exit(1)
	err = json.Unmarshal(respBytes, &uploadResponse)
	if err != nil {
		return uploadResponse, err
	}
	return uploadResponse, nil
}

func getFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	md5Hash := hash.Sum(nil)
	md5String := hex.EncodeToString(md5Hash)

	return md5String, nil
}

type Book118UploadChildDir struct {
	dirName string
	price   string
}

func main() {
	var uploadChildDirArr = []Book118UploadChildDir{
		{
			dirName: "finish.tikuvip（2023）.51test.net",
			price:   "88",
		},
		{
			dirName: "docx.topedu.ybep.com.cn/中考真题",
			price:   "88",
		},
		{
			dirName: "docx.topedu.ybep.com.cn/高考真题",
			price:   "88",
		},
		{
			dirName: "www.docx_shijuan1.com/中考试卷",
			price:   "88",
		},
		{
			dirName: "www.docx_shijuan1.com/高考试卷",
			price:   "88",
		},
		{
			dirName: "www.docx_shijuan1.com/小学试卷",
			price:   "88",
		},
		{
			dirName: "docx.51zjedu.com/一年级",
			price:   "88",
		},
		{
			dirName: "docx.51zjedu.com/二年级",
			price:   "88",
		},
		{
			dirName: "docx.51zjedu.com/三年级",
			price:   "88",
		},
		{
			dirName: "docx.51zjedu.com/四年级",
			price:   "88",
		},
		{
			dirName: "docx.51zjedu.com/五年级",
			price:   "88",
		},
		{
			dirName: "docx.51zjedu.com/六年级",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/小学/道德与法治",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/小学/美术",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/小学/数学",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/小学/信息技术",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/小学/音乐",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/小学/英语",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/小学/语文",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/初中/数学",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/初中/语文",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/初中/英语",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/初中/物理",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/初中/化学",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/初中/音乐",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/初中/道德与法治",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/初中/美术",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/初中/生物",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/初中/地理",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/初中/历史",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/高中/地理",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/高中/化学",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/高中/历史",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/高中/生物",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/高中/数学",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/高中/物理",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/高中/英语",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/高中/语文",
			price:   "88",
		},
		{
			dirName: "docx.zzstep.com/高中/政治",
			price:   "88",
		},
		{
			dirName: "docx.gzenxx.com",
			price:   "88",
		},
		{
			dirName: "docx.trjlseng.com",
			price:   "88",
		},
		{
			dirName: "yuwen.docx_chazidian.com",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/小学/政治",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/小学/数学",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/小学/英语",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/小学/语文",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/数学",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/语文",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/英语",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/物理",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/化学",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/政治",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/生物",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/地理",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/初中/历史",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/地理",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/化学",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/历史",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/生物",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/数学",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/物理",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/英语",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/语文",
			price:   "88",
		},
		{
			dirName: "bk.docx_cooco.net.cn/高中/政治",
			price:   "88",
		},
		{
			dirName: "sc.chinaz.com/ppt",
			price:   "88",
		},
		{
			dirName: "www.ppt818.com",
			price:   "88",
		},
		{
			dirName: "docx.jianlisheji.com",
			price:   "88",
		},
		{
			dirName: "docx.www.czsx.com.cn",
			price:   "88",
		},
		{
			dirName: "docx.www.word_1ppt.com",
			price:   "88",
		},
		{
			dirName: "www.finish_zip_1ppt.com",
			price:   "88",
		},
		{
			dirName: "docx.officeplus.cn",
			price:   "88",
		},
		{
			dirName: "ppt.officeplus.cn",
			price:   "88",
		},
		{
			dirName: "www.docx_ychedu.com/中考试题",
			price:   "88",
		},
		{
			dirName: "www.docx_ychedu.com/高考试题",
			price:   "88",
		},
		{
			dirName: "hao123.docx_eduy.net/小学",
			price:   "88",
		},
		{
			dirName: "hao123.docx_eduy.net/初中",
			price:   "88",
		},
		{
			dirName: "hao123.docx_eduy.net/高中",
			price:   "88",
		},
	}
	rootPath := "../upload.book118.com/"
	for _, childDir := range uploadChildDirArr {
		childDirPath := rootPath + childDir.dirName + "/"
		fmt.Println(childDirPath)
		files, err := ioutil.ReadDir(childDirPath)
		if err != nil {
			continue
		}
		id := 0
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			fileName := file.Name()
			if fileName == ".DS_Store" {
				continue
			}
			fileExt := path.Ext(fileName)
			fileExt = strings.ReplaceAll(fileExt, ".", "")

			filePath := childDirPath + fileName
			fmt.Println(filePath)

			price := childDir.price
			filePageNum := 0
			if fileExt == "pdf" {
				// 获取PDF文件，获取总页数
				if pdfFile, err := pdf.Open(filePath); err == nil {
					filePageNum = pdfFile.NumPage()
				}
			}
			// 根据页数设置价格
			if filePageNum > 0 {
				if filePageNum > 0 && filePageNum <= 5 {
					price = "28"
				} else if filePageNum > 5 && filePageNum <= 10 {
					price = "38"
				} else if filePageNum > 10 && filePageNum <= 15 {
					price = "48"
				} else if filePageNum > 15 && filePageNum <= 20 {
					price = "58"
				} else if filePageNum > 20 && filePageNum <= 25 {
					price = "68"
				} else if filePageNum > 25 && filePageNum <= 30 {
					price = "78"
				} else if filePageNum > 30 && filePageNum <= 35 {
					price = "88"
				} else if filePageNum > 35 && filePageNum <= 40 {
					price = "98"
				} else if filePageNum > 40 && filePageNum <= 45 {
					price = "108"
				} else if filePageNum > 45 && filePageNum <= 50 {
					price = "118"
				} else {
					price = "128"
				}
			}

			fmt.Println("==========开始上传==============")

			fileMD5, err := getFileMD5(filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(fileMD5)
			fmt.Println(fileName)
			fmt.Println(fileExt)
			// 验证是否可以上传
			isAllowUpload, err := VerifyUploadDocument(fileName, fileExt, price, fileMD5)
			if err != nil || isAllowUpload == false {
				fmt.Printf("isAllowUpload = %t, err = %s", isAllowUpload, err)
				break
			}
			fmt.Printf("isAllowUpload = %t\n", isAllowUpload)

			title := strings.ReplaceAll(fileName, "."+fileExt, "")
			// 获取文档所属分类
			systemCategory, err := GetDocCate(title)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(systemCategory)
			uploadResponseData, err := Book18Upload(filePath, strconv.Itoa(id), fileMD5, title, systemCategory.CateId, price)
			if err != nil {
				fmt.Println(err)
				// 删除源文件，继续
				err := os.Remove(filePath)
				if err != nil {
					return
				}
				continue
			}
			fmt.Println(uploadResponseData)
			fmt.Println("==========将已上传的文件转移到指定文件夹==============")

			// 将上传过文件移动到"../final-upload.book118.com/"
			finalDir := "../final-upload.book118.com/" + childDir.dirName
			if _, err = os.Stat(finalDir); err != nil {
				if os.MkdirAll(finalDir, 0777) != nil {
					fmt.Println(err)
					break
				}
			}

			// 将已上传的文件转移到指定文件夹
			fileFinal := finalDir + "/" + fileName
			err = os.Rename(filePath, fileFinal)
			if err != nil {
				fmt.Println(err)
				break
			}

			id++
			fmt.Println("==========上传完成==============")
		}
	}
}
