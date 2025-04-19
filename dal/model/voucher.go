package model

import (
	"fmt"
	"time"

	"doovvvDP/dal/mysql"
	"doovvvDP/dal/redis"

	"gorm.io/gorm"
)

// TbVoucher 代金券表
type TbVoucher struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                                              // 主键
	ShopID      *uint64   `gorm:"column:shop_id" json:"shopId"`                                                                              // 商铺id
	Title       string    `gorm:"column:title;type:varchar(255);not null" json:"title"`                                                      // 代金券标题
	SubTitle    *string   `gorm:"column:sub_title;type:varchar(255)" json:"subTitle"`                                                        // 副标题
	Rules       *string   `gorm:"column:rules;type:varchar(1024)" json:"rules"`                                                              // 使用规则
	PayValue    uint64    `gorm:"column:pay_value;not null" json:"payValue"`                                                                 // 支付金额，单位是分。例如200代表2元
	ActualValue int64     `gorm:"column:actual_value;not null" json:"actualValue"`                                                           // 抵扣金额，单位是分。例如200代表2元
	Type        uint8     `gorm:"column:type;not null;default:0" json:"type"`                                                                // 0,普通券；1,秒杀券
	Status      uint8     `gorm:"column:status;not null;default:1" json:"status"`                                                            // 1,上架; 2,下架; 3,过期
	CreateTime  time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP" json:"createTime"`                             // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updateTime"` // 更新时间
}

type TbSeckillVoucher struct {
	Id        uint64    `gorm:"column:id" json:"id"`
	Stock     int       `gorm:"column:stock" json:"stock"`
	BeginTime time.Time `gorm:"column:begin_time" json:"beginTime"`
	EndTime   time.Time `gorm:"column:end_time" json:"endTime"`
}

type DTOVoucher struct {
	Id          uint64    `gorm:"column:id" json:"id"`
	ShopID      *uint64   `gorm:"columb:shop_id" json:"shopId"`                         // 商铺id
	Title       string    `gorm:"column:title;type:varchar(255);not null" json:"title"` // 代金券标题
	SubTitle    *string   `gorm:"column:sub_title;type:varchar(255)" json:"subTitle"`   // 副标题
	Rules       *string   `gorm:"column:rules;type:varchar(1024)" json:"rules"`         // 使用规则
	PayValue    uint64    `gorm:"column:pay_value;not null" json:"payValue"`            // 支付金额，单位是分。例如200代表2元
	ActualValue int64     `gorm:"column:actual_value;not null" json:"actualValue"`      // 抵扣金额，单位是分。例如200代表2元
	Type        uint8     `gorm:"column:type;not null;default:0" json:"type"`           // 0,普通券；1,秒杀券
	Status      uint8     `gorm:"column:status;not null;default:1" json:"status"`       // 1,上架; 2,下架; 3,过期
	Stock       int       `gorm:"column:stock" json:"stock"`
	BeginTime   time.Time `gorm:"column:begin_time" json:"beginTime"`
	EndTime     time.Time `gorm:"column:end_time" json:"endTime"`
}

// TableName 设置表名
func (TbVoucher) TableName() string {
	return "tb_voucher"
}

func (TbSeckillVoucher) TableName() string {
	return "tb_seckill_voucher"
}

func QueryVoucherByShopId(shopId int) []DTOVoucher {
	vouchers := []DTOVoucher{}
	// 使用LEFT JOIN查询
	err := mysql.DB.
		Table("tb_voucher v").
		Select("v.*, s.stock as Stock, s.begin_time as begin_time, s.end_time as end_time").
		Joins("LEFT JOIN tb_seckill_voucher s ON v.id = s.id AND v.type = 1").
		Where("v.shop_id = ?", shopId).
		Scan(&vouchers).Error
	// fmt.Println(vouchers)
	// err := mysql.DB.Where("shop_id = ?",shopId).Find(&vouchers).Error
	if err != nil {
		return nil
	}

	return vouchers
}

func QueryVoucherById(id uint64) (DTOVoucher, error) {
	voucher := DTOVoucher{}

	err := mysql.DB.Table("tb_voucher v").
		Select("v.*, s.stock as Stock, s.begin_time as begin_time, s.end_time as end_time").
		Joins("LEFT JOIN tb_seckill_voucher s ON v.id = s.id AND v.type = 1").
		Where("v.id =?", id).
		Scan(&voucher).Error
	if err != nil {
		return DTOVoucher{}, err
	}
	return voucher, nil
}

func AddSeckillVoucher(dtoVoucher DTOVoucher) error {
	// 开始事务
	voucher := TbVoucher{
		ShopID:      dtoVoucher.ShopID,
		Title:       dtoVoucher.Title,
		SubTitle:    dtoVoucher.SubTitle,
		Rules:       dtoVoucher.Rules,
		PayValue:    dtoVoucher.PayValue,
		ActualValue: dtoVoucher.ActualValue,
		Type:        dtoVoucher.Type,
		Status:      dtoVoucher.Status,
	}
	tx := mysql.DB.Begin()
	// 先保存 TbVoucher 表数据
	if err := tx.Save(&voucher).Error; err != nil {
		tx.Rollback() // 如果 TbVoucher 保存失败，回滚事务
		return fmt.Errorf("failed to save voucher: %w", err)
	}

	// 只有当类型是秒杀券时，才保存 TbSeckillVoucher 表数据
	if voucher.Type == 1 {
		// 这里保存 TbSeckillVoucher 表数据
		seckillVoucher := TbSeckillVoucher{
			Id:        voucher.ID,
			Stock:     dtoVoucher.Stock,
			BeginTime: dtoVoucher.BeginTime,
			EndTime:   dtoVoucher.EndTime,
		}
		if err := tx.Save(&seckillVoucher).Error; err != nil {
			tx.Rollback() // 如果 TbSeckillVoucher 保存失败，回滚事务
			return fmt.Errorf("failed to save seckill voucher: %w", err)
		}
	}
	// 增加优惠券的同时保存库存到redis
	err := redis.RDB.Set(redis.RCtx, fmt.Sprintf("seckill:stock:%d", voucher.ID), dtoVoucher.Stock, -1).Err()
	if err != nil {
		tx.Rollback() // 如果 TbSeckillVoucher 保存失败，回滚事务
		return fmt.Errorf("failed to save seckill voucher: %w", err)
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// 库存减一
func DecreaseStock(db *gorm.DB, voucherId uint64) error {
	// 乐观锁（添加version字段）
	err := db.Model(&TbSeckillVoucher{}).
		Where("id = ? AND stock > 0", voucherId).
		Update("stock", gorm.Expr("stock - 1")).Error
	return err
}
