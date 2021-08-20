package controller

import (
	"github.com/gin-gonic/gin"
	"guduo/app/weixin_mini/internal/hepler/resp"
	"guduo/app/weixin_mini/internal/services"
)

func Publisher(c *gin.Context) {
	action := c.Query("action")
	if action == "" {
		resp.Fail(c, "参数错误")
		return
	}

	show, content := services.GetPublisher(action)

	ret := map[string]interface{}{
		"is_show": show,
		"content": content,
	}

	resp.Success(c, ret)
	return
}