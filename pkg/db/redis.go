package db

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// redis 链接，这是一个链接池的顶级对象
// 使用该对象处理redis操作，会自动从 redis pool 中获取链接
var rdb *redis.Client

func InitRedisConn(opt *RedisConfig) error {
	if rdb != nil {
		return nil
	}
	_rdb, err := NewRedisConn(opt)
	rdb = _rdb
	return err
}

// 新建一个redis链接池
// 链接池默认大小为 10/每个CPU
func NewRedisConn(opt *RedisConfig) (*redis.Client, error) {
	o := opt.toRedisOption()
	r := redis.NewClient(o)
	ctx := context.Background()
	if _, ok := r.Ping(ctx).Result(); ok != nil {
		return nil, ok
	}
	return r, nil
}

func GetRedisConn() *redis.Client {
	return rdb
}

type RedisConfig struct {
	Addr string
	Pass string
	Port int
	DB   int
}

func (r RedisConfig) toRedisOption() *redis.Options {
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.Addr, r.Port),
		Password: r.Pass,
		DB:       r.DB,
	}
}
