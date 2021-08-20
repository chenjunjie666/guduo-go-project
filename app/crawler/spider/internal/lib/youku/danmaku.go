package youku

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)



/*
优酷弹幕分页是每
page 弹幕分页，从0开始，step为1分钟
 */
func GetDanmakuUrl(vid string, page int) (string, map[string]string) {
	token, _ := GetToken()
	if strings.Contains(token, "_"){
		token = strings.Split(token, "_")[0]
	}

	api := "mopen.youku.danmu.list"
	postData := map[string]interface{}{
		"pid":0,
		"ctype":10004, // int
		"sver":"3.1.0",
		"cver":"v1.0",
		"ctime": int(time.Now().UnixNano() / 1e6), // 13位int
		"guid":"dUsHGVWmT2ICAXTtJ0bXd4wn",
		"vid":vid,
		"mat":page, // int,这是分页
		"mcount":1, // int
		"type":1, // int
	}


	pdataJson, _ := json.Marshal(postData)

	msg := base64.StdEncoding.EncodeToString(pdataJson)
	sign1 := md5.Sum([]byte(msg + key))

	pdata := string(pdataJson)
	pdata = strings.TrimRight(pdata, "}")
	pdata = fmt.Sprintf(`%s,"msg":"%s","sign":"%x"}`, pdata, msg, sign1)

	t := int(time.Now().Unix()*1000)
	str := fmt.Sprintf("%s&%d&%d&%s", token, t, appKey, pdata)
	sign2 := md5.Sum([]byte(str))
	sign2Str := fmt.Sprintf("%x", sign2)

	params := fmt.Sprintf("?jsv=2.5.1&appKey=%d&api=%s&v=1.0&type=originaljson&dataType=jsonp&timeout=20000&jsonpIncPrefix=utility&t=%d&sign=%s",
	appKey,
	api,
	t,
	sign2Str,
	)

	url := "https://acs.youku.com/h5/mopen.youku.danmu.list/1.0/" + params

	return url, map[string]string{"data": pdata}
}