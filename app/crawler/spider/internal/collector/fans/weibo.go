package fans

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func weiboHandle() {
	detailUrls := storage.Weibo.GetDetailUrl()

	for _, row := range detailUrls {
		wg.Add(1)
		ch.PushJob()
		go weiboFans(row.Url, row.ShowId)
	}
	wg.Done()
}

// 采集微博粉丝数主逻辑
func weiboFans(u string, showId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	c := common.Weibo.Collector(ModName)

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	findFlag := false
	// 粉丝数解析方法
	c.OnResponse(func(r *colly.Response) {
		htmlByte := r.Body
		// 正则匹配获取到粉丝数所在的文本段落
		reg := regexp.MustCompile(`W_f16\\">\d+<\\/strong><span class=\\"S_txt2\\">粉丝`)
		s := reg.Find(htmlByte)

		// 进一步找到其中的纯数字段落，也就是粉丝数
		reg = regexp.MustCompile(`>\d+<`)
		s = reg.Find(s)

		// 将正则找到的byte类型转为string后去掉左右的多余字符，再转换为int64类型
		fansCountTmp := string(s)
		fansCountStr := strings.Trim(fansCountTmp, "><")
		fansCount, err := strconv.ParseInt(fansCountStr, 10, 64)
		if err != nil {
			log.Warn(fmt.Sprintf("检测到粉丝数特征，但是提取粉丝数失败，失败原因：%s, 找到的特征为：%s", err, string(s)))
			return
		}

		findFlag = true
		// 存储微博粉丝数
		storage.Weibo.StoreFansCount(fansCount, JobAt, showId)
	})

	_ = c.Visit(u)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到粉丝数", u))
	}
}
