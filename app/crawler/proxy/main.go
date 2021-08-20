package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/proxy/controller"
	"strings"
	time2 "time"
)

func main() {
	//checkProxy()
	g := gin.Default()

	g.GET("/get_proxy", controller.FetchProxy)

	g.GET("/report_useless", controller.ReportErrorProxy)

	go controller.CheckProxy()

	g.Run("0.0.0.0:93")
	//fmt.Println(fetch)

	//idx := 1
	//for  {
	//	fmt.Println(idx)
	//	idx++
	//	fetch()
	//	time2.Sleep(time2.Second)
	//}
	//useOps()
}


//func useOps() {
//	auth := auth2.Auth{OrderID: "909005777021989", APIKey: "rxggagj4gdmdvwizjna5gnync1eh71qb"}
//	client2 := client.Client{Auth: auth}
//
//	// 获取订单到期时间, 返回时间字符串
//	expireTime, err := client2.GetOrderExpireTime(signtype.HmacSha1)
//	if err != nil {
//		log.Println(err)
//	}
//	fmt.Println("expire time: ", expireTime)
//
//	// 提取开放代理, 参数有: 提取数量、开放代理套餐种类(normal, vip, svip, ent
//	// 	分别对应opslevel.NORMAL, opslevel.VIP, opslevel.SVIP, opslevel.ENT)、
//	// 鉴权方式及其他参数(放入map[string]interface{}中, 若无则传入nil)
//	// (具体有哪些其他参数请参考帮助中心: "https://www.kuaidaili.com/doc/api/getops/")
//	//params := map[string]interface{}{"area": "北京,上海"}
//	ips, err := client2.GetProxy(2, opslevel.SVIP, signtype.HmacSha1, nil)
//	if err != nil {
//		log.Println(err)
//	}
//	fmt.Println("ips: ", ips)
//
//	// 检测开放代理有效性， 返回map[string]bool, ip:true/false
//	valids, err := client2.CheckOpsValid(ips, signtype.HmacSha1)
//	if err != nil {
//		log.Println(err)
//	}
//	fmt.Println("valids: ", valids)
//}


func checkProxy(){
	t := time2.Now().Unix()

	c2 := colly.NewCollector()
	urls := []string{
		"113.75.139.105:15787",
		"125.105.19.139:18378",
	}

	url := fmt.Sprintf("https://dps.kdlapi.com/api/checkdpsvalid?orderid=998752579990468&sign_type=simple&signature=ueqy8qfyrhwov9c8f5ksnkxlrnk8eo33&proxy=%s&timestamp=%d",
		strings.Join(urls, ","),
		t,
	)
	fmt.Println("222")
	c2.OnResponse(func(r *colly.Response) {
		fmt.Println(string(r.Body))
	})


	c2.Visit(url)
}

func fetch() {
	t := time2.Now().Unix()
	apiUrl := fmt.Sprintf("http://svip.kdlapi.com/api/getproxy/?orderid=909005777021989&sign_type=simple&signature=rxggagj4gdmdvwizjna5gnync1eh71qb" +
		"&timestamp=%d&num=100&protocol=1&method=1&quality=2&format=json&sep=1",
		t,
	)
	c := colly.NewCollector()

	urls := make([]string, 0, 100)
	c.OnResponse(func(r *colly.Response) {
		code, e := jsonparser.GetInt(r.Body, "code")
		msg, _ := jsonparser.GetString(r.Body, "msg")
		if e != nil || code != 0 {
			log.Error(fmt.Sprintf("获取代理失败，返回code：%d, msg:%s, 错误是：%s", code, msg, e))
			return
		}
		_, _ = jsonparser.ArrayEach(r.Body, func(v []byte, dataType jsonparser.ValueType, offset int, err error) {
			ip := string(v)
			url := fmt.Sprintf("http://%s", ip)
			// 用当前时间+90s做判断

			urls = append(urls, url)

			// 检测IP是否已经存在
		}, "data", "proxy_list")
	})

	_ = c.Visit(apiUrl)

	c2 := colly.NewCollector()
	url := fmt.Sprintf("https://dev.kdlapi.com/api/checkopsvalid?orderid=909005777021989&sign_type=simple&signature=rxggagj4gdmdvwizjna5gnync1eh71qb&proxy=%s&timestamp=%d",
		strings.Join(urls, ","),
		t,
		)

	c2.OnResponse(func(r *colly.Response) {
		jsonparser.ObjectEach(r.Body, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			if string(value) != "false" {
				fmt.Println(string(key), string(value))
			}
			return nil
		}, "data")
	})


	c2.Visit(url)
}