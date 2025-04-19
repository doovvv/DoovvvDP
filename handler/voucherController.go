package handler

import (
	"strconv"

	v1 "doovvvDP/api/v1"
	"doovvvDP/dal/model"
	"doovvvDP/dto"

	"github.com/gin-gonic/gin"
)

func VoucherHandlerInit() {
	v1.VoucherServiceInit()
}

func QueryVoucherByShopId(c *gin.Context) {
	result := &dto.Result{}
	shopIdStr, ok := (c.Params.Get("shopId"))
	if !ok {
		c.JSON(400, result.Fail("shopId is required"))
		return
	}
	shopId, err := strconv.Atoi(shopIdStr)
	if err != nil {
		c.JSON(400, result.Fail("shopId is invalid"))
		return
	}
	result = v1.QueryVoucherByShopId(shopId)
	c.JSON(200, result)
}

func AddSeckillVoucher(c *gin.Context) {
	result := &dto.Result{}
	var voucher model.DTOVoucher
	err := c.ShouldBind(&voucher)
	if err != nil {
		// body, _ := c.GetRawData()
		c.JSON(400, result.Fail("body err"))
		return
	}
	result = v1.AddSeckillVoucher(voucher)
	c.JSON(200, result)
}

func SeckillVoucher(c *gin.Context) {
	result := &dto.Result{}
	voucherIdStr := c.Param("voucherId")
	voucherId, err := strconv.ParseUint(voucherIdStr, 10, 64)
	if err != nil {
		result.Fail("代金券ID无效")
	}
	user, ok := c.Get("user")
	if ok {
		if u, ok := user.(map[string]string); ok {
			userId, err := strconv.ParseUint(u["id"], 10, 64)
			if err != nil {
				result.Fail("用户ID无效")
			}
			result = v1.SeckillVoucher(voucherId, userId)
			// fmt.Println(result)
		} else {
			result.Fail("用户信息无效")
		}
	} else {
		result.Fail("用户未登录")
	}
	c.JSON(200, result)
}
