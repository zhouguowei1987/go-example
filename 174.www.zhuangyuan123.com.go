package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	// "math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	ZhuangYuan123EnableHttpProxy = false
	ZhuangYuan123HttpProxyUrl    = "111.225.152.186:8089"
)

func ZhuangYuan123SetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(ZhuangYuan123HttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type ZhuangYuan123ListPeriod struct {
	periodId   int
	periodName string
	subject    []ZhuangYuan123ListSubject
}

type ZhuangYuan123ListSubject struct {
	periodId    int
	periodName  string
	subjectId   int
	subjectName string
}

var zhuangYuan123ListPeriodSubject = []ZhuangYuan123ListPeriod{
	// {
	// 	periodId:   1,
	// 	periodName: "小学",
	// 	subject: []ZhuangYuan123ListSubject{
	// 		/*// {subjectId: 43, subjectName: "语文"},
	// 		// {subjectId: 44, subjectName: "数*/学"},
	// 		{subjectId: 45, subjectName: "英语"},
	// 		{subjectId: 46, subjectName: "科学"},
	// 		{subjectId: 47, subjectName: "道德与法治"},
	// 		{subjectId: 48, subjectName: "音乐"},
	// 		{subjectId: 49, subjectName: "体育"},
	// 		{subjectId: 50, subjectName: "美术"},
	// 		{subjectId: 51, subjectName: "信息技术"},
	// 		{subjectId: 52, subjectName: "心理健康"},
	// 		{subjectId: 53, subjectName: "班会"},
	// 		{subjectId: 54, subjectName: "综合实践"},
	// 		{subjectId: 56, subjectName: "书法"},
	// 		{subjectId: 60, subjectName: "劳动技术"},
	// 		{subjectId: 61, subjectName: "专题教育"},
	// 	},
	// },
	{
		periodId:   2,
		periodName: "初中",
		subject: []ZhuangYuan123ListSubject{
			// {subjectId: 1, subjectName: "语文"},
			// {subjectId: 2, subjectName: "数学"},
			// {subjectId: 3, subjectName: "英语"},
			// {subjectId: 4, subjectName: "道德与法治"},
			// {subjectId: 5, subjectName: "历史"},
			// {subjectId: 6, subjectName: "物理"},
			// {subjectId: 7, subjectName: "生物"},
			// {subjectId: 8, subjectName: "化学"},
			{subjectId: 9, subjectName: "地理"},
			{subjectId: 10, subjectName: "科学"},
			{subjectId: 37, subjectName: "信息技术"},
			{subjectId: 36, subjectName: "历史与社会"},
			{subjectId: 38, subjectName: "音乐"},
			{subjectId: 39, subjectName: "美术"},
			{subjectId: 40, subjectName: "体育与健康"},
			{subjectId: 41, subjectName: "劳动技术"},
			{subjectId: 58, subjectName: "心理健康"},
			{subjectId: 42, subjectName: "综合"},
		},
	},
	// {
	// 	periodId:   3,
	// 	periodName: "高中",
	// 	subject: []ZhuangYuan123ListSubject{
	// 		{subjectId: 19, subjectName: "语文"},
	// 		{subjectId: 20, subjectName: "数学"},
	// 		{subjectId: 21, subjectName: "英语"},
	// 		{subjectId: 22, subjectName: "物理"},
	// 		{subjectId: 23, subjectName: "化学"},
	// 		{subjectId: 24, subjectName: "生物"},
	// 		{subjectId: 25, subjectName: "政治"},
	// 		{subjectId: 26, subjectName: "历史"},
	// 		{subjectId: 27, subjectName: "地理"},
	// 		{subjectId: 28, subjectName: "信息技术"},
	// 		{subjectId: 33, subjectName: "通用技术"},
	// 		{subjectId: 30, subjectName: "音乐"},
	// 		{subjectId: 32, subjectName: "体育与健康"},
	// 		{subjectId: 31, subjectName: "美术"},
	// 		{subjectId: 34, subjectName: "劳动技术"},
	// 		{subjectId: 57, subjectName: "心理健康"},
	// 		{subjectId: 35, subjectName: "拓展"},
	// 		{subjectId: 55, subjectName: "综合"},
	// 	},
	// },
}

// ychEduSpider 获取状元网试题
// @Title 获取状元网试题
// @Description http://www.zhuangyuan123.com/，获取状元网试题
func main() {
	for _, period := range zhuangYuan123ListPeriodSubject {
		fmt.Println("=======periodId：" + strconv.Itoa(period.periodId) + " ===periodName：" + period.periodName + "========")
		for _, subject := range period.subject {
			fmt.Println("=======subjectId：" + strconv.Itoa(subject.subjectId) + " ===subjectName：" + subject.subjectName + "========")
			pageNum := 1
			pageSize := 1000
			isPageListGo := true
			for isPageListGo {
				listUrl := fmt.Sprintf("http://iweb.zhuangyuan123.com/web/resources/list?period=%d&subjectId=%d&categoryId=&pointsId=&type=4&rank=&district=&province=&year=&grade=&examType=&status=2&pageNum=%d&orderByColumn=update_time&isAsc=desc&pageSize=%d", period.periodId, subject.subjectId, pageNum, pageSize)
				fmt.Println(listUrl)
				zhuangYuan123ListResponse, err := GetZhuangYuan123List(listUrl)
				if err != nil {
					fmt.Println(err)
					break
				}
				if len(zhuangYuan123ListResponse.Rows) > 0 {
					for _, row := range zhuangYuan123ListResponse.Rows {
						fmt.Println("===periodName：" + period.periodName + "===subjectName：" + subject.subjectName + "=====pageNum：" + strconv.Itoa(pageNum) + "========")
						title := row.Title
						title = strings.TrimSpace(title)
						title = strings.ReplaceAll(title, "/", "-")
						title = strings.ReplaceAll(title, " ", "")
						title = strings.ReplaceAll(title, "（", "(")
						title = strings.ReplaceAll(title, "）", ")")
						title = strings.ReplaceAll(title, "《", "")
						title = strings.ReplaceAll(title, "》", "")
						fmt.Println(title)

						previewUrl := row.PreviewUrl
						fmt.Println(previewUrl)
						// 查看是否是doc文件
						if strings.Index(previewUrl, ".doc") == -1 {
							fmt.Println("不是doc文件，跳过")
							continue
						}
						filePath := "../www.zhuangyuan123.com/2026-01-04/www.zhuangyuan123.com/" + period.periodName + "/" + subject.subjectName + "/" + title + ".doc"
						_, err = os.Stat(filePath)
						if err == nil {
							fmt.Println("文档已下载过，跳过")
							continue
						}
						fmt.Println("=======开始下载========")
						previewUrlArray := strings.Split(previewUrl, "furl=")
						downLoadUrl := previewUrlArray[1]
						fmt.Println(downLoadUrl)
						err = downloadZhuangYuan123(downLoadUrl, filePath)
						if err != nil {
							fmt.Println(err)
							continue
						}
						fmt.Println("=======开始完成========")
						// 设置倒计时
						DownLoadZhuangYuan123TimeSleep := 5
						// DownLoadZhuangYuan123TimeSleep := rand.Intn(3)
						for i := 1; i <= DownLoadZhuangYuan123TimeSleep; i++ {
							time.Sleep(time.Second)
							fmt.Println("暂停，periodName："+period.periodName+"===subjectName："+subject.subjectName+"===pageNumMax = ", strconv.Itoa((zhuangYuan123ListResponse.Total/pageSize)+1)+"=====pageNum："+strconv.Itoa(pageNum)+"，倒计时", i, "秒===========")
						}
					}
				}
				DownLoadZhuangYuan123PageTimeSleep := 8
				// DownLoadZhuangYuan123PageTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadZhuangYuan123PageTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("暂停，periodName："+period.periodName+"===subjectName："+subject.subjectName+"===pageNumMax = ", strconv.Itoa((zhuangYuan123ListResponse.Total/pageSize)+1)+"=====pageNum："+strconv.Itoa(pageNum)+"，倒计时", i, "秒===========")
				}
				pageNum++
				if pageNum > (zhuangYuan123ListResponse.Total/pageSize)+1 {
					fmt.Println("没有更多分页了")
					isPageListGo = false
					pageNum = 1
					break
				}

			}
		}
	}
}

type ZhuangYuan123ListResponse struct {
	Code  int                             `json:"code"`
	Msg   string                          `json:"msg"`
	Total int                             `json:"total"`
	Rows  []ZhuangYuan123ListResponseRows `json:"rows"`
}
type ZhuangYuan123ListResponseRows struct {
	PreviewUrl string `json:"previewUrl"`
	Title      string `json:"title"`
}

func GetZhuangYuan123List(requestUrl string) (zhuangYuan123ListResponse ZhuangYuan123ListResponse, err error) {
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
			ResponseHeaderTimeout: time.Second * 30,
		},
	}
	if ZhuangYuan123EnableHttpProxy {
		client = ZhuangYuan123SetHttpProxy()
	}
	zhuangYuan123ListResponse = ZhuangYuan123ListResponse{}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return zhuangYuan123ListResponse, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.zhuangyuan123.com")
	req.Header.Set("Origin", "http://www.zhuangyuan123.com/")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return zhuangYuan123ListResponse, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return zhuangYuan123ListResponse, err
	}
	err = json.Unmarshal(respBytes, &zhuangYuan123ListResponse)
	if err != nil {
		return zhuangYuan123ListResponse, err
	}
	return zhuangYuan123ListResponse, nil
}

func downloadZhuangYuan123(attachmentUrl string, filePath string) error {
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
			ResponseHeaderTimeout: time.Second * 30,
		},
	}
	if ZhuangYuan123EnableHttpProxy {
		client = ZhuangYuan123SetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.zhuangyuan123.com")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}

	// 创建一个文件用于保存
	fileDiv := filepath.Dir(filePath)
	if _, err = os.Stat(fileDiv); err != nil {
		if os.MkdirAll(fileDiv, 0777) != nil {
			return err
		}
	}
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
