package youku

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"guduo/pkg/util"
	"strings"
	"time"
)

func GetEpisodeUrl(vid, sid string, page int) string {
	tk_, _ := GetToken()
	token := strings.Split(tk_, "_")[0]

	params := epParamsTmp
	nextSession := epNextSessionTmp
	data := epDataTmp
	systemInfo := epSystemInfoTmp

	// 参数填充
	params["videoId"] = vid
	params["showId"] = sid

	start := (page - 1) * 30 + 1
	end := start + 29
	nextSession["itemStartStage"] = start
	nextSession["itemEndStage"] = end
	// 参数填充 END

	nextSessionJson, _ := json.Marshal(nextSession)
	params["nextSession"] = string(nextSessionJson)
	paramsJson, _ := json.Marshal(params)
	systemInfoJson, _ := json.Marshal(systemInfo)

	data["params"] = string(paramsJson)
	data["system_info"] = string(systemInfoJson)
	dataJson, _ := json.Marshal(data)

	t := int(time.Now().Unix()*1000)
	str := fmt.Sprintf("%s&%d&%d&%s", token, t, appKey, string(dataJson))
	sign := md5.Sum([]byte(str))

	apiUrl := fmt.Sprintf(
		"https://acs.youku.com/h5/mtop.youku.columbus.gateway.new.execute/1.0/?jsv=2.6.1&appKey=%d&t=%d&sign=%x&api=mtop.youku.columbus.gateway.new.execute&type=originaljson&v=1.0&ecode=1&dataType=json&data=%s",
		appKey,
		t,
		sign,
		util.UrlEncode(string(dataJson)),
	)

	return apiUrl
}


var epDataTmp = map[string]interface{}{
	"ms_codes": "2019030100",
	"params": "",
	"system_info": "",
}

var epParamsTmp = map[string]interface{}{
	"biz":true,
	"scene":"component",
	"componentVersion":"3",
	"ip":"",
	"debug":0,
	"utdid":"pcweb",
	"userId":"",
	"platform":"pc",
	"nextSession": "",
	"videoId":"",
	"showId":"",
}

var epNextSessionTmp = map[string]interface{}{
	"componentIndex":"3",
	"componentId":"61518",
	"level":"2",
	"itemPageNo": "0",
	"lastItemIndex":"0",
	"pageKey":"LOGICSHOW_LOGICTV_DEFAULT",
	"group":"0",
	"itemStartStage":1,
	"itemEndStage":30,
}

var epSystemInfoTmp = map[string]string{
	"os":"pc",
	"device":"pc",
	"ver":"1.0.0",
	"appPackageKey":"pcweb",
	"appPackageId":"pcweb",
}