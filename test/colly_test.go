package test

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"testing"
)

func TestCollyResponseCookie(t *testing.T) {
	url := "https://v.youku.com/v_show/id_XMjY3MTQ2MDE0OA==.html"

	c := colly.NewCollector()

	c.OnResponseHeaders(func(r *colly.Response) {
		fmt.Println(r.Headers)
	})

	c.Visit(url)

}
