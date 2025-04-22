package router

import (
	"encoding/gob"

	"doovvvDP/dto"
	"doovvvDP/handler"
	"doovvvDP/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	gob.Register(dto.UserVo{})
}

func RouterInit() {
	handler.VoucherHandlerInit()
	r := gin.Default()
	// 已弃用session
	// 创建一个简单的 CookieStore, 用于存储 session
	// store := cookie.NewStore([]byte("secret"))  // 用于加密和签名的密钥

	// 使用 session 中间件
	// r.Use(sessions.Sessions("mysession", store))

	r.Use(middleware.RefreshToken())
	authRouter := r.Group("")
	authRouter.Use(middleware.LoginInterceptor())
	userRouter := r.Group("/user")
	{
		userRouter.GET("/code", handler.SendCode)
		userRouter.POST("/login", handler.Login)
	}
	authUserRouter := authRouter.Group("/user")
	// authUserRouter.Use(middleware.LoginInterceptor())
	{
		authUserRouter.GET("/:id", handler.QueryUserById)
		authUserRouter.GET("/info/:id", handler.QueryUserById)
		authUserRouter.GET("/me", handler.Me)
	}
	shopRouter := r.Group("/shop")
	{
		shopRouter.GET("/:id", handler.QueryShopById)
		shopRouter.PUT("", handler.UpdateShop)
	}
	shopTypeRouter := r.Group("/shop-type")
	{
		shopTypeRouter.GET("/list", handler.QueryShopTypeList)
	}
	voucherRouter := r.Group("/voucher")
	{
		voucherRouter.GET("/list/:shopId", handler.QueryVoucherByShopId)
		voucherRouter.POST("/seckill", handler.AddSeckillVoucher)
	}
	voucherOrderRouter := authRouter.Group("/voucher-order")
	{
		voucherOrderRouter.POST("/seckill/:voucherId", handler.SeckillVoucher)
	}
	uploadRouter := r.Group("/upload")
	{
		uploadRouter.POST("/blog", handler.UploadFile)
	}
	blogRouter := r.Group("/blog")
	{
		blogRouter.POST("", handler.CreateBlog)
		blogRouter.GET("/hot", handler.QueryHotBlogs)
		blogRouter.GET("/:id", handler.QueryBlogById)
		blogRouter.PUT("/like/:id", handler.LikeBlog)
		blogRouter.GET("/likes/:id", handler.QueryBlogLikes)
		blogRouter.GET("/of/user", handler.QueryBlogByUserId)
	}
	blogAuthRouter := authRouter.Group("/blog")
	{
		blogAuthRouter.GET("of/follow", handler.QueryBlogOfFollow)
	}
	followRouter := authRouter.Group("/follow")
	{
		followRouter.PUT("/:id/:isFollow", handler.FollowUser)
		followRouter.GET("/or/not/:id", handler.QueryFollowByUserId)
		followRouter.GET("/common/:id", handler.FollowCommons)
	}
	r.Run(":8081")
}
