package router

import (
	"doovvvDP/handler"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)
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
	r.Run(":8081")
}