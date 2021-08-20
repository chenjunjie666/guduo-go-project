package daily

import (
	"github.com/yanyiwu/gojieba"
	"guduo/app/internal/model_clean/danmaku_word_cloud_daily_model"
	"guduo/app/internal/model_scrawler/danmaku_model"
	"guduo/app/internal/model_scrawler/show_model"
	"image/color"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)


var DefaultColors = []color.RGBA{
	{0x1b, 0x1b, 0x1b, 0xff},
	{0x48, 0x48, 0x4B, 0xff},
	{0x59, 0x3a, 0xee, 0xff},
	{0x65, 0xCD, 0xFA, 0xff},
	{0x70, 0xD6, 0xBF, 0xff},
}

func guduoWordCloudHandle() {
	showIds := show_model.GetActiveShows()

	for _, id := range showIds {
		guduoWordCloud(id)
	}
}

// 计算骨朵弹幕词云的词语与词频
func guduoWordCloud(sid uint64) {
	data := danmaku_model.GetDanmaku(sid)

	var weightMap map[string]int
	weightMap = make(map[string]int)

	logText := ""
	for _, sentence := range data {
		logText += "。" + sentence
		//fmt.Println(k)
		//dictDir := path.Join(filepath.Dir(os.Args[0]), "dict")
		//jiebaPath := path.Join(dictDir, "jieba.dict.utf8")
		//hmmPath := path.Join(dictDir, "hmm_model.utf8")
		//userPath := path.Join(dictDir, "user.dict.utf8")
		//idfPath := path.Join(dictDir, "idf.utf8")
		//stopPath := path.Join(dictDir, "stop_words.utf8")
		//
		//var seg = gojieba.NewJieba(jiebaPath, hmmPath, userPath, idfPath, stopPath)
	}

	dictDir := path.Join(filepath.Dir(os.Args[0]), "dict")
	jiebaPath := path.Join(dictDir, "jieba.dict.utf8")
	hmmPath := path.Join(dictDir, "hmm_model.utf8")
	userPath := path.Join(dictDir, "user.dict.utf8")
	idfPath := path.Join(dictDir, "idf.utf8")
	stopPath := path.Join(dictDir, "stop_words.utf8")

	var seg = gojieba.NewJieba(jiebaPath, hmmPath, userPath, idfPath, stopPath)
	//var seg = gojieba.NewJieba()

	logText = strings.Replace(logText, " ", "", -1)
	var resWords []string
	resWords = seg.Cut(logText, true)
	seg.Free()
	for _, resWord := range resWords {
		word := strings.Trim(resWord, " ")
		words := strings.Split(word, "")
		repeat := make(map[string]bool)
		for _, v := range words {
			repeat[v] = true
		}
		// 全部为重复字符，如哈哈哈，哈哈，666等
		if len(repeat) == 1 {
			continue
		}
		if utf8.RuneCountInString(word) > 1{
			weightMap[word] = weightMap[word] + 1
		}
	}


	// 取TOP50词频的词做弹幕
	weightMap = filterTop50CloudWords(weightMap)

	danmaku_word_cloud_daily_model.SaveWordCloud(weightMap, JobAt, sid)


	//colors := make([]color.Color, 0)
	//for _, c := range DefaultColors {
	//	colors = append(colors, c)
	//}
	//
	//w := wordclouds.NewWordcloud(
	//	weightMap,
	//	wordclouds.FontFile("/Users/zhaokun/Documents/Golang/src/guduo/build/black.ttf"),
	//	wordclouds.FontMaxSize(20),
	//	wordclouds.FontMinSize(10),
	//	wordclouds.Colors(colors),
	//	wordclouds.Height(500),
	//	wordclouds.Width(700),
	//)
	//
	//img := w.Draw()
	//
	//emptyBuff := bytes.NewBuffer(nil)
	//jpeg.Encode(emptyBuff, img, nil)
	//dist := base64.StdEncoding.EncodeToString(emptyBuff.Bytes())
	//picB64 := dist
	//weightJsonStr, _ := guduoJson.ConvertToJsonStr(weightMap)

}


func filterTop50CloudWords(ws map[string]int) map[string]int {
	if len(ws) <= 50 {
		return ws
	}

	cMap := make([]int, len(ws))
	wMap := make([]string, len(ws))

	i := 0
	for w, c := range ws {
		cMap[i] = c
		wMap[i] = w
		i++
	}

	resCMap, resWMap := top50(cMap, wMap, 50)

	res := make(map[string]int)

	for k := range resCMap {
		res[resWMap[k]] = resCMap[k]
	}
	return res
}

func top50(arr []int, wMap []string, nc int) ([]int, []string) {
	maxIdx := len(arr)
	if maxIdx == 1 {
		return arr, wMap
	}

	rand.Seed(time.Now().Unix())
	idx := rand.Intn(maxIdx)
	cnt := arr[idx]

	left := make([]int, 0, 50)
	leftMap := make([]string, 0, 50)
	mid := make([]int, 0, 50)
	midMap := make([]string, 0, 50)
	right := make([]int, 0, 50)
	rightMap := make([]string, 0, 50)
	for i := 0; i < maxIdx; i++ {
		if arr[i] > cnt {
			left = append(left, arr[i])
			leftMap = append(leftMap, wMap[i])
		}else if arr[i] == cnt {
			mid = append(mid, arr[i])
			midMap = append(midMap, wMap[i])
		}else {
			right = append(right, arr[i])
			rightMap = append(rightMap, wMap[i])
		}
	}

	lLen := len(left)
	mLen := len(mid)

	if lLen < nc && mLen + lLen >= nc {
		x := nc - lLen
		mTmp := mid[0:x]
		mMapTmp := midMap[0:x]

		left = append(left, mTmp...)
		leftMap = append(leftMap, mMapTmp...)
		return left, leftMap
	}else if lLen < nc && mLen + lLen < nc {
		nextCn := nc - (lLen + mLen)
		rTmp, rMapTmp := top50(right, rightMap, nextCn)

		left = append(left, mid...)
		leftMap = append(leftMap, midMap...)

		left = append(left, rTmp...)
		leftMap = append(leftMap, rMapTmp...)

		return left, leftMap
	}else if lLen == nc {
		return left, leftMap
	}else {
		return top50(left, leftMap, nc)
	}
}
