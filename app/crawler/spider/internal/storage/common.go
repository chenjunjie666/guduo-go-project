package storage

import (
	"context"
	"encoding/json"
	"guduo/app/internal/model_scrawler/actor_model"
	"guduo/app/internal/model_scrawler/show_actor_model"
	"guduo/app/internal/model_scrawler/show_detail_model"
	"guduo/app/internal/model_scrawler/show_model"
	"guduo/app/internal/test"
	"guduo/pkg/model"
	"strings"
)

// 这是一个 storage 内通用的空 context
// 当你不需要 context 发挥任何作用
// 使用此 context 变量
var ctx = context.TODO()



// 基本信息
func getNeedFetchBaseInfoUrl(platformId uint64) show_detail_model.DetailUrls {
	var sIds []*struct{
		ID model.PrimaryKey
	}
	show_model.Model().Select("id").
		Where("status", show_model.ShowStatStandard).
		Where("is_crawler_base_info", 0).
		Where("id IN ?", test.TestShowIds).
		Find(&sIds)

	sidArr := make([]model.PrimaryKey, len(sIds))
	for k, v := range sIds {
		sidArr[k] = v.ID
	}
	urls := show_detail_model.GetDetailUrl(platformId, sidArr)
	return urls
}


// 存储演员，导演，饰演角色信息
func storeBaseInfo(bim map[string]string, showId uint64){
	ckUpdate := show_model.Model().Where("is_crawler_base_info = ?", 0).Limit(1).Find(nil, showId)
	if ckUpdate.RowsAffected == 0{
		return
	}
	ckUpdate.Update("is_crawler_base_info", 1)
	// 保存导演信息
	if bim["Director"] != "" {
		dirs := strings.Split(bim["Director"], ",")
		dirTmp := make([]string, 0, 5)
		for _, dir := range dirs {
			dir = strings.Trim(dir, " ")// 去掉左右空格
			if dir == "" {
				continue
			}
			dirTmp = append(dirTmp, dir)
		}

		// 导演只需要存json数组即可
		dirJson, _ := json.Marshal(dirTmp)
		show_model.Model().First(nil, showId).Update("director", dirJson)
	}

	// 有饰演优先用饰演角色
	if bim["PlayRole"] != "" {
		var acts []*show_actor_model.Table
		var role map[string]string
		_ = json.Unmarshal([]byte(bim["PlayRole"]), &role)
		for actor, play := range role {
			a := strings.Trim(actor, " ")// 去掉左右空格
			if a == ""{
				continue
			}

			// 保存演员
			aid := actor_model.SaveActor(a)
			if aid == 0 {
				continue
			}
			act := &show_actor_model.Table{
				ShowId: showId,
				ActorId: aid,
				Name: actor,
				Play: play,
				PlayType: 1,
			}
			acts = append(acts, act)
		}

		if len(acts) > 0 {
			r := show_actor_model.Model().Where("show_id = ?", showId).Limit(1).Find(nil)
			if r.RowsAffected > 0 {
				return
			}
			// 保存剧集的出演演员信息
			show_actor_model.Model().Create(acts)
		}
	}else if bim["Actor"] != "" || bim["Guest"] != "" {
		var acts []*show_actor_model.Table

		var actors []string
		pt := show_actor_model.PlayTypeStar // play type
		if bim["Actor"] != "" {
			actors = strings.Split(bim["Actor"], ",")
		}else{
			pt = show_actor_model.PlayTypeGuest
			actors = strings.Split(bim["Guest"], ",")
		}

		for _, actor := range actors {
			a := strings.Trim(actor, " ")// 去掉左右空格
			if a == ""{
				continue
			}

			// 保存演员
			aid := actor_model.SaveActor(a)
			if aid == 0 {
				continue
			}
			act := &show_actor_model.Table{
				ShowId: showId,
				ActorId: aid,
				Name: actor,
				Play: "",
				PlayType: pt,
			}
			acts = append(acts, act)
		}

		if len(acts) > 0 {
			r := show_actor_model.Model().Where("show_id = ?", showId).Limit(1).Find(nil)
			if r.RowsAffected > 0 {
				return
			}
			// 保存剧集的出演演员信息
			show_actor_model.Model().Create(acts)
		}
	}
}

// 保存视频简介
func storeIntro(in string, sid uint64){
	var cnt int64
	show_model.Model().Where("id = ?", sid).
		Where("is_crawler_intro = ?", 0).
		Count(&cnt)
	if cnt == 0 {
		return // 已经被爬虫爬取过
	}

	show_model.Model().Where("id = ?", sid).
		Where("is_crawler_intro = ?", 0).
		Update("introduction", in).
		Update("is_crawler_intro", 1)
}

// 保存视频时长
func storeLength(l string, sid uint64){
	var cnt int64
	show_model.Model().Where("id = ?", sid).Where("is_crawler_len = ?", 0).Count(&cnt)
	if cnt == 0 {
		return // 已经被爬虫爬取过
	}

	show_model.Model().Where("id = ?", sid).
		Where("is_crawler_len = ?", 0).
		Update("length", l).
		Update("is_crawler_len", 1)
}