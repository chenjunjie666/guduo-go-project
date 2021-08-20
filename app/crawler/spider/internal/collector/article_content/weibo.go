package article_content

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func weiboHandle() {
	ret := storage.Weibo.GetFetchArticleUrl()

	log.Info("微博文章一共", len(ret), "个连接需要爬取")
	for _, row := range ret {
		wg.Add(1)
		ch.PushJob()
		go weiboArticleContent(row.Url, row.ID)
	}

	wg.Done() // 这里把外层的jobNum给Done掉
}

func weiboArticleContent(u string, showId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	c := common.Weibo.Collector(ModName)

	findFlag := false

	c.OnError(func(r *colly.Response, err error) {
		c.Retry(r, err)
	})

	tc := int64(0)
	uids := make([]string, 20)
	c.OnHTML("#pl_feedlist_index", func(e *colly.HTMLElement) {
		e.ForEach(".card-wrap", func(i int, ele *colly.HTMLElement) {
			findFlag = true
			uid := ele.Attr("mid")
			isTop := ele.DOM.Find(".card-top")
			if uid != "" {
				uids = append(uids, uid)
			}

			if isTop.Text() == "" || isTop.Find(".title").Text() != "热门" {
				return
			}

			ctHtml := ele.DOM.Find(".card") // 单条微博的 html
			name := ctHtml.Find(".name").Text() // 发布者名称
			_time := ctHtml.Find(".content > .from > a:nth-child(1)").Text() // 发布时间
			forward := ctHtml.Find(".card-act li:nth-child(2) > a").Text() // 转发数
			content := ctHtml.Find(".txt").Text() // 正文内容

			// 将数据多余的空格等清理，并转换时间，转发数等
			nameFormat := strings.Trim(name, " \n")
			timeFormat := weiboFilterTime(strings.Trim(_time, " \n"))
			forwardFormat := weiboFilterForward(strings.Trim(forward, " \n"))
			contentFormat := strings.Trim(content, " \n")

			tc += storage.Weibo.StoreArticleContent(uid, nameFormat, timeFormat, contentFormat, forwardFormat, JobAt, showId)
		})
	})

	b := ""
	c.OnResponse(func(r *colly.Response) {
		b = string(r.Body)
	})

	_ = c.Visit(u)

	// 如果没有找到，记录错误日志
	if !findFlag {
		log.Warn(fmt.Sprintf("获取相关微博失败，链接：%s, 代理IP:%s", u, c.GetProxyIp()))
	}

	storage.Weibo.StoreArticleNum(tc, JobAt, showId)
}

// todo 热门微博，手机端
func weiboArticleContentMobile() {
	//c := common.Weibo.MobileCollector(ModName)
	//
	//findFlag := false
	//
	//c.OnError(func(r *colly.Response, err error) {
	//	c.Retry(r, err)
	//})
}


func weiboFilterTime(t string) string {
	tFormat := ""
	if strings.Contains(t, "今天") {
		// 时间为 今天
		// 直接翻译成今天的日期
		if tt := strings.Trim(t, "今天 "); tt != "" {
			tFormat = time.Now().Format("2006-01-02") + " " + tt + ":00"
		}else{
			tFormat = time.Now().Format("2006-01-02 15:04:05")
		}
	}else if strings.Contains(t, "年") {
		// 时间为 xxxx年xx月xx日 xx:xx
		// 取其中的年月日
		tmp := strings.Split(t, " ")[0]
		reg := regexp.MustCompile(`(年|月|日|\s)`)
		tmp = reg.ReplaceAllString(tmp, "-")
		tFormat = strings.Trim(tmp, "-") + " 00:00:00"
	}else if strings.Contains(t, "月") {
		// 时间为 xx月xx日 xx:xx
		// 这种人为年为当前年，所以取 月日，在和当前年拼接
		curYear := time.Now().Format("2006")
		tmp := strings.Split(t, " ")[0]
		reg := regexp.MustCompile(`(月|日|\s)`)
		tmp = reg.ReplaceAllString(tmp, "-")
		tFormat = strings.Trim(tmp, "-")
		tFormat = fmt.Sprintf("%s-%s 00:00:00", curYear, tFormat)
	}else if strings.Contains(t, "分钟"){
		// 时间为 xx分钟前
		// 将分钟转换为秒，减去当前时间戳，将得到的时间戳格式化为年月日
		tmp := strings.Split(t, "分钟")[0]
		tmpNum, _ := strconv.ParseInt(tmp, 10, 64)
		sec := tmpNum * 60
		curTimestamp := time.Now().Unix()
		tmpTimestamp := curTimestamp - sec
		tFormat = time.Unix(tmpTimestamp, 0).Format("2006-01-02 15:04:05")
	}// 其他情况，暂时不考虑，还没有遇到

	return tFormat
}

func weiboFilterForward(f string) int64 {
	reg := regexp.MustCompile(`\d+`)
	fNumStr := reg.FindString(f)

	fNum, _ := strconv.ParseInt(fNumStr, 10, 64)
	return fNum
}
