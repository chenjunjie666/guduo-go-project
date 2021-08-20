package request

import (
	"github.com/gin-gonic/gin"
	"guduo/app/weixin_mini/internal/hepler/resp"
	"strconv"
)

func GetShowId(c *gin.Context) uint64 {
	sid_ := c.Query("show_id")
	if sid_ == "" {
		resp.Fail(c, "非法请求")
	}

	sid, _ := strconv.ParseInt(sid_, 10, 64)

	return uint64(sid)
}
