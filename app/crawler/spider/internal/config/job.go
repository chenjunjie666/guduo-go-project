package config

import (
	"guduo/app/crawler/spider/internal/collector/article_content"
	"guduo/app/crawler/spider/internal/collector/article_num"
	"guduo/app/crawler/spider/internal/collector/article_num_actor"
	"guduo/app/crawler/spider/internal/collector/attention"
	"guduo/app/crawler/spider/internal/collector/attention_actor"
	"guduo/app/crawler/spider/internal/collector/base_info"
	"guduo/app/crawler/spider/internal/collector/comment_count"
	"guduo/app/crawler/spider/internal/collector/danmaku"
	"guduo/app/crawler/spider/internal/collector/fans"
	"guduo/app/crawler/spider/internal/collector/hot"
	"guduo/app/crawler/spider/internal/collector/indicator"
	"guduo/app/crawler/spider/internal/collector/indicator_actor"
	"guduo/app/crawler/spider/internal/collector/introduction"
	"guduo/app/crawler/spider/internal/collector/length"
	"guduo/app/crawler/spider/internal/collector/news_num"
	"guduo/app/crawler/spider/internal/collector/play_count"
	"guduo/app/crawler/spider/internal/collector/rating_num"
	"guduo/app/crawler/spider/internal/collector/release_time"
	"guduo/app/crawler/spider/internal/collector/short_comment"
)

type JobFunc func()

var JobMap = map[string]JobFunc{
	// 文章内容
	"article_content": article_content.Run,

	// 文章数
	"article_num": article_num.Run,

	// 文章数-艺人
	"article_num_actor": article_num_actor.Run,

	// 关注度
	"attention": attention.Run,

	// 关注度-艺人
	"attention_actor": attention_actor.Run,

	// 主创&公司&演员
	"base_info": base_info.Run,

	//评论数
	"comment_count": comment_count.Run,

	// 弹幕数
	"danmaku": danmaku.Run,

	// 粉丝数
	"fans": fans.Run,

	// 热度
	"hot": hot.Run,

	// 指标
	"indicator": indicator.Run,
	//
	// 指标-艺人
	"indicator_actor": indicator_actor.Run,

	// 剧情简介
	"introduction": introduction.Run,

	// 视频时长
	"length": length.Run,

	// 新闻数
	"news_num": news_num.Run,

	// 播放量
	"play_count": play_count.Run,

	// 评分
	"rating_num": rating_num.Run,

	// 上线时间
	"release_time": release_time.Run,

	//短评数
	"short_comment": short_comment.Run,
}
