package router

import "guduo/app/weixin_mini/internal/core"

func InitRouter(){
	server := core.GetServer()

	apiRouter(server)
}
