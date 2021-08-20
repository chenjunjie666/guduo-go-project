package admin_model

import (
	"crypto/md5"
	"fmt"
	"gorm.io/gorm"
	"guduo/pkg/db"
	"guduo/pkg/errors"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCleanMysqlConn()
	}
	return m.Model(&Table{})
}

// 检查当前艺人用户名密码
func CheckUser(u, p string) error {
	var ad Table
	r := Model().Where("username", u).Find(&ad)
	if r.RowsAffected == 0 {
		return errors.CmsError("用户不存在")
	}

	pmd5Byte := md5.Sum([]byte(p))
	pmd5 := fmt.Sprintf("%x", pmd5Byte)
	if pmd5 != ad.Password {
		return errors.CmsError("用户名密码错误")
	}
	return nil
}