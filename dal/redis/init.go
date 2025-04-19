package redis

import (
	"context"
	"fmt"
	"strconv"

	"doovvvDP/config"

	"github.com/redis/go-redis/v9"
)

var (
	RDB  *redis.Client
	RCtx = context.Background()
)

func RedisInit() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.MyConfig.RedisConfig.Host + ":" + strconv.Itoa(config.MyConfig.RedisConfig.Port),
		Password: config.MyConfig.RedisConfig.Password,
		DB:       config.MyConfig.RedisConfig.DB,
	})

	if RDB != nil {
		fmt.Println("redis connect success")
	}
}
