package actor_info_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/model_scrawler/actor_model"
	"guduo/app/migration/model"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetLoliPopMysqlConn()
	}
	return m.Model(&Table{})
}

func Sync() {
	actor_model.Model().Where("id > ?", 0).Delete(nil)
	tn := actor_model.Table{}.TableName()
	db.GetCrawlerMysqlConn().Exec("alter table `" + tn + "` AUTO_INCREMENT=1;")

	var res []*Table
	Model().Select("*").Find(&res)

	arr := make([]*actor_model.Table, 0, 300)
	arrId := make([]uint64, 0, 300)
	for _, v := range res {
		row := &actor_model.Table{
			Name: v.Name,
			ActorSort: v.ActorSort,
			OtherName: v.OtherName,
			PhoneNumber: v.PhoneNumber,
			Gender: v.Gender,
			Birthday: model.StrToTime(v.Birthday),
			City: v.City,
			HomeTown: v.HomeTown,
			Height: v.Height,
			Weight: v.Weight,
			School: v.School,
			GraduationDate: model.StrToTime(v.GraduationDate),
			PrimaryImageUrl: v.PrimaryImageUrl,
			HeadImageUrl: v.HeadImageUrl,
			Talent: v.Talent,
			WechatAccount: v.WechatAccount,
			WeiboAccount: v.WeiboAccount,
			Bwh: v.Bwh,
			Experience: v.Experience,
			Sign: v.Sign,
			SelfEvaluation: v.SelfEvaluation,
			IndustryEvaluation: v.IndustryEvaluation,
			Major: v.Major,
			ActingMajor: v.ActingMajor,
			Awards: v.Awards,
			Disable: v.Disable,
			TagLevel: v.TagLevel,
			BrokerCompany: v.BrokerCompany,
			BirthYear: v.BirthYear,
			EnrollmentDate: model.StrToTime(v.EnrollmentDate),
		}
		row.ID = v.ActorId

		arr = append(arr, row)
		arrId = append(arrId, row.ID)

		if len(arr) >= 300 {
			actor_model.Model().Create(&arr)
			arr = make([]*actor_model.Table, 0, 300)
			arrId = make([]uint64, 0, 300)
		}
	}

	// 最后还要吧剩余的给塞进数据库
	if len(arr) > 0 {
		actor_model.Model().Where("id IN ?", arrId).Delete(nil)
		actor_model.Model().Create(&arr)
	}
}