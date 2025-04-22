package handler

import (
	"strconv"

	v1 "doovvvDP/api/v1"
	"doovvvDP/dto"

	"github.com/gin-gonic/gin"
)

func FollowUser(c *gin.Context) {
	result := &dto.Result{}
	userId, err := getUserId(c)
	if err != nil {
		result.Fail(err.Error())
		c.JSON(200, result)
		return
	}
	followId := c.Param("id")
	followIdInt, err := strconv.ParseUint(followId, 10, 64)
	if err != nil {
		result.Fail("ID无效")
		c.JSON(200, result)
		return
	}
	isFollow := c.Param("isFollow")
	isFollowBool, err := strconv.ParseBool(isFollow)
	if err != nil {
		result.Fail("关注状态无效")
		c.JSON(200, result)
		return
	}
	result = v1.FollowUser(userId, followIdInt, isFollowBool)
	c.JSON(200, result)
}

func QueryFollowByUserId(c *gin.Context) {
	result := &dto.Result{}
	userId, err := getUserId(c)
	if err != nil {
		result.Fail(err.Error())
		c.JSON(200, result)
		return
	}
	followId := c.Param("id")
	followIdInt, err := strconv.ParseUint(followId, 10, 64)
	if err != nil {
		result.Fail("ID无效")
		c.JSON(200, result)
		return
	}
	result = v1.QueryFollowByUserId(userId, followIdInt)
	c.JSON(200, result)
}

func FollowCommons(c *gin.Context) {
	result := &dto.Result{}
	otherId := c.Param("id")
	otherIdInt, err := strconv.ParseUint(otherId, 10, 64)
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

	result = v1.FollowCommons(userId, otherIdInt)
	c.JSON(200, result)
}
