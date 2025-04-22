package model

import (
	"time"

	"doovvvDP/dal/mysql"
	"doovvvDP/utils"

	"gorm.io/gorm"
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
	return "tb_blog"
}

func CreateBlog(blog *Blog) error {
	err := mysql.DB.Create(blog).Error
	return err
}

func GetBlogById(id uint64) (Blog, error) {
	var blog Blog
	err := mysql.DB.Where("id = ?", id).First(&blog).Error
	return blog, err
}

func QueryHotBlogs(current int) ([]Blog, error) {
	var blogs []Blog
	err := mysql.DB.Order("liked desc").
		Limit(10).Offset(utils.MAX_PAGE_SIZE * (current - 1)).Find(&blogs).Error
	return blogs, err
}

func LikeBlog(blogId uint64) error {
	err := mysql.DB.Model(&Blog{}).Where("id = ?", blogId).Update("liked", gorm.Expr("liked + ?", 1)).Error
	return err
}

func UnLikeBlog(blogId uint64) error {
	err := mysql.DB.Model(&Blog{}).Where("id =?", blogId).Update("liked", gorm.Expr("liked -?", 1)).Error
	return err
}

func QueryBlogOfUser(userId uint64, current int) ([]Blog, error) {
	var blogs []Blog
	err := mysql.DB.Where("user_id =?", userId).
		Limit(10).Offset(utils.MAX_PAGE_SIZE * (current - 1)).Find(&blogs).Error
	return blogs, err
}

func QueryBlogsByIDs(ids []uint64) ([]Blog, error) {
	var blogs []Blog
	// 这里没有调整blog的顺序
	err := mysql.DB.Where("id IN ?", ids).Find(&blogs).Error
	return blogs, err
}
