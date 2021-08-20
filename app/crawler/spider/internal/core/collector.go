package core

import (
	"context"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/core/proxy"
	"guduo/app/internal/constant"
	"guduo/pkg/util"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func NewCollector(cInfo *CollectorInfo, options ...colly.CollectorOption) *CollectorObj {
	c := colly.NewCollector(options...)
	cid := util.UniqueID()
	cObj := &CollectorObj{
		cid,
		cInfo,
		options,
		10, // 默认重试次数为10次
		make(map[string]int),
		c,
		false,
		"",
		false,
	}
	return cObj
}

type CollectorInfo struct {
	Name   string // 名称
	Host   string // 主机ID
	PltID  uint64 // 平台ID
	Module string // 指标名
}

// 这个struct继承了 colly.Collector
// 是实际上操作的爬虫实体
type CollectorObj struct {
	Cid              uint64         // 当前爬虫实体的唯一ID
	Info             *CollectorInfo // 爬虫基本信息
	option           []colly.CollectorOption
	maxRetry         int // 最大重试次数
	curRetryCount    map[string]int // 当前重试次数
	*colly.Collector     // colly爬虫实体采集器
	useProxy         bool
	curProxyIp       string // 当前使用的代理IP
	outLock          bool
}

// 重写了 colly 的 clone 方法
func (c *CollectorObj) Clone() *CollectorObj {
	cid := util.UniqueID()
	clone := &CollectorObj{
		cid,
		c.Info,
		c.option,
		c.maxRetry,
		make(map[string]int),
		c.Collector.Clone(),
		c.useProxy,
		c.curProxyIp,
		false,
	}

	return clone
}

func (c *CollectorObj) SetTimeout(sec time.Duration)  {
	// 默认3秒超时
	c.WithTransport(&http.Transport{
		DialContext: (&net.Dialer{
			Timeout: sec * time.Second,
		}).DialContext,
	})
}

// 指定采集器是否使用代理IP
func (c *CollectorObj) UseProxy() {
	c.useProxy = true
}

func (c *CollectorObj) GetProxyIp() string {
	return c.curProxyIp
}


func (c *CollectorObj) setProxy() {
	ip, err := proxy.FromProxyPoolSwitcher()
	// 未获取到代理的情况下，不使用代理
	if err != nil || ip.String() == "" {
		return
	}

	var getProxy = func(r *http.Request) (*url.URL, error) {
		ctx := context.WithValue(r.Context(), colly.ProxyURLKey, ip.String())
		*r = *r.WithContext(ctx)
		return ip, err
	}
	c.curProxyIp = ip.String()

	c.WithTransport(&http.Transport{
		Proxy: getProxy,
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second, // 最大10秒的超时
		}).DialContext,
	})

	// 该方法从代理池中获取一个IP，如果代理池IP为空，则不使用代理
}

func (c *CollectorObj) unsetProxy() {
	c.useProxy = false
	c.curProxyIp = "local IP"
	// 卸载代理
}


// 发起重试，在达到最大重试次数后，或者重试的时候返回错误，则直接返回error，并且停止重试
func (c *CollectorObj) Retry(r *colly.Response, e error) {
	url := r.Request.URL.String()
	//log.Error(fmt.Sprintf("爬虫ID：%d 出错：%s, url:%s, 使用代理IP为：%s", c.ID, e, r.Request.URL.String(), c.curProxyIp))

	if c.useProxy {
		proxy.ReportUselessProxy(c.curProxyIp)
		// 如果请求失败了，就换一个IP继续试
		c.setProxy()
	}

	if c.curRetryCount[url] < c.maxRetry {
		c.curRetryCount[url]++
		_ = c.ReVisit(url)
	} else {
		log.Warning("超出最大重试次数:%s", url)
	}
}

var curCollector = 0 // 当前并发数
var maxPoolNum = 2000 // 最大并发数
var startLock = &sync.Mutex{}
var doneLock = &sync.Mutex{}
var poolChannel = make(chan bool, maxPoolNum)
var outLockTimes = 0 // 程序卡顿后，绕过并发锁的collector数量
var maxOutLockTime = 200 // 最大允许快速生成的绕过并发锁的collector的数量
var outLockTimeWaitingSecond time.Duration = 3000 // 绕过并发锁时需要等待的时间，单位毫秒
var unLockTimeWaitingLimitTimes = 5 // outLockTimes 达到这个数字时，大幅缩短 outLockTimeWaitingSecond 的时间，让程序快速通过死锁阶段
var curWaitingThread = 0


func (c *CollectorObj) addWorker() {
	curWaitingThread++
	startLock.Lock()
	defer startLock.Unlock()
	//if c.Info.PltID == constant.PlatformIdIqiyi || c.Info.PltID == constant.PlatformIdBilibili {
	//	// 爱奇艺和b站最少70ms访问一次连接，他们需要频繁访问文件系统，所以慢一点，否则系统会挂
	//	time.Sleep(time.Millisecond * 80)
	//}else if c.Info.PltID == constant.PlatformIdTencent {
	//	// 腾讯在本地，60ms一次也不会被ban，所以腾讯就等的少一点
	//	// 这里可以以60为基数（本地测试无代理60ms不会被ban）
	//	// 除以代理IP池大小 * 3 （为了防止段时间获取到相同IP，稍微慢一点）
	//	// 代理池大小为 N， 这这里的 ms = 60 / N * 3  当前 N = 25, 则 ms ≈ 7
	//	time.Sleep(time.Millisecond * 7)
	//}else{
	//	// 其他平台不需要走IO所以，直接40ms启动一次collector
	//	time.Sleep(time.Millisecond * 60)
	//}
	if c.Info.PltID == constant.PlatformIdWeibo {
		time.Sleep(time.Millisecond * 200)
	}else if c.Info.PltID == constant.PlatformIdTencent {
		time.Sleep(time.Millisecond * 20)
	}else if c.Info.PltID == constant.PlatformIdBaidu{
		// 百度爬了一段时间莫名就被封了，不知道是账号还是评率问题，先降频试试 2021-05-25 22:09:57
		// 百度最多只有百度指数需要爬越4W个接口，耗时不超过5小时，为了防止被封
		// 这个频率已经降低到了单线程80ms极限的3倍
		time.Sleep(time.Millisecond * 250)
	}else {
		time.Sleep(time.Millisecond * 50)
	}

	select {
	case poolChannel <- true:
		outLockTimes = 0
	case <-time.After(time.Millisecond * outLockTimeWaitingSecond): // 最多等待5秒，就不计入最大count统计，防止死锁
		outLockTimes++
		c.outLock = true
		// 如果绕过并发控制的超过 unLockTimeWaitingLimitTimes，我们为人程序有哪里出现了死锁或者其他意外情况
		// 导致程序只能通过绕过并发锁来生成collector，这种情况，我们需要放开绕过并发锁的限制
		// 让他可以以最快150ms一次的速度生成collector
		if outLockTimes > unLockTimeWaitingLimitTimes {
			outLockTimeWaitingSecond = 100
		}
		
		// 我们最多允许程序以快速形式生成绕过并发锁的collector的数量为 maxOutLockTime
		// 这是为了防止触发了 unLockTimeWaitingLimitTimes 后的快速生成的collector过多，导致内存扛不住
		// 这时候就要把 outLockTimes 重置为0，重新开始计算绕过并发的collector数量，让他慢下来
		if outLockTimes > maxOutLockTime {
			outLockTimes = 0
			outLockTimeWaitingSecond = 3000
		}
	}
	curCollector++
	curWaitingThread--
}

func (c *CollectorObj) Done() {
	doneLock.Lock()
	defer doneLock.Unlock()
	if c.outLock == true {
		return
	}

	<-poolChannel

	if curCollector > 0 {
		curCollector--
	}
}


// retry的时候，有一个300ms的等待对略
func (c *CollectorObj) addRetryWorker() {
	startLock.Lock()
	defer startLock.Unlock()
	time.Sleep(time.Millisecond * 300)
}

type sepPool struct {
	lock *sync.Mutex
	interval time.Duration
}
func (s *sepPool) push() {
	s.lock.Lock()
	defer s.lock.Unlock()
	time.Sleep(time.Millisecond * s.interval)
}

var sepPoolMap = make(map[string]*sepPool) // 特别爬虫池
var createSepLock = &sync.Mutex{}
var delSepLock = &sync.Mutex{}
var sepChannel = make(map[uint64]chan bool)
func (c *CollectorObj) AddWorkVip(id string){
	createSepLock.Lock()
	target, ok := sepPoolMap[id]
	if !ok {
		var interval time.Duration = 110
		chanLen := 55
		switch c.Info.PltID {
		case constant.PlatformIdBilibili:
			interval = 900
			//chanLen = 50
		case constant.PlatformIdIqiyi:
			interval = 300
			//chanLen = 40
		case constant.PlatformIdMango:
			interval = 350
			//chanLen = 100
		case constant.PlatformIdSouhu:
			interval = 2000
			//chanLen = 20
		case constant.PlatformIdTencent:
			// 腾讯弹幕接口和评论接口不是同一个域名
			// 至少目前观察下来同时跑不会导致封IP
			// 所以不做降频处理
			interval = 150
			//chanLen = 100
		case constant.PlatformIdYouku:
			// 单线跑的极限是大约300ms，但是现实是弹幕和评论共享IP封锁容量
			// 所以为了房子弹幕和评论同时启动导致IP被封禁，这里对访问评率做降低处理 300 -> 400
			interval = 400
			//chanLen = 50
		}

		target = &sepPool{
			&sync.Mutex{},
			interval,
		}
		sepPoolMap[id] = target

		if _, ok2 := sepChannel[c.Info.PltID]; !ok2 {
			sepChannel[c.Info.PltID] = make(chan bool, chanLen)
		}
		sepChannel[c.Info.PltID] <- true
	}
	createSepLock.Unlock()
	target.push()
}

// 释放资源
func (c *CollectorObj) ReleaseVip(id string){
	delSepLock.Lock()
	defer delSepLock.Unlock()
	if _, ok := sepPoolMap[id]; ok {
		delete(sepPoolMap, id)
		<- sepChannel[c.Info.PltID]
	}
}

func ReleaseVip(id string, pid uint64){
	if _, ok := sepPoolMap[id]; ok {
		delete(sepPoolMap, id)
		<- sepChannel[pid]
	}
}


func (c *CollectorObj) Visit(url string, sep... string) error {
	sepStr := ""
	if len(sep) == 1 {
		sepStr = sep[0]
	}
	if sepStr != "" {
		c.AddWorkVip(sepStr)
	}else{
		c.addWorker()
	}
	// 防止用一个collector请求多个页面导致retry count次数公用
	// 所以visit视为新的collector
	// 并且重置一些其他必要的参数
	if c.useProxy {
		c.setProxy()
	}else{
		// 默认3秒超时
		c.SetTimeout(10)
	}
	//if sepStr != "" {
	//	log.Info(fmt.Sprintf("特别爬虫任务：%s, 爬虫：%s，正在请求连接：%s", sepStr, c.Info.Module, url))
	//}else{
	//	log.Info(fmt.Sprintf("爬虫：%s，正在请求连接：%s", c.Info.Module, url))
	//}
	if _, ok := c.curRetryCount[url]; !ok {
		// 该url的重试次数
		c.curRetryCount[url] = 0
	}
	e := c.Collector.Visit(url)
	if sepStr == "" {
		c.Done()
	}
	return e
}

func (c *CollectorObj) Post(url string, data map[string]string, sep... string) error {
	sepStr := ""
	if len(sep) == 1 {
		sepStr = sep[0]
	}
	if sepStr != "" {
		c.AddWorkVip(sepStr)
	}else{
		c.addWorker()
	}
	//if sepStr != "" {
	//	log.Info(fmt.Sprintf("特别爬虫任务：%s, 爬虫：%s，正在请求连接：%s", sepStr, c.Info.Module, url))
	//}else{
	//	log.Info(fmt.Sprintf("爬虫：%s，正在请求连接：%s", c.Info.Module, url))
	//}
	if c.useProxy {
		c.setProxy()
	}else{
		c.SetTimeout(10)
	}
	if _, ok := c.curRetryCount[url]; !ok {
		// 该url的重试次数
		c.curRetryCount[url] = 0
	}
	e := c.Collector.Post(url, data)
	if sepStr == "" {
		c.Done()
	}
	return e
}

func (c *CollectorObj) PostRaw(url string, reqData []byte, sep... string) error {
	sepStr := ""
	if len(sep) == 1 {
		sepStr = sep[0]
	}
	if sepStr != "" {
		c.AddWorkVip(sepStr)
	}else{
		c.addWorker()
	}
	if c.useProxy {
		c.setProxy()
	}else{
		c.SetTimeout(10)
	}
	//if sepStr != "" {
	//	log.Info(fmt.Sprintf("特别爬虫任务：%s, 爬虫：%s，正在请求连接：%s", sepStr, c.Info.Module, url))
	//}else{
	//	log.Info(fmt.Sprintf("爬虫：%s，正在请求连接：%s", c.Info.Module, url))
	//}
	if _, ok := c.curRetryCount[url]; !ok {
		// 该url的重试次数
		c.curRetryCount[url] = 0
	}
	e := c.Collector.PostRaw(url, reqData)
	if sepStr == "" {
		c.Done()
	}
	return e
}

func (c *CollectorObj) ReVisit(url string, sep... string) error {
	sepStr := ""
	if len(sep) == 1 {
		sepStr = sep[0]
	}
	if sepStr != "" {
		c.AddWorkVip(sepStr)
	}else{
		c.addRetryWorker()
	}
	if c.useProxy {
		c.setProxy()
	}else{
		c.SetTimeout(10)
	}
	c.AllowURLRevisit = true
	//log.Info(fmt.Sprintf(">>>>>爬虫：%s<<<<，正在尝试重新请求连接：%s", c.Info.Module, url))

	e := c.Collector.Visit(url)
	return e
}
