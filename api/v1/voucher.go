package v1

import (
	"time"

	"doovvvDP/dal/model"
	"doovvvDP/dal/mysql"
	"doovvvDP/dto"
	"doovvvDP/utils"
)

func QueryVoucherByShopId(shopId int) *dto.Result {
	result := &dto.Result{}
	vouchers := model.QueryVoucherByShopId(shopId)
	return result.OkWithData(vouchers)
}

func AddSeckillVoucher(voucher model.DTOVoucher) *dto.Result {
	result := &dto.Result{}
	model.AddSeckillVoucher(voucher)
	return result.Ok()
}

func SeckillVoucher(voucherId uint64, userId uint64) *dto.Result {
	result := &dto.Result{}
	voucher, err := model.QueryVoucherById(voucherId)
	if err != nil {
		return result.Fail("查询失败")
	}
	// 判断秒杀是否开始
	if time.Now().Before(voucher.BeginTime) {
		return result.Fail("秒杀还未开始")
	}

	if time.Now().After(voucher.EndTime) {
		return result.Fail("秒杀已经结束")
	}

	if voucher.Stock < 1 {
		return result.Fail("秒杀券库存不足")
	}
	tx := mysql.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			result.Fail("秒杀失败")
		}
	}()
	// 删除库存
	err = model.DecreaseStock(tx, voucherId)
	if err != nil {
		tx.Rollback()
		return result.Fail("库存不足")
	}
	// 订单id
	id, err := utils.MyIdWorker.Generate()
	if err != nil {
		tx.Rollback()
		return result.Fail("订单id生成错误")
	}
	voucherOrder := model.VoucherOrder{
		ID:        id,
		UserID:    userId,
		VoucherID: voucherId,
	}

	err = model.AddVoucherOrder(tx, voucherOrder)
	if err != nil {
		tx.Rollback()
		return result.Fail("订单创建失败")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return result.Fail("订单创建失败")
	}
	// fmt.Println(voucherOrder.ID)
	return result.OkWithData(voucherOrder.ID)
}
