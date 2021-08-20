package test

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"guduo/pkg/util"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestEpid(ttttttt *testing.T) {
	token := strings.Split(getCookies2()[0].Value, "_")[0]

	nextSessionJson, _ := json.Marshal(nextSession)
	params["nextSession"] = string(nextSessionJson)
	paramsJson, _ := json.Marshal(params)
	systemInfoJson, _ := json.Marshal(systemInfo)

	data["params"] = string(paramsJson)
	data["system_info"] = string(systemInfoJson)

	dataJson, _ := json.Marshal(data)

	//dataStr := `{"ms_codes":"2019030100","params":"{\"biz\":true,\"scene\":\"component\",\"componentVersion\":\"3\",\"ip\":\"116.233.61.29\",\"debug\":0,\"utdid\":\"yuP7GGf9hmkCAXTpPR3QoHbe\",\"userId\":\"\",\"platform\":\"pc\",\"nextSession\":\"{\\\"componentIndex\\\":\\\"3\\\",\\\"componentId\\\":\\\"61518\\\",\\\"level\\\":\\\"2\\\",\\\"itemPageNo\\\":\\\"0\\\",\\\"lastItemIndex\\\":\\\"0\\\",\\\"pageKey\\\":\\\"LOGICSHOW_LOGICTV_DEFAULT\\\",\\\"group\\\":\\\"0\\\",\\\"itemStartStage\\\":1,\\\"itemEndStage\\\":30}\",\"videoId\":\"XNTEwMzgzMDQ4MA==\",\"showId\":\"ceba4745ea10415eb791\"}","system_info":"{\"os\":\"pc\",\"device\":\"pc\",\"ver\":\"1.0.0\",\"appPackageKey\":\"pcweb\",\"appPackageId\":\"pcweb\"}"}`
	//fmt.Println("+++", string(dataJson))
	//fmt.Println("---", dataStr)
	t := int(time.Now().Unix()*1000)
	//t := 1619162367638
	str := fmt.Sprintf("%s&%d&%s&%s", token, t, appKey, string(dataJson))
	sign := md5.Sum([]byte(str))


	//str2 := fmt.Sprintf("%s&%d&%s&%s", token, t, appKey, string(dataJson))
	//sign2 := md5.Sum([]byte(str2))

	//fmt.Printf("%x\n", sign)
	//fmt.Printf("%x\n", sign2)
	apiUrl := fmt.Sprintf(
		"https://acs.youku.com/h5/mtop.youku.columbus.gateway.new.execute/1.0/?jsv=2.6.1&appKey=24679788&t=%d&sign=%x&api=mtop.youku.columbus.gateway.new.execute&type=originaljson&v=1.0&ecode=1&dataType=json&data=%s",
		t,
		sign,
		dataJson,
		)

	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Cookie", BuildCookie(getCookies2()))
		r.Headers.Add("Content-type", "application/x-www-form-urlencoded")
		r.Headers.Add("Refer", "https://v.youku.com/")
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(string(r.Body))
		//fmt.Println("-------------------")
		//fmt.Println(r.Request.Headers)
		//fmt.Println("-------------------")

		if len(r.Headers.Values("Set-Cookie")) > 0 {
			fmt.Println(r.Headers.Values("Set-Cookie")[0])
			fmt.Println(r.Headers.Values("Set-Cookie")[1])
		}
	})
	e := c.Visit(apiUrl)

	fmt.Println(e)
}

func getCookies2() []http.Cookie {
	cook := []http.Cookie{
		{
			Name:  "_m_h5_tk",
			Value: "73fba4de0468720d63df5691b03b9e68_1619196462116",
		},
		{
			Name:  "_m_h5_tk_enc",
			Value: "05f84730128389bae6fa519e80f146a9",
		},
	}

	return cook
}


var get_ = map[string]interface{}{
	"jsv": "2.6.1",
	"appKey": 24679788,
	"t": 1619162367638,
	"sign": "cb4a675349c8703b69ba23c2c9f3afb7",
	"api": "mtop.youku.columbus.gateway.new.execute",
	"type": "originaljson",
	"v": 1.0,
	"ecode": 1,
	"dataType": "json",
	"data": "",
}

var data = map[string]interface{}{
	"ms_codes": "2019030100",
	"params": "",
	"system_info": "",
}

var params = map[string]interface{}{
	"biz":true,
	"scene":"component",
	"componentVersion":"3",
	"ip":util.GenRandomIP(),
	"debug":0,
	"utdid":"yuP7GGf9hmkCAXTpPR3QoHbe",
	"userId":"",
	"platform":"pc",
	"nextSession": "",
	"videoId":"XNTEwMzgzMDQ4MA==",
	"showId":"ceba4745ea10415eb791",
}

var nextSession = map[string]interface{}{
	"componentIndex":"3",
	"componentId":"61518",
	"level":"2",
	"itemPageNo":"0",
	"lastItemIndex":"0",
	"pageKey":"LOGICSHOW_LOGICTV_DEFAULT",
	"group":"0",
	"itemStartStage":1,
	"itemEndStage":30,
}

var systemInfo = map[string]string{
	"os":"pc",
	"device":"pc",
	"ver":"1.0.0",
	"appPackageKey":"pcweb",
	"appPackageId":"pcweb",
}