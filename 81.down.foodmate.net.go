package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type internalCategory struct {
	id   int
	name string
}

var FoodMateCookie = "Hm_lvt_2aeaa32e7cee3cfa6e2848083235da9f=1730957037; HMACCOUNT=2CEC63D57647BCA5; Hm_lvt_d4fdc0f0037bcbb9bf9894ffa5965f5e=1730957038; __51cke__=; u_rdown=1; __gads=ID=6aa7bb34097db0c7:T=1730957036:RT=1730959091:S=ALNI_MZhgP1VBSkrn0nQ9pKk7tqwWsBCYQ; __gpi=UID=00000f776c9d87c1:T=1730957036:RT=1730959091:S=ALNI_MY1D8W2UisKk_PnJOK-lAlbgDC9WA; __eoi=ID=893729db1a75d12c:T=1730957036:RT=1730959091:S=AA-AfjbZRPAyoYfpui6zt8_ejNPv; Hm_lpvt_d4fdc0f0037bcbb9bf9894ffa5965f5e=1730959786; Hm_lpvt_2aeaa32e7cee3cfa6e2848083235da9f=1730959786; __tins__1484185=%7B%22sid%22%3A%201730959092377%2C%20%22vd%22%3A%204%2C%20%22expires%22%3A%201730961585874%7D; __51laig__=37"

// foodMateSpider 获取标准库Pdf文档
// @Title 获取标准库Pdf文档
// @Description http://down.foodmate.net/，获取标准库Pdf文档
func main() {
	// 国内标准列表
	var allCategory = []internalCategory{
		{id: 1, name: "国内标准"},
		{id: 2, name: "国外标准"},
	}
	for _, category := range allCategory {
		i := 1
		isPageGo := true
		for isPageGo {
			listUrl := fmt.Sprintf("http://down.foodmate.net/standard/sort/%d/index-%d.html", category.id, i)
			fmt.Println(listUrl)
			listDoc, _ := htmlquery.LoadURL(listUrl)
			liNodes := htmlquery.Find(listDoc, `//div[@class="bz_list"]/ul/li`)
			if len(liNodes) >= 1 {
				for _, liNode := range liNodes {
					fmt.Println(category.id, i, category.name)
					detailUrl := htmlquery.InnerText(htmlquery.FindOne(liNode, `./div[@class="bz_listl"]/ul[1]/a/@href`))
					fmt.Println(detailUrl)
					detailDoc, _ := htmlquery.LoadURL(detailUrl)
					downNodes := htmlquery.Find(detailDoc, `//div[@class="downk"]/a`)
					if len(downNodes) == 2 {
						title := htmlquery.InnerText(htmlquery.FindOne(detailDoc, `//div[@class="title2"]/span`))
						title = strings.ReplaceAll(title, "<font color=\"red\"></font>", "")
						title = strings.ReplaceAll(title, "/", "-")
						title = strings.ReplaceAll(title, " ", "")
						fmt.Println(title)

						authUrl := htmlquery.InnerText(htmlquery.FindOne(downNodes[1], `./@href`))
						fmt.Println(authUrl)
						// 获取请求Location
						downloadUrl, err := getFoodMateDownloadUrl(authUrl, detailUrl)
						if len(downloadUrl) == 0 {
							fmt.Println(err)
							continue
						}
						// 只下载pdf文件
						if strings.Index(downloadUrl, ".pdf") == -1 {
							fmt.Println("不是pdf文件")
							continue
						}
						fmt.Println(downloadUrl)
						filePath := "E:\\workspace\\down.foodmate.net\\" + category.name + "\\" + title + ".pdf"
						fmt.Println(filePath)
						if _, err := os.Stat(filePath); err != nil {
							fmt.Println("=======开始下载========")
							err = downloadFoodMatePdf(downloadUrl, filePath, detailUrl)
							if err != nil {
								fmt.Println(err)
							}
							fmt.Println("=======下载完成========")
						}
						downloadFoodMatePdfSleep := rand.Intn(20)
						for i := 1; i <= downloadFoodMatePdfSleep; i++ {
							time.Sleep(time.Second)
							fmt.Println("page="+strconv.Itoa(i)+"===========更新", title, "成功，暂停", downloadFoodMatePdfSleep, "秒，倒计时", i, "秒===========")
						}
					} else {
						continue
					}
				}
				i++
			} else {
				isPageGo = false
				i = 1
				break
			}
		}
	}
}

// 获取请求Location
func getFoodMateDownloadUrl(authUrl string, referer string) (downloadUrl string, err error) {
	// 初始化客户端
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("GET", authUrl, nil) //建立连接
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", FoodMateCookie)
	req.Header.Set("Host", "down.foodmate.net")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", referer)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return downloadUrl, err
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode == http.StatusOK {
		downloadUrl = authUrl
	} else if resp.StatusCode == http.StatusFound {
		downloadUrl = resp.Header.Get("Location")
	}
	return downloadUrl, nil
}

func downloadFoodMatePdf(pdfUrl string, filePath string, referer string) error {
	// 初始化客户端
	var client http.Client
	req, err := http.NewRequest("GET", pdfUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", FoodMateCookie)
	req.Header.Set("Host", "down.foodmate.net")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", referer)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
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
