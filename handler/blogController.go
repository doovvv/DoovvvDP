package handler

import (
	"path/filepath"
	"strconv"

	v1 "doovvvDP/api/v1"
	"doovvvDP/dal/model"
	"doovvvDP/dto"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	result := &dto.Result{}
	file, err := c.FormFile("file")
	if err != nil {
		// fmt.Println(err.Error())
		result.Fail(err.Error())
		c.JSON(200, result)
		return
	}
	fileName := file.Filename
	dst := filepath.Join("resources/nginx-1.18.0/html/hmdp/imgs", fileName)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		result.Fail(err.Error())
		c.JSON(200, result)
		return
	}
	result.OkWithData("/" + fileName)
	c.JSON(200, result)
}

func CreateBlog(c *gin.Context) {
	result := &dto.Result{}
	blog := model.Blog{}
	c.ShouldBind(&blog)
	user, ok := c.Get("user")
	if !ok {
		result.Fail("用户未登录")
		c.JSON(200, result)
		return
	}
	userMap := user.(map[string]string)
	userId, err := strconv.ParseUint(userMap["id"], 10, 64)
	if err != nil {
		result.Fail("用户ID无效")
		c.JSON(200, result)
		return
	}
	result = v1.CreateBlog(blog, userId)
	c.JSON(200, result)
}
