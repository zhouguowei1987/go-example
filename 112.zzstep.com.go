package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var ZZStepEnableHttpProxy = true
var ZZStepHttpProxyUrl = ""
var ZZStepHttpProxyUrlArr = make([]string, 0)

func ZZStepHttpProxy() error {
	pageMax := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	//pageMax := []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for _, page := range pageMax {
		freeProxyUrl := "https://www.beesproxy.com/free"
		if page > 1 {
			freeProxyUrl = fmt.Sprintf("https://www.beesproxy.com/free/page/%d", page)
		}
		beesProxyDoc, err := htmlquery.LoadURL(freeProxyUrl)
		if err != nil {
			return err
		}
		trNodes := htmlquery.Find(beesProxyDoc, `//figure[@class="wp-block-table"]/table[@class="table table-bordered bg--secondary"]/tbody/tr`)
		if len(trNodes) > 0 {
			for _, trNode := range trNodes {
				ipNode := htmlquery.FindOne(trNode, "./td[1]")
				if ipNode == nil {
					continue
				}
				ip := htmlquery.InnerText(ipNode)

				portNode := htmlquery.FindOne(trNode, "./td[2]")
				if portNode == nil {
					continue
				}
				port := htmlquery.InnerText(portNode)

				protocolNode := htmlquery.FindOne(trNode, "./td[5]")
				if protocolNode == nil {
					continue
				}
				protocol := htmlquery.InnerText(protocolNode)

				switch protocol {
				case "HTTP":
					ZZStepHttpProxyUrlArr = append(ZZStepHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					ZZStepHttpProxyUrlArr = append(ZZStepHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func ZZStepSetHttpProxy() (httpclient *http.Client) {
	if ZZStepHttpProxyUrl == "" {
		if len(ZZStepHttpProxyUrlArr) <= 0 {
			err := ZZStepHttpProxy()
			if err != nil {
				ZZStepSetHttpProxy()
			}
		}
		ZZStepHttpProxyUrl = ZZStepHttpProxyUrlArr[0]
		if len(ZZStepHttpProxyUrlArr) >= 2 {
			ZZStepHttpProxyUrlArr = ZZStepHttpProxyUrlArr[1:]
		} else {
			ZZStepHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(ZZStepHttpProxyUrl)
	ProxyURL, _ := url.Parse(ZZStepHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
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
	return httpclient
}

type ZZStepPaper struct {
	name string
	url  string
}

type ZZStepSubject struct {
	name   string
	id     int
	papers []ZZStepPaper
}

type ZZStepStudySectionSubjectsPapers struct {
	name     string
	id       int
	subjects []ZZStepSubject
}

var studySectionSubjectsPapers = []ZZStepStudySectionSubjectsPapers{
	//{
	//	name: "小学",
	//	subjects: []ZZStepSubject{
	//		{
	//			name: "语文",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=203&subject=29&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "数学",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=203&subject=30&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "英语",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=203&subject=31&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "道德与法治",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=203&subject=34&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "音乐",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=203&subject=41&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "美术",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=203&subject=42&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "信息技术",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=203&subject=43&page=1",
	//				},
	//			},
	//		},
	//	},
	//},

	//{
	//	name: "初中",
	//	subjects: []ZZStepSubject{
	//		{
	//			name: "语文",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=204&subject=29&page=1",
	//				},
	//				{
	//					name: "中考",
	//					url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=204&subject=29&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "数学",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=204&subject=30&page=1",
	//				},
	//				{
	//					name: "中考",
	//					url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=204&subject=30&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "英语",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=204&subject=31&page=1",
	//				},
	//				{
	//					name: "中考",
	//					url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=204&subject=31&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "物理",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=204&subject=32&page=1",
	//				},
	//				{
	//					name: "中考",
	//					url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=204&subject=32&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "化学",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=204&subject=33&page=1",
	//				},
	//				{
	//					name: "中考",
	//					url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=204&subject=33&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "生物",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=204&subject=37&page=1",
	//				},
	//				{
	//					name: "中考",
	//					url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=204&subject=37&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "道德与法治",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=204&subject=34&page=1",
	//				},
	//				{
	//					name: "中考",
	//					url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=204&subject=34&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "历史",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=204&subject=35&page=1",
	//				},
	//				{
	//					name: "中考",
	//					url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=204&subject=35&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "地理",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=204&subject=36&page=1",
	//				},
	//				{
	//					name: "中考",
	//					url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=204&subject=36&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "音乐",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=204&subject=41&page=1",
	//				},
	//				{
	//					name: "中考",
	//					url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=204&subject=41&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "美术",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=204&subject=42&page=1",
	//				},
	//				{
	//					name: "中考",
	//					url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=204&subject=42&page=1",
	//				},
	//			},
	//		},
	//		{
	//			name: "信息技术",
	//			papers: []ZZStepPaper{
	//				{
	//					name: "试卷",
	//					url:  "http://www2.zzstep.com/front/paper/index.html?studysection=204&subject=43&page=1",
	//				},
	//				{
	//					name: "中考",
	//					url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=204&subject=43&page=1",
	//				},
	//			},
	//		},
	//	},
	//},

	{
		name: "高中",
		subjects: []ZZStepSubject{
			{
				name: "语文",
				papers: []ZZStepPaper{
					{
						name: "试卷",
						url:  "http://www2.zzstep.com/front/paper/index.html?studysection=205&subject=29&page=1",
					},
					{
						name: "高考",
						url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=205&subject=29&page=1",
					},
				},
			},
			//{
			//	name: "数学",
			//	papers: []ZZStepPaper{
			//		{
			//			name: "试卷",
			//			url:  "http://www2.zzstep.com/front/paper/index.html?studysection=205&subject=30&page=1",
			//		},
			//		{
			//			name: "高考",
			//			url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=205&subject=30&page=1",
			//		},
			//	},
			//},
			//{
			//	name: "英语",
			//	papers: []ZZStepPaper{
			//		{
			//			name: "试卷",
			//			url:  "http://www2.zzstep.com/front/paper/index.html?studysection=205&subject=31&page=1",
			//		},
			//		{
			//			name: "高考",
			//			url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=205&subject=31&page=1",
			//		},
			//	},
			//},
			//{
			//	name: "物理",
			//	papers: []ZZStepPaper{
			//		{
			//			name: "试卷",
			//			url:  "http://www2.zzstep.com/front/paper/index.html?studysection=205&subject=32&page=1",
			//		},
			//		{
			//			name: "高考",
			//			url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=205&subject=32&page=1",
			//		},
			//	},
			//},
			//{
			//	name: "化学",
			//	papers: []ZZStepPaper{
			//		{
			//			name: "试卷",
			//			url:  "http://www2.zzstep.com/front/paper/index.html?studysection=205&subject=33&page=1",
			//		},
			//		{
			//			name: "高考",
			//			url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=205&subject=33&page=1",
			//		},
			//	},
			//},
			//{
			//	name: "生物",
			//	papers: []ZZStepPaper{
			//		{
			//			name: "试卷",
			//			url:  "http://www2.zzstep.com/front/paper/index.html?studysection=205&subject=37&page=1",
			//		},
			//		{
			//			name: "高考",
			//			url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=205&subject=37&page=1",
			//		},
			//	},
			//},
			//{
			//	name: "政治",
			//	papers: []ZZStepPaper{
			//		{
			//			name: "试卷",
			//			url:  "http://www2.zzstep.com/front/paper/index.html?studysection=205&subject=34&page=1",
			//		},
			//		{
			//			name: "高考",
			//			url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=205&subject=34&page=1",
			//		},
			//	},
			//},
			//{
			//	name: "历史",
			//	papers: []ZZStepPaper{
			//		{
			//			name: "试卷",
			//			url:  "http://www2.zzstep.com/front/paper/index.html?studysection=205&subject=35&page=1",
			//		},
			//		{
			//			name: "高考",
			//			url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=205&subject=35&page=1",
			//		},
			//	},
			//},
			//{
			//	name: "地理",
			//	papers: []ZZStepPaper{
			//		{
			//			name: "试卷",
			//			url:  "http://www2.zzstep.com/front/paper/index.html?studysection=205&subject=36&page=1",
			//		},
			//		{
			//			name: "高考",
			//			url:  "http://www2.zzstep.com/front/beikao/index.html?studysection=205&subject=36&page=1",
			//		},
			//	},
			//},
		},
	},
}

var NextDownloadSleep = 2

var randStringLength = 8

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 当前账号已下载文档数量
var eachUsernameDownloadCurrentCount = 0

// 每个账号最大下载数量
var eachUsernameDownloadMaxCount = 80
var password = "123456"
var refer = "http://www.zzstep.com/"
var ZZStepCookie = ""

// ychEduSpider 获取中国教育出版网文档
// @Title 获取中国教育出版网文档
// @Description http://www2.zzstep.com/，获取中国教育出版网文档
func main() {
	for _, studySection := range studySectionSubjectsPapers {
		for _, subject := range studySection.subjects {
			for _, paper := range subject.papers {
				current := 1
				isPageListGo := true
				for isPageListGo {
					subjectIndexUrl := paper.url
					subjectIndexUrl = strings.ReplaceAll(subjectIndexUrl, "page=1", "page="+strconv.Itoa(current))
					subjectIndexDoc, err := htmlquery.LoadURL(subjectIndexUrl)
					if err != nil {
						fmt.Println(err)
						current = 1
						isPageListGo = false
						continue
					}
					liNodes := htmlquery.Find(subjectIndexDoc, `//div[@class="zy-list fn-mt20"]/ul[@class="reslist"]/li[@class="fn-pt20 fn-pb20"]`)
					if len(liNodes) <= 0 {
						fmt.Println(err)
						current = 1
						isPageListGo = false
						continue
					}
					for _, liNode := range liNodes {
						fmt.Println("============================================================================")
						fmt.Println("主题：", studySection.name, subject.name, paper.name)
						fmt.Println("=======当前页URL", subjectIndexUrl, "========")

						// 所需智币
						pointsNode := htmlquery.FindOne(liNode, `./div[@class="btn-item fn-left"]/div[@class="money fn-pt10"]`)
						if pointsNode == nil {
							fmt.Println("没有智币div")
							continue
						}
						pointsText := htmlquery.InnerText(pointsNode)
						fmt.Println(pointsText)
						pointsText = strings.ReplaceAll(pointsText, "智币", "")

						points, err := strconv.Atoi(pointsText)
						if err != nil {
							fmt.Println(err)
							continue
						}
						if points > 0 {
							fmt.Println("需要智币下载", points)
							continue
						}

						// 当前文件类型
						fileExtTextNode := htmlquery.FindOne(liNode, `./div[@class="filetype fn-pl10 fn-left"]/img/@src`)
						if fileExtTextNode == nil {
							fmt.Println("没有文件类型div")
							continue
						}
						fileExtText := htmlquery.InnerText(fileExtTextNode)
						fileExtText = strings.ReplaceAll(fileExtText, "/public/front/images/", "")
						if !strings.Contains(fileExtText, "typeicon-word.png") {
							fmt.Println(fileExtText, "不在下载后缀列表")
							continue
						}

						fileName := htmlquery.InnerText(htmlquery.FindOne(liNode, `./div[@class="zy-box fn-left"]/div[@class="subject-t"]/a`))
						fileName = strings.TrimSpace(fileName)
						fileName = strings.ReplaceAll(fileName, "/", "-")
						fileName = strings.ReplaceAll(fileName, ":", "-")
						fileName = strings.ReplaceAll(fileName, "：", "-")
						fileName = strings.ReplaceAll(fileName, "（", "(")
						fileName = strings.ReplaceAll(fileName, "）", ")")
						fmt.Println(fileName)
						if strings.Contains(fileName, "图片") {
							fmt.Println("图片版，跳过")
							continue
						}

						if strings.Contains(fileName, "扫描") {
							fmt.Println("扫描版，跳过")
							continue
						}

						filePath := "../www2.zzstep.com/www2.zzstep.com/" + studySection.name + "/" + subject.name + "/" + fileName
						_, errDoc := os.Stat(filePath + ".doc")
						_, errDocx := os.Stat(filePath + ".docx")
						if errDoc != nil && errDocx != nil {
							viewUrl := "http://www2.zzstep.com" + htmlquery.InnerText(htmlquery.FindOne(liNode, `./div[@class="zy-box fn-left"]/div[@class="subject-t"]/a/@href`))
							fmt.Println(viewUrl)

							downLoadUrl := strings.ReplaceAll(viewUrl, "index", "download")
							fmt.Println(downLoadUrl)

							fmt.Println("=======开始下载" + strconv.Itoa(current) + "========")
							err = downloadZZStep(downLoadUrl, viewUrl, filePath)
							if err != nil {
								fmt.Println(err)
								//注册登陆新账号
								err = ZZStepProxyRegisterLoginUsername()
								if err != nil {
									// 将代理清空，重新获取
									ZZStepHttpProxyUrl = ""
									fmt.Println(err)
								}
								eachUsernameDownloadCurrentCount = 0
								continue
							}
							fmt.Println("=======下载完成========")
							for i := 1; i <= NextDownloadSleep; i++ {
								time.Sleep(time.Second)
								fmt.Println("===========操作结束，暂停", NextDownloadSleep, "秒，倒计时", i, "秒===========")
							}
							if eachUsernameDownloadCurrentCount++; eachUsernameDownloadCurrentCount >= eachUsernameDownloadMaxCount {
								//注册登陆新账号
								err = ZZStepProxyRegisterLoginUsername()
								if err != nil {
									// 将代理清空，重新获取
									ZZStepHttpProxyUrl = ""
									fmt.Println(err)
								}
								eachUsernameDownloadCurrentCount = 0
								continue
							}
						}
					}
					current++
					isPageListGo = true
				}
			}
		}
	}
}

func ZZStepProxyRegisterLoginUsername() error {
	// 注册新账号
	rand.Seed(time.Now().UnixNano()) // 设置随机种子
	// 生成长度为randStringLength的随机字符串
	username := randStringBytes(randStringLength)
	fmt.Println(username)
	err := ZZStepRegisterRandUsername(username, password, password, refer)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 登陆
	err = ZZStepLoginUsername(username, password, refer)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

type ZZStepRegisterRandUsernameResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func ZZStepRegisterRandUsername(username string, password string, password2 string, refer string) error {
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
	if ZZStepEnableHttpProxy {
		client = ZZStepSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("username", username)
	postData.Add("password", password)
	postData.Add("password2", password2)
	postData.Add("refer", refer)
	requestUrl := "http://www2.zzstep.com/front/regist/registbyusername.html"
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Host", "www2.zzstep.com")
	req.Header.Set("Origin", "http://www2.zzstep.com")
	req.Header.Set("Referer", "http://www2.zzstep.com/front/regist/index.html?refer=http://www.zzstep.com/")
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

	respBytes, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(respBytes))
	if err != nil {
		return err
	}

	zZStepRegisterRandAccountResp := ZZStepRegisterRandUsernameResp{}
	err = json.Unmarshal(respBytes, &zZStepRegisterRandAccountResp)
	if err != nil {
		return err
	}

	if zZStepRegisterRandAccountResp.Code != 1 {
		return errors.New(zZStepRegisterRandAccountResp.Msg)
	}
	return nil
}

type ZZStepLoginUsernameResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func ZZStepLoginUsername(username string, passwordu string, refer string) error {
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
	if ZZStepEnableHttpProxy {
		client = ZZStepSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("username", username)
	postData.Add("passwordu", passwordu)
	postData.Add("type", "username")
	postData.Add("refer", refer)
	requestUrl := "http://www2.zzstep.com/front/login/dologin.html"
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Host", "www2.zzstep.com")
	req.Header.Set("Origin", "http://www2.zzstep.com")
	req.Header.Set("Referer", "http://www2.zzstep.com/front/login/index.html?refer=http://www2.zzstep.com/")
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

	respBytes, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBytes))
	if err != nil {
		return err
	}

	zZStepLoginUsernameResp := ZZStepLoginUsernameResp{}
	err = json.Unmarshal(respBytes, &zZStepLoginUsernameResp)
	if err != nil {
		return err
	}

	if zZStepLoginUsernameResp.Code != 1 {
		return errors.New(zZStepLoginUsernameResp.Msg)
	}

	// 重新设置cookie
	ZZStepCookie = resp.Header.Get("Set-Cookie")
	return nil
}

func downloadZZStep(attachmentUrl string, referer string, filePath string) error {
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
	if ZZStepEnableHttpProxy {
		client = ZZStepSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", ZZStepCookie)
	req.Header.Set("Host", "www2.zzstep.com")
	req.Header.Set("Referer", referer)
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
		return err
	}
	// 检查HTTP响应头中的Content-Disposition字段获取文件名和后缀
	fileName := getZZStepFileNameFromHeader(resp)
	fileExtension := filepath.Ext(fileName) // 获取文件后缀
	fileExtArr := []string{".doc", ".docx"}
	fmt.Println("文件后缀:", fileExtension)
	if !StrInArrayZZStep(fileExtension, fileExtArr) {
		return errors.New("文件后缀：" + fileExtension + "不在下载后缀列表")
	}
	filePath += fileExtension
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

// 生成随机字符串
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// StrInArrayZZStep str in string list
func StrInArrayZZStep(str string, data []string) bool {
	if len(data) > 0 {
		for _, row := range data {
			if str == row {
				return true
			}
		}
	}
	return false
}

// 从HTTP响应头中获取文件名
func getZZStepFileNameFromHeader(resp *http.Response) string {
	contentDisposition := resp.Header.Get("Content-Disposition")
	fileName := ""
	if contentDisposition != "" {
		fileName = parseZZStepFileNameFromContentDisposition(contentDisposition)
	} else {
		fileName = filepath.Base(resp.Request.URL.Path) // 默认使用URL中的文件名作为本地文件名
	}
	return fileName
}

// 从Content-Disposition字段中解析文件名
func parseZZStepFileNameFromContentDisposition(contentDisposition string) string {
	// 参考：https://tools.ietf.org/html/rfc6266#section-4.3
	// 示例：attachment; filename="example.txt" -> example.txt
	fileNameStart := len("attachment; ") + len("filename=") + 1
	fileNameEnd := len(contentDisposition) - 1
	fileName := ""
	if fileNameStart <= fileNameEnd {
		fileName = contentDisposition[fileNameStart:fileNameEnd] // 提取文件名字符串
	}
	return fileName[:] // 去掉字符串开头的引号（如果存在）并返回结果
}
