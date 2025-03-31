package model

import (
	"doovvvDP/dal/mysql"
	"doovvvDP/utils"
	"time"

	"gorm.io/gorm"
)
type TbUser struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;comment:'主键'" json:"id"`
	Phone     string    `gorm:"type:varchar(11);not null;uniqueIndex:uniqe_key_phone;comment:'手机号码'" json:"phone"`
	Password  string    `gorm:"type:varchar(128);default:'';comment:'密码，加密存储'" json:"password"`
	NickName  string    `gorm:"type:varchar(32);default:'';comment:'昵称，默认是用户id'" json:"nick_name"`
	Icon      string    `gorm:"type:varchar(255);default:'';comment:'人物头像'" json:"icon"`
	CreateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
}
func (TbUser) TableName() string {
	return "tb_user"
}
func CheckUserExist(phone string) bool {
	var user TbUser
	err := mysql.DB.Where("phone = ?", phone).First(&user).Error
	return err == gorm.ErrRecordNotFound
}
func CreateUserWithPhone(phone string) error {
	//todo
	user := TbUser{
		Phone: phone,
		NickName: utils.USER_NICK_NAME_PREFIX+utils.RandomString(10),
	}
	err := mysql.DB.Create(&user).Error
	return err
}