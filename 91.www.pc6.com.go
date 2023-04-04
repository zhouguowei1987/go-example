package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	Pc6EnableHttpProxy = false
	Pc6HttpProxyUrl    = "111.225.152.186:8089"
)

func Pc6SetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(Pc6HttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type Pc6Subject struct {
	name string
	url  string
}

var AllPc6Subject = []Pc6Subject{
	//{
	//	name: "公文通知",
	//	url:  "https://www.pc6.com/pc/gwtzfwmb/",
	//},
	//{
	//	name: "规章制度",
	//	url:  "https://www.pc6.com/pc/gzzdwordmb/",
	//},
	//{
	//	name: "营销策划",
	//	url:  "https://www.pc6.com/pc/yxchwordmb/",
	//},
	//{
	//	name: "总结报告",
	//	url:  "https://www.pc6.com/pc/gzzjwordmb/",
	//},
	//{
	//	name: "演讲稿范文",
	//	url:  "https://www.pc6.com/pc/yjgfwmb/",
	//},
	//{
	//	name: "行政管理",
	//	url:  "https://www.pc6.com/pc/xzglwordmb/",
	//},
	//{
	//	name: "合同范本",
	//	url:  "https://www.pc6.com/pc/htfbwordmb/",
	//},
	//{
	//	name: "简历模板",
	//	url:  "https://www.pc6.com/pc/grjlwordmb/",
	//},
	{
		name: "读书分享",
		url:  "https://www.pc6.com/pc/dsfxpptmb/",
	},
	{
		name: "主题班会",
		url:  "https://www.pc6.com/pc/ztbhpptmb/",
	},
	{
		name: "活动策划",
		url:  "https://www.pc6.com/pc/hdchpptmb/",
	},
	{
		name: "企业宣传",
		url:  "https://www.pc6.com/pc/qyxcpptmb/",
	},
	{
		name: "工作计划",
		url:  "https://www.pc6.com/pc/gzjhpptmb/",
	},
	{
		name: "毕业答辩",
		url:  "https://www.pc6.com/pc/ztbhpptmb/",
	},
	{
		name: "企业培训",
		url:  "https://www.pc6.com/pc/qypxpptmb/",
	},
	{
		name: "节日庆典",
		url:  "https://www.pc6.com/pc/jrqdpptmb/",
	},
	{
		name: "个人简历",
		url:  "https://www.pc6.com/pc/grjlpptmb/",
	},
	{
		name: "教育教学",
		url:  "https://www.pc6.com/pc/jyjxpptmb/",
	},
	{
		name: "竞聘述职",
		url:  "https://www.pc6.com/pc/jpszbgpptmb/",
	},
	{
		name: "商业计划书",
		url:  "https://www.pc6.com/pc/syjhspptmb/",
	},
	{
		name: "工作总结",
		url:  "https://www.pc6.com/pc/gzzjpptmb/",
	},
}

// ychEduSpider 获取pc6文档
// @Title 获取pc6文档
// @Description https://www.pc6.com/，获取pc6文档
func main() {
	for _, subject := range AllPc6Subject {
		page := 1
		isPageListGo := true
		for isPageListGo {
			pageListUrl := subject.url
			if page > 1 {
				pageListUrl = fmt.Sprintf(subject.url+"%d/", page)
			}
			fmt.Println(pageListUrl)

			pageListDoc, err := htmlquery.LoadURL(pageListUrl)
			if err != nil {
				fmt.Println(err)
				break
			}

			dlNodes := htmlquery.Find(pageListDoc, `//div[@class="model_list"]/ul[@id="pullUp"]/li]`)
			if len(dlNodes) >= 1 {
				for _, dlNode := range dlNodes {

					fmt.Println("=================================================================================")
					// 文档详情URL
					fileName := htmlquery.InnerText(htmlquery.FindOne(dlNode, `./a/div[@class="txtbox"]/span`))
					fmt.Println(fileName)

					detailUrl := htmlquery.InnerText(htmlquery.FindOne(dlNode, `./a/@href`))
					detailUrl = "https://www.pc6.com" + detailUrl
					fmt.Println(detailUrl)

					detailDoc, _ := htmlquery.LoadURL(detailUrl)

					reg := regexp.MustCompile(`_downInfo={Address:"(.*?)"`)
					//将所有null替换为空字符串
					detailDocText := htmlquery.InnerText(detailDoc)
					regFindStingMatch := reg.FindStringSubmatch(detailDocText)

					// 下载文档URL
					downLoadUrl := "https://110pansoft.0098118.com" + regFindStingMatch[1]
					fmt.Println(downLoadUrl)

					// 文件格式
					attachmentFormat := strings.Split(downLoadUrl, ".")

					filePath := "../www.pc6.com/" + subject.name + "/"
					err = downloadPc6(downLoadUrl, detailUrl, filePath, fileName+"."+attachmentFormat[len(attachmentFormat)-1])
					if err != nil {
						fmt.Println(err)
						continue
					}
					time.Sleep(time.Second * 1)
				}
				page++
			} else {
				isPageListGo = false
				page = 1
				break
			}
		}
	}
}
func downloadPc6(attachmentUrl string, referer string, filePath string, fileName string) error {
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
	if Pc6EnableHttpProxy {
		client = Pc6SetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "https://110pansoft.0098118.com")
	req.Header.Set("Referer", referer)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
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
	out, err := os.Create(filePath + fileName)
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
