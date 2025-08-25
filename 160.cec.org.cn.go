package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var CecCookie = "Hm_lvt_49d543ff6a4932299291456dd99019b0=1754902653,1755956209; HMACCOUNT=1CCD0111717619C6; Hm_lpvt_49d543ff6a4932299291456dd99019b0=1756086545"

// CecSpider 获取中国电力企业联合会文档
// @Title 获取中国电力企业联合会文档
// @Description https://cec.org.cn/，将中国电力企业联合会文档入库
func main() {
	var startId = 177841
	var endId = 177999
	for id := startId; id <= endId; id++ {
		detailUrl := fmt.Sprintf("https://cec.org.cn/ms-mcms/mcms/content/detail?id=%d", id)
		fmt.Println(detailUrl)
		queryCecDetailResponseData, err := QueryCecDetail(detailUrl, fmt.Sprintf("https://cec.org.cn/detail/index.html?3-%d", id))
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 文档名称
		title := strings.TrimSpace(queryCecDetailResponseData.BasicTitle)
		title = strings.ReplaceAll(title, "/", "-")
		title = strings.ReplaceAll(title, "／", "-")
		title = strings.ReplaceAll(title, "—", "-")
		title = strings.ReplaceAll(title, " ", "-")
		title = strings.ReplaceAll(title, "　", "-")
		title = strings.ReplaceAll(title, "/", "-")
		title = strings.ReplaceAll(title, "：", "-")
		title = strings.ReplaceAll(title, "--", "-")
		fmt.Println(title)

		filePath := "../cec.org.cn/" + title + ".pdf"
		if _, err := os.Stat(filePath); err != nil {
			//fmt.Println(queryCecDetailResponseData.ArticleContent)
			fmt.Println("=======开始下载========")
			reg := regexp.MustCompile("href=\"(.*?).pdf\">")
			path2 := reg.Find([]byte(queryCecDetailResponseData.ArticleContent))
			path2Str := string(path2)
			path2StrHandle := strings.ReplaceAll(path2Str, "href=\"", "")
			path2StrHandle = strings.ReplaceAll(path2StrHandle, "\">", "")

			downloadUrl := "https://cec.org.cn" + path2StrHandle
			fmt.Println(downloadUrl)

			err = downCecPdf(downloadUrl, filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//复制文件
			tempFilePath := strings.ReplaceAll(filePath, "../cec.org.cn", "../upload.doc88.com/hbba.sacinfo.org.cn")
			err = copyCecFile(filePath, tempFilePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("=======完成下载========")

			// 设置倒计时
			DownLoadCecTimeSleep := 10
			for i := 1; i <= DownLoadCecTimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("id="+strconv.Itoa(id)+"===========操作完成，", "暂停", DownLoadCecTimeSleep, "秒，倒计时", i, "秒===========")
			}
		}
	}
}

type QueryCecDetailResponse struct {
	Data    QueryCecDetailResponseData `json:"data"`
	Msg     string                     `json:"msg"`
	Status  int                        `json:"status"`
	Success bool                       `json:"success"`
}

type QueryCecDetailResponseData struct {
	BasicTitle     string `json:"basicTitle"`
	ArticleContent string `json:"articleContent"`
}

func QueryCecDetail(requestUrl string, referer string) (queryCecDetailResponseData QueryCecDetailResponseData, err error) {
	// 初始化客户端
	client := &http.Client{}                            //初始化客户端
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return queryCecDetailResponseData, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", CecCookie)
	req.Header.Set("Host", "cec.org.cn")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		fmt.Println(err)
		return queryCecDetailResponseData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return queryCecDetailResponseData, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return queryCecDetailResponseData, err
	}
	queryCecDetailResponse := &QueryCecDetailResponse{}
	err = json.Unmarshal(respBytes, queryCecDetailResponse)
	if err != nil {
		return queryCecDetailResponseData, err
	}
	queryCecDetailResponseData = queryCecDetailResponse.Data
	return queryCecDetailResponseData, nil
}

func downCecPdf(pdfUrl string, filePath string) error {
	client := &http.Client{}                        //初始化客户端
	req, err := http.NewRequest("GET", pdfUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "__51cke__=; _gid=GA1.3.657608059.1670484462; ASP.NET_SessionId=3prbrx4xve3rhlmvwbexp3v5; __tins__18926186=%7B%22sid%22%3A%201670551578342%2C%20%22vd%22%3A%202%2C%20%22expires%22%3A%201670553390816%7D; __51laig__=73; _ga_34B604LFFQ=GS1.1.1670556735.6.1.1670558647.53.0.0; _ga=GA1.1.711340106.1670484462")
	req.Header.Set("Host", "cec.org.cn")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
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

func copyCecFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer func(in *os.File) {
		err := in.Close()
		if err != nil {
			return
		}
	}(in)

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)

	_, err = io.Copy(out, in)
	return
}
