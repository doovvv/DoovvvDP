package middleware

import (
	"doovvvDP/dal/redis"
	"doovvvDP/utils"

	"github.com/gin-gonic/gin"
)
func RefreshToken() gin.HandlerFunc{
	return func(c *gin.Context){
		//获取header中的token
		token := c.GetHeader("Authorization")
		if token == ""{
			//由loginInterceptor拦截器拦截
			c.Next()
			return
		}
		userMap := redis.RDB.HGetAll(redis.RCtx,utils.LOGIN_TOKEN_KEY+token).Val()
		if len(userMap) == 0{
			//由loginInterceptor拦截器拦截
			c.Next()
			return
		}
		c.Set("user",userMap)

		//刷新token有效期
		redis.RDB.Expire(redis.RCtx,utils.LOGIN_TOKEN_KEY+token,utils.LOGIN_TOKEN_TTL)
		c.Next()

	}
}