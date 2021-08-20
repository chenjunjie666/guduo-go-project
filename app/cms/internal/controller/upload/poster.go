package upload

import (
	"github.com/gin-gonic/gin"
	"guduo/app/cms/internal/hepler/resp"
	"guduo/app/cms/internal/hepler/upload"
	"io"
	"strconv"
	"strings"
	"time"
)

func Poster(c *gin.Context) {
	r, e := c.FormFile("file")
	if e != nil {
		resp.Fail(c, "文件不存在")
	}

	filename := r.Filename
	extArr := strings.Split(filename, ".")
	ext := extArr[len(extArr) - 1]
	f, e := r.Open()
	if e != nil {
		resp.Fail(c, "读取上传文件失败")
	}

	res, e := io.ReadAll(f)
	if e != nil {
		resp.Fail(c, "读取文件内容失败")
	}

	t := time.Now().UnixNano()
	filename = strconv.FormatInt(t, 10) + "." + ext
	uri, e :=upload.Image(res, filename)

	if e != nil {
		resp.Fail(c, "上传文件失败")
	}
	resp.Success(c, map[string]string{"url": uri})
}
