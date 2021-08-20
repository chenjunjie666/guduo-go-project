package proxy

import (
	"fmt"
	"guduo/app/internal/boot"
	"guduo/pkg/util"
	"io/ioutil"
	"net/http"
	"sync"
)

var Pool *proxyPool

type proxyUrls []*proxyUrl

func InitProxyPool() {
	if Pool == nil {
		Pool = newProxyPool()
	}
}

// 初始化代理池
func newProxyPool() *proxyPool {
	return &proxyPool{
		make(proxyUrls, 0, 20),
		make(proxyUrls, 0, 20),
		make(proxyUrls, 0, 20),
		20,
		false,
	}
}

type proxyPool struct {
	proxies proxyUrls
	using   proxyUrls
	// abandon只存储当前运行中被废弃的代理地址，他应该定期的被后端存储存入数据库
	// 以避免abandoned过大而导致内存爆掉
	abandoned       proxyUrls
	supplementLimit int         // 代理池内最少的IP数
	supplementFlag  bool        // 是否允许补充IP
}

// 补充代理池内的IP
func (p *proxyPool) supplementIp() {
	urls := fetchProxyUrls()
	p.proxies = append(p.proxies, urls...)
}

// 获取代理需要上锁
var proxyLock = &sync.Mutex{}
// 获取单个IP，这已经足够支撑当前业务
// 目前业务中单个 collector 只需要一个代理IP就行
// 具体切换规则由扩展决定
func (p *proxyPool) GetProxy() string {
	proxyLock.Lock()
	defer proxyLock.Unlock()

	resp, err := http.Get(fmt.Sprintf("http://%s:%s/get_proxy?secret=%s", boot.Cfg.Proxy.Host, boot.Cfg.Proxy.Port, boot.Cfg.Proxy.Secret))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	respByte, _ := ioutil.ReadAll(resp.Body)

	data := string(respByte)
	return data
	//pLen := len(p.proxies)
	//rand.Seed(time.Now().Unix()) // 设置随机数种子
	//var pu *proxyUrl
	//for {
	//	if pLen <= 50 {
	//		p.supplementIp()
	//		pLen = len(p.proxies)
	//	}
	//	if pLen == 0 {
	//		break
	//	}
	//	idx := rand.Intn(pLen)
	//	u := p.proxies[idx]
	//	// 看这个代理还能不能用
	//	if u.isUsable() {
	//		u.addCount()
	//		pu = u
	//		break
	//	}
	//
	//	tmp := make(proxyUrls, 0, pLen - 1)
	//	tmp = append(tmp, p.proxies[0: idx]...)
	//	tmp = append(tmp, p.proxies[idx + 1: pLen]...)
	//	p.proxies = tmp
	//	pLen--
	//}
	//if pu == nil {
	//	return ""
	//}
	//return pu.String()
}


func ReportUselessProxy(ip string) {
	http.Get(fmt.Sprintf("http://%s:%s/report_useless?ip=%s", boot.Cfg.Proxy.Host, boot.Cfg.Proxy.Port, util.UrlEncode(ip)))
}