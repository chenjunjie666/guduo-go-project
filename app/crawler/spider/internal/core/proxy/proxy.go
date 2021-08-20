package proxy

import (
	"net/url"
)

func FromProxyPoolSwitcher() (*url.URL, error) {
	ip := Pool.GetProxy()
	//ip = "http://31.192.107.162:3128"
	parsedU, err := url.Parse(ip)

	return parsedU, err
}