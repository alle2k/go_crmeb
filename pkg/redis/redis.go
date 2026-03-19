package redis

import (
	"boilerplate/config"
	"context"
	"log"
	"strconv"
	"time"

	_redis "github.com/redis/go-redis/v9"
)

var RDB *_redis.Client
var ctx = context.Background()

func Init() {

	RDB = _redis.NewClient(&_redis.Options{
		Addr:         config.AppConfig.Redis.Host + ":" + strconv.Itoa(config.AppConfig.Redis.Port),
		Password:     config.AppConfig.Redis.Password,
		DB:           config.AppConfig.Redis.Database,
		DialTimeout:  time.Duration(config.AppConfig.Redis.Timeout) * time.Millisecond,
		PoolSize:     config.AppConfig.Redis.MaxActive,
		MinIdleConns: config.AppConfig.Redis.MinIdle,
		MaxIdleConns: config.AppConfig.Redis.MaxIdle,
	})

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
}
