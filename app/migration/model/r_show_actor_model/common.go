package r_show_actor_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/model_scrawler/show_actor_model"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetLoliPopMysqlConn()
	}
	return m.Model(&Table{})
}

func Sync () {
	show_actor_model.Model().Where("id > ?", 0).Delete(nil)
	tn := show_actor_model.Table{}.TableName()
	db.GetCrawlerMysqlConn().Exec("alter table `" + tn + "` AUTO_INCREMENT=1;")

	var res []*Table
	Model().Select("*").Find(&res)


	arr := make([]*show_actor_model.Table, 0, 400)
	for _, v := range res {
		playType := getPlayType(v.Category)
		row := &show_actor_model.Table{
			ShowId:   v.ShowId,
			ActorId:  v.ActorId,
			Name:     v.ActorName,
			Avatar:   v.Avatar,
			Play:     v.Roles,
			PlayType: playType,
		}
		arr = append(arr, row)

		if len(arr) >= 400 {
			show_actor_model.Model().Create(&arr)
			arr = make([]*show_actor_model.Table, 0, 400)
		}
	}

	// 最后还要吧剩余的给塞进数据库
	if len(arr) > 0 {
		show_actor_model.Model().Create(&arr)
	}
}

func getPlayType(c string) int8 {
	switch c {
	case "FLYING_GUEST":
		return show_actor_model.PlayTypeTempGuest
		// 暂定嘉宾
	case "GENERAL":
		return show_actor_model.PlayTypeCame
		// 配角
	case "HERO":
		return show_actor_model.PlayTypeLead
		// 领衔主演
	case "HEROINE":
		return show_actor_model.PlayTypeLead
		// 领衔主演
	case "HOST":
		return show_actor_model.PlayTypeGuest
		// 常驻嘉宾
	case "JUDGES":
		return show_actor_model.PlayTypeGuest
		// 嘉宾
	case "LEAD_ACTOR":
		return show_actor_model.PlayTypeStar
		// 主演
	case "OTHER_LEADING_ACTOR":
		return show_actor_model.PlayTypeOtherLead
		// 其他领衔主演
	case "OTHER_MAJOR_ACTOR":
		return show_actor_model.PlayTypeOtherStar
		// 其他主演
	case "PERMANENT_GUEST":
		return show_actor_model.PlayTypeGuest
		// 常驻嘉宾
	case "VARIETY_PLAYER":
		return show_actor_model.PlayTypeGuest
		// 嘉宾
	}
	return 0
}