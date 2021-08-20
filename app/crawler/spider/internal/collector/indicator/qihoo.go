package indicator

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"guduo/pkg/util"
	"image"
	"image/png"
	"regexp"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func qihooHandle() {
	urls := storage.Qihoo.GetIndicatorUrl()

	for _, row := range urls {
		wg.Add(1)
		ch.PushJob()
		go qihooIndicator(row.Url, row.ShowId)
	}
	wg.Done()
}

// 爬取360指数
func qihooIndicator(u string, showId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	apiHost := common.Qihoo.ApiHosts.IndexHost
	words := common.Qihoo.ParseWords(u)
	startDate := common.Qihoo.ParseStartDate()
	endDate := common.Qihoo.ParseEndDate()
	// 经过观察以下参数是固定值
	area := util.UrlEncode("全国")
	//click := "2"
	//t := "index"

	// 格式化获取搜狗指数的链接
	u = fmt.Sprintf("%s/index/csssprite?q=%s&area=%sD&from=%s&to=%s&click=8&t=index",
		apiHost,
		words,
		area,
		startDate,
		endDate,

	)

	var requestData map[string]string
	requestData = make(map[string]string)
	//requestData["q"] = words
	//requestData["area"] = area
	//requestData["from"] = startDate
	//requestData["to"] = endDate
	//requestData["click"] = click
	//requestData["t"] = t
	c := common.Qihoo.Collector(ModName)

	c.OnError(func(r *colly.Response, err error) {
		c.Retry(r, err)
	})

	findFlag := false

	c.OnResponse(func(r *colly.Response) {
		findFlag = true
		ctxByte := r.Body
		keywords := util.UrlDecode(words) // urldecode
		cssOffsetData, _ := jsonparser.GetString(ctxByte, "data", keywords, "css")
		CssOffsetList := filterCssOffset(cssOffsetData)
		imgData, _ := jsonparser.GetString(ctxByte, "data", keywords, "img")

		if imgData == "" {
			log.Info("showid：", showId, "的360指数, 未解析到，源数据为：", string(ctxByte))
			return
		}

		imgStr := filterImgStr(imgData)
		imgList := generateImage(imgStr, CssOffsetList)
		lastDayDataInt := imgNumRecognition(imgList)
		log.Info("获取到showid：", showId, "的360指数，值为：", lastDayDataInt)
		storage.Qihoo.StoreIndicator(lastDayDataInt, JobAt, showId)
	})
	_ = c.Post(u, requestData)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到360指数", u))
	}
}

func filterCssOffset(s string) []int {
	reg := regexp.MustCompile(`background-position:-\d+\.000000px`)
	result := reg.FindAllString(s, -1)
	var offsetList []int
	for _, numStr := range result {
		numStr = strings.TrimLeft(numStr, "background-position")
		numStr = strings.TrimLeft(numStr, ":")
		numStr = strings.TrimRight(numStr, "000000px")
		numStr = strings.TrimRight(numStr, ".")
		num, _ := strconv.Atoi(numStr)
		offsetList = append(offsetList, num)
	}
	return offsetList
}

func filterImgStr(s string) string {
	imgStr := strings.TrimLeft(s, "data:image/png;base64")
	imgStr = strings.TrimLeft(imgStr, ",")
	return imgStr
}

func generateImage(s string, cssOffsetList []int) []image.Image {
	decodeBytes, _ := base64.StdEncoding.DecodeString(s)
	imgStreams := bytes.NewBuffer(decodeBytes)
	img, e := png.Decode(imgStreams)
	if e != nil {
		return nil
	}
	rgbImg := img.(*image.NRGBA)
	var imgList []image.Image
	for _, cssOffset := range cssOffsetList {
		cssOffset = -cssOffset
		subImg := rgbImg.SubImage(image.Rect(cssOffset, 0, cssOffset+6, 12)).(*image.NRGBA)
		imgList = append(imgList, subImg)
	}
	return imgList
}

func imgNumRecognition(imgList []image.Image) int64 {
	s := ""

	MD5ToNumStrMap :=
		map[string]string{
			"b356f39cdc9304e710a606f84772fa68": "0",
			"dbee25923ba34906bcd58d5ae296b2ec": "1",
			"f5c170293a2e31d3b766e6fd90fa7e4c": "2",
			"f28bdeaf80f1ed7ebd6988329db17320": "3",
			"629255230724f1dea68763429ce4366c": "4",
			"065a2b744fc4d8c6fa51dfb6c8013a35": "5",
			"c9ccbd52ffe1342da6c9311f13b15877": "6",
			"94ed6315455f4207c2e5f612c518704b": "7",
			"e9a58151002458bb886b1b75c1133434": "8",
			"a7bdc356788cdc0459fb55d2f46e6727": "9",
		}
	for _, subImg := range imgList {
		emptyBuff := bytes.NewBuffer(nil)
		png.Encode(emptyBuff, subImg)
		dist := base64.StdEncoding.EncodeToString(emptyBuff.Bytes())
		m := md5.Sum([]byte(dist))
		idx := fmt.Sprintf("%x", m)
		s += MD5ToNumStrMap[idx]
	}
	num, _ := strconv.ParseInt(s, 10, 64)
	return num
}
