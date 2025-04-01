package handler

import (
	v1 "doovvvDP/api/v1"
	"doovvvDP/dto"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)
func SendCode(c *gin.Context){
	// 发送短信验证码并保存验证码
	session := sessions.Default(c) // 获取当前请求的 session
	phone := c.Query("phone")
	result := v1.SendCode(phone,session)
	c.JSON(200, result)
}
func Login(c *gin.Context){
	var result *dto.Result = &dto.Result{
		
	}
	session := sessions.Default(c) // 获取当前请求的 session
	var userdto dto.UserDTO
	if err := c.ShouldBindJSON(&userdto); err != nil {
		result.Fail("参数错误")
	}
	result = v1.Login(userdto,session)
	c.JSON(200, result)

}
func Me(c *gin.Context){
	result := &dto.Result{}
	user,ok := c.Get("user")
	if !ok{
		result.Fail("未找到user信息")
		c.JSON(200,result)
		return
	}
	// fmt.Println(user)
	result.OkWithData(user)
	c.JSON(200,result)
}