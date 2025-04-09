package v1_test

import (
	v1 "doovvvDP/api/v1"
	"doovvvDP/config"
	"doovvvDP/dal/mysql"
	"doovvvDP/dal/redis"
	"testing"
)
func TestSaveShop2Redis(t *testing.T){
	config.ConfigInit()
	mysql.Init()
	redis.RedisInit()
	v1.SaveShop2Redis(1,10)
}