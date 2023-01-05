package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	FoodMateEnableHttpProxy = false
	FoodMateHttpProxyUrl    = "27.42.168.46:55481"
)

func FoodMateSetHttpProxy() (httpclient http.Client) {
	ProxyURL, _ := url.Parse(FoodMateHttpProxyUrl)
	httpclient = http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type internalCategory struct {
	id   int
	name string
}

// foodMateSpider 获取标准库Pdf文档
// @Title 获取标准库Pdf文档
// @Description http://down.foodmate.net/，获取标准库Pdf文档
func main() {
	// 国内标准列表
	var allCategory = []internalCategory{
		{id: 1, name: "国内标准"},
		{id: 2, name: "国外标准"},
		//{id: 3, name: "国家标准"},
		//{id: 4, name: "进出口行业标准"},
		//{id: 5, name: "农业标准"},
		//{id: 6, name: "商业标准"},
		//{id: 7, name: "水产标准"},
		//{id: 8, name: "轻工标准"},
		//{id: 9, name: "其它国内标准"},
		//{id: 12, name: "团体标准"},
		//{id: 14, name: "医药标准"},
		//{id: 15, name: "地方标准"},
		//{id: 16, name: "卫生标准"},
		//{id: 17, name: "化工标准"},
		//{id: 18, name: "烟草标准"},
		//{id: 19, name: "食品安全企业标准"},
		//{id: 46, name: "认证认可标准"},
	}
	for _, category := range allCategory {
		i := 1042
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

						downloadUrl := htmlquery.InnerText(htmlquery.FindOne(downNodes[1], `./@href`))
						fmt.Println(downloadUrl)
						downloadUrlArray, err := url.Parse(downloadUrl)
						filePath := "../down.foodmate.net/" + downloadUrlArray.Query().Get("auth") + "-" + title + ".pdf"
						fmt.Println(filePath)
						err = downloadFoodMatePdf(downloadUrl, filePath)
						if err != nil {
							fmt.Println(err)
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

func downloadFoodMatePdf(pdfUrl string, filePath string) error {
	// 初始化客户端
	var client http.Client
	if FoodMateEnableHttpProxy {
		client = FoodMateSetHttpProxy()
	}
	req, err := http.NewRequest("GET", pdfUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "u_rdown=1; __gads=ID=c9ebc3f3415762dc:T=1672378976:S=ALNI_MY7EFr7ppeoi71xEG5hP7xZsXQSsA; __51cke__=; Hm_lvt_d4fdc0f0037bcbb9bf9894ffa5965f5e=1672378976,1672406598; Hm_lvt_2aeaa32e7cee3cfa6e2848083235da9f=1672378976,1672406598; __gpi=UID=00000b9a75e47353:T=1672378976:RT=1672406598:S=ALNI_MaIREweBK2pNUDh9n_qnPNzmrr32g; Hm_lpvt_d4fdc0f0037bcbb9bf9894ffa5965f5e=1672419008; __tins__1484185=%7B%22sid%22%3A%201672417339674%2C%20%22vd%22%3A%206%2C%20%22expires%22%3A%201672420815540%7D; __51laig__=25; Hm_lpvt_2aeaa32e7cee3cfa6e2848083235da9f=1672419017")
	req.Header.Set("Host", "down.foodmate.net")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "http://down.foodmate.net/")
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
