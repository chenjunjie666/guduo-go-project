package services

import "guduo/app/internal/model_clean/admin_publish_model"

func GetPublisher(action string) (int8, string){
	var res *admin_publish_model.Table

	admin_publish_model.Model().Where("type", action).
		Find(&res)

	if res.Content == "" {
		res.IsShow = 0
	}

	return res.IsShow, res.Content
}
