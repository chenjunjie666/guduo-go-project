package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var Tencent = &tencent{
	PlatformId: storage.Tencent.PlatformId,
	Host:       storage.Tencent.Host,
	ApiHosts: struct {
		CommentHost     string
		VideoAccessHost string
		VideoTunnelHost string
		VideoMFMHost    string
	}{
		CommentHost:     "https://video.coral.qq.com",
		VideoAccessHost: "https://access.video.qq.com",
		VideoTunnelHost: "https://tunnel.video.qq.com",
		VideoMFMHost:    "https://mfm.video.qq.com",
	},
}

type tencent struct {
	PlatformId uint64
	Host       string
	ApiHosts   struct {
		CommentHost     string
		VideoAccessHost string
		VideoTunnelHost string
		VideoMFMHost    string
	}
}

// 腾讯视频采集器初始化
func (t tencent) Collector(mod string) *core.CollectorObj {
	cInfo := &core.CollectorInfo{
		"腾讯视频",
		t.Host,
		t.PlatformId,
		mod,
	}

	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫

	extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                            // 使用代理，如果代理没有则不使用
	c.DetectCharset = true                  // 非utf-8字符集支持
	return c
}

// 获取评论数接口所需要的callback参数
// cid 为评论id
func (t tencent) GenVarticleCallbackStr(cid string) string {
	cbStr := fmt.Sprintf("_varticle%scommentv2", cid)

	return cbStr

}

func (t tencent) GenJQueryCallbackStr() string {
	// 生成 callback 参数
	ns := time.Now().UnixNano()
	ms := time.Now().UnixNano() / 1e6
	cb := fmt.Sprintf("jQuery_%d_%d", ns, ms)

	return cb
}

func (t tencent) ParseCommentId(vid string) string {

	c := t.Collector("评论ID")

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	commentId := ""
	callback := t.GenJQueryCallbackStr()
	// 经观察以下为常数
	otype := "json"
	op := "3"
	vappid := "30645497"
	vsecret := "d38052bb634963e03eca5ce3aaf93525324d970f110f585f"
	u := fmt.Sprintf("%s/fcgi-bin/video_comment_id?otype=%s&callback=%s&op=%s&vappid=%s&vsecret=%s&vid=%s",
		t.ApiHosts.VideoAccessHost,
		otype,
		callback,
		op,
		vappid,
		vsecret,
		vid,
	)
	c.OnResponse(func(r *colly.Response) {
		ctxStr := string(r.Body)
		ctxStr = strings.TrimLeft(ctxStr, callback+"(")
		ctxStr = strings.TrimRight(ctxStr, ");")
		ctxByte := []byte(ctxStr)

		commentId, _ = jsonparser.GetString(ctxByte, "comment_id")
	})

	_ = c.Visit(u)

	return commentId
}

func (t tencent) ParseCommentIds(vids []string) []string {
	var commentIds []string
	for _, vid := range vids {
		commentId := t.ParseCommentId(vid)
		commentIds = append(commentIds, commentId)
	}
	return commentIds
}

func (t tencent) ParseCoverId(u string) string {

	uArr := strings.Split(u, "/")
	uArrLen := len(uArr)

	coverId := ""
	for i := uArrLen - 1; i >= 0; i-- {
		if uArr[i] == "cover" {
			coverIdTmp := uArr[i+1]
			coverIdTmpArr := strings.Split(coverIdTmp, ".")
			coverId = coverIdTmpArr[0]
			break
		}
	}
	return coverId
}

func (t tencent) ParseVid(u string) string {

	uArr := strings.Split(u, "/")
	uArrLen := len(uArr)

	Vid := ""
	vidTmp := uArr[uArrLen-1]
	VidTmpArr := strings.Split(vidTmp, ".")
	Vid = VidTmpArr[0]
	return Vid
}

func (t tencent) ParseVids(u string) []string {

	vids := make([]string, 0, 10)

	c := t.Collector("获取VIDs")
	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	c.OnResponse(func(r *colly.Response) {
		html := r.Body

		reg := regexp.MustCompile(`COVER_INFO\s*=.*?(\n|\r|\r\n|\n\r)`)

		coverInfoTmp := reg.Find(html)

		index := 0

		for i, ch := range coverInfoTmp {
			if string(ch) == "{" {
				index = i
				break
			}

		}

		_, _ = jsonparser.ArrayEach(coverInfoTmp[index:], func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			vid, _ := jsonparser.GetString(value, "V")
			vids = append(vids, vid)
		}, "nomal_ids")
	})

	_ = c.Visit(u)

	return vids
}

func (t tencent) ParseTargetId(vid string) string {

	type requestData struct {
		WRegistType   int64                        `json:"wRegistType"`
		VecIdList     []string                     `json:"vecIdList"`
		WSpeSource    int64                        `json:"wSpeSource"`
		BIsGetUserCfg int64                        `json:"bIsGetUserCfg"`
		MapExtData    map[string]map[string]string `json:"mapExtData"`
	}

	c := t.Collector("TargetID")
	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	// 经观察以下为常数
	vappid := "97767206"
	vsecret := "c0bdcbae120669fff425d0ef853674614aa659c605a613a4"
	raw := "1"
	var vecIdList []string
	vecIdList = append(vecIdList, vid)
	var mapExtData map[string]map[string]string
	mapExtData = make(map[string]map[string]string)
	var vidMap map[string]string
	vidMap = make(map[string]string)
	vidMap["strCid"] = ""
	vidMap["strLid"] = ""
	mapExtData[vid] = vidMap
	u := fmt.Sprintf("%s/danmu_manage/regist?vappid=%s&vsecret=%s&raw=%s",
		t.ApiHosts.VideoAccessHost,
		vappid,
		vsecret,
		raw,
	)
	targetId := ""
	requestDataTmp := &requestData{
		WRegistType:   2,
		VecIdList:     vecIdList,
		WSpeSource:    0,
		BIsGetUserCfg: 1,
		MapExtData:    mapExtData,
	}
	data, err := json.Marshal(requestDataTmp)
	if err != nil {
		log.Warn(err)
	}
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "text/plain")
		r.Headers.Set("Origin", "https://v.qq.com")
		r.Headers.Set("Refer", "https://v.qq.com/")
	})

	c.OnResponse(func(r *colly.Response) {
		ctx := r.Body
		vidFormatted := fmt.Sprintf("%s", vid)
		strDanMuKey, _ := jsonparser.GetString(ctx, "data", "stMap", vidFormatted, "strDanMuKey")
		reg := regexp.MustCompile(`targetid=\d+`)
		targetIdTmp := reg.FindString(strDanMuKey)
		targetId = strings.TrimLeft(targetIdTmp, `targetid=`)
	})

	_ = c.PostRaw(u, data)
	return targetId
}

func (t tencent) ParseTargetIds(vids []string) [][2]string {
	var targetIds [][2]string
	for _, vid := range vids {
		targetId := t.ParseTargetId(vid)
		targetIds = append(targetIds, [2]string{vid, targetId})
	}
	return targetIds
}


func (t tencent) ParseLength (coverId, tid string) int {
	url := fmt.Sprintf("https://v.qq.com/x/cover/%s/%s.html", coverId, tid)

	c := t.Collector("腾讯视频时长获取")

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	duration := 0
	c.OnResponse(func(r *colly.Response) {
		reg := regexp.MustCompile("VIDEO_INFO =.*?\n")
		tmp := reg.Find(r.Body)
		tmp = bytes.Trim(tmp, "VIDEO_INFO \n")
		durationStr, e := jsonparser.GetString(tmp, "duration")
		if e != nil {
			return
		}
		duration, _ = strconv.Atoi(durationStr)
	})

	c.Visit(url)

	return duration
}