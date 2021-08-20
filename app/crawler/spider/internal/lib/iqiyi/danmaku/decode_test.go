package danmaku

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestDecode(t *testing.T) {
	f := "test.z"

	/*
		https://cmts.iqiyi.com/bullet/[54/00/7973227714515400]_60_2_5f3b2e24.br
		取出方括号内的这段作为  参数1
		分段为               参数2

		爱奇艺每5分钟一个分段
		也就是说一个20分钟的视频 参数2 最大值为4
		未测试如果超过4的情况下回返回什么，可以去Decode文件，打印 x.Code看看结果，已知成功返回： A0000  （0的数量可能有差异，我记不清了）

		https://cmts.iqiyi.com/bullet/参数1_300_参数2.z
	*/

	b, _ := ioutil.ReadFile(f)
	res := Decode(b)

	// 结果为数据，其中每一个值为一条弹幕的结构体，弹幕结构体包含两个字段 Ctime 和 Content
	// res[0].Ctime    弹幕创建时间，int类型，为19位纳秒级时间戳，所以要 Ctime / 1e9 以获得秒级时间戳
	// res[0].Content  弹幕内容，string类型
	fmt.Println(res)
}
