package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
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

const (
	CaAcEnableHttpProxy = false
	CaAcHttpProxyUrl    = "111.225.152.186:8089"
)

func CaAcSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(CaAcHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

// 获取中国民航航空局标准
// @Title 获取中国民航航空局标准
// @Description https://www.caac.gov.cn/ 获取中国民航航空局标准
func main() {
	maxPage := 108
	page := 1
	isPageListGo := true
	for isPageListGo {
		requestUrl := fmt.Sprintf("https://www.caac.gov.cn/was5/web/search?page=%d&channelid=211383&orderby=-fabuDate&was_custom_expr=+PARENTID%3D%2715%27+or+CLASSINFOID%3D%2715%27+&perpage=10&outlinepage=7&orderby=-fabuDate&selST=All&fl=15", page)
		fmt.Println(requestUrl)

		pageDoc, err := htmlquery.LoadURL(requestUrl)
		if err != nil {
			fmt.Println(err)
		}
		// /html/body/div/div[2]/div/table/tbody/tr/td/table/tbody/tr[1]
		trNodes := htmlquery.Find(pageDoc, `//html/body/div/div[2]/div/table/tbody/tr/td/table/tbody/tr`)
		if len(trNodes) <= 0 {
			isPageListGo = false
			break
		}

		for _, trNode := range trNodes {

			aNode := htmlquery.FindOne(trNode, `./td[@class="t_l tdMC"]/a`)
			if aNode == nil {
				fmt.Println("未找到连接节点，跳过")
				continue
			}

			hrefNode := htmlquery.FindOne(aNode, `./@href`)
			detailUrl := htmlquery.InnerText(hrefNode)
			fmt.Println(detailUrl)

			title := strings.TrimSpace(htmlquery.InnerText(aNode))
			title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, "　", "-")
			title = strings.ReplaceAll(title, "/", "-")
			title = strings.ReplaceAll(title, "--", "-")
			fmt.Println(title)

			codeNode := htmlquery.FindOne(trNode, `./td[@class="strFL"]`)
			if codeNode == nil {
				fmt.Println("未找到标准号节点，跳过")
				continue
			}
			code := strings.TrimSpace(htmlquery.InnerText(codeNode))
			code = strings.ReplaceAll(code, "/", "-")
			code = strings.ReplaceAll(code, "—", "-")

			filePath := "../www.caac.gov.cn/" + title + "(" + code + ")" + ".pdf"
			if len(code) <= 0 {
				filePath = "../www.caac.gov.cn/" + title + ".pdf"
			}
			fmt.Println(filePath)
			if _, err := os.Stat(filePath); err != nil {
				detailDoc, err := htmlquery.LoadURL(detailUrl)
				if err != nil {
					fmt.Println(err)
					continue
				}

				downloadNode := htmlquery.FindOne(detailDoc, `//div[@class="wrap"]/div[@class="wrap_w p_t30"]/div[@class="clearfix"]/div[@class="a_left"]/div[@class="content"]/div[@id="id_tblAppendix"]/p/a/@href`)
				if downloadNode == nil {
					fmt.Println("未找到下载文件节点，跳过")
					continue
				}

				detailUrlArray := strings.Split(detailUrl, "/")
				downloadUrlArray := detailUrlArray[:len(detailUrlArray)-1]
				downloadUrl := strings.Join(downloadUrlArray, "/") + strings.ReplaceAll(htmlquery.InnerText(downloadNode), "./", "/")
				fmt.Println(downloadUrl)
				fmt.Println("=======开始下载" + title + "========")
				err = downloadCaAc(downloadUrl, detailUrl, filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "../www.caac.gov.cn", "../upload.doc88.com/www.caac.gov.cn")
				err = copyCaAcFile(filePath, tempFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("=======下载完成========")
				//DownLoadCaAcTimeSleep := 10
				DownLoadCaAcTimeSleep := rand.Intn(5)
				for i := 1; i <= DownLoadCaAcTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("page="+strconv.Itoa(page)+",filePath="+filePath+"===========下载成功 暂停", DownLoadCaAcTimeSleep, "秒 倒计时", i, "秒===========")
				}

			}
		}
		DownLoadCaAcPageTimeSleep := 10
		// DownLoadCaAcPageTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadCaAcPageTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("page="+strconv.Itoa(page)+"=========== 暂停", DownLoadCaAcPageTimeSleep, "秒 倒计时", i, "秒===========")
		}
		page++
		if page > maxPage {
			isPageListGo = false
			break
		}
	}
}

func downloadCaAc(attachmentUrl string, referer string, filePath string) error {
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
	if CaAcEnableHttpProxy {
		client = CaAcSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.caac.gov.cn")
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

func copyCaAcFile(src, dst string) (err error) {
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
