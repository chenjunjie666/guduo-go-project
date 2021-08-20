package config

import (
	"guduo/app/crawler/clean/task/daily"
	"guduo/app/crawler/clean/task/free"
	"guduo/app/crawler/clean/task/half_month"
	"guduo/app/crawler/clean/task/half_week"
)

type JobFunc func()

var JobMap = map[string]JobFunc{
	"daily": daily.Run,
	"weekly": half_week.Run,
	"monthly": half_month.Run,
	//"yearly": year.Run,
	"guduo_hot": free.GuduoHotHandle,
}
