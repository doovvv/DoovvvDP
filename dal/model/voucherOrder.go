package model

import (
	"time"

	"doovvvDP/dal/mysql"

	"gorm.io/gorm"
)

type VoucherOrder struct {
	ID         uint64     `gorm:"column:id;primaryKey;comment:主键" json:"id"`
	UserID     uint64     `gorm:"column:user_id;not null;comment:下单的用户id" json:"userId"`
	VoucherID  uint64     `gorm:"column:voucher_id;not null;comment:购买的代金券id" json:"voucherId"`
	PayType    uint8      `gorm:"column:pay_type;not null;default:1;comment:支付方式 1：余额支付；2：支付宝；3：微信" json:"payType"`
	Status     uint8      `gorm:"column:status;not null;default:1;comment:订单状态" json:"status"`
	CreateTime time.Time  `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP;comment:下单时间" json:"createTime"`
	PayTime    *time.Time `gorm:"column:pay_time;comment:支付时间" json:"payTime"`       // 允许 NULL
	UseTime    *time.Time `gorm:"column:use_time;comment:核销时间" json:"useTime"`       // 允许 NULL
	RefundTime *time.Time `gorm:"column:refund_time;comment:退款时间" json:"refundTime"` // 允许 NULL
	UpdateTime time.Time  `gorm:"column:update_time;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updateTime"`
}

func (v *VoucherOrder) TableName() string {
	return "tb_voucher_order"
}

func AddVoucherOrder(db *gorm.DB, voucherOrder *VoucherOrder) error {
	return db.Create(voucherOrder).Error
}

func CheckVoucherOrder(userId uint64, voucherId uint64) bool {
	var count int64
	err := mysql.DB.Model(&VoucherOrder{}).Where("user_id = ? AND voucher_id = ?", userId, voucherId).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}
