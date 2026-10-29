package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"rsc.io/pdf"
)

// Doc88Cookie 15238369929
const Doc88Cookie = "cdb_sys_sid=6nh9k3a70toarhknllurprqi75; cdb_RW_ID_2115708421=1; cdb_RW_ID_2115708341=1; cdb_RW_ID_2115709163=1; cdb_logined=1; cdb_RW_ID_2115713360=1; cdb_RW_ID_2115426657=1; cdb_RW_ID_2115802811=1; cdb_RW_ID_2115802808=1; cdb_RW_ID_2115706548=1; cdb_RW_ID_2115831923=1; cdb_RW_ID_2115677074=1; cdb_RW_ID_2115831614=1; cdb_RW_ID_2115887118=1; cdb_RW_ID_2115887105=1; cdb_RW_ID_2115887085=1; cdb_RW_ID_2115887070=1; cdb_RW_ID_2115888571=1; cdb_RW_ID_2115888567=1; cdb_RW_ID_2115888569=1; cdb_RW_ID_2115889760=1; cdb_RW_ID_2115889752=1; cdb_RW_ID_2107316235=1; cdb_RW_ID_2115889714=1; cdb_RW_ID_1976343520=4; cdb_RW_ID_2115628180=1; cdb_RW_ID_2115889696=1; cdb_login_if=1; cdb_uid=104598337; cdb_RW_ID_2115895696=1; cdb_RW_ID_862217681=2166; cdb_RW_ID_862880263=1250; cdb_RW_ID_862879695=95; cdb_RW_ID_862879371=17; cdb_RW_ID_2116223176=1; cdb_RW_ID_2116223677=1; cdb_RW_ID_21282751=350; cdb_RW_ID_2116443221=1; cdb_RW_ID_2116456789=1; cdb_RW_ID_2116543395=1; cdb_RW_ID_2116543389=1; cdb_RW_ID_2116543368=1; cdb_RW_ID_2116543405=1; cdb_RW_ID_2115889075=1; cdb_RW_ID_2116783626=1; cdb_RW_ID_1060444713=4; cdb_RW_ID_1060444681=2; cdb_RW_ID_2116943668=1; cdb_RW_ID_2117086316=1; cdb_RW_ID_2117176776=1; cdb_RW_ID_2117464100=1; cdb_RW_ID_2117540517=1; cdb_RW_ID_2117539921=2; cdb_RW_ID_2117768306=1; cdb_RW_ID_2079034316=537; PHPSESSID=6nh9k3a70toarhknllurprqi75; show_index=1; cdb_RW_ID_1349878984=261; cdb_H5R=1; cdb_RW_ID_1772212236=18; cdb_RW_ID_1762201897=5; cdb_RW_ID_24719283=94; cdb_RW_ID_49519562=9; cdb_RW_ID_2120242853=1; cdb_RW_ID_2020420508=11; cdb_RW_ID_28731639=60; Page_Y_090371409205=-191.41776315789474; Page_090371409205=6; cdb_RW_ID_2121999365=1; cdb_READED_PC_ID=%2C441441441443443441441720448593440443447702730593446730730730447440593443448440441440443; cdb_RW_ID_2117646269=1; Page_29654993878486=1; Page_Y_67519696333048=16.049342105263158; Page_67519696333048=1; cdb_RW_ID_2112326126=1; Page_Y_29237997072972=-119.39144736842105; Page_29237997072972=1; cdb_RW_ID_2098527909=3; Page_99554062143606=1; cdb_RW_ID_2122200363=1; Page_Y_90054944400585=1926.716447368421; Page_90054944400585=5; cdb_RW_ID_2020027997=7; Page_Y_18068484486336=206.29276315789474; Page_18068484486336=1; cdb_RW_ID_2122200364=1; Page_Y_74480700044158=139.7467105263158; Page_74480700044158=6; cdb_RW_ID_2122650415=1; ShowSkinTip_1=1; showAnnotateTipIf=1; Page_Y_91390900231493=-119.39144736842105; Page_91390900231493=1; cdb_pageType=2; cdb_RW_ID_2123084237=1; Page_Y_74480701498012=275.1875; Page_74480701498012=1; cdb_RW_ID_2123122389=1; Page_Y_78843731733159=-81.8125; Page_78843731733159=10; cdb_RW_ID_2123122363=1; Page_Y_98837970977020=-279.88486842105266; Page_98837970977020=3; cdb_RW_ID_2123122338=1; Page_Y_67719690699005=400.4506578947369; Page_67719690699005=3; cdb_RW_ID_2123143132=1; Page_Y_95920902972920=-119.39144736842105; Page_95920902972920=1; siftState=1; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d23425c000f13f2a30006737660ad8ab8c6d92f350c3f26f274; doc88_lt=wx; cdb_change_message=1; cdb_msg_num=0; cdb_msg_time=1793276416"

// Doc88Cookie 17539389687
// const Doc88Cookie = "cdb_sys_sid=vc0gvec15bo7h8rjq84llgd1h4; PHPSESSID=vc0gvec15bo7h8rjq84llgd1h4; cdb_change_message=1; cdb_msg_num=0; cdb_msg_time=1705637593; show_index=1; siftState=1; cdb_token=5176691bb4a2b7d6f48bef23f1c8c6bd6976e1452443c283d33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86c3480eba7246bba18f516cec2b411ba4e4de4720219466f830686a5776f5804dd60f82c713aca24f; cdb_tokenid=f4ec%2F1tvyaOMSK9UxFUgZJuEPiN3OwbisWAvNCi61AWsiQ9S50A2U2B66txBgWSEJVeGmWi8SkcMC8XLLbajoTPgd8Uy4U42fhoxFX2Ee%2FE; c_login_name=1056391268; cdb_login_if=1; cdb_system_sign=1; cdb_back[act]=get_intro; cdb_back[p]=2021%2F10%2F19%2F74687192717605.xml.gz; cdb_back[order]=0; cdb_back[pcid]=593; cdb_back[pid]=1905989379; cdb_back[curpage]=1; cdb_back[folder_id]=0; cdb_back[folder_page]=0; cdb_back[member_id]=6d4d209o34hyMZNLDyvpApTMemXbDwsEwc1QjGjbHoA%2FIZY1sHFLGSCJT7afUd0Ifg4; cdb_back[order_num]=1; cdb_back[show]=1; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[classify_id]=all; cdb_back[menuIndex]=2; cdb_back[show_index]=1; cdb_back[state]=all; cdb_back[t]=0; cdb_back[u]=1; cdb_back[log_state]=0; cdb_back[log_type_page]=all; cdb_back[captchaCode]=1; cdb_back[checkagreement]=on; cdb_back[txtPassword2]=abcdqq123456; cdb_back[txtPassword]=abcdqq123456; cdb_back[txtemail]=963094343%40qq.com; cdb_back[txtloginname]=woyoceo2023; cdb_back[username]=woyoceo2023; cdb_back[book_id]=3367789; cdb_back[classify]=8140; cdb_back[code]=66202dc48d76ef18; cdb_back[at]=0; cdb_back[doctype]=1; cdb_back[n]=6; cdb_back[p_id]=1401760345; cdb_back[p_name]=%E5%A4%A7%E5%AD%A6%E7%AA%81%E5%8F%91%E5%85%AC%E5%85%B1%E4%BA%8B%E4%BB%B6%E5%BA%94%E6%80%A5%E9%A2%84%E6%A1%88; cdb_back[rel_p_id]=1401760345; cdb_back[srlid]=128129RcKWcCMH%2F7zX0d+LXHyarhH5ko9rgGbBtxe8Y3NPDEfpvwVaPApe8esAntuluwmBC6fIe5K8RzK+HWKTRHm56Y4NVPn40Iskf0L84; cdb_READED_PC_ID=%2C; cdb_RW_ID_1401760345=3989; cdb_back[image_type]=9; cdb_back[login]=1; cdb_back[module_type]=7; cdb_back[email]=963094343%40qq.com; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dall"

// Doc88Cookie 15803889687
// const Doc88Cookie = "ew_user_task_degree=100; cdb_sys_sid=oojq5k1eldqc8n34n5l26n6mn0; PHPSESSID=oojq5k1eldqc8n34n5l26n6mn0; show_index=1; siftState=1; cdb_change_message=1; cdb_msg_num=0; cdb_msg_time=1793276552; cdb_token=5176691bb4a2b7d65cd02416684e41acd8ba6dc2254cdb98eeda5c6169a4ebe4cd6d0a7e57ad73f20ed2a9a8be13a295adb2f27dcfc55127add2c2984ea49c396d186f19a3268e3e49b7df87206f6cd1481948776dec36dcc954571ed2193126bbffde38826910734f58372cd96dc12da385c2d066ae3a7d; cdb_login_if=1; cdb_uid=113085378; doc88_lt=wx; Page_91390908425111=2; Page_Y_91390908425111=233.3026315789474; cdb_READED_PC_ID=%2C441449441440440720720720443720443593447593730593593442449447593; cdb_RW_ID_2123467000=1; Page_98637970664086=1; Page_Y_98637970664086=425.1118421052632; showAnnotateTipIf=1; ShowSkinTip_1=1; cdb_RW_ID_2123447304=1; cdb_RW_ID_2123240047=1; cdb_pageType=2; cdb_H5R=1; cdb_RW_ID_2121380377=4; cdb_RW_ID_2044001790=51; cdb_RW_ID_2119916979=1; cdb_RW_ID_2117644696=24; cdb_RW_ID_2116870530=2; cdb_RW_ID_2052251107=25; cdb_RW_ID_2117073097=1; cdb_RW_ID_2117073039=1; cdb_RW_ID_2116778940=1; cdb_RW_ID_2116778943=1; cdb_RW_ID_2116778967=1; cdb_RW_ID_2116440712=6; cdb_RW_ID_40298071=47; cdb_RW_ID_2116440720=1; cdb_RW_ID_2116440648=1; cdb_RW_ID_2116440686=1; cdb_RW_ID_2116439623=1; cdb_RW_ID_2116439739=1; cdb_RW_ID_2116440208=1; cdb_RW_ID_2116435911=1; cdb_RW_ID_2116435940=1; cdb_RW_ID_2116435943=1; cdb_RW_ID_2116435830=1; cdb_RW_ID_2116435842=2; cdb_RW_ID_2116435784=1; cdb_RW_ID_2116435781=1; cdb_RW_ID_2116435540=1; cdb_RW_ID_2116435545=1; cdb_RW_ID_2116435425=1; cdb_RW_ID_2116435437=1; cdb_RW_ID_2116435443=1; cdb_RW_ID_2116435457=1; cdb_RW_ID_285343137=27; cdb_RW_ID_2116435460=1; cdb_RW_ID_2116435475=1; cdb_RW_ID_2116435234=1; cdb_RW_ID_2116435369=1; cdb_RW_ID_2116435401=2; cdb_RW_ID_2116434995=1; cdb_RW_ID_2116435014=1; cdb_RW_ID_1656241560=15; cdb_RW_ID_2116435022=1; cdb_RW_ID_2116435026=1; cdb_RW_ID_2116435031=1; cdb_RW_ID_2116435038=1; cdb_RW_ID_2116435048=1; cdb_RW_ID_1306046245=255; cdb_RW_ID_2116434799=1; cdb_RW_ID_2116434809=1; cdb_RW_ID_2116434813=1; cdb_RW_ID_1094221962=10; cdb_RW_ID_2116434467=1; cdb_RW_ID_49704181=105; cdb_RW_ID_2116434484=1; cdb_RW_ID_2116434492=1; cdb_RW_ID_2116434498=1; cdb_RW_ID_2116434460=1; cdb_RW_ID_2116434476=1; cdb_RW_ID_2116434416=1; cdb_RW_ID_2116434421=1; cdb_RW_ID_2116434398=1; cdb_RW_ID_2116434382=1; cdb_RW_ID_427793=32; cdb_RW_ID_2116434365=1; cdb_RW_ID_2116223733=3; cdb_RW_ID_2116333201=1; cdb_RW_ID_2116333204=1; cdb_RW_ID_2116230471=1; cdb_RW_ID_2116181876=1; cdb_RW_ID_2116178726=1; cdb_RW_ID_2116028498=1; cdb_RW_ID_2116028480=1; cdb_RW_ID_2115838752=1; cdb_RW_ID_2115839036=1; cdb_RW_ID_2115838849=1; cdb_RW_ID_2115700919=1; cdb_RW_ID_2115701101=2; cdb_RW_ID_2115767598=1; cdb_RW_ID_1433501490=162; cdb_RW_ID_2115767600=1; cdb_RW_ID_2115767573=1; cdb_RW_ID_2115767575=1; cdb_RW_ID_1439547305=12; cdb_RW_ID_2115709156=1; cdb_RW_ID_1407212687=4; cdb_RW_ID_1939356669=4; cdb_RW_ID_1649845332=4; cdb_RW_ID_1923193898=7; cdb_RW_ID_2115700049=1; cdb_RW_ID_1981499336=13; cdb_RW_ID_2102445734=6; cdb_RW_ID_2089853740=1302; cdb_RW_ID_2115221698=10; cdb_RW_ID_2115221660=1; cdb_RW_ID_2115220406=2; cdb_RW_ID_1677417951=3; cdb_RW_ID_1981499314=5; cdb_RW_ID_2115220416=1; cdb_RW_ID_1967355765=3; cdb_RW_ID_2029517212=2; cdb_RW_ID_2115220407=1; cdb_RW_ID_2106140402=3; cdb_RW_ID_1054243817=7; cdb_RW_ID_2115168641=1; cdb_RW_ID_2114936516=3; cdb_RW_ID_2114936492=2; cdb_RW_ID_2114933745=4; cdb_RW_ID_2063277594=3; cdb_RW_ID_2104092537=1; cdb_RW_ID_2110573716=1; cdb_RW_ID_2114772208=1; cdb_RW_ID_2114614952=1; cdb_RW_ID_2114710248=2; cdb_RW_ID_2083648166=2; cdb_RW_ID_2083648173=1; cdb_RW_ID_2106961461=2; cdb_RW_ID_2051281103=8; cdb_RW_ID_2114235600=2; cdb_RW_ID_2114235555=1; cdb_RW_ID_2114235780=1; cdb_RW_ID_2050759977=7; cdb_RW_ID_2114235882=1; cdb_RW_ID_2094187836=285; cdb_RW_ID_2107836618=29; cdb_RW_ID_2098445643=108; cdb_RW_ID_2113537378=8; cdb_RW_ID_2068550049=7; cdb_RW_ID_2073100868=2; cdb_RW_ID_2073331961=25; cdb_RW_ID_2114010806=1; cdb_RW_ID_70518641=13; cdb_RW_ID_2113261389=4; cdb_RW_ID_2113579613=1; cdb_RW_ID_2113537382=1; cdb_RW_ID_2113537374=1; cdb_RW_ID_2113210536=6; cdb_RW_ID_2100617271=16; cdb_RW_ID_2112790468=2; cdb_RW_ID_2112452348=9; cdb_RW_ID_2111827243=6; cdb_RW_ID_2106742195=10; cdb_RW_ID_2112865813=1; cdb_RW_ID_84124975=57; cdb_RW_ID_2112452120=2; cdb_RW_ID_2112078824=3; cdb_RW_ID_2107932765=9; cdb_RW_ID_2112138084=1; cdb_RW_ID_1027981347=2; cdb_RW_ID_2112125862=1; cdb_RW_ID_2111314777=9; cdb_RW_ID_2111625054=1; cdb_RW_ID_2079937436=30; cdb_RW_ID_2110834308=5; cdb_RW_ID_2110481221=3; cdb_RW_ID_2110479084=1; cdb_RW_ID_2110479894=1; cdb_RW_ID_2100617852=2; cdb_RW_ID_2110059625=5; cdb_RW_ID_2109985387=2; cdb_RW_ID_1637407729=48; cdb_RW_ID_2109986442=1; cdb_RW_ID_2098989908=2; cdb_RW_ID_2103303050=4; cdb_RW_ID_2109783624=1; cdb_RW_ID_2109771495=1; cdb_RW_ID_2109256755=7; cdb_RW_ID_2109527955=2; cdb_RW_ID_2098445811=18; cdb_RW_ID_2104084969=7; cdb_RW_ID_2104031765=2; cdb_RW_ID_2104059616=5; cdb_RW_ID_2109360573=1; cdb_RW_ID_2109267220=4; cdb_RW_ID_227607347=233; cdb_RW_ID_2109188860=1; cdb_RW_ID_2109188870=2; cdb_RW_ID_2109188864=1; cdb_RW_ID_1034138937=48; cdb_RW_ID_2109079885=1; cdb_RW_ID_2109181458=1; cdb_RW_ID_2107932759=2; cdb_RW_ID_2108806828=7; cdb_RW_ID_2080738327=2; cdb_RW_ID_2108864069=1; cdb_RW_ID_2108753935=1; cdb_RW_ID_1654332059=5; cdb_RW_ID_834390721=9; cdb_RW_ID_2108686196=1; cdb_RW_ID_2099911462=15; cdb_RW_ID_2108320835=39; cdb_RW_ID_2108537993=1; cdb_RW_ID_2108537994=1; cdb_RW_ID_2090267006=2; cdb_RW_ID_2108538017=1; cdb_RW_ID_2044923914=65; cdb_RW_ID_2048014029=39; cdb_RW_ID_2108538019=1; cdb_RW_ID_2093694760=32; cdb_RW_ID_2108538020=1; cdb_RW_ID_2088107741=29; cdb_RW_ID_2071486117=1; cdb_RW_ID_2108422591=1; cdb_RW_ID_2107789572=1; cdb_RW_ID_2107524597=1; cdb_RW_ID_2108365488=1; cdb_RW_ID_2107343680=2; cdb_RW_ID_2108365486=1; cdb_RW_ID_2087566690=356; cdb_RW_ID_2071833247=1; cdb_RW_ID_1644494866=35; cdb_RW_ID_2108008086=2; cdb_RW_ID_2108008105=3; cdb_RW_ID_2049368485=189; cdb_RW_ID_2108008102=1; cdb_RW_ID_2099172703=2; cdb_RW_ID_2108008103=1; cdb_RW_ID_2106487686=3; cdb_RW_ID_2104880311=35; cdb_RW_ID_2104843668=3; cdb_RW_ID_2104836181=3; cdb_RW_ID_2104888188=2; cdb_RW_ID_2107255344=1; cdb_RW_ID_1153096564=4271; cdb_RW_ID_2016553459=3350; cdb_RW_ID_2079490803=134; cdb_logined=1; cdb_RW_ID_2105062836=12; cdb_RW_ID_2096068843=32; cdb_RW_ID_2106574141=3; cdb_RW_ID_2106071372=2; cdb_RW_ID_2106189452=2; cdb_RW_ID_2025958563=8; cdb_RW_ID_2106071367=1; cdb_RW_ID_2106071370=1; cdb_RW_ID_2106026786=1; cdb_RW_ID_2105025852=2; cdb_RW_ID_2105546204=1; cdb_RW_ID_2104671386=4; cdb_RW_ID_2105063821=1; cdb_RW_ID_49623644=57; cdb_RW_ID_2104810095=1; cdb_RW_ID_2104253233=1; cdb_RW_ID_2103951685=1; cdb_RW_ID_2098445911=269; cdb_RW_ID_2103924059=1; cdb_RW_ID_2103933475=1; cdb_RW_ID_2103933519=2; cdb_RW_ID_2103933547=1; cdb_RW_ID_2103933556=1; cdb_RW_ID_2071080681=69; cdb_RW_ID_1641551625=554; cdb_RW_ID_2103429778=1; cdb_RW_ID_2103205055=1; cdb_RW_ID_2103197728=1; cdb_RW_ID_2103197715=1; cdb_RW_ID_2103197703=1; cdb_RW_ID_2103117030=1; cdb_RW_ID_2103117090=1; cdb_RW_ID_2103117115=1; cdb_RW_ID_2103117127=1; cdb_RW_ID_2103117133=1; cdb_RW_ID_2102995720=1; cdb_RW_ID_2102995546=1; cdb_RW_ID_2102995524=1; cdb_RW_ID_2102918235=1; cdb_RW_ID_2102917936=1; cdb_RW_ID_2102917892=1; cdb_RW_ID_2102829316=1; cdb_RW_ID_2102569655=1; cdb_RW_ID_2102446467=1; cdb_RW_ID_2102446462=2; cdb_RW_ID_2102446465=2; cdb_RW_ID_2060083101=9; cdb_RW_ID_2102446205=1; cdb_RW_ID_2100959849=3; cdb_RW_ID_2102446200=1; cdb_RW_ID_1961805838=101; cdb_RW_ID_1061054128=16; cdb_RW_ID_2102282735=1; cdb_RW_ID_2102245257=1; cdb_RW_ID_2102245260=1; cdb_RW_ID_2102245261=1; cdb_RW_ID_2102245264=1; cdb_RW_ID_2102245265=1; cdb_RW_ID_2101241003=2; cdb_RW_ID_2101241006=1; cdb_RW_ID_2101241042=1; cdb_RW_ID_2101241037=1; cdb_RW_ID_2101220439=1; cdb_RW_ID_2101162155=1; cdb_RW_ID_2101162099=1; cdb_RW_ID_2101159289=1; cdb_RW_ID_2100253415=1; cdb_RW_ID_2098480023=1; cdb_RW_ID_2098985940=1; cdb_RW_ID_2099075055=1; cdb_RW_ID_2099838516=1; cdb_RW_ID_2099728766=1; cdb_RW_ID_2098816864=1; cdb_RW_ID_2098816042=1; cdb_RW_ID_2098530397=1; cdb_RW_ID_2098527001=1; cdb_RW_ID_2098481524=1; cdb_RW_ID_2098445808=1; cdb_RW_ID_2098113130=1; cdb_RW_ID_2097289310=1; cdb_RW_ID_2097289262=1; cdb_RW_ID_2097289237=1; cdb_RW_ID_2097289243=1; cdb_RW_ID_1000248695=1114; cdb_RW_ID_1398856507=19; cdb_RW_ID_1021305450=168; cdb_RW_ID_2096381365=1; cdb_RW_ID_1070336488=1001; cdb_RW_ID_2095955907=1; cdb_RW_ID_2095955375=1; cdb_RW_ID_2095791489=1; cdb_RW_ID_2095370303=1; cdb_RW_ID_2094978049=1; cdb_RW_ID_2094735877=1; cdb_RW_ID_2094735021=1; cdb_RW_ID_2094606879=1; cdb_RW_ID_2094603147=1; cdb_RW_ID_2094603016=1; cdb_RW_ID_2093653040=1; cdb_RW_ID_2093558888=1; cdb_RW_ID_2054082459=12; cdb_RW_ID_2093557630=1; cdb_RW_ID_2093557607=1; cdb_RW_ID_2093409660=3; cdb_RW_ID_2093409648=1; cdb_RW_ID_2092774793=1; cdb_RW_ID_2092773979=1; cdb_RW_ID_2002946432=1; cdb_RW_ID_2002913089=1; cdb_RW_ID_2092774082=1; cdb_RW_ID_2092632858=1; cdb_RW_ID_2092167051=1; cdb_RW_ID_2090223541=1; cdb_RW_ID_2090004000=1; cdb_RW_ID_2090004031=1; cdb_RW_ID_2089750151=1; cdb_RW_ID_2089370858=2; cdb_RW_ID_2051414775=2; cdb_RW_ID_2088494889=1; cdb_RW_ID_2088150740=1; cdb_RW_ID_2080284856=3; cdb_RW_ID_2087899576=1; cdb_RW_ID_2087886059=1; cdb_RW_ID_2087742144=1; cdb_RW_ID_1000625924=16; cdb_RW_ID_2087303907=1; cdb_RW_ID_2086920586=1; cdb_RW_ID_2086914432=1; cdb_RW_ID_1098251921=3; cdb_RW_ID_2086914441=1; cdb_RW_ID_42897619=38; cdb_RW_ID_2086914478=1; cdb_RW_ID_2086914544=1; cdb_RW_ID_1801942=163; cdb_RW_ID_2086914578=1; cdb_RW_ID_1017986246=826; cdb_RW_ID_2086914630=1; cdb_RW_ID_2086914646=1; cdb_RW_ID_755052518=9; cdb_RW_ID_2086914663=1; cdb_RW_ID_227313860=16; cdb_RW_ID_2086914666=1; cdb_RW_ID_2086915013=1; cdb_RW_ID_1434875925=7; cdb_RW_ID_2086915035=1; cdb_RW_ID_2086915134=1; cdb_RW_ID_2085080900=1; cdb_RW_ID_2044145934=7; cdb_RW_ID_2083654732=1; cdb_RW_ID_2083427874=2; cdb_RW_ID_2083427071=1; cdb_RW_ID_2083246914=1; cdb_RW_ID_2083106195=1; cdb_RW_ID_2083037294=1; cdb_RW_ID_2082459282=1; cdb_RW_ID_2077244680=1; cdb_RW_ID_2077244681=1; cdb_RW_ID_2082431887=1; cdb_RW_ID_2082261983=1; cdb_RW_ID_2014447320=2; cdb_RW_ID_2082262146=1; cdb_RW_ID_2082262163=1; cdb_RW_ID_2080370567=1; cdb_RW_ID_2082262152=1; cdb_RW_ID_2066793975=2; cdb_RW_ID_2081772457=1; cdb_RW_ID_2081041522=1; cdb_RW_ID_2061315675=1; cdb_RW_ID_2081041525=1; cdb_RW_ID_2081041558=1; cdb_RW_ID_2081040702=1; cdb_RW_ID_2080704952=9; cdb_RW_ID_2059879162=1; cdb_RW_ID_2069750564=1; cdb_RW_ID_2080873081=1; cdb_RW_ID_2080682698=1; cdb_RW_ID_2070959060=39"

func getKey() (uploadKey string, err error) {
	client := &http.Client{} //初始化客户端
	getKeyUrl := "https://www.doc88.com/uc/index.php?act=getkey"
	req, err := http.NewRequest("GET", getKeyUrl, nil) //建立连接
	if err != nil {
		return uploadKey, err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", Doc88Cookie)
	req.Header.Set("Host", "www.doc88.com")
	req.Header.Set("Referer", "https://www.doc88.com/uc/index.php?act=upload")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return uploadKey, err
	}
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return uploadKey, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return uploadKey, err
	}
	uploadKey = string(s)
	return uploadKey, nil
}

func reverseString(s string) string {
	b := []byte(s)
	n := len(b)
	for i := 0; i < n/2; i++ {
		b[i], b[n-i-1] = b[n-i-1], b[i]
	}
	return string(b)

}

type UploadResponseData struct {
	Result  string `json:"result"`
	Message string `json:"message"`
	DocCode string `json:"doccode"`
	PId     string `json:"p_id"`
	PCode   string `json:"p_code"`
}

// 上传文件
func uploadFile(filePath string, uploadKey string, ck string) (uploadResponseData UploadResponseData, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	uploadResponseData = UploadResponseData{}

	file, err := os.Open(filePath)
	if err != nil {
		return uploadResponseData, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	fileWriter, err := bodyWriter.CreateFormFile("upfile", filepath.Base(file.Name()))
	if err != nil {
		return uploadResponseData, err
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return uploadResponseData, err
	}

	contentType := bodyWriter.FormDataContentType()
	err = bodyWriter.Close()
	if err != nil {
		return UploadResponseData{}, err
	}

	err = bodyWriter.WriteField("act", "upload")
	if err != nil {
		return uploadResponseData, err
	}

	err = bodyWriter.WriteField("fileName", filepath.Base(file.Name()))
	if err != nil {
		return uploadResponseData, err
	}

	uploadUrl := fmt.Sprintf("https://upload.doc88.com/u.do?v=1&uploadkey=%s&ck=%s", uploadKey, ck)
	req, err := http.NewRequest("POST", uploadUrl, bodyBuf)
	if err != nil {
		return uploadResponseData, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Host", "upload.doc88.com")
	req.Header.Set("Origin", "https://www.doc88.com")
	req.Header.Set("Referer", "https://www.doc88.com/")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return uploadResponseData, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	respBytes, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBytes, &uploadResponseData)
	if err != nil {
		return uploadResponseData, err
	}
	return uploadResponseData, nil
}

type EditResponseData struct {
	Result     string `json:"result"`
	EditTitle  string `json:"edit_title"`
	Class      string `json:"class"`
	UpdateInfo string `json:"updateinfo"`
	State      string `json:"state"`
	SaveFile   string `json:"savefile"`
	Other      string `json:"other"`
}

// 编辑文件所属分类和下载所需积分
func editFile(docCode string, title string, intro string, pCid int, price int, pDocFormat string) (editResponseData EditResponseData, err error) {
	client := &http.Client{}
	editResponseData = EditResponseData{}
	postData := url.Values{}
	postData.Add("doccode", docCode)
	postData.Add("title", title)
	postData.Add("intro", intro)
	postData.Add("pcid", strconv.Itoa(pCid))
	postData.Add("keyword", "")
	postData.Add("sharetodoc", "1")
	postData.Add("download", "2")
	postData.Add("p_price", strconv.Itoa(price))
	postData.Add("p_default_points", "1")
	postData.Add("p_pagecount", "")
	postData.Add("p_doc_format", pDocFormat)
	postData.Add("act", "save_info")
	postData.Add("group_list", "")
	postData.Add("group_free_list", "")
	requestUrl := "https://www.doc88.com/uc/index.php"
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return editResponseData, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", Doc88Cookie)
	req.Header.Set("Host", "www.doc88.com")
	req.Header.Set("Origin", "https://www.doc88.com")
	req.Header.Set("Referer", "https://www.doc88.com/uc/index.php?act=upload")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return editResponseData, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return editResponseData, err
	}
	err = json.Unmarshal(respBytes, &editResponseData)
	if err != nil {
		return editResponseData, err
	}
	return editResponseData, nil
}

type UploadChildDir struct {
	dirName string
	pCid    int
	Price   int
}

func main() {
	var uploadChildDirArr = []UploadChildDir{
		{
			dirName: "dbba.sacinfo.org.cn",
			pCid:    8371,
			Price:   688,
		},
		{
			dirName: "hbba.sacinfo.org.cn",
			pCid:    8370,
			Price:   688,
		},
		{
			dirName: "docx.topedu.ybep.com.cn/小学试卷",
			pCid:    8243,
			Price:   188,
		},
		{
			dirName: "docx.topedu.ybep.com.cn/中考试卷",
			pCid:    8243,
			Price:   188,
		},
		{
			dirName: "docx.topedu.ybep.com.cn/高考试卷",
			pCid:    8244,
			Price:   188,
		},
		{
			dirName: "www.docx_shijuan1.com/中考试卷",
			pCid:    8243,
			Price:   188,
		},
		{
			dirName: "www.docx_shijuan1.com/高考试卷",
			pCid:    8244,
			Price:   188,
		},
		{
			dirName: "www.docx_shijuan1.com/小学试卷",
			pCid:    8242,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/小学/道德与法治",
			pCid:    8307,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/小学/美术",
			pCid:    8309,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/小学/数学",
			pCid:    8247,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/小学/信息技术",
			pCid:    8308,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/小学/音乐",
			pCid:    8309,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/小学/英语",
			pCid:    8248,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/小学/语文",
			pCid:    8246,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/初中/数学",
			pCid:    8156,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/初中/英语",
			pCid:    8157,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/初中/语文",
			pCid:    8155,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/初中/物理",
			pCid:    8158,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/初中/化学",
			pCid:    8160,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/初中/音乐",
			pCid:    55918,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/初中/道德与法治",
			pCid:    8163,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/初中/美术",
			pCid:    55918,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/初中/生物",
			pCid:    8162,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/初中/地理",
			pCid:    8159,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/初中/历史",
			pCid:    8161,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/高中/地理",
			pCid:    8256,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/高中/化学",
			pCid:    8183,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/高中/历史",
			pCid:    8257,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/高中/生物",
			pCid:    8186,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/高中/数学",
			pCid:    8180,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/高中/物理",
			pCid:    8182,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/高中/英语",
			pCid:    8181,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/高中/语文",
			pCid:    8179,
			Price:   188,
		},
		{
			dirName: "docx.zzstep.com/高中/政治",
			pCid:    8258,
			Price:   188,
		},
		{
			dirName: "docx.gzenxx.com",
			pCid:    55919,
			Price:   188,
		},
		{
			dirName: "docx.trjlseng.com",
			pCid:    8157,
			Price:   188,
		},
		{
			dirName: "sc.chinaz.com/ppt",
			pCid:    8120,
			Price:   188,
		},
		{
			dirName: "docx.www.czsx.com.cn",
			pCid:    8156,
			Price:   188,
		},
		{
			dirName: "www.finish_zip_1ppt.com",
			pCid:    8120,
			Price:   188,
		},
		{
			dirName: "www.docx_ychedu.com/中考试题",
			pCid:    8243,
			Price:   188,
		},
		{
			dirName: "www.docx_ychedu.com/高考试题",
			pCid:    8244,
			Price:   188,
		},
		{
			dirName: "www.docx_ychedu.com/中考试题",
			pCid:    8243,
			Price:   188,
		},
		{
			dirName: "www.docx_ychedu.com/高考试题",
			pCid:    8244,
			Price:   188,
		},
		{
			dirName: "hao123.docx_eduy.net/小学",
			pCid:    8242,
			Price:   188,
		},
		{
			dirName: "hao123.docx_eduy.net/初中",
			pCid:    8243,
			Price:   188,
		},
		{
			dirName: "hao123.docx_eduy.net/高中",
			pCid:    8244,
			Price:   188,
		},
		{
			dirName: "www.mee.gov.cn",
			pCid:    8368,
			Price:   188,
		},
		{
			dirName: "www.docx_czwlzx.cn",
			pCid:    8158,
			Price:   188,
		},
		{
			dirName: "htsfwb.samr.gov.cn",
			pCid:    8131,
			Price:   188,
		},
		{
			dirName: "flk.npc.gov.cn",
			pCid:    8095,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/小学/语文",
			pCid:    8246,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/小学/科学",
			pCid:    8253,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/小学/道德与法治",
			pCid:    8307,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/小学/数学",
			pCid:    8247,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/小学/英语",
			pCid:    8248,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/初中/数学",
			pCid:    8156,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/初中/英语",
			pCid:    8157,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/初中/语文",
			pCid:    8155,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/初中/化学",
			pCid:    8160,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/初中/体育与健康",
			pCid:    55918,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/初中/信息技术",
			pCid:    55918,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/初中/历史",
			pCid:    8161,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/初中/历史与社会",
			pCid:    8161,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/初中/地理",
			pCid:    8159,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/初中/物理",
			pCid:    8158,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/初中/生物",
			pCid:    8162,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/初中/科学",
			pCid:    55918,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/初中/综合",
			pCid:    55918,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/初中/道德与法治",
			pCid:    8163,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/初中/美术",
			pCid:    55918,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/初中/音乐",
			pCid:    55918,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/高中/地理",
			pCid:    8256,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/高中/化学",
			pCid:    8183,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/高中/历史",
			pCid:    8257,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/高中/生物",
			pCid:    8186,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/高中/数学",
			pCid:    8180,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/高中/通用技术",
			pCid:    55919,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/高中/物理",
			pCid:    8182,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/高中/信息技术",
			pCid:    55919,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/高中/英语",
			pCid:    8181,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/高中/语文",
			pCid:    8179,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/高中/政治",
			pCid:    8258,
			Price:   188,
		},
		{
			dirName: "docx.zhuangyuan123.com/高中/综合",
			pCid:    55919,
			Price:   188,
		},
		{
			dirName: "docx.guojiajiaoyu.com/学前",
			pCid:    8304,
			Price:   188,
		},
		{
			dirName: "docx.guojiajiaoyu.com/小学",
			pCid:    8309,
			Price:   188,
		},
		{
			dirName: "docx.guojiajiaoyu.com/初中",
			pCid:    55918,
			Price:   188,
		},
		{
			dirName: "docx.guojiajiaoyu.com/高中",
			pCid:    55919,
			Price:   188,
		},
		{
			dirName: "www.meewen.com/表格合同",
			pCid:    8131,
			Price:   188,
		},
		{
			dirName: "www.meewen.com/遴选题库",
			pCid:    8224,
			Price:   188,
		},
		{
			dirName: "www.meewen.com/领导讲话",
			pCid:    8138,
			Price:   188,
		},
		{
			dirName: "www.meewen.com/教辅资料",
			pCid:    55918,
			Price:   188,
		},
		{
			dirName: "tikuvip.51test.net/中考",
			pCid:    8243,
			Price:   188,
		},
		{
			dirName: "tikuvip.51test.net/初中一年级",
			pCid:    8243,
			Price:   188,
		},
		{
			dirName: "tikuvip.51test.net/初中二年级",
			pCid:    8243,
			Price:   188,
		},
		{
			dirName: "tikuvip.51test.net/初中三年级",
			pCid:    8243,
			Price:   188,
		},
		{
			dirName: "tikuvip.51test.net/高考",
			pCid:    8244,
			Price:   188,
		},
		{
			dirName: "tikuvip.51test.net/高一",
			pCid:    8244,
			Price:   188,
		},
		{
			dirName: "tikuvip.51test.net/高二",
			pCid:    8244,
			Price:   188,
		},
		{
			dirName: "tikuvip.51test.net/高三",
			pCid:    8244,
			Price:   188,
		},
		{
			dirName: "www.jiangyuanwang.com/word",
			pCid:    8131,
			Price:   188,
		},
		{
			dirName: "www.jiangyuanwang.com/excel",
			pCid:    8119,
			Price:   188,
		},
		{
			dirName: "www.jiangyuanwang.com/ppt",
			pCid:    8120,
			Price:   188,
		},
		{
			dirName: "www.docx_1shijuan.cn/一年级",
			pCid:    8251,
			Price:   188,
		},
		{
			dirName: "www.docx_1shijuan.cn/二年级",
			pCid:    8251,
			Price:   188,
		},
		{
			dirName: "www.docx_1shijuan.cn/三年级",
			pCid:    8251,
			Price:   188,
		},
		{
			dirName: "www.docx_1shijuan.cn/四年级",
			pCid:    8251,
			Price:   188,
		},
		{
			dirName: "www.docx_1shijuan.cn/五年级",
			pCid:    8251,
			Price:   188,
		},
		{
			dirName: "www.docx_1shijuan.cn/六年级",
			pCid:    8251,
			Price:   188,
		},
		{
			dirName: "www.docx_1shijuan.cn/七年级",
			pCid:    55918,
			Price:   188,
		},
		{
			dirName: "www.docx_1shijuan.cn/八年级",
			pCid:    55918,
			Price:   188,
		},
		{
			dirName: "www.docx_1shijuan.cn/九年级",
			pCid:    55918,
			Price:   188,
		},
		{
			dirName: "www.docx_1shijuan.cn/高一",
			pCid:    55919,
			Price:   188,
		},
		{
			dirName: "www.docx_1shijuan.cn/高二",
			pCid:    55919,
			Price:   188,
		},
		{
			dirName: "www.docx_1shijuan.cn/高三",
			pCid:    55919,
			Price:   188,
		},
		{
			dirName: "www.docx_1shijuan.cn/中考",
			pCid:    8243,
			Price:   188,
		},
		{
			dirName: "www.docx_1shijuan.cn/高考",
			pCid:    8244,
			Price:   188,
		},
	}
	rootPath := "../upload.doc88.com/"
	for _, childDir := range uploadChildDirArr {
		childDirPath := rootPath + childDir.dirName + "/"
		fmt.Println(childDirPath)
		files, err := ioutil.ReadDir(childDirPath)
		if err != nil {
			continue
		}
		shuffleArray(files)
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			fileName := file.Name()
			if fileName == ".DS_Store" {
				continue
			}
			fileExt := path.Ext(fileName)

			fmt.Println("==========开始上传==============")
			filePath := childDirPath + fileName
			fmt.Println(filePath)

			uploadKey, err := getKey()
			if err != nil {
				fmt.Println(err)
				continue
			}
			ck := reverseString(uploadKey[3:8])
			if err != nil {
				fmt.Println(err)
				continue
			}
			uploadResponseData, err := uploadFile(filePath, uploadKey, ck)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println(uploadResponseData)

			if result, _ := strconv.Atoi(uploadResponseData.Result); result != 0 {
				fmt.Println(uploadResponseData.Message)
				break
			}
			// 设置倒计时
			UploadDoc88TimeSleep := rand.Intn(10)
			if UploadDoc88TimeSleep <= 8 {
				UploadDoc88TimeSleep = 8 + rand.Intn(5)
			}
			for i := 1; i <= UploadDoc88TimeSleep; i++ {
				time.Sleep(time.Second)
				fmt.Println("暂停", UploadDoc88TimeSleep, "秒后编辑文件所属类别和下载积分，倒计时", i, "秒===========")
			}

			// 处理编辑文件时所需参数值
			docCode := uploadResponseData.DocCode
			title := strings.ReplaceAll(fileName, fileExt, "")
			title = strings.Replace(title, "\"", "“", 1)
			title = strings.Replace(title, "\"", "”", 1)
			intro := title
			pCid := childDir.pCid
			price := childDir.Price
			pDocFormat := ""
			switch fileExt {
			case ".pdf":
				pDocFormat = "PDF"
			case ".docx":
				pDocFormat = "DOCX"
			case ".doc":
				pDocFormat = "DOC"
			case ".ppt":
				pDocFormat = "PPT"
			case ".pptx":
				pDocFormat = "PPTX"
			case ".xls":
				pDocFormat = "XLS"
			case ".xlsx":
				pDocFormat = "XLSX"
			}

			filePageNum := 0
			if fileExt == ".pdf" {
				// 获取PDF文件，获取总页数
				if pdfFile, err := pdf.Open(filePath); err == nil {
					filePageNum = pdfFile.NumPage()
				}
			}
			//根据页数设置价格
			if filePageNum > 0 {
				if filePageNum <= 25 {
					price = 688
				} else if filePageNum > 25 && filePageNum <= 30 {
					price = 788
				} else if filePageNum > 30 && filePageNum <= 35 {
					price = 888
				} else {
					price = 988
				}
			}

			// 将已上传的文件转移到指定文件夹
			// 将上传过文件移动到"../final-upload.doc88.com/"
			finalDir := "../final-upload.doc88.com/" + childDir.dirName
			if _, err = os.Stat(finalDir); err != nil {
				if os.MkdirAll(finalDir, 0777) != nil {
					fmt.Println(err)
					continue
				}
			}
			fileFinal := finalDir + "/" + fileName
			if err := os.Rename(filePath, fileFinal); err != nil {
				// 如果文件正在被使用，等待一段时间后重试
				time.Sleep(time.Second)
				if err := os.Rename(filePath, fileFinal); err != nil {
					fmt.Println(err)
					continue
				}
			}

			// 编辑文件所需分类和下载所需积分
			editResponseData, err := editFile(docCode, title, intro, pCid, price, pDocFormat)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(editResponseData)

			fmt.Println("==========上传完成==============")
		}
	}
}

func shuffleArray(arr []fs.FileInfo) {
	rand.Seed(time.Now().UnixNano())
	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}
