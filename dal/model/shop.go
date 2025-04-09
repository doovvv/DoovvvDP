package model

import (
	"doovvvDP/dal/mysql"
	"time"
)

type Shop struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	Name       string    `gorm:"type:varchar(128);not null;comment:商铺名称" json:"name"`
	TypeID     uint64    `gorm:"column:type_id;type:bigint unsigned;not null;index;comment:商铺类型的id" json:"typeId"`
	Images     string    `gorm:"type:varchar(1024);not null;comment:商铺图片，多个图片以','隔开" json:"images"`
	Area       string    `gorm:"type:varchar(128);comment:商圈，例如陆家嘴" json:"area"`
	Address    string    `gorm:"type:varchar(255);not null;comment:地址" json:"address"`
	X          float64   `gorm:"type:double unsigned;not null;comment:经度" json:"x"`
	Y          float64   `gorm:"type:double unsigned;not null;comment:纬度" json:"y"`
	AvgPrice   *uint64   `gorm:"type:bigint unsigned;comment:均价，取整数" json:"avgPrice"` // 使用指针允许NULL
	Sold       uint32    `gorm:"type:int unsigned zerofill;not null;comment:销量" json:"sold"`
	Comments   uint32    `gorm:"type:int unsigned zerofill;not null;comment:评论数量" json:"comments"`
	Score      uint8     `gorm:"type:tinyint unsigned zerofill;not null;comment:评分，1~5分，乘10保存，避免小数" json:"score"`
	OpenHours  string    `gorm:"type:varchar(32);comment:营业时间，例如 10:00-22:00" json:"openHours"`
	CreateTime time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createTime"`
	UpdateTime time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updateTime"`
}

// TableName 自定义表名
func (Shop) TableName() string {
	return "tb_shop"
}
func GetShopById(id uint64) (Shop, error) {
	var shop Shop
	err := mysql.DB.Where("id = ?", id).First(&shop).Error
	//模拟重建时间
	time.Sleep(200*time.Microsecond)
	return shop, err
}
func UpdateShopById(shop Shop)(error){
	err := mysql.DB.Save(&shop).Error
	return err
}