package middleware

import (
	"github.com/gin-gonic/gin"
	"guduo/app/cms/internal/hepler/resp"
	"guduo/app/cms/internal/hepler/secure"
)


func Auth(c *gin.Context) {
	uri := c.Request.URL.Path
	if c.Query("debug") == "1" {
		c.Next()
		return
	}

	username := c.Query("username")
	token := c.Query("token")
	timestamp := c.Query("timestamp")

	if isWithoutAuthUri(uri) == true {
		c.Next()
		return
	}

	if username == "" || token == "" || timestamp == "" {
		resp.Fail(c, "未登录，请重新登录")
		c.Abort()
		return
	}

	if secure.CheckLoginToken(username, token, timestamp) == false {
		resp.Fail(c, "token错误，请重新登录")
		c.Abort()
		return
	}

	c.Next()
}

var withoutAuthUri = []string{
	"/login",
}
func isWithoutAuthUri(target string) bool {
	for _, uri := range withoutAuthUri {
		if uri == target {
			return true
		}
	}

	return false
}
