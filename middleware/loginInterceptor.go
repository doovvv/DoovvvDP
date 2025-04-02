package middleware

import (
	"github.com/gin-gonic/gin"
)
func LoginInterceptor() gin.HandlerFunc{
	return func(c *gin.Context){
		//判断context是否有用户信息，有代表已经登录
		if _,ok := c.Get("user");!ok{
			c.JSON(401,gin.H{
			})
			c.Abort()
			return
		}
		c.Next()

	}
}