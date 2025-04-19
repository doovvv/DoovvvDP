package redislock

import (
	"fmt"
	"os"
	"time"

	"doovvvDP/config"
	"doovvvDP/dal/redis"
)

type redislock struct {
	name string
}

var luaUnlockScript string

func init() {
	luaScriptPath := "resources/unLock.lua"
	scriptBytes, err := os.ReadFile(luaScriptPath)
	if err != nil {
		panic(err)
	}
	luaUnlockScript = string(scriptBytes)
}

func NewRedisLock(name string) *redislock {
	return &redislock{
		name: name,
	}
}

func (mutex *redislock) TryLock(timeout time.Duration) bool {
	value := fmt.Sprintf("%d-%d", config.MyConfig.MainConfig.WorkerId, os.Getpid())

	success, err := redis.RDB.SetNX(redis.RCtx, mutex.name, value, timeout).Result()
	if err != nil {
		return false
	}
	return success
}

func (mutex *redislock) Unlock() {
	// 如果锁到期，可能会发生误删其他进程的锁
	// 所以需要先判断锁是否存在，再进行删除
	redis.RDB.Eval(redis.RCtx, luaUnlockScript,
		[]string{mutex.name}, []string{fmt.Sprintf("%d-%d", config.MyConfig.MainConfig.WorkerId, os.Getpid())}).Result()
}
