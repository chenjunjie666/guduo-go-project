// 这个文件不是入口文件！

package main

import (
	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector()
	c.Async = true

	sliceRes := make([]map[string]string, 0, 10000)
	c.OnResponse(func(r *colly.Response) {
		var res map[string]string
		res = make(map[string]string)
		res["x"], _ = jsonparser.GetString(r.Body, "key1") // pprof show that there has 89% that was used on this line
		res["y"], _ = jsonparser.GetString(r.Body, "key2")
		res["z"], _ = jsonparser.GetString(r.Body, "key3")
		sliceRes = append(sliceRes, res)
	})

	for i:=0; i<100; i++ {
		c.Visit("link...")
	}

	c.Wait()
}

