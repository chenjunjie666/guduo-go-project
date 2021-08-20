package controller

import (
	"github.com/gin-gonic/gin"
	"guduo/app/cms/internal/hepler/resp"
)

func GetController(c *gin.Context){
	getParams := c.Query("id") // 获取get参数
	resp.Success(c, getParams)
}

func PostController(c *gin.Context){
	postParams := c.PostForm("id")

	c.JSON(200, postParams)
}
