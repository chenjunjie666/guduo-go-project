package main

import (
	"guduo/app/migration/internal"
	"guduo/app/migration/model/bi_shows_model"
	"guduo/app/migration/model/carl_shows_model"
)

// 从旧数据库迁移数据至新数据库
func main() {
	internal.InitDB()

	//actor_billboard_model.Sync()
	//barrage_daily_logs_model.Sync()
	//
	//
	//douban_logs_model.Sync()
	//inc_comment_logs_model.Sync()
	//inc_play_count_daily_model.Sync()
	//play_count_total_model.Sync()
	//show_gdi_logs_model.Sync()
	//
	//word_cloud_model.Sync()
	//
	//inc_weibo_model.Sync()


	// show表信息
	//r_show_actor_model.Sync()
	//actor_info_model.Sync()
	bi_shows_model.Sync()
	carl_shows_model.Sync()

	// 不用了
	//play_count_daily_inc_model.Sync()


	//fmt.Println("DONE!!!!!!!!!!!!!!!!!!")
}