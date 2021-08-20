package home

import (
	"github.com/gin-gonic/gin"
	"guduo/app/weixin_mini/internal/hepler/resp"
	"guduo/app/weixin_mini/internal/services/actor"
	"strconv"
)

func List(c *gin.Context) {
	dayAt_ := c.Query("day_at")
	if dayAt_ == "" {
		resp.Fail(c, "非法请求")
		return
	}
	dayAt, _ := strconv.Atoi(dayAt_)

	listType_ := c.Query("list_type")
	if listType_ == "" {
		resp.Fail(c, "非法请求")
		return
	}
	listType, _ := strconv.Atoi(listType_)

	rankType_ := c.Query("rank_type")
	if rankType_ == "" {
		resp.Fail(c, "非法请求")
		return
	}
	rankType, _ := strconv.Atoi(rankType_)

	playType := c.Query("play_type")
	if playType == "" {
		resp.Fail(c, "非法请求")
		return
	}


	ret := actor.List(dayAt, listType, rankType, playType)

	resp.Success(c, ret)
}