package show

import (
	"github.com/gin-gonic/gin"
	"guduo/app/cms/internal/hepler/resp"
	"guduo/app/cms/internal/services/show"
	"guduo/app/internal/constant"
	"guduo/app/internal/model_clean/tag_model"
	"guduo/app/internal/model_scrawler/show_actor_model"
	"guduo/app/internal/model_scrawler/show_model"
)

// poster 海报
func Add(c *gin.Context) {
	var data show.ShowDetailParams
	c.ShouldBindJSON(&data)

	if data.Name == "" {
		resp.Fail(c, "剧名必填")
		return
	}

	var isRepeat int64
	show_model.Model().Where("name", data.Name).Count(&isRepeat)
	if isRepeat > 0 {
		resp.Fail(c, "剧名重复")
		return
	}

	sid, e := show.AddNewShow(data)

	if e != nil {
		resp.Fail(c, e)
		return
	}

	ret := map[string]uint64{
		"show_id": sid,
	}

	resp.Success(c, ret)
}

func Edit(c *gin.Context) {
	var data show.ShowDetailParams
	c.ShouldBindJSON(&data)

	if data.Name == "" {
		resp.Fail(c, "剧名必填")
		return
	}

	if data.ShowId == 0 {
		resp.Fail(c, "剧ID不能为空")
		return
	}

	var isRepeat int64
	show_model.Model().Where("id != ?", data.ShowId).Where("name", data.Name).Count(&isRepeat)
	if isRepeat > 0 {
		resp.Fail(c, "剧名重复")
		return
	}

	e := show.Update(data)

	if e != nil {
		resp.Fail(c, e)
		return
	}

	resp.Success(c, nil)
}

func Detail(c *gin.Context) {
	type param struct {
		ShowId uint64 `json:"show_id"`
	}
	//
	var data param
	//var data map[string]interface{}
	_ = c.ShouldBindJSON(&data)

	if data.ShowId == 0 {
		resp.Fail(c, "剧ID不能为空")
		return
	}

	sid := data.ShowId

	detail := show.Detail(sid)

	resp.Success(c, detail)
}

func Delete(c *gin.Context) {
	type param struct {
		ShowId uint64 `json:"show_id"`
	}

	var data param
	c.ShouldBindJSON(&data)
	if data.ShowId == 0 {
		resp.Fail(c, "剧ID不能为空")
		return
	}

	sid := data.ShowId

	show_model.Model().Where("id", sid).Delete(nil)

	resp.Success(c, nil)
}

func Config(c *gin.Context) {
	plt := []map[string]interface{}{
		{
			"id": constant.PlatformIdTencent,
			"name": "腾讯",
		},
		{
			"id": constant.PlatformIdYouku,
			"name": "优酷",
		},
		{
			"id": constant.PlatformIdMango,
			"name": "芒果",
		},
		{
			"id": constant.PlatformIdIqiyi,
			"name": "爱奇艺",
		},
		{
			"id": constant.PlatformIdBilibili,
			"name": "bilibili",
		},
		{
			"id": constant.PlatformIdWeibo,
			"name": "微博",
		},
		{
			"id": constant.PlatformIdDouban,
			"name": "豆瓣",
		},
		{
			"id": constant.PlatformIdBaidu,
			"name": "百度",
		},
	}

	pltMicro := []map[string]interface{}{
		{
			"id": constant.PlatformIdTikTalk,
			"name": "抖音",
		},
		{
			"id": constant.PlatformIdKuaishou,
			"name": "快手",
		},
		{
			"id": constant.PlatformIdTxMicro,
			"name": "腾讯微视",
		},
		{
			"id": constant.PlatformIdKuaidianTV,
			"name": "快点TV",
		},
		{
			"id": constant.PlatformIdFanYue,
			"name": "番乐",
		},
	}

	playType := []map[string]interface{}{
		{
			"id": show_actor_model.PlayTypeLead,
			"name": "领衔主演",
		},
		{
			"id": show_actor_model.PlayTypeStar,
			"name": "主演",
		},
		{
			"id": show_actor_model.PlayTypeSupp,
			"name": "配角",
		},
		{
			"id": show_actor_model.PlayTypeCame,
			"name": "客串",
		},
		{
			"id": show_actor_model.PlayTypeGuest,
			"name": "嘉宾",
		},
	}

	subShowType := []map[string]interface{}{
		{
			"id": show_model.ShowTypeSeries,
			"name": "剧集",
			"sub_show_type": []map[string]interface{}{
				{
					"id": show_model.ShowSubTypeSeriesNet,
					"name": "网络剧",
				},
				{
					"id": show_model.ShowSubTypeSeriesTV,
					"name": "电视剧",
				},
			},
		},
		{
			"id": show_model.ShowTypeVariety,
			"name": "综艺",
			"sub_show_type": []map[string]interface{}{
				{
					"id": show_model.ShowSubTypeVarietyNet,
					"name": "网络综艺",
				},
				{
					"id": show_model.ShowSubTypeVarietyTV,
					"name": "电视综艺",
				},
			},
		},
		{
			"id": show_model.ShowTypeMovie,
			"name": "电影",
			"sub_show_type": []map[string]interface{}{
				{
					"id": show_model.ShowSubTypeMovieNet,
					"name": "网络电影",
				},
				{
					"id": show_model.ShowSubTypeMovieCinema,
					"name": "院线电影",
				},
			},
		},
		{
			"id": show_model.ShowTypeAmine,
			"name": "动漫",
			"sub_show_type": []map[string]interface{}{
				{
					"id": show_model.ShowSubTypeAmineChina,
					"name": "国漫",
				},
			},
		},
	}

	tag := tag_model.GetTag()

	cfg := map[string]interface{}{
		"platform": plt,
		"platform_micro": pltMicro,
		"play_type": playType,
		"show_type": subShowType,
		"tag": tag,
	}

	resp.Success(c, cfg)
}
