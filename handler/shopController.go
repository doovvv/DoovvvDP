package handler

import (
	"fmt"
	"strconv"

	v1 "doovvvDP/api/v1"
	"doovvvDP/dal/model"
	"doovvvDP/dto"

	"github.com/gin-gonic/gin"
)

func QueryShopById(c *gin.Context) {
	var result *dto.Result = &dto.Result{} // 防止参数错误时返回空指针
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		result.Fail(err.Error())
		c.JSON(200, result)
		return
	}
	result = v1.QueryShopById(uint64(id))
	c.JSON(200, result)
}

func UpdateShop(c *gin.Context) {
	var shop model.Shop
	var result *dto.Result = &dto.Result{}
	err := c.ShouldBindBodyWithJSON(&shop)
	fmt.Println(shop)
	if err != nil {
		result.Fail("数据错误")
	}
	result = v1.UpdateShop(shop)
	c.JSON(200, result)
}
