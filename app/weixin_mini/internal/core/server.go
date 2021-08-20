package core

import "github.com/gin-gonic/gin"

var coreServer *gin.Engine

func initServer() {
	if coreServer == nil {
		coreServer = gin.Default()
	}
}

func GetServer() *gin.Engine {
	if coreServer == nil {
		initServer()
	}
	return coreServer
}