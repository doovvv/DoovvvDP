package v1

import (
	"doovvvDP/dal/model"
	"doovvvDP/dal/mysql"
	"doovvvDP/dal/redis"
	"doovvvDP/dto"
	"doovvvDP/utils"
	"doovvvDP/utils/cacheClient"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)
func QueryShopById(id uint64) (result *dto.Result){
	result = &dto.Result{}

	// shop,err := cacheClient.QueryWithPassThrough(utils.CACHE_SHOP_KEY,id,model.GetShopById,utils.CACHE_SHOP_TTL)

	//互斥锁解决缓存击穿
	// shop,err := QueryShopWithMutex(id)

	//逻辑过期解决缓存击穿
	shop,err := cacheClient.QueryWithLogicalExpire(utils.CACHE_SHOP_KEY,id,model.GetShopById,utils.CACHE_SHOP_TTL)
	if err!= nil{
		return result.Fail(fmt.Sprintf("查询商铺信息失败：%v",err))
	}


	return result.OkWithData(shop)
}

//互斥锁解决缓存击穿
// func QueryShopWithMutex(id int) (model.Shop,error){
// 	//现在redis里面查询，如果存在直接返回
// 	key := utils.CACHE_SHOP_KEY+strconv.Itoa(id)
// 	cacheShop,err := redis.RDB.Get(redis.RCtx,key).Result()
// 	// fmt.Println(cacheShop)
// 	if err == nil && cacheShop != ""{
// 		var shop model.Shop
// 		err := json.Unmarshal([]byte(cacheShop),&shop)
// 		if err != nil{
// 			return shop,err
// 		}
// 		return shop,nil
// 	}
// 	//查到一个空值(防止缓存穿透)
// 	if(err != goredis.Nil && cacheShop == ""){
// 		return model.Shop{},errors.New("商铺信息不存在！")
// 	}

// 	//缓存未命中，需要实现缓存重建
// 	//1.1获取锁
// 	mutexKey := utils.CACHE_SHOP_MUTEX_KEY+strconv.Itoa(id)
// 	ok := tryLock(mutexKey)

// 	//1.2判断是否成功
// 	for (!ok){
// 		//1.3失败，休眠并重试
// 		time.Sleep(50*time.Millisecond)
// 		ok = tryLock(mutexKey)
// 	}
// 	//1.5 返回锁
// 	defer	unLock(mutexKey)
// 	//doublecheck
// 	cacheShop,err = redis.RDB.Get(redis.RCtx,key).Result()
// 	if err!= goredis.Nil && cacheShop!= ""{
// 		var shop model.Shop
// 		err := json.Unmarshal([]byte(cacheShop),&shop)
// 		if err!= nil{
// 			return shop,err
// 		}
// 		return shop,nil
// 	}
// 	if(err != goredis.Nil && cacheShop == ""){
// 		return model.Shop{},errors.New("商铺信息不存在！")
// 	}

// 	//1.4成功，根据id查询数据库
// 	//redis中没有，在mysql中查询
// 	shop,err := model.GetShopById(uint64(id))
// 	if err != nil{
// 		//将空值存入redis
// 		if err == gorm.ErrRecordNotFound{
// 			redis.RDB.Set(redis.RCtx,key,"",utils.CACHE_NULL_TTL)
// 		}
// 		return shop,err
// 	}

// 	//写回redis
// 	shopJson,err := json.Marshal(shop)
// 	if err!= nil{
// 		return shop,err
// 	}
// 	redis.RDB.Set(redis.RCtx,key,shopJson,utils.CACHE_SHOP_TTL)

// 	return shop,nil

// }

//逻辑时间解决缓击穿（已经提前预热）
// func QueryShopWithLogicalExpire(id uint64) (model.Shop,error){
// 	key := utils.CACHE_SHOP_KEY+strconv.FormatUint(id,10)
// 	cacheShop,err := redis.RDB.Get(redis.RCtx,key).Result()
// 	if (err == goredis.Nil || cacheShop == ""){
// 		return model.Shop{},errors.New("商铺信息不存在！")
// 	}

// 	//从缓存中取出shop和过期时间
// 	var redisData redisData.RedisData
// 	err = json.Unmarshal([]byte(cacheShop),&redisData)
// 	if err!= nil{
// 		return model.Shop{},err
// 	}
// 	shopJson,err := json.Marshal(redisData.Data)
// 	if err!= nil{
// 		return model.Shop{},err
// 	}
// 	var shop model.Shop
// 	err = json.Unmarshal(shopJson,&shop)
// 	if err!= nil{
// 		return model.Shop{},err
// 	}

// 	expireTime := redisData.ExpireTime
// 	if time.Now().Unix() <= expireTime{
// 		return shop,nil
// 	}
// 	//逻辑过期，需要重建
// 	//1.获取锁
// 	mutexKey := utils.CACHE_SHOP_MUTEX_KEY+strconv.FormatUint(id,10)
// 	ok := tryLock(mutexKey)
// 	if ok {
// 		//doublecheck
// 		cacheShop,_ = redis.RDB.Get(redis.RCtx,key).Result()
// 		json.Unmarshal([]byte(cacheShop),&redisData)
// 		shopJson,_ = json.Marshal(redisData.Data)
// 		json.Unmarshal(shopJson,&shop)
// 		expireTime = redisData.ExpireTime
// 		if time.Now().Unix() <= expireTime{
// 			return shop,nil
// 		}

// 		//用一个线程去缓存重建
// 		go func(){
// 			SaveShop2Redis(uint64(id),20)
// 			unLock(mutexKey)
// 		}()
// 	}
// 	return shop,nil
// }

func UpdateShop(shop model.Shop)*dto.Result{
	result := &dto.Result{}
	id := shop.ID
	if(id == 0){
		return result.Fail("缺少商品id")
	}
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Save(&shop).Error
		if err != nil{
			return err
		}
		//删除缓存
		_,err = redis.RDB.Del(redis.RCtx,utils.CACHE_SHOP_KEY+strconv.FormatUint(id,10)).Result()
		if(err != nil){
			return err
		}
		return nil
	})
	if(err != nil){
		return result.Fail("事务执行错误")
	}
	//更新数据库
	return result.Ok()
}	

func SaveShop2Redis(id uint64,expireTime int64){
	shop,err := model.GetShopById(id)
	if err!= nil{
		return
	}
	//设置逻辑过期时间
	cacheClient.SetWithLogicalExpire(utils.CACHE_SHOP_KEY+fmt.Sprintf("%v",id),shop,10*time.Second)
}