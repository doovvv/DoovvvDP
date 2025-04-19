package model

import (
	"time"

	"doovvvDP/dal/mysql"
)

type Blog struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ShopID     int64     `gorm:"not null" json:"shopId"`
	UserID     uint64    `gorm:"not null" json:"userId"`
	Title      string    `gorm:"type:varchar(255);not null" json:"title"`
	Images     string    `gorm:"type:varchar(2048);not null" json:"images"`
	Content    string    `gorm:"type:varchar(2048);not null" json:"content"`
	Liked      uint32    `gorm:"default:0" json:"liked"`
	Comments   uint32    `gorm:"default:0" json:"comments"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (Blog) TableName() string {
	return "blog"
}

func CreateBlog(blog Blog) error {
	err := mysql.DB.Create(&blog).Error
	return err
}
