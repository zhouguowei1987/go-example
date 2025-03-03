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
var SessionId = "8j5cbctli97mrdaipee47bjg75"
var Token = "8bbaa6e0e730a48ac48fc37851cf59bd"
var Cookie = "CLIENT_SYS_UN_ID=wKh2CmblP9zCxwkKfGiJAg==; UPLOAD_AGREEMENT_CHECKED=1; PHPSESSID=8j5cbctli97mrdaipee47bjg75; HMACCOUNT=2CEC63D57647BCA5; 94ca48fd8a42333b_code_getgraphcode=1739117139%2C1; max_u_token=125cff0af322e46301c57b92e2a0e10b; vip_left_recommend_doc_info=628987892; detail_show_similar=0; Hm_lvt_ed4f006fba260fb55ee1dfcb3e754e1c=1737643571,1739975137; a_8036113126007030=1; s_rfd=cdh%3D%3Ec865f4f0%7C%7C%7Ctrd%3D%3Emax.book118.com%7C%7C%7Cftrd%3D%3Ebook118.com; a_7010134064006105=1; 30d8fb61e609cac11=1739975175%2C1; a_7065105013010033=1; c4da14928424747de8b677208095de01=1740069471%2C1; Hm_lvt_f32e81852cb54f29133561587adb93c1=1739117135; s_v=cdh%3D%3Ec865f4f0%7C%7C%7Cvid%3D%3E1726300125295199813%7C%7C%7Cfsts%3D%3E1726300125%7C%7C%7Cdsfs%3D%3E160%7C%7C%7Cnps%3D%3E5; a_7132033111010036=1; Hm_lvt_27fe35f56bdde9c16e63129d678cd236=1740099639; Hm_lpvt_27fe35f56bdde9c16e63129d678cd236=1740099639; a_7105026003004042=1; a_5314332322004212=1; Hm_lvt_af8c54428f2dd7308990f5dd456fae6d=1737643628,1740099654; home_15045200=2; d6b93d63cc960c878126=1740099679%2C1; home_46465572=5; Hm_lpvt_af8c54428f2dd7308990f5dd456fae6d=1740099750; a_5341000301002021=1; a_5133202000000000=1; Hm_lvt_01a0a5632981ad913df7ee8d0d145f4c=1740099793; a_8037105117007033=1; a_8040025117007033=1; 5a9a221b83986f79ee93b689251380af=1740099869%2C2; Hm_lpvt_01a0a5632981ad913df7ee8d0d145f4c=1740099870; a_7042143142010036=1; s_m=34045279%3D%3Epreview_re%7C%7C%7C1189496135%3D%3Esecondcate_doclist_item_href%7C%7C%7C1872290310%3D%3Esecondcate_doclist_item_href%7C%7C%7Ccdh%3D%3Ec865f4f0%7C%7C%7C-1320480330%3D%3Esimilar%7C%7C%7C-734988587%3D%3Epreview_re%7C%7C%7C-997290743%3D%3Esecondcate_doclist_item_href%7C%7C%7C-1904860231%3D%3Esecondcate_doclist_item_href; s_s=cdh%3D%3Ec865f4f0%7C%7C%7Clast_req%3D%3E1740099901%7C%7C%7Csid%3D%3E1740099618721118748%7C%7C%7Cdsps%3D%3E1; Hm_lpvt_ed4f006fba260fb55ee1dfcb3e754e1c=1740099901; a_7042124142010036=1; PREVIEWHISTORYPAGES=727793067_2,727793080_3,727793169_2,211769600_4,430035420_2,727577224_1,720704027_2; ef7656dc08a0f1cf4c78acb87d97a1b9=1740184694%2C1; operation_user_center=1; Hm_lpvt_f32e81852cb54f29133561587adb93c1=1740242092"

// SessionId 15803889687
//var SessionId = "vb8kht61sj31fau20k1r4lejv3"
//var Token = "3cabba46df1b1ec810b2a7a65b828e01"
//var Cookie = "max_u_token=4a80ab01606cb39a35159e0376d51da9; Hm_lpvt_f32e81852cb54f29133561587adb93c1=1723646034; Hm_lvt_f32e81852cb54f29133561587adb93c1=1723645629; CLIENT_SYS_UN_ID=3rvhcma8vsEAxHvIa1ukAg==; HMACCOUNT=B9FF8ED9E02CA9F5; operation_user_center=1; 94ca48fd8a42333b_register_pc=1723645623%2C3; 94ca48fd8a42333b_register_checkgraphcode=1723645618%2C5; 94ca48fd8a42333b_register_checkusername=1723645603%2C2; 94ca48fd8a42333b_code_getgraphcode=1723645590%2C3; 94ca48fd8a42333b_register_checkmobilecode=1723645505%2C1; 94ca48fd8a42333b_code_getmobilecode=1723645474%2C1; 94ca48fd8a42333b_register_checkmobile=1723645473%2C1; PHPSESSID=vb8kht61sj31fau20k1r4lejv3; return_url=http%3A%2F%2Fmax.book118.com%2Fuser_center_v1%2Fdoc%2FIndex%2Findex.html"

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
