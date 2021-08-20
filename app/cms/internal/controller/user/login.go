package user

import (
	"github.com/gin-gonic/gin"
	"guduo/app/cms/internal/hepler/resp"
	"guduo/app/cms/internal/hepler/secure"
	"guduo/app/internal/model_clean/admin_model"
)

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var data loginParams

	c.ShouldBindJSON(&data)

	if data.Username == "" || data.Password == "" {
		resp.Fail(c, "用户名密码不能为空")
		return
	}

	e := admin_model.CheckUser(data.Username, data.Password)
	if e != nil {
		resp.Fail(c, e)
		return
	}

	t, token := secure.GetLoginToken(data.Username)

	ret := map[string]interface{}{
		"token": token,
		"timestamp": t,
		"username": data.Username,
	}

	resp.Success(c, ret)
}