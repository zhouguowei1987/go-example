package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// TbzSpider 获取广西标准化协会Pdf文档
// @Title 获取广西标准化协会Pdf文档
// @Description http://www.guangxibiaoxie.com/，将广西标准化协会Pdf文档入库
func main() {
	var startId = 1
	var endId = 2580
	goCh := make(chan int, endId-startId)
	for id := startId; id <= endId; id++ {
		go func(id int) {
			err := guangXiBiaoXieSpider(id)
			if err != nil {
				fmt.Println(err)
			}
			goCh <- id
		}(id)
		fmt.Println(<-goCh)
	}
	//guangXiBiaoXieSpider(116)
}

func guangXiBiaoXieSpider(id int) error {
	detailUrl := fmt.Sprintf("http://www.guangxibiaoxie.com/a/%d.html", id)
	fmt.Println(detailUrl)
	detailDoc, err := htmlquery.LoadURL(detailUrl)
	if err != nil {
		return err
	}
	// 查看是否有下载链接
	detailDocText := htmlquery.OutputHTML(detailDoc, true)
	regFile := regexp.MustCompile(`<a href="http://guangxibiaoxie.com/uploads/(.*?)"`)
	regFindStingMatch := regFile.FindStringSubmatch(detailDocText)
	if len(regFindStingMatch) < 2 {
		return errors.New("没有文档下载链接")
	}

	// 下载文档URL
	downLoadUrl := "http://guangxibiaoxie.com/uploads/" + regFindStingMatch[1]
	fmt.Println(downLoadUrl)

	// 文档名称
	titleNode := htmlquery.FindOne(detailDoc, `//div[@class="panel-body"]/div[@class="article-metas"]/h1[@class="metas-title"]`)
	title := htmlquery.InnerText(titleNode)
	title = strings.ReplaceAll(title, "（发布稿）", "")
	title = strings.ReplaceAll(title, "TGXAS", "T-GXAS")
	flagIndex := strings.Index(title, "团体标准")
	if flagIndex == -1 {
		return errors.New("不含有团体标准字样，跳过")
	}
	titleSplit := strings.Split(title, "团体标准")
	title = strings.ReplaceAll(titleSplit[1], " ", "")
	title = strings.ReplaceAll(title, "《", "")
	title = strings.ReplaceAll(title, "》", "")
	standardNo := strings.TrimSpace(titleSplit[0])
	title = title + "(" + standardNo + ")"

	filePath := "../www.guangxibiaoxie.com/" + title + ".pdf"
	if _, err := os.Stat(filePath); err != nil {
		fmt.Println("=======开始下载========")
		err = downGuangXiBiaoXiePdf(downLoadUrl, filePath)
		if err != nil {
			return err
		}
		fmt.Println("=======完成下载========")
	}
	return nil
}

func downGuangXiBiaoXiePdf(pdfUrl string, filePath string) error {
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
	req.Header.Set("Host", "www.guangxibiaoxie.com")
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
