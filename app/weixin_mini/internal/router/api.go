package router

import (
	"github.com/gin-gonic/gin"
	"guduo/app/weixin_mini/internal/controller"
	home2 "guduo/app/weixin_mini/internal/controller/actor/home"
	"guduo/app/weixin_mini/internal/controller/cfg"
	"guduo/app/weixin_mini/internal/controller/show/detail"
	"guduo/app/weixin_mini/internal/controller/show/home"
)

func apiRouter(server *gin.Engine) {
	server.GET("/get", controller.GetController)
	server.POST("/post", controller.PostController)

	server.GET("/publish", controller.Publisher)
	server.GET("/config", cfg.Config)
	server.GET("/date_config", cfg.DateConfig)

	server.GET("/show/home/search", home.Search)
	server.GET("/show/home/hot_search", home.HotSearch)
	server.GET("/show/home/list", home.List)
	server.GET("/actor/home/list", home2.List)

	server.GET("/show/detail/base_info", detail.BaseInfo)
	server.GET("/show/detail/net_hot", detail.NetHot)
	server.GET("/show/detail/douban", detail.Douban)
	server.GET("/show/detail/danmaku_comment", detail.DanmakuComment)
	server.GET("/show/detail/word_cloud", detail.WordCloud)
	server.GET("/show/detail/weibo", detail.Weibo)
	server.GET("/show/detail/hot_weibo", detail.WeiboHot)
	server.GET("/show/detail/analysis", detail.Analysis)
}
