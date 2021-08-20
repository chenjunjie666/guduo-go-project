package resp

import (
	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	Send(c, 0, data, "成功")
}

func Fail(c *gin.Context, msg interface{})  {
	s := ""
	if v, ok := msg.(error); ok {
		s = v.Error()
	}else if v, ok := msg.(string); ok {
		s = v
	}

	if s == "" {
		s = "失败"
	}
	Send(c, 1, nil, s)
}

func Send(c *gin.Context, code int, data interface{}, msg string)  {
	response := buildResponse(code, data, msg)
	c.JSON(200, response)
}

func buildResponse(code int, data interface{}, msg string) map[string]interface{} {
	ret := map[string]interface{}{
		"code": code,
		"data": data,
		"msg":  msg,
	}

	return ret
}