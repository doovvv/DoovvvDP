package router

import (
	"doovvvDP/dto"
	"doovvvDP/handler"
	"doovvvDP/middleware"
	"encoding/gob"

	"github.com/gin-gonic/gin"
)

func init(){
	gob.Register(dto.UserVo{})
}
func RouterInit(){
	r := gin.Default()
	//已弃用session
    // 创建一个简单的 CookieStore, 用于存储 session
    // store := cookie.NewStore([]byte("secret"))  // 用于加密和签名的密钥

    // 使用 session 中间件
    // r.Use(sessions.Sessions("mysession", store))

	r.Use(middleware.RefreshToken())
	userRouter := r.Group("/user")
	{
		userRouter.GET("/code",handler.SendCode)
		userRouter.POST("/login",handler.Login)
	}
	authUserRouter := r.Group("/user")
	authUserRouter.Use(middleware.LoginInterceptor())
	{
		authUserRouter.GET("/me",handler.Me)
	}
	r.Run(":8081")
}