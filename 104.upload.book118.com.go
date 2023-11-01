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

var SessionId = "f1tq9sahmckd7sr8ckht2glbe5"
var Token = "75efa22fb33bed7f455f09b0b160837e"
var Cookie = "__yjs_duid=1_1543d26121978a9cfb0ca147de19aa051678550479017; a_8120007071005116=1; a_6010135035005211=1; a_8024001040005117=1; a_6050221154005154=1; EXAMINE_CLOSE=true; a_7025115044005142=1; a_5034330112010304=1; a_5034204112010304=1; a_7166015111005142=1; a_8140135071005117=1; a_7125044014005143=1; a_8074113131005117=1; a_5222111324010304=1; a_7116014155005142=1; a_6140230225005211=1; a_6232002143005212=1; a_7014001035005144=1; a_5033113143010311=1; a_7025066066005144=1; a_5034204143010311=1; a_5034033143010311=1; a_7140063141005144=1; a_8037007005005122=1; a_8101114042004135=1; a_7043001005005145=1; a_7033115025004150=1; home_4081071=1; a_5242134134010234=1; a_8102010045005003=1; a_8051124045005122=1; a_8137032102005122=1; a_7046133066003134=1; a_5300033303010311=1; a_6125013015005215=1; 0ab7d33081eeb=%E5%B9%B6%E7%BD%91%E8%B0%83%E5%BA%A6%E5%8D%8F%E8%AE%AE%E7%A4%BA%E8%8C%83%E6%96%87%E6%9C%AC; input_search_logs=a%3A4%3A%7Bi%3A0%3Ba%3A2%3A%7Bs%3A8%3A%22keywords%22%3Bs%3A30%3A%22%E5%B9%B6%E7%BD%91%E8%B0%83%E5%BA%A6%E5%8D%8F%E8%AE%AE%E7%A4%BA%E8%8C%83%E6%96%87%E6%9C%AC%22%3Bs%3A4%3A%22time%22%3Bi%3A1692713602%3B%7Di%3A1%3Ba%3A2%3A%7Bs%3A8%3A%22keywords%22%3Bs%3A46%3A%222023%E5%B9%B4%E4%BA%91%E5%8D%97%E5%A4%A7%E7%90%86%E4%B8%AD%E8%80%83%E8%8B%B1%E8%AF%AD%E8%AF%95%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88%22%3Bs%3A4%3A%22time%22%3Bi%3A1691545978%3B%7Di%3A2%3Ba%3A2%3A%7Bs%3A8%3A%22keywords%22%3Bs%3A46%3A%222023%E5%B9%B4%E4%BA%91%E5%8D%97%E5%BE%B7%E5%AE%8F%E4%B8%AD%E8%80%83%E6%95%B0%E5%AD%A6%E8%AF%95%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88%22%3Bs%3A4%3A%22time%22%3Bi%3A1691545915%3B%7Di%3A3%3Ba%3A2%3A%7Bs%3A8%3A%22keywords%22%3Bs%3A46%3A%222023%E5%B9%B4%E4%BA%91%E5%8D%97%E5%BE%B7%E5%AE%8F%E4%B8%AD%E8%80%83%E8%8B%B1%E8%AF%AD%E8%AF%95%E9%A2%98%E5%8F%8A%E7%AD%94%E6%A1%88%22%3Bs%3A4%3A%22time%22%3Bi%3A1691545748%3B%7D%7D; d6b93d63cc960c878126=1692713602%2C1; a_8137034105005064=1; a_8064041013005123=1; a_8007104047005123=1; a_7010143054005146=1; a_8037076056005015=1; a_6155035242005215=1; a_7130200200005146=1; a_5040100103010314=1; a_8024011034005124=1; a_8057132103005123=1; a_5312314202010240=1; a_5332144320010314=1; a_8134140125005124=1; a_8051072131005035=1; home_25896707=1; a_8060022131005037=1; a_6021112032005213=1; home_4926016=2; home_46465572=32; a_8131114125005124=1; a_8002135131005035=1; home_21887109=1; PHPSESSID=f1tq9sahmckd7sr8ckht2glbe5; TRANSFORM_USER_CHECK_AGREEMENT=read; a_5332334320010314=1; a_5020110024010320=1; a_7033114052005130=1; a_6014015022005221=1; a_5310302202010240=1; a_8001140054005125=1; a_5001302134010320=1; a_8002003054005125=1; a_6142142045001243=1; home_4070593=1; a_5330103012010321=1; a_5330001012010321=1; return_url=http%3A%2F%2Fmax.book118.com%2Fuser_center_v1%2Fdoc%2Findex%2Findex.html; a_6135045102005222=1; a_8040015104005126=1; a_8037073104005126=1; a_8032107104005126=1; a_8051102001005127=1; a_7134005201004040=1; a_5241233110010322=1; a_8025023075005127=1; __51cke__=; a_5041234221010322=1; a_8002065133005127=1; a_7012134026005154=1; a_7012110026005154=1; a_6012210032005224=1; a_7011200026005154=1; Hm_lvt_f32e81852cb54f29133561587adb93c1=1694227655; a_6230135243005224=1; a_8132054143005130=1; a_8113006031005131=1; a_8124113054004013=1; a_5321041201010324=1; a_8135141116005131=1; a_7163061141005155=1; a_5100310012010330=1; a_8066006044005132=1; a_7151161055003050=1; a_5343021240003110=1; Hm_lvt_27fe35f56bdde9c16e63129d678cd236=1695001548; home_15265435=1; d6b93d4rgc960c878126=1695001563%2C1; __tins__21789007=%7B%22sid%22%3A%201695001547274%2C%20%22vd%22%3A%202%2C%20%22expires%22%3A%201695003371585%7D; a_6151204214002110=1; detail_show_similar=0; a_8027101133005132=1; a_7042124016002155=1; __tins__21784937=%7B%22sid%22%3A%201695025922409%2C%20%22vd%22%3A%201%2C%20%22expires%22%3A%201695027722409%7D; Hm_lpvt_27fe35f56bdde9c16e63129d678cd236=1695025923; a_7030033163003020=1; a_6043234030005231=1; reward_download_aid=591182794; tongji_46465572=1; Hm_lvt_b645044a3b9e8b6315c6fe7d4733b16c=1693234255,1695106630; Hm_lpvt_b645044a3b9e8b6315c6fe7d4733b16c=1695106794; a_8033033022005133=1; a_8017012061005074=1; a_5102314033010331=1; a_8137041073005133=1; a_117943075=1; 5a9a221b83986f79ee93b689251380af=1695175374%2C1; a_5340114214010331=1; a_6120202040005214=1; UPLOAD_AGREEMENT_CHECKED=1; Hm_lvt_ed4f006fba260fb55ee1dfcb3e754e1c=1695259307; a_8123100103005133=1; a_8122132103005133=1; a_6100033140005233=1; s_rfd=cdh%3D%3Ec865f4f0%7C%7C%7Ctrd%3D%3Emax.book118.com%7C%7C%7Cftrd%3D%3Ebook118.com; a_7051014114005162=1; a_7050155114005162=1; CLIENT_SYS_UN_ID=3rvhLGUYy0jCm+7BDSZmAg==; CRM_DETAIL_INFOS=[{\"aid\":7046156114005162,\"title\":\"2018å¹´æ²³åŒ—çœ\u0081ä¸­è€ƒæ–‡ç»¼(é\u0081“å¾·ä¸Žæ³•æ²»)è¯•å\u008D·(å\u0090«ç­”æ¡ˆ).pdf\",\"firstType\":\"622\",\"secondType\":\"631\"},{\"aid\":7050155114005162,\"title\":\"2018å¹´æµ™æ±Ÿçœ\u0081æ\u009D­å·žå¸‚ä¸­è€ƒåŒ–å­¦è¯•å\u008D·å\u008FŠè§£æž\u0090.pdf\",\"firstType\":\"622\",\"secondType\":\"631\"},{\"aid\":7051014114005162,\"title\":\"2018å¹´æµ™æ±Ÿçœ\u0081æ\u009D­å·žå¸‚ä¸­è€ƒç”Ÿç‰©è¯•å\u008D·å\u008FŠè§£æž\u0090.pdf\",\"firstType\":\"622\",\"secondType\":\"631\"}]; a_7046156114005162=1; s_v=cdh%3D%3Ec865f4f0%7C%7C%7Cvid%3D%3E1695983286868757343%7C%7C%7Cfsts%3D%3E1695983286%7C%7C%7Cdsfs%3D%3E2%7C%7C%7Cnps%3D%3E4; s_s=cdh%3D%3Ec865f4f0%7C%7C%7Clast_req%3D%3E1696172356%7C%7C%7Csid%3D%3E1696172356732712947%7C%7C%7Cdsps%3D%3E1; Hm_lpvt_ed4f006fba260fb55ee1dfcb3e754e1c=1696172356; 94ca48fd8a42333b=1696172355%2C1; PREVIEWHISTORYPAGES=593603490_4,593603589_2,591678290_1,591678364_1,591599534_2,582244874_1,591182727_2,314942124_2,590912365_1,590365406_3,589789397_1,589518621_1,589257506_1,588999059_2,588200974_1,587612119_1,587307168_1,587014166_1,586682671_2,586683213_3,586385929_2,586079028_1,585440196_1,585141011_2,531894818_1,529894158_1,583674790_1,583987098_1,583987123_2,583390768_3,582846882_5,374483473_1,581787745_1,581481918_1,581261101_1,580639202_1,579322001_2,580116832_1,579321961_1,579230659_1,579248702_2,579073158_1,579073270_2,578847592_1,578844627_3,578577974_1,570070289_1,578079653_1,577772647_3,576712400_1; c4da14928424747de8b677208095de01=1696229223%2C2; 94ca48fd8a42333b_code_getgraphcode=1696312823%2C1; max_u_token=617e85a3d2281143f809696f15e3e80f; operation_user_center=1; __tins__21784547=%7B%22sid%22%3A%201696312902807%2C%20%22vd%22%3A%201%2C%20%22expires%22%3A%201696314702807%7D; __51laig__=92; Hm_lpvt_f32e81852cb54f29133561587adb93c1=1696312903"

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
			dirName: "finish.tikuvip（2012-2013）.51test.net/中考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/小升初", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/成人高考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/考研", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/高中会考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2012-2013）.51test.net/高考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/专升本考试", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/中考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/小升初", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/成人高考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/考研", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/自考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/高中会考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2014-2015）.51test.net/高考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/专升本考试", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/中考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/小升初", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/成人高考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/考研", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/自考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/高中会考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2016-2017）.51test.net/高考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/专升本考试", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/中考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/小升初", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/成人高考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/考研", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/自考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/高中会考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2018-2019）.51test.net/高考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/专升本考试", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/中考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/小升初", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/成人高考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/考研", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/自考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/高中会考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2020-2022）.51test.net/高考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/专升本考试", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/中考", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/小升初", price: "2000",
		},
		{
			dirName: "finish.tikuvip（2023）.51test.net/成人高考", price: "2000",
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
