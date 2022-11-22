package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

func main() {
	match, _ := regexp.MatchString("^p([a-z]+)ch$", "peach")
	fmt.Println(match)

	r, _ := regexp.Compile("p([a-z]+)ch")

	fmt.Println(r.MatchString("peach"))

	fmt.Println(r.FindString("peach punch"))

	fmt.Println(r.FindStringIndex("peach punch"))

	fmt.Println(r.FindStringSubmatch("peach punch"))

	fmt.Println(r.FindStringSubmatchIndex("peach punch"))

	fmt.Println(r.FindAllString("peach punch pinch", -1))

	fmt.Println(r.FindAllStringIndex("peach punch pinch", -1))

	fmt.Println(r.FindAllString("peach punch pinch", 2))

	fmt.Println(r.Match([]byte("peach")))

	r = regexp.MustCompile("p([a-z]+)ch")
	fmt.Println(r)
	matches := r.FindAllSubmatch([]byte("a peach peach"), -1)
	//fmt.Println(len(matches))
	for _, m := range matches {
		fmt.Println(string(m[0]), string(m[1]))
	}

	text := `Hello 世界！123 Go.`
	reg := regexp.MustCompile(`[a-z]+`)             // 查找连续的小写字母
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // 输出结果["ello" "o"]

	reg = regexp.MustCompile(`[^a-z]+`)             // 查找连续的非小写字母
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["H" " 世界！123 G" "."]

	reg = regexp.MustCompile(`[\w]+`)               // 查找连续的单词字母
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello" "123" "Go"]

	reg = regexp.MustCompile(`[^\w\s]+`)            // 查找连续的非单词字母、非空白字符
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["世界！" "."]

	reg = regexp.MustCompile(`[[:upper:]]+`)        // 查找连续的大写字母
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["H" "G"]

	reg = regexp.MustCompile(`[[:^ascii:]]+`)       // 查找连续的非 ASCII 字符
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["世界！"]

	reg = regexp.MustCompile(`[\pP]+`)              // 查找连续的标点符号
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["！" "."]

	reg = regexp.MustCompile(`[\PP]+`)              // 查找连续的非标点符号字符
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello 世界" "123 Go"]

	reg = regexp.MustCompile(`[\p{Han}]+`)          // 查找连续的汉字
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["世界"]

	reg = regexp.MustCompile(`[\P{Han}]+`)          // 查找连续的非汉字字符
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello " "！123 Go."]

	reg = regexp.MustCompile(`Hello|Go`)            // 查找 Hello 或 Go
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello" "Go"]

	reg = regexp.MustCompile(`^H.*\s`)              // 查找行首以 H 开头，以空格结尾的字符串
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello 世界！123 "]

	reg = regexp.MustCompile(`(?U)^H.*\s`)          // 查找行首以 H 开头，以空白结尾的字符串（非贪婪模式）
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello "]

	reg = regexp.MustCompile(`(?i:^hello).*Go`)     // 查找以 hello 开头（忽略大小写），以 Go 结尾的字符串
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello 世界！123 Go"]

	reg = regexp.MustCompile(`\QGo.\E`)             // 查找 Go.
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Go."]

	reg = regexp.MustCompile(`(?U)^.* `)            // 查找从行首开始，以空格结尾的字符串（非贪婪模式）
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello "]

	reg = regexp.MustCompile(` [^ ]*$`)             // 查找以空格开头，到行尾结束，中间不包含空格字符串
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // [" Go."]

	reg = regexp.MustCompile(`(?U)\b.+\b`)          // 查找“单词边界”之间的字符串
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello" " 世界！" "123" " " "Go"]

	reg = regexp.MustCompile(`[^ ]{1,4}o`)          // 查找连续 1 次到 4 次的非空格字符，并以 o 结尾的字符串
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello" "Go"]

	reg = regexp.MustCompile(`(?:Hell|G)o`)         // 查找 Hello 或 Go
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello" "Go"]

	reg = regexp.MustCompile(`(Hell|G)o`)                     // 查找 Hello 或 Go，替换为 Hellooo、Gooo
	fmt.Printf("%q\n", reg.ReplaceAllString(text, "${n}ooo")) // "Hellooo 世界！123 Gooo."

	reg = regexp.MustCompile(`(Hello)(.*)(Go)`)              // 交换 Hello 和 Go
	fmt.Printf("%q\n", reg.ReplaceAllString(text, "$3$2$1")) // "Go 世界！123 Hello."

	//fmt.Println(r.ReplaceAllString("a peach", "<fruit>"))
	//
	//in := []byte("a peach")
	//out := r.ReplaceAllFunc(in, bytes.ToUpper)
	//fmt.Println(string(out))

	var filter1 = map[string]string{"name": "能源", "status": "bool"}
	fmt.Println(filter1)
	for k, v := range filter1 {
		println(k, v)
	}

	var filter = map[string]interface{}{"name": "能源", "status": true}
	fmt.Println(filter)
	for k, _ := range filter {
		fmt.Println(k, filter[k])
	}

	text = "CO2强化深部咸水开采与封存（CO2 enhanced water recovery and storage）"
	reg = regexp.MustCompile(`\(([^)]+)\)|（([^）]+)）`)
	regFindAllString := reg.FindAllString(text, -1)
	fmt.Println(regFindAllString)
	oldNew := make([]string, len(regFindAllString)*2)
	for _, v := range regFindAllString {
		fmt.Println(v)
		oldNew = append(oldNew, v, "")
	}
	fmt.Println(oldNew)
	replacer := strings.NewReplacer(oldNew...)
	fmt.Println(replacer.Replace(text))

	//text1 := `var itemPropertyValues = {"5926":{"item_name":"\u7164\u5236\u6cb9\u5e73\u5747\n\uff08Coal-to-liquids-average\uff09","click_count":640,"modify_time":"2021.12.08","modify_users":"","\u4e0a\u6e38\u6392\u653e\uff08Upstream emissions\uff09":{"property_id":67027,"subject":"10.55"},"\u4e0a\u6e38\u6392\u653e\u5355\u4f4d\uff08Unit\uff09":{"property_id":67028,"subject":"\u5428\u4e8c\u6c27\u5316\u78b3\u5f53\u91cf\/\u5428\n(t CO2-eq\/ t)"},"\u4e0b\u6e38\u6392\u653e\uff08Downstream emissions\uff09":{"property_id":67029,"subject":null},"\u4e0b\u6e38\u6392\u653e\u5355\u4f4d\uff08Unit\uff09":{"property_id":67030,"subject":null},"\u6392\u653e\u73af\u8282\uff08Emission processes\uff09":{"property_id":67031,"subject":"\u751f\u4ea7\uff1a10.55"},"\u6392\u653e\u6e29\u5ba4\u6c14\u4f53\u5360\u6bd4\uff08GHG percentage\uff09":{"property_id":67032,"subject":null},"\u6570\u636e\u65f6\u95f4\uff08Year\uff09":{"property_id":67033,"subject":"2013"},"\u4e0d\u786e\u5b9a\u6027\uff08Uncertainty\uff09":{"property_id":67034,"subject":null},"\u5176\u4ed6\uff08Others\uff09":{"property_id":67035,"subject":null},"\u53c2\u8003\u6587\u732e\/\u6570\u636e\u6765\u6e90\uff08Data source\uff09":{"property_id":67036,"subject":"[1] \u7164\u70ad\u79d1\u5b66\u6280\u672f\u7814\u7a76\u9662\u6709\u9650\u516c\u53f8.\u300a\u7164\u5316\u5de5\u5355\u4f4d\u4ea7\u54c1\u80fd\u6e90\u6d88\u8017\u9650\u989d\u300b.2020;\n[2] \u6731\u73b2,\u51af\u76f8\u662d,\u5b54\u4f73\u96ef.\u57fa\u4e8e\u751f\u547d\u5468\u671f\u8bc4\u4ef7LCA\u7684\u7164\u5236\u6cb9\u751f\u4ea7\u8fc7\u7a0b\u5206\u6790 [J]. \u6d01\u51c0\u7164\u6280\u672f, 2018,24(02):119-126."},"\u4fee\u6539\u4eba":{"property_id":0,"subject":""}},"5927":{"item_name":"\u7164\u5236\u6cb9\uff08\u76f4\u63a5\u6db2\u5316\uff09\n\uff08Direct-liquefied coal-to-liquids\uff09","click_count":299,"modify_time":"2021.12.08","modify_users":"","\u4e0a\u6e38\u6392\u653e\uff08Upstream emissions\uff09":{"property_id":67037,"subject":"5.8"},"\u4e0a\u6e38\u6392\u653e\u5355\u4f4d\uff08Unit\uff09":{"property_id":67038,"subject":"\u5428\u4e8c\u6c27\u5316\u78b3\u5f53\u91cf\/\u5428\n(t CO2-eq\/ t)"},"\u4e0b\u6e38\u6392\u653e\uff08Downstream emissions\uff09":{"property_id":67039,"subject":null},"\u4e0b\u6e38\u6392\u653e\u5355\u4f4d\uff08Unit\uff09":{"property_id":67040,"subject":null},"\u6392\u653e\u73af\u8282\uff08Emission processes\uff09":{"property_id":67041,"subject":"\u751f\u4ea7\uff1a5.80"},"\u6392\u653e\u6e29\u5ba4\u6c14\u4f53\u5360\u6bd4\uff08GHG percentage\uff09":{"property_id":67042,"subject":null},"\u6570\u636e\u65f6\u95f4\uff08Year\uff09":{"property_id":67043,"subject":null},"\u4e0d\u786e\u5b9a\u6027\uff08Uncertainty\uff09":{"property_id":67044,"subject":null},"\u5176\u4ed6\uff08Others\uff09":{"property_id":67045,"subject":null},"\u53c2\u8003\u6587\u732e\/\u6570\u636e\u6765\u6e90\uff08Data source\uff09":{"property_id":67046,"subject":"\u7164\u70ad\u79d1\u5b66\u6280\u672f\u7814\u7a76\u9662\u6709\u9650\u516c\u53f8.\u300a\u7164\u5316\u5de5\u5355\u4f4d\u4ea7\u54c1\u80fd\u6e90\u6d88\u8017\u9650\u989d\u300b.2020"},"\u4fee\u6539\u4eba":{"property_id":0,"subject":""}},"5928":{"item_name":"\u7164\u5236\u6cb9\uff08\u95f4\u63a5\u6db2\u5316\uff09\n\uff08Indirectly liquefied coal-to-liquids\uff09","click_count":251,"modify_time":"2021.12.08","modify_users":"","\u4e0a\u6e38\u6392\u653e\uff08Upstream emissions\uff09":{"property_id":67047,"subject":"6.35"},"\u4e0a\u6e38\u6392\u653e\u5355\u4f4d\uff08Unit\uff09":{"property_id":67048,"subject":"\u5428\u4e8c\u6c27\u5316\u78b3\u5f53\u91cf\/\u5428\n(t CO2-eq\/ t)"},"\u4e0b\u6e38\u6392\u653e\uff08Downstream emissions\uff09":{"property_id":67049,"subject":null},"\u4e0b\u6e38\u6392\u653e\u5355\u4f4d\uff08Unit\uff09":{"property_id":67050,"subject":null},"\u6392\u653e\u73af\u8282\uff08Emission processes\uff09":{"property_id":67051,"subject":"\u751f\u4ea7\uff1a6.35"},"\u6392\u653e\u6e29\u5ba4\u6c14\u4f53\u5360\u6bd4\uff08GHG percentage\uff09":{"property_id":67052,"subject":null},"\u6570\u636e\u65f6\u95f4\uff08Year\uff09":{"property_id":67053,"subject":null},"\u4e0d\u786e\u5b9a\u6027\uff08Uncertainty\uff09":{"property_id":67054,"subject":null},"\u5176\u4ed6\uff08Others\uff09":{"property_id":67055,"subject":null},"\u53c2\u8003\u6587\u732e\/\u6570\u636e\u6765\u6e90\uff08Data source\uff09":{"property_id":67056,"subject":"\u7164\u70ad\u79d1\u5b66\u6280\u672f\u7814\u7a76\u9662\u6709\u9650\u516c\u53f8.\u300a\u7164\u5316\u5de5\u5355\u4f4d\u4ea7\u54c1\u80fd\u6e90\u6d88\u8017\u9650\u989d\u300b.2020"},"\u4fee\u6539\u4eba":{"property_id":0,"subject":""}},"5929":{"item_name":"\u5176\u4ed6\u7164\u5236\u6cb9\n\uff08Other coal-to-liquids\uff09","click_count":216,"modify_time":"2021.12.08","modify_users":"","\u4e0a\u6e38\u6392\u653e\uff08Upstream emissions\uff09":{"property_id":67057,"subject":"19.51"},"\u4e0a\u6e38\u6392\u653e\u5355\u4f4d\uff08Unit\uff09":{"property_id":67058,"subject":"\u5428\u4e8c\u6c27\u5316\u78b3\u5f53\u91cf\/\u5428\n(t CO2-eq\/ t)"},"\u4e0b\u6e38\u6392\u653e\uff08Downstream emissions\uff09":{"property_id":67059,"subject":null},"\u4e0b\u6e38\u6392\u653e\u5355\u4f4d\uff08Unit\uff09":{"property_id":67060,"subject":null},"\u6392\u653e\u73af\u8282\uff08Emission processes\uff09":{"property_id":67061,"subject":"\u7164\u70ad\u5f00\u91c7\u52a0\u5de5\uff1a10.51\uff1b\u7164\u70ad\u8fd0\u8f93\uff1a0.19\uff1b\u5de5\u5382\u52a0\u5de5\uff1a8.81"},"\u6392\u653e\u6e29\u5ba4\u6c14\u4f53\u5360\u6bd4\uff08GHG percentage\uff09":{"property_id":67062,"subject":"CO2\uff1a0.61\uff1bCH4\uff1a0.38\uff1bN2O\uff1a0.01"},"\u6570\u636e\u65f6\u95f4\uff08Year\uff09":{"property_id":67063,"subject":"2013"},"\u4e0d\u786e\u5b9a\u6027\uff08Uncertainty\uff09":{"property_id":67064,"subject":"\uff081\uff09\u5178\u578b\u4e2a\u6848"},"\u5176\u4ed6\uff08Others\uff09":{"property_id":67065,"subject":null},"\u53c2\u8003\u6587\u732e\/\u6570\u636e\u6765\u6e90\uff08Data source\uff09":{"property_id":67066,"subject":"\u6731\u73b2,\u51af\u76f8\u662d,\u5b54\u4f73\u96ef.\u57fa\u4e8e\u751f\u547d\u5468\u671f\u8bc4\u4ef7LCA\u7684\u7164\u5236\u6cb9\u751f\u4ea7\u8fc7\u7a0b\u5206\u6790 [J]. \u6d01\u51c0\u7164\u6280\u672f, 2018,24(02):119-126."},"\u4fee\u6539\u4eba":{"property_id":0,"subject":""}}};`
	//regItemPropertyOptions := regexp.MustCompile(`var itemPropertyValues = ([{](.*?)[}]{3})`)
	//text1 = strings.Replace(text1, "\"subject\":null", "\"subject\":\"\"", -1)
	//fmt.Println(text1)
	//regItemPropertyOptionsFindAllSubMatch := regItemPropertyOptions.FindAllSubmatch([]byte(text1), -1)
	//for _, itemsMatch := range regItemPropertyOptionsFindAllSubMatch {
	//	var items map[string]interface{}
	//	_ = json.Unmarshal(itemsMatch[1], &items)
	//	for _, item := range items {
	//		fmt.Println(item)
	//		//itemMap := item.(map[string]interface{})
	//		//fmt.Println(itemMap)
	//		fmt.Println("========")
	//		//fmt.Println(itemMap["click_count"].(float64))
	//		//fmt.Println(itemMap["item_name"].(string))
	//		//fmt.Println(itemMap["modify_time"].(string))
	//		//fmt.Println(itemMap["modify_users"].(string))
	//		//UpstreamEmissions := itemMap["上游排放单位（Unit）"].(map[string]interface{})
	//		//fmt.Println(UpstreamEmissions["subject"].(string))
	//	}
	//}

	//UpstreamEmissionsUnit := "吨二氧化碳当量/吨 （t CO2-eq/ t）"
	//regUpstreamEmissionsUnit := regexp.MustCompile(`(.*?)([(（](.*?)[)）])`)
	//regUpstreamEmissionsUnitFindAllSubMatch := regUpstreamEmissionsUnit.FindAllSubmatch([]byte(UpstreamEmissionsUnit), -1)
	//for _, match := range regUpstreamEmissionsUnitFindAllSubMatch {
	//	fmt.Println(string(match[0]), string(match[1]), string(match[2]))
	//}

	//str := "家具产品(Furniture&nbsp;Products)"
	//reg = regexp.MustCompile(`\([\w\s\-,\/&nbsp;]+\)|（[\w\s\-,\/&nbsp;]+）|\([\w\s\-,\/&nbsp;]+）|（[\w\s\-,\/&nbsp;]+\)`)
	//regFindAllString = reg.FindAllString(str, -1)
	//fmt.Println(regFindAllString)
	//regReplacerOldNew := make([]string, len(regFindAllString)*2)
	//for _, v := range regFindAllString {
	//	regReplacerOldNew = append(regReplacerOldNew, v, "")
	//}
	//regReplacer := strings.NewReplacer(regReplacerOldNew...)
	//fmt.Println(strings.TrimSpace(regReplacer.Replace(str)))

	//str := "  dfgh  hj  "
	//str = strings.TrimSpace(str)
	//fmt.Println(str)

	//	text = `
	//                            泥炭及泥炭产物
	//（Peat）                            (1)`
	//	oldNew = []string{"\n", "", "\r", ""}
	//	replacer = strings.NewReplacer(oldNew...)
	//	level1FinalTitle := strings.TrimSpace(replacer.Replace(text))
	//	reg = regexp.MustCompile(`(.*?)\([0-9]+\)`)
	//	regFindStringSubMatch := reg.FindStringSubmatch(level1FinalTitle)
	//	fmt.Println(regFindStringSubMatch[1])

	//level1FinalTitle := "\n                            泥炭及泥炭产物\n（Peat）                            (1)"
	//
	//level1FinalTitleOldNew := []string{"\n", "", "\r", ""}
	//level1FinalTitleReplacer := strings.NewReplacer(level1FinalTitleOldNew...)
	//level1FinalTitle = level1FinalTitleReplacer.Replace(level1FinalTitle)
	//fmt.Println(level1FinalTitle)
	//
	//regLevel1FinalTitle := regexp.MustCompile(`(.*?)\([0-9]+\)`)
	//regFindStringSubMatch := regLevel1FinalTitle.FindStringSubmatch(level1FinalTitle)
	//level1FinalTitle = strings.TrimSpace(regFindStringSubMatch[1])
	//fmt.Println(level1FinalTitle)

	str := "ni124534"
	for _, r := range str {
		fmt.Println(unicode.IsDigit(r))
	}

}
