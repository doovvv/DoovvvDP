package middleware

import (
	"doovvvDP/dto"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)
func PreHandle() gin.HandlerFunc{
	return func(c *gin.Context){
		session := sessions.Default(c)
		userDto,ok := session.Get("user").(dto.UserVo)
		if !ok{
			//不存在就进行拦截
			c.JSON(401,gin.H{ });
			c.Abort()
			return
		}
		c.Set("user",userDto)
		c.Next()

	}
}