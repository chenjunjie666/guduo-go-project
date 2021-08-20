package actor_info_model

import "guduo/pkg/model"

type Table struct {
	ActorId            model.PrimaryKey `gorm:"primaryKey"`
	Name               model.Varchar
	ActorSort          model.Int // 重名次数
	OtherName          model.Varchar
	PhoneNumber        model.Varchar
	Gender             model.Tinyint // 0-男  1-女
	Birthday           model.Varchar
	City               model.Varchar
	HomeTown           model.Varchar // 老家
	Height             model.Float
	Weight             model.Float
	School             model.Varchar // 毕业学校
	GraduationDate     model.Varchar
	PrimaryImageUrl    model.Varchar
	HeadImageUrl       model.Varchar
	Talent             model.Varchar //特长
	WechatAccount      model.Varchar
	WeiboAccount       model.Varchar
	Bwh                model.Varchar
	Experience         model.Tinyint  // 是否有表演经验
	Sign               model.Tinyint  // 是否签约
	SelfEvaluation     model.Text     // 自我介绍
	IndustryEvaluation model.Text     // 行业介绍
	Major              model.Varchar  // 专业
	ActingMajor        model.Tinyint  // 是否表演专业
	Awards             model.Text     // 获奖奖项
	Disable            model.Tinyint  // 是否显示
	TagLevel           model.Tinyint  // ???
	BrokerCompany      model.Varchar  // 经纪公司
	BirthYear          model.Int // 出生年
	EnrollmentDate     model.Varchar // 出道时间
}

func (t Table) TableName() string {
	return "actor_info"
}
