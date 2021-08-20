package danmaku

import (
	"bytes"
	"compress/zlib"
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
)

type allXml struct {
	XMLName xml.Name `xml:"danmu"`
	Code    string   `xml:"code"`
	Data    struct {
		XMLName xml.Name `xml:"data"`
		Entry   []struct {
			XMLName xml.Name `xml:"entry"`
			Int     int      `xml:"int"`
			List    danmakuList `xml:"list"`
		} `xml:"entry"`
	} `xml:"data"`
}

type danmakuList []struct {
	XMLName    xml.Name `xml:"list"`
	BulletInfo Danmaku `xml:"bulletInfo"`
}

type Danmaku struct {
	XMLName xml.Name `xml:"bulletInfo"`
	Ctime   int64    `xml:"contentId"`
	Content string   `xml:"content"`
}

func Decode(b []byte) []Danmaku {
	nr := bytes.NewReader(b)
	var out bytes.Buffer
	r, _ := zlib.NewReader(nr)
	defer r.Close() // 释放内存
	if r == nil {
		return make([]Danmaku, 0)
	}
	io.Copy(&out, r)

	resByte := out.Bytes()
	var x *allXml
	xml.Unmarshal(resByte, &x)


	dmkList := make([]Danmaku, 0, 100)

	if x == nil {
		log.Warn(fmt.Sprintf("未解析到弹幕内容，原文长度为：%d", len(b)))
		return dmkList
	}

	for _, entry := range x.Data.Entry {
		for _, dmRow := range entry.List{
			dmkList = append(dmkList, dmRow.BulletInfo)
		}
	}

	return dmkList
}
