package util

import "net/http"

// 输入一个 http.cookie 类型，返回用于设置
// Header 的 cookie 字符串
func BuildCookie(cookie []http.Cookie) string {
	var c string

	for _, ck := range cookie {
		tmpCookie := ck.Name + "=" + ck.Value + ";"
		c += tmpCookie
	}

	return c
}
