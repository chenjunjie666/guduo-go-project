package controller

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	time2 "time"
)

type proxyUrls []*proxyUrl
type proxyUrl struct {
	url string
	expired int64 // 连接超时时间 - 已经没用了，用接口判断可用时间
	useCount int64
	errCount int64
}

var orderId = "902174654738615"
var secret = "dmudun2ujccp4eap5uikgkld4ln9vba4"

// 以 [ip]errCount 的形式存储一个代理ip的错误次数信息
var ipMap = make(map[string]int)

// ip池
var pool = make([]*proxyUrl, 0, 100)

func FetchProxy (ctx *gin.Context) {
	if ctx.Query("secret") != "ueqy8qfyrhwov9c8f5ksnkxlrnk8eo33" {
		ctx.String(403, "not allow")
		return
	}

	ip := getProxyUrl()

	ctx.String(200, ip)
}


var lastFetch int64 = 0
func fetch() {
	t := time2.Now().Unix()
	if t - lastFetch < 3 {
		time2.Sleep(time2.Second * 2)
	}
	lastFetch = t

	apiUrl := fmt.Sprintf("http://dps.kdlapi.com/api/getdps?orderid=%s&sign_type=simple&signature=%s" +
		"&timestamp=%d&num=60&dedup=1&format=json&sep=1",
		orderId,
		secret,
		t,
	)
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		code, e := jsonparser.GetInt(r.Body, "code")
		msg, _ := jsonparser.GetString(r.Body, "msg")
		if e != nil || code != 0 {
			log.Error(fmt.Sprintf("获取代理失败，返回code：%d, msg:%s, 错误是：%s", code, msg, e))
			return
		}
		_, _ = jsonparser.ArrayEach(r.Body, func(v []byte, dataType jsonparser.ValueType, offset int, err error) {
			ip := string(v)
			url := fmt.Sprintf("http://%s", ip)

			addProxyUrl(url)
		}, "data", "proxy_list")
	})

	_ = c.Visit(apiUrl)
}

var poolLock = &sync.Mutex{}
func ReportErrorProxy(ctx *gin.Context) {
	//暂不使用错误报告机制，就用代理检测机制即可
	//ip := ctx.Query("ip")
	//if _, ok := ipMap[ip]; ok == true && ip != ""{
	//	poolLock.Lock()
	//	ipMap[ip]++
	//	poolLock.Unlock()
	//}
}

func CheckProxy() {
	for {
		checkProxyValid()
		checkProxyLeftTime()
		checkErrCount()
		time2.Sleep(time2.Second * 5) // 每五秒检测一次IP可用性
	}
}

// 检测代理可用性
func checkProxyValid(){
	if len(pool) == 0 {
		return
	}
	t := time2.Now().Unix()

	urls := make([]string, 0, 100)
	for _, v := range pool {
		u := strings.Trim(v.url, "htps:/")
		urls = append(urls, u)
	}
	c2 := colly.NewCollector()
	url := fmt.Sprintf("https://dps.kdlapi.com/api/checkdpsvalid?orderid=%s&signature=%s&proxy=%s&timestamp=%d",
		orderId,
		secret,
		strings.Join(urls, ","),
		t,
	)

	c2.OnResponse(func(r *colly.Response) {
		jsonparser.ObjectEach(r.Body, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			if string(value) == "false" {
				deleteUrl(string(key))
			}
			return nil
		}, "data")
	})
	c2.Visit(url)
}

// 检测代理剩余时长
func checkProxyLeftTime() {
	if len(pool) == 0 {
		return
	}
	t := time2.Now().Unix()

	urls := make([]string, 0, 100)
	for _, v := range pool {
		u := strings.Trim(v.url, "htps:/")
		urls = append(urls, u)
	}
	c2 := colly.NewCollector()
	url := fmt.Sprintf("https://dps.kdlapi.com/api/getdpsvalidtime?orderid=998752579990468&signature=ueqy8qfyrhwov9c8f5ksnkxlrnk8eo33&proxy=%s&timestamp=%d",
		strings.Join(urls, ","),
		t,
	)

	c2.OnResponse(func(r *colly.Response) {
		jsonparser.ObjectEach(r.Body, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			leftTime, _ := strconv.Atoi(string(value))
			if leftTime <= 8 {
				deleteUrl(string(key))
			}
			return nil
		}, "data")
	})
	c2.Visit(url)
}

func checkErrCount() {
	for _, url := range pool{
		if url.errCount > 5 {
			deleteUrl(url.url)
		}
	}
}

func deleteUrl(u string) {
	poolLock.Lock()
	defer poolLock.Unlock()
	if !strings.Contains(u, "http") {
		u = "http://" + u
	}
	delete(ipMap, u)
	pLen := len(pool)
	for idx, v := range pool {
		if v.url == u {
			tmp := make(proxyUrls, 0, pLen - 1)
			tmp = append(tmp, pool[0: idx]...)
			tmp = append(tmp, pool[idx + 1: pLen]...)
			pool = tmp
			break
		}
	}
}


func getProxyUrl() string {
	poolLock.Lock()
	defer poolLock.Unlock()
	pLen := len(pool)
	rand.Seed(time2.Now().UnixNano()) // 设置随机数种子
	var pu *proxyUrl

	if pLen <= 60 {
		fetch()
		pLen = len(pool)
	}
	if pLen == 0 {
		return "" // 代理池为空就返回空字符串，爬虫如果返回空字符串就会不使用代理
	}
	// 从代理池随即取出一个IP
	idx := rand.Intn(pLen)
	u := pool[idx]
	pu = u

	ip := ""
	if pu != nil {
		ip = pu.url
	}
	return ip
}

// 添加url的方式只有一处，而该处已经在锁里面了，所以这里不要再加锁了，不然就会死锁
func addProxyUrl(url string)  {
	// 检测IP是否已经存在
	if _, ok := ipMap[url]; ok == true {
		return
	}

	ipMap[url] = 0
	row := &proxyUrl{url, 0, 0, 0}
	pool = append(pool, row)
}