package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "https://so2cimsfl2.execute-api.eu-central-1.amazonaws.com/prod/product-categories"
	client := &http.Client{}                     //初始化客户端
	req, err := http.NewRequest("GET", url, nil) //建立连接
	if err != nil {
		panic(err)
	}

	req.Header.Set("authority", "so2cimsfl2.execute-api.eu-central-1.amazonaws.com")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("authorization", "AWS4-HMAC-SHA256 Credential=ASIA5QCG2T6RXR6TW433/20221031/eu-central-1/execute-api/aws4_request, SignedHeaders=host;x-amz-date;x-amz-security-token, Signature=n8ba6ca7104b358397e7193f787403c286dc62e54d2dea0eb1f732f9ec3e3579c")
	req.Header.Set("if-none-match", "W/\"3f1-EwjOysAuoSj/f4/gFlyzTmjts7A\"")
	req.Header.Set("origin", "https://www.2030calculator.com")
	req.Header.Set("referer", "https://www.2030calculator.com/")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"104\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"104\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	req.Header.Set("x-amz-date", "20221031T080907Z")
	req.Header.Set("x-amz-security-token", "IQoJb3JpZ2luX2VjEIj//////////wEaDGV1LWNlbnRyYWwtMSJHMEUCIDetVe9SisyWpNd5wBiICRVOjQXmrWD+R8XnGIVicVhAAiEA9MeQxXFCsixESjv4xNxT8mZ/U3ScAx0TGQVUjHCeaTQq1gQIcRACGgw5Mjc4NjE0ODEzNzkiDFbAggrU+L7G8ewOeiqzBNGdPPrC/8FERZhaQIkLL7yfuQEzDezmP1k9zpwzX4KK6iAYp6GLGKFMahSJC29MywVTQLJTTvAD8aVlR+fzDGnrxA+Hbkfn59TZSRrb/8atiweayc0TTCBTHpOMmXXBgJYaOPP6ADEgV9eHNBOoluBPoqOGHNIa2B75PWTpMQS5OUXwr4lilyoOtWCjTF1iFRMif7wqGrRfDGfzZpR1sfo6cKzpeCvNiXkyiAwljvvARIVEioX5JCykR6StIw3tyxrZ2E92h0YGJG2xu2q8Br1fDW8mPxGrdUpNcDfyD7SD/YmcBQA+UtmRlMLeDuDU6UOq/1fcGMXn2rTF5MYgR2e40hLzB4JYUdHHRt+//HboWXetUlixuFvcA7BPuurcEdH6rCVfwjzssK64P1ScVXGBQzR3uUO20+7m8KhupFkVyaNniyNmniqfE20TfdQnE+uf1PbYH7aNP9ySbqFuxvd+aXcLoOtl1tT9hINGcGgENdYCMpRUmx1iH8NbAQA4qCy0FxZvQVXJMpP9FLEE0eWcMvvbELeWToEM57TJrGAy/I9wQCpuIurSETfGqGZPROIlUuVNKwJYHyZmqsbkkaEjve6ZgUJWqBcBWUbfoh9UKgwUPQ3QgynNJQuC0P5tkCPih4ChFzNKBWPVnxae7zOzv+JP4Xf5flk5RaUmRbV+c5ymckBWwr+tqk+q8t9Ux7j9yxmlMSqKXS7vBKAiRJcboe3p158/Vw8qS7wFAhxwmcXGML2E/poGOoUCTdqlzP7Q5M0Uva3uOeAws0Vkv+tfyrJ4lt3t7gsTwHcVWCC9s3v+Vo1RwBYwtQg3McXlyGM+RDXkbyQ/0FoUMdFBszTlM7ZARdHHmmA7iMZ8AZtjqWzcUfkNJTH5oQOnXQDbzogni2XPA5vGt3YWyuxYrb2ZcdMJCnEal/kLRZz5ZbzdMkxXDBigG5lrN3aavxWJuXEGkrVZmDZjLG53vzgBQrGK2wLHXh7X8YNLPcG2xxYVb4gyjRXlyQ5Bk/w+37w2zbj/v2wq6vXx8zmnDFmkgOHbwfkc1vHdJqbr+KkcyVuIbQQygjRWneTm9JBngr/QYOGkSxLjYGs35ejRtGietB0p")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	//if resp.StatusCode != http.StatusOK {
	//	log.Fatal("Http status code:", resp.StatusCode)
	//}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
