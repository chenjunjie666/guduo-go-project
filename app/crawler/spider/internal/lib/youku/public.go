package youku

import (
	"github.com/gocolly/colly/v2"
	"strconv"
	"strings"
	"sync"
	"time"
)

var key = "MkmC9SoIw6xCkSKHhJ7b5D2r51kBiREr"
var appKey = 24679788
var MH5Tk = ""
var MH5TkEnc = ""

func GetAppKey() int {
	return appKey
}


var tokenLock = &sync.Mutex{}

// 获取优酷token
func GetToken() (string, string){
	tokenLock.Lock()
	defer tokenLock.Unlock()

	if MH5Tk != "" && MH5TkEnc != "" {
		limit, _ := strconv.ParseInt(strings.Split(MH5Tk, "_")[1], 10, 64)
		now := time.Now().UnixNano() / 1e6
		// 如果过期时间大于20分钟，才返回
		if limit > 0 && now - limit < -1000 * 60 * 20 {
			return MH5Tk, MH5TkEnc
		}
	}

	MH5Tk, MH5TkEnc = genToken()
	return MH5Tk, MH5TkEnc
}


func genToken() (string, string) {
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
	c.Visit("https://acs.youku.com/h5/mtop.youku.play.ups.appinfo.get/1.1/?jsv=2.4.16&appKey=" + strconv.Itoa(appKey))

	return tk, tk_
}

