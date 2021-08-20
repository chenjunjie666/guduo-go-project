package show_actor_model

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}

func GetShowActors(sid uint64) ShowActors {
	var actors ShowActors

	r := Model().Where("show_id", sid).Find(&actors)

	if r.Error != nil {
		log.Warn(fmt.Sprintf("查询剧集演员失败，show_id:%d，错误原因:%s", sid, r.Error))
		return make(ShowActors, 0)
	}

	return actors
}

// 根据饰演角色的类型（主演，领衔主演等）
func GetShowIdsByType(aid []uint64, pt int8, sids []uint64) []*Table {
	var res []*Table
	Model().Select("show_id", "actor_id").
		Where("actor_id IN ?", aid).
		Where("play_type", pt).
		Where("show_id IN ?", sids).
		Find(&res)

	return res
}