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

// NewDocDeleteCookie 15238369929账号
var NewDocDeleteCookie = "BDUSS=JXc3R6TnFFd2dRcUp0dFBkd2d2V2RiU1BoYUc2M2dtdzFzV1hRQThiYXFOR3hsSVFBQUFBJCQAAAAAAAAAAAEAAADcjCMiYdbcufrOsAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKqnRGWqp0RlUX; BDUSS_BFESS=JXc3R6TnFFd2dRcUp0dFBkd2d2V2RiU1BoYUc2M2dtdzFzV1hRQThiYXFOR3hsSVFBQUFBJCQAAAAAAAAAAAEAAADcjCMiYdbcufrOsAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKqnRGWqp0RlUX; PSTM=1698999470; BIDUPSID=8B2C214BC9D17E56E153605938409B3E; BAIDUID=12BC5983BECCE097EB0A3B596F658A2C:SL=0:NR=10:FG=1; MCITY=-289%3A; H_WISE_SIDS=40211_40080_40364_40352_40373_40368_40415_40299_40467_40317_40506_40500_40511_40514; H_WISE_SIDS_BFESS=40211_40080_40364_40352_40373_40368_40415_40299_40467_40317_40506_40500_40511_40514; BDSFRCVID=KJ0OJexroG3KxNJt20Agjcryt257H43TDYrEOwXPsp3LGJLVY5VYEG0Pt_gaCz--oxcOogKKBeOTHgIF_2uxOjjg8UtVJeC6EG0Ptf8g0M5; H_BDCLCKID_SF=tRk8oDLhJIvDqTrP-trf5DCShUFs0hFOB2Q-XPoO3KOvJhr_555shTLq-p58-TDtJ2-JoMbgy4op8P3y0bb2DUA1y4vpX6QpymTxoUJ2Xb3HEPFzqtnW-U4ebPRi3tQ9QgbMalQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0hD0wD6Daj5PVKgTa54cbb4o2WbCQLpnU8pcN2b5oQTtU2-TZ-Tj3ab7B3-oXbI_hMpPmDpOUWJDkXpJvQnJjt2JxaqRCyC54oq5jDh3MBP-_LpCLe4ROK2Oy0hvc0J5cShnTyfjrDRLbXU6BK5vPbNcZ0l8K3l02V-bIe-t2XjQhDHt8JjKDJn3aQ5rtKRTffjrnhPF35xLdXP6-hnjy3b7iXprFbIosDn6F-n5B55j0jPAJJh3Ry6r42-39LPO2hpRjyxv4WjLZWtoxJpOJXKvK5-FEHR7WHp7vbURvL4Lg3-7MBx5dtjTO2bc_5KnlfMQ_bf--QfbQ0hOhqP-jBRIEoC0XtDLaMKvPKITD-tFO5eT22-us5C6i2hcHMPoosIOo2MTk55DJ0P5B5bj8bnut_hP2tMbUotoHXnJi0btQDPvxBf7pyRTphl5TtUJM_n3pLnCVqt4b-PnyKMni0DT9-pnjWlQrh459XP68bTkA5bjZKxtq3mkjbPbDfn028DKuDjtBD65LjaRabK6aKC5bL6rJabC3OpFxXU6q2bDeQNbihU6Zbnuj5R3_K-oBfK_GjUAMDp0vWq54WbbvLT7johRTWqR4oD5KWfonDh83BUJK3bILHCOO3h7O5hvvER5O3M7JyfKmDloOW-TB5bbPLUQF5l8-sq0x0bOte-bQXH_E5bj2qRCj_K_b3j; ZFY=Kvh3mlTDq08kcM6gCVjGPdJ8od8Jh1tm70qFFB33tSg:C; BAIDUID_BFESS=12BC5983BECCE097EB0A3B596F658A2C:SL=0:NR=10:FG=1; BDSFRCVID_BFESS=KJ0OJexroG3KxNJt20Agjcryt257H43TDYrEOwXPsp3LGJLVY5VYEG0Pt_gaCz--oxcOogKKBeOTHgIF_2uxOjjg8UtVJeC6EG0Ptf8g0M5; H_BDCLCKID_SF_BFESS=tRk8oDLhJIvDqTrP-trf5DCShUFs0hFOB2Q-XPoO3KOvJhr_555shTLq-p58-TDtJ2-JoMbgy4op8P3y0bb2DUA1y4vpX6QpymTxoUJ2Xb3HEPFzqtnW-U4ebPRi3tQ9QgbMalQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0hD0wD6Daj5PVKgTa54cbb4o2WbCQLpnU8pcN2b5oQTtU2-TZ-Tj3ab7B3-oXbI_hMpPmDpOUWJDkXpJvQnJjt2JxaqRCyC54oq5jDh3MBP-_LpCLe4ROK2Oy0hvc0J5cShnTyfjrDRLbXU6BK5vPbNcZ0l8K3l02V-bIe-t2XjQhDHt8JjKDJn3aQ5rtKRTffjrnhPF35xLdXP6-hnjy3b7iXprFbIosDn6F-n5B55j0jPAJJh3Ry6r42-39LPO2hpRjyxv4WjLZWtoxJpOJXKvK5-FEHR7WHp7vbURvL4Lg3-7MBx5dtjTO2bc_5KnlfMQ_bf--QfbQ0hOhqP-jBRIEoC0XtDLaMKvPKITD-tFO5eT22-us5C6i2hcHMPoosIOo2MTk55DJ0P5B5bj8bnut_hP2tMbUotoHXnJi0btQDPvxBf7pyRTphl5TtUJM_n3pLnCVqt4b-PnyKMni0DT9-pnjWlQrh459XP68bTkA5bjZKxtq3mkjbPbDfn028DKuDjtBD65LjaRabK6aKC5bL6rJabC3OpFxXU6q2bDeQNbihU6Zbnuj5R3_K-oBfK_GjUAMDp0vWq54WbbvLT7johRTWqR4oD5KWfonDh83BUJK3bILHCOO3h7O5hvvER5O3M7JyfKmDloOW-TB5bbPLUQF5l8-sq0x0bOte-bQXH_E5bj2qRCj_K_b3j; PSINO=2; delPer=0; __bid_n=18e37f202126128d0122e4; H_PS_PSSID=40299_40500_40080_60127_60138_60237; BDRCVFR[feWj1Vr5u3D]=I67x6TjHwwYf0; ab_sr=1.0.1_YmQ3YjAzMTFhZWE4ZjU5MjIyOTgwYWEwMGZiZWE5MjBhNmIyYTg1YjJiMmRkMjk3NDFkOWRlZDc3OTA1ZGJhYzk5MjI3Yjk3ZWI1YTY1MWEwOGNmY2IwZDdlZTJjMGE1ZGUyYWQxZGMwMjU1OGU1MDFhOWMwZTUzZWNiZjExY2NhNmRmYjJmYzUxMGE0NzhlMmYwYmQ3Y2Q1ZjM1MjA1MTk4Zjg5MzhlNzQ4ZWEzNGEwODI3NTZmZDRhMGUxMTIw"

// NewDocDeleteCookie 15803889687账号
//var NewDocDeleteCookie = "ab_sr=1.0.1_MGI2NjQ2NWY0NWIyNjEyMDZlZmVjOTQwOTY1OTA0ZWVkYmE3MWQyOTMyNjIxNzBjMWNlMWQxZmE2NTMxMDlhODU4Nzk5MzZhMzBlMzY2YTMzNjcyNDQwMWQxN2NmY2MxYjkyMjBlNTA4M2JiZGFlYjlmODBlZWVjN2MyNWM5ZTgxMWUxMzY0YmUyMmE5M2JmMWIxODFmMWM4ZWNmMGU5OTFiODFhYjg0MGYyMTllN2M4NGIyYTczYjMzNWY5NjM2; __bid_n=18f3c1a89884368d55ef7c; H_PS_PSSID=40010_39661_40207_40211_40217_40079_40365_40351_40304_40380_40368_40409_40415; BDUSS=lk3RTVyZi1jMzRDZ2NtYm1Nb3hhTDhLU2VDekMtb21xdGdOMjJzNFB4aDZiT1JsRVFBQUFBJCQAAAAAAAAAAAEAAAA80dNV1YzYmcDWytPD1gAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAHrfvGV637xlc; BIDUPSID=3BC54111662B04D6D3A343B0F6A2D2F9; PSTM=1702619109; BAIDUID=3BC54111662B04D6D3A343B0F6A2D2F9:FG=1"

// NewDocDeleteCookie 17539389687账号
//var NewDocDeleteCookie = "BIDUPSID=EACB60FE57432A4F592404B25D1A711F; PSTM=1614481586; __bid_n=185ec8d2a4a437e3294207; BAIDU_WISE_UID=wapp_1674718857532_689; FPTOKEN=sbPlQcNZbeFCbEYfdOrnJ6qcpID/06laPdXTjFrXDglLBJWQl30/+Vh86S1EJJGZykln4TEv76gRNguLer7/s8fjYd1niHQEUKngXw6X1pQ7RyNx5+dEADpLHxfkC32sRqL8Wn2aI8FluJ0NeZSKGCf/EeczG52bjrBDv5729Pl/WsOJSE7vKUF/5YEd1ES/8Zh71zWJpFtgxgfLjBvcLGf04w2GeR9xWToIPwPeZZ3xqbXD3welXNNjU5X2GR+7qG+orF9B23b8jQKdCeTMbdBHUJHqnWRNLrYPi28SY7CDEWGDuNjKrFn/OEWuP6PZ2ORQnyCnyAeAm6G/ijf2NokfIoT4l7WNTw8GnYbBzzYQ6H+RcszkplXU1qe71ip3adoGhbH5TGWaPnhC80XMxw==|bJtkSU/lYu1X41w73/OlDop2e9GzmK5gnjMunxaVfYo=|10|ca2fe9c56e49521b8d38c53cab4ce8b6; BAIDUID=340EBDC3BF189A84744F943C28903CCC:FG=1; H_PS_PSSID=60337_60375; BAIDUID_BFESS=340EBDC3BF189A84744F943C28903CCC:FG=1; ZFY=QGk8i9qm6WMpBqGQ5XH9RRpe44mo7tuvBzixH2xINFw:C; BDUSS_BFESS=2paTW1BRjBCYWcwT0JXZWNqSTh4QXkwcTc4RU9-OXhodGNSTy1ydVB5WG1acXRtRUFBQUFBJCQAAAAAAQAAAAEAAAAv1bKNTWlzs6zIy7K7u-G3yQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAObZg2bm2YNmT; RT=\"z=1&dm=baidu.com&si=d5db28c7-6ae3-4521-8424-7acd87ac8697&ss=ly4a77qt&sl=b&tt=tu0&bcn=https%3A%2F%2Ffclog.baidu.com%2Flog%2Fweirwood%3Ftype%3Dperf&ld=5zh2&ul=aqfr&hd=aqg3\"; ab_sr=1.0.1_M2ViMmE4ZjhhZjQzYTYxYzllMjVmZTdjNDU1MzZhYWUyODIyYWIwM2RkOWNiOTFkMDAxYjI5ZmU3N2JkZWNmOTI0NGZjOTU0ODk2MjJkYjYyYzBiNTFiNzg5ZWIxMDM2Y2JkMjlkOTBkMzhhMjYzZjhjY2FhNjFhOTA4MWE0NTcxYjdkNTNiODY0YzYwNDAwOTFmMmM0NmJkZTZkY2Y1YWFiYjczYWFiZjRjMjZiY2FhOGM2NDhjNmNiNWZkYzIy"

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
				oldTime := currentTime.AddDate(0, 0, -40)
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
