package v1

import (
	"doovvvDP/dal/model"
	"doovvvDP/dal/redis"
	"doovvvDP/dto"
	"doovvvDP/utils"
	"encoding/json"

	goredis "github.com/redis/go-redis/v9"
)
func QueryShopTypeList() (result *dto.Result){
	result = &dto.Result{}

	key := utils.CACHE_SHOP_TYPE_KEY
	cacheShoptypes,err := redis.RDB.Get(redis.RCtx,key).Result()
	if err != goredis.Nil{
		var shopTypes []model.ShopType
		err := json.Unmarshal([]byte(cacheShoptypes),&shopTypes)
		if err != nil{
			return result.Fail("商品类型解析失败")
		}
		return result.OkWithData(shopTypes)
	}
	shopTypes,err := model.GetShopTypeList()
	if err != nil{
		result.Fail("商品类型查询失败")
	}
	shopTypesJson,err := json.Marshal(shopTypes)
	if err!= nil{
		return result.Fail("商品类型解析失败")
	}
	err = redis.RDB.Set(redis.RCtx,key,shopTypesJson,utils.CACHE_SHOP_TYPE_TTL).Err()
	if err!= nil{
		return result.Fail("商品类型缓存失败")
	}
	return result.OkWithData(shopTypes)
}