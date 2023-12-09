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

var SessionId = "gpr9fvacpfbcgd21mvifgfl525"
var Token = "93c33831068c0cf18ccc08098fb7fab8"
var Cookie = "__yjs_duid=1_1543d26121978a9cfb0ca147de19aa051678550479017; d6b93d4rgc960c878126=1695001563%2C1; Hm_lvt_b645044a3b9e8b6315c6fe7d4733b16c=1693234255,1695106630; 5a9a221b83986f79ee93b689251380af=1695175374%2C1; CLIENT_SYS_UN_ID=3rvhcmU98Ma5RVqmh6UoAg==; TRANSFORM_USER_CHECK_AGREEMENT=read; PHPSESSID=gpr9fvacpfbcgd21mvifgfl525; __51cke__=; Hm_lvt_f32e81852cb54f29133561587adb93c1=1698989589,1699374285; __tins__21784547=%7B%22sid%22%3A%201699594064465%2C%20%22vd%22%3A%201%2C%20%22expires%22%3A%201699595864465%7D; __51laig__=5; a_5232041132011004=1; a_7013125125006004=1; a_5014310233011004=1; a_5001334004004312=1; a_7002033162006004=1; a_7154043016006004=1; Hm_lvt_ed4f006fba260fb55ee1dfcb3e754e1c=1698288918,1699111115,1699758173; a_8077122015006004=1; a_7120010016006004=1; a_7200013161006004=1; a_8001041135006004=1; a_7120155043004144=1; a_8077064026006005=1; a_7101144113005145=1; a_6041053210002011=1; a_6112235115010005=1; home_46465572=1; Hm_lvt_af8c54428f2dd7308990f5dd456fae6d=1697563736,1699944510; a_8054034057006005=1; a_7140030006005134=1; a_6104003202005135=1; home_24253368=1; a_8072064013006005=1; a_8027037040004116=1; home_41775754=1; d6b93d63cc960c878126=1700031660%2C1; Hm_lpvt_af8c54428f2dd7308990f5dd456fae6d=1700031662; a_8077137102005055=1; a_5014041214003124=1; a_8113032113005135=1; a_8034034014005124=1; a_6213113141004223=1; a_6100003024005204=1; a_8056055116005102=1; a_5121044132003314=1; detail_show_similar=0; a_7045036006006010=1; a_7043022165003140=1; a_6111154000010010=1; a_5033041012010310=1; PREVIEWHISTORYPAGES=606004370_1,605474428_1,605226352_1,604930224_3,603618514_3,602280266_1,601613133_2,231604277_3,600294702_2,598115535_3,598378098_1,597881355_1,597473571_2,597595624_4,213161945_2,407555420_3,269803840_5,596777297_1,597072420_2,597072455_2,579237318_1,596523778_1,594812159_1,593603490_4,593603589_2,591678290_1,591678364_1,591599534_2,582244874_1,591182727_2,314942124_2,590912365_1,590365406_3,589789397_1,589518621_1,589257506_1,588999059_2,588200974_1,587612119_1,587307168_1,587014166_1,586682671_2,586683213_3,586385929_2,586079028_1,585440196_1,585141011_2,531894818_1,529894158_1,583674790_1; a_6053014010010011=1; Hm_lvt_27fe35f56bdde9c16e63129d678cd236=1700492470; Hm_lpvt_27fe35f56bdde9c16e63129d678cd236=1700492470; a_7020150044004100=1; a_8041006006006007=1; s_m=577898479%3D%3Esimilar%7C%7C%7C982069690%3D%3Esimilar%7C%7C%7C1314443394%3D%3Esimilar%7C%7C%7C1473898992%3D%3Esimilar%7C%7C%7C2105215810%3D%3Esimilar%7C%7C%7Ccdh%3D%3Ec865f4f0%7C%7C%7C-300601523%3D%3Esimilar%7C%7C%7C-1514282584%3D%3Esimilar; a_6130102115003231=1; a_8040132006006007=1; s_rfd=cdh%3D%3Ec865f4f0%7C%7C%7Ctrd%3D%3Ewww.baidu.com%7C%7C%7Cftrd%3D%3Ebaidu.com; a_8017003131005132=1; a_7104121042004051=1; a_6043105014005225=1; CRM_DETAIL_INFOS=[{\"aid\":6050202043010010,\"title\":\"éƒ¨ç¼–ç‰ˆå°\u008Få­¦è¯­æ–‡ä¸‰å¹´çº§ä¸‹å†ŒæœŸä¸­æµ‹è¯•å\u008D·-æ–‡æ¡£åœ¨çº¿é¢„è§ˆ.pptx\",\"firstType\":\"612\",\"secondType\":\"617\"},{\"aid\":6043105014005225,\"title\":\"2021å¹´å®\u0081å¾·åˆ\u009Dä¸­æ•°å­¦ç¬¬ä¸€æ¬¡è´¨æ£€æ•°å­¦ç­”æ¡ˆ.pdf\",\"firstType\":\"622\",\"secondType\":\"625\"},{\"aid\":8017003131005132,\"title\":\"2021-2022å­¦å¹´æ¹–åŒ—çœ\u0081å®œæ˜Œå¸‚äº”å³°åŽ¿ä¸ƒå¹´çº§ï¼ˆä¸Šï¼‰æœŸæœ«æ•°å­¦è¯•å\u008D·.pdf\",\"firstType\":\"622\",\"secondType\":\"625\"}]; Hm_lpvt_ed4f006fba260fb55ee1dfcb3e754e1c=1700803242; a_6050202043010010=1; s_v=cdh%3D%3Ec865f4f0%7C%7C%7Cvid%3D%3E1695983286868757343%7C%7C%7Cfsts%3D%3E1695983286%7C%7C%7Cdsfs%3D%3E56%7C%7C%7Cnps%3D%3E58; s_s=cdh%3D%3Ec865f4f0%7C%7C%7Clast_req%3D%3E1700803242%7C%7C%7Csid%3D%3E1700803242809263845%7C%7C%7Cdsps%3D%3E1; 94ca48fd8a42333b=1700803242%2C2; 30d8fb61e609cac11=1700803247%2C1; c4da14928424747de8b677208095de01=1700877391%2C2; return_url=http%3A%2F%2Fmax.book118.com%2Fuser_center_v1%2Fdoc%2Findex%2Findex.html; 94ca48fd8a42333b_code_getgraphcode=1700899316%2C1; max_u_token=e5261e1daf5f3890b4a5ba0154e7e594; operation_user_center=1; Hm_lpvt_f32e81852cb54f29133561587adb93c1=1700928324"

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
			dirName: "docx.www.shijuan1.com/中考试卷",
			price:   "2000",
		},
		{
			dirName: "docx.www.shijuan1.com/高考试卷",
			price:   "2000",
		},
		{
			dirName: "docx.www.shijuan1.com/小学试卷",
			price:   "2000",
		},
		{
			dirName: "docx.51zjedu.com/一年级",
			price:   "2000",
		},
		{
			dirName: "docx.51zjedu.com/二年级",
			price:   "2000",
		},
		{
			dirName: "docx.51zjedu.com/三年级",
			price:   "2000",
		},
		{
			dirName: "docx.51zjedu.com/四年级",
			price:   "2000",
		},
		{
			dirName: "docx.51zjedu.com/五年级",
			price:   "2000",
		},
		{
			dirName: "docx.51zjedu.com/六年级",
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
		{
			dirName: "docx.zzstep.com/小学/道德与法治",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/小学/美术",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/小学/数学",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/小学/信息技术",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/小学/音乐",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/小学/英语",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/小学/语文",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/初中/数学",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/初中/语文",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/初中/英语",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/初中/物理",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/初中/化学",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/初中/音乐",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/初中/道德与法治",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/初中/美术",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/初中/生物",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/初中/地理",
			price:   "2000",
		},
		{
			dirName: "docx.zzstep.com/初中/历史",
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
