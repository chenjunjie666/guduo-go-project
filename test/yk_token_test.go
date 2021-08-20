package test

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestToken(tt *testing.T) {
	//c := colly.NewCollector()

}

func Test2(tt *testing.T) {
	token := strings.Split(getCookies()[0].Value, "_")[0]

	pdataJson, _ := json.Marshal(post)

	msg := base64.StdEncoding.EncodeToString(pdataJson)
	sign_1 := md5.Sum([]byte(msg + key))

	pdata := string(pdataJson)
	pdata = strings.TrimRight(pdata, "}")
	pdata = fmt.Sprintf(`%s,"msg":"%s","sign":"%x"}`, pdata, msg, sign_1)

	t := int(time.Now().Unix()*1000)
	str := fmt.Sprintf("%s&%d&%s&%s", token, t, appKey, pdata)
	sign_2 := md5.Sum([]byte(str))

	params := "?"
	for k, v := range get {
		if u, ok := v.(int); ok {
			params += fmt.Sprintf("%s=%d&", k, u)
		}else if u, ok := v.(string); ok {
			params += fmt.Sprintf("%s=%s&", k, u)
		}
	}

	params += fmt.Sprintf("t=%d&sign=%x", t, sign_2)

	apiUrl := "https://acs.youku.com/h5/mopen.youku.danmu.list/1.0/" + params

	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Cookie", BuildCookie(getCookies()))
		r.Headers.Add("Content-type", "application/x-www-form-urlencoded")
		r.Headers.Add("Refer", "https://v.youku.com/")
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(string(r.Body))
		fmt.Println("-------------------")
		fmt.Println(r.Request.Headers)
		fmt.Println("-------------------")

		if len(r.Headers.Values("Set-Cookie")) > 0 {
			fmt.Println(r.Headers.Values("Set-Cookie")[0])
			fmt.Println(r.Headers.Values("Set-Cookie")[1])
		}
	})
	_ = c.Post(apiUrl, map[string]string{"data": pdata})
}

func TestGetToken(t *testing.T) {
	tt, ttt := GenToken()

	fmt.Println(tt)
	fmt.Println(ttt)
}

func GenToken() (string, string) {
	c := colly.NewCollector()

	tk := ""
	tk_ := ""

	c.OnResponse(func(r *colly.Response) {
		cookie1 := r.Headers.Values("Set-Cookie")[0]
		cookie2 := r.Headers.Values("Set-Cookie")[1]

		cookie1 = strings.Split(cookie1, ";")[0]
		cookie2 = strings.Split(cookie2, ";")[0]

		name1 := strings.Split(cookie1, "=")[0]
		cookie1 = strings.Split(cookie1, "=")[1]

		cookie2 = strings.Split(cookie2, "=")[1]

		if name1 == "_m_h5_tk"{
			tk = cookie1
			tk_ = cookie2
		}else{
			tk_ = cookie1
			tk = cookie2
		}
	})
	c.Visit("https://acs.youku.com/h5/mtop.youku.play.ups.appinfo.get/1.1/?jsv=2.4.16&appKey=24679788")

	return tk, tk_
}

func getCookies() []http.Cookie {
	cook := []http.Cookie{
		{
			Name:  "_m_h5_tk",
			Value: "ddf2dee82b6506a3332b38f6a16a8ff6_1619278576399",
		},
		{
			Name:  "_m_h5_tk_enc",
			Value: "725d7bfc9d07d99e4477372b7eca19c7",
		},
	}

	return cook
}

func BuildCookie(cookie []http.Cookie) string {
	var c string

	for _, ck := range cookie {
		tmpCookie := ck.Name + "=" + ck.Value + ";"
		c += tmpCookie
	}

	return c
}

var get = map[string]interface{}{
	"jsv":"2.5.1",
	"appKey":24679788,
	"api":"mopen.youku.danmu.list",
	"v":"1.0",
	"type":"originaljson",
	"dataType":"jsonp",
	"timeout":20000,
	"jsonpIncPrefix":"utility",
}

var post = map[string]interface{}{
	"pid":0,
	"ctype":10004,
	"sver":"3.1.0",
	"cver":"v1.0",
	"ctime": int(time.Now().Unix()*1000),
	"guid":"dUsHGVWmT2ICAXTtJ0bXd4wn",
	"vid":"XNTEzMTk0ODIyNA==",
	"mat":22,
	"mcount":1,
	"type":1,
}

var key = "MkmC9SoIw6xCkSKHhJ7b5D2r51kBiREr"

var appKey = "24679788"
