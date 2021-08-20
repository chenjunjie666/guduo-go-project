package youku

import (
	"crypto/md5"
	"fmt"
	"guduo/pkg/util"
	"strings"
	"time"
)

var commentAppKey = "100-DDwODVkv"
var commentAppSecret = "6c4aa6af6560efff5df3c16c704b49f1"

func GetCommentApiUrl (vid string) string {
	data := commentDataTmp
	t := time.Now().UnixNano() / 1e6
	//t := 1619264926237
	preSign2Str := fmt.Sprintf("%s&%s&%d", commentAppKey, commentAppSecret, t)
	sign1 := md5.Sum([]byte(preSign2Str))
	sign1Str := fmt.Sprintf("%x", sign1)
	dataJson := fmt.Sprintf(data, t, vid, sign1Str)

	// 第一步加密完成

	token, _ := GetToken()
	if strings.Contains(token, "_") {
		token = strings.Split(token, "_")[0]
	}
	//token := "e89cb3d1411298c1aa0043af3b30ffff"

	t = time.Now().UnixNano() / 1e6
	//t = 1619264926251
	str := fmt.Sprintf("%s&%d&%d&%s", token, t, appKey, dataJson)
	sign2 := md5.Sum([]byte(str))
	sign2Str := fmt.Sprintf("%x", sign2)

	apiUrl := fmt.Sprintf(`https://acs.youku.com/h5/mtop.youku.ycp.comment.mainpage.get/1.0/?jsv=2.6.1&appKey=%d&t=%d&sign=%s&api=mtop.youku.ycp.comment.mainpage.get&type=originaljson&v=1.0&ecode=1&dataType=json&data=%s`,
		appKey,
		t,
		sign2Str,
		util.UrlEncode(dataJson),
		)
	return apiUrl
}



var commentDataTmp =`{"app":"100-DDwODVkv","time":%d,"objectCode":"%s","objectType":1,"sign":"%s"}`
var commentDataTmp2 =`{"app":"100-DDwODVkv","time":%d,"objectCode":"%s","objectType":1}`