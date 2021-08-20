package micro

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"guduo/app/cms/internal/hepler/resp"
	"guduo/app/cms/internal/services/micro"
	"guduo/app/internal/model_clean/micro_drama_rank_tmp_model"
)

type listParams struct {
	DayAt int `json:"day_at"`
}

func List(c *gin.Context) {
	var ReqData listParams
	_ = c.ShouldBindJSON(&ReqData)


	ret := micro.List(ReqData.DayAt)
	resp.Success(c, ret)
}

type saveParams struct {
	DayAt uint `json:"day_at"`
	Data []*saveRow `json:"data"`
}
type saveRow struct {
	Name string `json:"name"`
	Platform []uint64 `json:"platform"`
	Num float64 `json:"num"`
}


func Save(c *gin.Context) {
	var ReqData saveParams
	_ = c.ShouldBindJSON(&ReqData)

	save := make([]*micro_drama_rank_tmp_model.Table, len(ReqData.Data))
	for k, row := range ReqData.Data {
		pids, _ := json.Marshal(row.Platform)
		save[k] = &micro_drama_rank_tmp_model.Table{
			Name:       row.Name,
			PlatformId: string(pids),
			Num:        row.Num,
			DayAt:      ReqData.DayAt,
		}
	}

	micro.SaveRank(save, ReqData.DayAt)
	resp.Success(c, nil)
}
