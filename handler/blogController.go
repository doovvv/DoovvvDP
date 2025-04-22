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

func QueryBlogById(c *gin.Context) {
	result := &dto.Result{}
	id := c.Param("id")
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		result.Fail("ID无效")
		c.JSON(200, result)
		return
	}
	userId, err := getUserId(c)
	if err != nil {
		result.Fail(err.Error())
		c.JSON(200, result)
		return
	}
	result = v1.QueryBlogByID(idInt, userId)
	c.JSON(200, result)
}

func QueryHotBlogs(c *gin.Context) {
	result := &dto.Result{}
	current := c.Query("current")
	currentInt, err := strconv.Atoi(current)
	if err != nil {
		result.Fail("当前页无效")
		c.JSON(200, result)
		return
	}
	userId, err := getUserId(c)
	if err != nil {
		result.Fail(err.Error())
		c.JSON(200, result)
		return
	}
	result = v1.QueryHotBlogs(currentInt, userId)
	c.JSON(200, result)
}

func LikeBlog(c *gin.Context) {
	result := &dto.Result{}
	id := c.Param("id")
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		result.Fail("ID无效")
		c.JSON(200, result)
		return
	}
	userId, err := getUserId(c)
	if err != nil {
		result.Fail(err.Error())
		c.JSON(200, result)
		return
	}
	result = v1.LikeBlog(idInt, userId)
	c.JSON(200, result)
}

func QueryBlogLikes(c *gin.Context) {
	result := &dto.Result{}
	id := c.Param("id")
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		result.Fail("ID无效")
		c.JSON(200, result)
		return
	}
	result = v1.QueryBlogLikes(idInt)
	c.JSON(200, result)
}

func QueryBlogByUserId(c *gin.Context) {
	result := &dto.Result{}
	lastId := c.Query("id")
	lastIdInt, err := strconv.ParseUint(lastId, 10, 64)
	if err != nil {
		result.Fail("ID无效")
		c.JSON(200, result)
		return
	}
	current := c.Query("current")
	currentInt, err := strconv.Atoi(current)
	if err != nil {
		result.Fail("当前页无效")
		c.JSON(200, result)
		return
	}
	result = v1.QueryBlogByUserId(lastIdInt, currentInt)
	c.JSON(200, result)
}

func QueryBlogOfFollow(c *gin.Context) {
	result := &dto.Result{}
	maxIdStr := c.Query("lastId")
	maxId, err := strconv.ParseUint(maxIdStr, 10, 64)
	if err != nil {
		result.Fail("ID无效")
		c.JSON(200, result)
		return
	}

	offset := c.DefaultQuery("offset", "0")
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		result.Fail("偏移量无效")
		c.JSON(200, result)
		return
	}
	userId, err := getUserId(c)
	if err != nil {
		result.Fail(err.Error())
		c.JSON(200, result)
		return
	}
	result = v1.QueryBlogOfFollow(userId, maxId, offsetInt)
	c.JSON(200, result)
}
