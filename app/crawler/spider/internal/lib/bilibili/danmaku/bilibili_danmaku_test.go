package danmaku

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
)

// 测试用弹幕的 buffer 二进制文件
var path = "./seg.so.example"

func TestLoadDanmakuFile(t *testing.T) {
	f, e := os.Open(path)
	if e != nil {
		fmt.Println(e)
		return
	}

	defer f.Close()

	c, e := ioutil.ReadAll(f)

	//fmt.Println(string(c))

	// 弹幕池
	dm := &DmSegMobileReply{}
	// 解析 bytebuffer 到弹幕池中
	proto.Unmarshal(c, dm)

	for _, row := range dm.Elems {
		// 关键字段：
		// row.Content 正文
		// row.Ctime 弹幕发送时间
		fmt.Println(row)
	}
}
