package router

import (
	"github.com/gin-gonic/gin"
	"guduo/app/cms/internal/controller"
	"guduo/app/cms/internal/controller/micro"
	"guduo/app/cms/internal/controller/show"
	"guduo/app/cms/internal/controller/upload"
	"guduo/app/cms/internal/controller/user"
	"guduo/app/cms/internal/middleware"
)

func apiRouter(server *gin.Engine) {
	server.Use(middleware.Cors())

	server.GET("/get", controller.GetController)
	server.POST("/post", controller.PostController)
	server.POST("/login", user.Login)

	//server.Use(middleware.Auth) // 登录认证

	server.POST("/upload", func(c *gin.Context) {
		//r, e := c.FormFile("file")
		//fmt.Println(r, e)
		upload.Poster(c)
	})

	server.POST("/show/list", show.List)
	server.POST("/show/detail", show.Detail)
	server.POST("/show/detail/add", show.Add)
	server.POST("/show/detail/edit", show.Edit)
	server.POST("/show/detail/delete", show.Delete)
	server.POST("/show/detail/config", show.Config)


	server.POST("/micro/list", micro.List)
	server.POST("/micro/list/save", micro.Save)
	server.POST("/micro/list/config", show.Config)
	//authGroup.POST("/config", home.List)
}
