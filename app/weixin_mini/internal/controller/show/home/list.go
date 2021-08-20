package home

import (
	"guduo/app/internal/model_scrawler/show_model"
	"guduo/app/weixin_mini/internal/hepler/resp"
	"guduo/app/weixin_mini/internal/services/show"
	"guduo/pkg/util"
	"strconv"

	"github.com/gin-gonic/gin"
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

	typeTmp := c.Query("type")
	if typeTmp == "" {
		resp.Fail(c, "非法请求")
		return
	}
	type_, _ := strconv.Atoi(typeTmp)

	subTypeTmp := c.Query("sub_type")
	if subTypeTmp == "" {
		resp.Fail(c, "非法请求")
		return
	}
	subType, _ := strconv.Atoi(subTypeTmp)

	pidTmp := c.Query("platform_id")
	if pidTmp == "" {
		resp.Fail(c, "非法请求")
		return
	}
	pid, _ := strconv.Atoi(pidTmp)

	var ret []map[string]interface{}
	if type_ == int(show_model.ShowTypeMicro) {
		ret = show.GetMicroRank(dayAt, pid)
	}else {
		ret = show.List(dayAt, listType, type_, subType, pid)
	}

	resp.Success(c, ret)
}


func Search(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		resp.Fail(c, "请输入关键字")
		return
	}
	keyword = util.UrlDecode(keyword)
	ret := show.SearchShow(keyword)

	resp.Success(c, ret)
}


func HotSearch(c *gin.Context) {
	typeTmp := c.Query("type")
	if typeTmp == "" {
		resp.Fail(c, "非法请求")
		return
	}
	type_, _ := strconv.Atoi(typeTmp)

	ret := show.HotSearch(type_)

	resp.Success(c, ret)
}