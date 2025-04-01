package router

import (
	"doovvvDP/dto"
	"doovvvDP/handler"
	"doovvvDP/middleware"
	"encoding/gob"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func init(){
	gob.Register(dto.UserVo{})
}
func RouterInit(){
	r := gin.Default()
    // 创建一个简单的 CookieStore, 用于存储 session
    store := cookie.NewStore([]byte("secret"))  // 用于加密和签名的密钥

    // 使用 session 中间件
    r.Use(sessions.Sessions("mysession", store))
	userRouter := r.Group("/user")
	{
		userRouter.GET("/code",handler.SendCode)
		userRouter.POST("/login",handler.Login)
	}
	authUserRouter := r.Group("/user")
	authUserRouter.Use(middleware.PreHandle())
	{
		authUserRouter.GET("/me",handler.Me)
	}
	r.Run(":8081")
}