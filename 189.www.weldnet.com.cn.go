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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var WeldnetCookie = "Hm_lvt_435e560f9f9b8ee21820850608de7456=1778557819; HMACCOUNT=9C0CD19686802BBF; Hm_lpvt_435e560f9f9b8ee21820850608de7456=1778558188"

// WeldnetSpider 获取中国焊接协会标准文档
// @Title 获取中国焊接协会标准文档
// @Description http://www.weldnet.com.cn/，将中国焊接协会标准文档入库
func main() {
	var startId = 1
	var endId = 294
	for id := startId; id <= endId; id++ {
		detailUrl := fmt.Sprintf("http://admin.china-weldnet.com/api/SpecialDocument.asmx/WeldingStandardDetail?id=%d", id)
		detailReferer := "http://www.weldnet.com.cn/specialist-standard"
		tweldnetDetailResponse, err := GetWeldnetDetail(detailUrl, detailReferer)
		if err != nil {
			fmt.Println(err)
			break
		}
		title := tweldnetDetailResponse.Data[0].ChineseTitle
		title = strings.TrimSpace(title)
		title = strings.ReplaceAll(title, "/", "-")
		title = strings.ReplaceAll(title, "／", "-")
		title = strings.ReplaceAll(title, "/", "-")
		title = strings.ReplaceAll(title, "　", "-")
		title = strings.ReplaceAll(title, " ", "-")
		title = strings.ReplaceAll(title, "：", ":")
		title = strings.ReplaceAll(title, "—", "-")
		title = strings.ReplaceAll(title, "－", "-")
		title = strings.ReplaceAll(title, "（", "(")
		title = strings.ReplaceAll(title, "）", ")")
		title = strings.ReplaceAll(title, "《", "")
		title = strings.ReplaceAll(title, "》", "")
		fmt.Println(title)
		if len(title) <= 0 {
			fmt.Println("文档标题不存在，跳过")
			continue
		}

		StandardFile := tweldnetDetailResponse.Data[0].StandardFile
		if len(StandardFile) <= 0 {
			fmt.Println("文档文件不存在，跳过")
			continue
		}
		// 查看文档后缀
		fileExt := filepath.Ext(StandardFile)
		fileExt = strings.ToLower(fileExt)
		if strings.Index(fileExt, "doc") == -1 && strings.Index(fileExt, "pdf") == -1 {
			fmt.Println("文档不是doc、pdf文档，跳过")
			continue
		}

		filePath := "../www.weldnet.com.cn/" + title + fileExt
		fmt.Println(filePath)
		_, err = os.Stat(filePath)
		if err == nil {
			fmt.Println("文档已下载过，跳过")
			continue
		}
		downloadUrl := fmt.Sprintf("http://admin.china-weldnet.com/%s", StandardFile)
		fmt.Println(downloadUrl)

		fmt.Println("=======开始下载" + title + "========")

		err = downloadWeldnet(downloadUrl, detailUrl, filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("=======下载完成========")
		// 查看文件大小，如果是空文件，则删除
		fileInfo, err := os.Stat(filePath)
		if err == nil && fileInfo.Size() == 0 {
			fmt.Println("空文件删除")
			err = os.Remove(filePath)
			continue
		}
		//复制文件
		tempFilePath := strings.ReplaceAll(filePath, "www.weldnet.com.cn", "temp-hbba.sacinfo.org.cn")
		err = copyWeldnetFile(filePath, tempFilePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//DownLoadWeldnetTimeSleep := 10
		DownLoadWeldnetTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadWeldnetTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("filePath="+filePath+"===========下载成功 暂停", DownLoadWeldnetTimeSleep, "秒 倒计时", i, "秒===========")
		}
	}
}

type WeldnetDetailResponse struct {
	Code  int                         `json:"code"`
	Count int                         `json:"count"`
	Data  []WeldnetDetailResponseData `json:"data"`
	Msg   string                      `json:"msg"`
}
type WeldnetDetailResponseData struct {
	ChineseTitle string `json:"ChineseTitle"`
	StandardFile string `json:"StandardFile"`
}

func GetWeldnetDetail(requestUrl string, referer string) (tweldnetDetailResponse WeldnetDetailResponse, err error) {
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
	tweldnetDetailResponse = WeldnetDetailResponse{}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return tweldnetDetailResponse, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", WeldnetCookie)
	req.Header.Set("Host", "admin.china-weldnet.com")
	req.Header.Set("Origin", "http://www.weldnet.com.cn")
	req.Header.Set("Referer", referer)
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
		return tweldnetDetailResponse, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tweldnetDetailResponse, err
	}
	err = json.Unmarshal(respBytes, &tweldnetDetailResponse)
	if err != nil {
		return tweldnetDetailResponse, err
	}
	return tweldnetDetailResponse, nil
}

func downloadWeldnet(requestUrl string, referer string, filePath string) error {
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
	} //初始化客户端
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", WeldnetCookie)
	req.Header.Set("Host", "www.weldnet.com.cn")
	req.Header.Set("Origin", "http://www.weldnet.com.cn")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
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

func copyWeldnetFile(src, dst string) (err error) {
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

	// 创建一个文件用于保存
	fileDiv := filepath.Dir(dst)
	if _, err = os.Stat(fileDiv); err != nil {
		if os.MkdirAll(fileDiv, 0o777) != nil {
			return err
		}
	}
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
