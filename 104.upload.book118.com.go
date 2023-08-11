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

var SessionId = "ga4qcm6cbo7ghj6d4m4d4k0un2"
var Token = "c266bc25b785c93dffe6c493741f1cc7"
var Cookie = "__yjs_duid=1_1543d26121978a9cfb0ca147de19aa051678550479017; a_7045105123005010=1; a_8001130106005057=1; a_5211220343010142=1; a_5202134313010141=1; a_8027023013005057=1; a_8002052065005056=1; a_8141042073005055=1; a_7165005113005063=1; a_5024220320010133=1; a_7035023056005100=1; a_7020036105005102=1; a_7006035150005066=1; a_6012021104005124=1; a_7011146163005102=1; a_7142121125005103=1; a_7130145125005103=1; a_6133210224005124=1; reward_download_aid=552885778; user_download_check=505ec9ed820a40f3dd5f5b62f2f24eb0; a_8116127104005064=1; a_8116076104005064=1; a_8027101137005064=1; a_8023105137005064=1; a_8026107137005064=1; a_6040200122005050=1; a_6101045012004051=1; a_8122006065005061=1; home_14038769=1; a_6214042040005125=1; a_5224120321010203=1; a_5222341321010203=1; a_6132112222005125=1; a_7113161152005104=1; a_8127021056005066=1; a_8126004056005066=1; a_5012022032010204=1; a_8016022113005066=1; a_6021222203005130=1; a_8123075142005070=1; a_8123031142005070=1; a_6014054103005133=1; a_6210204222005140=1; a_7145136036005114=1; a_8105017025005072=1; a_5042213311010112=1; a_6130111042005140=1; a_8101063057005075=1; a_5244140311010002=1; home_4889072=1; a_6110042115005141=1; a_5102311142010221=1; a_8106074026005077=1; a_6242215034005143=1; Hm_lvt_b65e7ffe374e9dc8240fbd00b9336d29=1686761455; Hm_lpvt_b65e7ffe374e9dc8240fbd00b9336d29=1686761458; a_6104034040005143=1; a_6125051222005140=1; a_6015002134005140=1; a_6234125133005140=1; a_7024002144005120=1; a_8013143124005060=1; a_6045200242005115=1; a_6123102233002051=1; a_7160154030005125=1; a_6231133033005152=1; a_6134050034005150=1; PHPSESSID=ga4qcm6cbo7ghj6d4m4d4k0un2; UPLOAD_AGREEMENT_CHECKED=1; TRANSFORM_USER_CHECK_AGREEMENT=read; a_5340134041010234=1; a_8035053132005105=1; a_8130012105005104=1; a_3616942=1; home_42890=1; a_7010136053005130=1; a_7133145002005123=1; a_8002127046005106=1; a_5243300002010231=1; a_5034121243010040=1; a_8041116106004032=1; a_5242014030010240=1; a_7031104163004013=1; a_5203311002010231=1; Hm_lvt_f32e81852cb54f29133561587adb93c1=1689041030; a_5042104222010241=1; a_5202010001004103=1; Hm_lvt_ed4f006fba260fb55ee1dfcb3e754e1c=1689552477; a_8037064027005111=1; a_8143120010005111=1; a_7043064032005133=1; a_7023056126005133=1; a_8043007133005102=1; a_6232154110001214=1; a_5314224343010243=1; CLIENT_SYS_UN_ID=3rvgCmS4AE8X0wy/Ay6AAg==; a_8101112051004046=1; a_8011042103004130=1; a_6134123015005202=1; a_6024230103005202=1; a_8126131132005112=1; a_7150060031005135=1; a_8124042026005113=1; a_5340141113010234=1; a_7133062156005134=1; a_8112020113005113=1; a_8111113113005113=1; a_6024130232005203=1; a_7154046026005136=1; a_8055044063005114=1; a_7050161102005136=1; a_8004015062005114=1; a_5333123144010301=1; home_20672244=1; a_6105014121005204=1; 5a9a221b83986f79ee93b689251380af=1690521796%2C9; a_5240102300010300=1; a_6110035050004201=1; a_5244313141010301=1; a_8124042056005114=1; a_6034204243005204=1; a_6030004023005145=1; home_46465572=27; d6b93d4rgc960c878126=1690791644%2C4; a_8143115027005115=1; a_5344134043010302=1; s_m=168135363%3D%3Evipdocs%7C%7C%7C225152386%3D%3Esecondcate_doclist_item_href%7C%7C%7C629810427%3D%3Esecondcate_doclist_item_href%7C%7C%7C1717957616%3D%3Esecondcate_doclist_item_href%7C%7C%7Ccdh%3D%3E27a30245%7C%7C%7C-2065366703%3D%3Esimilar%7C%7C%7C-2129424524%3D%3Esecondcate_doclist_item_href%7C%7C%7C-272989638%3D%3Esecondcate_doclist_item_href%7C%7C%7C-852022213%3D%3Esecondcate_nav_href%7C%7C%7C-1139764172%3D%3Esimilar%7C%7C%7C-836792867%3D%3Erelate%7C%7C%7C-272279446%3D%3Esimilar; a_8050011043005106=1; a_5044000241010301=1; a_7035065140005140=1; return_url=http%3A%2F%2Fmax.book118.com%2Fuser_center_v1%2Fdoc%2FIndex%2Findex.html; max_u_token=01573740a0d711710fa7439681aa499f; a_5341203012010303=1; a_8002131007005106=1; a_5310131212010303=1; s_rfd=cdh%3D%3E27a30245%7C%7C%7Ctrd%3D%3Emax.book118.com%7C%7C%7Cftrd%3D%3Ebook118.com; a_8120046071005116=1; a_5310112212010303=1; a_8117142071005116=1; a_5304244212010303=1; c4da14928424747de8b677208095de01=1691370880%2C2; detail_show_similar=0; s_v=cdh%3D%3E27a30245%7C%7C%7Cvid%3D%3E1672140928782145224%7C%7C%7Cfsts%3D%3E1672140928%7C%7C%7Cdsfs%3D%3E223%7C%7C%7Cnps%3D%3E107; a_6203232220005210=1; operation_user_center=1; buy_vip_from_aid=578844627; Hm_lpvt_ed4f006fba260fb55ee1dfcb3e754e1c=1691378413; CRM_DETAIL_INFOS=[{\"aid\":8056033124005116,\"title\":\"2023å¹´æ²³å\u008D—ä¿¡é˜³ä¸­è€ƒç”Ÿç‰©çœŸé¢˜å\u008FŠç­”æ¡ˆ.pdf\",\"firstType\":\"622\",\"secondType\":\"631\"},{\"aid\":6203232220005210,\"title\":\"2023å¹´æ²³å\u008D—é©»é©¬åº—ä¸­è€ƒç”Ÿç‰©çœŸé¢˜å\u008FŠç­”æ¡ˆ.pdf\",\"firstType\":\"622\",\"secondType\":\"631\"},{\"aid\":5310131212010303,\"title\":\"2023å¹´æ²³å\u008D—æ¿®é˜³ä¸­è€ƒé\u0081“å¾·ä¸Žæ³•æ²»çœŸé¢˜å\u008FŠç­”æ¡ˆ.pdf\",\"firstType\":\"622\",\"secondType\":\"631\"}]; s_s=cdh%3D%3E27a30245%7C%7C%7Clast_req%3D%3E1691378413%7C%7C%7Csid%3D%3E1691378061146912756%7C%7C%7Cdsps%3D%3E0; a_8056033124005116=1; 94ca48fd8a42333b=1691378413%2C2; PREVIEWHISTORYPAGES=578844627_3,578577974_1,570070289_1,578079653_1,577772647_3,576712400_1,576468434_1,576467483_1,473304223_3,576494110_1,576499338_1,576500413_1,576513592_3,576514536_2,576208834_1,575921654_2,575757416_4,574907344_1,574908689_1,574391690_1,574115851_1,488670934_3,182429270_1,573691741_1,573233152_1,571622229_1,520731936_1,566027382_1,569902943_3,569219544_2,231935138_2,547982972_1,560579453_1,561472781_1,558216915_1,557391034_1,556988361_1,556988325_1,554751418_1,554170712_2,554468604_2,554468717_2,553865992_2,553866435_3,552952365_2,552951969_1,552687964_1,552687082_1,543851460_1,546530242_2; Hm_lpvt_f32e81852cb54f29133561587adb93c1=1691380838"

// MoldType 金币上传
var MoldType = "0"
var CoinScoreType = "application/pdf"

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
	postData.Add("mold_type", "4")
	postData.Add("type", "4")
	postData.Add("session_id", SessionId)
	postData.Add("title", title)
	postData.Add("format", format)
	postData.Add("systemCategory", "0")
	postData.Add("folder", "0")
	postData.Add("price", price)
	postData.Add("contentMD5", md5)
	requestUrl := "https://max.book118.com/user_center_v1/upload/Api/verifyUploadDocument.html"
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
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
	if err != nil {
		return false, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	respBytes, err := ioutil.ReadAll(resp.Body)
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

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
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
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
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
	format  string
	price   string
}

func main() {
	var uploadChildDirArr = []Book118UploadChildDir{
		{
			dirName: "finish.tikuvip（2023）.51test.net",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/专升本考试",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/初中一年级",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/中考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/小升初",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/成人高考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/考研",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/自考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/高中会考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/高考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/初中一年级",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/中考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/小升初",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/成人高考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/考研",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/高中会考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/高考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/专升本考试",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/中考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/小升初",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/成人高考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/考研",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/自考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/高中会考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/高考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/专升本考试",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/中考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/小升初",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/成人高考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/考研",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/自考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/高中会考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/高考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/专升本考试",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/中考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/小升初",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/成人高考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/考研",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/自考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/高中会考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/高考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/专升本考试",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/中考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/小升初",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/成人高考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/考研",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/自考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/高中会考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/高考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/专升本考试",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/中考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/小升初",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/成人高考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/考研",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/自考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/高中会考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/高考",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.topedu.ybep.com.cn/中考真题",
			format:  "pdf",
			price:   "2000",
		},
		{
			dirName: "finish.topedu.ybep.com.cn/高考真题",
			format:  "pdf",
			price:   "2000",
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
			fileExt := path.Ext(fileName)
			fileExt = strings.ReplaceAll(fileExt, ".", "")
			if fileExt != childDir.format {
				continue
			}

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
				if filePageNum > 0 && filePageNum <= 8 {
					price = "288"
				} else if filePageNum > 8 && filePageNum <= 18 {
					price = "388"
				} else if filePageNum > 18 && filePageNum <= 28 {
					price = "488"
				} else if filePageNum > 28 && filePageNum <= 38 {
					price = "588"
				} else if filePageNum > 38 && filePageNum <= 48 {
					price = "688"
				} else if filePageNum > 48 && filePageNum <= 58 {
					price = "788"
				} else {
					price = "888"
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
			// 验证是否可以上传
			isAllowUpload, err := VerifyUploadDocument(fileName, childDir.format, price, fileMD5)
			if err != nil || isAllowUpload == false {
				fmt.Printf("isAllowUpload = %t, err = %s", isAllowUpload, err)
				break
			}
			fmt.Printf("isAllowUpload = %t\n", isAllowUpload)

			title := strings.ReplaceAll(fileName, "."+childDir.format, "")
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
				break
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
