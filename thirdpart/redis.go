package thirdpart

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-counter/config"
)

var redisClient *redis.Client

func InitRedis() error {

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.GetConfig().RedisConf.Host, config.GetConfig().RedisConf.Port), // 本机redis地址
		Password: config.GetConfig().RedisConf.Passwd,                                                        // redis密码，没有则留空
		DB:       0,                                                                                          // 默认数据库，默认是0
	})

	return redisClient.Ping(context.Background()).Err()
}

func HMGetRedisData(key string, fields ...string) []interface{} {

	val := redisClient.HMGet(context.Background(), key, fields...).Val()
	return val
}

func HMSetRedisData(key string, fields ...interface{}) error {

	result, err := redisClient.HMSet(context.Background(), key, fields...).Result()
	if !result {

		return err
	}
	return nil

}
