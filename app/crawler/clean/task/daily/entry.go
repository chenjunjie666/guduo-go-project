package daily

import (
	log "github.com/sirupsen/logrus"
	"guduo/pkg/time"
	"sync"
	time2 "time"
)
const ModName = "每日指标"

var JobAt = time.Today()

var Today = JobAt

var Yesterday = Today - 86400

var beforeYesterday = Yesterday - 86400

var yearStart = uint(time2.Date(time2.Now().Year(), 1, 1, 0, 0, 0, 0, time2.Local).Unix())

var wg = &sync.WaitGroup{}

func Run() {
	//jobNum := 1
	//wg.Add(jobNum)

	// 骨朵剧集热度
	//go free.guduoHotHandle()

	log.Info("开始计算昨日骨朵热度排名升降情况")
	//昨日骨朵热度排行以及升降情况
	GuduoHotDailyRankHandle()
	log.Info("计算昨日骨朵热度排名升降情况结束")

	//wg.Wait()

	log.Info("开始计算平均骨朵热度")
	AvgGuduoHotHandle() // 计算平均骨朵热度
	log.Info("计算平均骨朵热度结束")

	log.Info("开始计算骨朵艺人热度/霸屏指数")
	// 骨朵艺人热度/霸屏指数
	GuduoActorHotHandle()
	GuduoActorDomiHandle()
	log.Info("计算骨朵艺人热度/霸屏指数结束")


	log.Info("开始计算总播放量")
	// 总榜播放量
	PlayCountTotalHandle()
	log.Info("总播放量计算结束")


	log.Info("开始计算年播放量")
	// 年榜播放量
	YearTotalPlayCountHandle()
	log.Info("计算年播放量结束")

	log.Info("开始计算电影昨日播放量排名")
	// 电影昨日播放量排名
	moviePlayCountHandle()
	log.Info("计算电影昨日播放量排名结束")



	log.Info("开始生成词云")
	// 生成词云
	guduoWordCloudHandle()
	log.Info("生成词云结束")

}