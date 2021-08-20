package util

// 反转字符串
func ReserveString(s string) string {
	ret := ""
	for _, w := range []rune(s) {
		defer func(r rune) { ret += string(r) }(w)
	}
	return ret
}

// 反转 []byte 类型字符串
func ReserveByte(b []byte) []byte {
	var ret []byte
	for _, w := range []rune(string(b)) {
		defer func(r rune) { ret = append(ret, []byte(string(r))...) }(w)
	}
	return ret
}
