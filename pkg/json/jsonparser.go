package json

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/buger/jsonparser"
	jsoniter "github.com/json-iterator/go"
)

func JsonParserGetArrayLength(j []byte, keys ...string) (int64, error) {
	jt, o, _, err := jsonparser.Get(j, keys...)
	if err != nil {
		return 0, err
	}

	if o != jsonparser.Array {
		return 0, errors.New("looks not like an array")
	}

	var arr []interface{}
	err = jsoniter.Unmarshal(jt, &arr)
	if err != nil {
		return 0, err
	}

	l := len(arr)
	return int64(l), nil
}

//获取指定字段的值(参数1：字段路径，参数2：原始json数据)
func GetFieldFromJson(pathStr string, value []byte) string {
	path := strings.Split(pathStr, "/")
	var temp jsoniter.Any
	for i, v := range path {
		if i == 0 {
			temp = jsoniter.Get(value, v)
			if temp == nil {
				return ""
			}
		} else {
			temp = temp.Get(v)
			if temp == nil {
				return ""
			}
		}
	}

	switch temp.ValueType() {
	case jsoniter.InvalidValue, jsoniter.NilValue, jsoniter.BoolValue, jsoniter.ArrayValue, jsoniter.ObjectValue:
		return ""
	case jsoniter.StringValue:
		return temp.ToString()
	case jsoniter.NumberValue:
		return strconv.Itoa(temp.ToInt())
	}
	return ""
}

// Convert map json string
func ConvertToJsonStr(in interface{}) (string, error) {
	jsonByte, err := jsoniter.Marshal(in)
	if err != nil {
		log.Warn(fmt.Sprintf("Marshal with error: %+v\n", err))
		return "", nil
	}
	return string(jsonByte), nil
}
