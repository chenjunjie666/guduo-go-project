package router

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(server *gin.Engine){
	apiRouter(server)
}
