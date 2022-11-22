package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Payload struct {
	Color        int     `json:"color"`
	Cough        int     `json:"cough"`
	Feverd       int     `json:"feverd"`
	Fuxie        int     `json:"fuxie"`
	Hsxx         int     `json:"hsxx"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Name         string  `json:"name"`
	Smell        int     `json:"smell"`
	StuId        string  `json:"stuId"`
	Temperatured int     `json:"temperatured"`
	Weakness     int     `json:"weakness"`
	Yanteng      int     `json:"yanteng"`
}

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Error   bool        `json:"error"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
}

func main() {
	var Stubs = map[string]string{
		"1561627150648348674": "周诗霖",
		"1561625329819086850": "周诗濛",
	}
	url := "http://jiaoyu.gcpy365.com:8213//api/sys/health/report/add"
	for stuId, name := range Stubs {
		payload := Payload{
			Color:        0,
			Cough:        0,
			Feverd:       0,
			Fuxie:        0,
			Hsxx:         0,
			Latitude:     31.29456,
			Longitude:    121.436295,
			Name:         name,
			Smell:        0,
			StuId:        stuId,
			Temperatured: 0,
			Weakness:     0,
			Yanteng:      0,
		}
		payloadJson, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJson))
		if err != nil {
			fmt.Println(err.Error())
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Encoding", "gzip, deflate")
		req.Header.Set("Accept-Languageg", "zh-CN,zh;q=0.9")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Content-Length", "196")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Host", "jiaoyu.gcpy365.com:8213")
		req.Header.Set("Origin", "http://jiaoyu2.gcpy365.com:8213")
		req.Header.Set("Referer", "http://jiaoyu2.gcpy365.com:8213/")
		req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")
		req.Header.Set("X-GcSoft-Token", "8657ff1f99e5ef8fb52fb4cc661bae4e")
		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(err.Error())
		}
		JoaoYuResp := Response{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&JoaoYuResp); err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(JoaoYuResp.Code, JoaoYuResp.Data, JoaoYuResp.Error, JoaoYuResp.Msg, JoaoYuResp.Success)
		time.Sleep(time.Second * 3)
	}
}
