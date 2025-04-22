package model

import (
	"time"

	"doovvvDP/dal/mysql"
)

type Follow struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint64    `gorm:"not null" json:"userId"`
	FollowUserID uint64    `gorm:"not null" json:"followUserId"`
	CreateTime   time.Time `gorm:"autoCreateTime" json:"createTime"`
}

func (Follow) TableName() string {
	return "tb_follow"
}

func FollowUser(follow Follow) error {
	err := mysql.DB.Create(&follow).Error
	return err
}

func UnFollowUser(follow Follow) error {
	err := mysql.DB.Where("user_id = ? AND follow_user_id = ?", follow.UserID, follow.FollowUserID).Delete(&follow).Error
	return err
}

func QueryFollowByUserId(follow Follow) (int64, error) {
	var count int64
	err := mysql.DB.Model(&follow).Where("user_id =? AND follow_user_id =?", follow.UserID, follow.FollowUserID).Count(&count).Error
	return count, err
}

func QueryFans(userId uint64) ([]Follow, error) {
	var fans []Follow
	err := mysql.DB.Where("follow_user_id =?", userId).Find(&fans).Error
	return fans, err
}
