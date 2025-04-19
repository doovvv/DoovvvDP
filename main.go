package main

import (
	"doovvvDP/config"
	"doovvvDP/dal/mysql"
	"doovvvDP/dal/redis"
	"doovvvDP/router"
)

func main() {
	config.ConfigInit()
	mysql.Init()
	redis.RedisInit()
	router.RouterInit()
}
