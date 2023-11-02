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

var SessionId = "dkco0tmnr8ui540b5crck40r97"
var Token = "74211dfcb74e901fd6574df78c76ed8e"
var Cookie = "__yjs_duid=1_1543d26121978a9cfb0ca147de19aa051678550479017; a_5020110024010320=1; a_7033114052005130=1; a_6014015022005221=1; a_5310302202010240=1; a_8001140054005125=1; a_5001302134010320=1; a_8002003054005125=1; a_6142142045001243=1; home_4070593=1; a_5330103012010321=1; a_5330001012010321=1; return_url=http%3A%2F%2Fmax.book118.com%2Fuser_center_v1%2Fdoc%2Findex%2Findex.html; a_6135045102005222=1; a_8040015104005126=1; a_8037073104005126=1; a_8032107104005126=1; a_8051102001005127=1; a_7134005201004040=1; a_5241233110010322=1; a_8025023075005127=1; __51cke__=; a_5041234221010322=1; a_8002065133005127=1; a_7012134026005154=1; a_7012110026005154=1; a_6012210032005224=1; a_7011200026005154=1; a_6230135243005224=1; a_8132054143005130=1; a_8113006031005131=1; a_8124113054004013=1; a_5321041201010324=1; a_8135141116005131=1; a_7163061141005155=1; a_5100310012010330=1; a_8066006044005132=1; a_7151161055003050=1; a_5343021240003110=1; home_15265435=1; d6b93d4rgc960c878126=1695001563%2C1; a_6151204214002110=1; a_8027101133005132=1; a_7042124016002155=1; a_7030033163003020=1; a_6043234030005231=1; reward_download_aid=591182794; tongji_46465572=1; Hm_lvt_b645044a3b9e8b6315c6fe7d4733b16c=1693234255,1695106630; Hm_lpvt_b645044a3b9e8b6315c6fe7d4733b16c=1695106794; a_8033033022005133=1; a_8017012061005074=1; a_5102314033010331=1; a_8137041073005133=1; a_117943075=1; 5a9a221b83986f79ee93b689251380af=1695175374%2C1; a_5340114214010331=1; a_6120202040005214=1; a_8123100103005133=1; a_8122132103005133=1; a_6100033140005233=1; a_7051014114005162=1; a_7050155114005162=1; a_7046156114005162=1; a_7031050144005163=1; a_8101064105003123=1; a_8025142121005136=1; a_6052220222003220=1; a_5041220311010334=1; a_6033135213005234=1; a_6034104213005234=1; a_8032032115005115=1; a_8026013121005136=1; a_5041244311010334=1; a_8106027102005136=1; a_7112013110005163=1; Hm_lvt_f32e81852cb54f29133561587adb93c1=1696826476; a_7066045001005165=1; a_7065110001005165=1; a_8057065001005140=1; a_5142000001010341=1; a_6112155001005240=1; a_7132065156005164=1; a_5331214213004213=1; a_7052141103005165=1; a_6101111124005240=1; a_5301143302010341=1; a_8113102115005140=1; a_5243132302010341=1; a_5301003302010341=1; a_5243033043010304=1; a_7125101030005142=1; a_6040233011005241=1; a_7043021123003140=1; a_6040133011005241=1; a_6040131011005241=1; a_5044040012010342=1; a_5242342302010341=1; a_8135005072005141=1; a_7015156114002025=1; a_8046050120002105=1; a_5204040210004012=1; a_6045151210001210=1; a_27447145=1; a_5020303312010010=1; a_5034140031002023=1; a_8134101072005141=1; a_7110033113005166=1; a_7050131065005166=1; a_6131203015005242=1; a_8015067130005141=1; home_46465572=33; Hm_lvt_af8c54428f2dd7308990f5dd456fae6d=1697563736; d6b93d63cc960c878126=1697563761%2C3; Hm_lpvt_af8c54428f2dd7308990f5dd456fae6d=1697563761; a_8132032045005142=1; a_8011042075004112=1; Hm_lvt_27fe35f56bdde9c16e63129d678cd236=1697427905; a_8131077045005142=1; a_8126133045005142=1; a_7151115052005200=1; a_7143200052005200=1; a_8067043013005142=1; a_8057005035006000=1; a_7065002041006000=1; Hm_lvt_ed4f006fba260fb55ee1dfcb3e754e1c=1698288918; 94ca48fd8a42333b_code_getgraphcode=1698293099%2C1; max_u_token=a604eafb185b2da9efd9e0d6e44243de; a_7016020032006001=1; CLIENT_SYS_UN_ID=3rvhcmU98Ma5RVqmh6UoAg==; PHPSESSID=dkco0tmnr8ui540b5crck40r97; a_8061121022005136=1; TRANSFORM_USER_CHECK_AGREEMENT=read; a_5132302220002111=1; __tins__21789007=%7B%22sid%22%3A%201698656820106%2C%20%22vd%22%3A%201%2C%20%22expires%22%3A%201698658620106%7D; __tins__21784937=%7B%22sid%22%3A%201698656820168%2C%20%22vd%22%3A%201%2C%20%22expires%22%3A%201698658620168%7D; Hm_lpvt_27fe35f56bdde9c16e63129d678cd236=1698656820; a_5223124112004104=1; detail_show_similar=0; a_8037041075006001=1; PREVIEWHISTORYPAGES=601613133_2,231604277_3,600294702_2,598115535_3,598378098_1,597881355_1,597473571_2,597595624_4,213161945_2,407555420_3,269803840_5,596777297_1,597072420_2,597072455_2,579237318_1,596523778_1,594812159_1,593603490_4,593603589_2,591678290_1,591678364_1,591599534_2,582244874_1,591182727_2,314942124_2,590912365_1,590365406_3,589789397_1,589518621_1,589257506_1,588999059_2,588200974_1,587612119_1,587307168_1,587014166_1,586682671_2,586683213_3,586385929_2,586079028_1,585440196_1,585141011_2,531894818_1,529894158_1,583674790_1,583987098_1,583987123_2,583390768_3,582846882_5,374483473_1,581787745_1; a_7134062136006001=1; s_v=cdh%3D%3Ec865f4f0%7C%7C%7Cvid%3D%3E1695983286868757343%7C%7C%7Cfsts%3D%3E1695983286%7C%7C%7Cdsfs%3D%3E32%7C%7C%7Cnps%3D%3E34; s_rfd=cdh%3D%3Ec865f4f0%7C%7C%7Ctrd%3D%3Emax.book118.com%7C%7C%7Cftrd%3D%3Ebook118.com; a_8110077114006001=1; a_8112121114006001=1; Hm_lpvt_ed4f006fba260fb55ee1dfcb3e754e1c=1698769456; s_m=741760905%3D%3Esimilar%7C%7C%7C895981424%3D%3Esimilar%7C%7C%7C1958952301%3D%3Esimilar%7C%7C%7Ccdh%3D%3Ec865f4f0%7C%7C%7C-756804232%3D%3Esimilar%7C%7C%7C-2018221611%3D%3Esimilar%7C%7C%7C-316425562%3D%3Esimilar; CRM_DETAIL_INFOS=[{\"aid\":6114225221004215,\"title\":\"2020ä¸‹å\u008DŠå¹´é™•è¥¿æ•™å¸ˆèµ„æ ¼é«˜ä¸­éŸ³ä¹\u0090å­¦ç§‘çŸ¥è¯†ä¸Žæ•™å­¦èƒ½åŠ›çœŸé¢˜å\u008FŠç­”æ¡ˆ.doc\",\"firstType\":\"669\",\"secondType\":\"674\"},{\"aid\":8112121114006001,\"title\":\"2020ä¸‹å\u008DŠå¹´é\u009D’æµ·æ•™å¸ˆèµ„æ ¼é«˜ä¸­éŸ³ä¹\u0090å­¦ç§‘çŸ¥è¯†ä¸Žæ•™å­¦èƒ½åŠ›çœŸé¢˜å\u008FŠç­”æ¡ˆ.pdf\",\"firstType\":\"622\",\"secondType\":\"637\"},{\"aid\":8110077114006001,\"title\":\"2020ä¸‹å\u008DŠå¹´é™•è¥¿æ•™å¸ˆèµ„æ ¼é«˜ä¸­ç‰©ç\u0090†å­¦ç§‘çŸ¥è¯†ä¸Žæ•™å­¦èƒ½åŠ›çœŸé¢˜å\u008FŠç­”æ¡ˆ.pdf\",\"firstType\":\"669\",\"secondType\":\"674\"}]; a_6114225221004215=1; s_s=cdh%3D%3Ec865f4f0%7C%7C%7Clast_req%3D%3E1698769469%7C%7C%7Csid%3D%3E1698769445568927539%7C%7C%7Cdsps%3D%3E1; 94ca48fd8a42333b=1698769468%2C4; c4da14928424747de8b677208095de01=1698887611%2C2; operation_user_center=1; __tins__21784547=%7B%22sid%22%3A%201698941009529%2C%20%22vd%22%3A%202%2C%20%22expires%22%3A%201698943090651%7D; __51laig__=248; Hm_lpvt_f32e81852cb54f29133561587adb93c1=1698941291"

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

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
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
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/专升本考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/初中一年级",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/中考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/小升初",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/成人高考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/考研",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/自考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/高中会考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2010-2011）.51test.net/高考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/初中一年级",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/中考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/小升初",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/成人高考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/考研",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/高中会考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/高考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/专升本考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/中考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/小升初",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/成人高考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/考研",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/自考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/高中会考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/高考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/专升本考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/中考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/小升初",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/成人高考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/考研",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/自考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/高中会考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/高考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/专升本考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/中考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/小升初",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/成人高考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/考研",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/自考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/高中会考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/高考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/专升本考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/中考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/小升初",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/成人高考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/考研",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/自考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/高中会考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/高考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/专升本考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/中考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/小升初",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/成人高考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/考研",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/自考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/高中会考",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/高考",
			price:   "2000",
		},
		{
			dirName: "finish.topedu.ybep.com.cn/中考真题",
			price:   "2000",
		},
		{
			dirName: "finish.topedu.ybep.com.cn/高考真题",
			price:   "2000",
		},
		{
			dirName: "finish.www.shijuan1.com/中考试卷",
			price:   "2000",
		},
		{
			dirName: "finish.www.shijuan1.com/高考试卷",
			price:   "2000",
		},
		{
			dirName: "docx.lvlin.baidu.com",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/二级建造师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/公务员考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/教师招聘",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/教师资格考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/事业单位招聘",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/一级建造师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/造价工程师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/证券从业资格考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2010-2011）.51test.net/注册安全工程师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/二级建造师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/公务员考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/教师招聘",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/教师资格考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/事业单位招聘",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/一级建造师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/造价工程师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/证券从业资格考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2012-2013）.51test.net/注册安全工程师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/二级建造师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/公务员考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/教师招聘",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/教师资格考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/事业单位招聘",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/消防工程师",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/一级建造师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/造价工程师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/证券从业资格考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2014-2015）.51test.net/注册安全工程师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/二级建造师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/公务员考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/教师招聘",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/教师资格考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/事业单位招聘",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/消防工程师",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/一级建造师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/造价工程师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/证券从业资格考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2016-2017）.51test.net/注册安全工程师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/二级建造师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/公务员考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/教师招聘",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/教师资格考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/事业单位招聘",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/消防工程师",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/一级建造师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/造价工程师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/证券从业资格考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2018-2019）.51test.net/注册安全工程师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/二级建造师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/公务员考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/教师招聘",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/教师资格考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/事业单位招聘",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/消防工程师",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/一级建造师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/造价工程师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/证券从业资格考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2020-2022）.51test.net/注册安全工程师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2023）.51test.net/二级建造师考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2023）.51test.net/公务员考试",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2023）.51test.net/教师招聘",
			price:   "2000",
		},
		{
			dirName: "finish.tikuvip-certification（2023）.51test.net/教师资格考试",
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
