package youku

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"guduo/app/crawler/spider/internal/util"
	"net/http"
	"testing"
)

func TestSignDanmakuParams(t *testing.T) {
	url, pdata := GetDanmakuUrl("XNDI2NTg3NDc4OA==", 1)

	MH5Tk, MH5TkEnc := GetToken()
	fmt.Println(MH5Tk)

	// 爬虫配置
	c := colly.NewCollector()

	// 设置 cookie
	cookies := getCookies(MH5Tk, MH5TkEnc)

	// 设置 cookie
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", util.BuildCookie(cookies))
		r.Headers.Add("Content-type", "application/x-www-form-urlencoded")
		r.Headers.Add("Refer", "https://v.youku.com/")
	})

	extensions.RandomUserAgent(c) // 随机设置 user-agent

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(string(r.Body))
	})

	e := c.Post(url, pdata)
	fmt.Println(e)
}


func getCookies(MH5Tk, MH5TkEnc string) []http.Cookie {
	cookies := []http.Cookie{
		{
			Name:  "_m_h5_tk",
			Value: MH5Tk,
		},
		{
			Name:  "_m_h5_tk_enc",
			Value: MH5TkEnc,
		},
	}

	return cookies
}