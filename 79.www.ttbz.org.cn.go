package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var TbzCookie = "safeline_bot_token=AG8+lRQAAAAAAAAAAAAAAAAa7q02mAEAAG4SzO1xkkO5pPxKx0xR2XpjfI5g; __51cke__=; ASP.NET_SessionId=kztywruoejdlibymyizdhk0z; __jsluid_s=4e74def0a093311d3c09d1caf84e0fb9; __jsluid_h=746830dcf9977a07d36ede64b0d40169; __51vcke__undefined=a3421553-a204-5f6a-8013-82062b198471; __51vuft__undefined=1745827959239; CCKF_visitor_id_830188=1404949643; _d_id=ce2e43c6bc70d5759609bc8d782ec0; __51uvsct__undefined=309; __vtins__undefined=%7B%22sid%22%3A%20%22cb421605-3c6f-5575-b7b6-0cf53a807703%22%2C%20%22vd%22%3A%205%2C%20%22stt%22%3A%202855225%2C%20%22dr%22%3A%201186%2C%20%22expires%22%3A%201751441809494%2C%20%22ct%22%3A%201751440009494%7D; __tins__18926186=%7B%22sid%22%3A%201751602316401%2C%20%22vd%22%3A%202%2C%20%22expires%22%3A%201751605017483%7D; __51laig__=5193; Hm_lvt_8c446e9fafe752e4975210bc30d7ab9d=1751681868; HMACCOUNT=00EDEFEA78E0441D; safeline_bot_challenge_ans=BAAAAABvPpUUAAAAAAAAAAAAAAAAOHGtNpgBAACIHEVDwJH3VnPrt5MgZqKj6byGo+8p13wdb1kZpxwOpkkBHZKsRgfjRP7tad0Jc/j45l7I169; Hm_lpvt_8c446e9fafe752e4975210bc30d7ab9d=1753264054"

// TbzSpider 获取全国团体标准信息平台Pdf文档
// @Title 获取全国团体标准信息平台Pdf文档
// @Description https://www.ttbz.org.cn/，将全国团体标准信息平台Pdf文档入库
func main() {
	var startId = 141900
	var endId = 141990
	goCh := make(chan int, endId-startId)
	for id := startId; id <= endId; id++ {
		go func(id int) {
			err := tbzSpider(id)
			if err != nil {
				fmt.Println(err)
			}
			goCh <- id
		}(id)
		fmt.Println(<-goCh)
	}
	//tbzSpider(132099)
}

func getTbz(url string) (doc *html.Node, err error) {
	// 创建一个自定义的http.Transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 忽略证书验证
		},
	}
	client := &http.Client{Transport: tr}        //初始化客户端
	req, err := http.NewRequest("GET", url, nil) //建立连接
	if err != nil {
		return doc, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", TbzCookie)
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	req.Header.Set("Host", "www.ttbz.org.cn")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return doc, err
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return doc, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	doc, err = htmlquery.Parse(resp.Body)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func downloadPdf(pdfUrl string, filePath string) error {
	// 创建一个自定义的http.Transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 忽略证书验证
		},
	}
	client := &http.Client{Transport: tr}           //初始化客户端                     //初始化客户端
	req, err := http.NewRequest("GET", pdfUrl, nil) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", TbzCookie)
	req.Header.Set("Host", "www.ttbz.org.cn")
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
		if os.MkdirAll(fileDiv, 0644) != nil {
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

func copyTbzFile(src, dst string) (err error) {
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

func tbzSpider(id int) error {
	detailUrl := fmt.Sprintf("https://www.ttbz.org.cn/StandardManage/Detail/%d", id)
	detailDoc, err := getTbz(detailUrl)
	if err != nil {
		return err
	}
	detailTableNodes := htmlquery.Find(detailDoc, `//table[@class="tctable"]`)
	if len(detailTableNodes) == 3 {
		// 标准详细信息table
		detailTableNode := detailTableNodes[1]
		tbodyNode := htmlquery.FindOne(detailTableNode, `./tbody`)
		trNodes := htmlquery.Find(tbodyNode, `./tr`)

		// 判断是否有废止日期
		abolitionDateTrNode := trNodes[10]
		abolitionDateTdNodes := htmlquery.Find(abolitionDateTrNode, `./td`)
		abolitionDateText := htmlquery.InnerText(abolitionDateTdNodes[0])
		// 标准文本URL tr index
		aHrefTrNodeIndex := 15
		if abolitionDateText == "废止日期  " {
			// 都要加1,因为在废止日期后面
			aHrefTrNodeIndex = 16
		}
		// 标准文本URL
		aHrefTrNode := trNodes[aHrefTrNodeIndex]
		aHrefTdNode := htmlquery.Find(aHrefTrNode, `./td`)
		aHrefText := strings.TrimSpace(htmlquery.InnerText(aHrefTdNode[1]))
		if aHrefText != "" {
			span2Text := htmlquery.InnerText(htmlquery.FindOne(aHrefTdNode[1], `./span[@id="Span2"]`))
			if span2Text == "不公开" {
				return errors.New("不公开文档，跳过")
			}
			aHref := htmlquery.InnerText(htmlquery.FindOne(aHrefTdNode[1], `./span[@id="Span2"]/a`))
			if aHref != "查看" {
				return errors.New("没有下载链接")
			}
			// 标准编号
			standardNoTrNode := trNodes[2]
			standardNoTdNodes := htmlquery.Find(standardNoTrNode, `./td`)
			standardNo := htmlquery.InnerText(htmlquery.FindOne(standardNoTdNodes[1], `./span[@id="r1_c5"]`))
			standardNo = strings.ReplaceAll(standardNo, "/", "-")
			fmt.Println(standardNo)

			// 中文标题
			chineseTitleTrNode := trNodes[3]
			chineseTitleTdNodes := htmlquery.Find(chineseTitleTrNode, `./td`)
			chineseTitle := htmlquery.InnerText(htmlquery.FindOne(chineseTitleTdNodes[1], `./span[@id="r1_c5"]`))
			chineseTitle = strings.ReplaceAll(chineseTitle, "/", "-")
			chineseTitle = strings.ReplaceAll(chineseTitle, " ", "")
			fmt.Println(chineseTitle)

			pdfsUrl := fmt.Sprintf("https://www.ttbz.org.cn/Pdfs/Index/?ftype=st&pms=%d", id)
			pdfsDoc, err := getTbz(pdfsUrl)
			if err != nil {
				return err
			}
			iframeSrcNode := htmlquery.FindOne(pdfsDoc, `//iframe[@id="myiframe"]/@src`)
			iframeSrc := htmlquery.InnerText(iframeSrcNode)
			fmt.Println(iframeSrc)

			// 下载pdf文件
			pdfUrl := strings.ReplaceAll(iframeSrc, "/Home/PdfView?file=", "https://www.ttbz.org.cn")
			// 移除rnd参数
			pdfUrl = strings.Split(pdfUrl, "&")[0]
			fmt.Println(pdfUrl)

			filePath := "../www.ttbz.org.cn/" + strconv.Itoa(id) + "-" + chineseTitle + "(" + strings.ReplaceAll(standardNo, "T/", "T") + ")" + ".pdf"
			if _, err := os.Stat(filePath); err != nil {
				fmt.Println("=======开始下载========")
				err = downloadPdf(pdfUrl, filePath)
				if err != nil {
					return err
				}
				//复制文件
				tempFilePath := strings.ReplaceAll(filePath, "www.ttbz.org.cn", "temp-www.ttbz.org.cn")
				err = copyTbzFile(filePath, tempFilePath)
				if err != nil {
					return err
				}
				fmt.Println("=======完成下载========")

				// 设置下载倒计时
				DownLoadTTbzTimeSleep := 10
				for i := 1; i <= DownLoadTTbzTimeSleep; i++ {
					time.Sleep(time.Second)
					fmt.Println("id="+strconv.Itoa(id)+"===========操作完成，", "暂停", DownLoadTTbzTimeSleep, "秒，倒计时", i, "秒===========")
				}
			}
		}
	}
	return nil
}
