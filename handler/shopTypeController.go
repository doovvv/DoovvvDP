package handler

import (
	v1 "doovvvDP/api/v1"

	"github.com/gin-gonic/gin"
)
func QueryShopTypeList(c *gin.Context){
	result := v1.QueryShopTypeList()
	c.JSON(200,result)
}