package show

import (
	"github.com/gin-gonic/gin"
	"guduo/app/cms/internal/hepler/resp"
	"guduo/app/cms/internal/services/show"
)

type listParams struct {
	Status int `json:"status"`
	Keyword string `json:"keyword"`
	ShowType int `json:"show_type"`
	Page int `json:"page"`
	Limit int `json:"limit"`
}

func List(c *gin.Context) {
	var ReqData listParams
	_ = c.ShouldBindJSON(&ReqData)

	if stat := c.PostForm("status"); stat == "" {
		ReqData.Status = -99
	}

	res, total := show.List(ReqData.ShowType, ReqData.Status, ReqData.Keyword, ReqData.Page, ReqData.Limit)

	ret := map[string]interface{}{
		"total": total,
		"data": res,
	}

	resp.Success(c, ret)
}