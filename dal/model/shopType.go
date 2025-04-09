package model

import (
	"doovvvDP/dal/mysql"
	"time"
)

type ShopType struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	Name       string    `gorm:"type:varchar(32);comment:类型名称" json:"name"`
	Icon       string    `gorm:"type:varchar(255);comment:图标" json:"icon"`
	Sort       uint      `gorm:"type:int unsigned;comment:顺序" json:"sort"`
	CreateTime time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createTime"`
	UpdateTime time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updateTime"`
}

// TableName 自定义表名（GORM 默认使用结构体名的蛇形复数作为表名）
func (ShopType) TableName() string {
	return "tb_shop_type"
}
func GetShopTypeList() ([]ShopType, error) {
	var shopTypes []ShopType
	err := mysql.DB.Order("sort asc").Find(&shopTypes).Error
	return shopTypes, err
}