package secure

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
)

func GetLoginToken(u string) (int64, string) {
	t := time.Now().Unix()
	token := genLoginToken(u, strconv.FormatInt(t, 10))
	return t, token
}

func CheckLoginToken(u string, token string, t string) bool {
	tk := genLoginToken(u, t)

	if tk == token {
		return true
	}

	return false
}


func genLoginToken(s1 string, s2 string) string {
	s := fmt.Sprintf("%s%s%s", s1, s2, "qwertyujngtyujn")
	res := md5.Sum([]byte(s))
	token := fmt.Sprintf("%x", res)

	return token
}