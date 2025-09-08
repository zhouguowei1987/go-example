package main

import (
	"errors"
	"fmt"
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

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const (
	DbSpPtEnableHttpProxy = false
	DbSpPtHttpProxyUrl    = "111.225.152.186:8089"
)

func DbSpPtSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(DbSpPtHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type DownloadDbSpPtFormData struct {
	task      string
	guid      string
	file_guid string
	num_tn    string
	type_temp string
	fact_name string
	filePath  string
	title_tip string
}

type DbSpPtAllData struct {
	AllData []DbSpPtData
}
type DbSpPtData struct {
	province      string
	id            string
	title         string
	rn            string
	standard_code string
}

var dbSpPtAllData = DbSpPtAllData{
	AllData: []DbSpPtData{
		{
			province:      "甘肃",
			id:            "4508E128-8D5A-4EE9-9186-3DD81489CA42",
			title:         "食品安全地方标准  肉苁蓉生产卫生规范",
			rn:            "1",
			standard_code: "DB S62/012-2021",
		}, {
			province:      "甘肃",
			id:            "10266EE8-7594-4D23-8F0A-2D1CF002284E",
			title:         "食品安全地方标准  党参生产卫生规范",
			rn:            "2",
			standard_code: "DB S62/010-2021",
		}, {
			province:      "重庆",
			id:            "84BC5745-4DDB-4D03-B048-B2FC54DC25D7",
			title:         "食品安全地方标准 山银花及其制品",
			rn:            "3",
			standard_code: "DB S50/31-2023",
		}, {
			province:      "山西",
			id:            "461F44FE-B32A-49F9-9EB7-9172DA80686F",
			title:         "食品安全地方标准 翅果仁",
			rn:            "4",
			standard_code: "DB S14/002-2020",
		}, {
			province:      "甘肃",
			id:            "580F2925-4380-4B6F-9B32-44548EE9C0DA",
			title:         "食品安全地方标准  苦水玫瑰",
			rn:            "5",
			standard_code: "DB S62/002-2020(补充）",
		}, {
			province:      "吉林",
			id:            "DB05BCBA-8492-4778-AEB0-80C7366BF1F9",
			title:         "食品安全地方标准 蓝靛果",
			rn:            "6",
			standard_code: "DB S22/038-2024",
		}, {
			province:      "广东",
			id:            "C866A28B-0AC5-452B-94DD-B3086CDA122C",
			title:         "食品安全地方标准 汕头牛肉丸",
			rn:            "7",
			standard_code: "DB S44/005-2024",
		}, {
			province:      "上海",
			id:            "8BB5321B-FBBF-4A56-BCE5-7AFEB7767130",
			title:         "《食品安全地方标准 现制饮料》",
			rn:            "8",
			standard_code: "DB 31/2007-2023",
		}, {
			province:      "上海",
			id:            "35850F53-37F8-425B-B769-3C10FF817A67",
			title:         "《食品安全地方标准 即食食品现场制售卫生规范》",
			rn:            "9",
			standard_code: "DB 31/2027-2023",
		}, {
			province:      "上海",
			id:            "545F500D-140A-48E4-9C2F-5C600070CC1F",
			title:         "食品安全地方标准 发酵肉制品生产卫生规范",
			rn:            "10",
			standard_code: "DB 31/2017-2013",
		}, {
			province:      "浙江",
			id:            "EB24B227-98B4-49C0-93E2-9E63331781BE",
			title:         "食品安全地方标准 藕粉生产卫生规范",
			rn:            "11",
			standard_code: "DB S33/3018-2024",
		}, {
			province:      "广东",
			id:            "B3FE6BD6-0984-4260-90C8-1FC930D2F375",
			title:         "食品安全地方标准 化橘红胎",
			rn:            "12",
			standard_code: "DB S44/021-2023",
		}, {
			province:      "广东",
			id:            "78764404-C97B-40D9-A8A4-6BE8A07221F6",
			title:         "食品安全地方标准 鱼露生产卫生规范",
			rn:            "13",
			standard_code: "DB S44/019-2023",
		}, {
			province:      "广东",
			id:            "3D1701B1-6951-4D9F-B6F8-C50956144424",
			title:         "食品安全地方标准 鸡蛋花",
			rn:            "14",
			standard_code: "DB S44/018-2022",
		}, {
			province:      "湖南",
			id:            "50CDB23A-FA0B-42BF-A30A-85E5F075A459",
			title:         "食品安全地方标准 牡丹木槿花",
			rn:            "15",
			standard_code: "DB S43/017-2024",
		}, {
			province:      "山西",
			id:            "EFEF0726-8FD6-485A-ABE2-172566F0058F",
			title:         "食品安全地方标准 连翘叶代用茶",
			rn:            "16",
			standard_code: "DB S14/009-2024",
		}, {
			province:      "甘肃",
			id:            "A24C6147-639C-4E55-9B24-2A9F93700AB7",
			title:         "食品安全地方标准  黄芪生产卫生规范",
			rn:            "17",
			standard_code: "DB S62/011-2021",
		}, {
			province:      "山西",
			id:            "9234A23B-1C45-438E-9FA0-A735037B3891",
			title:         "食品安全地方标准 黄芩叶",
			rn:            "18",
			standard_code: "DB S14/007-2024",
		}, {
			province:      "河北",
			id:            "E898D7A1-3745-499E-885D-1AAD583E2997",
			title:         "食品安全地方标准 预制李鸿章烩菜",
			rn:            "19",
			standard_code: "DB S13/022-2024",
		}, {
			province:      "河北",
			id:            "A6A39E9A-80B9-4D81-8756-244D52E494B9",
			title:         "食品安全地方标准 保定驴肉火烧",
			rn:            "20",
			standard_code: "DB S13/021-2024",
		}, {
			province:      "上海",
			id:            "77F2DB63-AFDE-463F-A47A-57AC9BF749AD",
			title:         "《食品安全地方标准 集体用餐配送膳食》",
			rn:            "21",
			standard_code: "DB 31/2023-2023",
		}, {
			province:      "上海",
			id:            "95BD6171-6F82-49E2-BC37-AF3A2CED9700",
			title:         "《食品安全地方标准 集体用餐配送膳食生产配送卫生规范》",
			rn:            "22",
			standard_code: "DB 31/2024-2023",
		}, {
			province:      "天津",
			id:            "47E49197-C626-4CEE-88C2-617694626D99",
			title:         "食品安全地方标准 集体用餐配送膳食",
			rn:            "23",
			standard_code: "DB S12/004—2024",
		}, {
			province:      "重庆",
			id:            "A12E56C2-6BFF-4385-8EAC-4B02609F3A74",
			title:         "食品安全地方标准 丰都麻辣鸡",
			rn:            "24",
			standard_code: "DB S50/034—2024",
		}, {
			province:      "西藏",
			id:            "D9A19BE9-951D-4891-BC29-DAD43EA82A3F",
			title:         "食品安全地方标准 酿造用藏曲",
			rn:            "25",
			standard_code: "DB S54/2006-2024",
		}, {
			province:      "广西",
			id:            "65C37F80-A4C3-40F3-86DF-99321342FCCC",
			title:         "食品安全地方标准 调制水牛乳",
			rn:            "26",
			standard_code: "DB S45/046-2024",
		}, {
			province:      "广西",
			id:            "F1FBFDF1-F0D8-4F9B-A3EA-E4822E6793C1",
			title:         "食品安全地方标准 灭菌水牛乳",
			rn:            "27",
			standard_code: "DB S45/037-2024",
		}, {
			province:      "广西",
			id:            "5D38BCCA-3FB1-4A67-922B-1F9246CE1A7B",
			title:         "食品安全地方标准 发酵水牛乳",
			rn:            "28",
			standard_code: "DB S45/024-2024",
		}, {
			province:      "广西",
			id:            "904F83CA-6A6B-4DED-B739-2562A12A7ED1",
			title:         "食品安全地方标准 生水牛乳",
			rn:            "29",
			standard_code: "DB S45/011-2024",
		}, {
			province:      "广西",
			id:            "1C1F9BEC-F8E9-45AB-985A-824D287FCA2C",
			title:         "食品安全地方标准 巴氏杀菌水牛乳",
			rn:            "30",
			standard_code: "DB S45/012-2024",
		}, {
			province:      "重庆",
			id:            "1175E964-DCA0-4538-9F25-155489D38E47",
			title:         "食品安全地方标准 保鲜花椒",
			rn:            "31",
			standard_code: "DB S50/003-2024",
		}, {
			province:      "上海",
			id:            "1E9D847F-FE2D-43F4-8103-6564A14F9418",
			title:         "食品安全地方标准 食品生产加工小作坊卫生规范",
			rn:            "32",
			standard_code: "DB 31/2019-2013",
		}, {
			province:      "广东",
			id:            "99F798A6-6B84-4953-BF8A-B0866B3E012A",
			title:         "食品安全地方标准 食品生产加工小作坊卫生规范",
			rn:            "33",
			standard_code: "DB S44/020-2023",
		}, {
			province:      "天津",
			id:            "DCD8FF7F-04C1-4CF0-93FA-C2D500BA3765",
			title:         "食品安全地方标准 冷藏即食食品生产卫生规范",
			rn:            "34",
			standard_code: "DB S12/003—2024",
		}, {
			province:      "天津",
			id:            "7C2CF15D-4151-4509-A517-9C68AAC27B97",
			title:         "食品安全地方标准 食品生产加工小作坊食品安全控制基本要求",
			rn:            "35",
			standard_code: "DB S12/002—2024",
		}, {
			province:      "青海",
			id:            "9678FA01-1ED6-4345-9B0D-1A6EDCD0F300",
			title:         "食品安全地方标准 食品生产加工小作坊卫生规范",
			rn:            "36",
			standard_code: "DB S63/0004—2024",
		}, {
			province:      "青海",
			id:            "D0977999-3603-4673-B511-506E44410AEF",
			title:         "食品安全地方标准 黑果枸杞",
			rn:            "37",
			standard_code: "DB S63/0010-2024",
		}, {
			province:      "黑龙江",
			id:            "5D7433EB-2EBB-4984-A9D0-6D0623CC232A",
			title:         "食品安全地方标准  淀粉制品小作坊生产卫生规范",
			rn:            "38",
			standard_code: "DB S23/025-2024",
		}, {
			province:      "黑龙江",
			id:            "36584CCF-3686-448D-93E5-04232B9BE151",
			title:         "食品安全地方标准 偃松籽",
			rn:            "39",
			standard_code: "DB S23/024-2023",
		}, {
			province:      "河南",
			id:            "D825C5CE-94CC-44CE-88F7-2792B7F4EEA5",
			title:         "食品安全地方标准 方便胡辣汤",
			rn:            "40",
			standard_code: "DB S41/006-2023",
		}, {
			province:      "贵州",
			id:            "A68337E6-3C32-4ADC-B44E-05212B9B25F3",
			title:         "食品安全地方标准 魔芋凝胶食品",
			rn:            "41",
			standard_code: "DB S52/077—2024",
		}, {
			province:      "贵州",
			id:            "4756B693-40F8-42C5-B513-7FBC2A0C6808",
			title:         "食品安全地方标准 酸菜蹄膀",
			rn:            "42",
			standard_code: "DB S52/076—2024",
		}, {
			province:      "贵州",
			id:            "594A240A-622F-4324-AFF4-4D6C3BD3DA8D",
			title:         "食品安全地方标准 辣子鸡",
			rn:            "43",
			standard_code: "DB S52/001—2024",
		}, {
			province:      "贵州",
			id:            "9C9817B1-C8CD-41BA-83A7-98EC92B15A4E",
			title:         "食品安全地方标准 黄粑",
			rn:            "44",
			standard_code: "DB S52/070—2023",
		}, {
			province:      "贵州",
			id:            "7DBF88CE-B15E-4F9C-B1A6-B0922CD6C069",
			title:         "食品安全地方标准 火锅底料",
			rn:            "45",
			standard_code: "DB S52/071—2023",
		}, {
			province:      "贵州",
			id:            "449F9CEC-27F2-4B0D-ADDC-75011772528A",
			title:         "食品安全地方标准 刺梨原汁",
			rn:            "46",
			standard_code: "DB S52/073—2023",
		}, {
			province:      "贵州",
			id:            "6F89079A-FD60-4E93-9BC5-B0731A5641CE",
			title:         "食品安全地方标准 香酥辣椒生产卫生规范",
			rn:            "47",
			standard_code: "DB S52/075—2023",
		}, {
			province:      "江西",
			id:            "4C07F7CB-F178-4B85-A80D-1DE09479E819",
			title:         "食品安全地方标准  九层皮（米糕）生产卫生规范 ",
			rn:            "48",
			standard_code: "DB 36/1683-2022",
		}, {
			province:      "重庆",
			id:            "6BDF5299-2D90-473D-8947-F5AC49F23601",
			title:         "食品安全地方标准 麻辣调料",
			rn:            "49",
			standard_code: "DB S50/021-2021",
		}, {
			province:      "重庆",
			id:            "01144047-5751-4102-8FA7-5AA45CA3C3E0",
			title:         "食品安全地方标准 火锅底料 ",
			rn:            "50",
			standard_code: "DB S50/022-2021",
		}, {
			province:      "重庆",
			id:            "52FF18FA-37EC-4D06-BFA1-B2782385A73F",
			title:         "食品安全地方标准 泡菜类调料",
			rn:            "51",
			standard_code: "DB S50/020-2021",
		}, {
			province:      "江苏",
			id:            "17957633-08F7-4E89-94B4-EB3C24D993F6",
			title:         "食品安全地方标准 牛蒡根制品",
			rn:            "52",
			standard_code: "DB S32/019—2021",
		}, {
			province:      "江苏",
			id:            "0C7C3869-CB7F-4CA2-9DEA-2609B3C62402",
			title:         "食品安全地方标准 方便菜肴",
			rn:            "53",
			standard_code: "DB S32/005-2021",
		}, {
			province:      "黑龙江",
			id:            "56487263-1C3C-4BE6-9F85-1CEEF011ED51",
			title:         "食品安全地方标准 龙江小烧酒小作坊生产卫生规范",
			rn:            "54",
			standard_code: "DB S23/009-2019",
		}, {
			province:      "新疆",
			id:            "60F752CF-A4B9-452E-B637-25F259670C64",
			title:         "食品安全地方标准 发酵乳粉",
			rn:            "55",
			standard_code: "DB S/65020-2023",
		}, {
			province:      "新疆",
			id:            "AA1BB418-D100-470E-AB23-CCD95D12A44D",
			title:         "食品安全地方标准 巴氏杀菌驴乳",
			rn:            "56",
			standard_code: "DB S65/018-2023",
		}, {
			province:      "新疆",
			id:            "BBD744EC-3CC7-421A-8276-CA32E563784C",
			title:         "食品安全地方标准 生驴乳",
			rn:            "57",
			standard_code: "DB S65/017-2023",
		}, {
			province:      "新疆",
			id:            "F0EB4742-9264-440E-AE32-2F7C327AFFFA",
			title:         "食品安全地方标准 生马乳",
			rn:            "58",
			standard_code: "DB S65/015-2023",
		}, {
			province:      "新疆",
			id:            "35AD3BC2-3643-47E3-B3A1-33CBFAABAB06",
			title:         "食品安全地方标准 发酵驼乳",
			rn:            "59",
			standard_code: "DB S65/013-2023",
		}, {
			province:      "云南",
			id:            "06E3756F-7CC3-4A28-BAA3-0E6DD3B45A5C",
			title:         "食品安全地方标准  米线、卷粉、饵丝（块）",
			rn:            "60",
			standard_code: "DB S53/017-2023",
		}, {
			province:      "新疆",
			id:            "2345AC16-E782-46C1-8EE1-1BE2D38BA22A",
			title:         "食品安全地方标准 巴氏杀菌驼乳",
			rn:            "61",
			standard_code: "DB S65/011-2023",
		}, {
			province:      "新疆",
			id:            "C276047E-A4BD-4D5B-A3E5-7F415E446F30",
			title:         "食品安全地方标准 生驼乳",
			rn:            "62",
			standard_code: "DB S65/010-2023",
		}, {
			province:      "黑龙江",
			id:            "4F3E456C-C63A-4694-B493-BDDCEFBC84D9",
			title:         "食品安全地方标准 冷冻粘豆包小作坊生产卫生规范",
			rn:            "63",
			standard_code: "DB S23/022-2023",
		}, {
			province:      "黑龙江",
			id:            "281E3F6D-CCE1-43D7-8D62-4069D4E15791",
			title:         "食品安全地方标准 糖葫芦小作坊生产卫生规范",
			rn:            "64",
			standard_code: "DB S23/021-2023",
		}, {
			province:      "黑龙江",
			id:            "9E67E690-9517-49D4-9DFA-CB4E31C6D965",
			title:         "食品安全地方标准  熟肉制品小作坊生产卫生规范",
			rn:            "65",
			standard_code: "DB S23/020-2023",
		}, {
			province:      "黑龙江",
			id:            "77869CFA-F650-44E5-BFED-6FC5C4483F20",
			title:         "食品安全地方标准 酱腌菜小作坊生产卫生规范",
			rn:            "66",
			standard_code: "DB S23/019-2023",
		}, {
			province:      "广西",
			id:            "1293DF10-2B35-4D82-B1C7-610CF3CEE18F",
			title:         "食品安全地方标准 金槐米",
			rn:            "67",
			standard_code: "DB S45/078-2022",
		}, {
			province:      "广西",
			id:            "841AC5A5-E97C-41A7-8380-79EE61895CA6",
			title:         "食品安全地方标准 柳州螺蛳粉",
			rn:            "68",
			standard_code: "DB S45/034-2021",
		}, {
			province:      "湖南",
			id:            "B39B7329-63D3-4D35-BAB0-35D4C8A04E60",
			title:         "食品安全地方标准 中央厨房卫生规范",
			rn:            "69",
			standard_code: "DB S43/015-2023",
		}, {
			province:      "湖南",
			id:            "FE8B34A7-B457-497F-8ABE-F9B01426CC7B",
			title:         "食品安全地方标准 茯苓",
			rn:            "70",
			standard_code: "DB S43/014-2022",
		}, {
			province:      "湖南",
			id:            "299F4AC4-1324-4F77-B222-DB7348D009BB",
			title:         "食品安全地方标准 猪血丸子",
			rn:            "71",
			standard_code: "DB S43/012-2022",
		}, {
			province:      "湖南",
			id:            "A6E9CE22-B955-4E3C-99D3-189EF85C8C52",
			title:         "食品安全地方标准 预制长沙臭豆腐生产卫生规范",
			rn:            "72",
			standard_code: "DB S43/011-2022",
		}, {
			province:      "河北",
			id:            "8995C14F-C333-4E95-913D-7F0410C99317",
			title:         "食品安全地方标准 干制文冠果叶（花）",
			rn:            "73",
			standard_code: "DB S13/017-2023",
		}, {
			province:      "海南",
			id:            "A41940E1-6D6C-4D84-B036-9C954DCB7FC7",
			title:         "食品安全地方标准 香露兜叶（粉）",
			rn:            "74",
			standard_code: "DB S46/004—2022",
		}, {
			province:      "宁夏",
			id:            "A8D18244-EB0F-4E50-8F34-5536DE02D155",
			title:         "食品安全地方标准 食品生产加工小作坊通用卫生规范",
			rn:            "75",
			standard_code: "DB S64/009-2022",
		}, {
			province:      "黑龙江",
			id:            "A43C34C0-34DE-4FF9-B0BC-B53412A481FA",
			title:         "食品安全地方标准 食品小作坊通用卫生规范",
			rn:            "76",
			standard_code: "DB S23/018-2022",
		}, {
			province:      "浙江",
			id:            "ECA61049-83B9-4EDD-A925-1FBFFBB27B88",
			title:         "食品安全地方标准 即食发酵火腿生产经营卫生规范",
			rn:            "77",
			standard_code: "DB S33/3014-2022",
		}, {
			province:      "江西",
			id:            "E0059792-1FE4-4858-9B9A-CFAE99EDC1C7",
			title:         "食品安全地方标准 搓菜生产卫生规范",
			rn:            "78",
			standard_code: "DB 36/1682-2022",
		}, {
			province:      "贵州",
			id:            "DF8063CB-DB8B-4973-8088-CD3D54284585",
			title:         "食品安全地方标准  食用畜禽血制品加工卫生规范",
			rn:            "79",
			standard_code: "DB S52/067-2022",
		}, {
			province:      "贵州",
			id:            "571FA796-0A0F-44DB-A5C6-4B282E3A372F",
			title:         "食品安全地方标准 米豆腐",
			rn:            "80",
			standard_code: "DB S52/030-2022",
		}, {
			province:      "宁夏",
			id:            "B5FE05D3-712F-4046-BFB5-D099B2CFB98E",
			title:         "食品安全地方标准  枸杞原浆",
			rn:            "81",
			standard_code: "DB S64/008-2022",
		}, {
			province:      "甘肃",
			id:            "1B3D1A38-0313-49C5-B0AD-A7B552C58D80",
			title:         "食品安全地方标准  当归",
			rn:            "82",
			standard_code: "DB S62/001-2022",
		}, {
			province:      "广东",
			id:            "1932FCCF-77AD-46A1-AF7B-B7249E734B24",
			title:         "食品安全地方标准 湿米粉生产和经营卫生规范",
			rn:            "83",
			standard_code: "DB S44/017-2021",
		}, {
			province:      "山东",
			id:            "5F755CD3-02E6-415E-977F-3AD8563062F7",
			title:         "食品安全地方标准 食品小作坊生产加工卫生规范",
			rn:            "84",
			standard_code: "DB S37/002-2022",
		}, {
			province:      "湖北",
			id:            "833C06B8-9C4A-447D-A93B-2579590A1BDB",
			title:         "食品安全地方标准 孝感麻糖生产卫生规范",
			rn:            "85",
			standard_code: "DB S42/011-2022",
		}, {
			province:      "浙江",
			id:            "DDEFFA5E-7A38-4197-A515-D5DC7AA57C2B",
			title:         "食品安全地方标准 酥饼生产卫生规范",
			rn:            "86",
			standard_code: "DB S33/3013-2022",
		}, {
			province:      "云南",
			id:            "37CB2F00-5D8E-49FF-870E-192562AB5CE0",
			title:         "食品安全地方标准 过桥米线餐饮加工卫生规范",
			rn:            "87",
			standard_code: "DB S53/032—2022",
		}, {
			province:      "上海",
			id:            "2AA506D0-C4D0-436A-B5EE-1366C2755192",
			title:         "食品安全地方标准 青团 第1号修改单",
			rn:            "88",
			standard_code: "DB 31/2001-2012",
		}, {
			province:      "上海",
			id:            "AAC992B7-8B40-4172-9C8A-EAD807A97E07",
			title:         "食品安全地方标准 预包装冷藏膳食生产经营卫生规范",
			rn:            "89",
			standard_code: "DB 31/2026—2021",
		}, {
			province:      "黑龙江",
			id:            "4EDB0A92-C02A-44E4-995D-5E11CD0794E9",
			title:         "食品安全地方标准 小油坊生产卫生规范",
			rn:            "90",
			standard_code: "DB S23/013-2021",
		}, {
			province:      "黑龙江",
			id:            "9343CB2A-B0BA-4F55-A9C3-1E65ADC9F05D",
			title:         "食品安全地方标准 糕点小作坊生产卫生规范",
			rn:            "91",
			standard_code: "DB S23/015-2021",
		}, {
			province:      "黑龙江",
			id:            "BB48A817-0BDD-4CC8-8D94-6AEC0D97EDEC",
			title:         "食品安全地方标准 豆制品小作坊生产卫生规范",
			rn:            "92",
			standard_code: "DB S23/014-2021",
		}, {
			province:      "黑龙江",
			id:            "93E64C20-1A4C-42C2-BF9E-28FB52401F05",
			title:         "食品安全地方标准 生湿面制品小作坊生产卫生规范",
			rn:            "93",
			standard_code: "DB S23/012-2021",
		}, {
			province:      "海南",
			id:            "32B57EB0-B9E6-4D03-924E-49C887DAC498",
			title:         "食品安全地方标准 海南黄花梨叶",
			rn:            "94",
			standard_code: "DB S46/003-2021",
		}, {
			province:      "吉林",
			id:            "53870800-F7E9-4403-AF8C-18E60B018CA0",
			title:         "食品安全地方标准 葵花盘",
			rn:            "95",
			standard_code: "DB S22/036-2021",
		}, {
			province:      "吉林",
			id:            "1262B9B4-67C6-4560-A9C1-BB13E6635D34",
			title:         "食品安全地方标准 刺五加鲜叶",
			rn:            "96",
			standard_code: "DB S22/035—2019",
		}, {
			province:      "西藏",
			id:            "BA94FF08-761B-4E49-9491-5FC634D109C4",
			title:         "食品安全地方标准 珠芽蓼果实粉",
			rn:            "97",
			standard_code: "DB S54/2003-2021",
		}, {
			province:      "新疆",
			id:            "1137F293-94A1-49ED-8B6A-8DCC45ECEA0D",
			title:         "食品安全地方标准 馕",
			rn:            "98",
			standard_code: "DB S65/022-2021",
		}, {
			province:      "青海",
			id:            "D4CD4601-82F2-43C8-8FB5-515E3206C64C",
			title:         "食品安全地方标准 黑果枸杞中花青素含量的测定",
			rn:            "99",
			standard_code: "DB S63/0011-2021",
		}, {
			province:      "青海",
			id:            "FA21F63F-56A7-4BA0-A502-09A830325B95",
			title:         "食品安全地方标准 黑果枸杞",
			rn:            "100",
			standard_code: "DB S63/0010-2021",
		}, {
			province:      "青海",
			id:            "05509E9D-074A-40BE-9986-579C28E84DE1",
			title:         "食品安全地方标准 牦牛奶酪",
			rn:            "101",
			standard_code: "DB S63/0008-2021",
		}, {
			province:      "青海",
			id:            "98D69266-58CF-4B8F-9780-5EE933E93079",
			title:         "食品安全地方标准 枸杞芽茶",
			rn:            "102",
			standard_code: "DB S63/0004-2021",
		}, {
			province:      "青海",
			id:            "85D4193D-9E57-4F52-B168-142E9B55F142",
			title:         "食品安全地方标准 牦牛生乳",
			rn:            "103",
			standard_code: "DB S63/0001-2019",
		}, {
			province:      "甘肃",
			id:            "77C5CCF3-4207-443B-B725-EEAFF0301F50",
			title:         "食品安全地方标准 餐饮食品外卖卫生规范",
			rn:            "104",
			standard_code: "DB S62/006-2020",
		}, {
			province:      "甘肃",
			id:            "D494263D-F0FE-4E9B-9731-3DE2BB088E8B",
			title:         "食品安全地方标准 兰州牛肉面（煮食型）",
			rn:            "105",
			standard_code: "DB S62/003-2020",
		}, {
			province:      "甘肃",
			id:            "24AB56C3-4B4F-4812-BF90-9DBBCEAAC38F",
			title:         "食品安全地方标准 金银花",
			rn:            "106",
			standard_code: "DB S62/005-2020",
		}, {
			province:      "甘肃",
			id:            "465E6BDA-DDA8-4AB2-9F5C-F6AB664E8160",
			title:         "食品安全地方标准 当归生产卫生规范",
			rn:            "107",
			standard_code: "DB S62/004-2020",
		}, {
			province:      "云南",
			id:            "352A131B-9406-4F2D-BFB4-64D1F106453F",
			title:         "食品安全地方标准 三七须根",
			rn:            "108",
			standard_code: "DB S53/029-2020",
		}, {
			province:      "贵州",
			id:            "0555076D-B17D-468F-9D99-6D11E2573679",
			title:         "食品安全地方标准 米粉（米皮）",
			rn:            "109",
			standard_code: "DB S52/051-2021",
		}, {
			province:      "贵州",
			id:            "0C4233FD-3E76-4587-A241-6015E6CEA9C1",
			title:         "食品安全地方标准 金钗石斛叶（干制品）",
			rn:            "110",
			standard_code: "DB S52/050-2021",
		}, {
			province:      "贵州",
			id:            "F67F5FD7-08BC-484A-A04F-3CB23B681F17",
			title:         "食品安全地方标准 金钗石斛花（干制品）",
			rn:            "111",
			standard_code: "DB S52/049-2021",
		}, {
			province:      "贵州",
			id:            "CD78E381-D3F0-4C73-AE34-4F442C0F8BDD",
			title:         "食品安全地方标准 食品摊贩卫生规范",
			rn:            "112",
			standard_code: "DB S52/044—2020",
		}, {
			province:      "贵州",
			id:            "B838D958-AAA6-4579-9AF2-2FDE8F20B9FC",
			title:         "食品安全地方标准 食品生产加工小作坊卫生规范",
			rn:            "113",
			standard_code: "DB S52/043—2020",
		}, {
			province:      "贵州",
			id:            "4B1372A3-54F9-411C-ADF4-B8CFB5795633",
			title:         "食品安全地方标准 米豆腐、豌豆凉粉加工小作坊卫生规范",
			rn:            "114",
			standard_code: "DB S52/040—2019",
		}, {
			province:      "贵州",
			id:            "9DB61268-F76B-4507-923C-FA5AA97B1066",
			title:         "食品安全地方标准 豆制品加工小作坊卫生规范",
			rn:            "115",
			standard_code: "DB S52/039—2019",
		}, {
			province:      "四川",
			id:            "A21FEEAF-0635-4DE9-A29B-3F3F4273FCB1",
			title:         "食品安全地方标准 自热式方便火锅生产卫生规范",
			rn:            "116",
			standard_code: "DB S51/009-2020",
		}, {
			province:      "广西",
			id:            "D1E2DFF7-B2A8-49AE-BDA7-67DA99A70FC7",
			title:         "食品安全地方标准 食品工业用冷冻水果浆（汁）",
			rn:            "117",
			standard_code: "DB S45/059-2019",
		}, {
			province:      "广西",
			id:            "4FCD5312-95DD-48C4-A0EE-55D17518646C",
			title:         "食品安全地方标准 红糟酸",
			rn:            "118",
			standard_code: "DB S45/061-2019",
		}, {
			province:      "广西",
			id:            "9673CD36-C2CA-48D7-9C21-5D2A7CD63C00",
			title:         "食品安全地方标准 粉利",
			rn:            "119",
			standard_code: "DB S45/063-2019",
		}, {
			province:      "广西",
			id:            "E69964AB-285D-4672-9C93-E7A1B64AF647",
			title:         "食品安全地方标准 螺蛳鸭脚煲",
			rn:            "120",
			standard_code: "DB S45/066-2020",
		}, {
			province:      "广东",
			id:            "359379B9-D1E1-4E25-873D-0E855524FBD1",
			title:         "食品安全地方标准 广东省食品安全地方标准 橄榄菜",
			rn:            "121",
			standard_code: "DB S44/014-2019",
		}, {
			province:      "湖北",
			id:            "A1BEA55D-A843-4DA3-BAEB-43B9A03172B4",
			title:         "食品安全地方标准 孝感米酒生产技术规范",
			rn:            "122",
			standard_code: "DB S42/012-2020",
		}, {
			province:      "湖北",
			id:            "40438F5A-7E7F-4E56-9F26-A75D0447D857",
			title:         "食品安全地方标准 熟卤制品气调包装要求",
			rn:            "123",
			standard_code: "DB S42/008-2021",
		}, {
			province:      "湖北",
			id:            "362B8314-F4D8-45B9-A93A-0C3A663C784D",
			title:         "食品安全地方标准 魔芋膳食纤维",
			rn:            "124",
			standard_code: "DB S42/007-2021",
		}, {
			province:      "湖北",
			id:            "FB378B82-11F9-4CD5-99EF-BDDDCBE2D350",
			title:         "食品安全地方标准 现制饮料加工操作卫生规范",
			rn:            "125",
			standard_code: "DB S42/015-2021",
		}, {
			province:      "湖北",
			id:            "870780A0-4D55-4E4A-AE7F-35FF0B652000",
			title:         "食品安全地方标准 脐橙蒸馏酒生产技术规范",
			rn:            "126",
			standard_code: "DB S42/013-2021",
		}, {
			province:      "河南",
			id:            "FA2A4C30-43A2-4795-8C4D-E9F48482CC1D",
			title:         "食品安全地方标准 油茶",
			rn:            "127",
			standard_code: "DB S41/005-2020",
		}, {
			province:      "河南",
			id:            "A5BACDF9-75CE-4BF8-85D7-1A6E120BCCA0",
			title:         "食品安全地方标准 食品小作坊通用卫生规范",
			rn:            "128",
			standard_code: "DB S41/012-2020",
		}, {
			province:      "浙江",
			id:            "39EA8012-FD54-481E-8782-478A671E22E7",
			title:         "食品安全地方标准 粽子生产卫生规范",
			rn:            "129",
			standard_code: "DB 33/3010-2020",
		}, {
			province:      "安徽",
			id:            "2798D0AE-D73C-4BB4-8C99-B248A3165886",
			title:         "食品安全地方标准 食品小作坊卫生规范",
			rn:            "130",
			standard_code: "DB S34/003-2021",
		}, {
			province:      "安徽",
			id:            "1711D6FD-908A-4549-ACD2-7BD886AE1201",
			title:         "食品安全地方标准 霍山石斛茎（人工种植）",
			rn:            "131",
			standard_code: "DB S34/002—2019",
		}, {
			province:      "内蒙古",
			id:            "9A764E1D-203C-4E60-BEE6-13B7E01B6BAF",
			title:         "食品安全地方标准 蒙古族传统乳制品 策格（酸马奶）",
			rn:            "132",
			standard_code: "DB S15/013-2019",
		}, {
			province:      "内蒙古",
			id:            "B137A19C-4F0B-4165-A9EB-F1BF7B867C5D",
			title:         "食品安全地方标准 亚麻籽粉",
			rn:            "133",
			standard_code: "DB S15/014-2020",
		}, {
			province:      "重庆",
			id:            "DB0878C2-D47E-4562-BCFB-0D904F80D5A2",
			title:         "食品安全地方标准 食品生产加工小作坊通用卫生规范",
			rn:            "134",
			standard_code: "DB S50/029-2020",
		}, {
			province:      "吉林",
			id:            "032C1F56-739A-4EE4-B0BB-DFC5D0612AF8",
			title:         "食品安全地方标准 代用茶",
			rn:            "135",
			standard_code: "DB S22/032-2018",
		}, {
			province:      "贵州",
			id:            "268BCE45-78A7-406E-BAE1-9C0880500641",
			title:         "食品安全地方标准 贵州辣子鸡",
			rn:            "136",
			standard_code: "DB S52/001—2014",
		}, {
			province:      "天津",
			id:            "AFCACB49-C1AF-4344-8C6E-1F8C7094AF8F",
			title:         "食品安全地方标准 工业化豆芽生产卫生规范",
			rn:            "137",
			standard_code: "DB S12/001—2014",
		}, {
			province:      "海南",
			id:            "2475EDED-99C4-4731-B8C9-7B94EBF7CC41",
			title:         "食品安全地方标准 鹧鸪茶",
			rn:            "138",
			standard_code: "DB S46/001—2018",
		}, {
			province:      "贵州",
			id:            "A4021369-BDAA-46D9-9846-39CAD823D4F3",
			title:         "食品安全地方标准 贵州苕丝糖生产卫生规范",
			rn:            "139",
			standard_code: "DB S52/034—2018",
		}, {
			province:      "贵州",
			id:            "BA12A7B5-4E5C-46C0-8755-55B97267E17D",
			title:         "食品安全地方标准 贵州苕丝糖",
			rn:            "140",
			standard_code: "DB S52/033—2018",
		}, {
			province:      "贵州",
			id:            "BD25493B-E6F8-4DDA-ACF4-67D2D120F2FD",
			title:         "食品安全地方标准 鱼酱酸调味料生产卫生规范",
			rn:            "141",
			standard_code: "DB S52/032—2018",
		}, {
			province:      "贵州",
			id:            "034EE203-38BD-49DF-BE4B-400AE08BCBB3",
			title:         "食品安全地方标准 鱼酱酸调味料",
			rn:            "142",
			standard_code: "DB S52/031—2018",
		}, {
			province:      "贵州",
			id:            "7527A157-EF46-4B5D-B0C6-08CA21DDAE44",
			title:         "食品安全地方标准 贵州省现榨饮品加工卫生规范",
			rn:            "143",
			standard_code: "DB S52/028—2017",
		}, {
			province:      "贵州",
			id:            "192CCF68-0A12-409B-8105-11D168C1FB35",
			title:         "食品安全地方标准 贵州省食用冰加工卫生规范",
			rn:            "144",
			standard_code: "DB S52/027—2017",
		}, {
			province:      "贵州",
			id:            "931434FC-7DE5-4E9F-937B-FD2169A7E9AE",
			title:         "食品安全地方标准 贵州省凉拌菜加工卫生规范",
			rn:            "145",
			standard_code: "DB S52/026—2017",
		}, {
			province:      "贵州",
			id:            "EFAA7CE3-6A0E-46C5-8B8D-68F083FEB6BA",
			title:         "食品安全地方标准 贵州米粉（米皮）加工卫生规范",
			rn:            "146",
			standard_code: "DB S52/025—2017",
		}, {
			province:      "贵州",
			id:            "C0B08DF8-A07F-487B-BEB8-6EC5E926879A",
			title:         "食品安全地方标准 豆沙粑",
			rn:            "147",
			standard_code: "DB S52/018—2016",
		}, {
			province:      "贵州",
			id:            "04B0B16C-20CB-4B2B-BF0C-62CC71C3ECCD",
			title:         "食品安全地方标准 贵州小米鲊",
			rn:            "148",
			standard_code: "DB S52/017—2016",
		}, {
			province:      "贵州",
			id:            "E2380B30-B76D-4619-871C-AE862A4963F6",
			title:         "食品安全地方标准 鸡枞油",
			rn:            "149",
			standard_code: "DB S52/016—2016",
		}, {
			province:      "贵州",
			id:            "F5A97D00-ED57-4368-852E-0266F08B949B",
			title:         "食品安全地方标准 贵州素辣椒",
			rn:            "150",
			standard_code: "DB S52/015—2016",
		}, {
			province:      "贵州",
			id:            "96C34725-C72E-4610-829C-23C4C8546B78",
			title:         "食品安全地方标准 贵州辣椒面",
			rn:            "151",
			standard_code: "DB S52/011—2016",
		}, {
			province:      "贵州",
			id:            "0024C88D-B4F3-4F7C-99EC-327CF15B580D",
			title:         "食品安全地方标准 贵州鲊辣椒",
			rn:            "152",
			standard_code: "DB S52010-2016",
		}, {
			province:      "贵州",
			id:            "1BC02152-CB00-49EA-A12B-64CE0B636605",
			title:         "食品安全地方标准 代用茶",
			rn:            "153",
			standard_code: "DB S52/002—2014",
		}, {
			province:      "云南",
			id:            "94BDD56E-0999-4971-89BE-22BE3A309047",
			title:         "食品安全地方标准 食品生产加工小作坊卫生规范",
			rn:            "154",
			standard_code: "DB S53/028-2018",
		}, {
			province:      "云南",
			id:            "AEA2848D-4601-4F9C-A3B4-242EF0206283",
			title:         "食品安全地方标准 食用玫瑰花馅料",
			rn:            "155",
			standard_code: "DB S53/025-2017",
		}, {
			province:      "云南",
			id:            "82BEF838-6C27-4FC0-949E-BD8930EE7AB9",
			title:         "食品安全地方标准 鲜花饼",
			rn:            "156",
			standard_code: "DB S53/019-2014",
		}, {
			province:      "云南",
			id:            "DD5C94DF-F2BE-4F7B-8C1E-29FAE69130D2",
			title:         "食品安全地方标准 昌宁红茶",
			rn:            "157",
			standard_code: "DB S53/012-2013",
		}, {
			province:      "云南",
			id:            "E007B2B1-B7A8-454C-8EEC-4721EA64764B",
			title:         "食品安全地方标准 乳扇",
			rn:            "158",
			standard_code: "DB S53/010-2016",
		}, {
			province:      "云南",
			id:            "86F2AB93-28D2-49A4-A775-2F978E7213E4",
			title:         "食品安全地方标准 乳饼",
			rn:            "159",
			standard_code: "DB S53/009-2016",
		}, {
			province:      "云南",
			id:            "9FA59A05-F82D-43A4-9E27-3F281FB2E3D5",
			title:         "食品安全地方标准 云南小曲清香型白酒",
			rn:            "160",
			standard_code: "DB S53/007-2015",
		}, {
			province:      "广西",
			id:            "C9072EE3-840B-4AB3-A02D-AEA60946A5F6",
			title:         "食品安全地方标准 小油坊压榨花生油黄曲霉毒素B1控制规范",
			rn:            "161",
			standard_code: "DB S45/045-2017",
		}, {
			province:      "广西",
			id:            "67CDB3A3-4FB1-4327-9EEC-D925C0A8AECF",
			title:         "食品安全地方标准 金花茶叶茶",
			rn:            "162",
			standard_code: "DB S45/033-2016",
		}, {
			province:      "广西",
			id:            "F4E7EC6D-7D0F-4D93-B5F7-E30AD00BD441",
			title:         "食品安全地方标准 黑凉粉（干粉）",
			rn:            "163",
			standard_code: "DB S45/013-2014",
		}, {
			province:      "广西",
			id:            "33CF1D6F-8FD2-4971-BB2F-DAD5FE149345",
			title:         "食品安全地方标准 油茶",
			rn:            "164",
			standard_code: "DB S45/003-2018",
		}, {
			province:      "广西",
			id:            "EE34D574-0A1A-4091-AB32-5E7B9E9FDF59",
			title:         "食品安全地方标准 柠檬鸭",
			rn:            "165",
			standard_code: "DB S45/056-2018",
		}, {
			province:      "西藏",
			id:            "7C4855BD-042E-45FD-91B0-E33C0467EC6F",
			title:         "食品安全地方标准 糌粑",
			rn:            "166",
			standard_code: "DB S54/2002-2017",
		}, {
			province:      "福建",
			id:            "281B30AB-50BA-4C8C-9BAD-420A3B9E7B0D",
			title:         "食品安全地方标准 红曲黄酒",
			rn:            "167",
			standard_code: "DB S35/003-2017",
		}, {
			province:      "福建",
			id:            "4FE44CAD-5439-4C37-BADD-49DFC619170F",
			title:         "食品安全地方标准 酿造用红曲",
			rn:            "168",
			standard_code: "DB S35/002-2017",
		}, {
			province:      "福建",
			id:            "1FCB5574-2452-4FB4-A655-E4B0E63768CE",
			title:         "食品安全地方标准 连城地瓜干系列产品",
			rn:            "169",
			standard_code: "DB S35/001-2017",
		}, {
			province:      "宁夏",
			id:            "99222FDB-F11E-4215-9935-CEC3961106FF",
			title:         "食品安全地方标准 凉皮",
			rn:            "170",
			standard_code: "DB S64/004-2019",
		}, {
			province:      "宁夏",
			id:            "70152B13-3775-4CC5-B8E2-53CFEA55F49E",
			title:         "食品安全地方标准 火锅底料",
			rn:            "171",
			standard_code: "DB S64/003-2018",
		}, {
			province:      "宁夏",
			id:            "D2005D7C-D0E2-4C0B-86D2-D8F1CB79646E",
			title:         "食品安全地方标准 八宝茶",
			rn:            "172",
			standard_code: "DB S64/002-2018",
		}, {
			province:      "内蒙古",
			id:            "AED1A79C-A7A5-44AA-BD48-8C03D577740D",
			title:         "食品安全地方标准 蒙古族传统乳制品 第2部分：奶皮子",
			rn:            "173",
			standard_code: "DB S15/001.2-2016",
		}, {
			province:      "内蒙古",
			id:            "328757D9-138A-4E48-AEB7-934B811F3C18",
			title:         "食品安全地方标准 蒙古族传统乳制品 第3部分：奶豆腐",
			rn:            "174",
			standard_code: "DB S15/001.3-2017",
		}, {
			province:      "内蒙古",
			id:            "B9999B50-A806-4DA2-9502-EA2CAC70AEFC",
			title:         "食品安全地方标准 含乳固态成型制品",
			rn:            "175",
			standard_code: "DB S15/002-2013",
		}, {
			province:      "内蒙古",
			id:            "075FD5B0-CBDB-4C2F-BC78-2B5E3CDF9419",
			title:         "食品安全地方标准 蒙古族传统乳制品 毕希拉格",
			rn:            "176",
			standard_code: "DB S15/005-2017",
		}, {
			province:      "内蒙古",
			id:            "1D7E7FCF-1919-46BC-8C30-45BE258A22E8",
			title:         "食品安全地方标准 蒙古族传统乳制品 酸酪蛋（奶干）",
			rn:            "177",
			standard_code: "DB S15/006-2016",
		}, {
			province:      "内蒙古",
			id:            "25E4B95B-BCC8-474D-B9B7-2F2FBBDB5441",
			title:         "食品安全地方标准 蒙古族传统乳制品 楚拉",
			rn:            "178",
			standard_code: "DB S15/007-2016",
		}, {
			province:      "内蒙古",
			id:            "DB3F7926-66C4-49EC-B7B3-6C8A40C6D939",
			title:         "食品安全地方标准 蒙古族传统乳制品生产卫生规范",
			rn:            "179",
			standard_code: "DB S15/008-2016",
		}, {
			province:      "内蒙古",
			id:            "4C2D2396-8E09-4563-8705-F0CBA59D505F",
			title:         "食品安全地方标准 炒米",
			rn:            "180",
			standard_code: "DB S15/010-2016",
		}, {
			province:      "内蒙古",
			id:            "5C3B7ADB-E98C-4A93-8DC8-859DF198FE89",
			title:         "食品安全地方标准 生马乳",
			rn:            "181",
			standard_code: "DB S15/011-2019",
		}, {
			province:      "内蒙古",
			id:            "7E0012FB-A761-4A8E-806F-20A85F2074B6",
			title:         "食品安全地方标准 蒙古族传统乳制品 嚼克",
			rn:            "182",
			standard_code: "DB S15/012-2019",
		}, {
			province:      "内蒙古",
			id:            "8B4AEBAD-ED64-4460-B063-A92E65C635DF",
			title:         "食品安全地方标准 生驼乳",
			rn:            "183",
			standard_code: "DB S15/015-2019",
		}, {
			province:      "内蒙古",
			id:            "0B39EB81-2CC1-48D2-B6A3-D8502AA3C37B",
			title:         "食品安全地方标准 灭菌驼乳",
			rn:            "184",
			standard_code: "DB S15/017-2019",
		}, {
			province:      "内蒙古",
			id:            "2A809186-7F33-4A7D-9D80-2DC6357C25F5",
			title:         "食品安全地方标准 奶茶粉",
			rn:            "185",
			standard_code: "DB S15/001.1-2019",
		}, {
			province:      "内蒙古",
			id:            "0F40F8F0-C4B0-449F-BBD2-9DA3285F2FE0",
			title:         "食品安全地方标准 奶片",
			rn:            "186",
			standard_code: "DB S15/009-2019",
		}, {
			province:      "河北",
			id:            "6D5D7BFC-B212-41B3-AB45-6BB0C2B5EA79",
			title:         "食品安全地方标准 代用茶",
			rn:            "187",
			standard_code: "DB S13/002-2015",
		}, {
			province:      "河北",
			id:            "B4C27068-5D14-4E0D-99D6-E0CD2F0B23B6",
			title:         "食品安全地方标准 龙凤贡面生产卫生规范",
			rn:            "188",
			standard_code: "DB S13/011-2018",
		}, {
			province:      "河北",
			id:            "8C9C102B-BD22-44ED-B1E8-75588196A441",
			title:         "食品安全地方标准 速冻草莓生产卫生规范",
			rn:            "189",
			standard_code: "DB S13/012-2018",
		}, {
			province:      "河南",
			id:            "2BA316F1-791C-45C2-B492-3E47AFE0CDA7",
			title:         "食品安全地方标准 食用畜禽血制品",
			rn:            "190",
			standard_code: "DB S41/011—2016",
		}, {
			province:      "河南",
			id:            "B8488144-A951-48E0-BD94-068C08BBE414",
			title:         "食品安全地方标准 代用茶",
			rn:            "191",
			standard_code: "DB S41/010—2016",
		}, {
			province:      "山西",
			id:            "DCA22C03-D831-4DD6-AC6D-18F857153E0E",
			title:         "食品安全地方标准 连翘叶",
			rn:            "192",
			standard_code: "DB S14/001-2017",
		}, {
			province:      "江西",
			id:            "5209A840-50EF-46AA-8D96-D4FC83F0E582",
			title:         "食品安全地方标准 南酸枣糕生产卫生规范",
			rn:            "193",
			standard_code: "DB 36/1090-2018",
		}, {
			province:      "江西",
			id:            "043A75A1-4B42-4680-A97D-442F0968A5CB",
			title:         "食品安全地方标准 鲜湿类米粉生产卫生规范",
			rn:            "194",
			standard_code: "DB 36/1091-2018",
		}, {
			province:      "四川",
			id:            "FF96D567-4433-40F0-A011-9A8F162A79BE",
			title:         "食品安全地方标准 苦荞茶",
			rn:            "195",
			standard_code: "DB S51/004-2017",
		}, {
			province:      "安徽",
			id:            "E20BBBCF-A02D-436C-8B09-2E7136F8BD45",
			title:         "食品安全地方标准 代用茶",
			rn:            "196",
			standard_code: "DB S34/2607-2016",
		}, {
			province:      "重庆",
			id:            "7B9020B0-8AC5-4072-8AAE-6521EB1297E1",
			title:         "食品安全地方标准 保鲜花椒",
			rn:            "197",
			standard_code: "DB S50/003-2014",
		}, {
			province:      "重庆",
			id:            "30FA8C26-CC1C-48D3-B319-0F419BFD4B67",
			title:         "食品安全地方标准 食用畜血产品（血旺）",
			rn:            "198",
			standard_code: "DB S50/017-2014",
		}, {
			province:      "浙江",
			id:            "87B7A109-0DE5-47D8-8752-AA0F1514DCD7",
			title:         "食品安全地方标准 现榨果蔬汁、五谷杂粮饮品",
			rn:            "199",
			standard_code: "DB 33/3005-2015",
		}, {
			province:      "浙江",
			id:            "2B7100D0-2BE0-4FF5-97A8-2B023EE3A6DF",
			title:         "食品安全地方标准 火腿生产卫生规范",
			rn:            "200",
			standard_code: "DB 33/3008-2016",
		}, {
			province:      "浙江",
			id:            "5FA997D8-E219-4AF0-A11E-4618D6BE43DA",
			title:         "食品安全地方标准 食品小作坊通用卫生规范",
			rn:            "201",
			standard_code: "DB 33/3009-2018",
		}, {
			province:      "江苏",
			id:            "C00CBFA0-3E1C-47BB-AF16-80EEE7D2F8B7",
			title:         "食品安全地方标准 工业化豆芽生产卫生规范",
			rn:            "202",
			standard_code: "DB S32/018-2018",
		}, {
			province:      "江苏",
			id:            "2EC731F0-0FDC-4DEC-A03D-F3F4CC429776",
			title:         "食品安全地方标准 甜油（调味品）",
			rn:            "203",
			standard_code: "DB S32/016-2018",
		}, {
			province:      "江苏",
			id:            "F2CC4AC7-6330-4153-9D5F-3C42C6BC43D5",
			title:         "食品安全地方标准 食品小作坊卫生规范",
			rn:            "204",
			standard_code: "DB S32/013-2017",
		}, {
			province:      "江苏",
			id:            "D9D13249-E065-471F-B715-0F0E658A1B4B",
			title:         "食品安全地方标准 耳叶牛皮消制品",
			rn:            "205",
			standard_code: "DB S32/009-2016",
		}, {
			province:      "江苏",
			id:            "16193E5A-2BA4-499B-BF02-62B008C67280",
			title:         "食品安全地方标准 熟制鸡胚蛋（活珠子）",
			rn:            "206",
			standard_code: "DB S32/007-2016",
		}, {
			province:      "广东",
			id:            "5D6E7AA4-F312-49FE-86DD-6276926F277D",
			title:         "食品安全地方标准 纳豆粉",
			rn:            "207",
			standard_code: "DB S44/013-2019",
		}, {
			province:      "广东",
			id:            "88BE2AA8-18D2-4818-94C3-45796F102FAE",
			title:         "食品安全地方标准 白木香叶",
			rn:            "208",
			standard_code: "DB S44/011-2018",
		}, {
			province:      "广东",
			id:            "0C2AF22B-FD73-4348-829A-E788AC6D3364",
			title:         "食品安全地方标准 新会柑皮含茶制品",
			rn:            "209",
			standard_code: "DB S44/010-2018",
		}, {
			province:      "广东",
			id:            "34EFAE8C-1184-414F-947A-21FA9E3AA481",
			title:         "食品安全地方标准 簕菜及干制品",
			rn:            "210",
			standard_code: "DB S44/009-2018",
		}, {
			province:      "广东",
			id:            "BE31A545-82C1-44E7-9D81-0BEAA2A97A08",
			title:         "食品安全地方标准 预包装冷藏、冷冻膳食生产经营卫生规范",
			rn:            "211",
			standard_code: "DB S44/008-2017",
		}, {
			province:      "广东",
			id:            "E85BF5A7-CECA-4042-9712-2D7354DF690E",
			title:         "食品安全地方标准 预包装冷藏、冷冻膳食",
			rn:            "212",
			standard_code: "DB S44/007-2017",
		}, {
			province:      "湖南",
			id:            "A2C973F6-76EC-44B8-8753-C172CAA86E55",
			title:         "食品安全地方标准 气调包装酱卤肉制品生产卫生规范",
			rn:            "213",
			standard_code: "DB S43/009-2018",
		}, {
			province:      "湖南",
			id:            "93A48D8B-9082-41D7-8E4E-D40E41EA767C",
			title:         "食品安全地方标准 湿面生产卫生规范",
			rn:            "214",
			standard_code: "DB S43/008-2018",
		}, {
			province:      "湖南",
			id:            "B9FBBA7E-E5F8-406E-952A-1E4861C44846",
			title:         "食品安全地方标准 米粉生产卫生规范",
			rn:            "215",
			standard_code: "DB S43/007-2018",
		}, {
			province:      "青海",
			id:            "091E5261-1B99-4351-94CA-0783FBAD2942",
			title:         "食品安全地方标准 茶卡盐",
			rn:            "216",
			standard_code: "DB S/630001-2019",
		}, {
			province:      "上海",
			id:            "7C5501EF-0388-477B-AA22-D1D8A3B9C5A8",
			title:         "食品安全地方标准 即食食品自动售卖（制售）卫生规范",
			rn:            "217",
			standard_code: "DB 31/2028-2019",
		},
	},
}

var DbSpPtCookie = "cookieName=cookieValue; name=value; cookieName=cookieValue; name=value; JSESSIONID=29AF28AB20A00B1FC9511C71EFC36B75"

// ychEduSpider 获取食品安全地方标准数据文档文档
// @Title 获取食品安全地方标准数据文档文档
// @Description https://sppt.cfsa.net.cn:8087/，获取食品安全地方标准数据文档文档
func main() {
	requestListUrl := "https://sppt.cfsa.net.cn:8087/db"
	for in_index, allData := range dbSpPtAllData.AllData {
		fmt.Println("===========in_index=" + strconv.Itoa(in_index) + "============")
		// 中文标题
		title := allData.title
		title = strings.TrimSpace(title)
		title = strings.ReplaceAll(title, "/", "-")
		title = strings.ReplaceAll(title, "／", "-")
		title = strings.ReplaceAll(title, "　", "")
		title = strings.ReplaceAll(title, " ", "")
		title = strings.ReplaceAll(title, "：", ":")
		title = strings.ReplaceAll(title, "—", "-")
		title = strings.ReplaceAll(title, "－", "-")
		title = strings.ReplaceAll(title, "（", "(")
		title = strings.ReplaceAll(title, "）", ")")
		title = strings.ReplaceAll(title, "《", "")
		title = strings.ReplaceAll(title, "》", "")
		title = strings.ReplaceAll(title, "()", "")
		fmt.Println(title)

		code := allData.standard_code
		code = strings.TrimSpace(code)
		code = strings.ReplaceAll(code, "/", "-")
		code = strings.ReplaceAll(code, "—", "-")
		fmt.Println(code)

		filePath := "../sppt.cfsa.net.cn/" + title + "(" + code + ")" + ".pdf"
		fmt.Println(filePath)
		_, err := os.Stat(filePath)
		if err == nil {
			fmt.Println("文档已下载过，跳过")
			continue
		}

		fmt.Println("=======开始下载========")

		detailUrl := fmt.Sprintf("https://sppt.cfsa.net.cn:8087/staticPages/%s.html", allData.id)
		fmt.Println(detailUrl)

		DbSpPtDetailDoc, err := DbSpPtHtmlDoc(detailUrl, requestListUrl)
		if err != nil {
			fmt.Println("获取文档详情失败，跳过")
			continue
		}

		titleTipNode := htmlquery.FindOne(DbSpPtDetailDoc, `//html/body/div[2]/div[2]/div[1]/span/i`)
		if titleTipNode == nil {
			fmt.Println("标题节点不存在，跳过")
			continue
		}
		titleTip := htmlquery.InnerText(titleTipNode)

		detailDownloadNode := htmlquery.FindOne(DbSpPtDetailDoc, `//html/body/div[2]/div[2]/div[2]/div[2]/span/b[2]/a`)
		// load('11FA1C00-8498-4AD5-800E-638099CE36CF');
		detailClickText := htmlquery.SelectAttr(detailDownloadNode, "onclick")
		fileGuid := strings.ReplaceAll(detailClickText, "load('", "")
		fileGuid = strings.ReplaceAll(fileGuid, "');", "")
		fmt.Println(fileGuid)

		downloadDbSpPtUrl := "https://sppt.cfsa.net.cn:8087/cfsa_aiguo"
		fmt.Println(downloadDbSpPtUrl)
		downloadDbSpPtFormData := DownloadDbSpPtFormData{
			task:      "d_p",
			guid:      "",
			file_guid: fileGuid,
			num_tn:    "",
			type_temp: "",
			fact_name: fileGuid,
			filePath:  "/opt/tomcat/Res_upload/swfupload",
			title_tip: titleTip,
		}
		err = downloadDbSpPt(downloadDbSpPtUrl, detailUrl, downloadDbSpPtFormData, filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//复制文件
		tempFilePath := strings.ReplaceAll(filePath, "../sppt.cfsa.net.cn", "../upload.doc88.com/hbba.sacinfo.org.cn")
		err = copyDbSpPtFile(filePath, tempFilePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("=======完成下载========")
		//DownLoadDbSpPtTimeSleep := 10
		DownLoadDbSpPtTimeSleep := rand.Intn(5)
		for i := 1; i <= DownLoadDbSpPtTimeSleep; i++ {
			time.Sleep(time.Second)
			fmt.Println("title="+title+"===========下载", title, "成功，暂停", DownLoadDbSpPtTimeSleep, "秒，倒计时", i, "秒===========")
		}
	}
}

func DbSpPtHtmlDoc(requestUrl string, referer string) (doc *html.Node, err error) {
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
	if DbSpPtEnableHttpProxy {
		client = DbSpPtSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接
	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", DbSpPtCookie)
	req.Header.Set("Host", "sppt.cfsa.net.cn:8087")
	req.Header.Set("Origin", "https://sppt.cfsa.net.cn:8087/")
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
		return doc, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
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

func downloadDbSpPt(requestUrl string, referer string, downloadDbSpPtFormData DownloadDbSpPtFormData, filePath string) error {
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
	if DbSpPtEnableHttpProxy {
		client = DbSpPtSetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("task", downloadDbSpPtFormData.task)
	postData.Add("guid", downloadDbSpPtFormData.guid)
	postData.Add("file_guid", downloadDbSpPtFormData.file_guid)
	postData.Add("num_tn", downloadDbSpPtFormData.num_tn)
	postData.Add("type", downloadDbSpPtFormData.type_temp)
	postData.Add("fact_name", downloadDbSpPtFormData.fact_name)
	postData.Add("filePath", downloadDbSpPtFormData.filePath)
	postData.Add("title_tip", downloadDbSpPtFormData.title_tip)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", DbSpPtCookie)
	req.Header.Set("Host", "sppt.cfsa.net.cn:8087")
	req.Header.Set("Origin", "https://sppt.cfsa.net.cn:8087")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
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

func copyDbSpPtFile(src, dst string) (err error) {
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
