package util

import "net/url"

// url encode 方法
func UrlEncode(s string) string {
	encode := url.QueryEscape(s)
	return encode
}

// url decode 方法，在无法 decode 时会返回原本的字符串
func UrlDecode(s string) string {
	decode, err := url.QueryUnescape(s)
	if err != nil {
		return s
	}
	return decode
}
