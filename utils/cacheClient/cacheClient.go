package cacheClient

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"doovvvDP/dal/redis"
	"doovvvDP/utils"
	"doovvvDP/utils/redisData"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Set(key string, value any, timeCount time.Duration) {
	valueJson, err := json.Marshal(value)
	if err != nil {
		return
	}
	redis.RDB.Set(redis.RCtx, key, valueJson, timeCount)
}

func SetWithLogicalExpire(key string, value any, timeCount time.Duration) {
	redisData := &redisData.RedisData{
		Data:       value,
		ExpireTime: time.Now().Unix() + int64(timeCount.Seconds()),
	}
	valueJson, err := json.Marshal(redisData)
	if err != nil {
		return
	}
	redis.RDB.Set(redis.RCtx, key, valueJson, -1)
}

func QueryWithPassThrough[T any, ID any](keyPrefix string, id ID,
	dbFallback func(ID) (T, error), timeCount time.Duration,
) (T, error) {
	key := keyPrefix + fmt.Sprintf("%v", id)
	valueJson, err := redis.RDB.Get(redis.RCtx, key).Result()
	if err == nil && valueJson != "" {
		var res T
		json.Unmarshal([]byte(valueJson), &res)
		return res, nil
	}
	if err != goredis.Nil && valueJson == "" {
		var nullValue T
		return nullValue, errors.New("no this data")
	}

	r, err := dbFallback(id)
	if err != nil {
		// 将空值存入redis
		if err == gorm.ErrRecordNotFound {
			redis.RDB.Set(redis.RCtx, key, "", utils.CACHE_NULL_TTL)
		}
		return r, err
	}
	Set(key, r, timeCount)
	return r, nil
}

func getDataFromRedisData[model any](redisData redisData.RedisData, data model) {
	dataJson, err := json.Marshal(redisData.Data)
	if err != nil {
		return
	}
	err = json.Unmarshal(dataJson, data)
	if err != nil {
		return
	}
	return
}

func QueryWithLogicalExpire[T any, ID any](keyPrefix string, id ID,
	dbFallback func(ID) (T, error), timeCount time.Duration,
) (T, error) {
	var nullValue T
	var resValue T
	key := keyPrefix + fmt.Sprintf("%v", id)
	cacheJson, err := redis.RDB.Get(redis.RCtx, key).Result()
	if err == goredis.Nil || cacheJson == "" { // 需要提前预热
		return nullValue, errors.New("商铺信息不存在！")
	}

	// 从缓存中取出shop和过期时间
	var redisData redisData.RedisData
	err = json.Unmarshal([]byte(cacheJson), &redisData)
	if err != nil {
		return nullValue, err
	}
	getDataFromRedisData(redisData, &resValue)

	expireTime := redisData.ExpireTime
	if time.Now().Unix() <= expireTime { // 未过期
		return resValue, nil
	}
	// 逻辑过期，需要重建
	// 1.获取锁
	mutexKey := utils.CACHE_SHOP_MUTEX_KEY + fmt.Sprintf("%v", id)
	ok := tryLock(mutexKey)
	if ok {
		// doublecheck
		cacheJson, _ = redis.RDB.Get(redis.RCtx, key).Result()
		json.Unmarshal([]byte(cacheJson), &redisData)
		getDataFromRedisData(redisData, &resValue)
		expireTime = redisData.ExpireTime
		if time.Now().Unix() <= expireTime {
			return resValue, nil
		}

		// 用一个线程去缓存重建
		go func() {
			r, err := dbFallback(id)
			if err != nil {
				return
			}
			SetWithLogicalExpire(key, r, timeCount)
			unLock(mutexKey)
		}()
	}
	return resValue, nil
}

func tryLock(key string) bool {
	flag, err := redis.RDB.SetNX(redis.RCtx, key, "1", 10*time.Second).Result()
	if err != nil {
		fmt.Println("")
		return false
	}
	return flag
}

func unLock(key string) {
	_, err := redis.RDB.Del(redis.RCtx, key).Result()
	if err != nil {
		fmt.Println("")
	}
}
