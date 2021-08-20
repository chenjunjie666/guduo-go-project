package detail

import (
	"github.com/gin-gonic/gin"
	"guduo/app/weixin_mini/internal/hepler/resp"
	"guduo/app/weixin_mini/internal/services/show"
	"strconv"
)

func BaseInfo(c *gin.Context) {
	sid_ := c.Query("show_id")
	if sid_ == "" {
		resp.Fail(c, "非法请求")
		return
	}

	sid, _ := strconv.ParseInt(sid_, 10, 64)
	if sid <= 0 {
		resp.Fail(c, "剧综")
		return
	}
	bInfo := show.GetBaseInfo(uint64(sid))

	resp.Success(c, bInfo)
}
