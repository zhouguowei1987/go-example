package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	NewDocDeleteEnableHttpProxy = false
	NewDocDeleteHttpProxyUrl    = "111.225.152.186:8089"
)

func NewDocDeleteSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(NewDocDeleteHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

var NewDocDeleteCookie = "BIDUPSID=8B2C214BC9D17E56E153605938409B3E; PSTM=1672058384; BAIDUID=8B2C214BC9D17E56DBA0FA4612B6B0C8:SL=0:NR=10:FG=1; hotDocPack=1; __yjs_duid=1_6b1f64fbff8e95977f656dcf822314451684723017094; BDUSS=3lhQWlqVHdwd20wUVI4ZkdCeWdxMnVnN1hmblJ6c3F5QjU3QmQ0bHRxSUtBLWRrSUFBQUFBJCQAAAAAAAAAAAEAAADcjCMiYdbcufrOsAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAp2v2QKdr9kd; BDUSS_BFESS=3lhQWlqVHdwd20wUVI4ZkdCeWdxMnVnN1hmblJ6c3F5QjU3QmQ0bHRxSUtBLWRrSUFBQUFBJCQAAAAAAAAAAAEAAADcjCMiYdbcufrOsAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAp2v2QKdr9kd; SIGNIN_UC=70a2711cf1d3d9b1a82d2f87d633bd8a04408683111x0F1M1l80pEF32Fh89wHsqGS%2Fy59%2FaM64G%2BShCzUfqSyP62oN4GLTBQD2uCSqH2bwcJeCSf7cipgJlUufvelPnxvd8ryLFfyCrRRE6zI5azjPxNaQMXUc0aYecLCG%2B2JnPYBY2uBFXlgzmAhmODMrlfysZOivcq%2FDJGToxtKZj1RN0aG3UlHNBWbSXS59DLOrxHDu5veXaaCPYl%2FAH0xuOVPYh1U42XsCzclCFyFJfambEJ95zvesH2D7g7SS65nv7hKU%2FtGqGrDH78tK5WyfLJxJ7IUYKRyYyuTIx48LILyFhx8b9iOEb9YogbFHHcV03266907366597309466153115985604; H_WISE_SIDS=213352_214793_110085_244714_261716_236312_256419_265615_265881_266368_259033_268030_265986_259642_269235_256154_269731_269780_268237_269328_269905_270084_267066_256739_270442_270460_270318_270548_271020_271169_271177_269771_271226_267659_269297_271322_265032_271271_271266_271690_270102_271882_271675_271813_271939_271954_269564_269666_234296_234207_179347_272280_266566_267596_272366_272462_272507_253022_271688_272611_272822_272816_272802_260335_271284_273066_273095_273147_273265_273230_273301_273370_273400_273396_271158_270055_273520_273603_273199_271562_272475_271147_273670_273704_272765_264170_270186_273735_263619_273165_273921_273931_274140_274158_269609_273788_273044_273594_263750_272855_274324_274279_274356_272319_197096_274411_272562; H_WISE_SIDS_BFESS=213352_214793_110085_244714_261716_236312_256419_265615_265881_266368_259033_268030_265986_259642_269235_256154_269731_269780_268237_269328_269905_270084_267066_256739_270442_270460_270318_270548_271020_271169_271177_269771_271226_267659_269297_271322_265032_271271_271266_271690_270102_271882_271675_271813_271939_271954_269564_269666_234296_234207_179347_272280_266566_267596_272366_272462_272507_253022_271688_272611_272822_272816_272802_260335_271284_273066_273095_273147_273265_273230_273301_273370_273400_273396_271158_270055_273520_273603_273199_271562_272475_271147_273670_273704_272765_264170_270186_273735_263619_273165_273921_273931_274140_274158_269609_273788_273044_273594_263750_272855_274324_274279_274356_272319_197096_274411_272562; MCITY=-289%3A; BDSFRCVID=blKOJeC62RkRrHnq5fHR2QqtOx3jE2TTH6q2Q8MlZ4GED-PexVwSEG0PgM8g0KAMRS_FogKK0eOTHkCF_2uxOjjg8UtVJeC6EG0Ptf8g0x5; H_BDCLCKID_SF=tRk8oDLhJIvDqTrP-trf5DCShUFsXUujB2Q-XPoO3KJabKQkKfv1jqDq-NjZWjJK05cObxbgy4op8P3y0bb2DUA1y4vp0f0DbgTxoUJ2bUFMO-QqqtnW-U4ebPRi3tQ9QgbMalQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0hD0wD6Daj5PVKgTa54cbb4o2WbCQ2qcT8pcN2b5oQT8LeajXbx77XbrB3KQFBDOHEPQ90qOUWJDkXpJvQnJjt2JxaqRC3q-afh5jDh3M5hT3yhjCe4ROK2Oy0hvc0J5cShnTyfjrDRLbXU6BK5vPbNcZ0l8K3l02V-bIe-t2XjQhDHt8JjKDJn3aQ5rtKRTffjrnhPF32hFPXP6-hnjy3bRDVIQ8-lR1sMjF-TjkX6DUypOaQl3Ry6r42-39LPO2hpRjyxv4Qx3L5-oxJpOJ3nkHXt5cHR7WHp7vbURvL4Lg3-7MBx5dtjTO2bc_5KnlfMQ_bf--QfbQ0hOhqP-jBRIEoKtytI8KMILr24rSMt_eqxbyJ6kOHD7yWCvtJIJcOR5Jj65b3f_e3GJDBhvl3mjyhPTwtUJ2jb4G3MA--t415xOeKTjR3j7xBfTs5pvosq0x0-Tte-bQyPbaqxCtBDOMahkM5l7xObvJQlPK5JkgMx6MqpQJQeQ-5KQN3KJmfbL9bT3tjjT3eHAqJ6_Htb3fL-08KJTHDn6Nh4Jb-tCsqxby26naag39aJ5nJDoCjUo4y-TajtKn0G5IbR39Mn7BKComQpP-qDoheMcSW601Qt8LQPQOQjvgKl0MLn3Ybb0xyUQY0ltmQxnMBMPj5mOnanvn3fAKftnOM46JehL3346-35543bRTLnLy5KJWMDcnK4-XD6oyjG3P; delPer=0; PSINO=5; BDSFRCVID_BFESS=blKOJeC62RkRrHnq5fHR2QqtOx3jE2TTH6q2Q8MlZ4GED-PexVwSEG0PgM8g0KAMRS_FogKK0eOTHkCF_2uxOjjg8UtVJeC6EG0Ptf8g0x5; H_BDCLCKID_SF_BFESS=tRk8oDLhJIvDqTrP-trf5DCShUFsXUujB2Q-XPoO3KJabKQkKfv1jqDq-NjZWjJK05cObxbgy4op8P3y0bb2DUA1y4vp0f0DbgTxoUJ2bUFMO-QqqtnW-U4ebPRi3tQ9QgbMalQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0hD0wD6Daj5PVKgTa54cbb4o2WbCQ2qcT8pcN2b5oQT8LeajXbx77XbrB3KQFBDOHEPQ90qOUWJDkXpJvQnJjt2JxaqRC3q-afh5jDh3M5hT3yhjCe4ROK2Oy0hvc0J5cShnTyfjrDRLbXU6BK5vPbNcZ0l8K3l02V-bIe-t2XjQhDHt8JjKDJn3aQ5rtKRTffjrnhPF32hFPXP6-hnjy3bRDVIQ8-lR1sMjF-TjkX6DUypOaQl3Ry6r42-39LPO2hpRjyxv4Qx3L5-oxJpOJ3nkHXt5cHR7WHp7vbURvL4Lg3-7MBx5dtjTO2bc_5KnlfMQ_bf--QfbQ0hOhqP-jBRIEoKtytI8KMILr24rSMt_eqxbyJ6kOHD7yWCvtJIJcOR5Jj65b3f_e3GJDBhvl3mjyhPTwtUJ2jb4G3MA--t415xOeKTjR3j7xBfTs5pvosq0x0-Tte-bQyPbaqxCtBDOMahkM5l7xObvJQlPK5JkgMx6MqpQJQeQ-5KQN3KJmfbL9bT3tjjT3eHAqJ6_Htb3fL-08KJTHDn6Nh4Jb-tCsqxby26naag39aJ5nJDoCjUo4y-TajtKn0G5IbR39Mn7BKComQpP-qDoheMcSW601Qt8LQPQOQjvgKl0MLn3Ybb0xyUQY0ltmQxnMBMPj5mOnanvn3fAKftnOM46JehL3346-35543bRTLnLy5KJWMDcnK4-XD6oyjG3P; BAIDUID_BFESS=8B2C214BC9D17E56DBA0FA4612B6B0C8:SL=0:NR=10:FG=1; ZFY=jk0EPTVb56DBYRiCWgTD1r1uOziyCoYjV774gIX1uRA:C; BDRCVFR[feWj1Vr5u3D]=I67x6TjHwwYf0; H_PS_PSSID=39310_39398_39396_39419_39415_39438_39433_39480_39307_39233_39403_39487_26350_39426; ab_sr=1.0.1_Y2UzZTZlNDVkOTc5ZTg0MjU1N2ZlMmUwNTBlN2JhYjViNzA5OWUwMjQ3MWNkNmEyN2E0YzgxNGJlN2UxNWVmYWViMGJmOWIzMjI5NWNiNjVhMTgzNjk0ODlmODY3NmYyZDhkNzkyMGI5YWVmYTU3ZjRiNTFkMzI3MTQzYjRjMzUzNDZlYmY3NWViNmFhY2Y4NDY0MDFkY2FhMjJkMzg0MmMzMzAyYTcwNTgxYmQxZDNmNDQ3ODExOWQ3NzliNTlh"

type GetListResponse struct {
	Data   GetListResponseData   `json:"data"`
	Status GetListResponseStatus `json:"status"`
}
type GetListResponseData struct {
	Token      string                       `json:"token"`
	DocList    []GetListResponseDataDocList `json:"doc_list"`
	TotalCount int                          `json:"total_count"`
}

type GetListResponseDataDocList struct {
	DocId      string `json:"doc_id"`
	CreateTime string `json:"create_time"`
	DocStatus  int    `json:"doc_status"`
	Title      string `json:"title"`
}

type GetListResponseStatus struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// ychEduSpider 删除未通过审核的文档
// @Title 删除未通过审核的文档
// @Description https://cuttlefish.baidu.com/，删除未通过审核的文档
func main() {
	NextDocDeleteSleep := 6
	pn := 0
	rn := 10
	isPageListGo := true
	for isPageListGo {
		hasDeleteFlag := false
		requestUrl := fmt.Sprintf("https://cuttlefish.baidu.com/nshop/doc/getlist?sub_tab=1&pn=%d&rn=%d&query=&doc_id_str=&time_range=&buyout_show_type=1&needDayUploadUserCount=1", pn, rn)
		fmt.Println(requestUrl)
		getListResponse, err := GetList(requestUrl)
		if err != nil {
			fmt.Println(err)
			break
		}
		if getListResponse.Status.Code == 0 && len(getListResponse.Data.DocList) > 0 {
			token := getListResponse.Data.Token
			fmt.Println("token：", token)
			for _, doc := range getListResponse.Data.DocList {
				fmt.Println("=======当前页为：" + strconv.Itoa(pn) + "========")
				title := doc.Title
				fmt.Println(title)

				currentTime := time.Now()
				oldTime := currentTime.AddDate(0, 0, -20)
				oldTimeStr := oldTime.Format("2006-01-02")

				// 文档状态为4可以删除
				if doc.DocStatus == 4 || (doc.DocStatus == 1 && doc.CreateTime <= oldTimeStr) {
					docIdStr := doc.DocId
					fmt.Println("=======开始删除" + strconv.Itoa(pn) + "========")
					docDeleteUrl := fmt.Sprintf("https://cuttlefish.baidu.com/user/submit/newdocdelete?token=%s&new_token=%s&fold_id_str=0&doc_id_str=%s&skip_fold_validate=1", token, token, docIdStr)
					newDocDeleteResponse, err := NewDocDelete(docDeleteUrl)
					if err == nil && newDocDeleteResponse.ErrorNo == "0" {
						hasDeleteFlag = true
						fmt.Println("=======删除成功========")
					} else {
						fmt.Println("=======删除失败========")
					}
					for i := 1; i <= NextDocDeleteSleep; i++ {
						time.Sleep(time.Second)
						fmt.Println("===========操作结束，当前是", pn, "页，暂停", NextDocDeleteSleep, "秒，倒计时", i, "秒===========")
					}
				}
			}
		}
		// 如果当前页没有任何文档删除，则请求下一页
		if hasDeleteFlag == false {
			pn++
			if pn > (getListResponse.Data.TotalCount/rn)+1 {
				fmt.Println("没有更多分页了")
				isPageListGo = false
				pn = 1
				break
			}
		}
		time.Sleep(time.Second)
	}
}

func GetList(requestUrl string) (getListResponse GetListResponse, err error) {
	// 初始化客户端
	var client *http.Client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*3)
				if err != nil {
					fmt.Println("dail timeout", err)
					return nil, err
				}
				return c, nil

			},
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 3,
		},
	}
	if NewDocDeleteEnableHttpProxy {
		client = NewDocDeleteSetHttpProxy()
	}
	getListResponse = GetListResponse{}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return getListResponse, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", NewDocDeleteCookie)
	req.Header.Set("Host", "cuttlefish.baidu.com")
	req.Header.Set("Origin", "https://cuttlefish.baidu.com/")
	req.Header.Set("Referer", "https://cuttlefish.baidu.com/shopmis?_wkts_=1697418873962")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return getListResponse, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return getListResponse, err
	}
	err = json.Unmarshal(respBytes, &getListResponse)
	if err != nil {
		return getListResponse, err
	}
	return getListResponse, nil
}

type NewDocDeleteResponse struct {
	ErrorNo string `json:"error_no"`
}

func NewDocDelete(docDeleteUrl string) (newDocDeleteResponse NewDocDeleteResponse, err error) {
	// 初始化客户端
	var client *http.Client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*3)
				if err != nil {
					fmt.Println("dail timeout", err)
					return nil, err
				}
				return c, nil

			},
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 3,
		},
	}
	if NewDocDeleteEnableHttpProxy {
		client = NewDocDeleteSetHttpProxy()
	}

	newDocDeleteResponse = NewDocDeleteResponse{}
	req, err := http.NewRequest("GET", docDeleteUrl, nil) //建立连接
	if err != nil {
		return newDocDeleteResponse, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", NewDocDeleteCookie)
	req.Header.Set("Host", "cuttlefish.baidu.com")
	req.Header.Set("Referer", "https://cuttlefish.baidu.com/shopmis?_wkts_=1697418873962")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return newDocDeleteResponse, err
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return
		}
	} else {
		reader = resp.Body
	}
	respBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return newDocDeleteResponse, err
	}
	err = json.Unmarshal(respBytes, &newDocDeleteResponse)
	if err != nil {
		return newDocDeleteResponse, err
	}
	return newDocDeleteResponse, nil
}
