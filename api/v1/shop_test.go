package v1_test

import (
	"testing"

	v1 "doovvvDP/api/v1"
	"doovvvDP/config"
	"doovvvDP/dal/mysql"
)

func TestSaveShop2Redis(t *testing.T) {
	config.ConfigInit()
	mysql.Init()
	// redis.init()
	v1.SaveShop2Redis(1, 10)
}
