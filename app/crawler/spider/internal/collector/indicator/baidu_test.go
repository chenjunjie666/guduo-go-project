package indicator

import (
	"fmt"
	"guduo/app/crawler/spider/internal/core"
	"testing"
)

func TestBaiduIndicator(t *testing.T) {
	core.Init()

	wg.Add(3)
	ch.PushJob()
	ch.PushJob()
	ch.PushJob()

	name := "壮志高飞"
	indiUrl := fmt.Sprintf("https://index.baidu.com/v2/main/index.html#/trend/%s?words=%s", name, name)
	//genAgeUrl := fmt.Sprintf("https://index.baidu.com/v2/main/index.html#/crowd/%s?words=%s", name, name)


	//baiduGenderIndicator(indiUrl, 0)
	//baiduAgeIndicator(genAgeUrl, 0)
	baiduIndicator(indiUrl, 0)

	//wg.Wait()
}
