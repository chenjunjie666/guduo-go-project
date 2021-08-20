package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"guduo/app/cms/internal/boot"
	"guduo/pkg/log"
	"net/http"
	"net/http/httptest"
	"strings"
)

var s *gin.Engine
func init()  {
	log.InitLogger()
	boot.InitDB()
	s = boot.InitServer()
}

func GetRequest(uri string, data map[string]string) string {
	queryStr := ""

	for key, val := range data {
		queryStr += fmt.Sprintf("%s=%s&", key, val)
	}

	queryStr = strings.Trim(queryStr, "&")
	queryStr = "?" + queryStr

	uri = uri + queryStr
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", uri, nil)

	s.ServeHTTP(w, req)
	return w.Body.String()
}

func PostRequest(uri string, data map[string]string) string {
	queryStr := ""

	for key, val := range data {
		queryStr += fmt.Sprintf("%s=%s&", key, val)
	}

	queryStr = strings.Trim(queryStr, "&")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer([]byte(queryStr)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	s.ServeHTTP(w, req)
	return w.Body.String()
}

func PrintRespJson(s string) {
	var str bytes.Buffer

	_ = json.Indent(&str, []byte(s), "", "\t")

	fmt.Println("+++++++++++++++++++++ response json area +++++++++++++++++++++")
	fmt.Println(str.String())
	fmt.Println("+++++++++++++++++++++ response json area +++++++++++++++++++++")
}
