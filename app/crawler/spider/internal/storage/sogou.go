package storage

import (
	"fmt"
	"guduo/app/internal/constant"
	indicator_model "guduo/app/internal/model_scrawler/indicator_model"
	show_model "guduo/app/internal/model_scrawler/show_model"
	"guduo/pkg/model"
	"guduo/pkg/util"
)

var Sogou = &sogou{
	PlatformId: constant.PlatformIdSogou,
	Host:       "https://www.sogou.com/",
}

type sogou struct {
	PlatformId uint64
	Host       string
}

type SogouIndicatorUrls struct {
	ShowId model.PrimaryKey
	Url    model.Varchar
}

// 获取搜狗指数页面链接
func (s *sogou) GetIndicatorUrl() []*SogouIndicatorUrls {
	shows := show_model.GetActiveShowsName()
	ret := make([]*SogouIndicatorUrls, len(shows))

	for k, show := range shows {
		name := util.UrlEncode(show.Name)
		url := fmt.Sprintf("http://index.sogou.com/index/searchHeat?kwdNamesStr=%s&timePeriodType=MONTH&dataType=SEARCH_ALL&queryType=INPUT", name)
		ret[k] = &SogouIndicatorUrls{
			ShowId: show.Id,
			Url:    url,
		}
	}

	return ret
}

// TODO
// 存储搜狗指数信息
func (s *sogou) StoreIndicator(i int64, jobAt uint, showId uint64) {
	indicator_model.SaveIndicator(float64(i), jobAt, showId, s.PlatformId)
}
