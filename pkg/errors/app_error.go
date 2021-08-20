package errors

import (
	"errors"
	"fmt"
)

func CmsError(s string) error {
	return errors.New(s)
}

// 应用内产生的错误，都调用此方法
func AppError(mod, s string) error {
	ae := &appError{
		mod,
		s,
	}
	return ae
}

type appError struct {
	module string
	es     string
}

func (e appError) Error() string {
	msg := fmt.Sprintf(`程序模块"%s"发生异常：%s`, e.module, e.es)
	return msg
}
