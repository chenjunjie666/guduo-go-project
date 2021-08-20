package proxy

import (
	"net/url"
	"time"
)

func newProxyUrl(u string, expired int64) (*proxyUrl, error) {
	parsedU, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	return &proxyUrl{parsedU, expired, 0}, nil

}

type proxyUrl struct {
	*url.URL
	expired int64 // 连接超时时间
	useCount int
}

func (p *proxyUrl) isUsable() bool {
	now := time.Now().Unix()
	// 过期时间-5秒抵御程序执行误差
	if p.useCount > 100 || (p.expired - 5) < now {
		return false
	}

	return true
}

func (p *proxyUrl) addCount() {
	p.useCount++
}
