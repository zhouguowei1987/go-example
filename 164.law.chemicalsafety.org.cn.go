package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	_ "os"
	"path/filepath"
	_ "path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	_ "golang.org/x/net/html"
)

var LawEnableHttpProxy = false
var LawHttpProxyUrl = "111.225.152.186:8089"
var LawHttpProxyUrlArr = make([]string, 0)

func LawHttpProxy() error {
	pageMax := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
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
					LawHttpProxyUrlArr = append(LawHttpProxyUrlArr, "http://"+ip+":"+port)
				case "HTTPS":
					LawHttpProxyUrlArr = append(LawHttpProxyUrlArr, "https://"+ip+":"+port)
				}
			}
		}
	}
	return nil
}

func LawSetHttpProxy() (httpclient *http.Client) {
	if LawHttpProxyUrl == "" {
		if len(LawHttpProxyUrlArr) <= 0 {
			err := LawHttpProxy()
			if err != nil {
				LawSetHttpProxy()
			}
		}
		LawHttpProxyUrl = LawHttpProxyUrlArr[0]
		if len(LawHttpProxyUrlArr) >= 2 {
			LawHttpProxyUrlArr = LawHttpProxyUrlArr[1:]
		} else {
			LawHttpProxyUrlArr = make([]string, 0)
		}
	}

	fmt.Println(LawHttpProxyUrl)
	ProxyURL, _ := url.Parse(LawHttpProxyUrl)
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
			ResponseHeaderTimeout: time.Second * 30,
		},
	}
	return httpclient
}

type QueryLawListRequestFormData struct {
	Command     string                               `json:"_command"`
	ContentType string                               `json:"_contentType"`
	Args        []QueryLawListRequestFormDataArgsMap `json:"_args"`
}

type QueryLawListRequestFormDataArgsMap struct {
	Index           int      `json:"index"`
	Count           int      `json:"count"`
	SortName        string   `json:"sortName"`
	SortOrder       string   `json:"sortOrder"`
	ModuleId        int      `json:"moduleId"`
	NodeId          int      `json:"nodeId"`
	KeyContent      string   `json:"keyContent"`
	Type            string   `json:"type"`
	Status1         string   `json:"status1"`
	Status2         string   `json:"status2"`
	Status3         string   `json:"status3"`
	RegulationType1 string   `json:"regulationType1"`
	RegulationType2 string   `json:"regulationType2"`
	RegulationType3 string   `json:"regulationType3"`
	RegulationType4 string   `json:"regulationType4"`
	RegulationType5 string   `json:"regulationType5"`
	TechType1       string   `json:"techType1"`
	TechType2       string   `json:"techType2"`
	TechType3       string   `json:"techType3"`
	TechType4       string   `json:"techType4"`
	Group           []string `json:"group"`
	Industry        []string `json:"industry"`
	Local           string   `json:"local"`
	StartDate       string   `json:"startDate"`
	EndDate         string   `json:"endDate"`
	ScStartDate     string   `json:"scStartDate"`
	ScEndDate       string   `json:"scEndDate"`
	SearchType      string   `json:"searchType"`
}

var LawListCookie = "home_page2=https://law.chemicalsafety.org.cn:443/compliance/guild/main/Main.jsp; JSESSIONID=778FAFBB43A113496FA464B506C141A0; Hm_lvt_8f8c4c4dbc9ac98f44dcd28c17e60d71=1756479886"

// 下载化学品安全法规标准文档
// @Title 下载化学品安全法规标准文档
// @Description https://law.chemicalsafety.org.cn/，下载化学品安全法规标准文档
func main() {
	pageListUrl := "https://law.chemicalsafety.org.cn/compliance/global/callService.action"
	page := 0
	maxPage := 1440
	isPageListGo := true
	for isPageListGo {
		queryLawListRequestFormData := QueryLawListRequestFormData{
			Command:     "CMMF70EDB9F095022BA2F04F37BA700898DA00E2FD4BC353A49FC4211D223FF93C584B02D4B2C8572F6DAC5426E89D38FD0D09DA370F23625E2",
			ContentType: "json",
			Args: []QueryLawListRequestFormDataArgsMap{
				{
					Index:           page,
					Count:           10,
					SortName:        "",
					SortOrder:       "",
					ModuleId:        2,
					NodeId:          0,
					KeyContent:      "",
					Type:            "全部",
					Status1:         "",
					Status2:         "",
					Status3:         "",
					RegulationType1: "",
					RegulationType2: "",
					RegulationType3: "",
					RegulationType4: "",
					RegulationType5: "",
					TechType1:       "",
					TechType2:       "",
					TechType3:       "",
					TechType4:       "",
					Group:           []string{},
					Industry:        []string{},
					Local:           "",
					StartDate:       "",
					EndDate:         "",
					ScStartDate:     "",
					ScEndDate:       "",
					SearchType:      "精确",
				},
			},
		}
		queryLawListResponseRows, err := QueryLawList(pageListUrl, queryLawListRequestFormData)
		if err != nil {
			fmt.Println(err)
			break
		}
		for _, law := range queryLawListResponseRows {
			fmt.Println("=========开始处理数据==============")

			code := law.Number
			code = strings.ReplaceAll(code, "/", "-")
			code = strings.ReplaceAll(code, "—", "-")
			fmt.Println(code)

			title := law.Title
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "--", "-")
			fmt.Println(title)

			filePath := "../law.chemicalsafety.org.cn/" + title + "(" + code + ")" + ".pdf"
			filePath = strings.ReplaceAll(filePath, "()", "")
			fmt.Println(filePath)

			_, err = os.Stat(filePath)
			if err == nil {
				fmt.Println("文档已下载过，跳过")
				continue
			}

			fmt.Println("=======开始下载========")
			downloadRefererUrl := fmt.Sprintf("https://law.chemicalsafety.org.cn/compliance/guild/tech/TechBrowse.jsp?moduleId=2&recordId=%s&libraryId=%s&attachmentId=%s", law.RecordId, law.LibraryId, law.AttachmentId)
			//fmt.Println(downloadRefererUrl)
			downloadUrl := fmt.Sprintf("https://law.chemicalsafety.org.cn/compliance/rmm/function/record/downloadFile.action?attachmentId=%s", law.AttachmentId)
			fmt.Println(downloadUrl)

			fmt.Println("=======开始下载" + title + "========")

			err = downloadLaw(downloadUrl, downloadRefererUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "../law.chemicalsafety.org.cn", "../upload.doc88.com/hbba.sacinfo.org.cn")
			err = copyLawFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======下载完成========")
			//DownLoadLawTimeSleep := 10
			DownLoadLawTimeSleep := rand.Intn(5)
			for i := 1; i <= DownLoadLawTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("filePath="+filePath+"===========下载成功 暂停", DownLoadLawTimeSleep, "秒 倒计时", i, "秒===========")
			}
		}
		DownLoadLawCategoryTimeSleep := 10
		// DownLoadLawCategoryTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadLawCategoryTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("=========== 暂停", DownLoadLawCategoryTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

type QueryLawListResponse struct {
	KeysAll string                     `json:"keysAll"`
	Rows    []QueryLawListResponseRows `json:"rows"`
	Total   int                        `json:"total"`
}

type QueryLawListResponseRows struct {
	AttachmentId string `json:"attId"`
	RecordId     string `json:"id"`
	LibraryId    string `json:"libraryId"`
	Number       string `json:"number"`
	Title        string `json:"title"`
}

func QueryLawList(requestUrl string, queryLawListRequestFormData QueryLawListRequestFormData) (queryLawListResponseRows []QueryLawListResponseRows, err error) {
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
	if LawEnableHttpProxy {
		client = LawSetHttpProxy()
	}

	// 将数据编码为JSON格式
	queryLawListRequestFormDataArgsJson, err := json.Marshal(queryLawListRequestFormData.Args)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	postData := url.Values{}
	postData.Add("_command", queryLawListRequestFormData.Command)
	postData.Add("_contentType", queryLawListRequestFormData.ContentType)
	postData.Add("_args", string(queryLawListRequestFormDataArgsJson))
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	queryLawListResponse := QueryLawListResponse{}
	if err != nil {
		return queryLawListResponseRows, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", LawListCookie)
	req.Header.Set("Host", "law.chemicalsafety.org.cn")
	req.Header.Set("Origin", "https://law.chemicalsafety.org.cn")
	req.Header.Set("Referer", "https://law.chemicalsafety.org.cn/compliance/guild/customer/SearchAllCustomer.jsp?moduleId=2&type=all")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryLawListResponseRows, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryLawListResponseRows, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryLawListResponseRows, err
	}
	err = json.Unmarshal(respBytes, &queryLawListResponse)
	if err != nil {
		return queryLawListResponseRows, err
	}
	queryLawListResponseRows = queryLawListResponse.Rows
	return queryLawListResponseRows, nil
}

func downloadLaw(attachmentUrl string, referer string, filePath string) error {
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
	if LawEnableHttpProxy {
		client = LawSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", LawListCookie)
	req.Header.Set("Host", "law.chemicalsafety.org.cn")
	//req.Header.Set("Referer", referer)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 如果访问失败 就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}

	// 创建一个文件用于保存
	fileDiv := filepath.Dir(filePath)
	if _, err = os.Stat(fileDiv); err != nil {
		if os.MkdirAll(fileDiv, 0o777) != nil {
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

func copyLawFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(in *os.File) {
		err := in.Close()
		if err != nil {
			return
		}
	}(in)

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)

	_, err = io.Copy(out, in)
	return nil
}
