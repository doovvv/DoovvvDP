package handler

import (
	"strconv"

	v1 "doovvvDP/api/v1"
	"doovvvDP/dto"

	"github.com/gin-gonic/gin"
)

func SendCode(c *gin.Context) {
	// 发送短信验证码并保存验证码
	// session := sessions.Default(c) // 获取当前请求的 session
	phone := c.Query("phone")
	result := v1.SendCode(phone)
	c.JSON(200, result)
}

func Login(c *gin.Context) {
	var result *dto.Result = &dto.Result{}
	// session := sessions.Default(c) // 获取当前请求的 session
	var userdto dto.UserDTO
	if err := c.ShouldBindJSON(&userdto); err != nil {
		result.Fail("参数错误")
	}
	result = v1.Login(userdto)
	c.JSON(200, result)
}

func Me(c *gin.Context) {
	// 从gin.context中获取用户信息，目前类型是map
	result := &dto.Result{}
	user, ok := c.Get("user")
	if !ok {
		result.Fail("未找到user信息")
		c.JSON(200, result)
		return
	}
	// fmt.Println(user)

	// 给前端返回用户信息
	result.OkWithData(user)
	c.JSON(200, result)
}

func getUserId(c *gin.Context) (uint64, error) {
	user, ok := c.Get("user")
	if !ok {
		return 0, nil
	}
	userMap := user.(map[string]string)
	userId, err := strconv.ParseUint(userMap["id"], 10, 64)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func QueryUserById(c *gin.Context) {
	result := &dto.Result{}
	id := c.Param("id")
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		result.Fail("ID无效")
		c.JSON(200, result)
		return
	}
	result = v1.QueryUserById(idInt)
	c.JSON(200, result)
}

func Sign(c *gin.Context) {
	result := &dto.Result{}
	userId, err := getUserId(c)
	if err != nil {
		result.Fail("未找到user信息")
		c.JSON(200, result)
		return
	}
	result = v1.Sign(userId)
	c.JSON(200, result)
}

func SignCount(c *gin.Context) {
	result := &dto.Result{}
	userId, err := getUserId(c)
	if err != nil {
		result.Fail("未找到user信息")
		c.JSON(200, result)
		return
	}
	result = v1.SignCount(userId)
	c.JSON(200, result)
}
