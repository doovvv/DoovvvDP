package redisData
type RedisData struct{
	Data any `json:"data"`
	ExpireTime int64 `json:"expire_time"`
}